package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 配置结构
type Config struct {
	WeatherAPI struct {
		Key      string `json:"key"`
		CityCode string `json:"city_code"`
		CityName string `json:"city_name"`
	} `json:"weather_api"`
	Feishu struct {
		WebhookURL string `json:"webhook_url"`
	} `json:"feishu"`
	Schedule struct {
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
	} `json:"schedule"`
}

// 飞书消息结构
type FeishuMessage struct {
	MsgType string                 `json:"msg_type"`
	Content map[string]interface{} `json:"content"`
}

// 高德天气API响应结构
type WeatherResponse struct {
	Status    string        `json:"status"`
	Count     string        `json:"count"`
	Info      string        `json:"info"`
	InfoCode  string        `json:"infocode"`
	Lives     []WeatherInfo `json:"lives,omitempty"`
	Forecasts []Forecast    `json:"forecasts,omitempty"`
}

type WeatherInfo struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`
	Temperature   string `json:"temperature"`
	WindDirection string `json:"winddirection"`
	WindPower     string `json:"windpower"`
	Humidity      string `json:"humidity"`
	ReportTime    string `json:"reporttime"`
}

type Forecast struct {
	City       string        `json:"city"`
	Adcode     string        `json:"adcode"`
	Province   string        `json:"province"`
	ReportTime string        `json:"reporttime"`
	Casts      []WeatherCast `json:"casts"`
}

type WeatherCast struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"dayweather"`
	NightWeather string `json:"nightweather"`
	DayTemp      string `json:"daytemp"`
	NightTemp    string `json:"nighttemp"`
	DayWind      string `json:"daywind"`
	NightWind    string `json:"nightwind"`
	DayPower     string `json:"daypower"`
	NightPower   string `json:"nightpower"`
}



// 全局配置变量
var config Config

// 加载配置文件
func loadConfig(configPath string) error {
	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON配置
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证必要的配置项
	if config.WeatherAPI.Key == "YOUR_AMAP_API_KEY" {
		fmt.Println("警告: 您尚未设置高德API密钥，请在config.json中设置有效的API密钥")
	}

	return nil
}

// 获取天气信息
func getWeather() (string, error) {
	// 使用高德开放平台API
	return getAmapWeather()
}

// 使用高德开放平台API获取天气
func getAmapWeather() (string, error) {
	// 构建请求URL
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%s&extensions=all", 
		config.WeatherAPI.Key, config.WeatherAPI.CityCode)

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("获取天气信息失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 解析JSON响应
	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return "", fmt.Errorf("解析天气数据失败: %v", err)
	}

	// 检查API响应状态
	if weatherResp.Status != "1" {
		return "", fmt.Errorf("天气API返回错误: %s", weatherResp.Info)
	}

	// 获取实时天气
	url = fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%s&extensions=base", 
		config.WeatherAPI.Key, config.WeatherAPI.CityCode)
	
	resp, err = http.Get(url)
	if err != nil {
		return "", fmt.Errorf("获取实时天气信息失败: %v", err)
	}
	defer resp.Body.Close()
	
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取实时天气响应内容失败: %v", err)
	}
	
	var liveWeatherResp WeatherResponse
	if err := json.Unmarshal(body, &liveWeatherResp); err != nil {
		return "", fmt.Errorf("解析实时天气数据失败: %v", err)
	}

	// 格式化天气信息
	var weatherText string
	
	// 添加标题
	weatherText += fmt.Sprintf("【%s天气信息】\n\n", config.WeatherAPI.CityName)
	
	// 实时天气信息
	if len(liveWeatherResp.Lives) > 0 {
		live := liveWeatherResp.Lives[0]
		weatherText += "【实时天气】\n"
		weatherText += fmt.Sprintf("天气: %s\n", live.Weather)
		weatherText += fmt.Sprintf("温度: %s℃\n", live.Temperature)
		weatherText += fmt.Sprintf("风向: %s\n", live.WindDirection)
		weatherText += fmt.Sprintf("风力: %s级\n", live.WindPower)
		weatherText += fmt.Sprintf("湿度: %s%%\n", live.Humidity)
		weatherText += fmt.Sprintf("发布时间: %s\n\n", live.ReportTime)
	}
	
	// 天气预报
	if len(weatherResp.Forecasts) > 0 && len(weatherResp.Forecasts[0].Casts) > 0 {
		weatherText += "【未来天气预报】\n"
		
		// 获取今天和明天的预报
		forecasts := weatherResp.Forecasts[0].Casts
		
		// 今天预报
		if len(forecasts) > 0 {
			today := forecasts[0]
			weatherText += fmt.Sprintf("今天 (%s):\n", today.Date)
			weatherText += fmt.Sprintf("白天: %s %s℃ %s风 %s级\n", 
				today.DayWeather, today.DayTemp, today.DayWind, today.DayPower)
			weatherText += fmt.Sprintf("夜间: %s %s℃ %s风 %s级\n\n", 
				today.NightWeather, today.NightTemp, today.NightWind, today.NightPower)
		}
		
		// 明天预报
		if len(forecasts) > 1 {
			tomorrow := forecasts[1]
			weatherText += fmt.Sprintf("明天 (%s):\n", tomorrow.Date)
			weatherText += fmt.Sprintf("白天: %s %s℃ %s风 %s级\n", 
				tomorrow.DayWeather, tomorrow.DayTemp, tomorrow.DayWind, tomorrow.DayPower)
			weatherText += fmt.Sprintf("夜间: %s %s℃ %s风 %s级\n\n", 
				tomorrow.NightWeather, tomorrow.NightTemp, tomorrow.NightWind, tomorrow.NightPower)
		}
	}
	
	// 添加温馨提示
	weatherText += "【温馨提示】\n"
	
	// 根据天气状况给出建议
	if len(liveWeatherResp.Lives) > 0 {
		live := liveWeatherResp.Lives[0]
		
		if contains(live.Weather, []string{"雨", "阵雨", "雷阵雨", "暴雨"}) {
			weatherText += "今天有雨，出门请记得带伞！\n"
		}
		
		if contains(live.Weather, []string{"雪", "阵雪", "暴雪"}) {
			weatherText += "今天有雪，注意保暖，路面可能湿滑，出行注意安全！\n"
		}
		
		if contains(live.Weather, []string{"雾", "霾"}) {
			weatherText += "今天有雾霾，建议戴口罩出行，减少户外活动时间！\n"
		}
		
		// 根据温度给出建议
		temp := 0
		fmt.Sscanf(live.Temperature, "%d", &temp)
		
		if temp <= 5 {
			weatherText += "天气寒冷，注意保暖，多穿衣服！\n"
		} else if temp >= 30 {
			weatherText += "天气炎热，注意防暑降温，多喝水！\n"
		}
	}
	
	weatherText += "祝您一天愉快！\n"
	
	return weatherText, nil
}

