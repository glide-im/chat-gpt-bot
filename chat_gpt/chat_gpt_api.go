package chat_gpt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/glide-im/glide/pkg/logger"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/sashabaranov/go-openai"
	"io/ioutil"
	"net/http"
	"net/url"
)

var _token = ""
var _client = &http.Client{}
var _openAi *openai.Client
var _cache *lru.Cache[string, *MessageList]

type MessageList struct {
	m []openai.ChatCompletionMessage
}

func init() {
	var err error
	_cache, err = lru.New[string, *MessageList](200)
	if err != nil {
		panic(err)
	}
}

func ApiToken(token string) {
	_token = token

	_openAi = openai.NewClient(_token)
}

func SetProxy(httpProxy string) {
	url_, err := url.Parse(httpProxy)
	if err != nil {
		panic(err)
	}
	_client.Transport = &http.Transport{
		Proxy: http.ProxyURL(url_),
	}
	config := openai.DefaultConfig(_token)
	config.HTTPClient = _client
	_openAi = openai.NewClientWithConfig(config)
}

func textCompletion(param *ChatGPTRequest) ([]byte, error) {

	marshal, _ := json.Marshal(param)

	reader := bytes.NewReader(marshal)

	request, err2 := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", reader)
	if err2 != nil {
		return nil, err2
	}

	request.Header.Add("Authorization", "Bearer "+_token)
	request.Header.Add("Content-Type", "application/json")

	response, err := _client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200 {
		return ioutil.ReadAll(response.Body)

	}
	return nil, errors.New(response.Status)
}

func TextCompletion(msg string, userId string) (string, error) {

	m := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: msg}
	history := loadHistory(userId)
	logger.D("load %d history for user %s", len(history), userId)
	history = append(history, m)
	request := openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo,
		Messages:         history,
		MaxTokens:        500,
		Temperature:      0.7,
		TopP:             0.5,
		N:                1,
		Stream:           false,
		Stop:             nil,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.1,
		LogitBias:        nil,
		User:             userId,
	}
	response, err := _openAi.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return "", err
	}
	choice := response.Choices[0]
	content := choice.Message.Content

	updateChatHistory(userId, m, choice.Message)

	return content, nil
}

func updateChatHistory(id string, user openai.ChatCompletionMessage, bot openai.ChatCompletionMessage) {

	value, ok := _cache.Get(id)
	if !ok {
		_cache.Add(id, &MessageList{m: []openai.ChatCompletionMessage{
			user,
			bot,
		}})
	} else {
		value.m = append(value.m, user)
		value.m = append(value.m, bot)
	}
}

func loadHistory(id string) []openai.ChatCompletionMessage {
	value, ok := _cache.Get(id)
	var result []openai.ChatCompletionMessage
	if ok {
		result = make([]openai.ChatCompletionMessage, len(value.m))
		copy(result[:], value.m[:])
	}
	return result

}
