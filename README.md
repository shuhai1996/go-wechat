#go-wechat
###go 聊天室服务端

Dockerfile
````docker
# 基础镜像 golang
FROM golang
# 运行环境
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=1
# 工作目录
WORKDIR [你的工作目录]
# 将文件copy到镜像相同目录中
COPY . .
# 暴露端口
EXPOSE 8000
# 执行命令
RUN go build socket/main.go
# docker run 时执行
CMD ["./main"]
````