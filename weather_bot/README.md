# 天气机器人

这是一个简单的Go语言工具，用于定时获取指定城市的天气信息，并通过飞书机器人发送到飞书群聊中。

## 功能特点

- 定时获取指定城市的天气信息
- 支持实时天气和天气预报
- 根据天气情况提供穿衣、出行建议
- 通过飞书机器人webhook发送消息
- 可配置发送时间
- 支持配置文件，易于修改设置

## 使用方法

### 前置条件

1. 安装Go语言环境（1.16或更高版本）或Docker环境
2. 注册高德开放平台账号并创建应用，获取API Key
3. 在飞书管理后台创建自定义机器人，获取Webhook URL

### 配置

编辑`config.json`文件，修改以下配置：

```json
{
  "weather_api": {
    "key": "YOUR_AMAP_API_KEY",  // 替换为你的高德API密钥
    "city_code": "110000",      // 城市编码，默认为北京
    "city_name": "北京"         // 城市名称
  },
  "feishu": {
    "webhook_url": "https://open.feishu.cn/open-apis/bot/v2/hook/a059c069-a6bf-4937-9e61-cf9b2c9c667d"
  },
  "schedule": {
    "hour": 8,                  // 每天发送时间：小时
    "minute": 0                 // 每天发送时间：分钟
  }
}
```

城市编码可以在[高德开放平台城市编码表](https://lbs.amap.com/api/webservice/download)中查询。

### 编译

#### 使用Makefile（推荐）

```bash
cd weather_bot

# 编译
make

# 测试运行（发送一次天气信息）
make test

# 安装到系统（需要管理员权限）
sudo make install

# 查看所有可用命令
make help
```

#### 手动编译

```bash
cd weather_bot
go build
```

### 运行

#### 使用Docker（推荐）

项目提供了Dockerfile和docker-compose.yml，方便在容器环境中运行：

```bash
# 使用Docker Compose构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

或者直接使用Docker命令：

```bash
# 构建镜像
docker build -t weather-bot .

# 运行容器
docker run -d --name weather-bot -v $(pwd)/config.json:/app/config.json weather-bot

# 测试模式（立即发送一次）
docker run --rm -v $(pwd)/config.json:/app/config.json weather-bot /app/weather_bot -test
```

#### 使用启动脚本

提供了便捷的启动脚本，自动处理编译和参数设置：

```bash
# Linux/macOS
./start.sh          # 正常模式启动
./start.sh -test    # 测试模式（发送后退出）
./start.sh -test -daemon  # 测试模式（发送后继续运行）

# Windows
start.bat          # 正常模式启动
start.bat -test    # 测试模式（发送后退出）
start.bat -test -daemon  # 测试模式（发送后继续运行）
```

#### 直接运行可执行文件

```bash
# 使用默认配置文件（config.json）
./weather_bot

# 指定配置文件路径
./weather_bot -config /path/to/your/config.json

# 测试模式：立即发送一次天气信息并退出
./weather_bot -test

# 测试模式但不退出（继续运行定时任务）
./weather_bot -test daemon
```

## 定时任务

如果希望将此工具作为系统服务长期运行，可以：

### 在Linux上使用systemd服务

项目提供了一个systemd服务文件模板 `weather-bot.service`：

1. 编辑服务文件，修改路径：
   ```bash
   sudo nano /etc/systemd/system/weather-bot.service
   ```
   将文件内容复制到此处，并修改路径。

2. 重新加载systemd配置：
   ```bash
   sudo systemctl daemon-reload
   ```

3. 启动服务：
   ```bash
   sudo systemctl start weather-bot
   ```

4. 设置开机自启：
   ```bash
   sudo systemctl enable weather-bot
   ```

5. 查看服务状态：
   ```bash
   sudo systemctl status weather-bot
   ```

### 在Linux/macOS上使用crontab

```bash
# 编辑crontab
crontab -e

# 添加以下内容（每天早上7:30运行）
30 7 * * * cd /path/to/weather_bot && ./weather_bot -test
```

### 在Windows上使用计划任务

可以使用Windows任务计划程序设置定时任务，执行`weather_bot.exe -test`命令。

1. 打开任务计划程序（Task Scheduler）
2. 创建基本任务
3. 设置每天触发
4. 设置操作为启动程序
5. 程序路径设置为`C:\path\to\weather_bot\weather_bot.exe`
6. 添加参数`-test`

## 注意事项

- 程序默认不会立即发送天气信息，除非使用`-test`参数
- 确保网络连接正常
- 高德API有调用次数限制，请合理使用
- 配置文件默认从程序所在目录读取，也可以通过`-config`参数指定绝对路径
- 使用启动脚本前，请确保脚本有执行权限：`chmod +x start.sh`
- systemd服务文件需要根据实际安装路径修改后才能使用
- 在Windows系统中，如果使用计划任务，建议设置工作目录为程序所在目录
- 使用Docker时，请确保配置文件已正确设置API Key和Webhook URL
- Docker容器默认使用亚洲/上海时区，如需修改请在docker-compose.yml中调整