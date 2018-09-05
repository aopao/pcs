package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"uugu.org/pcs/common"
	"time"
)

type MixCategory struct {
	Id         int64     `json:"id"`
	Key        string    `json:"key"`
	Trait      string    `json:"trait"`  //特点
	Career     string    `json:"career"` //事业，职业
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Valid      int       `json:"valid"`
}

func (m *MixCategory) TableName() string {
	return TableName("pcs_mix_category")
}

func GetMixCategoryById(id int64) (*MixCategory, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM MixCategory WHERE id IN (?, ?, ?)", ids)
	category := MixCategory{Id: id}

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

func GetMixCategoryByKey(mixKey string) (*MixCategory, error) {
	var mixCategory MixCategory
	err := o.QueryTable(new(MixCategory).TableName()).Filter("key", mixKey).One(&mixCategory)
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
	return &mixCategory, nil
}

func GetMixCategories(param *MixCategory) ([]MixCategory, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM category WHERE id IN (?, ?, ?)", ids)
	var categorys []MixCategory
	num, err := o.Raw("SELECT * FROM category WHERE key = ?", param.Key).QueryRows(&categorys)
	if err != nil {
		return nil, err
	}
	logs.Info("Get MixCategory nums: ", num)
	//categorys := []*MixCategory{}
	//o.QueryTable(new(MixCategory).TableName()).All(&categorys)
	return categorys, nil
}

func AddMixCategory(param *MixCategory) (*common.CommonResponse) {
	commonResponse := common.CommonResponse{}
	//o := orm.NewOrm()
	result, ormerr := o.Insert(param)
	if nil != ormerr {
		logs.Error("Add MixCategory Error:", ormerr)
		return commonResponse.ToFail(ormerr.Error())
	} else {
		return commonResponse.ToSuccess("Add MixCategory Success! Id:" + string(result))
	}
}
