package main

import (
	"fmt"
	"github.com/glide-im/chat-gpt-bot/openai"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/robotic"
	"github.com/google/uuid"
	"strings"
	"time"
)

type GptBot struct {
	Commands *commands
	Config   *BotConfig
	BotX     *robotic.BotX

	Gpt *openai.ChatGPT

	server string
}

func New(c *BotConfig, gpt *openai.ChatGPT) *GptBot {

	var botX *robotic.BotX
	botX = robotic.NewBotX(config.Common.BotServer)

	return &GptBot{
		Gpt:      gpt,
		Commands: &commands{bot: botX},
		Config:   c,
		BotX:     botX,
		server:   config.Common.BotServer,
	}
}

func (g *GptBot) Run() {

	func() {
		err := recover()
		if err != nil {
			println(err.(error).Error())
			go g.Run()
		}
	}()

	g.Commands.initCommand()

	// 处理聊天消息
	g.BotX.HandleChatMessage(g.MessageHandler)

	// 启动
	err := g.BotX.RunAndLogin(g.Config.Email, g.Config.Password, func(m *messages.GlideMessage) {
		// 处理所有消息
	})
	panic(err)
}

func (g *GptBot) MessageHandler(m *messages.GlideMessage, cm *messages.ChatMessage) {

	logger.I("handler chat message >> %s", m.GetAction())

	if cm.From == "100000" || cm.From == "543852" {
		return
	}
	if m.GetAction() == robotic.ActionChatMessage {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					err2, ok := err.(error)
					if ok {
						logger.ErrE("robot error", err2)
					} else {
						logger.ErrE("robot error", fmt.Errorf("%v", err))
					}
				}
			}()

			if g.Config.Type == 2 {
				g.imageGen(cm)
			} else if g.Config.Type == 1 {
				g.handleStream(cm)
			}
		}()
	}
	if m.GetAction() == robotic.ActionGroupMessage {
		g.handleGroupMessage(m.To, cm)
	}
}

func (g *GptBot) imageGen(cm *messages.ChatMessage) {

	var reply = ""
	var err error
	var replyType int32 = 2

	reply, err = g.Gpt.ImageGen(cm.Content)

	if err != nil {
		logger.ErrE("robot error", err)
		if strings.Contains(err.Error(), "status code: 400") {
			reply = "**请不要试图生成包含非法的内容的图片**"
		} else {
			reply = "机器人出错啦"
		}
		replyType = 11
	}
	replyMsg := messages.ChatMessage{
		CliMid:  uuid.New().String(),
		Mid:     0,
		From:    g.BotX.Id,
		To:      cm.From,
		Type:    replyType,
		Content: reply,
		SendAt:  time.Now().Unix(),
	}
	err2 := g.BotX.Send(cm.From, robotic.ActionChatMessage, &replyMsg)
	if err2 != nil {
		logger.ErrE("send error", err2)
	}
}

func (g *GptBot) handleGroupMessage(gid string, cm *messages.ChatMessage) {
	if cm.Type == 100 && cm.Content != g.BotX.Id {
		go g.greetingTo(cm.Content)
	}
	logger.I("Receive Group Message: %s", gid)
	if strings.HasPrefix(cm.Content, "@openai ") {

		go func() {
			reply, err := g.Gpt.TextCompletion(cm.Content, cm.From)
			if err != nil {
				reply = "机器人出错啦"
				logger.ErrE("robot error", err)
			}

			msgType := 1
			if g.Config.Type == 2 {
				msgType = 11
			}
			replyMsg := messages.ChatMessage{
				CliMid:  uuid.New().String(),
				From:    g.BotX.Id,
				To:      cm.To,
				Type:    int32(msgType),
				Content: fmt.Sprintf("@%s %s", cm.From, reply),
				SendAt:  time.Now().Unix(),
			}
			err2 := g.BotX.Send(gid, robotic.ActionGroupMessage, &replyMsg)
			if err2 != nil {
				logger.ErrE("send error", err2)
			}

		}()
	}
}

func (g *GptBot) greetingTo(uid string) {

	greeting := messages.ChatMessage{
		CliMid:  uuid.New().String(),
		From:    g.BotX.Id,
		To:      uid,
		Mid:     time.Now().Unix(),
		Type:    11,
		Content: g.Config.Greetings,
		SendAt:  time.Now().Unix(),
	}

	err := g.BotX.Send(uid, robotic.ActionChatMessage, &greeting)
	if err != nil {
		logger.E("%v", err)
	}
}
