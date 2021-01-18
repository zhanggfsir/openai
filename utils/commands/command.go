package commands

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	beegoCache "github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/cache/redis"
	beelogs "github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"log"
	"net/url"
	"openai-backend/models"
	"openai-backend/services"
	"openai-backend/utils"
	"openai-backend/utils/cacheUtil"
	"openai-backend/utils/casbinUtil"
	"openai-backend/utils/commonUtil"
	"openai-backend/utils/httpUtil"
	"openai-backend/utils/stringUtil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	if configPath, err := filepath.Abs(utils.ConfigurationFile); err == nil {
		utils.ConfigurationFile = configPath
	}
}

// RunCommand 注册orm命令行工具
func Bootstrap() {
	macCheck := beego.AppConfig.DefaultBool("mac_check", false)
	if macCheck {
		if !InitMac() {
			os.Exit(1)
		}
	}

	ResolveCommand()
	kafkaOpen := beego.AppConfig.DefaultBool("kafka_open_flag", false)
	if kafkaOpen {
		//beelogs.Info("----start kafka----")
		services.StartKafka()
	}


	window:=beego.AppConfig.DefaultInt64("open_ai_window", 10)
	go func() {
		for  {
			startTime:=time.Now().Add(-time.Second*time.Duration(window)).Format("2006-01-02 15:04:05")
			//time.Sleep(time.Duration(10)*time.Second)
			endTime:=time.Now().Format("2006-01-02 15:04:05")
			// need time
			go services.DataStatistics(startTime,endTime)
			time.Sleep(time.Duration(window)*time.Second)
		}
	}()


	//window:=beego.AppConfig.DefaultInt64("open_ai_window", 10)
	//go func() {
	//	for  {
	//		startTime:=time.Now().Format("2006-01-02 15:04:05")
	//		time.Sleep(time.Duration(window)*time.Second)
	//		endTime:=time.Now().Format("2006-01-02 15:04:05")
	//
	//
	//	}
	//}()

}

func InitMac() bool {
	mac:=beego.AppConfig.DefaultString("mac", "00:ff:97:8d:c2:16")
	macEnhance:=mac+"Copyright cubigdata.cn All Rights Reserved"
	macMd5 := stringUtil.Md5(macEnhance)

	url_ := beego.AppConfig.DefaultString("macUrl", "https://172.16.222.253:9031/mac")
	domain := ""
	requestData := make(map[string]interface{})
	response := httpUtil.PostRequest(url_, domain, requestData, "application/json; charset=utf-8")
	//beelogs.Info("---response-->",response)
	if response["mac"] == macMd5{
		return true
	}else {
		beelogs.Error("-- MAC校验错误,请联系管理员 --")
		return false
	}

}
//解析命令
func ResolveCommand() {
	if err := beego.LoadAppConfig("ini", utils.ConfigurationFile); err != nil {
		log.Fatal("An error occurred:", err)
	}
	if utils.LogFile == "" {
		logPath, err := filepath.Abs(beego.AppConfig.DefaultString("log_path",
			utils.WorkingDir("runtime", "logs")))
		if err == nil {
			utils.LogFile = logPath
		} else {
			utils.LogFile = utils.WorkingDir("runtime", "logs")
		}
	}
	utils.AutoLoadDelay = beego.AppConfig.
		DefaultInt("config_auto_delay", 0)
	beego.BConfig.WebConfig.StaticDir["/static"] = filepath.
		Join(utils.WorkingDirectory, "static")
	beego.BConfig.WebConfig.ViewsPath = utils.WorkingDir("views")

	beego.BConfig.MaxMemory=1<<28   // 文件上传默认内存缓存大小，默认值是 1 << 26(64M)。
	RegisterDataBase()
	RegisterLogger(utils.LogFile)
	RegisterCache()
	RegisterModel()
	initFunc()
}

// RegisterDataBase 注册数据库
func RegisterDataBase() {
	adapter := beego.AppConfig.String("db_adapter")
	orm.DefaultTimeLoc = time.Local

	if strings.EqualFold(adapter, "mysql") {
		host := beego.AppConfig.String("db_host")
		database := beego.AppConfig.String("db_database")
		username := beego.AppConfig.String("db_username")
		password := beego.AppConfig.String("db_password")

		timezone := beego.AppConfig.String("timezone")
		location, err := time.LoadLocation(timezone)
		if err == nil {
			orm.DefaultTimeLoc = location
		} else {
			beelogs.Error("加载时区配置信息失败,请检查是否存在 ZONEINFO 环境变量->", err)
		}

		port := beego.AppConfig.String("db_port")

		dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s", username, password, host, port, database, url.QueryEscape(timezone))
		fmt.Println(port, password, dataSource)
		if err := orm.RegisterDataBase("default", "mysql", dataSource); err != nil {
			beelogs.Error("注册默认数据库失败->", err)
			os.Exit(1)
		}

	} else {
		beelogs.Error("不支持的数据库类型.")
		os.Exit(1)
	}

	beelogs.Info("数据库初始化完成.")
}

