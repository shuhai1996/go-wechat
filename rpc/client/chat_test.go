package client

import (
	"flag"
	"fmt"
	"go-wechat/config"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestGetChatUserInfo(t *testing.T) {
	config.Setup("development") //加载配置文件
	var addr = flag.String("addr", config.Get("rpc.url"), "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())

	if err != nil {
		log.Fatal("dialing" + err.Error())
	}
	fmt.Println("server state:" + conn.GetState().String())
	defer conn.Close()
	GetChatUserInfo(conn, "test")
	// Check if the message has been written
}
