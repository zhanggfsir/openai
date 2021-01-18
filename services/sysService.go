package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/dchest/captcha"
	"io/ioutil"
	"net/http"
	"openai-backend/models"
	"openai-backend/utils"
	"openai-backend/utils/cacheUtil"
	"openai-backend/utils/casbinUtil"
	"openai-backend/utils/conf"
	"openai-backend/utils/httpUtil"
	"openai-backend/utils/redis"
	"regexp"
	"strings"
	"time"
)

var GlobalStore captcha.Store

type respLogin struct {
	userId   string
	userName string
	password string
	domain   string
}

func InitCaptcha() {
	collectNum := captcha.CollectNum
	expiration := 3 * time.Minute
	GlobalStore = captcha.NewMemoryStore(collectNum, expiration)
	captcha.SetCustomStore(GlobalStore)
}

func GenCaptcha(captchaId string) (*bytes.Buffer, string) {
	var content bytes.Buffer
	// 保存验证码
	GlobalStore.Set(captchaId, captcha.RandomDigits(captcha.DefaultLen))

	captchaWidth := captcha.StdWidth
	captchaHeight := captcha.StdHeight

	captcha.WriteImage(&content, captchaId, captchaWidth, captchaHeight)
	saveRedisCaptcha(captchaId)

	return &content, captchaId + ".png"
}

func SmsValidate(smsType, phone string) error {
	valid := validation.Validation{}

	valid.Phone(phone, "phone").Message("手机号码格式错误")
	valid.Length(smsType, 1, "smsType").Message("类型错误")

	if valid.HasErrors() {
		var errStr string
		for _, err := range valid.Errors {
			errStr = errStr + err.Message
		}
		return errors.New(errStr)
	}

	//根据手机号码获得用户信息
	if smsType == "2" {
		if _, err := models.NewUserModel().FetchByPhone(phone); err != nil {
			return errors.New("手机号码不存在")
		}
	}

	if smsType == "1" {
		if  models.NewUserModel().QueryTable().Filter("phone",phone).Exist(){
			return errors.New("手机号码已注册")
		}
	}
	return nil
}

// 生成短信验证码
func GenSmsCaptcha(smsType, phone string) (string, error) {
	var (
		err             error
		captcha, smsKey string
	)
	smsKey = genSmsKey(smsType, phone)

	captcha = utils.GenerateRandom("num", 6)
	if err = smsPost(captcha, phone); err != nil {
		smsLog(phone, "0", captcha, smsType)
		return "", errors.New(httpUtil.SYSTEM_ERROR)
	}

	smsLog(phone, "1", captcha, smsType)

	if err = saveSmsRedis(smsKey, captcha); err != nil {
		return "", errors.New(httpUtil.SYSTEM_ERROR)
	}

	return smsKey, nil
}

func LoginValidate(tokenSec, refreshSec string, body map[string]string) error {
	userType := body["type"]
	hostUser := body["hostUser"]
	userName := body["userName"]
	password := body["password"]
	captchaId := body["captchaId"]
	captcha := body["captcha"]
	valid := validation.Validation{}

	valid.Length(userType, 1, "userType").Message("请输入用户类型,")
	valid.Required(userName, "userName").Message("请输入用户名,")
	valid.Required(password, "password").Message("请输入密码,")
	valid.Required(captchaId, "captchaId").Message("请输入验证码id,")
	valid.Required(captcha, "captcha").Message("请输入验证码,")
	//valid.Required(tokenSec, "tokenSec").Message("请输入token过期时间,")
	//valid.Match(tokenSec, regexp.MustCompile(`^[1-9][0-9]{1,8}`), "tokenSec").
	//	Message("token过期时间证格式错误,")
	//valid.Required(refreshSec, "refreshSec").Message("请输入refresh token过期时间,")
	//valid.Match(refreshSec, regexp.MustCompile(`^[1-9][0-9]{1,8}`), "refreshSec").
	//	Message("refresh token过期时间证格式错误,")

	if valid.HasErrors() {
		var errStr string
		for _, err := range valid.Errors {
			errStr = errStr + err.Message
		}
		return errors.New(errStr)
	}

	if userType == "2" && hostUser == "" {
		return errors.New("请输入主账号")
	}

	//if !VerifyCaptcha(captchaId, captcha) {
	//	return errors.New("验证码错误")
	//}

	return nil
}

