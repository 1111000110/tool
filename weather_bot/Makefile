# 天气机器人 Makefile

.PHONY: build clean test install

# 默认目标
all: build

# 编译
build:
	go build -o weather_bot

# 清理
clean:
	rm -f weather_bot
	@echo "清理完成"

# 测试运行
test:
	go build -o weather_bot
	./weather_bot -test

# 安装（Linux/macOS）
install: build
	mkdir -p $(DESTDIR)/usr/local/bin
	cp weather_bot $(DESTDIR)/usr/local/bin/
	mkdir -p $(DESTDIR)/etc/weather_bot
	@if [ ! -f $(DESTDIR)/etc/weather_bot/config.json ]; then \
		cp config.json $(DESTDIR)/etc/weather_bot/; \
	fi
	@echo "安装完成，请修改 /etc/weather_bot/config.json 配置文件"

# 卸载（Linux/macOS）
uninstall:
	rm -f $(DESTDIR)/usr/local/bin/weather_bot
	@echo "卸载完成"

# 帮助
help:
	@echo "可用命令："
	@echo "  make build     - 编译程序"
	@echo "  make clean     - 清理编译文件"
	@echo "  make test      - 测试运行（发送一次天气信息）"
	@echo "  make install   - 安装到系统（需要管理员权限）"
	@echo "  make uninstall - 从系统卸载（需要管理员权限）"