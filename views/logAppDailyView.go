package views

import (
	"fmt"
	"openai-backend/models"
	"strconv"
)

//返回数据集组列表和其下数据集列表JSON
//func GetLogAppDailyView(logAppDailyViewTempList []*models.LogAppDailyViewTemp) []*models.LogAppDailyView {
func GetLogAppDailyView(logAppDailyViewTempList []*models.LogAppDailyViewTemp) []*models.LogAppDailyView {

	var logAppDailyViewList = make([]*models.LogAppDailyView, 0)

	for _, logAppDailyTemp := range logAppDailyViewTempList {

		var failedRate  float64
		total:=logAppDailyTemp.Successful + logAppDailyTemp.Failed
		if total==0{
			failedRate=0.0
		}else {
			failedRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(logAppDailyTemp.Failed)/float64(total)), 64)
		}

		var logAppDaily = &models.LogAppDailyView{

			AppId:      logAppDailyTemp.AppId,
			AppName:    logAppDailyTemp.AppName,
			Total:      total,
			Successful: logAppDailyTemp.Successful,
			Failed:     logAppDailyTemp.Failed,
			FailedRate: failedRate,
		}

		logAppDailyViewList = append(logAppDailyViewList, logAppDaily)
	}
	return logAppDailyViewList
}

//func AbilityCategoryList(list []*models.AbilityCategoryModel) []*models.AbilityCategoryModelTemp {
//	var abilityCategoryModelTempList = make([]*models.AbilityCategoryModelTemp, 0)
//	/*	var ids =make([]int,0)
//		for _, model := range list {
//			ids = append(ids,model.DatasetGroupId)
//		}
//		logs.Info(ids)*/
//
//	//读取AbilityCategory下的Ability
//	var abilitys = make([]*models.AbilityModel, 0)
//	var abilityModelTempList = make([]*models.AbilityModelTemp, 0)
//	var abilityCategoryModelTemp = &models.AbilityCategoryModelTemp{}
//	for _, abilityCategory := range list {
//
//		_, _ = models.NewAbilityModel().QueryTable().Filter("category_id", abilityCategory.Id).OrderBy("-id").All(&abilitys)
//
//		abilityModelTempList = AbilityList(abilitys)
//
//		abilityCategoryModelTemp.CategoryId = abilityCategory.Id
//		abilityCategoryModelTemp.CategoryName = abilityCategory.CategoryName
//		abilityCategoryModelTemp.AbilityList = abilityModelTempList
//
//		abilityCategoryModelTempList = append(abilityCategoryModelTempList, abilityCategoryModelTemp)
//	}
//	return abilityCategoryModelTempList
//}
