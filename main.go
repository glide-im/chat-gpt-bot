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

func main() {

	// 机器人的名字
	robotName := "chat_gpt"
	// 机器人 token
	token := ""
	chat_gpt.ApiToken("")
	// 设置代理
	chat_gpt.SetProxy("http://127.0.0.1:7890")
	// 使用这个服务器, 在这 http://im.dengzii.com/#/im/session/the_world_channel 可以看到机器人
	botX := robotic.NewBotX("ws://intercom.ink/ws", token)

	// 处理聊天消息
	botX.HandleChatMessage(func(m *messages.GlideMessage, cm *messages.ChatMessage) {
		logger.I("handler chat message >> %s", m.GetAction())
		if m.GetAction() == robotic.ActionChatMessage {
			go func() {
				reply, err := chat_gpt.Chat(cm.Content)
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
			logger.I("Receive Group Message: %s", m.To)
			if strings.HasPrefix(cm.Content, "@"+robotName) {

				go func() {
					reply, err := chat_gpt.Chat(cm.Content)
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
					err2 := botX.Send(m.To, robotic.ActionGroupMessage, &replyMsg)
					if err2 != nil {
						logger.ErrE("send error", err2)
					}

				}()
			}
		}
	})

	// 启动
	err := botX.Start(func(m *messages.GlideMessage) {
		// 处理所有消息
	})
	panic(err)
}