func DoLogin(tokenSec, refreshSec string, body map[string]string) (map[string]interface{}, error) {
	var (
		err      error
		userInfo *respLogin
	)
	userType := body["type"]
	hostUser := body["hostUser"]
	userName := body["userName"]
	password := body["password"]

	if userType == "1" {
		userInfo, err = doHostLogin(userName)
	} else {
		userInfo, err = doChildLogin(hostUser, userName)
	}

	if err != nil {
		logs.Error(err)
		return nil, errors.New(httpUtil.ERROR_USER_NOT_EXIST)
	}

	if err = verifyUserPassword(password, userInfo.password); err != nil {
		logs.Error(err)
		return nil, err
	}

	userId := userInfo.userId
	domain := userInfo.domain
	tokenSecInt64, _ := utils.A2Int64(tokenSec)
	refreshSecInt64, _ := utils.A2Int64(refreshSec)

	if tokenSecInt64 == 0 {
		tokenSecInt64 = utils.JWT_DEFAULT_EXPIRE_SECONDS
	}
	if refreshSecInt64 == 0 {
		refreshSecInt64 = utils.JWT_DEFAULT_LONG_EXPIRE_SECONDS
	}

	token, refreshToken := GenTokenAndRefreshToken(userId, userType, domain, tokenSecInt64, refreshSecInt64)

	restData := map[string]interface{}{
		"userId":       userId,
		"userName":     userName,
		"domain":       domain,
		"userType":     userType,
		"token":        token,
		"refreshToken": refreshToken,
	}

	return restData, nil
}

func doHostLogin(userName string) (*respLogin, error) {
	var (
		err        error
		userInfo   *models.UserModel
		domainInfo *models.DomainModel
	)
	if utils.VerifyMobile(userName) {
		phone := userName
		if userInfo, err = models.NewUserModel().FetchByPhone(phone); err != nil {
			goto ERR
		}
	} else {
		if userInfo, err = models.NewUserModel().FetchByUserName(userName); err != nil {
			goto ERR
		}
	}

	if domainInfo, err = models.NewDomainModel().FetchByUserId(userInfo.UserId); err != nil {
		goto ERR
	}

	return &respLogin{
		userId:   userInfo.UserId,
		domain:   domainInfo.Domain,
		userName: userInfo.UserName,
		password: userInfo.Password,
	}, err

ERR:
	logs.Error("doHostLogin", utils.Fields{
		"err":      err.Error(),
		"userName": userName,
	})

	return nil, err
}

func doChildLogin(hostUser, userName string) (*respLogin, error) {
	var (
		err        error
		userInfo   *models.UserModel
		childInfo  *models.ChildUserModel
		domainInfo *models.DomainModel
	)
	if utils.VerifyMobile(hostUser) {
		phone := hostUser
		if userInfo, err = models.NewUserModel().FetchByPhone(phone); err != nil {
			logs.Error(err)
			goto ERR
		}
	} else {
		if userInfo, err = models.NewUserModel().FetchByUserName(hostUser); err != nil {
			logs.Error(err)
			goto ERR
		}
	}

	if domainInfo, err = models.NewDomainModel().FetchByUserId(userInfo.UserId); err != nil {
		logs.Error(err)
		goto ERR
	}

	if childInfo, err = models.NewChildUserModel().FetchByUserName(domainInfo.Domain, userName); err != nil {
		logs.Error("domain", domainInfo.Domain)
		goto ERR
	}

	return &respLogin{
		userId:   childInfo.UserId,
		userName: childInfo.UserName,
		password: childInfo.Password,
		domain:   childInfo.Domain,
	}, err

ERR:
	logs.Error("doChildLogin", utils.Fields{
		"err":      err.Error(),
		"hostUser": hostUser,
		"userName": userName,
	})

	return nil, err
}

func RegisterValidate(body map[string]string) error {
	var err error
	phone := body["phone"]
	userName := body["userName"]
	password := body["password"]
	captcha := body["captcha"]
	valid := validation.Validation{}

	valid.Phone(phone, "phone").Message("手机号码格式错误")
	valid.Required(userName, "userName").Message("请输入用户名")
	valid.Required(password, "password").Message("请输入密码,")
	valid.Required(captcha, "captcha").Message("请输入验证码,")

	if valid.HasErrors() {
		var errStr string
		for _, err := range valid.Errors {
			errStr = errStr + err.Message
		}
		err = errors.New(errStr)
		goto ERR
	}

	// 判断手机号码唯一
	if _, err = models.NewUserModel().FetchByPhone(phone); err == nil {
		err = errors.New("手机号已经注册")
		goto ERR
	}

	// 判断公司名称唯一
	if _, err = models.NewUserModel().FetchByUserName(userName); err == nil {
		err = errors.New("用户名已经存在")
		goto ERR
	}

	if _, err = base64.StdEncoding.DecodeString(password); err != nil {
		err = errors.New("密码加密错误")
		goto ERR
	}

	if !VerifyCaptcha(genSmsKey("1", phone), captcha) {
		err = errors.New("验证码错误")
		goto ERR
	}

	return nil

ERR:
	logs.Error("register error", utils.Fields{
		"password": password,
		"phone":    phone,
		"userName": userName,
		"captcha":  captcha,
		"err":      err.Error(),
	})

	return err
}

