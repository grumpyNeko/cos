package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
)

func MustLLM(url string, apikey string, reqBody openai.ChatCompletionRequest) openai.ChatCompletionResponse {
	if reqBody.Stream {
		panic(fmt.Sprintf(`[payload.Stream != false] payload.Stream=true`))
	}
	reqBodyBytes := MustMarshal(reqBody)
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		panic(err)
	}
	req.Header = make(http.Header)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if code := res.StatusCode; code != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		panic(fmt.Sprintf("code not 200, body: %s", string(bodyBytes)))
	}
	resBody := openai.ChatCompletionResponse{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		panic(err)
	}
	return resBody
}

func Resp(resp openai.ChatCompletionResponse) string {
	return resp.Choices[0].Message.Content
}

// 这个是试图在开源基础上微操
func MustLLM0(url string, apikey string, payload openai.ChatCompletionRequest) openai.ChatCompletionResponse {
	config := openai.DefaultConfig(apikey)
	config.BaseURL = url
	c := openai.NewClientWithConfig(config)
	if payload.User != "" {
		panic(fmt.Sprintf(`[payload.User != ""] payload.User=%s`, payload.User))
	}
	payload.User = apikey
	if payload.Stream {
		panic(fmt.Sprintf(`[payload.Stream != false] payload.Stream=true`))
	}
	resp, err := c.CreateChatCompletion(
		context.TODO(),
		payload,
	)
	if err != nil {
		panic(err)
	}
	return resp
}
