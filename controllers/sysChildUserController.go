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
@Time: 2020-08-24 14:02
@Description: 子账号管理
*/

type ChildUserController struct {
	BaseController
}

func (this *ChildUserController) Post() {
	var (
		err error
		childUserInfo *models.ChildUserModel
	)

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	if childUserInfo, err = services.ChildUserValidate(domain, this.Ctx.Input.RequestBody); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if err = childUserInfo.Insert(); err != nil {
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

func (this *ChildUserController) GetList() {
	var (
		err error
		pageNo, pageSize, totalCount int
		list []*models.ChildUserModel
		queryFilter = make(map[string]string)
		restData = make(map[string]interface{})
	)

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
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

	if value := this.GetString("userName"); value != "" {
		queryFilter["user_name"] = value
	}
	queryFilter["domain"] = domain

	if list, totalCount, err = models.NewChildUserModel().FindToPager(pageNo, pageSize, queryFilter); err != nil {
		logs.Error("child user list mysql= ", utils.Fields{
			"err": err.Error(),
		})

		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	datalist := views.ChildUserList(list)
	restData["list"] = datalist
	restData["pageNo"] = pageNo
	restData["pageSize"] = pageSize
	restData["total"] = totalCount

	this.JsonResult(httpUtil.SUCCESS, restData)
}

func (this *ChildUserController) Delete() {
	var contentBody map[string]string
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &contentBody)
	if err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "json格式错误")
		return
	}

	id := contentBody["id"]
	err = models.NewChildUserModel().DeleteById(id)
	if err != nil {
		logs.Error("policy delete", utils.Fields{
			"ids": id,
			"err": err.Error(),
		})
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

//子账号角色管理
func (this *ChildUserController) ChildRoleManager() {
	var (
		err error
		userId, roleListStr string
		addList = make([]string, 0)
		delList = make([]string, 0)
	)

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	if addList, delList, userId, roleListStr, err = services.DoChildRoleManager(domain, this.Ctx.Input.RequestBody); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if err =models.NewChildUserModel().UpdateUserRole(domain,userId,roleListStr,addList,delList); err != nil {
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}