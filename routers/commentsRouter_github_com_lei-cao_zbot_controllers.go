package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["github.com/lei-cao/zbot/controllers:ChatController"] = append(beego.GlobalControllerRouter["github.com/lei-cao/zbot/controllers:ChatController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

}
