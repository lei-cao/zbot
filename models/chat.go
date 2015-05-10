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

func (this *Chat) testingChat() []byte {
	switch this.Reply.Type {
	case "q":
		return []byte(`{ "status": "ok", "responses": [{"type": "q","value": "shoes"}], "sessionid": 50510 }`)
	case "top":
		return []byte(`{ "status": "ok", "responses": ["{\n      \"type\": \"top\",\n      \"key\": \"nike\"\n    }"], "sessionid": 50624 }`)
	case "raw":
		return []byte(`{ "status": "ok", "responses": [{"type": "raw","value": "If you are not 100% satisfied with your purchase, you can return your item to us for a full refund. Returns must be done within 30 days of receipt together with the Returns slip at a SingPost's counter or POPStation, un-used with tags on, in original packaging and must not fall under the list of non-refundable brands/items <a href=\"http://www.zalora.sg/faq-non-refundable/\">HERE</a>.\n\nFor more information regarding the Return policy, please view the steps <a href=\"http://www.zalora.sg/faq-returns/\">HERE</a>"}], "sessionid": 50510 }`)
	}
	return []byte(`{ "status": "ok", "responses": ["{\n      \"type\": \"top\",\n      \"key\": \"nike\"\n    }"], "sessionid": 50624 }`)
}

func (this *Chat) ChatWithBot(testing bool) error {
	answer := []byte{}
	fmt.Println(testing)
	if testing {
		answer = this.testingChat()
	} else {
		answer = SendChat(this.Msg)
	}
	fmt.Println(string(answer))
	chatBot := &ChatBot{}
	err := json.Unmarshal(answer, chatBot)
	fmt.Println(chatBot)
	if err != nil {
		beego.Alert(err)
		return errors.New("Bot didn't reply: " + err.Error())
	}
	reply := &Reply{}
	if len(chatBot.Responses) == 1 {
		replyInterface := chatBot.Responses[0]
		//		pretty.Println("xxx")
		//		pretty.Println(replyInterface)
		//		pretty.Println("xxx")
		if replyMap, ok := replyInterface.(map[string]interface{}); ok {
			//			pretty.Println("yyy")
			//			pretty.Println(replyMap)
			//			pretty.Println("yyy")
			if t, ok := replyMap["type"]; ok {
				reply.Type = t.(string)
			}
			if k, ok := replyMap["key"]; ok {
				reply.Key = k.(string)
			}
			if v, ok := replyMap["value"]; ok {
				reply.Value = v.(string)
			}
		} else if replyStr, ok := replyInterface.(string); ok {
			err := json.Unmarshal([]byte(replyStr), reply)
			if err != nil {
				//				pretty.Println("HAH")
				//				pretty.Println(err)
				//				pretty.Println(reply)
				//				pretty.Println("HAH")
				reply.Value = replyStr
			} else {
				//				pretty.Println("FK")
				//				pretty.Println(replyStr)
				//				pretty.Println(reply)
				//				pretty.Println("FK")
			}
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

func Say(msg string, testing bool, typeOfTest string) (*Chat, error) {
	chat := &Chat{Msg: msg, Reply: &Reply{Type: typeOfTest}}
	err := chat.ChatWithBot(testing)
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
	if r.Key == "" {
		topUrl = "https://api.zalora.sg/v1/products/?limit=3"
	}
	fmt.Println(topUrl)
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
