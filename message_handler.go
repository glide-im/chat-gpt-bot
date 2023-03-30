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
	config *Bot
	bot    *robotic.BotX
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
			var replyType int32 = 11
			if h.config.Type == 2 {
				replyType = 2
				reply, err = chat_gpt.ImageGen(cm.Content)
				if err != nil {
					logger.ErrE("robot error", err)
					if strings.Contains(err.Error(), "status code: 400") {
						reply = "**请不要试图生成包含非法的内容的图片**"
					} else {
						reply = "机器人出错啦"
					}
					replyType = 11
				}
			} else if h.config.Type == 1 {
				reply, err = chat_gpt.TextCompletion(cm.Content, cm.From)
				if err != nil {
					reply = "机器人出错啦"
					logger.ErrE("robot error", err)
				}
			}
			replyMsg := messages.ChatMessage{
				CliMid:  uuid.New().String(),
				Mid:     0,
				From:    h.bot.Id,
				To:      cm.From,
				Type:    replyType,
				Content: reply,
				SendAt:  time.Now().Unix(),
			}
			err2 := h.bot.Send(cm.From, robotic.ActionChatMessage, &replyMsg)
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
	if cm.Type == 100 && cm.Content != h.bot.Id {
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
			if h.config.Type == 2 {
				msgType = 11
			}
			replyMsg := messages.ChatMessage{
				CliMid:  uuid.New().String(),
				From:    h.bot.Id,
				To:      cm.To,
				Type:    int32(msgType),
				Content: fmt.Sprintf("@%s %s", cm.From, reply),
				SendAt:  time.Now().Unix(),
			}
			err2 := h.bot.Send(gid, robotic.ActionGroupMessage, &replyMsg)
			if err2 != nil {
				logger.ErrE("send error", err2)
			}

		}()
	}
}

func (h *MsgHandler) greetingTo(uid string) {

	greeting := messages.ChatMessage{
		CliMid:  uuid.New().String(),
		From:    h.bot.Id,
		To:      uid,
		Type:    11,
		Content: h.config.Greetings,
		SendAt:  time.Now().Unix(),
	}

	err := h.bot.Send(uid, robotic.ActionChatMessage, &greeting)
	if err != nil {
		logger.E("%v", err)
	}
}
