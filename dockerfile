# 第一阶段：构建 (Builder)
# 使用官方 Go 镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置国内代理，加速依赖下载 (国内环境必备)
ENV GOPROXY=https://goproxy.cn,direct

# 1. 先拷贝依赖定义文件，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 2. 拷贝所有源代码
COPY . .

# 3. 编译 Go 程序
# -ldflags="-s -w" 可以减小二进制文件体积
RUN go build -ldflags="-s -w" -o shortener shortener.go

# 第二阶段：运行 (Runner)
# 使用轻量级 Alpine 镜像
FROM alpine:latest

WORKDIR /app

# 安装必要的时区数据 (否则日志时间可能是 UTC)
RUN apk add --no-cache tzdata

# 从构建阶段拷贝编译好的二进制文件
COPY --from=builder /app/shortener .

# 拷贝配置文件目录
# 注意：这里拷贝的是 .yaml 文件，而不是 .local
# 在实际运行时，我们会通过 docker volume 挂载覆盖它，或者直接用这个默认的
COPY etc/shortener-api.yaml etc/shortener-api.yaml

# 暴露端口
EXPOSE 8888

# 启动命令
CMD ["./shortener", "-f", "etc/shortener-api.yaml"]