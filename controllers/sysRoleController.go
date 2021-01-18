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
@Time: 2020-08-06 20:02
@Description: 角色管理
*/

type RoleController struct {
	BaseController
}

func (this *RoleController) Post() {
	var (
		err error
		body map[string]string
		roleInfo *models.RoleModel
	)

	domain, err := this.getDomain()
	if err != nil || domain == "" {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	//获得body中的参数
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &body); err != nil {
		this.JsonResult(httpUtil.ERROR_JSON_PATTERN, nil)
		return
	}

	if roleInfo, err = services.RoleAddValidate(domain, body); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if err = roleInfo.Insert(); err != nil {
		logs.Error("role insert err", utils.Fields{
			"body": body,
			"err": err.Error(),
		})
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

func (this *RoleController) Put() {
	var (
		err error
		domain, id string
		body map[string]string
		restData map[string]interface{}
	)

	if domain, err = this.getDomain(); err != nil {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}
	//获得body中的参数
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &body); err != nil {
		this.JsonResult(httpUtil.ERROR_JSON_PATTERN, nil)
		return
	}

	if restData, id, err = services.RoleUpdateValidate(body); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if err = models.NewRoleModel().SaveFilter(domain, id, restData); err != nil {
		logs.Error("role update err", utils.Fields{
			"body": body,
			"err": err.Error(),
		})
		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

func (this *RoleController) GetList() {
	var (
		err error
		domain string
		pageNo, pageSize, totalCount int
		list []*models.RoleModel
		queryFilter = make(map[string]string)
		restData = make(map[string]interface{})
	)

	if domain, err = this.getDomain(); err != nil {
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
	queryFilter["domain"] = domain
	if list, totalCount, err = models.NewRoleModel().FindToPager(pageNo, pageSize, queryFilter); err != nil {
		logs.Error("policy list mysql= ", utils.Fields{
			"err": err.Error(),
		})

		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	datalist := views.RoleList(list)
	restData["list"] = datalist
	restData["pageNo"] = pageNo
	restData["pageSize"] = pageSize
	restData["total"] = totalCount

	this.JsonResult(httpUtil.SUCCESS, restData)
}

//TODO 删除角色的时候判断是是否有人
func (this *RoleController) Delete() {
	var contentBody map[string]string
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &contentBody)
	if err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "json格式错误")

		return
	}

	id := contentBody["id"]
	if err = models.NewRoleModel().DeleteById(id); err != nil {
		logs.Error("role delete", utils.Fields{
			"ids": id,
			"err": err.Error(),
		})

		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

// role/authority  list
func (this *RoleController) RoleAuthList() {
	var (
		err error
		jwtUser utils.JwtUser
		roles []string
		permissions map[string][]string
		restData = make(map[string]interface{})
	)

	if jwtUser, err = utils.VerifySession(this.Ctx.Request.Header); err != nil {
		this.JsonResult(err.Error(), nil)
		return
	}

	userId := jwtUser.Id
	userType := jwtUser.UserType
	domain := jwtUser.Domain

	logs.Info("get roles", utils.Fields{
		"userId": userId,
		"userType": userType,
		"domain": domain,
	})

	if permissions, roles, err = services.DoRoleAuthList(domain, userId, userType); err != nil {
		logs.Error("RoleAuthList err", utils.Fields{
			"userId": userId,
			"userType": userType,
			"domain": domain,
			"err": err.Error(),
		})

		this.JsonResult(httpUtil.SYSTEM_ERROR, nil)
		return
	}

	restData["roles"] = roles
	restData["permissions"] = permissions
 	this.JsonResult(httpUtil.SUCCESS, restData)
}

//获得role auth list
func (this *RoleController) GetApiRoleAuthList()  {
	var (
		err error
		RequestRoleId string
		jwtUser utils.JwtUser
	)

	if jwtUser, err = utils.VerifySession(this.Ctx.Request.Header); err != nil {
		this.JsonResult(err.Error(), nil)
		return
	}

	if RequestRoleId = this.GetString("roleId"); RequestRoleId == "" {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "请输入角色Id")
		return
	}

	datalist := services.GetRoleApiList(RequestRoleId, jwtUser)

	var restData = make(map[string]interface{})
	restData["list"] = datalist

	this.JsonResult(httpUtil.SUCCESS, restData)
}

func (this *RoleController) UpdateApiRoleAuthList() {
	var (
		err error
		domain string
	)

	if domain, err = this.getDomain(); err != nil {
		logs.Error("err", "域名错误")
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, "域名错误")
		return
	}

	//获得body中的参数
	if err = services.DoUpdateApiRoleAuthList(domain, this.Ctx.Input.RequestBody); err != nil {
		this.JsonResult(err.Error(), nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

