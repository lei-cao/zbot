package models

import (
	"encoding/json"
	"errors"
	"fmt"
    "github.com/astaxie/beego"
)

const ALICE_QUERY_URL = "http://www.zalora.sg/catalog/?q="

type Chat struct {
	Msg   string `json:"msg"`
	Reply *Reply `json:"reply"`
}

func (this *Chat) BotInfo() {
	body, err := GetBot()
	if err != nil {

	}
	fmt.Println(string(body))
}

func (this *Chat) TestChatWithBot() {
	SendChat(this.Msg)
}

type ChatBot struct {
	Status    string   `json:"status"`
	Responses []*Reply `json:"responses"`
}

func (this *Chat) ChatWithBot() error {
	// TODO call bot
	answer := []byte(`{ "status": "ok", "responses": [{"type": "raw","value": "If you are not 100% satisfied with your purchase, you can return your item to us for a full refund. Returns must be done within 30 days of receipt together with the Returns slip at a SingPost's counter or POPStation, un-used with tags on, in original packaging and must not fall under the list of non-refundable brands/items <a href=\"http://www.zalora.sg/faq-non-refundable/\">HERE</a>.\n\nFor more information regarding the Return policy, please view the steps <a href=\"http://www.zalora.sg/faq-returns/\">HERE</a>"}], "sessionid": 50510 }`)
	chatBot := &ChatBot{}
	reply := &Reply{}
	err := json.Unmarshal(answer, chatBot)
	if err != nil {
        beego.Alert(err)
		return errors.New("Bot didn't reply: " + err.Error())
	}
	if len(chatBot.Responses) == 1 {
		reply = chatBot.Responses[0]
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
	switch this.Type {
	case "q":
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