func RetrievePasswordValidate(body map[string]string) error {
	var err error
	var decodePd []byte
	var passwordNew string
	phone := body["phone"]
	password := body["password"]
	captcha_ := body["captcha"]
	valid := validation.Validation{}
	data := make(map[string]interface{})

	valid.Phone(phone, "phone").Message("手机号码格式错误")
	valid.Required(password, "password").Message("请输入密码,")
	valid.Required(captcha_, "captcha").Message("请输入验证码,")

	if valid.HasErrors() {
		var errStr string
		for _, err := range valid.Errors {
			errStr = errStr + err.Message
		}
		err = errors.New(errStr)
		goto ERR
	}

	// 判断手机号码
	if !models.NewUserModel().QueryTable().Filter("phone", phone).Exist() {
		err = errors.New("手机号不存在")
		goto ERR
	}

	if _, err = base64.StdEncoding.DecodeString(password); err != nil {
		err = errors.New("密码加密错误")
		goto ERR
	}

	if !VerifyCaptcha(genSmsKey("2", phone), captcha_) {
		err = errors.New("验证码错误")
		goto ERR
	}

	//updata password
	decodePd, _ = base64.StdEncoding.DecodeString(password)
	passwordNew = utils.GenPassword(string(decodePd))
	data["password"] = passwordNew

	if err := models.NewUserModel().UpdateByPhone(phone, data); err != nil {
		err = errors.New("更新密码错误")
		goto ERR
	}

	return nil

ERR:
	logs.Error("RetrievePassword error", utils.Fields{
		"password": password,
		"phone":    phone,
		"captcha":  captcha_,
		"err":      err.Error(),
	})

	return err
}

func DoRegisterDomain(body map[string]string) (map[string]interface{}, error) {
	phone := body["phone"]
	userName := body["userName"]
	password := body["password"]
	userId := utils.GenerateUuid()
	roleId := casbinUtil.RoleAdmin
	domain, _ := generateDomain()

	domainStruct := models.DomainModel{
		UserId: userId,
		Domain: domain,
		Phone:  phone,
	}

	if err := domainStruct.DomainRegister(password, roleId, userName); err != nil {
		return nil, errors.New(httpUtil.SYSTEM_ERROR)
	}

	restData := map[string]interface{}{
		"domain": domain,
	}
	return restData, nil
}

func RefreshTokenValidate(tokenSec, refreshSec string) error {
	if tokenSec != "" {
		valid := validation.Validation{}
		valid.Match(tokenSec, regexp.MustCompile(`^[1-9][0-9]{1,8}`), "tokenSec").
			Message("token过期时间证格式错误,")
		valid.Match(refreshSec, regexp.MustCompile(`^[1-9][0-9]{1,8}`), "refreshSec").
			Message("refresh token过期时间证格式错误,")

		if valid.HasErrors() {
			var errStr string
			for _, err := range valid.Errors {
				errStr = errStr + err.Message
			}
			return errors.New(errStr)
		}
	}

	return nil
}

func DoRefreshToken(jwtUser utils.JwtUser, tokenSec, refreshSec string) (string, string, error) {
	var (
		errToken, errRe                error
		tokenSecInt64, refreshSecInt64 int64
	)
	if tokenSec != "" {
		tokenSecInt64, errToken = utils.A2Int64(tokenSec)
		refreshSecInt64, errRe = utils.A2Int64(refreshSec)
		if errToken != nil || errRe != nil {
			return "", "", errors.New("token 过期时间格式错误")
		}
	} else {
		tokenSecInt64 = conf.JWT_DEFAULT_EXPIRE_SECONDS
		refreshSecInt64 = conf.JWT_DEFAULT_LONG_EXPIRE_SECONDS
	}

	token := utils.GenerateToken(jwtUser, tokenSecInt64)
	refreshToken := utils.GenerateToken(jwtUser, refreshSecInt64)

	return token, refreshToken, nil
}

