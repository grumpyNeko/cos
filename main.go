package main

import (
	"cos/conf"
	"fmt"
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

type ModelPayload struct {
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
			s := ModelPayload{}
			MustUnmarshal([]byte(arg.Model), &s)
			if s.MaskName != "translate" {
				panic(`s.MaskName != "translate"`)
			}
			// todo: s.SessionId
			last := arg.Messages[len(arg.Messages)-1] // todo: 只取最后一个user message
			llmRes := MustLLM(
				"https://aihubmix.com/v1/chat/completions",
				"sk-4JGSa4uexQfH6VIjD366C77c11F74bC6Bd919dEb6055Dd31",
				openai.ChatCompletionRequest{
					Model:    "gpt-5",
					Stream:   false,
					Messages: Translate(last.Content),
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: openai.ChatCompletionResponseFormatTypeJSONObject,
					},
					ReasoningEffort: "minimal",
				},
			)
			println(Resp(llmRes))
			type translateLLMResponse struct {
				Literal string `json:"literal"`
				Free    string `json:"free"`
			}
			payload := translateLLMResponse{}
			MustUnmarshal([]byte(Resp(llmRes)), &payload)
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
							Content: string(MustMarshal(payload)),
						},
						FinishReason: "",
					},
				},
				Usage:               openai.Usage{},
				SystemFingerprint:   "",
				PromptFilterResults: nil,
			})
		})
		v.POST("argue", argueHandler)
	}
	return r
}

func argueHandler(context *gin.Context) {
	var arg openai.ChatCompletionRequest
	if err := context.BindJSON(&arg); err != nil {
		panic(err)
	}
	s := ModelPayload{}
	MustUnmarshal([]byte(arg.Model), &s)
	if s.MaskName != "argue" {
		panic(`s.MaskName != "argue"`)
	}

	session, ok := sessionStore[s.SessionId]
	if !ok {
		println(fmt.Sprintf(`sessionid=%s not found in sessionStore`, s.SessionId))
		session = Session{msgList: make([]string, 16)}
	}
	session.msgList = append(session.msgList, arg.Messages[len(arg.Messages)-1].Content)
	sessionStore[s.SessionId] = session
	llmRes := MustLLM(
		"https://aihubmix.com/v1/chat/completions",
		"sk-4JGSa4uexQfH6VIjD366C77c11F74bC6Bd919dEb6055Dd31",
		openai.ChatCompletionRequest{
			Model:    "gpt-5",
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
}
