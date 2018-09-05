package common

import (
	"github.com/astaxie/beego/logs"
	"reflect"
)

/*
* 公共应答对象
*
* Author: Zhiguo.Chen
 */
type CommonResponse struct {
	Success bool                   `json:"success"`
	Status  string                 `json:"status,omitempty"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    interface{}            `json:"data,omitempty"`
	MapData map[string]interface{} `json:"mapData,omitempty"`
}

func (this *CommonResponse) IsSuccess() bool {
	return this.Success
}

func (this *CommonResponse) ToSuccess(message string) *CommonResponse {
	this.Success = true
	this.Message = message
	return this
}

func (this *CommonResponse) ToFail(message string) *CommonResponse {
	this.Success = false
	this.Message = message
	return this
}

func (this *CommonResponse) SetData(data interface{}) *CommonResponse {
	this.Data = data
	return this
}

func (this *CommonResponse) SetMapData(key string, data interface{}) *CommonResponse {
	if nil == this.MapData {
		keyType := reflect.TypeOf(key)
		dataType := reflect.TypeOf(data)
		this.MapData = make(map[string]interface{})
		logs.Info("MapData未初始化，先初始化！keyType:{}, dataType:{}", keyType, dataType)
	}
	//(*this.MapData)[key] = data
	this.MapData[key] = data
	return this
}

func (this *CommonResponse) GetMapData(key string) interface{} {
	if nil == this.MapData {
		logs.Error("MapData未初始化，无数据可取！")
		return nil
	}
	return this.MapData[key]
}
