package controllers

import (
	"github.com/astaxie/beego"
	"uugu.org/pcs/models"
	"github.com/astaxie/beego/logs"
	"uugu.org/pcs/common"
	"uugu.org/pcs/utils"
	"time"
	"strings"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["IsLogin"] = utils.IsLogin(c.Ctx)
	c.TplName = "index.html"
}

func (c *MainController) Login() {
	c.TplName = "login.html"
}

func (c *MainController) DoLogin() {
	response := common.CommonResponse{}
	u := models.User{
		Valid:  1,
		Status: 1,
	}
	// if err := c.ParseForm(&u); err != nil {
	// 	logs.Error(err)
	// }
	u.UserName = c.GetString("username")
	u.Password = utils.GetMD5(common.PASSWORD_SALT + c.GetString("password"))
	if len(strings.TrimSpace(u.UserName)) == 0 {
		logs.Error("恶意构造的访问请求，请注意！IP：", c.Ctx.Request.RemoteAddr)
		c.Data["json"] = response.ToFail("用户名不可为空！")
		c.ServeJSON()
		return
	}
	if len(strings.TrimSpace(u.Password)) == 0 {
		logs.Error("恶意构造的访问请求，请注意！IP：", c.Ctx.Request.RemoteAddr)
		c.Data["json"] = response.ToFail("密码不可为空！")
		c.ServeJSON()
		return
	}
	users, err := models.GetUsers(&u)
	if nil != err {
		logs.Error(err.Error())
		c.Data["json"] = response.ToFail("登录服务未知错误，请稍后重试！")
		c.ServeJSON()
		return
	}
	if len(users) == 0 {
		logs.Error("用户登录，未成功查询到该用户！")
		c.Data["json"] = response.ToFail("用户名或密码错误，请重试！")
		c.ServeJSON()
		return
	} else {
		logs.Info("用户登录成功，用户名：", u.UserName)
		c.Data["json"] = response.ToSuccess("用户登录成功！")
		sessionKey := utils.GetUUID()
		utils.SetOrUpdateSession(&users[0], sessionKey)
		c.Ctx.SetCookie(utils.SEESION_KEY, sessionKey)
		c.ServeJSON()
	}
}

func (c *MainController) Sign() {
	c.TplName = "sign.html"
}

func (c *MainController) DoSign() {
	response := common.CommonResponse{}
	token, err := models.GetTokenByKey(c.GetString("token"))
	if err != nil || len(strings.TrimSpace(token.Token)) == 0 {
		logs.Error(err)
		c.Data["json"] = response.ToFail("该卡密无效，请联系销售商解决！")
		c.ServeJSON()
		return
	}
	u := models.User{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Valid:      1,
		Status:     1,
	}
	//if err := c.ParseForm(&u); err != nil {
	//	logs.Error(err)
	//}
	u.UserName = c.GetString("username")
	u.Name = c.GetString("name")
	u.Password = utils.GetMD5(common.PASSWORD_SALT + c.GetString("password"))
	u.Tel = c.GetString("tel")
	u.Email = c.GetString("email")
	u.Qq = c.GetString("qq")
	u.Province = c.GetString("province")
	u.City = c.GetString("city")
	u.Country = c.GetString("country")
	u.School = c.GetString("school")
	//u.Address = c.GetString("school")
	//u.Status, err = c.GetInt("status")
	if len(strings.TrimSpace(u.UserName)) == 0 {
		logs.Error("恶意构造的访问请求，请注意！IP：", c.Ctx.Request.RemoteAddr)
		c.Data["json"] = response.ToFail("用户名不可为空！")
		c.ServeJSON()
		return
	}
	if len(strings.TrimSpace(u.Password)) == 0 {
		logs.Error("恶意构造的访问请求，请注意！IP：", c.Ctx.Request.RemoteAddr)
		c.Data["json"] = response.ToFail("密码不可为空！")
		c.ServeJSON()
		return
	}
	if len(strings.TrimSpace(u.Tel)) == 0 {
		logs.Error("恶意构造的访问请求，请注意！IP：", c.Ctx.Request.RemoteAddr)
		c.Data["json"] = response.ToFail("电话不可为空！")
		c.ServeJSON()
		return
	}
	commonResponse := models.AddUser(&u)
	if commonResponse.IsSuccess() {
		token.UserName = u.UserName
		token.Status = 2 //更改为2已使用状态
		models.UpdateToken(token)
		c.Data["json"] = response.ToSuccess("恭喜您，开通账户成功！")
	} else {
		logs.Error(commonResponse.Message)
		c.Data["json"] = response.ToFail(commonResponse.Message)
	}
	c.ServeJSON()
}

func (c *MainController) Logout() {
	//response := common.CommonResponse{}

	userName := c.GetString("userName")

	if len(strings.TrimSpace(userName)) == 0 {
		logs.Error("注销时，请求中无用户名！恶意构造的访问请求，请注意！IP：", c.Ctx.Request.RemoteAddr)
		//c.Data["json"] = response.ToFail("用户名不可为空！")
		//c.ServeJSON()
		//return
		c.TplName = "login.html"
	}
	sessionKey := c.Ctx.GetCookie(utils.SEESION_KEY)
	if sessionKey == "" {
		logs.Error("注销时，请求中无SessionKey！恶意构造的访问请求，请注意！IP：", c.Ctx.Request.RemoteAddr)
		//c.Data["json"] = response.ToFail("SessionKey为空！")
		//c.ServeJSON()
		//return
	}
	utils.LogOut(sessionKey, userName)
	//c.Data["json"] = response.ToSuccess("用户注销成功！")
	//c.DelSession(utils.SEESION_KEY)
	c.Ctx.SetCookie(utils.SEESION_KEY, "")
	//c.ServeJSON()
	c.Redirect("/login", 302)
}
