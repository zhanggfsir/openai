package controllers

import (
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
	"time"
)

type LogApiDailyHistoryAddUpController struct {
	BaseController
}



// api 调用量日统计 log_hourly_quotas_statistic
func (this *LogApiDailyHistoryAddUpController) ApiDailyHistoryAddUp() {

	//domain, err := this.getDomain()
	//if err != nil || domain == "" {
	//	logs.Error("err", "域名错误")
	//	this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
	//	return
	//}

	var queryFilter = make(map[string]interface{})
	queryFilter["start_time"] = this.GetString("startTime")
	queryFilter["end_time"] = this.GetString("endTime")

	if queryFilter["start_time"]==""{
		oldTime := time.Now().AddDate(0, 0, -7)
		queryFilter["start_time"]=oldTime.Format("2006-01-02" )
	}
	if queryFilter["end_time"]!=""{
		queryFilter["end_time"]=time.Now().Format("2006-01-02")
	}
	//获取当前页需要的数据
	JSONData, err := services.GetApiDailyHistoryAddUp(queryFilter)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}






