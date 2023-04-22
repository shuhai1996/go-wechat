package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sashabaranov/go-openai"
	"go-wechat/config"
	"go-wechat/middleware"
	"go-wechat/socket"
	"io"
	"log"
)

type GptServices struct {
}

func (GptServices) GetText(str string) string {
	client := openai.NewClient(config.Get("open.api.token"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: str,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "ChatCompletion error:" + err.Error()
	}

	return resp.Choices[0].Message.Content
}

func (GptServices) GetStream(user *middleware.User, conn *websocket.Conn, str string) (err error, res string) {
	client := openai.NewClient(config.Get("open.api.token"))
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 2048,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: str,
			},
		},
		Stream:           true,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		TopP:             1,
		Temperature:      0,
	}
	res = ""
	stream, er := client.CreateChatCompletionStream(context.Background(), req)
	if er != nil {
		err = er
		log.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, er := stream.Recv()
		if errors.Is(er, io.EOF) {
			data := &socket.Response{
				Type:     "customer",
				Username: user.Nickname,
				Message:  "!!!Words finished!!!",
			}

			bytes, e := json.Marshal(data)
			if e != nil {
				log.Printf("\njson decode error: %v\n", err)
				err = e
				return
			}

			_ = conn.WriteMessage(websocket.TextMessage, bytes)

			return
		}

		if er != nil {
			log.Printf("\nStream error: %v\n", err)

			data := &socket.Response{
				Type:     "customer",
				Username: user.Nickname,
				Message:  "!!!Words finished!!!",
			}

			bytes, e := json.Marshal(data)
			if e != nil {
				log.Printf("\njson decode error: %v\n", err)
				err = e
				return
			}

			_ = conn.WriteMessage(websocket.TextMessage, bytes)
			err = er
			return
		}

		re := response.Choices[0].Delta.Content

		data := &socket.Response{
			Type:     "customer",
			Username: user.Nickname,
			Message:  re,
		}

		bytes, e := json.Marshal(data)
		if e != nil {
			fmt.Printf("\njson decode error: %v\n", err)
			err = e
			return
		}

		_ = conn.WriteMessage(websocket.TextMessage, bytes)
		log.Println("Stream finished:"+re)
		res += re
	}
}
