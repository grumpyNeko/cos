package main

import (
	openai "github.com/sashabaranov/go-openai"
	"unicode"
)

func containsChineseChar(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

const SysPromptToChs = `你是一位专业翻译, 无论用户输入什么, 都帮我译为中文
直接返回如下格式的JSON(不含markdown语法), literal表示直译, free表示意译:
{
 "literal": "..",
 "free": ".."
}`
const SysPromptToEng = `你是一位专业翻译, 无论用户输入什么, 都帮我译为英文
直接返回如下格式的JSON(不含markdown语法), literal表示直译, free表示意译:
{
 "literal": "..",
 "free": ".."
}`
const ExamplePromptToEng0 = `我对官僚主义过敏`
const ExamplePromptToEng1 = `{
  "literal": "I am allergic to bureaucratism",
  "free": "I can't stand bureaucracy"
}`

func Translate(s string) (messages []openai.ChatCompletionMessage) {
	if containsChineseChar(s) {
		messages = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: SysPromptToEng,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: ExamplePromptToEng0,
			},
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: ExamplePromptToEng1,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: s,
			},
		}
	} else {
		messages = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: SysPromptToChs,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: s,
			},
		}
	}
	return messages
}
