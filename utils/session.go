package utils

import (
	"github.com/astaxie/beego"
	"uugu.org/pcs/models"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"github.com/astaxie/beego/context"
)

const SEESION_KEY = "sessionKey"
const USER_SEESION_KEY = "user_session_key-"

func GetUser(c *beego.Controller) (*models.User, error) {
	sessionKey := c.Ctx.GetCookie(SEESION_KEY)
	userJson, err := RedisGet(sessionKey)
	if nil != err {
		logs.Error(err)
		return nil, err
	}
	user := models.User{}
	if err := json.Unmarshal([]byte(userJson), &user); nil != err {
		logs.Error(err)
		return nil, err
	}
	return &user, nil
}

func GetUserName(c *context.Context) (string, error) {
	sessionKey := c.GetCookie(SEESION_KEY)
	userJson, err := RedisGet(sessionKey)
	if nil != err {
		logs.Error(err)
		return "", err
	}
	user := models.User{}
	if err := json.Unmarshal([]byte(userJson), &user); nil != err {
		logs.Error(err)
		return "", err
	}
	return user.UserName, nil
}

func VerifyLogin(c *context.Context) {
	requestUrl := c.Request.RequestURI
	logs.Info("开始验证访问权限，访问路径：" + requestUrl)

	sessionKey := c.GetCookie(SEESION_KEY)
	if len(sessionKey) == 0 {
		logs.Error("未发现用户的SessionKey!")
		c.Redirect(302, "/login")
		return
	}
	logs.Info("该用户的SessionKey:", sessionKey)
	userJson, err := RedisGet(sessionKey)
	if nil != err {
		logs.Error(err)
		c.Redirect(302, "/login")
		return
	}
	logs.Info("该用户的UserJson:", userJson)
	user := models.User{}
	if err := json.Unmarshal([]byte(userJson), &user); nil != err {
		logs.Error(err)
		c.Redirect(302, "/login")
		return
	}
	//单点登录验证
	sessionKey2, err := RedisGet(USER_SEESION_KEY + user.UserName)
	if nil != err || sessionKey2 != sessionKey {
		logs.Error("单点登录验证，SessionKey匹配失败！UserName：", user.UserName)
		c.Redirect(302, "/login")
		return
	}
	c.Input.SetParam("userName", user.UserName)
	c.Input.SetParam("name", user.Name)

	SetOrUpdateSession(&user, sessionKey)
}

func IsLogin(c *context.Context) bool {
	sessionKey := c.GetCookie(SEESION_KEY)
	_, err := RedisGet(sessionKey)
	if nil != err {
		logs.Error(err)
		return false
	}
	return true
}

func SetOrUpdateSession(user *models.User, sessionKey string) {
	userJson, err := json.Marshal(user)
	if nil != err {
		logs.Error(err)
		return
	}
	RedisSet(sessionKey, string(userJson), 3600*6)
	RedisSet(USER_SEESION_KEY+user.UserName, sessionKey, 3600*6)
}

func LogOut(sessionKey, username string) {
	RedisDel(sessionKey)
	RedisDel(USER_SEESION_KEY + username)
}
