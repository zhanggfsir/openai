package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type LogApiDailyStatisticTemp struct {

	ApiId           uint64 `description:"apiId" json:"api_id"`
	ApiName         string `description:"month" json:"api_name"`
	//CategoryId      uint64 `description:"apiId" json:"category_id"`
	//CategoryName    string `description:"month" json:"category_name"`
	//AbilityId       uint64 `description:"apiId" json:"ability_id"`
	//AbilityName     string `description:"month" json:"ability_name"`

	Year            string `description:"month" json:"year"`
	YearSuccessful  uint64 `description:"month" json:"year_successful"`
	YearFailed      uint64 `description:"month" json:"year_failed"`
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

type LogApiDailyStatistic struct {
	Id          	uint64 `orm:"auto" description:"ID" json:"id"`
	ApiId           uint64 `description:"apiId" json:"api_id"`
	ApiName         string `description:"month" json:"api_name"`
	//AbilityId       uint64 `description:"apiId" json:"ability_id"`
	//AbilityName     string `description:"month" json:"ability_name"`
	//CategoryId      uint64 `description:"apiId" json:"category_id"`
	//CategoryName    string `description:"month" json:"category_name"`


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

func (this *LogApiDailyStatistic) TableName() string {
	return "log_api_daily_statistic"
}

// TableEngine 获取数据使用的引擎.
func (this *LogApiDailyStatistic) TableEngine() string {
	return "INNODB"
}

func (this *LogApiDailyStatistic) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewLogApiDailyStatistic() *LogApiDailyStatistic {
	return &LogApiDailyStatistic{}
}

func (this *LogApiDailyStatistic) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix()) //.Filter("domain", domain)
}

