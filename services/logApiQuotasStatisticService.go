package services

import (
	"github.com/astaxie/beego/logs"
	"openai-backend/models"
	"openai-backend/views"
)

// 概览
func GetApisCalledStatics(queryFilter map[string]interface{},domain string) (map[string]interface{}, error) {

	/*	// 1.overViewLists API调用总量
		overViewLists, _, err := models.NewQuotasStatisticModel().GetQuotasStatisticModel(queryFilter)
		if err != nil {
			logs.Error("get monitor chart error", err)
			return nil, err
		}
		myOverViewList := views.QuotasStatisticModelTempView(overViewLists)*/

	//  API调用情况统计 2.1 totalSuccessFailed
	totalSuccessFailed, err := models.NewLogApiQuotasStatisticModel().GetCalledTotal(queryFilter,domain)
	if err != nil {
		return nil, err
	}
	// API调用情况统计 2.2 success	failed
	successFailed, err := models.NewLogApiQuotasStatisticModel().GetSuccessAndFailed(queryFilter,domain)
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



//获取api列表
func CategoryPieChart(queryFilter map[string]interface{},domain string) (map[string]interface{}, error) {

	categoryList,err:=models.NewAbilityCategoryModel().CategoryPieChart(domain)
	if err!=nil{
		logs.Error(err)
		return nil, err
	}
	//for i, i2 := range categoryList {
	//	logs.Info(i,i2)
	//}
	//logs.Info(category)
	//model := &models.AbilityCategoryModel{}
	//lists, Count, err := model.AbilityCategoryList(queryFilter)
	//if err != nil {
	//	logs.Error("data group lists error", err)
	//	return nil, err
	//}
	////logs.Info(lists, Count)
	//
	datalist := views.CategoryList(categoryList)
	////logs.Info(datalist)
	var jsonData = make(map[string]interface{})
	jsonData["data"] = datalist
	//jsonData["total"] = Count
	return jsonData, nil
}
