package services

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"openai-backend/models"
	"strconv"
	"time"
)

func InsertAccessLogKafka(appId, apiId, status, ts uint64, url, appName,apiName ,domain string) error {
	if apiId==0{
		logs.Info("--- api is zero ---",appId,appName, apiId,apiName, status, ts, "url-->",url)
	}
	// ts
	timeTemplateYear := "2006"                //其他类型	 月
	timeTemplateMonth := "2006-01"            //其他类型	 月
	timeTemplateDay := "2006-01-02"           //其他类型	 月
	timeTemplateHour := "2006-01-02 15:00:00" //其他类型	 月
	timeTemplateMinute:= "2006-01-02 15:04:05"
	year := time.Unix(int64(ts), 0).Format(timeTemplateYear)
	month := time.Unix(int64(ts), 0).Format(timeTemplateMonth)
	day := time.Unix(int64(ts), 0).Format(timeTemplateDay)
	hour := time.Unix(int64(ts), 0).Format(timeTemplateHour)
	minute := time.Unix(int64(ts), 0).Format(timeTemplateMinute)

	createdTime, _ := time.ParseInLocation(timeTemplateMinute, minute, time.Local)

	// 1.原始日志写入 log_access_detail，后期继续写入 es
	//defaultQuotas
	accessLog := &models.LogAccessDetail{
		AppId:        appId,
		ApiId:        apiId,
		AppName:      appName,
		ApiName:      apiName,
		CalledStatus: strconv.FormatUint(status, 10),
		Year:         year,
		Month:        month,
		Day:          day,
		Hour:         hour,
		Url:          url,
		Domain: 	  domain,
		CreatedTime:   createdTime,
		UpdateTime:		time.Now(),

	}

	values:=make( []interface{},0)
	values=append(values,  appId)
	values=append(values,  apiId)
	values=append(values,  appName)
	values=append(values,  apiName)
	values=append(values,  year)
	values=append(values,  month)
	values=append(values,  day)
	values=append(values,  hour)
	values=append(values,  strconv.FormatUint(status, 10))
	values=append(values,  url)
	values=append(values,	domain)
	values=append(values,   createdTime)
	values=append(values,   time.Now())
	//data["updated_time"] = time.Now()


	//插入es LogAccessDetail
	if err := openAiDetailLogSerializationFromKafka(accessLog,values); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}



func openAiDetailLogSerializationFromKafka(accessLog *models.LogAccessDetail,values []interface{}) error {
	//插入es logAccessDetail
	isOpenElastic := beego.AppConfig.DefaultBool("elastic", false)
	if isOpenElastic {
		err := accessLog.InsertElasticFromKafka(values)
		return err
	} else {
		err := accessLog.InsertMySql(values)
		//err := accessLog.Insert(accessLog,values)
		if err != nil {
			logs.Info(err)
			return err

		}
		return nil
	}

}







