package httpUtil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// JsonResult 响应 json 结果
func JsonResult(ctx *context.Context, errCode string, data map[string]interface{}) {
	jsonData := make(map[string]interface{})

	jsonData["code"] = errCode
	jsonData["msg"] = GetMsg(errCode)

	if len(data) > 0 {
		jsonData["data"] = data
	}

	returnJSON, err := json.Marshal(jsonData)
	if err != nil {
		logs.Error(err)
	}

	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
	io.WriteString(ctx.ResponseWriter, string(returnJSON))
}


// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如："application/json; charset=utf-8"
// content：     请求放回的内容

//var Client = &http.Client{Timeout: 120 * time.Second}

func PostRequest(url, domain string, data interface{}, contentType string) map[string]interface{} {

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	Client := &http.Client{Transport: tr,Timeout: 120 * time.Second}

	//http.Post()
	// 超时时间：5秒

	jsonStr, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	req.Header.Add("Domain", domain)
	req.Header.Set("Content-Type", contentType)
	resp, err := Client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	//logs.Info(string(result))
	resultJson := make(map[string]interface{})
	err = json.Unmarshal(result, &resultJson)
	if err != nil {
		logs.Error(err)
		return nil
	}
	//logs.Info(resultJson)
	return resultJson
}
//
//func PostRequest(url, domain string, data interface{}, contentType string) map[string]interface{} {
//logs.Info(1)
//	//http.Post()
//	// 超时时间：5秒
//
//	jsonStr, _ := json.Marshal(data)
//	//resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
//	if err != nil {
//		logs.Error("--err-->",err)
//		//panic(err)
//	}
//	req.Header.Add("Domain", domain)
//	req.Header.Set("Content-Type", contentType)
//	resp, err := Client.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//
//	result, _ := ioutil.ReadAll(resp.Body)
//	resultJson := make(map[string]interface{})
//	err = json.Unmarshal(result, &resultJson)
//	if err != nil {
//		logs.Error(err)
//		return nil
//	}
//	return resultJson
//}