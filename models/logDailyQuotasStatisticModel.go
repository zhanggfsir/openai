package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type LogDailyQuotasStatisticModelTemp struct {
	Id             uint64 `orm:"auto" description:"ID" json:"id"`
	AppId          uint64 `description:"appId" json:"app_id"`
	ApiId          uint64 `description:"apiId" json:"api_id"`
	AppName        string `description:"month" json:"app_name"`
	ApiName        string `description:"month" json:"api_name"`
	Year           string `description:"month" json:"year"`
	YearSuccessful uint64 `description:"month" json:"year_successful"`
	YearFailed     uint64 `description:"month" json:"year_failed"`
	Month          string `description:"month" json:"month"`
	MonthCount     uint64 `description:"month" json:"month_count"`
	Day            string `description:"day" json:"day"`
	DayCount       uint64 `description:"day" json:"day_count"`
	FailedCalled   uint64 `description:"失败调用次数" json:"failed_called"`
	Domain         string `description:"域名" json:"domain"`

	CreateTime string `json:"createdTime"`
	ModifyTime string `json:"updatedTime"`
}

type LogDailyQuotasStatisticModel struct {
	Id              uint64 `orm:"auto" description:"ID" json:"id"`
	AppId           uint64 `description:"appId" json:"app_id"`
	ApiId           uint64 `description:"apiId" json:"api_id"`
	AppName         string `description:"month" json:"app_name"`
	ApiName         string `description:"month" json:"api_name"`
	Year            string `description:"month" json:"year"`
	YearSuccessful  uint64 `description:"month" json:"year_successful"`
	YearFailed      uint64 `description:"month" json:"year_failed"`
	Month           string `description:"month" json:"month"`
	MonthSuccessful uint64 `description:"month" json:"month_successful"`
	MonthFailed     uint64 `description:"month" json:"month_failed"`
	Day             string `description:"day" json:"day"`
	DaySuccessful   uint64 `description:"day" json:"day_successful"`
	DayFailed       uint64 `description:"day" json:"day_failed"`
	Domain          string `description:"域名" json:"domain"`

	PublicModel
}

func (this *LogDailyQuotasStatisticModel) TableName() string {
	return "log_daily_quotas_statistic"
}

// TableEngine 获取数据使用的引擎.
func (this *LogDailyQuotasStatisticModel) TableEngine() string {
	return "INNODB"
}

func (this *LogDailyQuotasStatisticModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewLogDailyQuotasStatisticModel() *LogDailyQuotasStatisticModel {
	return &LogDailyQuotasStatisticModel{}
}

func (this *LogDailyQuotasStatisticModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix()) //.Filter("domain", domain)
}

