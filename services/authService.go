package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"openai-backend/utils/redis"
	//"github.com/go-redis/redis/v8"
	"net/url"
	"openai-backend/models"
	"openai-backend/utils/httpUtil"
	"strconv"
	"strings"
	"time"
)

//var Rdb = redis.NewClient(&redis.Options{
//	Addr:     beego.AppConfig.DefaultString("cache_redis_host", ""),
//	Password: beego.AppConfig.DefaultString("cache_redis_password", ""), // no password set
//	DB:       4,                                                         // use default DB
//})
var ctx = context.Background()

var hashKey = "__openai__authcache__"

type AuthCacheResp struct {
	App    uint64 `json:"app"`
	Api    uint64 `json:"api"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Qps    uint64 `json:"qps"`
	Quota  uint64 `json:"quota"`
}

var authCacheResp AuthCacheResp

func CacheHandler(key, url string) (AuthCacheResp, error) {
	// 1. app->url
	//app
	appModel := &models.AppModel{}
	err := models.NewAppModel().QueryTable().Filter("api_key", key).One(appModel)
	if err != nil {
		logs.Error(err)
		return authCacheResp, err
	}
	authCacheResp.App = appModel.Id
	authCacheResp.Key = key
	authCacheResp.Secret = appModel.SecretKey
	// app->url
	appAbilityList := make([]*models.AppAbilityModel, 0)
	_, err = models.NewAppAbilityModel().QueryTable().Filter("app_id", appModel.Id).All(&appAbilityList, "ability_id")
	if err != nil {
		logs.Error(err)
		return authCacheResp, err
	}
	abilityList := make([]uint64, 0)
	for _, appAbility := range appAbilityList {
		abilityList = append(abilityList, appAbility.AbilityId)
	}
	//  解决 缓存穿透问题
	if len(abilityList)==0{

		redisKey := key + "_" + "no_url_can_access"
		redisCache := map[string]interface{}{
			"app":    appModel.Id,
			"key":    key,
			"secret": appModel.SecretKey,
		}
		bytes, _ := json.Marshal(&redisCache)
		_, err = redis.C().HSet(ctx, hashKey, map[string]interface{}{
			redisKey: string(bytes),
		}).Result()

		return authCacheResp, errors.New("key:" + key + " no_url_can_access")
	}
	ok := models.NewApiModel().QueryTable().Filter("ability_id__in", abilityList).Filter("url", url).Exist()

	if !ok { // no auth
		logs.Error("key:" + key + " 没有权限访问" + url)
		return authCacheResp, errors.New("key:" + key + " 没有权限访问" + url)
	}

	// 2.app、api -> quotas
	// api
	apiModel := &models.ApiModel{}
	err = models.NewApiModel().QueryTable().Filter("url", url).One(apiModel)
	if err != nil {
		logs.Error(err)
		return authCacheResp, err
	}
	//quotas
	quotas := &models.QuotasModel{}
	err = models.NewQuotasModel().QueryTable().Filter("app_id", appModel.Id).Filter("api_id", apiModel.Id).One(quotas)
	if err != nil {
		logs.Error(err)
		return authCacheResp, err
	}

	authCacheResp.Qps = quotas.MaxQps
	authCacheResp.Quota = quotas.Quotas
	authCacheResp.Api = apiModel.Id

	redisKey := key + "_" + url
	//var value redisCache
	//rdb.HSet(ctx,hashKey,redisKey,value)

	//rdb.HSet(ctx,hashKey,"app",appModel.Id)
	//rdb.HSet(ctx,hashKey,"key",key)
	//rdb.HSet(ctx,hashKey,"secret",appModel.SecretKey)
	//rdb.HSet(ctx,hashKey,"qps",quotas.MaxQps)
	//rdb.HSet(ctx,hashKey,"quota",quotas.Quotas)

	logs.Info( " redisKey ",redisKey," app ",appModel.Id, " api ",apiModel.Id, " key ",key, " secret ",appModel.SecretKey, " qps ",quotas.MaxQps, " quota ",quotas.Quotas)
	redisCache := map[string]interface{}{
		"app":    appModel.Id,
		"api":    apiModel.Id,
		"key":    key,
		"secret": appModel.SecretKey,
		"qps":    quotas.MaxQps,
		"quota":  quotas.Quotas,
	}
	logs.Info("redisCache---->", redisCache, "-->redisCache:", redisKey)
	bytes, _ := json.Marshal(&redisCache)
	//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
	_, err = redis.C().HSet(ctx, hashKey, map[string]interface{}{
		redisKey: string(bytes),
	}).Result()

	if err != nil {
		logs.Error("---->", err)
		return authCacheResp, err
	}
	return authCacheResp, err
}









/*
{
  "iss": "1234567890",
  "exp": 2147483647,
  "iat": 1516239022
}
*/
// MyCustomClaims 自定义Claims
type MyCustomClaim struct {
	jwt.StandardClaims
}

type AuthResp struct {
	UbdApi uint64 `json:"ubd_api"`
	UbdKey string `json:"ubd_key"`
	UbdApp uint64 `json:"ubd_app"`
}

var authResp AuthResp

func VerifyToken(headers map[string]string, path, rowQuery, method string) (AuthResp, bool, error) {
	token, err := getToken(headers, rowQuery, method)
	if err != nil {
		logs.Error("Token error=", err)
		return authResp, false, err
	}

	authResp, flag, err := TokenHandle(path, token)
	if err != nil {
		logs.Error("Token error=", err)
		return authResp, false, err
	}
	return authResp, flag, err
}

func getToken(headers map[string]string, rowQuery, method string) (string, error) {
	var token string
	if strings.EqualFold(method, "get") {
		data, err := url.ParseQuery(rowQuery)
		if err != nil {
			return token, err
		}
		return data.Get("token"), nil
	} else if strings.EqualFold(method, "post") {
		if authorization, ok := headers["authentication"]; ok {
			if strings.Contains(authorization, "Bearer") {
				if len(authorization) > 8 {
					token = authorization[7:]
					return token, nil
				}
			} else {
				return token, errors.New(httpUtil.ERROR_TOKEN_NOT_EXIST)
			}
		} else {
			return token, errors.New(httpUtil.ERROR_TOKEN_NOT_EXIST)
		}
	}
	return token, errors.New("check your 'method' param")
}
func TokenHandle(url, token string) (AuthResp, bool, error) {
	logs.Info("---url ---", url, "-----token-----", token)
	jt, err := jwt.ParseWithClaims(token, &MyCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		myClaims, ok := t.Claims.(*MyCustomClaim)
		if !ok {
			logs.Error(" ok :", ok)
			return nil, fmt.Errorf("claims parse error")
		}
		//myClaims.StandardClaims.ExpiresAt //16xx  1s
		return getSecretByKey(url, myClaims.StandardClaims.Issuer) // sk
	})
	logs.Info("---jt.Valide ---", jt.Valid)
	if jt.Valid {
		return authResp, true, nil
	}
	return authResp, false, err
}


//顺序 key是否存在：是->取出sk，app_id->取出ability的集合abilityList-->取出api的集合，判断url是否存在。即该key是否有权限访问此url。
//secret_key: url--> api_info 得到 ability_id --> app_ability_info 得到 app_id_List --> 判断 api_key是否存在 ，存在，获取对应的 secret_key
func getSecretByKey(url string, apiKey string) ([]byte, error) { // secret_key
	logs.Info("----enter getSecretByKey ---")
	apiModel := &models.ApiModel{}
	err := models.NewApiModel().QueryTable().Filter("url", url).One(apiModel, "ability_id", "id")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	var appAbilityList = make([]models.AppAbilityModel, 0)
	_, err = models.NewAppAbilityModel().QueryTable().Filter("ability_id", apiModel.AbilityId).All(&appAbilityList, "app_id")
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	var appIds = make([]uint64, 0)
	for _, appAbility := range appAbilityList {
		appIds = append(appIds, appAbility.AppId)
	}
	appModel := &models.AppModel{}
	ok := models.NewAppModel().QueryTable().Filter("id__in", appIds).Filter("api_key", apiKey).Exist()
	if !ok {
		return nil, errors.New("获取secret key 失败")
	}
	err = models.NewAppModel().QueryTable().Filter("id__in", appIds).Filter("api_key", apiKey).One(appModel, "secret_key", "id")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	//var  authResp  AuthResp
	authResp.UbdApi = apiModel.Id
	authResp.UbdApp = appModel.Id
	authResp.UbdKey = apiKey

	return []byte(appModel.SecretKey), nil

}

func QpsAndQuotasLimitCheck(authResp AuthResp) error {
	// 1. [鉴权中拿到]app_id，api_id 。查询 quotas_info，取出 max_qps 、 quotas
	quotas := &models.QuotasModel{}
	err := models.NewQuotasModel().QueryTable().Filter("app_id", authResp.UbdApp).Filter("api_id", authResp.UbdApi).One(quotas)
	if err != nil {
		logs.Error(err, "app_id:", authResp.UbdApp, "api_id", authResp.UbdApi)
		return err
	}
	// 从mysql中拿到 最大qps限制 和 计次
	maxQpsMysql := int(quotas.MaxQps)
	quotasMysql := int(quotas.Quotas)

	qpsRedisKey := strconv.FormatUint(quotas.AppId, 10) + "_" + strconv.FormatUint(quotas.ApiId, 10) + "_qps"
	quotaRedisKey := strconv.FormatUint(quotas.AppId, 10) + "_" + strconv.FormatUint(quotas.ApiId, 10) + "_quota"

	// qps
	err = validateQps(qpsRedisKey, maxQpsMysql, time.Second)
	if err != nil {
		return err
	}

	// quota
	err = validateQps(quotaRedisKey, quotasMysql, time.Hour*24)
	if err != nil {
		return err
	}
	return nil
}

func validateQps(redisKey string, quotas int, exp time.Duration) error {
	intCmd := redis.C().Incr(ctx, redisKey)

	if intCmd.Err() != nil {
		//panic(err)
		logs.Error(intCmd.Err())
		return intCmd.Err()
	}
	logs.Info("redisKey->", redisKey, "   v->", intCmd.Val())
	if intCmd.Val() > int64(quotas) {
		logs.Error("qps超出最大值限制")
		return errors.New("qps超出最大值限制")
	}
	if intCmd.Val() == 1 {
		redis.C().Expire(ctx, redisKey, exp)
	}
	return nil

}






//var mu sync.Mutex


//// 分布式锁
//redisLockKey:=strconv.FormatUint(appId,10)+strconv.FormatUint(apiId,10)
//logs.Info("------------enter Lock -------------------------------")
//logs.Info(redisLockKey)
//if apiId==126{
//	logs.Info(appId , apiId , status, ts)
//}
//unlock,err := redis.Lock(ctx,redisLockKey,time.Millisecond*200)
//logs.Info("------------exit Lock -------------------------------")
//if err !=nil{
//	logs.Error(err)
//	return err
//}
//defer unlock()

//mu.Lock()
//defer mu.Unlock()




