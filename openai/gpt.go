package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/glide-im/glide/pkg/logger"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/sashabaranov/go-openai"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type MessageList struct {
	totalToken int
	messages   []openai.ChatCompletionMessage
}

type ChatGPT struct {
	token  string
	proxy  string
	client *http.Client
	openAi *openai.Client
	cache  *lru.Cache[string, *MessageList]
}

func New(token string, proxy string) *ChatGPT {

	cache, err := lru.New[string, *MessageList](200)
	if err != nil {
		panic(err)
	}

	gpt := &ChatGPT{
		token:  token,
		proxy:  proxy,
		client: &http.Client{},
		openAi: openai.NewClient(token),
		cache:  cache,
	}

	if proxy != "" {
		gpt.SetProxy(proxy)
	}

	return gpt
}

func (c *ChatGPT) ApiToken(token string) {
	c.token = token
	c.openAi = openai.NewClient(c.token)
}

func (c *ChatGPT) SetProxy(httpProxy string) {
	url_, err := url.Parse(httpProxy)
	if err != nil {
		panic(err)
	}
	c.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(url_),
	}
	config := openai.DefaultConfig(c.token)
	config.HTTPClient = c.client
	c.openAi = openai.NewClientWithConfig(config)
}

func (c *ChatGPT) textCompletion(param *ChatGPTRequest) ([]byte, error) {

	marshal, _ := json.Marshal(param)

	reader := bytes.NewReader(marshal)

	request, err2 := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", reader)
	if err2 != nil {
		return nil, err2
	}

	request.Header.Add("Authorization", "Bearer "+c.token)
	request.Header.Add("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200 {
		return ioutil.ReadAll(response.Body)

	}
	return nil, errors.New(response.Status)
}

func (c *ChatGPT) ImageGen(prompt string) (string, error) {
	ctx := context.Background()

	// Sample image by link
	reqUrl := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize512x512,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	respUrl, err := c.openAi.CreateImage(ctx, reqUrl)
	if err != nil {
		return "", err
	}
	return respUrl.Data[0].URL, nil
}

func (c *ChatGPT) TextCompletion(msg string, userId string) (string, error) {

	m := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: msg}
	history := c.loadHistory(userId)
	logger.D("load %d history for user %s", len(history), userId)
	history = append(history, m)
	request := openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo0301,
		Messages:         history,
		MaxTokens:        2000,
		Temperature:      1.0,
		TopP:             0.5,
		N:                1,
		Stream:           false,
		Stop:             nil,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.1,
		LogitBias:        nil,
		User:             userId,
	}
	response, err := c.openAi.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return "", err
	}
	choice := response.Choices[0]
	content := choice.Message.Content

	c.updateChatHistory(userId, m, choice.Message)

	return content, nil
}

func (c *ChatGPT) TextCompletionSteam(msg string, userId string) (<-chan string, error) {

	m := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: msg}
	history := c.loadHistory(userId)
	logger.D("load %d history for user %s", len(history), userId)
	history = append(history, m)
	request := openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo0301,
		Messages:         history,
		MaxTokens:        2000,
		Temperature:      1.0,
		TopP:             0.5,
		N:                1,
		Stream:           true,
		Stop:             nil,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.1,
		LogitBias:        nil,
		User:             userId,
	}
	response, err := c.openAi.CreateChatCompletionStream(context.Background(), request)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.E("error in stream %v", err)
			}
		}()
		all := ""
		for {
			recv, err2 := response.Recv()
			if err2 != nil {
				if err2 != io.EOF {
					logger.E("error in stream %v", err2)
				}
				response.Close()
				close(ch)
				c.updateChatHistory(userId, m, openai.ChatCompletionMessage{
					Role:    "assistant",
					Content: all,
				})
				return
			}
			content := recv.Choices[0].Delta.Content
			ch <- content
			all += content
		}
	}()

	return ch, nil
}

func (c *ChatGPT) updateChatHistory(userId string, user openai.ChatCompletionMessage, bot openai.ChatCompletionMessage) {

	history, hasHistory := c.cache.Get(userId)
	if !hasHistory {
		history = &MessageList{messages: []openai.ChatCompletionMessage{}}
		c.cache.Add(userId, history)
	}

	history.totalToken += len(user.Content) + len(bot.Content)

	history.messages = append(history.messages, user)
	history.messages = append(history.messages, bot)

	for history.totalToken > 4000 {
		history.totalToken -= len(history.messages[0].Content)
		history.messages = history.messages[1:]
	}
}

func (c *ChatGPT) loadHistory(id string) []openai.ChatCompletionMessage {
	value, ok := c.cache.Get(id)
	var result []openai.ChatCompletionMessage
	if ok {
		result = make([]openai.ChatCompletionMessage, len(value.messages))
		copy(result[:], value.messages[:])
	}
	return result

}
