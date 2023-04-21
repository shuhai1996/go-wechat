# go-wechat
### go 聊天室服务端

Dockerfile
```dockerfile
# 基础镜像 golang
FROM golang
# 运行环境
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=1
# 工作目录
WORKDIR /go/src/go-wechat

# 将依赖项下载到缓存中
COPY go.mod go.sum ./
RUN go mod download

# 将文件copy到镜像相同目录中
COPY . /go/src/go-wechat
# 暴露端口
EXPOSE 8000
# 执行命令
RUN go build main.go
# docker run 时执行
CMD ["./main"]
```

创建conf.ini文件
```
[development]
jwt.secret=
open.api.token=
rpc.url=
```

在根目录执行
```
docker build -t go-wechat:v1 .
docker run -d -p 8000:8000 --name go-wechat go-wechat:v1 
```