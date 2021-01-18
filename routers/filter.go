package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"openai-backend/utils"
	"openai-backend/utils/casbinUtil"
	"openai-backend/utils/httpUtil"
	"strings"
	"time"
)

/*
@Author:
@Time: 2020-07-30 10:02
@Description: filter
*/

func filter() {
	beego.InsertFilter("/api/:ver/*", beego.BeforeRouter, checkAuth)

	beego.InsertFilter("/*", beego.AfterExec, finishExec, false)
}

// 权限检查
var checkAuth = func(ctx *context.Context) {
	//设置开始时间
	setStartTime(ctx)

	subPath := ctx.Input.Query(":splat")
	model := strings.Split(subPath, "/")[0]

	if subPath == "token/refresh" || subPath == "thirdParty/token" {
		return
	}
	if model == "login" || model == "register" || model == "captcha" ||
		model == "sms" || model == "gateway" || model == "retrievePassword" {
		return
	}

	// todo auth
	//doAuth(ctx)

	logs.Info("beforeRouter:", utils.Fields{
		"subPath": subPath,
		"model":   model,
		"domain":  ctx.Input.Data()["domain"],
	})

}

var finishExec = func(ctx *context.Context) {
	startTime, _ := utils.A2Int64(ctx.ResponseWriter.Header().Get("StartTime"))
	endTime := time.Now().UnixNano()
	execTimeLen := (endTime - startTime) / 1e6

	requestPath := ctx.Request.URL.Path
	requestMethod := ctx.Request.Method
	requestStatus := ctx.ResponseWriter.Status

	logs.Warning("afterExec", utils.Fields{
		"path":     requestPath,
		"method":   requestMethod,
		"execTime": execTimeLen,
		"status":   requestStatus,
	})
}

func setStartTime(ctx *context.Context) {
	startTime := utils.A2String(time.Now().UnixNano())
	ctx.ResponseWriter.Header().Add("StartTime", startTime)
}

func doAuth(ctx *context.Context) {
	var (
		err                                    error
		jwtUser                                utils.JwtUser
		userId, domain, userType, method, path string
	)

	// jwt token 认证
	if jwtUser, err = utils.VerifySession(ctx.Request.Header); err != nil {
		httpUtil.JsonResult(ctx, err.Error(), nil)
		goto ERR
	}

	userId = jwtUser.Id
	domain = jwtUser.Domain
	userType = jwtUser.UserType
	method = strings.ToLower(ctx.Request.Method)
	path = ctx.Request.URL.Path

	//ctx.Request.Header.Add("userId", userId)
	//ctx.Request.Header.Add("userType", userType)
	ctx.Input.SetData("userId", userId)
	ctx.Input.SetData("userType", userType)
	ctx.Input.SetData("domain", domain)

	if userType == "1"{
		logs.Info("主账户不需限制权限")
		return
	}

	if !casbinUtil.AuthorizationDomain(userId, path, method, domain) {
		httpUtil.JsonResult(ctx, httpUtil.ERROR_NO_AUTHORITY, nil)
		goto ERR
	}

	return

ERR:
	logs.Error("auth err", utils.Fields{
		"userId": userId,
		"err":    err,
		"domain": domain,
		"method": method,
		"path":   path,
	})

	panic(beego.ErrAbort)
}
