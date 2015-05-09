package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/url"

	"github.com/astaxie/beego"
	"github.com/kr/pretty"
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
	Status    string        `json:"status"`
	Responses []interface{} `json:"responses"`
}

func (this *Chat) ChatWithBot() error {
	r := &Reply{}
	r.Key = "shoes"
	// TODO call bot
	answer := []byte(`{ "status": "ok", "responses": [{"type": "top","key":"shoes","value": "If you are not 100% satisfied with your purchase, you can return your item to us for a full refund. Returns must be done within 30 days of receipt together with the Returns slip at a SingPost's counter or POPStation, un-used with tags on, in original packaging and must not fall under the list of non-refundable brands/items <a href=\"http://www.zalora.sg/faq-non-refundable/\">HERE</a>.\n\nFor more information regarding the Return policy, please view the steps <a href=\"http://www.zalora.sg/faq-returns/\">HERE</a>"}], "sessionid": 50510 }`)
	chatBot := &ChatBot{}
	err := json.Unmarshal(answer, chatBot)
	if err != nil {
		beego.Alert(err)
		return errors.New("Bot didn't reply: " + err.Error())
	}
	reply := &Reply{}
	if len(chatBot.Responses) == 1 {
		replyInterface := chatBot.Responses[0]
		if replyMap, ok := replyInterface.(map[string]interface{}); ok {
			reply.Type = replyMap["type"].(string)
			reply.Value = replyMap["value"].(string)
		} else if replyStr, ok := replyInterface.(string); ok {
			reply.Value = replyStr
		}
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
	case "top":
		this.Value = getTop(this)
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

type Product struct {
	ConfigSku    string `json:"config_sku"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Url          string `json:"url"`
	Price        string `json:"price"`
	MainImageUrl string `json:"main_image_url"`
	SpecialPrice string `json:"special_price"`
}

func getTop(r *Reply) string {
	var urlObj *url.URL
	topUrl := fmt.Sprintf("https://api.zalora.sg/v1/products/?limit=3&query=%s", r.Key)
	urlObj, err := url.Parse(topUrl)
	if err != nil {
	}
	parameters := url.Values{}
	parameters.Add("query", r.Key)
	urlObj.RawQuery = parameters.Encode()

	res, err := GetThrift(urlObj.String())
	if err != nil {
		beego.Alert(err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(res, &result)
	if err != nil {

	}
	var pList = []*Product{}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if products, ok := data["products"].([]interface{}); ok {
			for _, pM := range products {
				pretty.Println(pM)
				if _, ok := pM.(map[string]interface{}); ok {
					b, err := json.Marshal(pM)
					if err != nil {
						continue
					}
					p := Product{}
					json.Unmarshal(b, &p)
					pList = append(pList, &p)
				}

			}
		}
	}
	pretty.Print(pList)
	str, err := json.Marshal(pList)
	if err != nil {
		return ""
	}

	return string(str)
}
