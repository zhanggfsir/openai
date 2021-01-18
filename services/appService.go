package services

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"openai-backend/models"
	"openai-backend/utils"
	"openai-backend/utils/redis"
	"openai-backend/views"
)

/*
@Author:
@Time: 2020-11-05 15:02
@Description:APPmodel
*/

//新建 服务
func ValidatePostApp(domain string, requestBody []byte) (*models.AppModel, error) {
	//获取requestBody中的参数
	type requestJSON struct {
		AppName     string   `description:"App名称" json:"appName"`
		AbilityList []uint64 `description:"api列表" json:"abilityList"`
		Desc        string   `description:"备注" json:"desc"`
	}
	var jSON requestJSON

	if err := json.Unmarshal(requestBody, &jSON); err != nil {
		logs.Error("err=", err)
		return nil, errors.New("json 格式错误" + err.Error())
	}

	logs.Info(jSON)
	//校验数据
	valid := validation.Validation{}
	valid.Required(jSON.AppName, "appName").Message("未输入名称 ")
	valid.Required(jSON.AbilityList, "dataType").Message("未选择能力 ")

	valid.MaxSize(jSON.Desc, 100, "desc").Message("描述过长 ")

	if valid.HasErrors() {
		var errStr string
		for _, err := range valid.Errors {
			errStr = errStr + err.Message
		}
		//	logs.Error("err", errStr)
		return nil, errors.New(errStr)

	}

	//分配KEY
	var (
		apiKey    string
		secretKey string
	)
	for {
		apiKey = utils.GenerateRandom("mix", 24)
		secretKey = utils.GenerateRandom("mix", 32)
		if models.NewAppModel().QueryTable().Filter("api_key", apiKey).Exist() != true &&
			models.NewAppModel().QueryTable().Filter("secret_key", secretKey).Exist() != true {
			break
		}
	}

	logs.Info(apiKey, len(apiKey), "    ", secretKey, "    ", len(secretKey))
	//将json数据复制给model
	app := &models.AppModel{
		AppName:   jSON.AppName,
		ApiKey:    apiKey,
		SecretKey: secretKey,
		Desc:      jSON.Desc,
		Domain:    domain,
	}

	if app.QueryTable().Filter("app_name", app.AppName).Exist() == true {
		return nil, errors.New("APP名称重复")
	}

	ability := models.NewAbilityModel()
	for _, abilityID := range jSON.AbilityList {
		if ability.QueryTable().Filter("id", abilityID).Exist() != true {
			return nil, errors.New("能力不存在")
		}
	}

	o := orm.NewOrm()

	if err := o.Begin(); err != nil {
		logs.Error(err)
		return nil, errors.New("开启事务错误")
	}

	defer o.Commit()
	//插入app_info数据库
	if _, err := o.Insert(app); err != nil {
		_ = o.Rollback()
		logs.Error(err)
		return nil, errors.New("插入数据库失败 " + err.Error())
	}

	//1.将ability列表中的API写入数据库quotas_info

	for _, abilityID := range jSON.AbilityList {
		//写入APP 能力关系表
		appAbilityModel := &models.AppAbilityModel{}
		appAbilityModel.AppId = app.Id
		appAbilityModel.AbilityId = abilityID
		if _, err := o.Insert(appAbilityModel); err != nil {
			_ = o.Rollback()
			logs.Error(err)
			return nil, errors.New("插入数据库失败 " + err.Error())
		}
		//写入quotas_info表
		api := &models.ApiModel{}
		apiList := make([]*models.ApiModel, 0)
		_, _ = api.QueryTable().Filter("ability_id", abilityID).All(&apiList)

		for _, api := range apiList {
			defaultQuotas := &models.DefaultQuotasModel{}
			_, _ = defaultQuotas.QueryTable().Filter("api_id", api.Id).All(defaultQuotas)

			quotas := &models.QuotasModel{}
			quotas.AppId = app.Id
			quotas.ApiId = api.Id
			quotas.MaxQps = defaultQuotas.MaxQps
			quotas.Quotas = defaultQuotas.Quotas
			quotas.QuotasPeriod = defaultQuotas.QuotasPeriod
			quotas.Status = "0"

			if _, err := o.Insert(quotas); err != nil {
				_ = o.Rollback()
				logs.Error(err)
				return nil, errors.New("插入数据库失败 " + err.Error())
			}
			// 写入redis
			//Rdb.HSet()
		}
	}

	return app, nil
}

func GetAppDetailService(domain string, appId uint64) (map[string]interface{}, error) {
	app := &models.AppModel{}
	if app.QueryTable().Filter("domain", domain).Filter("id", appId).Exist() != true {
		return nil, errors.New("app does not exist ")
	}
	_, _ = app.QueryTable().Filter("domain", domain).Filter("id", appId).All(app)
	logs.Info(app)

	abilityList := views.AppsAbilityDetailInfo(app)
	var jsonData = make(map[string]interface{})
	jsonData["appName"] = app.AppName
	jsonData["apiKey"] = app.ApiKey
	jsonData["secretKey"] = app.SecretKey
	jsonData["desc"] = app.Desc
	jsonData["createdTime"] = app.CreateTime.Format("2006-01-02 15:04:05")
	jsonData["updatedTime"] = app.ModifyTime.Format("2006-01-02 15:04:05")
	jsonData["abilityList"] = abilityList
	return jsonData, nil
}

