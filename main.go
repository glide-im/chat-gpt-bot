package main

import (
	"github.com/glide-im/chat-gpt-bot/chat_gpt"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/robotic"
	"github.com/spf13/viper"
)

type Bot struct {
	Greetings string
	BotToken  string
	Type      int32
}

type Common struct {
	OpenAiApiKey string
	Proxy        string
	BotServer    string
}

type Config struct {
	Bot1   *Bot
	Bot2   *Bot
	Common *Common
}

func init() {

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	config = &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		panic(err)
	}
}

var config *Config

func main() {
	chat_gpt.ApiToken(config.Common.OpenAiApiKey)
	if config.Common.Proxy != "" {
		chat_gpt.SetProxy(config.Common.Proxy)
	}
	go startBot(config.Bot2)
	startBot(config.Bot1)
}

func startBot(bot *Bot) {

	var botX *robotic.BotX

	botX = robotic.NewBotX(config.Common.BotServer, bot.BotToken)
	// 处理聊天消息
	h := &MsgHandler{bot: botX, config: bot}
	botX.HandleChatMessage(h.MessageHandler)

	// 启动
	err := botX.Start(func(m *messages.GlideMessage) {
		// 处理所有消息
	})
	panic(err)
}
