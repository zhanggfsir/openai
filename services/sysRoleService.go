package services

import (
	"openai-backend/models"
	"openai-backend/utils"
	"openai-backend/utils/casbinUtil"
	"openai-backend/utils/httpUtil"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

type ApiStruct struct {
	Op string `json:"op"`
	Flag string `json:"flag"`
}

func RoleAddValidate(domain string, body map[string]string) (*models.RoleModel,error) {
	roleId := body["roleId"]
	name := body["name"]
	remark := body["remark"]
	valid := validation.Validation{}

	valid.Required(roleId, "roleId").Message("请输入角色,")
	valid.Match(roleId, regexp.MustCompile(`^[a-z]{1,8}`),"roleId").Message("角色只能是小写字母")
	valid.Required(name, "name").Message("请输入角色名称,")

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

	// 判断role_id 是否存在
	if _, err := models.NewRoleModel().FetchById(domain+roleId); err == nil {
		logs.Error("source insert err", utils.Fields{
			"body": body,
			"err": "资源Id已经存在",
		})

		return nil, errors.New("角色Id已经存在")
	}
	roleInfo := &models.RoleModel{
		Id: domain+roleId,
		RoleId: roleId,
		RoleName: name,
		Domain: domain,
		Remark: remark,
	}
	return roleInfo, nil
}

func RoleUpdateValidate(body map[string]string) (map[string]interface{}, string, error) {
	id := body["id"]
	name := body["name"]
	remark := body["remark"]
	valid := validation.Validation{}

	valid.Required(id, "id").Message("请输入角色,")
	valid.Required(name, "name").Message("请输入角色名称,")

	if valid.HasErrors() {
		var errStr string
		for _,err := range valid.Errors {
			errStr = errStr + err.Message
		}
		logs.Error("policy insert err", utils.Fields{
			"body": body,
			"err": errStr,
		})

		return nil, "", errors.New(errStr)
	}

	resultMap := map[string]interface{}{
		"remark": remark,
		"role_name": name,
	}

	return resultMap, id , nil
}

func DoRoleAuthList(domain, userId, userType string) (map[string][]string, []string, error) {
	var (
		err error
		roleList []string
		adminBool = false
		userAndRoleList [] *models.UserRoleModel
	)

	if userAndRoleList, err = models.NewUserRoleModel().FetchByUserId(userId); err != nil || len(userAndRoleList) == 0{
		return nil, nil, err
	}

	for _, userRoleInfo := range userAndRoleList {
		if userRoleInfo.RoleId == casbinUtil.RoleAdmin {
			adminBool = true
		}
		roleList = append(roleList, userRoleInfo.RoleId)
	}

	if adminBool {
		permissions, err := doGetAuthByDB()
		return permissions, roleList, err
	}

	if userType == "1" {
		permissions, err := doGetAuthByDB()
		return permissions, roleList, err
	}

	return doGetAuthByCasbin(domain, userAndRoleList), roleList, err
}

//TODO 更新role api list
func DoUpdateApiRoleAuthList(domain string, contentBody []byte) error {
	var (
		err error
		roleId string
		addApiMap, delApiMap   []string
		reqApiMap map[string]string
		saveApiMap = make(map[string]string)
	)

	if reqApiMap, roleId, err = parseUpdateApiRoleAuthJson(contentBody); err != nil {
		return err
	}

	DoGetRoleAuthByCasbin(domain, roleId, saveApiMap)

	for reqKey, _ := range reqApiMap {
		if _, ok := saveApiMap[reqKey]; !ok {
			addApiMap = append(addApiMap, reqKey)
		}
	}

	for saveKey, _ := range saveApiMap {
		if _, ok := reqApiMap[saveKey]; !ok {
			delApiMap = append(delApiMap, saveKey)
		}
	}

	for _, value := range addApiMap {
		splitKey := strings.Split(value, "-")
		policy := splitKey[0]
		action := splitKey[1]
		if _, err = casbinUtil.Enforcer.AddPolicy(roleId, domain, policy, action); err != nil {
			goto ERR
		}
	}

	for _, value := range delApiMap {
		splitKey := strings.Split(value, "-")
		policy := splitKey[0]
		action := splitKey[1]
		if _, err = casbinUtil.Enforcer.RemoveFilteredPolicy(0, roleId, domain, policy, action); err != nil {
			goto ERR
		}
	}

	return err

ERR:
	logs.Error("update api role auth list", utils.Fields{
		"domain" : domain,
		"roleId": roleId,
		"err": err.Error(),
		"contentBody": string(contentBody),
	})

	return errors.New(httpUtil.SYSTEM_ERROR)
}

//func parseUpdateApiRoleAuthJson(contentBody []byte) (map[string]string, string, error) {
//	mapData := make(map[string]string)
//
//	roleId := gjson.GetBytes(contentBody, "roleId").String()
//	if roleId == "" {
//		return nil, "", errors.New(httpUtil.PARAMETER_ERROR)
//	}
//
//	policyList := gjson.GetBytes(contentBody, "policyList")
//	flag := policyList.IsArray()
//	if !flag {
//		return nil, "", errors.New(httpUtil.PARAMETER_ERROR)
//	}
//	for _, value := range policyList.Array() {
//		policy := value.Map()["module"].String()
//		actionList := value.Map()["operation"]
//		if !actionList.IsArray() {
//			return nil, "", errors.New(httpUtil.PARAMETER_ERROR)
//		}
//
//		for _, action := range actionList.Array() {
//			key := policy + "-" + action.String()
//			mapData[key] = "1"
//		}
//
//	}
//	return mapData, roleId, nil
//}


func parseUpdateApiRoleAuthJson(contentBody []byte) (map[string]string, string, error) {
	mapData := make(map[string]string)

	roleId := gjson.GetBytes(contentBody, "roleId").String()
	if roleId == "" {
		return nil, "", errors.New(httpUtil.PARAMETER_ERROR)
	}

	policyMap := gjson.GetBytes(contentBody, "policyList").Map()
	for policy, value := range policyMap {
		if !value.IsArray() {
			return nil, "", errors.New(httpUtil.PARAMETER_ERROR)
		}
		for _, action := range value.Array() {
			key := policy + "-" + action.String()
			mapData[key] = "1"
		}
	}

	return mapData, roleId, nil
}


func doGetAuthByDB() (map[string][]string, error) {
	var(
		err error
		list []*models.PolicyModel
		roleListMap = make(map[string][]string)
	)
	if list, err = models.NewPolicyModel().FetchAllInfo(); err != nil {
		return nil, err
	}

	for _, value := range list {
		if data, ok := roleListMap[value.PolicyGroup]; ok {
			roleListMap[value.PolicyGroup] = append(data, value.Action)
		} else {
			roleListMap[value.PolicyGroup] = []string{value.Action}
		}
	}

	return roleListMap, nil
}

func doGetAuthByCasbin(domain string, list []*models.UserRoleModel) map[string][]string {
	var roleListMap = make(map[string][]string)
	for _, userRoleInfo := range list {
		for _, casbin := range casbinUtil.Enforcer.GetFilteredPolicy(0, userRoleInfo.RoleId, domain) {
			mode := casbin[2]
			action := casbin[3]
			if data, ok := roleListMap[mode]; ok {
				roleListMap[mode] = append(data, action)
			} else {
				roleListMap[mode] = []string{action}
			}
		}
	}

	return roleListMap
}

func GetRoleApiList(RequestRoleId string, jwtUser utils.JwtUser) map[string][]ApiStruct {
	requestApiList := make(map[string]string)
	domain := jwtUser.Domain
	hostApiList, _ := getUserAuthList(domain, jwtUser.Id)
	DoGetRoleAuthByCasbin(domain, RequestRoleId, requestApiList)

	for key, _ := range hostApiList {
		if _, ok := requestApiList[key]; ok {
			hostApiList[key] = "1"
		}
	}

	return genApiList(hostApiList)
}

func genApiList(apiList map[string]string) map[string][]ApiStruct {
	retMap := make(map[string][]ApiStruct)
	for key, flag := range apiList {
		splitKey := strings.Split(key, "-")
		source := splitKey[0]
		action := splitKey[1]
		api := ApiStruct{
			Op:action,
			Flag: flag,
		}
		if data, ok := retMap[source]; ok {
			data = append(data, api)
			retMap[source] = data
		} else {
			retMap[source] = []ApiStruct{api}
		}
	}
	return retMap
}

//获得userId的所有的权限
func getUserAuthList(domain, userId string) (map[string]string, error) {
	var (
		err error
		userRoleList []*models.UserRoleModel
	)
	if userRoleList, err = models.NewUserRoleModel().FetchByUserId(userId); err != nil || len(userRoleList) == 0{
		return nil, err
	}

	var roleListMap = make(map[string]string)
	for _, userRoleInfo := range userRoleList {
		DoGetRoleAuthByCasbin(domain, userRoleInfo.RoleId, roleListMap)
	}

	return roleListMap, nil
}

func DoGetRoleAuthByCasbin(domain, roleId string, roleListMap map[string]string){
	for _, casbin := range casbinUtil.Enforcer.GetFilteredPolicy(0, roleId, domain) {
		key := casbin[2] + "-" + casbin[3]
		roleListMap[key] = "0"
	}
}

