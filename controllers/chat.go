package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lei-cao/zbot/models"
)

// Operations about object
type ChatController struct {
	beego.Controller
}

// @Title Get
// @Description chat with my bot
// @Param	msg		query 	string	true		"the chat msg user sent"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (this *ChatController) Get() {
	msg := this.GetString("msg")
	chat, err := models.Say(msg)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"message": err}
	} else {
		this.Data["json"] = map[string]interface{}{"data": &chat}
	}
	this.ServeJson(false)
}
