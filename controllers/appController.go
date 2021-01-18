package controllers

/*
@Author:
@Time: 2020-11-05 15:02
@Description:APPmodel
*/

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
)

type AppController struct {
	BaseController
}

//App列表
func (this *AppController) GetAppList() {
	var queryFilter = make(map[string]interface{})

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}
	// 获取页码和每页条数
	pageNo, err := this.GetInt("pageNo")
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取页码错误", err) {
		return
	}

	pageSize, err := this.GetInt("pageSize")
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取页码错误", err) {
		return
	}

	//获取当前页需要的数据
	JSONData, err := services.GetAppListService(domain, pageNo, pageSize, queryFilter)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}

//app detail
func(this *AppController) Get() {
	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	appId, err := this.GetUint64("appId")
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取id错误", err) {
		return
	}

	//获取当前页需要的数据
	JSONData, err := services.GetAppDetailService(domain,appId)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}


//添加app
func (this *AppController) Post() {
	//domain, err := this.getDomain()
	//if err != nil || domain == "" {
	//	logs.Error("err=", "域名错误")
	//	this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
	//	return
	//}
	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	//获取requestBody中的参数 校验数据
	_, err = services.ValidatePostApp(domain,this.Ctx.Input.RequestBody)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		return
	}



	//返回结果
	this.JsonResult(httpUtil.SUCCESS, nil)
	return
}

//删除app
func (this *AppController) Delete() {
	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	//校验数据
	app, Ids, err := services.ValidateDeleteApp(domain,this.Ctx.Input.RequestBody)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		logs.Error(err)
		return
	}

	//删除设备和设备下的数据
	o := orm.NewOrm()
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		logs.Error(err)
		return
	}


	err = o.Begin()
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		logs.Error(err)
		return
	}
	err = app.DeleteById(o, Ids)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		logs.Error(err)
		return
	}

	err = o.Commit()
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		logs.Error(err)
		return
	}

	//
	this.JsonResult(httpUtil.SUCCESS, nil)
	return
}


//修改app
func (this *AppController) Modify() {
	//domain, err := this.getDomain()
	//if err != nil || domain == "" {
	//	logs.Error("err=", "域名错误")
	//	this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
	//	return
	//}
	//this.Ctx.Input.Context
	//获取requestBody中的参数 校验数据
	_, err := services.ValidateModifyApp(this.Ctx.Input.RequestBody)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "", err) {
		return
	}


	//返回结果
	this.JsonResult(httpUtil.SUCCESS, nil)
	return
}