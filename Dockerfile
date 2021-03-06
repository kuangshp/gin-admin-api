# 获取go的版本
FROM golang:1.15

# 为镜像配置环境
ENV GO115MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"

# 这个为代码目录
WORKDIR /home/gin-admin-api

# 将代码复制到容器中
COPY . .
# 代码编译成二进制文件,名字为app
RUN go build -o app .

# 声明端口
EXPOSE 9090

# 启动容器的命令
CMD ["./app"]