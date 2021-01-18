package views

import (
	"openai-backend/models"
)

//返回数据集组列表和其下数据集列表JSON
func AbilityCategoryList(list []*models.AbilityCategoryModel) []*models.AbilityCategoryModelTemp {
	var abilityCategoryModelTempList = make([]*models.AbilityCategoryModelTemp, 0)
	/*	var ids =make([]int,0)
		for _, model := range list {
			ids = append(ids,model.DatasetGroupId)
		}
		logs.Info(ids)*/

	//读取AbilityCategory下的Ability
	var abilitys = make([]*models.AbilityModel, 0)
	for _, abilityCategory := range list {

		_, _ = models.NewAbilityModel().QueryTable().Filter("category_id",
			abilityCategory.Id).OrderBy("-id").All(&abilitys)


		abilityModelTempList := AbilityList(abilitys)

		var abilityCategoryModelTemp = &models.AbilityCategoryModelTemp{
			CategoryId:   abilityCategory.Id,
			CategoryName: abilityCategory.CategoryName,
			AbilityList:  abilityModelTempList,
		}
		abilityCategoryModelTempList = append(abilityCategoryModelTempList, abilityCategoryModelTemp)
	}
	return abilityCategoryModelTempList
}

type CategoryModelView struct {
	Id           uint64 `orm:"auto;pk" description:"ID" json:"id"`
	CategoryName string `description:"能力种类名称" json:"categoryName"`
	SuccessfulCalled       uint64 `description:"apiId" json:"successful_called"`
	FailedCalled       uint64 `description:"apiId" json:"failed_called"`
}

func CategoryList(categoryList []*models.CategoryModelTemp) []interface{} {
	//var apiModelTempList = make([]*models.ApiModelTemp, 0)
	var dataList=make([]interface{},0)
	for _, category := range categoryList {
		var parent=make(map[string]interface{})
		parent["label"]=category.CategoryName
		parent["value"]=category.SuccessfulCalled+category.FailedCalled

		var children=make([]interface{},0)
		var childrenSuccess=make(map[string]interface{})
		childrenSuccess["label"]="成功"
		childrenSuccess["value"]=category.SuccessfulCalled

		var childrenFail=make(map[string]interface{})
		childrenFail["label"]="失败"
		childrenFail["value"]=category.FailedCalled

		children=append(children, childrenSuccess)
		children=append(children, childrenFail)
		parent["children"]=children

		dataList=append(dataList,parent)




		////查询api的defaultQuotas
		//defaultQuotas := models.DefaultQuotasModel{}
		//_, _ = defaultQuotas.QueryTable().Filter("api_id", api.Id).All(&defaultQuotas)
		////
		//var apiModelTemp = &models.ApiModelTemp{
		//	Id:           api.Id,
		//	ApiName:      api.ApiName,
		//	Url:          api.Url,
		//	MaxQps:       defaultQuotas.MaxQps,
		//	Quotas:       defaultQuotas.Quotas,
		//	QuotasPeriod: defaultQuotas.QuotasPeriod,
		//	CreateTime:   api.CreateTime.Format("2006-01-02 15:04:05"),
		//	ModifyTime:   api.ModifyTime.Format("2006-01-02 15:04:05"),
		//}
		//apiModelTempList = append(apiModelTempList, apiModelTemp)
	}
	return dataList
}
//}
