package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"openai-backend/utils/httpUtil"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type BaseController struct {
	beego.Controller
	Option                map[string]string
	EnableAnonymous       bool
	EnableDocumentHistory bool
}

// datasetImportTaskController
var (
	uploadFileParameter         = "file"
	defaultMaxPictureSize int64 = 100
)

var (
	defaultRootPath = "static/" //默认的文件保存路径
	noonCasbin      = "noon"    //不需要使用casbin认证
	specialChars    = []string{"%", "_", "*", ".", "#", "@", "!", "^", "&"}
)

var (
	accessAuthSource = "permission" //权限下发资源管里
)

//获得参数中的域名
func (c *BaseController)getDomain() (string, error) {
	if value, ok := c.Ctx.Request.Header["Domain"]; ok && value[0] != "" {
		return value[0], nil
	} else {//
		inputData := c.Ctx.Input.Data()
		if inputData["domain"] != nil{
			domain := inputData["domain"].(string)
			return domain,nil
		}else {
			return "", errors.New("域名不存在")
		}
	}
}

func (c *BaseController) JsonResultObj(errCode int, o interface{}, msg string) {
	jsonData := make(map[string]interface{})

	jsonData["code"] = errCode
	jsonData["msg"] = msg

	if o != nil {
		jsonData["data"] = o
	}

	returnJSON, err := json.Marshal(jsonData)
	if err != nil {
		logs.Error(err)
	}

	//logs.Error("returnJSON", string(returnJSON))

	fmt.Println("返回数据 *************** ", string(returnJSON))

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
	_, _ = io.WriteString(c.Ctx.ResponseWriter, string(returnJSON))

	//c.StopRun()
}

// JsonResult 响应 json 结果
func (c *BaseController) JsonResult(errCode string, data map[string]interface{}) {
	jsonData := make(map[string]interface{})

	jsonData["code"] = errCode
	jsonData["msg"] = httpUtil.GetMsg(errCode)

	if len(data) > 0 {
		jsonData["data"] = data
	}

	returnJSON, err := json.Marshal(jsonData)
	if err != nil {
		logs.Error(err)
	}

	//logs.Error("returnJSON", string(returnJSON))

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
	_, _ = io.WriteString(c.Ctx.ResponseWriter, string(returnJSON))

	//c.StopRun()
}

// JsonResult 响应 json 结果
func (c *BaseController) JsonResultList(errCode string, data []map[string]interface{}) {
	jsonData := make(map[string]interface{}, 3)

	jsonData["code"] = errCode
	jsonData["msg"] = httpUtil.GetMsg(errCode)
	jsonData["data"] = data

	returnJSON, err := json.Marshal(jsonData)
	if err != nil {
		logs.Error(err)
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
	_, _ = io.WriteString(c.Ctx.ResponseWriter, string(returnJSON))

	//c.StopRun()
}

// JsonResult 响应 json 结果
func (c *BaseController) JsonResultOther(errCode string, msg string) {
	jsonData := make(map[string]interface{})

	jsonData["code"] = errCode
	jsonData["msg"] = msg

	returnJSON, err := json.Marshal(jsonData)
	if err != nil {
		logs.Error(err)
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
	_, _ = io.WriteString(c.Ctx.ResponseWriter, string(returnJSON))

	//c.StopRun()
}

// 验证码返回
func CaptchaResp(c *SysController, fileName string, content *bytes.Buffer) {
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Ctx.ResponseWriter.Header().Set("Pragma", "no-cache")
	c.Ctx.ResponseWriter.Header().Set("Expires", "0")
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")

	http.ServeContent(c.Ctx.ResponseWriter, c.Ctx.Request, fileName, time.Time{}, bytes.NewReader(content.Bytes()))
}

func (this *BaseController) checkErr(errcode, msg string, err error) bool {
	if err != nil {
		this.JsonResultOther(errcode, msg+" "+err.Error())
		return true
	}
	return false
}
