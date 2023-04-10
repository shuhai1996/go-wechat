package service

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"go-wechat/config"
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
