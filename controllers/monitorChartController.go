package controllers

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
)

type MonitorChartController struct {
	BaseController
}



// 概览  API调用情况统计 log_app_quotas_statistic
func (this *MonitorChartController) AppsCalledStatics() {


	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	var queryFilter = make(map[string]interface{})

	queryFilter["start_time"] = this.GetString("startTime")
	queryFilter["end_time"] = this.GetString("endTime")

	//获取当前页需要的数据
	JSONData, err := services.GetAppsCalledStatics(queryFilter,domain)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}




////   应用调用总量 AppCalledTotal  3条线
//func (this *MonitorChartController) AppCalledTotal() {
//
//	var queryFilter = make(map[string]interface{})
//
//	queryFilter["start_time"] = this.GetString("startTime")
//
//	queryFilter["end_time"] = this.GetString("endTime")
//	//获取当前页需要的数据
//	JSONData, err := services.GetLogAppDailyCalledTotal(queryFilter)
//	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
//		return
//	}
//
//	this.JsonResult(httpUtil.SUCCESS, JSONData)
//	return
//}










