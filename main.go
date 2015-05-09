package main

import (
	_ "github.com/lei-cao/zbot/docs"
	_ "github.com/lei-cao/zbot/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
