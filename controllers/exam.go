package controllers

import (
	"github.com/astaxie/beego"
	"uugu.org/pcs/models"
	"github.com/astaxie/beego/logs"
	"time"
	"math/rand"
	"uugu.org/pcs/common"
	"strings"
	"uugu.org/pcs/services"
	"strconv"
	"encoding/json"
	"uugu.org/pcs/utils"
)

type ExamController struct {
	beego.Controller
}

const HollandQuestionNum = 60

func init() {
	rand.Seed(time.Now().UnixNano())
}

type HollandResult struct {
	Id          int64  `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Trait       string `json:"trait"`       //特点
	Career      string `json:"career"`      //事业，职业
	Description string `json:"description"` //`json:"description,omitempty"`
	Score       int    `json:"score"`
}

func (c *ExamController) Index() {
	//baseCategories, err := models.GetAllBaseCategories()
	//if nil != err {
	//	logs.Error(err)
	//	c.Data["message"] = "获取基础分类失败！"
	//	c.TplName = "error.html"
	//	c.Abort("500")
	//	return
	//}
	//hollandQuestions := []models.HollandQuestion{}
	//for _, baseCategory := range baseCategories {
	//	questions, err := models.GetHollandQuestions(&models.HollandQuestion{BaseId: baseCategory.Id})
	//	if nil != err {
	//		logs.Error("Get holland questions error, baseId:", baseCategory.Id)
	//		c.Data["message"] = "获取评测试题失败，未查询到试题类型：" + string(baseCategory.Id)
	//		c.TplName = "error.html"
	//		c.Abort("500")
	//	}
	//	hollandQuestions = append(hollandQuestions, getQuestions(&questions, HollandQuestionNum/len(baseCategories))...)
	//}
	//c.Data["hollandQuestions"] = sortQuestions(&hollandQuestions)
	userName := c.GetString("userName")
	hollandAnswerResult, err := models.GetHollandAnswerResultByUserName(userName)
	if nil != err {
		logs.Info("未获取到答题结果，可以进行测试！", err)
	}
	if nil != hollandAnswerResult {
		logs.Info("该用户已作答过该测试，直接跳转至结果分析页面！UserName:", userName)
		c.Redirect("/exam/result", 302)
	}
	c.Data["Name"] = c.GetString("name")
	c.Data["IsLogin"] = true
	c.TplName = "exam.html"
}

func (c *ExamController) GetQuestions() {
	response := common.CommonResponse{}
	baseCategories, err := models.GetAllBaseCategories()
	if nil != err {
		logs.Error(err)
		//c.Data["message"] = "获取基础分类失败！"
		//c.TplName = "error.html"
		//c.Abort("500")
		c.Data["json"] = response.ToFail("获取基础分类失败！")
		c.ServeJSON()
		return
	}
	hollandQuestions := []models.HollandQuestion{}
	for _, baseCategory := range baseCategories {
		questions, err := models.GetHollandQuestions(&models.HollandQuestion{BaseKey: baseCategory.Key})
		if nil != err {
			logs.Error("Get holland questions error, baseKey:", baseCategory.Key)
			//c.Data["message"] = "获取评测试题失败，未查询到试题类型：" + string(baseCategory.Key)
			//c.TplName = "error.html"
			//c.Abort("500")
			c.Data["json"] = response.ToFail("获取评测试题失败，未查询到试题类型：" + baseCategory.Key)
			c.ServeJSON()
			return
		}
		hollandQuestions = append(hollandQuestions, getSomeQuestions(questions, HollandQuestionNum/len(baseCategories))...)
	}
	c.Data["json"] = response.SetData(sortQuestions(&hollandQuestions)).ToSuccess("查询成功！")
	c.ServeJSON()
}

func (c *ExamController) DoSubmit() {
	response := common.CommonResponse{}
	answerStr := c.GetString("answerStr")
	if answerStr == "" {
		logs.Error("The answers you submitted was empty!")
		c.Data["json"] = response.ToFail("提交的答案为空！")
		c.ServeJSON()
		return
	}

	username, err := utils.GetUserName(c.Ctx)
	if nil != err {
		logs.Error("Get username from session fail!")
		c.Data["json"] = response.ToFail("Session过期，请重新登录！")
		c.ServeJSON()
		return
	}

	hollandModel := services.HollandModel{UserName: username, Answers: make(map[int64]int)}

	keys := make([]int64, 0)
	for _, answerPairStr := range strings.Split(answerStr, ",") {
		answerPair := strings.Split(answerPairStr, "|")
		if len(answerPair) != 2 {
			logs.Error("The submitted answers has error format!")
			c.Data["json"] = response.ToFail("提交的答案格式有误，请重试！")
			c.ServeJSON()
			return
		}
		questionId, err1 := strconv.ParseInt(answerPair[0], 10, 64)
		answer, err2 := strconv.Atoi(answerPair[1])

		if nil != err1 || nil != err2 {
			logs.Error("The submitted answer pair has error format!")
			c.Data["json"] = response.ToFail("提交的答案格式有误，请重试！")
			c.ServeJSON()
			return
		}
		hollandModel.Answers[questionId] = answer
		keys = append(keys, questionId)
	}

	hollandQuestions, err := models.GetHollandQuestionByIds(&keys)
	if nil != err {
		logs.Error(err)
		//c.Data["message"] = "获取基础分类失败！"
		//c.TplName = "error.html"
		//c.Abort("500")
		c.Data["json"] = response.ToFail("获取评测试题失败！")
		c.ServeJSON()
		return
	}
	hollandModel.Questions = *hollandQuestions
	// Compute holland model
	computeResponse := hollandModel.DoCompute()
	if computeResponse.IsSuccess() {
		logs.Error(computeResponse.Message)
		//c.Data["message"] = "获取基础分类失败！"
		//c.TplName = "error.html"
		//c.Abort("500")
		c.Data["json"] = computeResponse
		c.ServeJSON()
		return
	}
	//for _, baseCategory := range *baseCategories {
	//	logs.Info(baseCategory)
	//
	//}
	answerJson, err := json.Marshal(hollandModel.Answers)
	resultJson, err2 := json.Marshal(hollandModel.Result)
	if nil != err || nil != err2 {
		logs.Error(err)
		c.Data["json"] = response.ToFail("答题数据转Json格式失败！")
		c.ServeJSON()
		return
	}
	hollandAnswerResult := models.HollandAnswerResult{
		UserName:       username,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		OriginResult:   string(answerJson),
		Result:         string(resultJson),
		PreHollandType: hollandModel.PreHollandType,
		HollandType:    hollandModel.HollandType,
		Valid:          1,
	}
	saveResultResponse := models.AddHollandAnswerResult(&hollandAnswerResult)
	if saveResultResponse.IsSuccess() {
		c.Data["json"] = response.ToSuccess("提交答卷成功，模型数据分析成功！")
		c.ServeJSON()
	} else {
		c.Data["json"] = response.ToFail(saveResultResponse.Message)
		c.ServeJSON()
	}

}

// Get some questions
func getSomeQuestions(questions *[]models.HollandQuestion, u int) []models.HollandQuestion {
	questionsNum := len(*questions)
	if questionsNum < u {
		logs.Error("获取测试试题数量不足！存在数量：?，需要数量：?", len(*questions), u)
		return nil
	}
	if questionsNum == u {
		return *questions
	}
	result := []models.HollandQuestion{}
	//logs.Info(questions)
	for i := 0; i < u; i++ {
		randNum := rand.Intn(questionsNum)
		//logs.Info("获取的题号随机值为：", randNum)
		question := (*questions)[randNum]
		result = append(result, question)
		*questions = append((*questions)[:randNum], (*questions)[randNum+1:]...)
		questionsNum--
	}
	return result
}

// Sort the questions
func sortQuestions(question *[]models.HollandQuestion) []models.HollandQuestion {
	for i := len(*question); i > len(*question)/2; i-- {
		randIndex := rand.Intn(i)
		(*question)[i-1], (*question)[randIndex] = (*question)[randIndex], (*question)[i-1]
		//*question = (*question)[:n-1]
	}
	return *question
}

func (c *ExamController) ToResult() {
	//if !utils.IsLogin(c.Ctx) {
	//	c.Data["message"] = "未登录系统，请先登录！"
	//	c.TplName = "error.html"
	//	c.Abort("302")
	//}
	c.Data["Name"] = c.GetString("name")
	c.Data["IsLogin"] = true
	c.TplName = "result.html"
}

func (c *ExamController) GetHollandResult() {
	response := common.CommonResponse{}
	userName := c.GetString("userName")
	if userName == "" {
		logs.Error("The answers you submitted was empty!")
		c.Data["json"] = response.ToFail("提交的答案为空！")
		c.ServeJSON()
		return
	}
	hollandAnswerResult, err := models.GetHollandAnswerResultByUserName(userName)
	if nil != err {
		logs.Error("Get Answer result error!", err)
		c.Data["json"] = response.ToFail("获取答题结果失败！")
		c.ServeJSON()
		return
	}
	var originResult map[string]int
	err2 := json.Unmarshal([]byte(hollandAnswerResult.Result), &originResult)
	if nil != err2 {
		logs.Error("Unmarshal Answer result error!", err2)
		c.Data["json"] = response.ToFail("解析答题结果失败！")
		c.ServeJSON()
		return
	}
	baseCategories, err3 := models.GetAllBaseCategories()
	if nil != err3 {
		logs.Error("Get Base Category error!", err3)
		c.Data["json"] = response.ToFail("获取霍兰德基本类型失败！")
		c.ServeJSON()
		return
	}
	var result []HollandResult
	var scores []int
	for _, baseCategory := range baseCategories {
		hollandResult := HollandResult{
			Id:    baseCategory.Id,
			Key:   baseCategory.Key,
			Name:  baseCategory.Name,
			Score: originResult[baseCategory.Key],
		}
		if hollandResult.Score >= 35 {
			hollandResult.Description = baseCategory.LevelHigh
		}
		if hollandResult.Score < 35 && hollandResult.Score >= 25 {
			hollandResult.Description = baseCategory.LevelMiddle
		}
		if hollandResult.Score < 25 {
			hollandResult.Description = baseCategory.LevelLow
		}
		result = append(result, hollandResult)
		scores = append(scores, hollandResult.Score)
	}
	response.SetMapData("results", result)
	response.SetMapData("scores", scores)
	c.Data["json"] = response.ToSuccess("获取结果成功！")
	c.ServeJSON()
}
