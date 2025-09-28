package main

import (
	"cos/conf"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"net/http"
	"os"
	"strings"
)

var apikey = os.Getenv("OPENAI_API_KEY")
var baseURL = os.Getenv("OPENAI_API_KEY")

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	gin.SetMode(gin.DebugMode)
	r := NewRouter()

	err = r.Run(conf.GlobalConfig.Server.Addr)
	if err != nil {
		panic(err)
	}
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
			last := arg.Messages[len(arg.Messages)-1]
			if last.Role != openai.ChatMessageRoleUser {
				panic(`last.Role != openai.ChatMessageRoleUser`)
			}
			llmRes := MustLLM(
				baseURL,
				apikey,
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

const drawGunCommand = "drawgun: "

func argueHandler(context *gin.Context) {
	var arg openai.ChatCompletionRequest
	if err := context.BindJSON(&arg); err != nil {
		panic(err)
	}
	s := ModelPayload{}
	MustUnmarshal([]byte(arg.Model), &s)
	if s.MaskName != "Charlie" {
		panic(fmt.Sprintf(`s.MaskName != "Charlie" // s.MaskName is %s`, s.MaskName))
	}

	type finalResp struct {
		//UserMsgEng string `json:"user_msg_eng"`
		ReplyEng string `json:"reply_eng"`
		ReplyChs string `json:"reply_chs"`
	}
	last := arg.Messages[len(arg.Messages)-1]
	if last.Role != openai.ChatMessageRoleUser {
		panic(`last.Role != openai.ChatMessageRoleUser`)
	}
	if strings.HasPrefix(last.Content, drawGunCommand) {
		forceList = append(forceList, strings.TrimPrefix(last.Content, drawGunCommand))
		context.JSON(http.StatusOK, openai.ChatCompletionResponse{
			ID:      "",
			Object:  "",
			Created: 0,
			Model:   "",
			Choices: []openai.ChatCompletionChoice{
				{
					Index: 0,
					Message: openai.ChatCompletionMessage{
						Role: "assistant",
						Content: string(MustMarshal(finalResp{
							ReplyEng: "ok",
							ReplyChs: "好好好",
						})),
					},
					FinishReason: "",
				},
			},
			Usage:               openai.Usage{},
			SystemFingerprint:   "",
			PromptFilterResults: nil,
		})
		return
	}

	history = append(history, Msg{
		Role:    openai.ChatMessageRoleAssistant,
		Content: last.Content,
	})
	defer println(fmt.Sprintf(" %s", HistoryToString()))
	defer println(fmt.Sprintf(" %s", ForceListToString()))
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
		baseURL,
		apikey,
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
		UserMsgEng string `json:"user_msg_eng"`
		DebateMode bool   `json:"debate_mode"`
		ReplyLen   int    `json:"reply_len"`
		ReplyEng   string `json:"reply_eng"`
		ReplyChs   string `json:"reply_chs"`
	}
	genResp := ArgueGenResp{}
	MustUnmarshal([]byte(Resp(llmRes)), &genResp)
	if !genResp.DebateMode {
		context.JSON(http.StatusOK, openai.ChatCompletionResponse{
			ID:      "",
			Object:  "",
			Created: 0,
			Model:   "",
			Choices: []openai.ChatCompletionChoice{
				{
					Index: 0,
					Message: openai.ChatCompletionMessage{
						Role: "assistant",
						Content: string(MustMarshal(finalResp{
							//UserMsgEng: genResp.UserMsgEng,
							ReplyEng: genResp.ReplyEng,
							ReplyChs: genResp.ReplyChs,
						})),
					},
					FinishReason: "",
				},
			},
			Usage:               openai.Usage{},
			SystemFingerprint:   "",
			PromptFilterResults: nil,
		})
		return
	}

	payload0 := ArgueRefine(genResp.ReplyEng, genResp.ReplyLen)
	println(fmt.Sprintf(" %+v", payload0))
	llmRes0 := MustLLM(
		baseURL,
		apikey,
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
		Reply0        string `json:"reply0"`
		Reply1Eng     string `json:"reply1_eng"`
		Reply1Chs     string `json:"reply1_chs"`
	}
	refineResp := ArgueRefineResp{}
	MustUnmarshal([]byte(Resp(llmRes0)), &refineResp)

	history = append(history, Msg{
		Role:    openai.ChatMessageRoleAssistant,
		Content: refineResp.Reply1Eng,
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
					Role: "assistant",
					Content: string(MustMarshal(finalResp{
						//UserMsgEng: genResp.UserMsgEng,
						ReplyEng: refineResp.Reply1Eng,
						ReplyChs: refineResp.Reply1Chs,
					})),
				},
				FinishReason: "",
			},
		},
		Usage:               openai.Usage{},
		SystemFingerprint:   "",
		PromptFilterResults: nil,
	})
}
