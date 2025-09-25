package main

import (
	"cos/conf"
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"net/http"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := NewRouter()

	err := r.Run(conf.GlobalConfig.Server.Addr)
	if err != nil {
		panic(err)
	}
}

type TestMaskPayload struct {
	M0 string `json:"m0"`
	M1 string `json:"m1"`
}

type Session struct {
	SessionId string `json:"sessionid"`
	MaskName  string `json:"maskname"`
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CORS())
	v := r.Group("api/v") // no auth

	{
		v.POST("translate", func(context *gin.Context) {
			var arg openai.ChatCompletionRequest
			if err := context.BindJSON(&arg); err != nil {
				panic(err)
			}
			println(arg.Model)
			s := Session{}
			MustUnmarshal([]byte(arg.Model), &s)
			if s.MaskName != "translate" {
				panic(`s.MaskName != "translate"`)
			}
			// todo: s.SessionId
			last := arg.Messages[len(arg.Messages)-1]
			llmRes := MustLLM(
				"https://yourapi.cn/v1/chat/completions",
				"sk-n26nHrTj7lmzvQmPRxVLVLR5SqdXDbaaHptd07wo1ul4yMuF",
				openai.ChatCompletionRequest{
					Model:    "gpt-5-2025-08-07", //"gpt-5-all",
					Stream:   false,
					Messages: Translate(last.Content),
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: openai.ChatCompletionResponseFormatTypeJSONObject,
					},
					ReasoningEffort: "minimal",
				},
			)
			// TODO: assert json
			payload := string(MustMarshal(TestMaskPayload{
				M0: Resp(llmRes),
				M1: "extra info",
			}))
			println(payload)
			context.JSON(http.StatusOK, openai.ChatCompletionResponse{
				ID:      "",
				Object:  "",
				Created: 0,
				Model:   "",
				Choices: []openai.ChatCompletionChoice{
					{
						Index: 0,
						Message: openai.ChatCompletionMessage{
							Role:    "assistant",
							Content: payload,
						},
						FinishReason: "",
					},
				},
				Usage:               openai.Usage{},
				SystemFingerprint:   "",
				PromptFilterResults: nil,
			})
		})
		v.POST("argue", func(context *gin.Context) {
			var arg openai.ChatCompletionRequest
			if err := context.BindJSON(&arg); err != nil {
				panic(err)
			}
			s := Session{}
			MustUnmarshal([]byte(arg.Model), &s)
			if s.MaskName != "argue" {
				panic(`s.MaskName != "argue"`)
			}
			// todo: s.SessionId
			llmRes := MustLLM(
				"https://yourapi.cn/v1/chat/completions",
				"sk-n26nHrTj7lmzvQmPRxVLVLR5SqdXDbaaHptd07wo1ul4yMuF",
				openai.ChatCompletionRequest{
					Model:    "gpt-5-2025-08-07", //"gpt-5-all",
					Stream:   false,
					Messages: Argue(arg.Messages),
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: openai.ChatCompletionResponseFormatTypeJSONObject,
					},
					ReasoningEffort: "minimal",
				},
			)
			payload := string(MustMarshal(TestMaskPayload{
				M0: Resp(llmRes),
				M1: "extra info",
			}))
			println(payload)
			context.JSON(http.StatusOK, openai.ChatCompletionResponse{
				ID:      "",
				Object:  "",
				Created: 0,
				Model:   "",
				Choices: []openai.ChatCompletionChoice{
					{
						Index: 0,
						Message: openai.ChatCompletionMessage{
							Role:    "assistant",
							Content: payload,
						},
						FinishReason: "",
					},
				},
				Usage:               openai.Usage{},
				SystemFingerprint:   "",
				PromptFilterResults: nil,
			})
		})
	}
	return r
}
