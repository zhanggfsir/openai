package services

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"openai-backend/models"
	"openai-backend/views"
)

/*
@Author:
@Time: 2020-11-05 15:02
@Description:apiService
*/

//新建ability数据校验
func
ValidatePostApi(requestBody []byte) (*models.ApiModel, error) {
	//获取requestBody中的参数
	type requestJSON struct {
		ApiName      string `description:"api名称" json:"apiName"`
		Url          string `description:"url" json:"url"`
		AbilityId    uint64 `description:"能力ID" json:"abilityId"`
		MaxQps       uint64 `description:"最大Qps" json:"maxQps"`
		Quotas       uint64 `description:"配额" json:"quotas"`
		QuotasPeriod string `description:"配额单位 h d m y" json:"quotasPeriod"`
	}
	var jSON requestJSON

	if err := json.Unmarshal(requestBody, &jSON); err != nil {
		logs.Error("err=", err)
		return nil, errors.New("json 格式错误" + err.Error())
	}
	//校验数据
	valid := validation.Validation{}
	valid.Required(jSON.ApiName, "categoryName").Message("未输入名称 ")
	valid.Required(jSON.Url, "url").Message("未输入url ")
	valid.Required(jSON.AbilityId, "abilityId").Message("未输入能力ID ")
	valid.Required(jSON.MaxQps, "maxQps").Message("未输入最大Qps ")
	valid.Required(jSON.Quotas, "quotas").Message("未输入配额 ")
	valid.Required(jSON.QuotasPeriod, "quotasPeriod").Message("未输入配额单位 ")

	if valid.HasErrors() {
		var errStr string
		for _, err := range valid.Errors {
			errStr = errStr + err.Message
		}
		//	logs.Error("err", errStr)
		return nil, errors.New(errStr)

	}

	//将json数据复制给model
	model := &models.ApiModel{
		ApiName:   jSON.ApiName,
		AbilityId: jSON.AbilityId,
		Url:       jSON.Url,
	}

	if model.QueryTable().Filter("api_name", model.ApiName).Exist() == true {
		return nil, errors.New("名称重复")
	}

	if models.NewAbilityModel().QueryTable().Filter("id", model.AbilityId).Exist() != true {
		return nil, errors.New("能力不存在")
	}

	//api插入数据库
	if err := model.Insert(); err != nil {
		logs.Error("err= ", err)
		return nil, err
	}

	//defaultQuotas
	defaultQuotas := &models.DefaultQuotasModel{
		ApiId:        model.Id,
		MaxQps:       jSON.MaxQps,
		Quotas:       jSON.Quotas,
		QuotasPeriod: jSON.QuotasPeriod,
	}
	//defaultQuotas插入数据库
	if err := defaultQuotas.Insert(); err != nil {
		logs.Error("err= ", err)
		return nil, err
	}

	return model, nil
}

//获取api列表
func GetApiListService(queryFilter map[string]interface{}) (map[string]interface{}, error) {
	model := &models.AbilityCategoryModel{}
	lists, Count, err := model.AbilityCategoryList(queryFilter)
	if err != nil {
		logs.Error("data group lists error", err)
		return nil, err
	}
	//logs.Info(lists, Count)

	datalist := views.AbilityCategoryList(lists)
	//logs.Info(datalist)
	var jsonData = make(map[string]interface{})
	jsonData["list"] = datalist
	jsonData["total"] = Count
	return jsonData, nil
}
