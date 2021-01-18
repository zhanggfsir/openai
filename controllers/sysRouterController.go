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
@Time: 2020-11-30 14:02
@Description: router管理
*/

type RouterController struct {
	BaseController
}

// Post 添加router
func (this *RouterController) Post() {
	var (
		err error
		routerInfo *models.RouterModel
	)

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	if routerInfo, err = services.AddRouterValidate(domain, this.Ctx.Input.RequestBody); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if err = routerInfo.Insert(); err != nil {
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

// Delete 删除路由
func (this *RouterController) Delete() {
	var contentBody models.DelRequest
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &contentBody)
	if err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "json格式错误")
		return
	}

	ids := contentBody.Id
	if err = models.NewRouterModel().DeleteById(ids); err != nil {
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

// Put 修改信息
func (this *RouterController) Put() {
	var (
		id int
		err error
		domain string
		updateData map[string]interface{}
	)

	if domain, err = this.getDomain(); err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	if updateData, id, err = services.PutRouterValidate(domain, this.Ctx.Input.RequestBody); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if err = models.NewRouterModel().SaveById(id, updateData); err != nil {
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

// RouterList 获得路由列表
func (this *RouterController) RouterList() {
	var (
		err error
		domain string
		pageNo, pageSize, totalCount int
		list []*models.RouterModel
		queryFilter = make(map[string]string)
		restData = make(map[string]interface{})
	)
	// 获得请求参数
	if domain, err = this.getDomain(); err != nil || domain == "" {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	if pageNo, err = this.GetInt("pageNo"); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "请输入pageNo")
		return
	}
	if pageSize, err = this.GetInt("pageSize"); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "请输入pageSize")
		return
	}

	if value := this.GetString("path"); value != "" {
		queryFilter["path"] = value
	}
	if value := this.GetString("routerGroup"); value != "" {
		queryFilter["router_group"] = value
	}
	//queryFilter["domain"] = domain

	if list, totalCount, err = models.NewRouterModel().FindToPager(pageNo, pageSize, queryFilter); err != nil {
		logs.Error("child user list mysql= ", utils.Fields{
			"err": err.Error(),
		})

		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	datalist := views.RouterViewList(list)
	restData["list"] = datalist
	restData["pageNo"] = pageNo
	restData["pageSize"] = pageSize
	restData["total"] = totalCount

	this.JsonResult(httpUtil.SUCCESS, restData)
}