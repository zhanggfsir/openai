package services

import "openai-backend/models"

//  log_quotas_statistic_app
func GetAppsCalledStatics(queryFilter map[string]interface{}, domain string) (map[string]interface{}, error) {

	/*	// 1.overViewLists API调用总量
		overViewLists, _, err := models.NewQuotasStatisticModel().GetQuotasStatisticModel(queryFilter)
		if err != nil {
			logs.Error("get monitor chart error", err)
			return nil, err
		}
		myOverViewList := views.QuotasStatisticModelTempView(overViewLists)*/

	//  API调用情况统计 2.1 totalSuccessFailed
	totalSuccessFailed, err := models.NewLogAppQuotasStatisticModel().GetAppCalledTotal(queryFilter,domain)
	if err != nil {
		return nil, err
	}
	// API调用情况统计 2.2 success	failed
	successFailed, err := models.NewLogAppQuotasStatisticModel().GetAppSuccessAndFailed(queryFilter,domain)
	if err != nil {
		return nil, err
	}

	var jsonData = make(map[string]interface{})
	//jsonData["overViewLists"] = myOverViewList
	jsonData["totalSuccessFailed"] = totalSuccessFailed
	jsonData["success"] = successFailed.Success
	jsonData["failed"] = successFailed.Failed
	//jsonData["calledTotal"] =calledTotal
	return jsonData, nil

}

//