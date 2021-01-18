package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"openai-backend/services"
	"strconv"
)

type AuthController struct {
	BaseController
}


/*
  令牌 token 是由头部(Header)、载荷(Payload)和签名(Signature)三部分通过.连接起来的字符串。
  Header.Payload.Signature

*/



func (this *AuthController) Get() {
	logs.Info("---------开始----------------")
	key:=this.GetString("key")
	path:=this.GetString("path")

	authCacheResp,err:=services.CacheHandler(key,path)
	if err!=nil{
		logs.Error(err)
		this.JsonResultObj(1101, nil, err.Error())
		return
	}
	logs.Info("---------end----------------")
	//logs.Info("app->",authCacheResp.App,"   api->",authCacheResp.Api)
	this.JsonResultObj(0, authCacheResp, "")
	return
}








func (this *AuthController) Post() {
	var req struct {
		Method       string       	   `json:"method"`
		Headers 	 map[string]string `json:"headers"`
		Path   		 string            `json:"path"`
		RawQurery    string            `json:"rawQurery"`
	}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); nil != err {
		this.JsonResultObj(2, nil, "parse body 2 json error")
		return
	}
	// 鉴权
	authResp ,_, err := services.VerifyToken(req.Headers, req.Path,req.RawQurery,req.Method)
	if err != nil {
		logs.Error(err)
		this.JsonResultObj(2, nil, fmt.Sprintf("验证失败: %s", err))
		return
	}
	logs.Info("----鉴权完成-----","app->",authResp.UbdApp,"  api->",authResp.UbdApi," key->",authResp.UbdKey)
	// 限流
	err = services.QpsAndQuotasLimitCheck(authResp)
	if err !=nil {
		logs.Error(err)
		this.JsonResultObj(1011, nil, err.Error())
		return
	}
	logs.Info("----限流完成-----")
	logs.Info("----限流完成-----","x-ubd-api",strconv.FormatUint(authResp.UbdApi,10))
	logs.Info("----限流完成-----","x-ubd-key",authResp.UbdKey)
	logs.Info("----限流完成-----","x-ubd-app",strconv.FormatUint(authResp.UbdApp,10))


	this.Ctx.ResponseWriter.Header().Set("x-ubd-api",strconv.FormatUint(authResp.UbdApi,10))
	this.Ctx.ResponseWriter.Header().Set("x-ubd-key",authResp.UbdKey)
	this.Ctx.ResponseWriter.Header().Set("x-ubd-app",strconv.FormatUint(authResp.UbdApp,10))
	//this.Ctx.ResponseWriter
	this.JsonResultObj(0, nil, "")
	return
}
