package controllers

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/models"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
)

type AbilityCategoryController struct {
	BaseController
}

//添加AbilityCategory
func (this *AbilityCategoryController) Post() {

	//获取requestBody中的参数 校验数据
	app, err := services.ValidatePostAbilityCategory(this.Ctx.Input.RequestBody)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		return
	}

	//插入数据库
	err = app.Insert()
	if this.checkErr(httpUtil.PARAMETER_ERROR, "插入数据库失败", err) {
		return
	}

	//返回结果
	this.JsonResult(httpUtil.SUCCESS, nil)
	return
}


func (this *AbilityCategoryController) Get() {

	category:=&models.AbilityCategoryModel{}
	categoryList:=make([]*models.AbilityCategoryModel,0)
	_,err:=category.QueryTable().All(&categoryList)
	if err!=nil{
		this.JsonResult(httpUtil.SUCCESS, nil)
		logs.Error(err)
		return
	}
	var jsonData = make(map[string]interface{})
	jsonData["list"] = categoryList
	this.JsonResult(httpUtil.SUCCESS, jsonData)

	return
}








