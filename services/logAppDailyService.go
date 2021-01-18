package services

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/models"
	"openai-backend/views"
)

// APP调用总量
//func GetLogAppDailyCalledTotal(queryFilter map[string]interface{}) (map[string]interface{}, error) {
//
//	// 1.overViewLists API调用总量
//	logAppDailyModel, _, err := models.NewLogAppDailyStatisticModel().GetLogAppDailyModel(queryFilter)
//	if err != nil {
//		logs.Error("get monitor chart error", err)
//		return nil, err
//	}
//	myOverViewList := views.GetLogAppDailyView(logAppDailyModel)
//
//	////  API调用情况统计 2.1 totalSuccessFailed
//	//totalSuccessFailed, err := models.NewQuotasStatisticModel().GetCalledTotal(queryFilter)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//// API调用情况统计 2.2 success	failed
//	//successFailed, err := models.NewQuotasStatisticModel().GetSuccessAndFailed(queryFilter)
//	if err != nil {
//		return nil, err
//	}
//	//appModelList:=make([]* models.AppModel,0)
//	appModel:=&models.AppModel{}
//	appTotalNum,err:=appModel.QueryTable().Distinct().Count()
//	if err != nil {
//		return nil, err
//	}
//
//	var jsonData = make(map[string]interface{})
//	jsonData["overViewLists"] = myOverViewList	//todo appCalledTotal
//	jsonData["appTotalNum"]   = appTotalNum
//	//jsonData["totalSuccessFailed"] = totalSuccessFailed
//	//jsonData["success"] = successFailed.Success
//	//jsonData["failed"] = successFailed.Failed
//	//jsonData["calledTotal"] =calledTotal
//	return jsonData, nil
//
//}






// APP调用总量
func GetLogAppDailyCalledTotal( queryFilter map[string]interface{} , domain string) (map[string]interface{}, error) {

	// 1.overViewLists API调用总量
	logAppDailyViewTempList, _, err := models.NewLogAppDailyStatisticModel().GetLogAppDailyModel(queryFilter,domain)
	if err != nil {
		logs.Error("get monitor chart error", err)
		return nil, err
	}
	myOverViewList := views.GetLogAppDailyView(logAppDailyViewTempList)

	////  API调用情况统计 2.1 totalSuccessFailed
	//totalSuccessFailed, err := models.NewQuotasStatisticModel().GetCalledTotal(queryFilter)
	//if err != nil {
	//	return nil, err
	//}
	//// API调用情况统计 2.2 success	failed
	//successFailed, err := models.NewQuotasStatisticModel().GetSuccessAndFailed(queryFilter)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	//app:=models.AppModel{}
	//appList:=make([]*models.AppModel,0)
	////
	//appTotal,err:=app.QueryTable().Filter("domain",domain).All(&appList)

	if err != nil {
		logs.Error(err)
		return nil, err
	}

	var jsonData = make(map[string]interface{})
	jsonData["overViewLists"] = myOverViewList
	jsonData["appTotal"] =  len(myOverViewList)

	//jsonData["totalSuccessFailed"] = totalSuccessFailed
	//jsonData["success"] = successFailed.Success
	//jsonData["failed"] = successFailed.Failed
	//jsonData["calledTotal"] =calledTotal
	return jsonData, nil

}
