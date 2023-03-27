package main

import (
	"github.com/glide-im/chat-gpt-bot/chat_gpt"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/robotic"
	"github.com/spf13/viper"
)

type Config struct {
	BotName      string
	Greetings    string
	BotToken     string
	OpenAiApiKey string
	Proxy        string
	BotServer    string
	Type         int32
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

var botX *robotic.BotX

func main() {

	chat_gpt.ApiToken(config.OpenAiApiKey)
	if config.Proxy != "" {
		chat_gpt.SetProxy(config.Proxy)
	}
	botX = robotic.NewBotX(config.BotServer, config.BotToken)

	// 处理聊天消息
	botX.HandleChatMessage(MessageHandler)

	// 启动
	err := botX.Start(func(m *messages.GlideMessage) {
		// 处理所有消息
	})
	panic(err)
}
