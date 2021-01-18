package services

import (
	"openai-backend/models"
	"openai-backend/utils"
	"encoding/base64"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/tidwall/gjson"
	"strings"
)

func ChildUserValidate(domain string, body []byte) (*models.ChildUserModel, error) {
	var (
		err error
		userTypeInt int64
		passwordDecode []byte
		appKey, appSecret, savePd string
	)

	userName := gjson.GetBytes(body, "userName").String()
	password := gjson.GetBytes(body, "password").String()
	resetPd := gjson.GetBytes(body, "resetPd").String()
	remark := gjson.GetBytes(body, "remark").String()
	//roleIds := gjson.GetBytes(body, "roleIds").String()
	userType := gjson.GetBytes(body, "userType").String()
	valid := validation.Validation{}

	valid.Required(userName, "userName").Message("请输入用户名,")
	valid.Required(userType, "userType").Message("请输入用户类型,")
	valid.Required(resetPd, "resetPd").Message("请输入是否重置密码,")
	//valid.Required(roleIds, "roleIds").Message("请输入角色,")

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


	if userTypeInt, appKey, appSecret, err = getChildUserType(userType); err != nil {
		return nil, err
	}

	if userTypeInt >0 {
		if password == "" {
			return nil, errors.New("请输入密码")
		}
		if passwordDecode, err = base64.StdEncoding.DecodeString(password); err != nil {
			return nil, errors.New("base密码错误")
		}
		savePd = utils.GenPassword(string(passwordDecode))
 	}

	// 判断用户名是否存在 是否存在
	if _, err := models.NewChildUserModel().FetchByUserNameDomain(domain, userName); err == nil {
		logs.Error("child user insert err", utils.Fields{
			"body": body,
			"err": "userName已经存在",
		})

		return nil, errors.New("用户名已经存在")
	}

	childUserInfo := &models.ChildUserModel{
		UserId: utils.GenerateUuid(),
		UserName: userName,
		Password: savePd,
		AppKey: appKey,
		AppSecret: appSecret,
		Domain: domain,
		Remark: remark,
		UserType: userTypeInt,
		UserTypeStr: userType,
		//RoleIds: roleIds,
		ResetPd: resetPd,
		DelFlag: "0",
	}
	return childUserInfo, nil
}

func getChildUserType(userType string) (int64, string, string, error) {
	var (
		appKey = utils.GenerateUuid()
		appSecret = utils.GenerateUuid()
	)
	splitStr := strings.Split(userType, ",")
	length := len(splitStr)
	if length == 1 {
		loginType := splitStr[0]
		if loginType == "0" {
			return models.ChildTypeProgram, appKey, appSecret,  nil
		} else if loginType == "1" {
			return models.ChildTypeConsole, "", "", nil
		} else {
			return 0, "", "", errors.New("用户类型错误")
		}
	} else if length == 2 {
		return models.ChildTypeAll, appKey, appSecret,  nil
	} else {
		return 0, "", "", errors.New("用户类型错误")
	}
}

func DoChildRoleManager(domain string, body []byte) ([]string, []string, string, string, error) {
	var (
		err error
		roleSaveMap = make(map[string]string)
		roleReqMap = make(map[string]string)
		roleListStr string
		addList = make([]string, 0)
		delList = make([]string, 0)
	)

	userId := gjson.GetBytes(body, "id").String()
	roleIds := gjson.GetBytes(body, "roleIds")
	if !roleIds.IsArray() {
		return nil, nil, "","",errors.New("角色列表错误")
	}

	if roleSaveMap, err = getSaveRoleList(userId); err != nil {
		return nil, nil, "", "",err
	}

	for _, roleId := range roleIds.Array() {
		roleIdStr := roleId.String()
		if roleListStr == "" {
			roleListStr = roleListStr + roleIdStr
		} else {
			roleListStr = roleListStr + "," + roleIdStr
		}
		roleReqMap[roleIdStr] = "1"
		if _, ok := roleSaveMap[roleIdStr]; !ok {
			addList = append(addList, roleIdStr)
		}
	}

	for roleId, _ := range roleSaveMap {
		if _, ok := roleReqMap[roleId]; !ok {
			delList = append(delList, roleId)
		}
	}

	return addList, delList, userId , roleListStr,nil
}

func getSaveRoleList(userId string) (map[string]string, error) {
	var (
		err error
		userRoleInfos []*models.UserRoleModel
		roleMap = make(map[string]string)
	)
	if userRoleInfos, err = models.NewUserRoleModel().FetchByUserId(userId); err != nil {
		logs.Error("fetch user role list", utils.Fields{
			"userId": userId,
			"err": err,
		})
		return nil, err
	}

	for _, userRoleInfo := range userRoleInfos {
		roleMap[userRoleInfo.RoleId]= "1"
	}

	return roleMap, nil
}