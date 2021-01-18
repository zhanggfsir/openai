package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"openai-backend/models/request"
	"openai-backend/services"
	"openai-backend/utils/httpUtil"
)

/*
@Author:
@Time: 2020-12-07 14:02
@Description: api管理
*/

type SysApiController struct {
	BaseController
}

// RouterList 获得路由列表
func (this *SysApiController) GetRoleApiList() {
	var (
		err      error
		domain   string
		//list     []*models.ApiResponse
		restData = make(map[string]interface{})
	)
	// 获得请求参数
	if domain, err = this.getDomain(); err != nil || domain == "" {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	roleId := this.GetString("roleId")
	if roleId == "" {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "roleId为空")
		return
	}

	list, err := services.SerGetRoleApis(domain, roleId)
	if err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	restData["list"] = list
	this.JsonResult(httpUtil.SUCCESS, restData)
}

// UpdateCasbin 更新casbin信息
func (this *SysApiController) UpdateCasbin() {
	var	(
		err error
		domain string
		cmr request.CasbinInReceive
	)
	// 获得请求参数
	if domain, err = this.getDomain(); err != nil || domain == "" {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &cmr); err != nil {
		logs.Error("err=", string(this.Ctx.Input.RequestBody))
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "请求的参数错误")
		return
	}

	logs.Info("sdfsfdsfdsf", cmr)

	if err = services.UpdateCasbin(domain, cmr.RoleId, cmr.CasbinInfos); err != nil {
		this.JsonResultOther(httpUtil.SYSTEM_ERROR, err.Error())
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}