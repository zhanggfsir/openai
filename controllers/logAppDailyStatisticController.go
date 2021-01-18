package controllers

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
)

type LogAppDailyStatisticController struct {
	BaseController
}


//   应用调用总量
func (this *LogAppDailyStatisticController) AppDailyCalledTotalList() {

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
	JSONData, err := services.GetLogAppDailyCalledTotal(queryFilter,domain)
	if this.checkErr(httpUtil.PARAMETER_ERROR, "获取数据失败", err) {
		logs.Error(err)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, JSONData)
	return
}