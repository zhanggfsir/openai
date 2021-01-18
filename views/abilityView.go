package views

import (
	"openai-backend/models"
)

//返回ability列表JSON
func AbilityList(abilitys []*models.AbilityModel) []*models.AbilityModelTemp {
	var abilityModelTempList = make([]*models.AbilityModelTemp, 0)
	for _, ability := range abilitys {

		//获取ability下api
		var apis = make([]*models.ApiModel, 0)
		_, _ = models.NewApiModel().QueryTable().Filter("ability_id",
			ability.Id).OrderBy("-id").All(&apis)


		var apiModelTempList = make([]*models.ApiModelTemp, 0)
		apiModelTempList = ApiList(apis)

		var abilityModelTemp = &models.AbilityModelTemp{
			AbilityId: ability.Id,
			AbilityName: ability.AbilityName,
			ApiList: apiModelTempList,
		}
		abilityModelTempList = append(abilityModelTempList, abilityModelTemp)
	}
	return abilityModelTempList
}

//func AbilityList(abilitys []*models.AbilityModel) []*models.AbilityModelTemp {
//	var abilityModelTempList = make([]*models.AbilityModelTemp, 0)
//	apiModel := models.NewApiModel()
//	var apis = make([]*models.ApiModel, 0)
//	var apiModelTempList = make([]*models.ApiModelTemp, 0)
//	var abilityModelTemp = &models.AbilityModelTemp{}
//	for _, ability := range abilitys {
//
//		//获取ability下api
//
//		_, _ = apiModel.QueryTable().Filter("ability_id", ability.Id).OrderBy("-id").All(&apis)
//
//		apiModelTempList = ApiList(apis)
//
//		abilityModelTemp.AbilityId = ability.Id
//		abilityModelTemp.AbilityName = ability.AbilityName
//		abilityModelTemp.ApiList = apiModelTempList
//
//		abilityModelTempList = append(abilityModelTempList, abilityModelTemp)
//	}
//	return abilityModelTempList
//}
