package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"uugu.org/pcs/common"
	"time"
)

//实用型R、研究型I、艺术型A、社会型S、企业型E、常规型C
const (
	_ = iota
	R  // Seconds field, default 0
	I  // Minutes field, default 0
	A  // Hours field, default 0
	S  // Day of month field, default *
	E  // Month field, default *
	C  // Day of week field, default *
)

type BaseCategory struct {
	Id          int64     `json:"id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Trait       string    `json:"trait"`       //特点
	Career      string    `json:"career"`      //事业，职业
	LevelHigh   string    `json:"levelHigh"`   //高级
	LevelMiddle string    `json:"levelMiddle"` //中级
	LevelLow    string    `json:"levelLow"`    //低级
	Description string    `json:"description"` //`json:"description,omitempty"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
	Valid       int       `json:"valid"`
}

func (m *BaseCategory) TableName() string {
	return TableName("pcs_base_category")
}

func GetBaseCategoryById(id int64) (*BaseCategory, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM BaseCategory WHERE id IN (?, ?, ?)", ids)
	category := BaseCategory{Id: id}

	//err := o.QueryTable("user").Filter("name", "slene").One(&user)

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

func GetBaseCategoryByKey(baseKey string) (*BaseCategory, error) {
	var baseCategory BaseCategory
	err := o.QueryTable(new(BaseCategory).TableName()).Filter("key", baseKey).One(&baseCategory)
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
	return &baseCategory, nil
}

func GetBaseCategories(param *BaseCategory) ([]BaseCategory, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM category WHERE id IN (?, ?, ?)", ids)
	var categorys []BaseCategory
	//num, err := o.Raw("SELECT * FROM ? as bc WHERE bc.key = ?", param.TableName(), param.Key).QueryRows(&categorys)
	num, err := o.QueryTable(param.TableName()).Filter("key", param.Key).All(&categorys)
	if err != nil {
		return nil, err
	}
	logs.Info("BaseCategory nums: ", num)
	//categorys := []*BaseCategory{}
	//o.QueryTable(new(BaseCategory).TableName()).All(&categorys)
	return categorys, nil
}

func GetAllBaseCategories() ([]BaseCategory, error) {
	//o := orm.NewOrm()
	var categorys []BaseCategory
	//num, err := o.Raw("SELECT * FROM ? as bc WHERE bc.key = ?", param.TableName(), param.Key).QueryRows(&categorys)
	num, err := o.QueryTable(new(BaseCategory).TableName()).All(&categorys)
	if err != nil {
		return nil, err
	}
	logs.Info("Get BaseCategory nums: ", num)
	//categorys := []*BaseCategory{}
	//o.QueryTable(new(BaseCategory).TableName()).All(&categorys)
	return categorys, nil
}

func AddBaseCategory(param *BaseCategory) (*common.CommonResponse) {
	commonResponse := common.CommonResponse{}
	o := orm.NewOrm()
	result, ormerr := o.Insert(param)
	if nil != ormerr {
		logs.Error("Add BaseCategory Error:", ormerr)
		return commonResponse.ToFail(ormerr.Error())
	} else {
		return commonResponse.ToSuccess("Add BaseCategory Success! Id:" + string(result))
	}
}
