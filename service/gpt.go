package service

import (
	"bytes"
	"encoding/json"
	"github.com/thedevsaddam/gojsonq"
	"go-wechat/config"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// PostBody post 请求
type PostBody struct {
	Model       string `json:"model"`
	Messages  []map[string]string `json:"messages"`
	Temperature int    `json:"temperature"`
	MaxTokens   int    `json:"max_tokens"`
}

type GptServices struct {
}

func (GptServices) GetText(str string) string {
	d, err := PostJsonWithHeaders(str)
	if err != nil {
		return err.Error()
	}
	gq := gojsonq.New().FromString(string(d))
	district := gq.Find("choices.[0].message.content")
	if district == nil {
		return "Remote API returns error."
	}
	return district.(string)
}

func PostJsonWithHeaders(str string) (b []byte, e error) {
	url := "https://api.openai.com/v1/chat/completions"
	// 构造POST请求
	postBody := &PostBody{
		Model:       "gpt-3.5-turbo",
		Messages: []map[string]string{
			{"role": "user", "content": str},
		},
		Temperature: 0,
		MaxTokens:   2048,
	}
	// struct 转json
	body, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+config.Get("open.api.token"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		e = err
		return
	}
	defer resp.Body.Close()

	b, e = ioutil.ReadAll(resp.Body)
	log.Println(string(b))
	return
}