//注册缓存管道
func RegisterCache() {
	isOpenCache := beego.AppConfig.DefaultBool("cache", false)
	if !isOpenCache {
		cacheUtil.Init(&cacheUtil.NullCache{})
		return
	}
	cacheProvider := beego.AppConfig.String("cache_provider")
	if cacheProvider == "redis" {
		//设置Redis前缀
		if key := beego.AppConfig.
			DefaultString("cache_redis_prefix", ""); key != "" {
			redis.DefaultKey = key
		}
		redis.DefaultKey = beego.AppConfig.DefaultString("cache_redis_prefix", "")
		var redisConfig struct {
			Conn     string `json:"conn"`
			Password string `json:"password"`
			DbNum    int    `json:"dbNum"`
		}
		redisConfig.DbNum = 0
		redisConfig.Conn = beego.AppConfig.DefaultString("cache_redis_host", "")
		if pwd := beego.AppConfig.DefaultString("cache_redis_password", ""); pwd != "" {
			redisConfig.Password = pwd
		}
		if dbNum := beego.AppConfig.DefaultInt("cache_redis_db", 0); dbNum > 0 {
			redisConfig.DbNum = dbNum
		}

		bc, err := json.Marshal(&redisConfig)
		if err != nil {
			beelogs.Error("初始化Redis缓存失败:", err)
			os.Exit(1)
		}
		redisCache, err := beegoCache.NewCache("redis", string(bc))

		if err != nil {
			beelogs.Error("初始化Redis缓存失败:", err)
			os.Exit(1)
		}
		cacheUtil.Init(redisCache)
	} else {
		cacheUtil.Init(&cacheUtil.NullCache{})
		beelogs.Error("不支持的缓存管道,缓存将禁用 ->", cacheProvider)
		return
	}
	beelogs.Info("缓存初始化完成.")
}

// RegisterLogger 注册日志
func RegisterLogger(log string) {
    logs := beelogs.GetBeeLogger()
	logs.SetLogger("console")
	//如果你期望输出调用的文件名和文件行号
	logs.SetLogFuncCallDepth(3)
	logs.Async(1e3)

	if beego.AppConfig.DefaultBool("log_is_async", true) {
		logs.Async(1e3)
	}
	if log == "" {
		logPath, err := filepath.Abs(beego.AppConfig.DefaultString("log_path",
			utils.WorkingDir("runtime", "logs")))
		if err == nil {
			log = logPath
		} else {
			log = utils.WorkingDir("runtime", "logs")
		}
	}

	logPath := filepath.Join(log, "log.log")

	if _, err := os.Stat(log); os.IsNotExist(err) {
		//_ = os.MkdirAll(log, 0755)
		_ = commonUtil.CreateDir(log)
	}

	config := make(map[string]interface{}, 1)

	config["filename"] = logPath
	config["perm"] = "0755"
	config["rotate"] = true

	if maxLines := beego.AppConfig.DefaultInt("log_maxlines",
		1000000); maxLines > 0 {
		config["maxLines"] = maxLines
	}
	config["maxsize"] = 1 << 28
	if !beego.AppConfig.DefaultBool("log_daily", true) {
		config["daily"] = false
	}
	if maxDays := beego.AppConfig.DefaultInt("log_maxdays",
		7); maxDays > 0 {
		config["maxdays"] = maxDays
	}

	if level := beego.AppConfig.DefaultString("log_level", "Trace"); level != "" {
		switch level {
		case "Emergency":
			config["level"] = beelogs.LevelEmergency
			break
		case "Alert":
			config["level"] = beelogs.LevelAlert
			break
		case "Critical":
			config["level"] = beelogs.LevelCritical
			break
		case "Error":
			config["level"] = beelogs.LevelError
			break
		case "Warning":
			config["level"] = beelogs.LevelWarning
			break
		case "Notice":
			config["level"] = beelogs.LevelNotice
			break
		case "Informational":
			config["level"] = beelogs.LevelInformational
			break
		case "Debug":
			config["level"] = beelogs.LevelDebug
			break
		}
	}
	b, err := json.Marshal(config)
	if err != nil {
		panic(err)
	} else {
		logs.SetLogger(beelogs.AdapterFile, string(b))
	}
}

func RegisterModel() {
	///sys
	orm.RegisterModel(new(casbinUtil.CasbinRule))
	orm.RegisterModel(new(models.RoleModel))
	orm.RegisterModel(new(models.UserModel))
	orm.RegisterModel(new(models.ChildUserModel))
	orm.RegisterModel(new(models.DomainModel))
	//orm.RegisterModel(new(models.PolicyModel))
    orm.RegisterModel(new(models.SmsModel))
	orm.RegisterModel(new(models.UserRoleModel))
	orm.RegisterModel(new(models.RouterModel))

	/////
	orm.RegisterModel(new(models.AbilityCategoryModel))
	orm.RegisterModel(new(models.AbilityModel))
	orm.RegisterModel(new(models.ApiModel))
	orm.RegisterModel(new(models.DefaultQuotasModel))
	orm.RegisterModel(new(models.AppModel))
	orm.RegisterModel(new(models.AppAbilityModel))
	orm.RegisterModel(new(models.QuotasModel))
	orm.RegisterModel(new(models.LogHourlyQuotasStatisticModel))
	orm.RegisterModel(new(models.LogDailyQuotasStatisticModel))
	orm.RegisterModel(new(models.LogAppQuotasStatisticModel))
	orm.RegisterModel(new(models.LogApiQuotasStatisticModel))
	orm.RegisterModel(new(models.LogAccessDetail))
	orm.RegisterModel(new(models.LogAppDailyStatisticModel))
	orm.RegisterModel(new(models.LogApiDailyStatistic))
	orm.RegisterModel(new(models.LogApiDailyHistoryAddUp))



	_ = orm.RunSyncdb("default", false, false)
}

func initFunc() {
	casbinUtil.RegisterCasbin()
	models.InitElastic()
	services.InitCaptcha()
	models.InitCeph()

}
