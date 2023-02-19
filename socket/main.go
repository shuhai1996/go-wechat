package main

import (
	"encoding/json"
	"fmt"
	"go-wechat/middleware"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Response struct {
	Type string `json:"type"`
	Username string `json:"username"`
	Message string `json:"message"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 校验请求头部信息
	authHeader := r.Header.Get("Sec-WebSocket-Protocol")
	fmt.Println(authHeader)
	user, err := middleware.JwtAuth(authHeader)
	if err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//user := &middleware.User{}
	header := make(http.Header)
	header.Set("Sec-WebSocket-Protocol", authHeader)
	conn, err := upgrader.Upgrade(w, r, header)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("client connected!")

	defer conn.Close()

	// 处理消息
	if err := conn.WriteMessage(websocket.TextMessage, writeResponseMessage(user, nil, true)); err != nil {
		log.Println(err)
		return
	}

	for {
		// 读取客户端消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}


		// 处理消息
		if err := conn.WriteMessage(messageType, writeResponseMessage(user, p, false)); err != nil {
			log.Println(err)
			return
		}
	}
}

func writeResponseMessage(user *middleware.User, p []byte, isConnect bool)  []byte{
	var data *Response
	if isConnect {
		//转json
		data = &Response{
			Type: "server",
			Username: user.Nickname,
			Message: "join the chatroom",
		}
	} else {

		data = &Response{
			Type: "customer",
			Username: user.Nickname,
			Message: string(p),
		}
	}

	bytes,_ := json.Marshal(data)

	return bytes
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	log.Println("Server started at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

