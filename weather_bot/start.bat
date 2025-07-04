@echo off
REM 天气机器人启动脚本 (Windows版)

REM 切换到脚本所在目录
cd /d "%~dp0"

REM 检查是否有可执行文件
if not exist weather_bot.exe (
    echo 正在编译天气机器人...
    go build
    if errorlevel 1 (
        echo 编译失败，请检查错误信息
        exit /b 1
    )
    echo 编译成功
)

REM 检查配置文件
if not exist config.json (
    echo 错误：找不到配置文件 config.json
    exit /b 1
)

REM 解析命令行参数
set TEST_MODE=false
set DAEMON_MODE=false

:parse_args
if "%~1"=="" goto run_program
if "%~1"=="-test" set TEST_MODE=true
if "%~1"=="--test" set TEST_MODE=true
if "%~1"=="-daemon" set DAEMON_MODE=true
if "%~1"=="--daemon" set DAEMON_MODE=true
shift
goto parse_args

:run_program
REM 运行程序
if "%TEST_MODE%"=="true" (
    if "%DAEMON_MODE%"=="true" (
        echo 以测试模式启动天气机器人（守护进程模式）
        weather_bot.exe -test daemon
    ) else (
        echo 以测试模式启动天气机器人（发送后退出）
        weather_bot.exe -test
    )
) else (
    echo 启动天气机器人（定时发送模式）
    weather_bot.exe
)