package views

import (
	"openai-backend/models"
)

//返回数据集组列表和其下数据集列表JSON
func QuotasStatisticModelTempView(logQuotasStatisticModelList []*models.LogApiQuotasStatisticModel) []*models.LogApiQuotasStatisticModelTemp {

	var quotasStatisticModelTempList = make([]*models.LogApiQuotasStatisticModelTemp, 0)
	/*	var ids =make([]int,0)
		for _, model := range list {
			ids = append(ids,model.DatasetGroupId)
		}
		logs.Info(ids)*/


	for _, quotasStatisticModel := range logQuotasStatisticModelList {

		var quotasStatisticModelTemp = &models.LogApiQuotasStatisticModelTemp{

			ApiId   			 :quotasStatisticModel.ApiId,
			ApiName			 	:quotasStatisticModel.ApiName	,
			CalledTotal			 :quotasStatisticModel.SuccessfulCalled+quotasStatisticModel.FailedCalled  	,
			SuccessfulCalled     :quotasStatisticModel.SuccessfulCalled,
			FailedCalled		 :quotasStatisticModel.FailedCalled,

		}
		quotasStatisticModelTempList = append(quotasStatisticModelTempList, quotasStatisticModelTemp)
	}
	return quotasStatisticModelTempList
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
