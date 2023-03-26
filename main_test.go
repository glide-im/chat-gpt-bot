package main

import (
	"context"
	"encoding/json"
	"github.com/glide-im/chat-gpt-bot/chat_gpt"
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func TestGPT2(t *testing.T) {
	client := openai.NewClient("")
	request := openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo,
		Messages:         []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: "你好"}},
		MaxTokens:        500,
		Temperature:      0.5,
		TopP:             0.5,
		N:                1,
		Stream:           false,
		Stop:             nil,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.1,
		LogitBias:        nil,
		User:             "",
	}
	response, err := client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(&response)
	t.Log(string(marshal))
}

func TestGPT(t *testing.T) {

	chat_gpt.ApiToken("")
	chat_gpt.SetProxy("http://127.0.0.1:10808/")
	chat, err := chat_gpt.TextCompletion("System.out.", "")
	if err != nil {
		t.Error(err)
	}
	t.Log(chat)
}

func TestGenerateSecret(t *testing.T) {
	jwtAuthorize := jwt_auth.NewAuthorizeImpl("")
	token, _ := jwtAuthorize.GetToken(&jwt_auth.JwtAuthInfo{
		UID:         "",
		Device:      "1",
		ExpiredHour: 23,
	})
	t.Log(token.Token)
}
