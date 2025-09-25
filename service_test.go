package main

import (
	"github.com/sashabaranov/go-openai"
	"testing"
)

func use(args ...interface{}) {}

func Test_must(t *testing.T) {
	res := MustLLM(
		"https://yourapi.cn/v1/chat/completions",
		"sk-n26nHrTj7lmzvQmPRxVLVLR5SqdXDbaaHptd07wo1ul4yMuF",
		openai.ChatCompletionRequest{
			Model:    "gpt-5-all",
			Stream:   false,
			Messages: Translate("还挣了18万的差价是吧"),
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			ReasoningEffort: "minimal",
		},
	)
	println(Resp(res))
	println(res.Usage.PromptTokens)
	println(res.Usage.CompletionTokens)
}
