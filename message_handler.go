package main

import (
	"fmt"
	"github.com/glide-im/chat-gpt-bot/chat_gpt"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/robotic"
	"github.com/google/uuid"
	"strings"
	"time"
)

const greetings = `
您好，我是 ChatGPT Bot，是一个由 OpenAI GPT-3.5 开发的智能对话机器人。我可以回答各种问题，例如科学、技术、历史、文化等方面的问题, 也可以模拟面试, 文字互动游戏等。
`

func MessageHandler(m *messages.GlideMessage, cm *messages.ChatMessage) {

	logger.I("handler chat message >> %s", m.GetAction())

	if m.GetAction() == robotic.ActionChatMessage {
		go func() {
			reply, err := chat_gpt.TextCompletion(cm.Content, cm.From)
			if err != nil {
				reply = "机器人出错啦"
				logger.ErrE("robot error", err)
			}
			replyMsg := messages.ChatMessage{
				CliMid:  uuid.New().String(),
				Mid:     0,
				From:    botX.Id,
				To:      cm.From,
				Type:    cm.Type,
				Content: reply,
				SendAt:  time.Now().Unix(),
			}
			err2 := botX.Send(cm.From, robotic.ActionChatMessage, &replyMsg)
			if err2 != nil {
				logger.ErrE("send error", err2)
			}
		}()
	}
	if m.GetAction() == robotic.ActionGroupMessage {
		handleGroupMessage(m.To, cm)
	}
}

func handleGroupMessage(gid string, cm *messages.ChatMessage) {
	if cm.Type == 100 && cm.Content != botX.Id {
		go greetingTo(cm.Content)
	}
	logger.I("Receive Group Message: %s", gid)
	if strings.HasPrefix(cm.Content, "@"+config.BotName) {

		go func() {
			reply, err := chat_gpt.TextCompletion(cm.Content, cm.From)
			if err != nil {
				reply = "机器人出错啦"
				logger.ErrE("robot error", err)
			}

			replyMsg := messages.ChatMessage{
				CliMid:  uuid.New().String(),
				From:    botX.Id,
				To:      cm.To,
				Type:    cm.Type,
				Content: fmt.Sprintf("@%s %s", cm.From, reply),
				SendAt:  time.Now().Unix(),
			}
			err2 := botX.Send(gid, robotic.ActionGroupMessage, &replyMsg)
			if err2 != nil {
				logger.ErrE("send error", err2)
			}

		}()
	}
}

func greetingTo(uid string) {

	greeting := messages.ChatMessage{
		CliMid:  uuid.New().String(),
		From:    botX.Id,
		To:      uid,
		Type:    1,
		Content: greetings,
		SendAt:  time.Now().Unix(),
	}

	err := botX.Send(uid, robotic.ActionChatMessage, &greeting)
	if err != nil {
		logger.E("%v", err)
	}
}
