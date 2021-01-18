package services

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/models"
	"openai-backend/views"
)







// 概览
func GetAppCalledTotal( queryFilter map[string]interface{}) (map[string]interface{}, error) {

	// 1.overViewLists API调用总量
	overViewLists, _, err := models.NewLogApiQuotasStatisticModel().GetAppQuotasStatisticModel(queryFilter)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}
	myOverViewList := views.QuotasStatisticModelTempView(overViewLists)

	////  API调用情况统计 2.1 totalSuccessFailed
	//totalSuccessFailed, err := models.NewQuotasStatisticModel().GetCalledTotal(queryFilter)
	//if err != nil {
	//	return nil, err
	//}
	//// API调用情况统计 2.2 success	failed
	//successFailed, err := models.NewQuotasStatisticModel().GetSuccessAndFailed(queryFilter)
	if err != nil {
		return nil, err
	}

	var jsonData = make(map[string]interface{})
	jsonData["overViewLists"] = myOverViewList
	//jsonData["totalSuccessFailed"] = totalSuccessFailed
	//jsonData["success"] = successFailed.Success
	//jsonData["failed"] = successFailed.Failed
	//jsonData["calledTotal"] =calledTotal
	return jsonData, nil

}







func GetMonitorChartService(queryFilter map[string]interface{},domain string) (map[string]interface{}, error) {

	var jsonData = make(map[string]interface{})

	// 1.得到所有label
	labelList, err := models.NewHourlyQuotasStatisticModel().GetLabel(queryFilter,domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}

	jsonData["label"] = labelList

	// 2.得到所有data [内]
	//   2.1 得到 调用成功
	//var allData = make(map[string]interface{})
	var allData = make([]interface{}, 0)

	var success = make(map[string]interface{})
	calledSuccess, err := models.NewHourlyQuotasStatisticModel().GetCalledSuccess(queryFilter,domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}

	//var jsonData = make(map[string]interface{})

	success["value"] = calledSuccess
	success["key"] = "调用成功"
	success["type"] = "line"
	allData = append(allData, success)
	//allData["total"] = Count
	//allData["calledSuccessData"] = success

	//   2.2 得到 调用失败
	var failed = make(map[string]interface{})
	calledFailed, err := models.NewHourlyQuotasStatisticModel().GetCalledFailed(queryFilter,domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}

	//var jsonData = make(map[string]interface{})
	failed["value"] = calledFailed
	failed["key"] = "调用失败"
	failed["type"] = "line"
	//successAndFaiedData["total"] = Count
	allData = append(allData, failed)
	//allData["calledFailedData"] = failed

	//   2.3 得到 调用总数
	// GetCalledSuccessAndFailed
	var successAndFailed = make(map[string]interface{})

	calledSuccessAndFailed, err := models.NewHourlyQuotasStatisticModel().GetCalledSuccessAndFailed(queryFilter,domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}

	//var jsonData = make(map[string]interface{})
	successAndFailed["value"] = calledSuccessAndFailed
	successAndFailed["key"] = "调用总数"
	successAndFailed["type"] = "line"

	allData = append(allData, successAndFailed)
	//allData["calledFailedAndSuccess"] = successAndFailed

	//allData["total"] = Count

	jsonData["data"] = allData

	return jsonData, nil
}

