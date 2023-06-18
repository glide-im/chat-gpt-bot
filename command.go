package main

import (
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/robotic"
	"github.com/google/uuid"
	"time"
)

type Role int
type User interface{}

const (
	RoleAdmin = 1
	RoleVip   = 2
	RoleNone  = 0
)

var roles = map[User]Role{}

type commands struct {
	bot *robotic.BotX
}

func (c *commands) initCommand() {
	command, _ := robotic.NewCommand(RoleNone, "login", c.handleCommandLogin)

	_ = c.bot.AddCommand(command)
}

func (c *commands) handleCommandLogin(message *robotic.ResolvedChatMessage, value string) error {

	if value == config.Common.AdminPassword {
		roles[message.ChatMessage.From] = RoleAdmin
		_ = c.bot.Send(message.ChatMessage.From, robotic.ActionChatMessage, c.createReply(message.ChatMessage.From, "管理员登录成功"))
		return nil
	}
	if value == config.Common.VipPassword {
		roles[message.ChatMessage.From] = RoleVip
		_ = c.bot.Send(message.ChatMessage.From, robotic.ActionChatMessage, c.createReply(message.ChatMessage.From, "VIP 登录成功"))
		return nil
	}
	_ = c.bot.Send(message.ChatMessage.From, robotic.ActionChatMessage, c.createReply(message.ChatMessage.From, "登录失败"))
	return nil
}

func (c *commands) handleCommandRestart(message *messages.ChatMessage, value string) error {

	return nil
}

func (c *commands) handleCommandStop(message *messages.ChatMessage, value string) error {

	return nil
}

func (c *commands) handleCommandStatus(message *messages.ChatMessage, value string) error {

	return nil
}

func (c *commands) createReply(to, content string) *messages.ChatMessage {
	return &messages.ChatMessage{
		CliMid:  uuid.New().String(),
		Mid:     0,
		From:    c.bot.Id,
		To:      to,
		Type:    1,
		Content: content,
		SendAt:  time.Now().Unix(),
	}
}
