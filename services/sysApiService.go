package services

import (
	"openai-backend/models"
	"openai-backend/models/request"
	"openai-backend/utils/casbinUtil"
	"errors"
	"github.com/astaxie/beego/logs"
)

// SerGetRoleApis 获得api list
func SerGetRoleApis(domain, roleId string) ([]map[string]interface{}, error) {
	var (
		err error
		apiList []*models.RouterModel
		//apiResponses []*models.ApiResponse
		casbinMap = make(map[string]string)
	)
	casbinList := casbinUtil.Enforcer.GetFilteredPolicy(0, roleId, domain)
	// value [admin 50604129 account add]
	for _, value :=  range casbinList {
		key := value[2]+value[3]
		casbinMap[key] = "1"
	}

	if apiList, err = models.NewRouterModel().FetchAllRouters(map[string]string{}); err != nil {
		return nil, err
	}

	apiResponses :=make([]map[string]interface{},0)
	for _, info := range apiList {
		var flag = "0"
		if _, status := casbinMap[info.Path+info.Method]; status {
			flag = "1"
		}

		apiResponse := make(map[string]interface{})
		apiResponse["path"] = info.Path
		apiResponse["routerGroup"] = info.RouterGroup
		apiResponse["method"] = info.Method
		apiResponse["description"] = info.Description
		apiResponse["domain"] = info.Domain
		apiResponse["flag"] = flag
		//apiResponse := &models.ApiResponse{
		//	Path : info.Path,
		//	RouterGroup :info.RouterGroup,
		//	Method :     info.Method,
		//	Description  :info.Description,
		//	Domain :info.Domain,
		//	Flag  :flag,
		//}

		apiResponses = append(apiResponses, apiResponse)
	}

	return apiResponses, nil
}

func UpdateCasbin(domain, roleId string, casbinInfos []request.CasbinInfo) error {
	// 清空casbin
	//if !ClearCasbin(domain, roleId) {
	//	return errors.New("存在相同api,添加失败,请联系管理员")
	//}

	ClearCasbin(domain, roleId)

	// 添加casbin
	for _, v := range casbinInfos {
		cm := casbinUtil.CasbinRule{
			PType: "p",
			V0:    roleId,
			V1:    domain,
			V2:    v.Path,
			V3:    v.Method,
			V4:    "",
			V5:    "",
		}
		addFlag := casbinUtil.AddCasbin(cm)
		if !addFlag {
			return errors.New("存在相同api,添加失败,请联系管理员")
		}
	}
	return nil
}

func ClearCasbin(domain, roleId string) bool {
	success, err := casbinUtil.Enforcer.RemoveFilteredPolicy(0, roleId, domain)
	if err != nil {
		logs.Error("ClearCasbin err=", err, domain, roleId)
	}
	return success
}
