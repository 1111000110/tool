FROM golang:1.20-alpine AS builder

WORKDIR /app

# 复制Go模块文件
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY main.go ./

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o weather_bot

# 使用轻量级基础镜像
FROM alpine:latest

# 安装CA证书，用于HTTPS请求
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为亚洲/上海
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/weather_bot /app/

# 复制配置文件（注意：需要在构建时提供config.json）
COPY config.json /app/

# 运行非root用户
RUN adduser -D appuser
RUN chown -R appuser:appuser /app
USER appuser

# 容器启动命令
CMD ["/app/weather_bot"]