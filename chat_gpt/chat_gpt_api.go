package chat_gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var _token = ""
var _client = &http.Client{}

func ApiToken(token string) {
	_token = token
}

func SetProxy(httpProxy string) {
	url_, err := url.Parse(httpProxy)
	if err != nil {
		panic(err)
	}
	_client.Transport = &http.Transport{
		Proxy: http.ProxyURL(url_),
	}
}

func requestChatGPT(msg string) ([]byte, error) {

	marshal, _ := json.Marshal(&ChatGPTRequest{
		Prompt:    msg,
		MaxTokens: 2048,
		Model:     "text-davinci-003",
	})

	reader := bytes.NewReader(marshal)

	request, err2 := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", reader)
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
