package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"uugu.org/pcs/common"
	"time"
)

type HollandAnswerResult struct {
	Id             int64     `json:"id"`
	UserName       string    `json:"userName"`
	OriginResult   string    `json:"originResult"`
	Result         string    `json:"result"`
	PreHollandType string    `json:"preHollandType"`
	HollandType    string    `json:"hollandType"`
	CreateTime     time.Time `json:"createTime"`
	UpdateTime     time.Time `json:"updateTime"`
	Valid          int       `json:"valid"`
}

func (m *HollandAnswerResult) TableName() string {
	return TableName("pcs_holland_answer_result")
}

func GetHollandAnswerResultById(id int64) (*HollandAnswerResult, error) {
	category := HollandAnswerResult{Id: id}
	ormerr := o.Read(&category)
	if ormerr == orm.ErrNoRows {
		logs.Error("查询不到,主键:", category.Id)
	} else if ormerr == orm.ErrMissPK {
		logs.Error("找不到主键:", category.Id)
	} else {
		return &category, nil
	}
	return nil, ormerr
}

func GetHollandAnswerResultByUserName(userName string) (*HollandAnswerResult, error) {
	var result HollandAnswerResult
	//num, err := o.Raw("SELECT * FROM category WHERE key = ?", param.Key).QueryRows(&categorys)
	err := o.QueryTable(new(HollandAnswerResult).TableName()).Filter("user_name", userName).One(&result)
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
	//categorys := []*HollandAnswerResult{}
	//o.QueryTable(new(HollandAnswerResult).TableName()).All(&categorys)
	return &result, nil
}

func GetHollandAnswerResults(param *HollandAnswerResult) (*[]HollandAnswerResult, error) {
	var results []HollandAnswerResult
	//num, err := o.Raw("SELECT * FROM category WHERE key = ?", param.Key).QueryRows(&categorys)
	num, err := o.QueryTable(param.TableName()).Filter("user_name", param.UserName).All(&results)
	if err != nil {
		return nil, err
	}
	logs.Info("Get HollandAnswerResult nums: ", num)
	//categorys := []*HollandAnswerResult{}
	//o.QueryTable(new(HollandAnswerResult).TableName()).All(&categorys)
	return &results, nil
}

func AddHollandAnswerResult(param *HollandAnswerResult) (*common.CommonResponse) {
	commonResponse := common.CommonResponse{}
	//o := orm.NewOrm()
	result, ormerr := o.Insert(param)
	if nil != ormerr {
		logs.Error("Add HollandAnswerResult Error:", ormerr)
		return commonResponse.ToFail(ormerr.Error())
	} else {
		return commonResponse.ToSuccess("Add HollandAnswerResult Success! Id:" + string(result))
	}
}
