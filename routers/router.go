package routers

import (
	"github.com/astaxie/beego"
	"uugu.org/pcs/controllers"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/", &controllers.MainController{}, "*:Get")
	beego.Router("/index", &controllers.MainController{}, "*:Get")
	beego.Router("/sign", &controllers.MainController{}, "*:Sign")
	beego.Router("/doSign", &controllers.MainController{}, "*:DoSign")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/doLogin", &controllers.MainController{}, "*:DoLogin")
	//beego.Router("/index", &controllers.IndexController{}, "*:Get")
	beego.Router("/exam", &controllers.ExamController{}, "*:Index")
	beego.Router("/exam/getQuestions", &controllers.ExamController{}, "*:GetQuestions")
	beego.Router("/exam/doSubmit", &controllers.ExamController{}, "*:DoSubmit")
	beego.Router("/exam/result", &controllers.ExamController{}, "*:ToResult")
	beego.Router("/exam/getHollandAnswerResult", &controllers.ExamController{}, "*:GetHollandResult")
	beego.Router("/item/?:id:int", &controllers.ItemController{})
	//beego.AutoRouter(&controllers.ItemController{})
}
