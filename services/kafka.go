package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/url"
	"openai-backend/models"
	"openai-backend/utils/commonUtil"
	"strconv"
	"strings"
	//"sync"
	"time"

	"github.com/Shopify/sarama"
)

// StartKafka 启动kafka
func StartKafka() {
	GetFaceRecognitionMap()
	go func() {
		config := sarama.NewConfig()
		config.ClientID = commonUtil.GenUuid()
		config.Version = sarama.MaxVersion
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
		config.Consumer.Offsets.Initial = sarama.OffsetNewest  // sarama.OffsetOldest  // OffsetNewest
		config.Consumer.MaxProcessingTime = 2 * time.Second

		for {
			startKafka(config)
			time.Sleep(time.Second)
		}

	}()

}

func startKafka(config *sarama.Config) {

	defer func() {
		if rcv := recover(); nil != rcv {
			logs.Error("start kafka reconver: %s", rcv)
		}
	}()
	//var addrs []string
	addr := beego.AppConfig.DefaultString("kafka_brokers", "")
	addrs := strings.Split(addr, ",")

	client, err := sarama.NewClient(addrs, config)
	if nil != err {
		logs.Error("new kafka client error: %s", err)
		return
	}

	topic := beego.AppConfig.DefaultString("kafka_topic", "")
	//partitions, err := client.Partitions(topic)
	if nil != err {
		logs.Error("get partitons error: %s", err)
		return
	}

	consumerGroup := beego.AppConfig.DefaultString("kafka_consumer_group", "")

	consumer, err := sarama.NewConsumerGroupFromClient(consumerGroup, client)
	if err != nil {
		logs.Error(fmt.Errorf("Error creating consumer group client: %v", err))
		return
	}
	defer client.Close()

	errCh := consumer.Errors()
	go func() {
		for {
			select {
			case e := <-errCh:
				fmt.Println("err", e)
			}
		}
	}()

	if err = consumer.Consume(context.Background(), []string{topic}, &KafkaConsumeHandler{}); nil != err {
		fmt.Errorf("consumer异常退出: %s", err)
	}


}