func VerifyCaptcha(captchaId, captchaValue string) bool {
	var redisCaptchaValue string
	if captchaId == "" || captchaValue == "" {
		return false
	}
	ctx := context.Background()
	redisCaptchaValue, err := redis.C().Get(ctx, captchaId).Result()
	if err == redis.Nil {
		logs.Info("key does not exist")
	} else if err != nil {
		logs.Info(err)
		return false
	} else {
		logs.Info(captchaId, redisCaptchaValue)
	}

	//if err := cacheUtil.Get(captchaId, &redisCaptchaValue); err != nil {
	//	logs.Error("redis captcha:", utils.Fields{
	//		"requestCaptcha": captchaValue,
	//		"err":            err.Error(),
	//	})
	//
	//	return false
	//}

	if captchaValue != redisCaptchaValue {
		logs.Error("verify captcha:", utils.Fields{
			"requestCaptcha": captchaValue,
			"saveCaptcha":    redisCaptchaValue,
		})

		return false
	}

	//cacheUtil.Delete(captchaId)
	redis.C().Del(ctx, captchaId)
	return true
}

func ThirdPartyTokenVerity(appKey, appSecret string) (*models.ChildUserModel, error) {
	var (
		err           error
		childUserInfo *models.ChildUserModel
	)
	if appKey == "" || appSecret == "" {
		return nil, errors.New(httpUtil.PARAMETER_ERROR)
	}

	if childUserInfo, err = models.NewChildUserModel().FetchByAppKey(appKey); err != nil {
		logs.Error("err=", err)
		return nil, errors.New(httpUtil.SYSTEM_ERROR)
	}

	if appSecret != childUserInfo.AppSecret {
		logs.Error("third party get token err", utils.Fields{
			"appKey":     appKey,
			"appSecret":  appSecret,
			"saveSecret": childUserInfo.AppSecret,
		})
		return nil, errors.New(httpUtil.PARAMETER_ERROR)
	}

	return childUserInfo, nil
}

func saveRedisCaptcha(captchaId string) {
	isOpenCache := beego.AppConfig.DefaultBool("cache", false)
	if !isOpenCache {
		return
	}
	strCaptchaValue := getCaptchaValue(captchaId)
	logs.Info("strCaptchaValue= ", strCaptchaValue)

	ctx := context.Background()
	if err := redis.C().Set(ctx, captchaId, strCaptchaValue, 5*time.Minute).Err(); err != nil {
		logs.Error("", err)
	}
	//if err := cacheUtil.Put(captchaId, strCaptchaValue, time.Minute*5); err != nil {
	//	logs.Error("", err)
	//}
}

func getCaptchaValue(captchaId string) string {
	isOpenCache := beego.AppConfig.DefaultBool("cache", false)
	if !isOpenCache {
		cacheUtil.Init(&cacheUtil.NullCache{})
		return ""
	}

	//captchaValue, err := redis.C().Get(ctx, captchaId).Result()
	//if err == redis.Nil {
	//	logs.Info("key does not exist")
	//} else if err != nil {
	//	logs.Info(err)
	//} else {
	//	logs.Info(captchaId, captchaValue)
	//}

	captchaValue := GlobalStore.Get(captchaId, true)
	ns := make([]byte, len(captchaValue))
	for i := range ns {
		d := captchaValue[i]
		switch {
		case 0 <= d && d <= 9:
			ns[i] = d + '0'
		default:
			return ""
		}
	}
	return string(ns)
}

func verifyUserPassword(requestPassword, savePassword string) error {
	var (
		err            error
		scriptPassword string
	)

	if scriptPassword, err = decodeAndScriptPassword(requestPassword); err != nil {
		return errors.New(httpUtil.ERROR_USER_LOGIN_ERROR)
	}

	logs.Info("scriptPassword:", scriptPassword, "savePassword:", savePassword)
	if scriptPassword != savePassword {
		logs.Error("base64 decode", utils.Fields{
			"err":            "密码不对",
			"scriptPassword": scriptPassword,
			"savePassword":   savePassword,
		})

		return errors.New(httpUtil.ERROR_USER_LOGIN_ERROR)
	}

	return nil
}

func decodeAndScriptPassword(password string) (string, error) {
	passwordDecode, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		logs.Error("base64 decode", utils.Fields{
			"password": password,
			"err":      err.Error(),
		})

		return "", err
	}

	scriptPassword := utils.GenPassword(string(passwordDecode))

	return scriptPassword, nil
}

