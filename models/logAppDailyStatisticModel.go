package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type LogAppDailyStatisticModelTemp struct {

	AppId           uint64 `description:"appId" json:"app_id"`
	AppName           string `description:"apiId" json:"app_name"`

	Year            string `description:"month" json:"year"`
	YearSuccessful       uint64 `description:"month" json:"year_successful"`
	YearFailed       uint64 `description:"month" json:"year_failed"`
	Month           string `description:"month" json:"month"`
	MonthSuccessful uint64 `description:"month" json:"month_successful"`
	MonthFailed     uint64 `description:"month" json:"month_failed"`
	Day             string `description:"day" json:"day"`
	DaySuccessfulAndFailed  uint64 `description:"day" json:"day_successful_failed"`
	DaySuccessful   uint64 `description:"day" json:"day_successful"`
	DayFailed       uint64 `description:"day" json:"day_failed"`
	DayFailedRate   float64 `description:"day" json:"day_failed_rate"`
	Domain      string  `description:"域名" json:"domain"`



	//Domain     string `description:"所属域" json:"domain"`
	//CreateTime string `json:"createdTime"`
	//ModifyTime string `json:"updatedTime"`
}

type LogAppDailyStatisticModel struct {

	Id          	uint64 `orm:"auto" description:"ID" json:"id"`
	AppId           uint64 ` description:"appId" json:"app_id"`
	AppName         string `description:"apiId" json:"app_name"`

	Year            string `description:"month" json:"year"`
	YearSuccessful  uint64 `description:"month" json:"year_successful"`
	YearFailed      uint64 `description:"month" json:"year_failed"`

	Month           string `description:"month" json:"month"`
	MonthSuccessful uint64 `description:"month" json:"month_successful"`
	MonthFailed     uint64 `description:"month" json:"month_failed"`

	Day             string `description:"day" json:"day"`
	DaySuccessful   uint64 `description:"day" json:"day_successful"`
	DayFailed       uint64 `description:"day" json:"day_failed"`
	Domain      string  `description:"域名" json:"domain"`


	PublicModel
}

func (this *LogAppDailyStatisticModel) TableName() string {
	return "log_app_daily_statistic"
}

// TableEngine 获取数据使用的引擎.
func (this *LogAppDailyStatisticModel) TableEngine() string {
	return "INNODB"
}

func (this *LogAppDailyStatisticModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewLogAppDailyStatisticModel() *LogAppDailyStatisticModel {
	return &LogAppDailyStatisticModel{}
}

func (this *LogAppDailyStatisticModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix()) //.Filter("domain", domain)
}

//添加设备
func (this *LogAppDailyStatisticModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//修改设备
func (this *LogAppDailyStatisticModel) Update(id string, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *LogAppDailyStatisticModel) UpdateById(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

//多个过滤条件
func (this *LogAppDailyStatisticModel) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}




/*
	SELECT api_id,api_name,SUM(day_successful) AS successful,SUM(day_failed) AS failed
	FROM log_api_daily_statistic WHERE category_id=1 AND created_time>="20201207"  AND created_time<="20201218"  GROUP BY api_id,api_name ;
*/
func (this *LogAppDailyStatisticModel) GetLogAppDailyModel(
	q map[string]interface{},domain string) (list []*LogAppDailyViewTemp, Count int, err error) {

	o := orm.NewOrm()
	var values = make([]interface{}, 0)

	sql := `SELECT app_id,app_name,SUM(day_successful) AS successful,SUM(day_failed) AS failed FROM log_app_daily_statistic WHERE 1=1 and domain = ? `

	values=append(values,domain)


	if q["start_time"]!=""{
		curTimeSuffix:=time.Now().Format(" 15:04:05" )
		values = append(values, q["start_time"].(string)+curTimeSuffix)
		sql = sql + " and created_time >= ? "
	}
	if q["end_time"]!=""{

		curTimeSuffix:=time.Now().Format(" 15:04:05" )
		//logs.Info(q["end_time"],"---","---end_time-->",q["end_time"].(string)+curTimeSuffix)
		values = append(values, q["end_time"].(string)+curTimeSuffix)
		sql = sql + " and created_time <= ? "
	}



	//if value, ok := q["start_time"]; ok {
	//	values = append(values, value)
	//	sql = sql + " and created_time >= ? "
	//}
	//if value, ok := q["end_time"]; ok {
	//	values = append(values, value)
	//	sql = sql + " and created_time <= ? "
	//}
	sql=sql+" GROUP BY app_id,app_name"
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	//sql := `SELECT * FROM  log_quotas_statistic WHERE created_time>=? AND end_time<=?`

	//logs.Info(sql)
	_, err = o.Raw(sql, values).QueryRows(&list)

	//logs.Info(list)
	return

}


type LogAppDailyViewTemp struct {
	AppId           uint64 `description:"apiId" json:"app_id"`
	AppName         string `description:"month" json:"app_name"`
	Total 			uint64 `description:"total" 	json:"total"`
	Successful  	 uint64 `description:"successful" json:"successful"`
	Failed      	 uint64 `description:"failed" json:"failed"`
	PublicModel
}

type LogAppDailyView struct {
	AppId           uint64 `description:"apiId" json:"app_id"`
	AppName         string `description:"month" json:"app_name"`
	Total 			uint64 `description:"total" 	json:"total"`
	Successful  	 uint64 `description:"successful" json:"successful"`
	Failed      	 uint64 `description:"failed" json:"failed"`
	FailedRate      float64 `description:"day" json:"failedRate"`

	PublicModel
}






func (this *LogAppDailyStatisticModel) GetLogAppDailyModel_(
	q map[string]interface{}) (list []*LogAppDailyStatisticModel, Count int, err error) {

	o := orm.NewOrm()
	var values = make([]interface{}, 0)
	// todo sql 拼接
	sql := `SELECT * from log_app_daily_statistic WHERE 1=1  `
	// sql := `SELECT * FROM  log_quotas_statistic WHERE 1=1   `
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
		sql = sql + " and created_time >= ? "
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
		sql = sql + " and created_time <= ? "
	}

	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	//sql := `SELECT * FROM  log_quotas_statistic WHERE created_time>=? AND end_time<=?`

	_, err = o.Raw(sql, values).QueryRows(&list)

	return

	////////////////////////////////////////////////////////////////////
	//var c int64
	//
	//c, err = this.filter(queryFilter).Count()
	//
	////查询sql语句
	//_, err = this.filter(queryFilter).All(&list)
	//if err != nil {
	//	return
	//}
	//Count = int(c)
	//return
}




func (this *LogAppDailyStatisticModel) UpdateByIdAndDay(appId uint64,day string, domain string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("app_id",appId).Filter("day",day).Filter("domain",domain).Update(data)

	return err
}




