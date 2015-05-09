package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	"github.com/kr/pretty"
)

const (
	PANDORABOTS_API = "https://aiaas.pandorabots.com"
	BOT_ID          = "1409611904055"
	BOTNAME         = "nathalie"
	USER_KEY        = "412000ec3f11e9af1ad002b7218e25d9"
)

func GetBot() ([]byte, error) {
	botUri := fmt.Sprintf("%s/bot/%s?user_key=%s", PANDORABOTS_API, BOT_ID, USER_KEY)
	fmt.Println(botUri)
	return Get(botUri)
}

func SendChat(msg string) {
	var urlObj *url.URL
	chatUri := fmt.Sprintf("%s/talk/%s/%s", PANDORABOTS_API, BOT_ID, BOTNAME)
	urlObj, err := url.Parse(chatUri)
	if err != nil {
	}
	parameters := url.Values{}
	parameters.Add("user_key", USER_KEY)
	parameters.Add("input", msg)
	urlObj.RawQuery = parameters.Encode()

	//    https://aiaas.pandorabots.com/talk/APP_ID/BOTNAME?user_key=USER_KEY&input=INPUT
    pretty.Println(urlObj.String())
	Post(urlObj.String())
}

func Post(u string) ([]byte, error) {
	form := url.Values{}
	req, err := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	if err != nil {
		beego.Alert(err)
	}
	//	req.Header.Set("X-Custom-Header", "myvalue")
	//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(resp)
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body, nil
}

func Get(url string) ([]byte, error) {
	// request http api
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// read body
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return body, nil
}

func QueryBot(url string) {

}