// 判断字符串是否包含指定字符串列表中的任意一个
func contains(s string, substrs []string) bool {
	for _, substr := range substrs {
		if bytes.Contains([]byte(s), []byte(substr)) {
			return true
		}
	}
	return false
}

// 发送消息到飞书
func sendToFeishu(text string) error {
	// 构建消息内容
	message := FeishuMessage{
		MsgType: "text",
		Content: map[string]interface{}{
			"text": text,
		},
	}

	// 将消息转换为JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息JSON编码失败: %v", err)
	}

	// 发送POST请求
	resp, err := http.Post(config.Feishu.WebhookURL, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return fmt.Errorf("发送消息到飞书失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取飞书响应失败: %v", err)
	}

	// 打印响应
	fmt.Printf("飞书响应: %s\n", string(body))

	return nil
}

// 检查是否到了发送时间
func shouldSendNow(hour, minute int) bool {
	now := time.Now()
	return now.Hour() == hour && now.Minute() == minute
}

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config.json", "配置文件路径")
	testMode := flag.Bool("test", false, "测试模式：立即发送一次天气信息")
	flag.Parse()

	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("获取可执行文件路径失败: %v\n", err)
		os.Exit(1)
	}
	execDir := filepath.Dir(execPath)

	// 如果配置文件路径是相对路径，则相对于可执行文件所在目录
	if !filepath.IsAbs(*configPath) {
		*configPath = filepath.Join(execDir, *configPath)
	}

	// 加载配置
	if err := loadConfig(*configPath); err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("天气机器人已启动，将在每天 %02d:%02d 发送天气信息\n",
		config.Schedule.Hour, config.Schedule.Minute)

	// 测试模式：立即发送一次天气信息
	if *testMode {
		fmt.Println("测试模式：立即发送天气信息")
		weatherText, err := getWeather()
		if err != nil {
			fmt.Printf("获取天气信息失败: %v\n", err)
		} else {
			fmt.Println("获取到的天气信息:")
			fmt.Println(weatherText)

			err = sendToFeishu(weatherText)
			if err != nil {
				fmt.Printf("发送到飞书失败: %v\n", err)
			} else {
				fmt.Println("测试消息已发送到飞书")
			}
		}

		// 测试模式下发送完就退出
		if len(flag.Args()) == 0 || flag.Arg(0) != "daemon" {
			os.Exit(0)
		}
	}

	// 定时任务循环
	for {
		if shouldSendNow(config.Schedule.Hour, config.Schedule.Minute) {
			weatherText, err := getWeather()
			if err != nil {
				fmt.Printf("获取天气信息失败: %v\n", err)
			} else {
				err = sendToFeishu(weatherText)
				if err != nil {
					fmt.Printf("发送到飞书失败: %v\n", err)
				} else {
					fmt.Printf("%s 天气信息已发送到飞书\n", time.Now().Format("2006-01-02 15:04:05"))
				}
			}

			// 等待一分钟，避免在同一分钟内重复发送
			time.Sleep(time.Minute)
		}

		// 每30秒检查一次时间
		time.Sleep(30 * time.Second)
	}
}
