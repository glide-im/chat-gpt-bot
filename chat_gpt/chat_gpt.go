package chat_gpt

import "encoding/json"

func Chat(msg string) (string, error) {

	resp := &ChatGPTResponse{}
	reply, err := requestChatGPT(msg)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(reply, &resp)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}
