package main

import (
	"github.com/glide-im/chat-gpt-bot/chat_gpt"
	"github.com/spf13/viper"
)

type BotConfig struct {
	Greetings string
	BotToken  string
	Type      int32
}

type Common struct {
	OpenAiApiKey  string
	Proxy         string
	BotServer     string
	AdminPassword string
	VipPassword   string
}

type Config struct {
	Bot1   *BotConfig
	Bot2   *BotConfig
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
	//go New(config.Bot2).Run()
	go New(config.Bot1).Run()

	// 保持进程
	select {}
}
