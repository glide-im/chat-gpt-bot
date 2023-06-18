package main

import (
	"github.com/glide-im/chat-gpt-bot/openai"
	"github.com/glide-im/robotic"
	"github.com/spf13/viper"
)

type BotConfig struct {
	Greetings string
	Email     string
	Password  string
	Type      int32
}

type Common struct {
	OpenAiApiKey  string
	Proxy         string
	BotServer     string
	ApiBaseUrl    string
	AdminPassword string
	VipPassword   string
}

type Config struct {
	Bot    *BotConfig
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
	gpt := openai.New(config.Common.OpenAiApiKey, config.Common.Proxy)

	robotic.ApiBaseUrl = config.Common.ApiBaseUrl
	go New(config.Bot, gpt).Run()

	select {}
}
