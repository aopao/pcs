package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"uugu.org/pcs/models"
	"strconv"
)

/**
 * Item Controller
 */
type ItemController struct {
	beego.Controller
}

func (c *ItemController) Get() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Param convert error!")
	}
	c.Data["Id"] = id

	category := models.BaseCategory{
		Name:        "Test",
		Description: "Hello",
	}
	category1, err := models.GetBaseCategoryById(id)
	if nil != err {
		logs.Error(err)
	} else {
		logs.Info(category1)
		c.Data["category"] = category1
	}
	categories, err := models.GetBaseCategories(&category)
	if nil != err {
		logs.Error(err)
	}
	logs.Info(categories)
	c.TplName = "item.html"
}
