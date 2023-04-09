package main

import (
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/robotic"
	"github.com/google/uuid"
	"time"
)

const (
	StreamMarkdown = 1011
	StreamText     = 1001

	StatusStart    = 1
	StatusSending  = 2
	StatusFinish   = 3
	StatusCanceled = 4
)

type ChatMessage2 struct {
	*messages.ChatMessage
	Status int32 `json:"status,omitempty"`
}

func (g *GptBot) handleStream(cm *messages.ChatMessage) {

	go func() {

		to := cm.From

		id, _ := uuid.NewUUID()
		sendAt := time.Now().UnixMilli()

		g.sendStreamMessage(to, id, sendAt, StatusStart, 0, "")
		time.Sleep(time.Second * 1)

		ch, err := g.Gpt.TextCompletionSteam(cm.Content, cm.From)
		if err != nil || ch == nil {
			g.sendStreamMessage(to, id, sendAt, StatusCanceled, 0, "机器人出错了")
			logger.ErrE("robot error", err)
			return
		}

		seq := 0
		for s := range ch {
			seq++
			g.sendStreamMessage(to, id, sendAt, StatusSending, int64(seq), s)
		}

		time.Sleep(time.Second * 1)
		g.sendStreamMessage(to, id, sendAt, StatusFinish, 0, "")

	}()
}

func (g *GptBot) sendStreamMessage(to string, id uuid.UUID, sendAt int64, status int32, seq int64, content string) {
	m := ChatMessage2{
		Status: status,
		ChatMessage: &messages.ChatMessage{
			Type:    StreamMarkdown,
			Mid:     int64(id.ID()),
			CliMid:  id.String(),
			Content: content,
			From:    g.BotX.Id,
			To:      to,
			SendAt:  sendAt,
			Seq:     seq,
		},
	}
	_ = g.BotX.Send(to, robotic.ActionClientCustom, &m)
}
