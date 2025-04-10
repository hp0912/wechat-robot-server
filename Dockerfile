# 构建阶段
FROM golang:alpine AS builder

# 设置工作目录和构建环境
WORKDIR /build

ENV CGO_ENABLED=0

# 复制并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并构建
COPY . .
RUN go build -ldflags="-s -w" -o wechat-robot-server

FROM debian:stable-slim

# 预定义参数
ENV TZ="Asia/Shanghai"
ENV WECHAT_PORT=9001

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build ./

EXPOSE ${WECHAT_PORT}

# 开始运行
CMD ["/app/wechat-robot-server"]
