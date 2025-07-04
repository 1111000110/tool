#!/bin/bash

# 天气机器人启动脚本

# 获取脚本所在目录
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR"

# 检查是否有可执行文件
if [ ! -f "./weather_bot" ]; then
    echo "正在编译天气机器人..."
    go build
    if [ $? -ne 0 ]; then
        echo "编译失败，请检查错误信息"
        exit 1
    fi
    echo "编译成功"
fi

# 检查配置文件
if [ ! -f "./config.json" ]; then
    echo "错误：找不到配置文件 config.json"
    exit 1
fi

# 检查API Key是否已设置
API_KEY=$(grep -o '"key": "[^"]*"' ./config.json | cut -d '"' -f 4)
if [ "$API_KEY" = "YOUR_AMAP_API_KEY" ]; then
    echo "警告：您尚未设置高德API密钥，请在config.json中设置有效的API密钥"
    exit 1
fi

# 解析命令行参数
TEST_MODE=false
DAEMON_MODE=false

for arg in "$@"; do
    case $arg in
        -test|--test)
            TEST_MODE=true
            ;;
        -daemon|--daemon)
            DAEMON_MODE=true
            ;;
    esac
done

# 运行程序
if [ "$TEST_MODE" = true ]; then
    if [ "$DAEMON_MODE" = true ]; then
        echo "以测试模式启动天气机器人（守护进程模式）"
        ./weather_bot -test daemon
    else
        echo "以测试模式启动天气机器人（发送后退出）"
        ./weather_bot -test
    fi
else
    echo "启动天气机器人（定时发送模式）"
    ./weather_bot
fi