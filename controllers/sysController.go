package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"openai-backend/models"
	"openai-backend/services"
	"openai-backend/utils"
	"openai-backend/utils/conf"
	"openai-backend/utils/httpUtil"
)

type SysController struct {
	BaseController
}

// TODO 获得验证码
func (this *SysController) Captcha() {
	var captchaId string
	if captchaId = this.GetString(":captchaId"); captchaId == "" {
		this.JsonResult(httpUtil.PARAMETER_ERROR, nil)

		return
	}

	picContent, fileName := services.GenCaptcha(captchaId)
	CaptchaResp(this, fileName, picContent)
}

// TODO 获得短信验证码
func (this *SysController) SmsCaptcha() {
	smsType := this.GetString("type", "1")
	phone := this.GetString("phone")

	if err := services.SmsValidate(smsType, phone); err != nil {
		logs.Error("sms err=", utils.Fields{
			"err": err.Error(),
			"phone": phone,
			"smsType": smsType,
		})

		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if _, err := services.GenSmsCaptcha(smsType, phone); err != nil {
		this.JsonResult(err.Error(), nil)
	}

	this.JsonResult(httpUtil.SUCCESS, nil)
}

func (this *SysController) Login() {
	var (
		err error
		body map[string]string
		restData map[string]interface{}
		tokenSec, refreshSec string
	)

	tokenSec = this.Ctx.Request.Header.Get("Tokentimeoutsec")
	refreshSec = this.Ctx.Request.Header.Get("Retokentimeoutsec")

	//获得body中的参数
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &body); err != nil {
		logs.Error("body=", string(this.Ctx.Input.RequestBody))
		this.JsonResult(httpUtil.ERROR_JSON_PATTERN, nil)
		return
	}

	if err = services.LoginValidate(tokenSec, refreshSec, body); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if restData, err = services.DoLogin(tokenSec, refreshSec, body); err != nil {
		this.JsonResult(err.Error(), nil)
		return
	}

	this.JsonResult(httpUtil.SUCCESS, restData)
}

// register
func (this *SysController) Register() {
	var (
		err error
		contentBody map[string]string
		restData map[string]interface{}
	)

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &contentBody)
	if err != nil {
		this.JsonResult(httpUtil.ERROR_JSON_PATTERN, nil)

		return
	}

	if err := services.RegisterValidate(contentBody); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())

		return
	}

	if restData, err = services.DoRegisterDomain(contentBody); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())

		return
	}

	this.JsonResult(httpUtil.SUCCESS, restData)
}


// RetrievePassword
func (this *SysController) RetrievePassword() {
	var (
		err error
		contentBody map[string]string
		restData map[string]interface{}
	)

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &contentBody)
	if err != nil {
		this.JsonResult(httpUtil.ERROR_JSON_PATTERN, nil)
		return
	}

	if err := services.RetrievePasswordValidate(contentBody); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())

		return
	}


	this.JsonResult(httpUtil.SUCCESS, restData)
}

// logout
func (this *SysController) Logout() {
	jwtUser, _ := utils.VerifySession(this.Ctx.Request.Header)

	logs.Info("user logout", utils.Fields{
		"domain": jwtUser.Domain,
		"userId": jwtUser.Id,
		"type": jwtUser.Type,
	})

	this.JsonResult(httpUtil.SUCCESS, nil)
}

//刷新token
func (this *SysController) RefreshToken() {
	var (
		err error
		token, refreshToken string
		restData = make(map[string]interface{})
	)
	headerContent := this.Ctx.Request.Header

	logs.Info("///////////////////////////////REFRESH_TOKEN//////////////////////")

	jwtUser, err:= utils.VerifySession(this.Ctx.Request.Header)
	if err!=nil{
		if err.Error() == httpUtil.ERROR_TOKEN_EXPIRED {
			logs.Info("///////////////////////////////ERROR_REFRESH_TOKEN_EXPIRED//////////////////////")
			this.JsonResult(httpUtil.ERROR_REFRESH_TOKEN_EXPIRED, nil)
			return
		}
	}


	tokenSec := headerContent.Get("Tokentimeoutsec")
	refreshSec := headerContent.Get("Retokentimeoutsec")

	if err = services.RefreshTokenValidate(tokenSec, refreshSec); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	if token, refreshToken, err = services.DoRefreshToken(jwtUser, tokenSec, refreshSec); err != nil {
		this.JsonResultOther(httpUtil.PARAMETER_ERROR, err.Error())
		return
	}

	restData["token"] = token
	restData["refreshToken"] = refreshToken
	this.JsonResult(httpUtil.SUCCESS, restData)
}

//第三方获得token
func (this *SysController) ThirdPartyToken() {
	var (
		err error
		childUserInfo *models.ChildUserModel
		restData = make(map[string]interface{})
	)

	appKey := this.GetString("appKey", "")
	appSecret := this.GetString("appSecret", "")

	if childUserInfo, err = services.ThirdPartyTokenVerity(appKey, appSecret); err != nil {
		this.JsonResult(err.Error(), nil)
		return
	}

	token, refreshToken := services.GenTokenAndRefreshToken(childUserInfo.UserId, models.UserTypeChild, childUserInfo.Domain,
		conf.JWT_DEFAULT_EXPIRE_SECONDS, conf.JWT_DEFAULT_LONG_EXPIRE_SECONDS)

	restData["accessToken"] = token
	restData["refreshToken"] = refreshToken
	restData["expireTime"] = conf.JWT_DEFAULT_EXPIRE_SECONDS
	this.JsonResult(httpUtil.SUCCESS, restData)
}