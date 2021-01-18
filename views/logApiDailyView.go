package views

import (
	"fmt"
	"openai-backend/models"
	"strconv"
)

//返回数据集组列表和其下数据集列表JSON
func GetLogApiDailyView(logApiDailyViewTempList  []*models.LogApiDailyViewTemp) []*models.LogApiDailyView {

	var logApiDailyViewList = make([]*models.LogApiDailyView, 0)

	for _, logApiDailyTemp := range logApiDailyViewTempList {
		var failedRate  float64
		total:=logApiDailyTemp.Successful+logApiDailyTemp.Failed
		if total==0{
			failedRate=0.0
		}else {
			failedRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(logApiDailyTemp.Failed)/float64(total)), 64)
		}

		var logApiDaily = &models.LogApiDailyView{

			ApiId			:logApiDailyTemp.ApiId,
			ApiName			:logApiDailyTemp.ApiName,
			Total			:logApiDailyTemp.Successful+logApiDailyTemp.Failed,
			Successful		:logApiDailyTemp.Successful,
			Failed			:logApiDailyTemp.Failed,
			FailedRate		:failedRate,

		}
		//logs.Info("---logApiDaily-->",logApiDaily)
		logApiDailyViewList = append(logApiDailyViewList, logApiDaily)
	}

	return logApiDailyViewList
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
