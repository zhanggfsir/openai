package services

import (
	"openai-backend/models"
	"openai-backend/utils"
	"openai-backend/utils/httpUtil"

	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
)

func AddRouterValidate(domain string, body []byte) (*models.RouterModel, error) {
	var (
		err error
		routerInfo models.RouterModel
	)
	if err := json.Unmarshal(body, &routerInfo); err != nil {
		logs.Error("err=", err)
		return nil, errors.New("json格式错误"+err.Error())
	}

	valid := validation.Validation{}

	valid.Required(routerInfo.Path, "path").Message("请输入路由path,")
	valid.Required(routerInfo.Method, "method").Message("请输入请求方法,")
	valid.Required(routerInfo.RouterGroup, "routerGroup").Message("请输入路由组,")
	valid.Required(routerInfo.Description, "description").Message("请输入路由描述")

	if valid.HasErrors() {
		var errStr string
		for _,err := range valid.Errors {
			errStr = errStr + err.Message
		}
		logs.Error("router insert err", utils.Fields{
			"body": body,
			"err": errStr,
		})

		return nil, errors.New(errStr)
	}

	// 判断path + method是否存在
	if _, err = models.NewRouterModel().FetchByPathAndMethod(routerInfo.Path, routerInfo.Method); err == nil {
		return nil, errors.New("路径和方法已经存在")
	}
	routerInfo.Domain = domain
	return &routerInfo, nil
}

func PutRouterValidate(domain string, body []byte) (map[string]interface{}, int, error) {
	var (
		err error
		requestInfo models.RouterModel
	)
	//获得body中的参数
	if err = json.Unmarshal(body, &requestInfo); err != nil {
		logs.Error("parse json err", err)
		return nil, 0, errors.New(httpUtil.ERROR_JSON_PATTERN)
	}

	valid := validation.Validation{}

	valid.Required(requestInfo.Path, "path").Message("请输入路由path,")
	valid.Required(requestInfo.Method, "method").Message("请输入请求方法,")
	valid.Required(requestInfo.RouterGroup, "routerGroup").Message("请输入路由组,")
	valid.Required(requestInfo.Description, "description").Message("请输入路由描述")

	if valid.HasErrors() {
		var errStr string
		for _,err := range valid.Errors {
			errStr = errStr + err.Message
		}
		logs.Error("router insert err", utils.Fields{
			"body": body,
			"err": errStr,
		})

		return nil, 0, errors.New(errStr)
	}

	// 判断path + method是否存在
	if routerInfo, err := models.NewRouterModel().FetchByPathAndMethod(requestInfo.Path, requestInfo.Method);
	err == nil && requestInfo.Id != routerInfo.Id{
		return nil, 0, errors.New("路径和方法已经存在")
	}

	resultMap := map[string]interface{}{
		"path":         requestInfo.Path,
		"router_group": requestInfo.RouterGroup,
		"method":       requestInfo.Method,
		"description":  requestInfo.Description,
	}

	return resultMap, requestInfo.Id, nil
}