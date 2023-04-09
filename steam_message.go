package main

import (
	"github.com/glide-im/chat-gpt-bot/chat_gpt"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/robotic"
	"github.com/google/uuid"
	"time"
)

const (
	Steam          = 1000
	StreamMarkdown = 1011
	SteamFinish    = 1001
	SteamCanceled  = 1002
)

func (g *GptBot) handleStream(cm *messages.ChatMessage) {

	go func() {

		newUUID, _ := uuid.NewUUID()
		id := newUUID.ID()
		seq := 0
		m := messages.ChatMessage{
			Type:    Steam,
			Mid:     int64(id),
			CliMid:  newUUID.String(),
			Content: "",
			From:    cm.To,
			To:      cm.From,
			SendAt:  time.Now().Unix(),
			Seq:     int64(seq),
		}

		ch, err := chat_gpt.TextCompletionSteam(cm.Content, cm.From)
		if err != nil || ch == nil {
			m.Type = 11
			m.Content = "机器人出错啦"
			_ = g.BotX.Send(cm.From, robotic.ActionChatMessage, &m)
			logger.ErrE("robot error", err)
			return
		}
		for s := range ch {
			if s != "" {
				m.Content = s
				m.Seq = int64(seq)
				seq++
				_ = g.BotX.Send(cm.From, robotic.ActionClientCustom, &m)
				time.Sleep(time.Millisecond * 30)
			}
		}
		m.Content = "."
		m.Type = SteamFinish
		_ = g.BotX.Send(cm.From, robotic.ActionClientCustom, &m)
	}()
}
