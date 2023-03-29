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

type MsgHandler struct {
	bot *Bot
}

func (h *MsgHandler) MessageHandler(m *messages.GlideMessage, cm *messages.ChatMessage) {

	logger.I("handler chat message >> %s", m.GetAction())

	if cm.From == "100000" || cm.From == "543852" {
		return
	}
	if m.GetAction() == robotic.ActionChatMessage {
		go func() {
			var reply string
			var err error
			if h.bot.Type == 2 {
				reply, err = chat_gpt.ImageGen(cm.Content)
			} else if h.bot.Type == 1 {
				reply, err = chat_gpt.TextCompletion(cm.Content, cm.From)
			}
			if err != nil {
				reply = "机器人出错啦"
				logger.ErrE("robot error", err)
			}
			replyMsg := messages.ChatMessage{
				CliMid:  uuid.New().String(),
				Mid:     0,
				From:    botX.Id,
				To:      cm.From,
				Type:    11,
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
		h.handleGroupMessage(m.To, cm)
	}
}

func (h *MsgHandler) handleGroupMessage(gid string, cm *messages.ChatMessage) {
	if cm.Type == 100 && cm.Content != botX.Id {
		go h.greetingTo(cm.Content)
	}
	logger.I("Receive Group Message: %s", gid)
	if strings.HasPrefix(cm.Content, "@openai ") {

		go func() {
			reply, err := chat_gpt.TextCompletion(cm.Content, cm.From)
			if err != nil {
				reply = "机器人出错啦"
				logger.ErrE("robot error", err)
			}

			msgType := 1
			if h.bot.Type == 2 {
				msgType = 11
			}
			replyMsg := messages.ChatMessage{
				CliMid:  uuid.New().String(),
				From:    botX.Id,
				To:      cm.To,
				Type:    int32(msgType),
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

func (h *MsgHandler) greetingTo(uid string) {

	greeting := messages.ChatMessage{
		CliMid:  uuid.New().String(),
		From:    botX.Id,
		To:      uid,
		Type:    11,
		Content: h.bot.Greetings,
		SendAt:  time.Now().Unix(),
	}

	err := botX.Send(uid, robotic.ActionChatMessage, &greeting)
	if err != nil {
		logger.E("%v", err)
	}
}
