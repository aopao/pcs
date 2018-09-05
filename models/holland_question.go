package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"uugu.org/pcs/common"
	"time"
)

type HollandQuestion struct {
	Id         int64     `json:"id"`
	BaseKey    string    `json:"_"`
	Content    string    `json:"content"`
	Answer     int       `json:"_"`
	CreateTime time.Time `json:"_"`
	UpdateTime time.Time `json:"_"`
	Valid      int       `json:"valid"`
}

func (m *HollandQuestion) TableName() string {
	return TableName("pcs_holland_questions")
}

func GetHollandQuestionById(id int64) (*HollandQuestion, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM HollandQuestion WHERE id IN (?, ?, ?)", ids)
	category := HollandQuestion{Id: id}

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

func GetHollandQuestionByIds(ids *[]int64) (*[]HollandQuestion, error) {
	var results []HollandQuestion
	//num, err := o.Raw("SELECT * FROM category WHERE key = ?", param.Key).QueryRows(&categorys)
	num, err := o.QueryTable(new(HollandQuestion).TableName()).Filter("id__in", ids).Filter("valid", 1).All(&results)
	if err != nil {
		return nil, err
	}
	logs.Info("Get HollandQuestion nums: ", num)
	//categorys := []*HollandQuestion{}
	//o.QueryTable(new(HollandQuestion).TableName()).All(&categorys)
	return &results, nil
}

func GetHollandQuestions(param *HollandQuestion) (*[]HollandQuestion, error) {
	var results []HollandQuestion
	//num, err := o.Raw("SELECT * FROM category WHERE key = ?", param.Key).QueryRows(&categorys)
	num, err := o.QueryTable(param.TableName()).Filter("base_key", param.BaseKey).Filter("valid", 1).All(&results)
	if err != nil {
		return nil, err
	}
	logs.Info("Get HollandQuestion nums: ", num)
	//categorys := []*HollandQuestion{}
	//o.QueryTable(new(HollandQuestion).TableName()).All(&categorys)
	return &results, nil
}

func AddHollandQuestion(param *HollandQuestion) (*common.CommonResponse) {
	commonResponse := common.CommonResponse{}
	//o := orm.NewOrm()
	result, ormerr := o.Insert(param)
	if nil != ormerr {
		logs.Error("Add HollandQuestion Error:", ormerr)
		return commonResponse.ToFail(ormerr.Error())
	} else {
		return commonResponse.ToSuccess("Add HollandQuestion Success! Id:" + string(result))
	}
}
