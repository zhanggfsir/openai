package controllers

import (
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
)

type AbilityController struct {
	BaseController
}

//添加AbilityCategory
func (this *AbilityController) Post() {

	//获取requestBody中的参数 校验数据
	app, err := services.ValidatePostAbility(this.Ctx.Input.RequestBody)
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