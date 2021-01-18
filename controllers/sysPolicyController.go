package controllers

import (
	"openai-backend/models"
	"openai-backend/services"
	"openai-backend/utils"
	"openai-backend/utils/httpUtil"
	"openai-backend/views"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

/*
@Author:
@Time: 2020-07-30 10:02
@Description: 策略管理
*/

type PolicyController struct {
	BaseController
}

func (this *PolicyController) Post() {
	var (
		err error
		body map[string]string
		policyInfo *models.PolicyModel
	)

	//获得body中的参数
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &body); err != nil {
		this.JsonResult(httpUtil.ERROR_JSON_PATTERN, nil)
		return
	}

	if policyInfo, err = services.PolicyAddValidate(body); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if err = policyInfo.Insert(); err != nil {
		logs.Error("source insert err", utils.Fields{
			"body": body,
			"err": err.Error(),
		})
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

func (this *PolicyController) Put() {
	var (
		err error
		id string
		body map[string]string
		restData map[string]interface{}
	)

	//获得body中的参数
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &body); err != nil {
		this.JsonResult(httpUtil.ERROR_JSON_PATTERN, nil)
		return
	}

	if restData, id, err = services.PolicyUpdateValidate(body); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if err = models.NewPolicyModel().SaveFilter(id, restData); err != nil {
		logs.Error("policy update err", utils.Fields{
			"body": body,
			"err": err.Error(),
		})
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

func (this *PolicyController) Delete() {
	var contentBody map[string][]string
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &contentBody)
	if err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "json格式错误")

		return
	}

	ids := contentBody["ids"]
	err = models.NewPolicyModel().DeleteById(ids)
	if err != nil {
		logs.Error("policy delete", utils.Fields{
			"ids": ids,
			"err": err.Error(),
		})

		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

func (this *PolicyController) GetList() {
	var (
		err error
		pageNo, pageSize, totalCount int
		list []*models.PolicyModel
		queryFilter = make(map[string]string)
		restData = make(map[string]interface{})
	)

	if pageNo, err = this.GetInt("pageNo"); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "请输入pageNo")

		return
	}
	if pageSize, err = this.GetInt("pageSize"); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "请输入pageSize")

		return
	}

	if value := this.GetString("policyGroup"); value != "" {
		queryFilter["policy_group"] = value
	}

	if list, totalCount, err = models.NewPolicyModel().FindToPager(pageNo, pageSize, queryFilter); err != nil {
		logs.Error("policy list mysql= ", utils.Fields{
			"err": err.Error(),
		})

		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	datalist := views.PolicyList(list)
	restData["list"] = datalist
	restData["pageNo"] = pageNo
	restData["pageSize"] = pageSize
	restData["total"] = totalCount

	this.JsonResult(httpUtil.SUCCESS, restData)
}

