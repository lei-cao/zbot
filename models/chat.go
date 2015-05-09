package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

const ALICE_QUERY_URL = "http://www.zalora.sg/catalog/?q="

type Chat struct {
	Msg   string `json:"msg"`
	Reply *Reply `json:"reply"`
}

func (this *Chat) ChatWithBot() error {
	// TODO call bot
	answer := []byte(`{"type":"q","key":"shoes"}`)
	reply := &Reply{}
	err := json.Unmarshal(answer, reply)
	if err != nil {
		return errors.New("Bot didn't reply")
	}
	reply.SetValue()
	this.Reply = reply
	return nil
}

type Reply struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (this *Reply) SetValue() {
	if this.Type == "q" {
		this.Value = fmt.Sprintf(`See all <a href="%s">%s</a>`, ALICE_QUERY_URL+this.Key, this.Key)
	}
}

func Say(msg string) (*Chat, error) {
	chat := &Chat{Msg: msg}
	err := chat.ChatWithBot()
	if err != nil {
		return nil, err
	}
	return chat, nil
}
