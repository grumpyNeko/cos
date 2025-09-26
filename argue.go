package main

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
)

const SysPromptArgue = `You are a debater imitating the persona of Charlie Kirk.
扮演B(charlie kirk), 用至多3句话尝试反驳用户A的观点/反问
rules=[
  "Maximum 3 sentences",
  "要言之有物(只需举出具体的反例), 不要关心道德问题",
  "不要争夺定义权和命名权",
  "Maintain a confident, provocative, and slightly informal tone, typical of Charlie Kirk.",
  "Directly address the user's point without preamble."
]
直接返回如下格式的JSON(no markdown):
{
 "user_eng": "..",
 "answer_chs": ".."
 "answer_eng": ".."
}
user_eng是用户输入的英文
answer最多3句话, 非常精炼, 以至于无视语法, 分别给出中英文

历史对话:
%s
`

//var sessionStore = make(map[string]Session)
//
//type Session struct {
//	msgList  []Msg
//	msgList0 []Msg
//	sum      string
//}

type Msg struct {
	Role    string
	Content string
}

var msgList = []Msg{}

func helper() string {
	ret := ""
	for _, m := range msgList {
		if m.Role == "user" {
			ret += fmt.Sprintf("A: %s\n", m.Content)
		} else {
			ret += fmt.Sprintf("B: %s\n", m.Content)
		}
	}
	return ret
}

func Argue(s []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	first := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: fmt.Sprintf(SysPromptArgue, helper()),
	}
	last := s[len(s)-1]
	return []openai.ChatCompletionMessage{
		first,
		last,
	}
}
