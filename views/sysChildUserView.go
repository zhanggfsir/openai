package views

import (
	"openai-backend/models"
	"strings"
)

func ChildUserList(list []*models.ChildUserModel) []map[string]interface{} {
	var restData = make([]map[string]interface{}, 0)
	for _, tempValue := range list {
		var tempData = make(map[string]interface{})
		tempData["id"] = tempValue.UserId
		tempData["roleIds"] = strings.Split(tempValue.RoleIds, ",")
		tempData["userType"] = tempValue.UserTypeStr
		tempData["appKey"] = tempValue.AppKey
		tempData["appSecret"] = tempValue.AppSecret
		tempData["userName"] = tempValue.UserName
		tempData["resetPd"] = tempValue.ResetPd
		tempData["createTime"] = tempValue.CreateTime.Format("2006-01-02 15:04:05")
		restData = append(restData, tempData)
	}

	return restData
}
