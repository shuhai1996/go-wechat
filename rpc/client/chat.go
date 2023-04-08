package client

import (
	"context"
	pb "go-wechat/rpc/proto/chat"
	"google.golang.org/grpc"
	"log"
	"time"
)

func GetChatUserInfo(conn *grpc.ClientConn, name string) int {

	c := pb.NewChatServiceClient(conn)

	//Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetUserInfo(ctx, &pb.UserInfoRequest{Name: name})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	log.Printf("Name: %s Count: %v", r.GetName(), r.GetCount())

	return int(r.GetCount())
}

func SendMessage(conn *grpc.ClientConn, userId int, text string, room int) {
	c := pb.NewChatServiceClient(conn)

	//Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := c.SendMessage(ctx, &pb.Message{Userid: int32(userId), Text: text, RoomId: int32(room)})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	return
}