//添加设备
func (this *LogApiDailyStatistic) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//修改设备
func (this *LogApiDailyStatistic) Update(id string, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *LogApiDailyStatistic) UpdateById(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

//多个过滤条件
func (this *LogApiDailyStatistic) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
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
由于加了时间顾虑，需要从 log_api_daily_statistic统计数据。

SELECT * FROM log_api_quotas_statistic WHERE api_id
IN (
select t2.api_id
from
(select  id as ability_id, category_id from ability_info  where category_id=6 ) t1
inner join
( select id as api_id,ability_id from api_info)t2
ON  t1.ability_id=t2.ability_id
) ;

当前统计逻辑
SELECT api_id,api_name,SUM(day_successful) AS successful,SUM(day_failed) AS failed FROM log_api_daily_statistic
WHERE 1 =1 AND api_id IN() GROUP BY api_id,api_name

SELECT api_id,api_name,SUM(day_successful) AS successful,SUM(day_failed) AS failed FROM log_api_daily_statistic
WHERE 1 =1 AND api_id IN(
	select t2.api_id
	from
	(select  id as ability_id, category_id from ability_info  where category_id=6 ) t1
	inner join
	( select id as api_id,ability_id from api_info)t2
	ON  t1.ability_id=t2.ability_id
) AND created_time>="2020-12-20 10:36:22" AND created_time<="2020-12-24 10:36:22" GROUP BY api_id,api_name



//	mysql:=`
//SELECT api_id,api_name,SUM(day_successful)  successful,SUM(day_failed)  failed
//FROM log_api_daily_statistic WHERE 1 =1
//AND api_id IN ( SELECT  t2.id
//FROM
//ability_info t1
//INNER JOIN api_info t2
//ON  t1.id=t2.ability_id  WHERE  t1.category_id = 1  )
//AND DAY <= "2020-12-31"
//GROUP BY api_id,api_name`
 */





func (this *LogApiDailyStatistic) GetLogApiDailyModel(categoryId uint64,
	q map[string]interface{} ,domain string) (list []*LogApiDailyViewTemp, Count int, err error) {

	o := orm.NewOrm()
	var values = make([]interface{}, 0)
	//logs.Info("--categoryId-->",categoryId)
	sql := `
	SELECT api_id,api_name,SUM(day_successful) AS successful,SUM(day_failed) AS failed FROM log_api_daily_statistic 
	WHERE 1 =1  and domain = ? AND api_id IN(
	select t2.api_id
	from
	(select  id as ability_id, category_id from ability_info  WHERE 1=1 `

	values = append(values, domain)

	if categoryId != 0 {	// all
		values = append(values, categoryId)
		sql = sql + ` and category_id = ? `
	}
	sql = sql+ ` ) t1 inner join ( select id as api_id,ability_id from api_info ) t2 ON  t1.ability_id=t2.ability_id ) `
	if q["start_time"]!=""{

		startTime:=q["start_time"].(string)[:10]
		values = append(values, startTime)

		sql = sql + ` and day >= ? `
	}

	if q["end_time"]!=""{

		endTime:=q["end_time"].(string)[:10]
		values = append(values, endTime)

		sql = sql + ` and day <= ? `
	}

	sql=sql+` GROUP BY api_id,api_name `
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值

	_, err = o.Raw(sql, values).QueryRows(&list)
	if err!=nil{
		logs.Error(err)
	}

	return

}


//type LogApiDailyViewTemp struct {
//	ApiId           uint64 `description:"apiId" json:"api_id"`
//	ApiName         string `description:"month" json:"api_name"`
//	Successful  	 uint64 `description:"successful" json:"successful"`
//	Failed      	 uint64 `description:"failed" json:"failed"`
//	PublicModel
//}





type LogApiDailyViewTemp struct {
	ApiId           uint64 `description:"apiId" json:"api_id"`
	ApiName         string `description:"month" json:"api_name"`
	Successful  	 uint64 `description:"successful" json:"successful"`
	Failed      	 uint64 `description:"failed" json:"failed"`
	PublicModel
}

type LogApiDailyView struct {
	ApiId           uint64 `description:"apiId" json:"api_id"`
	ApiName         string `description:"month" json:"api_name"`
	//SuccessAndFailed 			uint64 `description:"success_and_failed" 	json:"success_and_failed"`
	Total 			uint64 `description:"total" 	json:"total"`
	Successful  	 uint64 `description:"successful" json:"successful"`
	Failed      	 uint64 `description:"failed" json:"failed"`
	FailedRate      float64 `description:"day" json:"failedRate"`

	PublicModel
}



func (this *LogApiDailyStatistic) GetLogQuotasStatistic(startTime,endTime string) (list []*LogApiDailyStatistic, err error) {

	o := orm.NewOrm()
	var values = make([]interface{}, 0)

	values=append(values, startTime)
	values=append(values, endTime)


	//  ApiIdName       this.TableName()
	sql :=`
SELECT api_id ,api_name as api_id_name,day_successful AS successful_called,day_failed as failed_called
FROM log_api_daily_statistic
where created_time>= ? and created_time< ?
`
	_, err = o.Raw(sql, values).QueryRows(&list)

	for i, i2 := range list {
		logs.Info(i,i2)

	}
	//logs.Info(list)
	return

}



func (this *LogApiDailyStatistic) UpdateByIdAndDay(apiId uint64,day string, domain string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id",apiId).Filter("day",day).Filter("domain",domain).Update(data)

	return err
}



type LogApiYesterday struct {

	DaySuccessful   uint64 `description:"day" json:"day_successful"`
	DayFailed       uint64 `description:"day" json:"day_failed"`
}


//  todo  不加domain
func (this *LogApiDailyStatistic) GetLogApiYesterdayTotal(yesterday string) (logApiYesterday *LogApiYesterday, err error) {

	o := orm.NewOrm()
	var values = make([]interface{}, 0)
	values=append(values, yesterday)

	sql :=`
SELECT  SUM(day_successful) as day_successful,SUM(day_failed) as day_failed FROM log_api_daily_statistic WHERE DAY = ?
`
	err = o.Raw(sql, values).QueryRow(&logApiYesterday)
	if err!=nil{
		logs.Error(err)
	}

	return

}





func (this *LogApiDailyStatistic) GetApiLabel(q map[string]interface{},domain string) (label []*string, err error) {
	label = make([]*string, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)
	//q["app_id"],q["api_id"],q["start_time"],q["end_time"]

	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}
	values=append(values,domain)
	//logs.Info("---->",q["app_id"],q["api_id"],q["start_time"],q["end_time"])
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	var sql string
	sql = `SELECT day  FROM  log_api_daily_statistic where  day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day;`

	_, err = o.Raw(sql, values).QueryRows(&label)

	return

}








func (this *LogApiDailyStatistic) GetApiCalledSuccess(q map[string]interface{},domain string) (calledSuccess []*uint64, err error) {
	calledSuccess = make([]*uint64, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)
	//q["app_id"],q["api_id"],q["start_time"],q["end_time"]
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}
	values=append(values,domain)
	//logs.Info("---->",q["app_id"],q["api_id"],q["start_time"],q["end_time"])
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	var sql string
	sql = `select cnt from ( SELECT day,sum(day_successful) cnt   FROM  log_api_daily_statistic where  day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day) t ;`
	_, err = o.Raw(sql, values).QueryRows(&calledSuccess)

	return

}






func (this *LogApiDailyStatistic) GetApiCalledFailed(q map[string]interface{},domain string) (calledFailed []*uint64, err error) {
	calledFailed = make([]*uint64, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)

	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}
	values=append(values,domain)

	//logs.Info("---->",q["app_id"],q["api_id"],q["start_time"],q["end_time"])
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值

	var sql string
	sql = `select cnt from ( SELECT day,sum(day_failed) cnt   FROM  log_api_daily_statistic where  day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day  ) t  ;`
	_, err = o.Raw(sql, values).QueryRows(&calledFailed)
	//}

	return

}




