package chat_gpt

import "testing"

func TestTextCompletion(t *testing.T) {
	re, err := TextCompletion("Hello", "1")
	if err != nil {
		t.Error(err)
	}
	t.Log(re)
}

func TestSteam(t *testing.T) {
	ApiToken("sk-")
	SetProxy("http://127.0.0.1:7890/")
	steam, err := TextCompletionSteam("Hello", "1")
	if err != nil {
		t.Error(err)
	}
	for s := range steam {
		t.Log(s)
	}
}