//添加设备
func (this *LogDailyQuotasStatisticModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//修改设备
func (this *LogDailyQuotasStatisticModel) Update(id string, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *LogDailyQuotasStatisticModel) UpdateById(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

//多个过滤条件
func (this *LogDailyQuotasStatisticModel) filter(domain string, queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}

func (this *LogDailyQuotasStatisticModel) GetLogDailyQuotasStatistic(startTime, endTime string) (list []*LogDailyQuotasStatisticModel, err error) {
	o := orm.NewOrm()
	var values = make([]interface{}, 0)

	values = append(values, startTime)
	values = append(values, endTime)

	//logs.Info("---- this.TableName()-->", this.TableName())
	//  ApiIdName       this.TableName()  //  openai_detail_log_20201221
	sql := `
SELECT  app_id,api_id ,app_name, api_name ,
year,successful_called as year_successful ,failed_called as year_failed,
month,successful_called as month_successful ,failed_called as month_failed,
DAY,successful_called as day_successful ,failed_called as day_failed,
hour,successful_called as hour_successful ,failed_called as hour_called
FROM (
SELECT app_id,api_id ,app_name, api_name ,year,month,day,hour,SUM(successful) as successful_called ,SUM(failed) as failed_called
FROM (
SELECT app_id,api_id ,app_name, api_name ,year,month,day,HOUR,
case when called_status ="200" then 1 ELSE 0 END  AS successful,
case when called_status!="200"  then 1 ELSE 0 END AS failed
FROM ` + this.TableName() + ` where created_time>= ? and created_time< ?
) t  GROUP BY app_id,api_id ,app_name, api_name ,year,month,day,HOUR
)tt `
	//logs.Info("------2------->",sql,values)

	_, err = o.Raw(sql, values).QueryRows(&list)

	if err != nil {
		logs.Error(err)
		return nil, err
	}
	//for i, i2 := range list {
	//	logs.Info(i,"------>",i2)
	//}
	//logs.Info(list)
	return

}

func (this *LogDailyQuotasStatisticModel) UpdateByIdAndDay(apiId, appId uint64, day string, domain string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id", apiId).Filter("app_id", appId).Filter("day", day).Filter("domain", domain).Update(data)

	return err
}

type ApiCalledByAppHistogram struct {
	AppId           uint64 `description:"appId" json:"app_id"`
	AppName         string `description:"month" json:"app_name"`
	Successful   	uint64 `description:"month" json:"successful"`
	Failed      	uint64 `description:"month" json:"failed"`
	Sum 			uint64 `description:"month" json:"sum"`

}

/*
select app_id,app_name,sum(day_successful) AS successful , sum(day_failed) AS failed,  sum(day_successful)+sum(day_failed) AS  SUM
from log_daily_quotas_statistic where api_id=30  and domain ="default" group by app_id,app_name ORDER BY SUM DESC LIMIT 10
 */
func (this *LogDailyQuotasStatisticModel) GetApiCalledByAppHistogramList(q map[string]interface{}, domain string) (histogram []*ApiCalledByAppHistogram, err error) {
	//calledSuccess = make([]*uint64, 0)
	o := orm.NewOrm()

	var values = make([]interface{}, 0)

	if value, ok := q["api_id"]; ok {
		values = append(values, value)
	}

	values = append(values, domain)

	if value, ok := q["top"]; ok {
		values = append(values, value)
	}
	//if value, ok := q["start_time"]; ok {
	//	values = append(values, value)
	//}
	//if value, ok := q["end_time"]; ok {
	//	values = append(values, value)
	//}

	logs.Info(values)
	sql := `
select app_id,app_name,sum(day_successful) AS successful , sum(day_failed) AS failed,  sum(day_successful)+sum(day_failed) AS  sum 
from log_daily_quotas_statistic where api_id = ?  and domain = ? group by app_id,app_name ORDER BY SUM DESC LIMIT ? `

	_, err = o.Raw(sql, values).QueryRows(&histogram)

	return

}



type AppCallApiLineChart struct {
	ApiId           uint64 `description:"api_id" json:"api_id"`
	ApiName         string `description:"api_name" json:"api_name"`
	Day             string `description:"day" json:"day"`
	Successful 		uint64 `description:"successful" json:"successful"`
	Failed     		uint64 `description:"failed" json:"failed"`
	Sum 			uint64 `description:"sum" json:"sum"`

}

/*
select api_id,api_name,day,sum(day_successful) AS successful , sum(day_failed) AS failed,  sum(day_successful)+sum(day_failed) AS  SUM
from log_daily_quotas_statistic where app_id=119  and domain ="default" group by api_id,api_name,day

ORDER BY SUM DESC LIMIT 10

 */
func (this *LogDailyQuotasStatisticModel) GetAppCallApiLineChartList(q map[string]interface{}, domain string) (appCallApiLineChart []*AppCallApiLineChart, err error) {
	//calledSuccess = make([]*uint64, 0)
	o := orm.NewOrm()

	var values = make([]interface{}, 0)

	if value, ok := q["app_id"]; ok {
		values = append(values, value)
	}

	values = append(values, domain)

	//if value, ok := q["top"]; ok {
	//	values = append(values, value)
	//}
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}

	sql := `
select api_id,api_name,day,sum(day_successful) AS successful , sum(day_failed) AS failed,  sum(day_successful)+sum(day_failed) AS  sum 
from log_daily_quotas_statistic where app_id= ?  and domain = ? and day >= ? and day <= ? and api_id!=0  group by api_id,api_name,day 
 `
	//
	_, err = o.Raw(sql, values).QueryRows(&appCallApiLineChart)

	//for i, i2 := range appCallApiLineChart {
	//	logs.Info(i,i2)
	//}

	return

}




