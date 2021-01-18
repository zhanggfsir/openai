package views

import (
	"openai-backend/models"
)

func RouterViewList(list []*models.RouterModel) []map[string]interface{} {
	var restData = make([]map[string]interface{}, 0)
	for _, tempValue := range list {
		var tempData = make(map[string]interface{})
		tempData["id"] = tempValue.Id
		tempData["path"] = tempValue.Path
		tempData["mothod"] = tempValue.Method
		tempData["routerGroup"] = tempValue.RouterGroup
		tempData["description"] = tempValue.Description
		tempData["domain"] = tempValue.Domain
		tempData["createTime"] = tempValue.CreateTime.Format("2006-01-02 15:04:05")
		restData = append(restData, tempData)
	}
	return restData
}
