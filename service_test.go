package main

import (
	"context"
	"fmt"
	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/genai"
	"io"
	"testing"
)

// models/gemini-2.5-flash-preview-05-20
// models/gemini-2.5-pro-preview-05-06
func Test_gemini(t *testing.T) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  "AIzaSyDaHj03PPPJ3HsM_GwtDqnNjrxhXpq5zzk",
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(err)
	}
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-pro-preview-05-06",
		genai.Text("Explain how AI works in a few words"),
		&genai.GenerateContentConfig{
			HTTPOptions: &genai.HTTPOptions{
				BaseURL: "https://api-proxy.me/gemini",
			},
			SystemInstruction:    nil,
			Temperature:          nil,
			TopP:                 nil,
			TopK:                 nil,
			CandidateCount:       0,
			MaxOutputTokens:      0,
			StopSequences:        nil,
			ResponseLogprobs:     false,
			Logprobs:             nil,
			PresencePenalty:      nil,
			FrequencyPenalty:     nil,
			Seed:                 nil,
			ResponseMIMEType:     "",
			ResponseSchema:       nil,
			RoutingConfig:        nil,
			ModelSelectionConfig: nil,
			SafetySettings:       nil,
			Tools:                nil,
			ToolConfig:           nil,
			Labels:               nil,
			CachedContent:        "",
			ResponseModalities:   nil,
			MediaResolution:      "",
			SpeechConfig:         nil,
			AudioTimestamp:       false,
			ThinkingConfig:       nil,
		},
	)
	if err != nil {
		panic(err)
	}
	use(result)
}

func Test_ds(t *testing.T) {
	client, _ := deepseek.NewClient("sk-c48b624ad64948a5b51e4baf1064da81")
	temperature := float32(1.5)
	req := &request.ChatCompletionsRequest{
		Model: deepseek.DEEPSEEK_REASONER_MODEL,
		Messages: []*request.Message{
			{
				Role:    "system",
				Content: MustReadAll("./word.template"),
			},
			{
				Role:    "user",
				Content: "frisk",
			},
		},
		Stream:      true,
		Temperature: &temperature,
	}
	sr, err := client.StreamChatCompletionsReasoner(context.Background(), req)
	if err != nil {
		panic(err)
	}

	// process response
	for {
		resp, err := sr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if resp.Choices[0].Delta.Content != "" {
			fmt.Print(resp.Choices[0].Delta.Content)
		} else {
			fmt.Print(resp.Choices[0].Delta.ReasoningContent)
		}
	}
	fmt.Println()
}

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
	//res := MustLLM(
	//	"https://api-proxy.me/xai/v1/chat/completions",
	//	"xai-4Kag7Eqy8UNK1zCoXEUuELtwLIgTJ4DmXqWDryuVzSsAf30YgsZ05wRPTtCqmoVkJXqwMsC75A4mIgyR",
	//	openai.ChatCompletionRequest{
	//		Model:    "grok-4",
	//		Stream:   false,
	//		Messages: Translate("佬们有没有什么和互联网相关的副业推荐一下"),
	//		ResponseFormat: &openai.ChatCompletionResponseFormat{
	//			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
	//		},
	//		//ReasoningEffort: "minimal",
	//	},
	//)
	println(Resp(res))
	println(res.Usage.PromptTokens)
	println(res.Usage.CompletionTokens)
}
