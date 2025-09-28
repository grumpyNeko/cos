package main

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type ArgueNamespaceType struct{}

var ArgueNamespace ArgueNamespaceType

const PromptArgue_GenReply = `
扮演charlie kirk
---
debate_rules=[
  "要言之有物(只需举出具体的反例), 不要关心道德问题",
  "不要争夺定义权和命名权",
  "Maintain a confident, provocative, and slightly informal tone, typical of Charlie Kirk.",
  "Directly address the user's point without preamble."
]

user_msg_eng = translate(user_msg)
debate_mode = 最新user_msg是否需要反驳
reply_len = 需回复几句话, 最多为3
if debate_mode
  reply_eng = 反驳或反问(debate_rules)
else
  reply_eng = 正常回复
reply_chs = translate(reply_eng)
---
历史对话=%s
注意事项=%s
---
直接返回如下格式的JSON(no markdown):
{
  "user_msg_eng": "..",
  "debate_mode": true,
  "reply_len": 2,
  "reply_eng": "..",
  "reply_chs": "..",
}
`

const PromptArgue_RefineReply0 = `
一个rag正在扮演charlie kirk, 请对优化rag生成的对话, 返回json

func refine(reply, history)
  weakness, shouldConcede = 找出Assistant论据中最大的隐患 // weakness是1句话, shouldConcede表示user的观点难以反驳
  reply0 = 根据weakness和shouldConcede修改reply
  reply1_eng = reply0口语化(模仿charlie尖锐的语气), 浓缩信息, 言简意赅, 根据情况在末尾增加反问
  reply1_chs = translate(reply1_eng)
  return {
    weakness, shouldConcede, reply0, reply1_eng, reply1_chs
  }

example:
"没错, 共同抚养确实不普遍, 我承认这一点。但看看更普遍的: 长期同居。他们一起租房、互为保险受益人。这种深度的财务捆绑和风险共担, 难道是普通朋友会做的吗？"
简化为
"好, 你说共同抚养不普遍, 那一起租房同居呢?"
`

const PromptArgue_RefineReply1 = `
reply=%s
expect_reply_len=%d
history=%s
注意事项=%s
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

var history = []Msg{
	//{
	//	Role:    openai.ChatMessageRoleUser,
	//	Content: "如果男女朋友没有性行为, 和普通朋友有什么本质区别",
	//},
	//{
	//	Role:    openai.ChatMessageRoleAssistant,
	//	Content: "承诺、排他性、共同规划未来，这些是普通朋友没有的；订婚情侣长期异地也可能零性行为，但关系本质不同。税务与法律层面也不同：伴侣常共同理财、共同租房、紧急联系人、医疗决策授权。看看数据：婚前禁欲的宗教情侣依然有明显更高的同居率、资源整合与家庭网络绑定，这不是普通友情。",
	//},
	//{
	//	Role:    openai.ChatMessageRoleUser,
	//	Content: "如果两人不结婚也没性关系, 请问两人是什么关系",
	//},
	//{
	//	Role:    openai.ChatMessageRoleAssistant,
	//	Content: "合伙创业的情侣、共同抚养孩子的伴侣、签署医疗与财务授权的伴侣，全部不是“普通朋友”。同居伴侣在多国法律下享有税务与继承安排，却可能选择禁欲",
	//},
	//{
	//	Role:    openai.ChatMessageRoleUser,
	//	Content: `所谓"合伙创业"没有说服力, "共同抚养"和"医疗与财务授权"并不普遍`,
	//},
	//{
	//	Role:    openai.ChatMessageRoleAssistant,
	//	Content: `不需要普遍才构成反例：长期同居的未婚伴侣在房贷联名、租约、紧急联系人、保险受益人上极其常见，已足以把他们与普通朋友区分。美国有数百万同居伴侣提交联合财务、共同抚养宠物与分工照护老人，这些资源整合与风险共担不是友情`,
	//},
}

var forceList = []string{
	//"只讨论不打算结婚的情况",
}

func HistoryToString() string {
	if len(history) == 0 {
		return "[]"
	}
	ret := "[\n"
	for _, m := range history {
		if m.Role == "user" {
			ret += fmt.Sprintf("user: %s\n", m.Content)
		} else {
			ret += fmt.Sprintf("charlie: %s\n", m.Content)
		}
	}
	ret += `]`
	return ret
}

func ForceListToString() string {
	if len(forceList) == 0 {
		return "[]"
	}
	ret := "[\n"
	for _, e := range forceList {
		ret += fmt.Sprintf("%s\n", e)
	}
	ret += `]`
	return ret
}

func ArgueGen(s []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	first := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: fmt.Sprintf(PromptArgue_GenReply, HistoryToString(), ForceListToString()),
	}
	last := s[len(s)-1]
	if last.Role != openai.ChatMessageRoleUser {
		panic(`last.Role != openai.ChatMessageRoleUser`)
	}
	history = append(history, Msg{
		Role:    last.Role,
		Content: last.Content,
	})
	return []openai.ChatCompletionMessage{
		first,
		last,
	}
}

func ArgueRefine(reply string, expect_reply_len int) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: PromptArgue_RefineReply0,
		},
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(PromptArgue_RefineReply1, reply, expect_reply_len, HistoryToString(), ForceListToString()),
		},
	}
}
