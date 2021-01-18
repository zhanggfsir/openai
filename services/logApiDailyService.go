package services

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/models"
	"openai-backend/views"
)

func GetApiChartService(queryFilter map[string]interface{},domain string) (map[string]interface{}, error) {

	// 1.得到所有label
	labelList, err := models.NewLogApiDailyStatistic().GetApiLabel(queryFilter,domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}

	var jsonData = make(map[string]interface{})
	jsonData["label"] = labelList

	// 2.得到所有data [内]
	//   2.1 得到 调用成功
	//var allData = make(map[string]interface{})
	var allData = make([]interface{}, 0)

	var success = make(map[string]interface{})
	calledSuccess, err := models.NewLogApiDailyStatistic().GetApiCalledSuccess(queryFilter,domain)
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

	calledFailed, err := models.NewLogApiDailyStatistic().GetApiCalledFailed(queryFilter,domain)
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

	calledSuccessAndFailed, err := models.NewHourlyQuotasStatisticModel().GetApiCalledSuccessAndFailed(queryFilter,domain)
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







// API调用总量
func GetLogApiDailyCalledTotal(categoryId uint64, queryFilter map[string]interface{},domain string) (map[string]interface{}, error) {

	// 1.overViewLists API调用总量
	LogApiDailyViewTempList, _, err := models.NewLogApiDailyStatistic().GetLogApiDailyModel(categoryId,queryFilter,domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}


	myOverViewList := views.GetLogApiDailyView(LogApiDailyViewTempList)

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
	//logs.Info("1")
	return jsonData, nil

}




