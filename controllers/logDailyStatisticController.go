package controllers

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
	"time"
)



type LogDailyStatisticController struct {
	BaseController
}


//  api被app调用情况 单个api调用情况（柱状图）
func (this *LogDailyStatisticController) ApiCalledByAppHistogram() {

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	var queryFilter = make(map[string]interface{})

	apiId, err := this.GetUint64("api")
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取api1错误", err) {
		logs.Error(err)
		return
	}

	top, err := this.GetUint64("top")
	if err!=nil {
		top=10
	}

	queryFilter["api_id"] = apiId
	queryFilter["top"] = top


	if this.GetString("startTime")!=""{
		queryFilter["start_time"] = this.GetString("startTime")
	}else{
		currentTime := time.Now()
		oldTime := currentTime.AddDate(0, 0, -7)
		queryFilter["start_time"] =oldTime.Format("2006-01-02")
	}

	if this.GetString("endTime")!=""{
		queryFilter["end_time"] = this.GetString("endTime")
	}else{
		queryFilter["end_time"] =time.Now().Format("2006-01-02")
	}


	//dimension := this.GetString("dimension")
	//if dimension == "y" {
	//	queryFilter["year"] = dimension
	//}else if dimension == "m" {
	//	queryFilter["month"] = dimension
	//} else if dimension == "d" {
	//	queryFilter["day"] = dimension
	//} else {
	//	queryFilter["hour"] = dimension
	//}

	//获取当前页需要的数据
	JSONData, err := services.GetApiCalledByAppHistogram(queryFilter,domain)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}






//app调用api情况 单个app调用情况（折线图）
func (this *LogDailyStatisticController) AppCallApiLineChart() {

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	var queryFilter = make(map[string]interface{})

	appId, err := this.GetUint64("app")
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取api1错误", err) {
		logs.Error(err)
		return
	}

	top, err := this.GetUint64("top")
	if err!=nil {
		top=10
	}

	queryFilter["app_id"] = appId
	queryFilter["top"] = top


	if this.GetString("startTime")!=""{
		queryFilter["start_time"] = this.GetString("startTime")
	}else{
		currentTime := time.Now()
		oldTime := currentTime.AddDate(0, 0, -7)
		queryFilter["start_time"] =oldTime.Format("2006-01-02")
	}

	if this.GetString("endTime")!=""{
		queryFilter["end_time"] = this.GetString("endTime")
	}else{
		queryFilter["end_time"] =time.Now().Format("2006-01-02")
	}

	//logs.Info(queryFilter["start_time"] ,"   -->",queryFilter["end_time"] )
	//获取当前页需要的数据
	JSONData, err := services.GetAppCallApiLineChart(queryFilter,domain)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}







