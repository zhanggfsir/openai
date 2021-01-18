package services

import (
	"openai-backend/models"
	"openai-backend/utils"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"errors"
)

func PolicyAddValidate(body map[string]string) (*models.PolicyModel,error) {
	policyGroup := body["policyGroup"]
	action := body["action"]
	url := body["url"]
	remark := body["remark"]
	valid := validation.Validation{}

	valid.Required(policyGroup, "policyGroup").Message("请输入策略组,")
	valid.Required(action, "action").Message("请输入操作,")
	valid.Required(url, "url").Message("请输入请求的url,")

	if valid.HasErrors() {
		var errStr string
		for _,err := range valid.Errors {
			errStr = errStr + err.Message
		}
		logs.Error("policy insert err", utils.Fields{
			"body": body,
			"err": errStr,
		})

		return nil, errors.New(errStr)
	}

	// 判断sourceID 是都存在
	if _, err := models.NewPolicyModel().FetchById(policyGroup+action); err == nil {
		logs.Error("source insert err", utils.Fields{
			"body": body,
			"err": "资源Id已经存在",
		})

		return nil, errors.New("资源Id已经存在")
	}
	policyInfo := &models.PolicyModel{
		Id: policyGroup+action,
		PolicyGroup: policyGroup,
		Action: action,
		Url: url,
		Remark: remark,
	}
	return policyInfo, nil
}

func PolicyUpdateValidate(body map[string]string) (map[string]interface{}, string, error) {
	id := body["id"]
	url := body["url"]
	action := body["action"]
	remark := body["remark"]
	policyGroup := body["policyGroup"]
	valid := validation.Validation{}

	valid.Required(id, "id").Message("请输入资源Id,")
	valid.Required(url, "url").Message("请输入请求的url,")
	valid.Required(action, "action").Message("请输入操作,")
	valid.Required(policyGroup, "policyGroup").Message("请输入策略组,")

	if valid.HasErrors() {
		var errStr string
		for _,err := range valid.Errors {
			errStr = errStr + err.Message
		}
		logs.Error("source update err", utils.Fields{
			"body": body,
			"err": errStr,
		})

		return nil, "", errors.New(errStr)
	}

	resultMap := map[string]interface{}{
		"policy_group": policyGroup,
		"url": url,
		"action": action,
		"remark": remark,
	}

	return resultMap, id , nil
}

func PolicyListValidate(header map[string]string) (int, int, map[string]string, error) {
	var (
		err error
		pageNoInt, pageSizeInt int
		queryFilter = make(map[string]string)
	)

	pageNo := header["pageNo"]
	pageSize := header["pageSize"]
	policyGroup := header["policyGroup"]
	valid := validation.Validation{}

	valid.Required(pageSize, "pageSize").Message("请输入pageSize,")
	valid.Required(pageNo, "pageNo").Message("请输入pageNo,")

	if valid.HasErrors() {
		var errStr string
		for _,err := range valid.Errors {
			errStr = errStr + err.Message
		}
		logs.Error("source list err", utils.Fields{
			"header": header,
			"err": errStr,
		})

		return 0, 0, nil, errors.New(errStr)
	}

	if pageNoInt, err = utils.A2Int(pageNo); err != nil {
		logs.Error("source list err", utils.Fields{
			"pageNo": pageNo,
			"err": err.Error(),
		})

		return 0, 0, nil, errors.New("pageNo数字")
	}

	if pageSizeInt, err = utils.A2Int(pageSize); err != nil {
		logs.Error("source list err", utils.Fields{
			"pageSize": pageSize,
			"err": err.Error(),
		})

		return 0, 0, nil, errors.New("pageSize数字")
	}

	if policyGroup != "" {
		queryFilter["policy_group"] = policyGroup
	}

	return pageNoInt, pageSizeInt, queryFilter, nil
}