func generateDomain() (string, error) {
	var domain string
	flag := true
	for i := 0; i <= 100; i++ {
		domain = utils.GenerateRandom("num", 8)
		_, err := models.NewDomainModel().FetchByDomain(domain)
		if err != nil {
			flag = false
			break
		}
	}
	if !flag {
		return domain, nil
	}
	return "", errors.New("domain repeat")
}

/*
@Description: 短信相关的操作
*/
func genSmsKey(smsType, phone string) string {
	return smsType + "-" + phone
}

// TODO 调用发送短信接口JSON
func smsPost(captcha, phone string) error {
	dataList := []string{phone}
	smsUrl := smsUrl()
	jsonData := make(map[string]interface{})
	jsonData["wordId"] = getSmsWordId()
	jsonData["variableOne"] = captcha
	jsonData["dataList"] = dataList
	jsonBody, err := json.Marshal(jsonData)
	if err != nil {
		logs.Error("sms api json=", utils.Fields{
			"err":      err,
			"jsonData": jsonData,
		})

		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", smsUrl, strings.NewReader(string(jsonBody)))
	if err != nil {
		logs.Error("sms api NewRequest=", utils.Fields{
			"err":      err,
			"jsonData": jsonData,
		})

		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Sig", getSmsSign())

	resp, err := client.Do(req)
	if err != nil {
		logs.Error("sms api NewRequest=", utils.Fields{
			"err":      err,
			"jsonData": jsonData,
		})

		return err
	}
	defer resp.Body.Close()

	var contentBody map[string]interface{}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(resBody, &contentBody); err != nil {
		logs.Error("sms api response=", utils.Fields{
			"err":      err,
			"jsonData": jsonData,
			"resBody":  string(resBody),
		})

		return err
	}

	logs.Warning("sms api response=", utils.Fields{
		"jsonData": jsonData,
		"resBody":  string(resBody),
	})

	if code, ok := contentBody["code"]; ok && code.(string) == "02000" {
		return nil
	} else {
		return errors.New(httpUtil.SYSTEM_ERROR)
	}

}

func smsUrl() string {
	companyId := getSmsCompanyId()
	url := "/sdyxinterface/20190426/msg/sendMsgTel/" + companyId
	smsHost := beego.AppConfig.DefaultString("sms_host", "http://120.52.23.243:10080")
	return smsHost + url
}

func getSmsCompanyId() string {
	return beego.AppConfig.DefaultString("sms_company", "D1NLCN")
}

func getSmsWordId() string {
	return beego.AppConfig.DefaultString("sms_id", "D1NLCN_001")
}

func getSmsSign() string {
	return beego.AppConfig.DefaultString("sms_token", "")
}

func getSmsTimeout() int64 {
	return beego.AppConfig.DefaultInt64("sms_timeout", 60)
}

func smsLog(phone, status, captcha, smsType string) {
	smsInfo := &models.SmsModel{
		Phone:   phone,
		Status:  status,
		Captcha: captcha,
		Type:    smsType,
	}
	if err := smsInfo.Insert(); err != nil {
		logs.Error("insert sms mysql=", utils.Fields{
			"smsInfo": smsInfo,
			"err":     err,
		})
	}
}

func saveSmsRedis(key, captcha string) error {
	isOpenCache := beego.AppConfig.DefaultBool("cache", false)
	if !isOpenCache {
		return nil
	}

	timeout := getSmsTimeout()

	ctx := context.Background()
	if err := redis.C().Set(ctx, key, captcha, time.Second*time.Duration(timeout)).Err(); err != nil {
		logs.Error("save redis err", utils.Fields{
			"err":     err,
			"key":     key,
			"captcha": captcha,
		})
		return err
	}

	//if err := cacheUtil.Put(key, captcha, time.Second*time.Duration(timeout)); err != nil {
	//	logs.Error("save redis err", utils.Fields{
	//		"err":     err,
	//		"key":     key,
	//		"captcha": captcha,
	//	})
	//
	//	return err
	//}

	logs.Info("///////////save Redis//////////")
	return nil
}

func GenTokenAndRefreshToken(userId, userType, domain string, tokenSec, refreshSec int64) (string, string) {
	jwtUser := utils.JwtUser{
		Id:       userId,
		Type:     userType,
		Domain:   domain,
		UserType: userType,
	}

	tokenSecInt64, _ := utils.A2Int64(tokenSec)
	refreshSecInt64, _ := utils.A2Int64(refreshSec)
	token := utils.GenerateToken(jwtUser, tokenSecInt64)
	refreshToken := utils.GenerateToken(jwtUser, refreshSecInt64)

	return token, refreshToken
}
