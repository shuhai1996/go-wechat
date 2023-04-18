package main

import (
	"encoding/json"
	"fmt"
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

	rpcConn, re := grpc.Dial(config.Get("rpc.url"), grpc.WithInsecure())
	if re != nil {
		log.Println("dialing rpc " + re.Error())
	}
	log.Println(rpcConn.GetState())


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
			response = er.Error()
			return
		}
		//time.Sleep(2000 * time.Millisecond)
		//data = &Response{
		//	Type: "customer",
		//	Username: user.Nickname,
		//	Message: "以下是一个简单的Go语言爬虫程序，用于爬取指定网站的所有链接：\n\n```go\npackage main\n\nimport (\n    \"fmt\"\n    \"net/http\"\n    \"io/ioutil\"\n    \"regexp\"\n)\n\nfunc main() {\n    url := \"https://www.example.com\"\n    visited := make(map[string]bool)\n    crawl(url, visited)\n}\n\nfunc crawl(url string, visited map[string]bool) {\n    if visited[url] {\n        return\n    }\n    visited[url] = true\n    fmt.Println(\"Crawling:\", url)\n    resp, err := http.Get(url)\n    if err != nil {\n        fmt.Println(\"Error:\", err)\n        return\n    }\n    defer resp.Body.Close()\n    body, err := ioutil.ReadAll(resp.Body)\n    if err != nil {\n        fmt.Println(\"Error:\", err)\n        return\n    }\n    re := regexp.MustCompile(`<a\\s+(?:[^>]*?\\s+)?href=\"([^\"]*)\"`)\n    matches := re.FindAllStringSubmatch(string(body), -1)\n    for _, match := range matches {\n        link := match[1]\n        if len(link) == 0 || link[0] == '#' {\n            continue\n        }\n        if link[0] == '/' {\n            link = url + link\n        }\n        crawl(link, visited)\n    }\n}\n```\n\n该程序使用了Go语言的标准库中的http和regexp包，通过正则表达式匹配页面中的链接，并递归爬取所有链接。在爬取过程中，使用了一个visited map来记录已经访问过的链接，避免重复访问。",
		//}

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

