# 使用多阶段构建减小镜像大小
# 第一阶段：构建
FROM golang:1.23 AS builder

# 设置工作目录
WORKDIR /app

# 设置国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOSUMDB=sum.golang.google.cn

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建静态链接的可执行文件（重要修改）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/runServer.go ./cmd/main.go

# 第二阶段：运行
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录并创建应用目录
WORKDIR /app

# 从构建阶段复制可执行文件（路径修改）
COPY --from=builder /app/main /app/

# 确保可执行权限（重要添加）
RUN chmod +x /app/main

# 暴露端口
EXPOSE 9000

# 运行应用（路径修改）
CMD ["/app/main", "run"]