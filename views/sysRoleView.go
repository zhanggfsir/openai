package views

import "openai-backend/models"

func RoleList(list []*models.RoleModel) []map[string]interface{} {
	var restData = make([]map[string]interface{}, 0)
	for _, tempValue := range list {
		var tempData = make(map[string]interface{})
		tempData["id"] = tempValue.Id
		tempData["roleId"] = tempValue.RoleId
		tempData["name"] = tempValue.RoleName
		tempData["remark"] = tempValue.Remark
		tempData["updateTime"] = tempValue.ModifyTime.Format("2006-01-02 15:04:05")
		tempData["createTime"] = tempValue.CreateTime.Format("2006-01-02 15:04:05")
		restData = append(restData, tempData)
	}

	return restData
}