func DataStatistics(startTime, endTime string) error {
	//// todo
	//domain := "domain"
	//  log_hourly_quotas_statistic LogHourlyQuotasStatistic WindowStatisticModel
	windowStatisticList, err := models.NewLogAccessDetail().GetWindowStatistic(startTime, endTime)
	for _, windowStatistic := range windowStatisticList {
		//  1.  log_hourly_quotas_statistic LogHourlyQuotasStatistic
		logHourlyQuotasStatisticModel := &models.LogHourlyQuotasStatisticModel{}

		if models.NewHourlyQuotasStatisticModel().QueryTable().Filter("api_id", windowStatistic.ApiId).Filter("app_id", windowStatistic.AppId).Filter("hour", windowStatistic.Hour).
			Filter("domain",windowStatistic.Domain).Exist() != true {

			logHourlyQuotasStatisticModel.AppId=windowStatistic.AppId
			logHourlyQuotasStatisticModel.ApiId=windowStatistic.ApiId
			logHourlyQuotasStatisticModel.AppName=windowStatistic.AppName
			logHourlyQuotasStatisticModel.ApiName=windowStatistic.ApiName
			logHourlyQuotasStatisticModel.Year=windowStatistic.Year
			logHourlyQuotasStatisticModel.YearSuccessful=windowStatistic.YearSuccessful
			logHourlyQuotasStatisticModel.YearFailed=windowStatistic.YearFailed
			logHourlyQuotasStatisticModel.Month=windowStatistic.Month
			logHourlyQuotasStatisticModel.MonthSuccessful=windowStatistic.MonthSuccessful
			logHourlyQuotasStatisticModel.MonthFailed=windowStatistic.MonthFailed
			logHourlyQuotasStatisticModel.Day=windowStatistic.Day
			logHourlyQuotasStatisticModel.DaySuccessful=windowStatistic.DaySuccessful
			logHourlyQuotasStatisticModel.DayFailed=windowStatistic.DayFailed
			logHourlyQuotasStatisticModel.Hour=windowStatistic.Hour
			logHourlyQuotasStatisticModel.HourSuccessful=windowStatistic.HourSuccessful
			logHourlyQuotasStatisticModel.HourFailed=windowStatistic.HourFailed

			logHourlyQuotasStatisticModel.Domain=windowStatistic.Domain


			if err := logHourlyQuotasStatisticModel.Insert(); err != nil {
				logs.Error("err= ", err)
				return err
			}
		} else {

			// mysql exist. get mysql exist
			err = models.NewHourlyQuotasStatisticModel().QueryTable().Filter("api_id", windowStatistic.ApiId).Filter("app_id", windowStatistic.AppId).Filter("hour", windowStatistic.Hour).
				Filter("domain",windowStatistic.Domain).One(logHourlyQuotasStatisticModel)
			if err != nil {
				logs.Error(err)
			}
			data := make(map[string]interface{})
			data["year_successful"] = logHourlyQuotasStatisticModel.YearSuccessful + windowStatistic.YearSuccessful
			data["year_failed"] = logHourlyQuotasStatisticModel.YearFailed + windowStatistic.YearFailed

			data["month_successful"] = logHourlyQuotasStatisticModel.MonthSuccessful + windowStatistic.MonthSuccessful
			data["month_failed"] = logHourlyQuotasStatisticModel.MonthFailed + windowStatistic.MonthFailed

			data["day_successful"] = logHourlyQuotasStatisticModel.DaySuccessful + windowStatistic.DaySuccessful
			data["day_failed"] = logHourlyQuotasStatisticModel.DayFailed + windowStatistic.DayFailed

			data["hour_successful"] = logHourlyQuotasStatisticModel.HourSuccessful + windowStatistic.HourSuccessful
			data["hour_failed"] = logHourlyQuotasStatisticModel.HourFailed + windowStatistic.HourFailed

			err = logHourlyQuotasStatisticModel.UpdateByIdAndHour(logHourlyQuotasStatisticModel.ApiId, logHourlyQuotasStatisticModel.AppId, logHourlyQuotasStatisticModel.Hour, windowStatistic.Domain,data)
			if err != nil {
				logs.Error(err)
			}
		}

		// 2.  log_daily_quotas_statistic logDailyQuotasStatistic
		daily := &models.LogDailyQuotasStatisticModel{}
		if models.NewLogDailyQuotasStatisticModel().QueryTable().Filter("api_id", windowStatistic.ApiId).Filter("app_id", windowStatistic.AppId).Filter("day", windowStatistic.Day).
			Filter("domain",windowStatistic.Domain).Exist() != true {
			daily.AppId = windowStatistic.AppId
			daily.ApiId = windowStatistic.ApiId
			daily.AppName = windowStatistic.AppName
			daily.ApiName = windowStatistic.ApiName
			daily.Year = windowStatistic.Year
			daily.YearSuccessful = windowStatistic.YearSuccessful
			daily.YearFailed = windowStatistic.YearFailed
			daily.Month = windowStatistic.Month
			daily.MonthSuccessful = windowStatistic.MonthSuccessful
			daily.MonthFailed = windowStatistic.MonthFailed
			daily.Day = windowStatistic.Day
			daily.DaySuccessful = windowStatistic.DaySuccessful
			daily.DayFailed = windowStatistic.DayFailed
			daily.Domain = windowStatistic.Domain


			if err := daily.Insert(); err != nil {
				logs.Error("err= ", err)
				return err

			}
		} else {
			// mysql exist. get mysql exist
			err = models.NewLogDailyQuotasStatisticModel().QueryTable().Filter("api_id", windowStatistic.ApiId).Filter("app_id", windowStatistic.AppId).Filter("day", windowStatistic.Day).
				Filter("domain",windowStatistic.Domain).One(daily)
			if err != nil {
				logs.Error(err)
			}
			data := make(map[string]interface{})
			data["year_successful"] = daily.YearSuccessful + windowStatistic.YearSuccessful
			data["year_failed"] = daily.YearFailed + windowStatistic.YearFailed
			data["month_successful"] = daily.MonthSuccessful + windowStatistic.MonthSuccessful
			data["month_failed"] = daily.MonthFailed + windowStatistic.MonthFailed
			data["day_successful"] = daily.DaySuccessful + windowStatistic.DaySuccessful
			data["day_failed"] = daily.DayFailed + windowStatistic.DayFailed

			//logs.Info("-------data----------", data, "api_id-->", myLogHourlyQuotasStatisticModel)
			err = daily.UpdateByIdAndDay(daily.ApiId, daily.AppId, daily.Day, windowStatistic.Domain, data)
			if err != nil {
				logs.Error(err)
			}

		}

		//  3 .log_app_daily_statistic LogAppDailyStatistic	LogAppDailyStatisticModel
		logAppDaily:=&models.LogAppDailyStatisticModel{}

		if logAppDaily.QueryTable().Filter("app_id", windowStatistic.AppId).Filter("day", windowStatistic.Day).Filter("domain",windowStatistic.Domain).Exist() != true {
			logAppDaily.AppId = windowStatistic.AppId
			logAppDaily.AppName = windowStatistic.AppName
			logAppDaily.Year = windowStatistic.Year
			logAppDaily.YearSuccessful = windowStatistic.YearSuccessful
			logAppDaily.YearFailed = windowStatistic.YearFailed
			logAppDaily.Month = windowStatistic.Month
			logAppDaily.MonthSuccessful = windowStatistic.MonthSuccessful
			logAppDaily.MonthFailed = windowStatistic.MonthFailed
			logAppDaily.Day = windowStatistic.Day
			logAppDaily.DaySuccessful = windowStatistic.DaySuccessful
			logAppDaily.DayFailed = windowStatistic.DayFailed

			logAppDaily.Domain = windowStatistic.Domain

			if err := logAppDaily.Insert(); err != nil {
				logs.Error("err= ", err)
				return err
			}

		}else{

			err = models.NewLogAppDailyStatisticModel().QueryTable().Filter("app_id", windowStatistic.AppId).Filter("day", windowStatistic.Day).
				Filter("domain",windowStatistic.Domain).One(logAppDaily)
			if err != nil {
				logs.Error(err)
			}

			data := make(map[string]interface{})

			data["year_successful"] = logAppDaily.YearSuccessful + windowStatistic.YearSuccessful
			data["year_failed"] = logAppDaily.YearFailed + windowStatistic.YearFailed
			data["month_successful"] = logAppDaily.MonthSuccessful + windowStatistic.MonthSuccessful
			data["month_failed"] = logAppDaily.MonthFailed + windowStatistic.MonthFailed
			data["day_successful"] = logAppDaily.DaySuccessful + windowStatistic.DaySuccessful
			data["day_failed"] = logAppDaily.DayFailed + windowStatistic.DayFailed

			//appId uint64,day string, data map[string]interface{}, domain string
			err = logAppDaily.UpdateByIdAndDay(logAppDaily.AppId, logAppDaily.Day, windowStatistic.Domain, data)
			if err != nil {
				logs.Error(err)
			}
		}



		//   4.log_api_daily_statistic log_api_daily_statistic	LogAppDailyStatisticModel
		logApiDaily:=&models.LogApiDailyStatistic{}
		if logApiDaily.QueryTable().Filter("api_id", windowStatistic.ApiId).Filter("day", windowStatistic.Day).Filter("domain",windowStatistic.Domain).Exist() != true {
			logApiDaily.ApiId = windowStatistic.ApiId
			logApiDaily.ApiName = windowStatistic.ApiName
			logApiDaily.Year = windowStatistic.Year
			logApiDaily.YearSuccessful = windowStatistic.YearSuccessful
			logApiDaily.YearFailed = windowStatistic.YearFailed
			logApiDaily.Month = windowStatistic.Month
			logApiDaily.MonthSuccessful = windowStatistic.MonthSuccessful
			logApiDaily.MonthFailed = windowStatistic.MonthFailed
			logApiDaily.Day = windowStatistic.Day
			logApiDaily.DaySuccessful = windowStatistic.DaySuccessful
			logApiDaily.DayFailed = windowStatistic.DayFailed

			logApiDaily.Domain = windowStatistic.Domain

			if err := logApiDaily.Insert(); err != nil {
				logs.Error("err= ", err)
				return err
			}
		}else{
			err = models.NewLogApiDailyStatistic().QueryTable().Filter("api_id", windowStatistic.ApiId).Filter("day", windowStatistic.Day).
				Filter("domain",windowStatistic.Domain).One(logApiDaily)
			if err != nil {
				logs.Error(err)
			}

			data := make(map[string]interface{})

			data["year_successful"] = logApiDaily.YearSuccessful + windowStatistic.YearSuccessful
			data["year_failed"] = logApiDaily.YearFailed + windowStatistic.YearFailed
			data["month_successful"] = logApiDaily.MonthSuccessful + windowStatistic.MonthSuccessful
			data["month_failed"] = logApiDaily.MonthFailed + windowStatistic.MonthFailed
			data["day_successful"] = logApiDaily.DaySuccessful + windowStatistic.DaySuccessful
			data["day_failed"] = logApiDaily.DayFailed + windowStatistic.DayFailed

			err = logApiDaily.UpdateByIdAndDay(logApiDaily.ApiId, logApiDaily.Day, windowStatistic.Domain, data)
			if err != nil {
				logs.Error(err)
			}

		}
		//  5.log_app_quotas_statistic
		logApp:=&models.LogAppQuotasStatisticModel{}

		if logApp.QueryTable().Filter("app_id", windowStatistic.AppId).Filter("domain",windowStatistic.Domain).Exist() != true {
			logApp.AppId = windowStatistic.AppId
			logApp.AppName = windowStatistic.AppName
			logApp.SuccessfulCalled = windowStatistic.HourSuccessful
			logApp.FailedCalled = windowStatistic.HourFailed

			logApp.Domain = windowStatistic.Domain

			if err := logApp.Insert(); err != nil {
				logs.Error("err= ", err)
				return err
			}
		}else{
			err = models.NewLogAppQuotasStatisticModel().QueryTable().Filter("app_id", windowStatistic.AppId).Filter("domain",windowStatistic.Domain).One(logApp)
			if err != nil {
				logs.Error(err)
			}

			data := make(map[string]interface{})

			data["successful_called"] = logApp.SuccessfulCalled + windowStatistic.HourSuccessful
			data["failed_called"] = logApp.FailedCalled + windowStatistic.HourFailed

			//appId uint64,day string, data map[string]interface{}, domain string
			err = logApp.UpdateByIdAndDay(logApp.AppId, windowStatistic.Domain,  data)
			if err != nil {
				logs.Error(err)
			}

		}

		//  6. log_api_quotas_statistic
		logApi:=&models.LogApiQuotasStatisticModel{}
		if logApi.QueryTable().Filter("api_id", windowStatistic.ApiId).Filter("domain",windowStatistic.Domain).Exist() != true {

			logApi.ApiId = windowStatistic.ApiId
			logApi.ApiName = windowStatistic.ApiName
			logApi.SuccessfulCalled = windowStatistic.HourSuccessful
			logApi.FailedCalled = windowStatistic.HourFailed
			logApi.Domain = windowStatistic.Domain
			if err := logApi.Insert(); err != nil {
				logs.Error("err= ", err)
				return err
			}
		}else{
			err = models.NewLogApiQuotasStatisticModel().QueryTable().Filter("api_id", windowStatistic.ApiId).Filter("domain",windowStatistic.Domain).One(logApi)
			if err != nil {
				logs.Error(err)
			}

			data := make(map[string]interface{})

			data["successful_called"] = logApi.SuccessfulCalled + windowStatistic.HourSuccessful
			data["failed_called"] = logApi.FailedCalled + windowStatistic.HourFailed

			//appId uint64,day string, data map[string]interface{}, domain string
			err = logApi.UpdateByIdAndDay(logApi.ApiId, windowStatistic.Domain,  data)
			if err != nil {
				logs.Error(err)
			}
		}



	}



	/*
		-------  只需调用一次即可  区别 用hour更新 ---------------
		1. 计算出昨天[一天]的总数 A
		2. 取前天所有 B   A+B即昨天全量
	*/

	//   7.log_api_daily_statistic log_api_daily_statistic	LogAppDailyStatisticModel
	if startTime[0:10]!=endTime[0:10] {
		currentTime := time.Now()
		// 昨天
		yesterday := currentTime.AddDate(0, 0, -1).Format("2006-01-02")
		apiDailyStatistic:=&models.LogApiDailyStatistic{}
		logApiYesterday,_:=apiDailyStatistic.GetLogApiYesterdayTotal(yesterday)

		// 前天
		beforeYesterday := currentTime.AddDate(0, 0, -2).Format("2006-01-02")

		addUpBY :=&models.LogApiDailyHistoryAddUp{}
		err=addUpBY.QueryTable().Filter("day",beforeYesterday).One(addUpBY)
		if err!=nil{
			logs.Error(err)
		}
		logs.Info(addUpBY.DaySuccessful,"----addUpBY----->",addUpBY,"   ",addUpBY.DailyHistoryAddUp,addUpBY.DailySuccessfulHistoryAddUp,beforeYesterday)

		// addUpYesterday.DailyHistoryAddUp=logApiYesterday.DaySuccessful+logApiYesterday.DayFailed+
		// 累计
		// 今日数据=今日数据
		addUp:=&models.LogApiDailyHistoryAddUp{}
		addUp.Day=yesterday
		addUp.DayTotal=logApiYesterday.DaySuccessful+logApiYesterday.DayFailed
		addUp.DaySuccessful=logApiYesterday.DaySuccessful
		addUp.DayFailed=logApiYesterday.DayFailed
		// 今日统计=今日数据+昨日统计
		addUp.DailyHistoryAddUp=logApiYesterday.DaySuccessful+logApiYesterday.DayFailed+addUpBY.DailyHistoryAddUp
		addUp.DailySuccessfulHistoryAddUp=logApiYesterday.DaySuccessful+addUpBY.DailySuccessfulHistoryAddUp
		addUp.DailyFailedHistoryAddUp=logApiYesterday.DayFailed+addUpBY.DailyFailedHistoryAddUp

		addUp.Insert()

	}

	return nil

}


