package main

import (
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/robotic"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type CliCustom struct {
	Type    int    `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
	Id      string `json:"id,omitempty"`
	From    string `json:"from,omitempty"`
}

const (
	Steam         = 100
	SteamFinish   = 101
	SteamCanceled = 102
)

func (h *MsgHandler) handleStream(cm *messages.ChatMessage) {

	go func() {
		newUUID, _ := uuid.NewUUID()
		id := newUUID.String()
		for i := 0; i < 20; i++ {
			_ = h.bot.Send(cm.From, robotic.ActionClientCustom, CliCustom{
				Type:    Steam,
				Id:      id,
				Content: strconv.Itoa(i),
				From:    cm.To,
			})
			time.Sleep(time.Millisecond * 300)
		}
	}()
}
