package casbinUtil

import (
	"openai-backend/utils"
	"errors"
	"github.com/astaxie/beego/logs"
	"net/http"
)

/*
@Author:
@Time: 2020-07-31 10:02
@Description: casbin
*/

type GroupRole struct {
	userId string     //用户ID
	roleId string     //roleId
	domain string     //域名
}

var (
	defaultDomain = "default"
	RoleAdmin    = "admin"
	//RolesId       = map[string] string{
	//	RoleAdmin:     "系统管理员",
	//}
	DefaultActionList = []string{"add","read","update","delete"}
	DefaultPolicyList = []string{"account","token","datasetGroup","dataset",
		"annotate","label","device","trainingTask","role","childUser","datafileupload"}
)

//获得资源的方法
func CasbinAction(method string) string  {
	switch method {
	case http.MethodGet:
		return "read"
	case http.MethodPost:
		return "add"
	case http.MethodPut:
		return "update"
	case http.MethodDelete:
		return "delete"
	default:
		return  "other"
	}
}

func AddCasbin(cm CasbinRule) bool {
	success, err := Enforcer.AddPolicy(cm.V0, cm.V1, cm.V2, cm.V3)
	if err != nil {
		logs.Error("Add err=", err, cm)
	}
	return success
}

//获得casbin中所有的用户/角色/域名
func FetchAllGroup() ([]GroupRole, error) {
	var listRole []GroupRole
	GetGroupingPolicy := Enforcer.GetGroupingPolicy()
	for _, value := range GetGroupingPolicy {
		if len(value) == 3 {
			tempValue := GroupRole{
				userId: value[0],
				roleId: value[1],
				domain: value[2],
			}
			listRole = append(listRole, tempValue)
		} else {
			return nil, errors.New("错误")
		}
	}
	return listRole, nil
}

func AddRolesGroupPolicy(domain, userId, roleId string) error {
	// 添加admin用户role权限
	if flag, err := Enforcer.AddGroupingPolicy(userId, roleId, domain); err != nil && flag {
		logs.Error("AddRolesGroupPolicy=", utils.Fields{
			"flag": flag,
			"domain": domain,
			"userId": userId,
		})
		return errors.New("err")
	}

	return nil
}

func DelRolesGroupPolicy(domain, userId, roleId string) error {
	// 添加admin用户role权限
	if flag, err := Enforcer.RemoveGroupingPolicy(userId, roleId, domain); err != nil && flag {
		logs.Error("DelRolesGroupPolicy", utils.Fields{
			"flag": flag,
			"domain": domain,
			"userId": userId,
		})

		return errors.New("err")
	}

	return nil
}

func AuthorizationDomain(userId, source, method, domain string) bool {
	e := Enforcer
	status, err := e.Enforce(userId, domain, source, method)
	if status {
		return true
	} else {
		logs.Error("first auth=", utils.Fields{
			"err": err,
			"userId": userId,
			"domain": domain,
			"source": source,
			"method": method,
		})

		return false
	}
}

func AuthorizationDomainAgain(userId, source, method, domain string) bool {
	if err := Enforcer.LoadPolicy(); err != nil {
		logs.Error("casbin load policy= ", utils.Fields{
			"err": err,
			"userId": userId,
			"domain": domain,
			"source": source,
			"method": method,
		})

		return true
	}

	if status, err := Enforcer.Enforce(userId, domain, source, method); !status {
		logs.Error("second auth", utils.Fields{
			"err": err,
			"userId": userId,
			"domain": domain,
			"source": source,
			"method": method,
		})

		return false
	}

	return true
}