package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go-wechat/config"
	"go-wechat/middleware"
	"go-wechat/rpc/client"
	"go-wechat/service"
	"go-wechat/socket"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

var sev = &service.GptServices{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 校验请求头部信息
	authHeader := r.Header.Get("Sec-WebSocket-Protocol")
	user, err := middleware.JwtAuth(authHeader)
	if err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	header := make(http.Header)
	header.Set("Sec-WebSocket-Protocol", authHeader)
	conn, err := upgrader.Upgrade(w, r, header)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("client connected!")

	rpcConn, re := grpc.Dial(config.Get("rpc.url"), grpc.WithInsecure())
	if re != nil {
		log.Println("dialing rpc " + re.Error())
	}

	defer conn.Close()
	defer rpcConn.Close()

	// 处理消息
	// 处理消息
	if err := writeResponseMessage(conn, user, nil, true, rpcConn); err != nil {
		log.Println(err)
		return
	}

	for {
		// 读取客户端消息
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if client.GetChatUserInfo(rpcConn, user.Username) >= 50 { // 超过50限制不让继续
			log.Println("超过最大限制")
			return
		}

		// 处理消息
		if err := writeResponseMessage(conn, user, p, false, rpcConn); err != nil {
			log.Println(err)
			return
		}
	}
}

func writeResponseMessage(conn *websocket.Conn, user *middleware.User, p []byte, isConnect bool, rpcConn *grpc.ClientConn) (err error){
	var data *socket.Response
	if isConnect {
		//转json
		data = &socket.Response{
			Type: "server",
			Username: user.Nickname,
			Message: "Hello！你可以提任何你想提的问题!",
		}
		bytes,_ := json.Marshal(data)

		err = conn.WriteMessage(websocket.TextMessage, bytes)
	} else {
		text := string(p)
		client.SendMessage(rpcConn, user.ID, text, 0)
		er, response := sev.GetStream(user, conn, text)
		if er != nil {
			log.Println("Stream error:"+ er.Error())
			return
		}
		log.Println("Stream finished:"+response)
		client.SendMessage(rpcConn, 0, response, 0)
	}

	return
}

func main() {
	config.Setup("development") //加载配置文件
	http.HandleFunc("/ws", wsHandler)

	log.Println("Server started at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

