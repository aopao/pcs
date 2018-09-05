package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"uugu.org/pcs/common"
	"time"
)

type Token struct {
	Id         int64     `json:"id"`
	Token      string    `json:"token"`
	Type       int       `json:"type"`
	Status     int       `json:"status"`
	UserName   string    `json:"userName"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Valid      int       `json:"valid"`
}

func (m *Token) TableName() string {
	return TableName("pcs_tokens")
}

func GetTokenById(id int64) (*Token, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM Token WHERE id IN (?, ?, ?)", ids)
	token := Token{Id: id}

	//err := o.QueryTable("token").Filter("name", "slene").One(&token)

	ormerr := o.Read(&token)

	if ormerr == orm.ErrNoRows {
		logs.Error("查询不到,主键:", token.Id)
	} else if ormerr == orm.ErrMissPK {
		logs.Error("找不到主键:", token.Id)
	} else {
		return &token, nil
	}
	return nil, ormerr
}

func GetTokenByKey(tokenKey string) (*Token, error) {
	var token Token
	err := o.QueryTable(new(Token).TableName()).Filter("token", tokenKey).Filter("status", 1).Filter("valid", 1).One(&token)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		logs.Error("Returned Multi Rows Not One")
		return nil, err
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		logs.Error("Not row found")
		return nil, err
	}
	return &token, nil
}

func GetTokens(param *Token) ([]Token, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM token WHERE id IN (?, ?, ?)", ids)
	var tokens []Token
	//num, err := o.Raw("SELECT * FROM ? as bc WHERE bc.key = ?", param.TableName(), param.Key).QueryRows(&tokens)
	num, err := o.QueryTable(param.TableName()).Filter("user_name", param.UserName).All(&tokens)
	if err != nil {
		return nil, err
	}
	logs.Info("Token nums: ", num)
	//tokens := []*Token{}
	//o.QueryTable(new(Token).TableName()).All(&tokens)
	return tokens, nil
}

func GetAllTokens() ([]Token, error) {
	//o := orm.NewOrm()
	var tokens []Token
	//num, err := o.Raw("SELECT * FROM ? as bc WHERE bc.key = ?", param.TableName(), param.Key).QueryRows(&tokens)
	num, err := o.QueryTable(new(Token).TableName()).All(&tokens)
	if err != nil {
		return nil, err
	}
	logs.Info("Get Token nums: ", num)
	//tokens := []*Token{}
	//o.QueryTable(new(Token).TableName()).All(&tokens)
	return tokens, nil
}

func AddToken(token *Token) (*common.CommonResponse) {
	commonResponse := common.CommonResponse{}
	//o := orm.NewOrm()
	result, ormerr := o.Insert(token)
	if nil != ormerr {
		logs.Error("Add Token Error:", ormerr)
		return commonResponse.ToFail(ormerr.Error())
	} else {
		return commonResponse.ToSuccess("Add Token Success! Id:" + string(result))
	}
}

func UpdateToken(token *Token) bool {
	if num, err := o.Update(token); err == nil {
		logs.Info("更新Token成功，操作token：", token.Token, num)
		return true
	}
	return false
}