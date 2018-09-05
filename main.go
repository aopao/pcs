package main

import (
	"uugu.org/pcs/models"
	_ "uugu.org/pcs/routers"
	"github.com/astaxie/beego"
	"uugu.org/pcs/utils"
)

func init() {
	models.Init()
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.TemplateLeft = "<<<"
	beego.BConfig.WebConfig.TemplateRight = ">>>"
}

func main() {
	beego.InsertFilter("/exam/*", beego.BeforeRouter, utils.VerifyLogin)
	beego.InsertFilter("/logout", beego.BeforeRouter, utils.VerifyLogin)
	beego.Run()
}
