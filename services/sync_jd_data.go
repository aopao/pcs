package services

import (
	"github.com/astaxie/beego"
	"time"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"uugu.org/pcs/utils"
)

//var jd_access_token string
//var jd_app_key string
//var jd_app_secret string
var basePram map[string]string

func Init() {
	jdAccessToken := beego.AppConfig.String("jd_access_token")
	jdAppKey := beego.AppConfig.String("jd_app_key")
	jdAppSecret := beego.AppConfig.String("jd_app_secret")
	basePram = map[string]string{"v": "2.0", "access_token": jdAccessToken, "app_key": jdAppKey, "app_secret": jdAppSecret}
}
func SyncJDCategoryData() {
	Init()
	basePram["method"] = "jingdong.union.search.goods.category.query"
	basePram["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	basePram["sign"] = "9C611C67E1AA2E498E24B7553887B537"
	//360buy_param_json={"parent_id":"0","grade":"0"}&timestamp=2017-12-31 00:55:54&sign=9C611C67E1AA2E498E24B7553887B537
	paramJson := make(map[string]string)
	paramJson["parent_id"] = "0"
	paramJson["grade"] = "0"

	paramJsonStr, err := json.Marshal(paramJson)
	if nil != err {
		logs.Error(err)
		return
	}
	basePram["360buy_param_json"] = string(paramJsonStr)
	response, err := utils.GetRequestWithParams("https://api.jd.com/routerjson", basePram)
	if nil != err {
		logs.Error(err)
	}
	var mapResult map[string]interface{}
	if err := json.Unmarshal(response, &mapResult); err != nil {
		logs.Error(err)
	}
	//logs.Info(mapResult)
	categoryResponce := mapResult["jingdong_union_search_goods_category_query_responce"].(map[string]interface{})
	//logs.Info(categoryResponce["querygoodscategory_result"])
	var categoryResponseMap map[string]interface{}
	if err := json.Unmarshal([]byte(categoryResponce["querygoodscategory_result"].(string)), &categoryResponseMap); err != nil {
		logs.Error(err)
	}
	//logs.Info(categoryResponseMap["data"])
	//categoryMaps := categoryResponseMap["data"].([]interface{})

	//for _, categoryMap := range categoryMaps {
	//	//logs.Info("Name=%v, Id=%v\n", categoryMap.(map[string]interface{})["name"], categoryMap.(map[string]interface{})["id"])
	//	category := models.Category{
	//		Id:          int64(categoryMap.(map[string]interface{})["id"].(float64)),
	//		ParentId:    int64(categoryMap.(map[string]interface{})["parentId"].(float64)),
	//		Grade:       int64(categoryMap.(map[string]interface{})["grade"].(float64)),
	//		Name:        categoryMap.(map[string]interface{})["name"].(string),
	//		Description: "",
	//	}
	//	logs.Info(category)
	//	models.AddCategory(&category)
	//}

}
