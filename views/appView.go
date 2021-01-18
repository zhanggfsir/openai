package views

import "openai-backend/models"

/*
@Author:
@Time: 2020-11-05 15:02
@Description:APPmodel
*/


func AppsAbilityDetailInfo(app *models.AppModel) []*models.AppsAbilityModel {

	var appsAbilityModelList = make([]*models.AppsAbilityModel, 0)

	appAbilityModel := &models.AppAbilityModel{}
	appAbilityList := make([]*models.AppAbilityModel, 0)
	_, _ = appAbilityModel.QueryTable().Filter("app_id", app.Id).All(&appAbilityList)

	for _,appAbilityModel := range appAbilityList{
		ability := models.NewAbilityModel()
		_, _ = ability.QueryTable().Filter("id", appAbilityModel.AbilityId).All(ability)
		appsAbilityModel :=&models.AppsAbilityModel{
			AbilityId:   appAbilityModel.AbilityId,
			AbilityName: ability.AbilityName,
			ApiList:     AppsApiDetailInfo(app,ability),
		}

		appsAbilityModelList = append(appsAbilityModelList,appsAbilityModel)
	}

	return appsAbilityModelList
}


func AppsApiDetailInfo(app *models.AppModel,ability *models.AbilityModel) []*models.AppsApiModel {

	var AppsApiModelList = make([]*models.AppsApiModel, 0)


	apiList := make([]*models.ApiModel, 0)
	_, _ = models.NewApiModel().QueryTable().Filter("ability_id", ability.Id).All(&apiList)

	for _,apiModel := range apiList{

		quotas := &models.QuotasModel{}
		_, _ = quotas.QueryTable().Filter("app_id", app.Id).Filter("api_id", apiModel.Id).All(quotas)
		appsApiModel :=&models.AppsApiModel{
			ApiId:        apiModel.Id,
			ApiName:      apiModel.ApiName,
			Url:          apiModel.Url,
			MaxQps:       quotas.MaxQps,
			Quotas:       quotas.Quotas,
			QuotasPeriod: quotas.QuotasPeriod,
		}

		AppsApiModelList = append(AppsApiModelList,appsApiModel)
	}

	return AppsApiModelList
}

//AppList
func AppList(list []*models.AppModel) []*models.AppModelTemp {
	var AppListTemp = make([]*models.AppModelTemp, 0)

	for _, app := range list {
		abilityIdList := make([]uint64, 0)
		appAbilityModel := &models.AppAbilityModel{}

		appAbilityList := make([]*models.AppAbilityModel, 0)
		_, _ = appAbilityModel.QueryTable().Filter("app_id", app.Id).All(&appAbilityList)

		for _,appAbilityModel := range appAbilityList{
			abilityIdList = append(abilityIdList,appAbilityModel.AbilityId)
		}
		var appTemp = &models.AppModelTemp{
			Id:          app.Id,
			//AppId:       app.AppId,
			AppName:     app.AppName,
			ApiKey:      app.ApiKey,
			SecretKey:   app.SecretKey,
			AbilityList: abilityIdList,
			AppType:     app.AppType,
			Desc:        app.Desc,
			CreateTime:  app.CreateTime.Format("2006-01-02 15:04:05"),
			ModifyTime:  app.ModifyTime.Format("2006-01-02 15:04:05"),
		}
		AppListTemp = append(AppListTemp, appTemp)

	}
	return AppListTemp
}