func GetAppListService(domain string, pageNo, pageSize int, queryFilter map[string]interface{}) (map[string]interface{}, error) {
	model := &models.AppModel{}
	lists, Count, err := model.FindToPager(domain, pageNo, pageSize, queryFilter)
	if err != nil {
		logs.Error("List error", err)
		return nil, err
	}
	//logs.Info(lists, Count)

	datalist := views.AppList(lists)
	var jsonData = make(map[string]interface{})
	jsonData["list"] = datalist
	jsonData["pageNo"] = pageNo
	jsonData["pageSize"] = pageSize
	jsonData["total"] = Count
	return jsonData, nil
}

//删除APP服务
func ValidateDeleteApp(domain string, requestBody []byte) (*models.AppModel, []uint64, error) {
	//获取requestBody中的参数
	type requestJSON struct {
		Ids []uint64 `description:"APPID" json:"id"`
	}
	var jSON requestJSON

	if err := json.Unmarshal(requestBody, &jSON); err != nil {
		logs.Error("err=", err)
		return nil, nil, errors.New("json 格式错误" + err.Error())
	}
	//校验数据

	logs.Info(jSON)
	if len(jSON.Ids) == 0 {
		return nil, nil, errors.New("未输入ID ")
	}

	app := &models.AppModel{}

	for _, id := range jSON.Ids {

		if app.QueryTable().Filter("domain",domain).Filter("id", id).Exist() != true {
			return nil, nil, errors.New("ID不存在")
		}
	}

	return app, jSON.Ids, nil
}

//修改 服务
func ValidateModifyApp(requestBody []byte) (*models.AppModel, error) {
	//获取requestBody中的参数
	type requestJSON struct {
		AppID       uint64   `description:"AppId" json:"appId"`
		AbilityList []uint64 `description:"列表"   json:"abilityList"`
	}
	var jSON requestJSON

	if err := json.Unmarshal(requestBody, &jSON); err != nil {
		logs.Error("err=", err)
		return nil, errors.New("json 格式错误" + err.Error())
	}

	logs.Info(jSON)
	//校验数据
	valid := validation.Validation{}
	valid.Required(jSON.AbilityList, "dataType").Message("未选择能力 ")

	if valid.HasErrors() {
		var errStr string
		for _, err := range valid.Errors {
			errStr = errStr + err.Message
		}
		//	logs.Error("err", errStr)
		return nil, errors.New(errStr)

	}

	//将json数据复制给model
	app := &models.AppModel{
		Id: jSON.AppID,
	}

	ability := models.NewAbilityModel()
	for _, abilityID := range jSON.AbilityList {
		if ability.QueryTable().Filter("id", abilityID).Exist() != true {
			return nil, errors.New("能力不存在")
		}
	}

	if app.QueryTable().Filter("id", app.Id).Exist() != true {
		return nil, errors.New("app不存在")
	}
	err := app.QueryTable().Filter("id", app.Id).One(app)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("开启事务错误")
	}

	o := orm.NewOrm()

	if err := o.Begin(); err != nil {
		logs.Error(err)
		return nil, errors.New("开启事务错误")
	}

	defer o.Commit()

	//1.将ability列表中的API写入数据库quotas_info

	for _, abilityID := range jSON.AbilityList {
		if models.NewAppAbilityModel().QueryTable().Filter("app_id", app.Id).Filter("ability_id", abilityID).Exist() != true {
			//写入APP 能力关系表 app_ability_info
			appAbilityModel := &models.AppAbilityModel{}
			appAbilityModel.AppId = app.Id
			appAbilityModel.AbilityId = abilityID
			if _, err := o.Insert(appAbilityModel); err != nil {
				_ = o.Rollback()
				logs.Error(err)
				return nil, errors.New("插入数据库失败 " + err.Error())
			}
			//写入quotas_info表
			api := &models.ApiModel{}
			apiList := make([]*models.ApiModel, 0)
			_, _ = api.QueryTable().Filter("ability_id", abilityID).All(&apiList)

			for _, api := range apiList {
				defaultQuotas := &models.DefaultQuotasModel{}
				_, _ = defaultQuotas.QueryTable().Filter("api_id", api.Id).All(defaultQuotas)

				quotas := &models.QuotasModel{}
				quotas.AppId = app.Id
				quotas.ApiId = api.Id
				quotas.MaxQps = defaultQuotas.MaxQps
				quotas.Quotas = defaultQuotas.Quotas
				quotas.QuotasPeriod = defaultQuotas.QuotasPeriod
				quotas.Status = "0"

				if _, err := o.Insert(quotas); err != nil {
					_ = o.Rollback()
					logs.Error(err)
					return nil, errors.New("插入数据库失败 " + err.Error())
				}
				// redis
				redisKey := app.ApiKey + "_" + api.Url
				redisCache := map[string]interface{}{
					"app":    app.Id,
					"api":    api.Id,
					"key":    app.ApiKey,
					"secret": app.SecretKey,
					"qps":    quotas.MaxQps,
					"quota":  quotas.Quotas,
				}
				logs.Info("redisCache---->", redisCache, "-->redisCache:", redisKey)
				bytes, _ := json.Marshal(&redisCache)
				_, err = redis.C().HSet(ctx, hashKey, map[string]interface{}{
					redisKey: string(bytes),
				}).Result()

				if err != nil {
					logs.Error("---->", err)
					return app, err
				}

			}
		}
	}

	return app, nil
}
