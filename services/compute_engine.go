package services

import (
	"uugu.org/pcs/models"
	"github.com/astaxie/beego/logs"
	"sort"
	"uugu.org/pcs/common"
)

// Compute Engine Interface
type ComputeEngine interface {
	DoCompute() *common.CommonResponse
}

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	return p
}

type HollandModel struct {
	UserName       string
	Answers        map[int64]int
	Questions      []models.HollandQuestion
	Result         map[string]int
	PreHollandType string
	HollandType    string
}

// Holland model key size
const HollandTypeKeySize = 3

// Compute the Holland model result
func (m *HollandModel) DoCompute() *common.CommonResponse {
	commonResponse := common.CommonResponse{}
	if nil == m.Answers || nil == m.Questions || len(m.Answers) == 0 || len(m.Questions) == 0 {
		logs.Error("Compute Error! No Answers or Question!")
		return commonResponse.ToFail("Compute Error! No Answers or Question!")
	}

	for _, hollandQuestion := range m.Questions {
		if nil == m.Result {
			m.Result = make(map[string]int)
		}
		if m.Result[hollandQuestion.BaseKey] == 0 {
			m.Result[hollandQuestion.BaseKey] = m.Answers[hollandQuestion.Id]
		} else {
			m.Result[hollandQuestion.BaseKey] = m.Result[hollandQuestion.BaseKey] + m.Answers[hollandQuestion.Id]
		}
	}

	result := sortMapByValue(m.Result)
	if result.Len() < HollandTypeKeySize {
		logs.Error("Result data to less, data num：", result.Len())
		return commonResponse.ToFail("Result data to less!")
	}
	//baseCategories := make([]models.BaseCategory, 0)
	//logs.Info(result)
	preHollandType := ""
	typeIndex, preBaseCategoryValue := 0, 100
	for _, pair := range result {
		baseCategory, err := models.GetBaseCategoryByKey(pair.Key)
		if nil != err {
			logs.Error("Can't find the BaseCategory, Key:", pair.Key)
			return commonResponse.ToFail("Can't find the BaseCategory, Key:" + string(pair.Key))
		}
		logs.Info("BaseCategory Key:", baseCategory.Key)
		//baseCategories = append(baseCategories, *baseCategory)

		// 26,25,20,20,6,4
		if typeIndex < HollandTypeKeySize || preBaseCategoryValue == pair.Value {
			preHollandType = preHollandType + baseCategory.Key
			preBaseCategoryValue = pair.Value
		}
		typeIndex ++
	}
	if len(preHollandType) > HollandTypeKeySize {
		logs.Info("出现不可确定混合类型：", preHollandType)
		m.PreHollandType = preHollandType
	} else {
		mixCategory, err := models.GetMixCategoryByKey(preHollandType)
		if nil != err {
			logs.Error(err)
			return commonResponse.ToFail(err.Error())
		}
		if nil == mixCategory {
			logs.Error("评测失败，未有匹配类型，Key：", preHollandType)
			return commonResponse.ToFail("评测失败，未有匹配类型，请重新作答！")
		}
		m.HollandType = preHollandType
	}

	return commonResponse.ToSuccess("引擎计算成功！")
}
