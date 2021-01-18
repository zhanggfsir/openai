package controllers

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
	"time"
)



type LogApiDailyStatisticController struct {
	BaseController
}



// api 调用量日统计 log_hourly_quotas_statistic
func (this *LogApiDailyStatisticController) ApiDailyCalledLine() {

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	var queryFilter = make(map[string]interface{})
	queryFilter["start_time"] = this.GetString("startTime")
	queryFilter["end_time"] = this.GetString("endTime")

	if queryFilter["start_time"]==""{
		oldTime := time.Now().AddDate(0, 0, -7)
		queryFilter["start_time"]=oldTime.Format("2006-01-02 15:04:05" )
	}
	if queryFilter["end_time"]!=""{
		queryFilter["end_time"]=time.Now().Format("2006-01-02 15:04:05")
	}
	//获取当前页需要的数据
	JSONData, err := services.GetApiChartService(queryFilter,domain)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}



//  API用量统计  API调用总量
func (this *LogApiDailyStatisticController) ApiCalledTotalList() {

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	categoryId, err := this.GetUint64("categoryId")
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取页码错误", err) {
		logs.Error(err)
		return
	}

	var queryFilter = make(map[string]interface{})

	queryFilter["start_time"] = this.GetString("startTime")

	queryFilter["end_time"] = this.GetString("endTime")

	//获取当前页需要的数据
	JSONData, err := services.GetLogApiDailyCalledTotal(categoryId,queryFilter,domain)
	//logs.Info("2")
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		logs.Error(err)
		return
	}


	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}
