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
	last := arg.Messages[len(arg.Messages)-1]
	if last.Role != openai.ChatMessageRoleUser {
		panic(`last.Role != openai.ChatMessageRoleUser`)
	}
	history = append(history, Msg{
		Role:    openai.ChatMessageRoleAssistant,
		Content: last.Content,
	})
	//session, ok := sessionStore[s.SessionId]
	//if !ok {
	//	println(fmt.Sprintf(`sessionid=%s not found in sessionStore`, s.SessionId))
	//	session = Session{msgList: make([]string, 16)}
	//}
	//session.msgList = append(session.msgList, arg.Messages[len(arg.Messages)-1].Content)
	//sessionStore[s.SessionId] = session

	payload := ArgueGen(arg.Messages)
	println(fmt.Sprintf(" %+v", payload))
	llmRes := MustLLM(
		"https://aihubmix.com/v1/chat/completions",
		"sk-4JGSa4uexQfH6VIjD366C77c11F74bC6Bd919dEb6055Dd31",
		openai.ChatCompletionRequest{
			Model:    "gpt-5",
			Stream:   false,
			Messages: payload,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			ReasoningEffort: "minimal",
		},
	)
	println(fmt.Sprintf(" %+v", llmRes))
	type ArgueGenResp struct {
		DebateMode bool   `json:"debate_mode"`
		ReplyLen   int    `json:"reply_len"`
		UserMsgEng string `json:"user_msg_eng"`
		Reply      string `json:"reply"`
	}
	genResp := ArgueGenResp{}
	MustUnmarshal([]byte(Resp(llmRes)), &genResp)

	payload0 := ArgueRefine(genResp.Reply, genResp.ReplyLen)
	println(fmt.Sprintf(" %+v", payload0))
	llmRes0 := MustLLM(
		"https://aihubmix.com/v1/chat/completions",
		"sk-4JGSa4uexQfH6VIjD366C77c11F74bC6Bd919dEb6055Dd31",
		openai.ChatCompletionRequest{
			Model:    "gpt-5",
			Stream:   false,
			Messages: payload0,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			ReasoningEffort: "minimal",
		},
	)
	println(fmt.Sprintf(" %+v", llmRes0))
	type ArgueRefineResp struct {
		Weakness      string `json:"weakness"`
		ShouldConcede bool   `json:"shouldConcede"`
		TooLong       bool   `json:"tooLong"`
		Reply         string `json:"reply"`
	}
	refineResp := ArgueRefineResp{}
	MustUnmarshal([]byte(Resp(llmRes0)), &refineResp)

	history = append(history, Msg{
		Role:    openai.ChatMessageRoleAssistant,
		Content: refineResp.Reply,
	})
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
					Content: refineResp.Reply,
				},
				FinishReason: "",
			},
		},
		Usage:               openai.Usage{},
		SystemFingerprint:   "",
		PromptFilterResults: nil,
	})
}
