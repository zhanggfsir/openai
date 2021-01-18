package controllers

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
)

type ApiController struct {
	BaseController
}

//添加Api
func (this *ApiController) Post() {

	//获取requestBody中的参数 校验数据
	_, err := services.ValidatePostApi(this.Ctx.Input.RequestBody)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		return
	}

	//返回结果
	this.JsonResult(httpUtil.SUCCESS, nil)
	return
}

//显示数据集组列表
func (this *ApiController) GetApiList() {
	//userId, _ := this.preAuth(noonCasbin)
	var queryFilter = make(map[string]interface{})

	//获取当前页需要的数据
	jsonData, err := services.GetApiListService(queryFilter)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		logs.Error(err)
		return
	}

	//返回结果
	this.JsonResult(httpUtil.SUCCESS, jsonData)
	return
}











