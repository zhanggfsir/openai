package services

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"openai-backend/models"
)

/*
@Author:
@Time: 2020-11-05 15:02
@Description:abilityService
*/

//新建ability数据校验
func ValidatePostAbility(requestBody []byte) (*models.AbilityModel, error) {
	//获取requestBody中的参数
	type requestJSON struct {
		AbilityName     string `description:"Ability名称" json:"abilityName"`
		CategoryId     uint64 `description:"能力种类ID" json:"categoryId"`
	}
	var jSON requestJSON

	if err := json.Unmarshal(requestBody, &jSON); err != nil {
		logs.Error("err=", err)
		return nil, errors.New("json 格式错误" + err.Error())

	}
	//校验数据
	valid := validation.Validation{}
	valid.Required(jSON.AbilityName, "categoryName").Message("未输入名称 ")
	valid.Required(jSON.CategoryId, "categoryId").Message("未输入能力种类ID ")

	if valid.HasErrors() {
		var errStr string
		for _, err := range valid.Errors {
			errStr = errStr + err.Message
		}
		//	logs.Error("err", errStr)
		return nil, errors.New(errStr)

	}

	//将json数据复制给model
	model := &models.AbilityModel{
		AbilityName:     jSON.AbilityName,
		CategoryId:  jSON.CategoryId,

	}

	if model.QueryTable().Filter("ability_name", model.AbilityName).Exist() == true {
		return nil, errors.New("名称重复")
	}

	if models.NewAbilityCategoryModel().QueryTable().Filter("id", model.CategoryId).Exist() != true {
		return nil, errors.New("能力种类不存在")
	}

	return model, nil
}