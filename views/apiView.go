package views

import "openai-backend/models"

func ApiList(apis []*models.ApiModel) []*models.ApiModelTemp {
	var apiModelTempList = make([]*models.ApiModelTemp, 0)
	for _, api := range apis {
		//查询api的defaultQuotas
		defaultQuotas := models.DefaultQuotasModel{}
		_, _ = defaultQuotas.QueryTable().Filter("api_id", api.Id).All(&defaultQuotas)
		//
		var apiModelTemp = &models.ApiModelTemp{
			Id:           api.Id,
			ApiName:      api.ApiName,
			Url:          api.Url,
			MaxQps:       defaultQuotas.MaxQps,
			Quotas:       defaultQuotas.Quotas,
			QuotasPeriod: defaultQuotas.QuotasPeriod,
			CreateTime:   api.CreateTime.Format("2006-01-02 15:04:05"),
			ModifyTime:   api.ModifyTime.Format("2006-01-02 15:04:05"),
		}
		apiModelTempList = append(apiModelTempList, apiModelTemp)
	}
	return apiModelTempList
}


//func ApiList(apis []*models.ApiModel) []*models.ApiModelTemp {
//	var apiModelTempList = make([]*models.ApiModelTemp, 0)
//	defaultQuotas := models.DefaultQuotasModel{}
//	var apiModelTemp = &models.ApiModelTemp{}
//	for _, api := range apis {
//		//查询api的defaultQuotas
//		_, _ = defaultQuotas.QueryTable().Filter("api_id", api.Id).All(&defaultQuotas)
//		//
//
//		apiModelTemp.Id = api.Id
//		apiModelTemp.ApiName = api.ApiName
//		apiModelTemp.Url = api.Url
//		apiModelTemp.MaxQps = defaultQuotas.MaxQps
//		apiModelTemp.Quotas = defaultQuotas.Quotas
//		apiModelTemp.QuotasPeriod = defaultQuotas.QuotasPeriod
//		apiModelTemp.CreateTime = api.CreateTime.Format("2006-01-02 15:04:05")
//		apiModelTemp.ModifyTime = api.ModifyTime.Format("2006-01-02 15:04:05")
//
//
//		apiModelTempList = append(apiModelTempList, apiModelTemp)
//	}
//	return apiModelTempList
//}
