package controllers

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
	"time"
)

type LogHourlyQuotasStatisticController struct {
	BaseController
}





//   调用总量趋势统计 AppCalledTotalTrend  app 调用详情
//   todo 仅 按天查询
func (this *LogHourlyQuotasStatisticController) AppCalledTotalTrend() {

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	var queryFilter = make(map[string]interface{})

	appId, err := this.GetUint64("appId")
	if err != nil {
		queryFilter["app_id"] = 0
		//return
	} else {
		queryFilter["app_id"] = appId
	}


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


	dimension := this.GetString("dimension")
	dimension="d"
	if dimension == "y" {
		queryFilter["year"] = dimension
	}else if dimension == "m" {
		queryFilter["month"] = dimension
	} else if dimension == "d" {
		queryFilter["day"] = dimension
	} else {
		queryFilter["hour"] = dimension
	}

	//获取当前页需要的数据
	JSONData, err := services.GetAppCalledTotalTrend(queryFilter,domain)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}






// monitorChart log_hourly_quotas_statistic  监控报表		api 调用详情
func (this *LogHourlyQuotasStatisticController) AppApiMonitorChart() {

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	var queryFilter = make(map[string]interface{})

	appId, err := this.GetUint64("appId")

	if err != nil {
		queryFilter["app_id"] = 0
		//return
	} else {
		queryFilter["app_id"] = appId
	}

	apiId, err := this.GetInt("apiId")
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取apiId错误", err) {
		return
	}
	queryFilter["api_id"] = apiId

	//queryFilter["start_time"] = this.GetString("startTime")
	//queryFilter["end_time"] = this.GetString("endTime")

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


	dimension := this.GetString("dimension")
	if dimension == "y" {
		queryFilter["year"] = dimension
	}else if dimension == "m" {
		queryFilter["month"] = dimension
	} else if dimension == "d" {
		queryFilter["day"] = dimension
	} else {
		queryFilter["hour"] = dimension
	}

	//获取当前页需要的数据
	JSONData, err := services.GetMonitorChartService(queryFilter,domain)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}