// KafkaConsumeHandler handler
type KafkaConsumeHandler struct{}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (*KafkaConsumeHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (*KafkaConsumeHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

var countMessage int64

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (*KafkaConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine


	for msg := range claim.Messages() {
		//logs.Info("---- from kafka--->", msg.Partition, "  ") //  msg.Value,
		//countMessage++
		//if countMessage % 10000==0{
		//	logs.Info("------countMessage------->",countMessage)
		//}

		HandlerMsg(msg.Value)
		session.MarkMessage(msg, "")
	}
	return nil
}

var smartDoor = "/smartdoor"
var smartCommunity = "/smartcommunity"
var smartCampus = "/smartcampus"

var maskFace = "/maskface"
var search = "/search"
var maskFaceSearch = "/maskface/search"

var face = "/face"
var faceSearch = "/face/search"

var deviceReport = "/device/report"
var other = "/other"

type ReqHeader struct {
	Authentication string `json:"authentication"`
	Host           string `json:"host"`
}

type Request struct {
	Url        string    `json:"url"`
	ReqHeaders ReqHeader `json:"headers"`
	Method     string    `json:"method"`
}

type RespHeader struct {
	XUbdApi string `json:"x-ubd-api"` // ubd_api  x-ubd-api  x_ubd_api
	UbdKey  string `json:"x-ubd-key"`
	UbdApp  string `json:"x-ubd-app"`
	UbdAuth string `json:"x-ubd-auth"` //0 未鉴权 string

}
type Response struct {
	RespHeaders RespHeader `json:"headers"`
	Status      uint64     `json:"status"`
}

type ReqAndRes struct {
	MyResponse Response `json:"response"`
	MyRequest  Request  `json:"request"`
	Ts         uint64   `json:"ts"`
}

var count int64
var  categoryId=1  // 人脸识别不鉴权
var  faceRecognitionMap = make(map[string]uint64)

func GetFaceRecognitionMap() map[string]uint64 {

	ability:=&models.AbilityModel{}
	abilityList:=make([]*models.AbilityModel,10)


	_,err:=ability.QueryTable().Filter("category_id",categoryId).All(&abilityList)
	if err!=nil{
		logs.Error(err)
	}
	abilityIdList:=make([]uint64,0)
	for _, myAbility := range abilityList {
		abilityIdList=append(abilityIdList,myAbility.Id)
	}

	//for _, id := range abilityIdList {
	//	logs.Info("--id---",id)
	//}
	//logs.Info("-------------------------------")
	apiInfo:=&models.ApiModel{}
	apiList:=make([]*models.ApiModel,10)
	_,err=apiInfo.QueryTable().Filter("ability_id__in",abilityIdList).All(&apiList)
	if err!=nil{
		logs.Error(err)
	}

	for _, api := range apiList {
		faceRecognitionMap[strings.TrimSpace(api.Url)]=api.Id
	}
	//logs.Info(  " ")
	//for i, i2 := range faceRecognitionMap {
	//	logs.Info(i," -",i2)
	//}
	//logs.Info(faceRecognitionMap["/openapi/cv/face/v1/compare"])


	return faceRecognitionMap

}

//handerMsg
func HandlerMsg(bts []byte) {
	//logs.Info("-------------- enter handler -----------------------")
	var values ReqAndRes
	if err := json.Unmarshal(bts, &values); nil != err {
		logs.Error("parse kafka msg error: %s", err)
		return
	}

	apiModel := &models.ApiModel{}
	appModel := &models.AppModel{}

	isAuth := values.MyResponse.RespHeaders.UbdAuth

	//logs.Info("--------------------------------",isAuth)
	if isAuth != "1" { // 未鉴权
		//logs.Info(1)
		myUrl := values.MyRequest.Url

		// api 处理
		urlParse, err := url.Parse(myUrl)
		if err != nil {
			logs.Error(err)
		}
		path := urlParse.Path
		//logs.Info("---path--->",path)
		//logs.Info("app->",app,"  path-->",path,"---myUrl--->",myUrl,"---path--->",path)

		if strings.Index(path, smartDoor) >= 0 {
			apiModel = getApiModel(path, smartDoor, apiModel)

			if apiModel.Id==0{
				return
			}

			appModel.QueryTable().Filter("desc", smartDoor).One(appModel, "id", "app_name","domain")


		} else if strings.Index(path, smartCommunity) >= 0 {
			apiModel = getApiModel(path, smartCommunity, apiModel)

			if apiModel.Id==0{
				return
			}

			appModel.QueryTable().Filter("desc", smartCommunity).One(appModel, "id", "app_name","domain")

		} else if strings.Index(path, smartCampus) >= 0 {
			apiModel = getApiModel(path, smartCampus, apiModel)

			if apiModel.Id==0{
				return
			}

			appModel.QueryTable().Filter("desc", smartCampus).One(appModel, "id", "app_name","domain")


		} else { // 人脸识别 鉴权
				if faceRecognitionMap[path]>0 {

					apiModel.QueryTable().Filter("id",faceRecognitionMap[path]).One(apiModel,"id","api_name")
					appModel.QueryTable().Filter("desc", "/face_recognition").One(appModel, "id", "app_name","domain")
					//logs.Info(appModel.Id, apiModel.Id, values.MyResponse.Status, values.Ts, values.MyRequest.Url, appModel.AppName, apiModel.ApiName,appModel.Domain)
				} else{  //  非原子能力的api不统计
					// 暂时保留吧
					return
					//if  path=="/openapi/cv/id/v1/detect" {
					//	return
					//}

					//appInt, _ := strconv.Atoi(values.MyResponse.RespHeaders.UbdApp)
					//apiInt, _ := strconv.Atoi(values.MyResponse.RespHeaders.XUbdApi)
					//apiModel.QueryTable().Filter("id",faceRecognitionMap[path]).One(apiModel,"id","api_name")
					//apiModel.QueryTable().Filter("url", path).One(apiModel, "id", "api_name")
					//logs.Info("  ---- undefined ------",path,myUrl," app=",values.MyResponse.RespHeaders.UbdApp,"  api=",values.MyResponse.RespHeaders.XUbdApi,
					//	"  isAuth=",isAuth," status=",values.MyResponse.Status,"       from model-->",apiModel.Id,apiModel.ApiName)
					//return
			}

			//logs.Info("  ------- undefined -----------------------",path)
			// 人脸的不鉴权
			//  非原子能力的api不统计
			//err = apiModel.QueryTable().Filter("url_desc", "/aicubigdatacn/undefined").One(apiModel, "id", "api_name")
			//if err != nil {
			//	logs.Error(err)
			//}
			//appModel.QueryTable().Filter("desc", "/aicubigdatacn/undefined").One(appModel, "id", "app_name","domain")
		}

	} else {	// 鉴权
		//return
		appInt, _ := strconv.Atoi(values.MyResponse.RespHeaders.UbdApp)
		apiInt, _ := strconv.Atoi(values.MyResponse.RespHeaders.XUbdApi)

		apiModel.QueryTable().Filter("id", apiInt).One(apiModel, "id", "api_name")
		appModel.QueryTable().Filter("id", appInt).One(appModel, "id", "app_name","domain")

		//logs.Info(values.MyResponse.RespHeaders.XUbdApi,"  apiInt =",apiInt,"  apiModel.id=",apiModel.Id)
	}
	//  "app=",values.MyResponse.RespHeaders.UbdApp,
	//logs.Info(appModel.Id, apiModel.Id, values.MyResponse.Status, values.Ts, values.MyRequest.Url, appModel.AppName, apiModel.ApiName,appModel.Domain," isAuth=",isAuth)

	err := InsertAccessLogKafka(appModel.Id, apiModel.Id, values.MyResponse.Status, values.Ts, values.MyRequest.Url, appModel.AppName, apiModel.ApiName,appModel.Domain)
	if err != nil {
		logs.Error("err:", err)
		return
	}

	//logs.Info("-------------- exit handler -----------------------")

	return
}


func getApiModel(path, business string, apiModel *models.ApiModel) *models.ApiModel {
	if strings.Index(path, maskFace) >= 0 && strings.Index(path, search) >= 0 {
		apiModel.QueryTable().Filter("url_desc", business+maskFaceSearch).One(apiModel, "id", "api_name")
	} else if strings.Index(path, face) >= 0 && strings.Index(path, search) >= 0 {
		apiModel.QueryTable().Filter("url_desc", business+faceSearch).One(apiModel, "id", "api_name")
	} else if strings.Index(path, deviceReport) >= 0 {
		apiModel.QueryTable().Filter("url_desc", business+deviceReport).One(apiModel, "id", "api_name")
	} else {  //  心跳不统计
		//return
		//apiModel.QueryTable().Filter("url_desc", business+other).One(apiModel, "id", "api_name")
	}
	return apiModel
}











// FunMsg 算法回掉消息结构体
type FunMsg struct {
	CameraID   string        `json:"camera_id"`
	FrameIndex int           `json:"frame_index"`
	FunType    int           `json:"function_type"`
	FunName    string        `json:"function_name"`
	ImgPath    string        `json:"image_path"`
	Items      []*ActionItem `json:"result"`
	TimeLen    int64         `json:"time_sec"`
	StartTime  string        `json:"start_time"`
	EndTime    string        `json:"end_time"`
}

// ActionItem 单条报警信息
type ActionItem struct {
	ID       string `json:"id"`
	Type     int    `json:"action_type"`
	Name     string `json:"action_name"`
	Location struct {
		Left   float64 `json:"left"`
		Top    float64 `json:"top"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"location"`
	Score float64 `json:"score"`
}



//handerMsg
func HandlerMsgNoSql(bts []byte) {
	//logs.Info("-------------- enter handler -----------------------")
	var values ReqAndRes
	if err := json.Unmarshal(bts, &values); nil != err {
		logs.Error("parse kafka msg error: %s", err)
		return
	}
	//logs.Info(values)
	//logs.Info(values, "--status-->", values.MyResponse.Status)
	//logs.Info("--- ts ---",values.Ts)
	ts := uint64(time.Now().Unix())


	var app uint64
	var api uint64
	//apiModel := &models.ApiModel{}
	//appModel := &models.AppModel{}

	isAuth := values.MyResponse.RespHeaders.UbdAuth

	if isAuth != "1" { // 未认证
		myUrl := values.MyRequest.Url

		// api 处理
		urlParse, err := url.Parse(myUrl)
		if err != nil {
			logs.Error(err)
		}
		path := urlParse.Path
		//logs.Info("app->",app,"  path-->",path,"---myUrl--->",myUrl,"---path--->",path)

		if strings.Index(path, smartDoor) >= 0 {

			if strings.Index(path, maskFace) >= 0 && strings.Index(path, search) >= 0 {
				api=114
			} else if strings.Index(path, face) >= 0 && strings.Index(path, search) >= 0 {
				api=115
			} else if strings.Index(path, deviceReport) >= 0 {
				api=116
			} else {
				api=117
			}
			app=64


		} else if strings.Index(path, smartCommunity) >= 0 {

			if strings.Index(path, maskFace) >= 0 && strings.Index(path, search) >= 0 {
				api=118
			} else if strings.Index(path, face) >= 0 && strings.Index(path, search) >= 0 {
				api=119
			} else if strings.Index(path, deviceReport) >= 0 {
				api=120
			} else {
				api=121
			}
			app=65

		} else if strings.Index(path, smartCampus) >= 0 {

			if strings.Index(path, maskFace) >= 0 && strings.Index(path, search) >= 0 {
				api=122
			} else if strings.Index(path, face) >= 0 && strings.Index(path, search) >= 0 {
				api=123
			} else if strings.Index(path, deviceReport) >= 0 {
				api=124
			} else {
				api=125
			}
			app=66

		} else {
			api=126
			app=67

		}
		//api = apiModel.Id
		//app = appModel.Id

	} else {


		apiInt, _ := strconv.Atoi(values.MyResponse.RespHeaders.XUbdApi)
		api=uint64(apiInt)

		appInt, _ := strconv.Atoi(values.MyResponse.RespHeaders.UbdApp)
		app=uint64(appInt)


		//apiModel.QueryTable().Filter("id", appInt).One(apiModel, "id", "api_name")
		//appModel.QueryTable().Filter("id", apiInt).One(appModel, "id", "app_name")

	}

	//logs.Info("--  InsertAccessLogKafka -- ", appModel.Id, apiModel.Id, values.MyResponse.Status, ts, values.MyRequest.Url, appModel.AppName, apiModel.ApiName)
	//logs.Info(app,api)

	logs.Info(app,api)
	logs.Info(ts)
	//err := InsertAccessLogKafka(app, api, values.MyResponse.Status, ts, values.MyRequest.Url)
	//if err != nil {
	//	logs.Error("err:", err)
	//	return
	//}

	//logs.Info("-------------- exit handler -----------------------")

	return
}

