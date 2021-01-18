package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type LogHourlyQuotasStatisticModelTemp struct {
	Id               uint64 `orm:"auto" description:"ID" json:"id"`
	AppId            uint64 `description:"appId" json:"app_id"`
	ApiId            uint64 `description:"apiId" json:"api_id"`
	AppName          string `description:"month" json:"app_name"`
	ApiName          string `description:"month" json:"api_name"`
	Year             string `description:"month" json:"year"`
	YearSuccessful   uint64 `description:"month" json:"year_successful"`
	YearFailed       uint64 `description:"month" json:"year_failed"`
	Month            string `description:"month" json:"month"`
	MonthCount       uint64 `description:"month" json:"month_count"`
	Day              string `description:"day" json:"day"`
	DayCount         uint64 `description:"day" json:"day_count"`
	Hour             string `description:"hour" json:"hour"`
	HourCount        uint64 `description:"hour" json:"hour_count"`
	SuccessfulCalled uint64 `description:"成功调用次数" json:"successful_called"`
	FailedCalled     uint64 `description:"失败调用次数" json:"failed_called"`
	Domain      string  `description:"域名" json:"domain"`


	CreateTime string `json:"createdTime"`
	ModifyTime string `json:"updatedTime"`
}

type LogHourlyQuotasStatisticModel struct {
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
	Hour            string `description:"hour" json:"hour"`
	HourSuccessful  uint64 `description:"hour" json:"hour_successful"`
	HourFailed      uint64 `description:"hour" json:"hour_failed"`
	Domain      string  `description:"域名" json:"domain"`


	PublicModel
}

func (this *LogHourlyQuotasStatisticModel) TableName() string {
	return "log_hourly_quotas_statistic"
}

// TableEngine 获取数据使用的引擎.
func (this *LogHourlyQuotasStatisticModel) TableEngine() string {
	return "INNODB"
}

func (this *LogHourlyQuotasStatisticModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewHourlyQuotasStatisticModel() *LogHourlyQuotasStatisticModel {
	return &LogHourlyQuotasStatisticModel{}
}

func (this *LogHourlyQuotasStatisticModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix()) //.Filter("domain", domain)
}

//添加设备
func (this *LogHourlyQuotasStatisticModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//修改设备
func (this *LogHourlyQuotasStatisticModel) Update(id string, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *LogHourlyQuotasStatisticModel) UpdateById(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

//多个过滤条件
func (this *LogHourlyQuotasStatisticModel) filter(domain string, queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}

func (this *LogHourlyQuotasStatisticModel) GetCalledSuccessAndFailed(q map[string]interface{},domain string) (calledSuccessAndFailed []*uint64, err error) {
	calledSuccessAndFailed = make([]*uint64, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)
	//q["app_id"],q["api_id"],q["start_time"],q["end_time"]
	if value, ok := q["app_id"]; ok {
		values = append(values, value)
	}
	if value, ok := q["api_id"]; ok {
		values = append(values, value)
	}
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}

	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值

	var sql string
	//if q["api_id"] == nil {
	//	if q["hour"] != nil {
	//		sql = `select cnt from ( SELECT HOUR,sum(hour_failed)+SUM(hour_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY HOUR  ORDER BY HOUR ) t;`
	//	} else if q["day"] != nil {
	//		sql = `select cnt from ( SELECT day,sum(day_failed)+SUM(day_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY day  ORDER BY day ) t ;`
	//	} else  if q["month"] != nil{
	//		sql = `select cnt from ( SELECT month,sum(month_failed)+SUM(month_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY month  ORDER BY month ) t;`
	//	}else {
	//		sql = `select cnt from ( SELECT year,sum(year_failed)+SUM(year_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY year  ORDER BY year ) t;`
	//	}
	//	_, err = o.Raw(sql, q["app_id"], q["start_time"], q["end_time"]).QueryRows(&calledSuccessAndFailed)
	//}else

	if q["app_id"] == 0 {
		if q["hour"] != nil {
			sql = `select cnt from ( SELECT HOUR,sum(hour_failed)+SUM(hour_successful) cnt   FROM  log_hourly_quotas_statistic where  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR ) t;`
		} else if q["day"] != nil {
			sql = `select cnt from ( SELECT day,sum(day_failed)+SUM(day_successful) cnt   FROM  log_daily_quotas_statistic where  api_id = ? and day >= ? and day <= ?   and domain= ?  GROUP BY day  ORDER BY day ) t ;`
		} else if q["month"] != nil {
			sql = `select cnt from ( SELECT month,sum(month_failed)+SUM(month_successful) cnt   FROM  log_daily_quotas_statistic where  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
		} else {
			sql = `select cnt from ( SELECT year,sum(year_failed)+SUM(year_successful) cnt   FROM  log_daily_quotas_statistic where  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
		}


		_, err = o.Raw(sql, q["api_id"], q["start_time"], q["end_time"],domain).QueryRows(&calledSuccessAndFailed)
	} else {
		if q["hour"] != nil {
			sql = `select cnt from ( SELECT HOUR,sum(hour_failed)+SUM(hour_successful) cnt   FROM  log_hourly_quotas_statistic where app_id = ? and  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR ) t;`
		} else if q["day"] != nil {
			sql = `select cnt from ( SELECT day,sum(day_failed)+SUM(day_successful) cnt   FROM  log_daily_quotas_statistic where app_id = ? and  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day ) t ;`
		} else if q["month"] != nil {
			sql = `select cnt from ( SELECT month,sum(month_failed)+SUM(month_successful) cnt   FROM  log_daily_quotas_statistic where  app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
		} else {
			sql = `select cnt from ( SELECT year,sum(year_failed)+SUM(year_successful) cnt   FROM  log_daily_quotas_statistic where  app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
		}

		values=append(values,domain)

		_, err = o.Raw(sql, values).QueryRows(&calledSuccessAndFailed)
	}

	return

}

func (this *LogHourlyQuotasStatisticModel) GetCalledFailed(q map[string]interface{},domain string) (calledFailed []*uint64, err error) {
	calledFailed = make([]*uint64, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)
	//q["app_id"],q["api_id"],q["start_time"],q["end_time"]
	if value, ok := q["app_id"]; ok {
		values = append(values, value)
	}
	if value, ok := q["api_id"]; ok {
		values = append(values, value)
	}
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}

	var sql string

	if q["app_id"] == 0 {
		if q["hour"] != nil {
			sql = `select cnt from ( SELECT HOUR,sum(hour_failed) cnt   FROM  log_hourly_quotas_statistic where  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR ) t;`
		} else if q["day"] != nil {
			sql = `select cnt from ( SELECT day,sum(day_failed) cnt   FROM  log_daily_quotas_statistic where  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day ) t ;`
		} else if q["month"] != nil {
			sql = `select cnt from ( SELECT month,sum(month_failed) cnt   FROM  log_daily_quotas_statistic where api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
		} else {
			sql = `select cnt from ( SELECT year,sum(year_failed) cnt   FROM  log_daily_quotas_statistic where api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
		}


		_, err = o.Raw(sql, q["api_id"], q["start_time"], q["end_time"],domain).QueryRows(&calledFailed)

	} else {
		if q["hour"] != nil {
			sql = `select cnt from ( SELECT HOUR,sum(hour_failed) cnt   FROM  log_hourly_quotas_statistic where app_id = ?  and  api_id = ? and day >= ? and day <= ?   and domain= ?  GROUP BY HOUR  ORDER BY HOUR ) t;`
		} else if q["day"] != nil {
			sql = `select cnt from ( SELECT day,sum(day_failed) cnt   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?   and domain= ?  GROUP BY day  ORDER BY day ) t ;`
		} else if q["month"] != nil {
			sql = `select cnt from ( SELECT month,sum(month_failed) cnt   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
		} else {
			sql = `select cnt from ( SELECT year,sum(year_failed) cnt   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
		}

		values=append(values,domain)

		_, err = o.Raw(sql, values).QueryRows(&calledFailed)
	}

	return

}

func (this *LogHourlyQuotasStatisticModel) GetCalledSuccess(q map[string]interface{},domain string) (calledSuccess []*uint64, err error) {
	calledSuccess = make([]*uint64, 0)
	o := orm.NewOrm()

	var values = make([]interface{}, 0)
	if value, ok := q["app_id"]; ok {
		values = append(values, value)
	}
	if value, ok := q["api_id"]; ok {
		values = append(values, value)
	}
	//logs.Info("======================")
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}

	var sql string

	if q["app_id"] == 0 {

		//logs.Info(q["hour"], q["day"], q["month"])
		if q["hour"] != nil {
			sql = `select cnt from ( SELECT HOUR,sum(hour_successful) cnt   FROM  log_hourly_quotas_statistic where api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR ) t;`
		} else if q["day"] != nil {
			//logs.Info(q["day"])
			sql = `select cnt from ( SELECT day,sum(day_successful) cnt   FROM  log_daily_quotas_statistic where  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day ) t ;`
			//logs.Info("sql->", sql)
		} else if q["month"] != nil {
			sql = `select cnt from ( SELECT month,sum(month_successful) cnt   FROM  log_daily_quotas_statistic where api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
		} else {
			sql = `select cnt from ( SELECT year,sum(year_successful) cnt   FROM  log_daily_quotas_statistic where api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
		}

		_, err = o.Raw(sql, q["api_id"], q["start_time"], q["end_time"],domain).QueryRows(&calledSuccess)

	} else {
		if q["hour"] != nil {
			sql = `select cnt from ( SELECT HOUR,sum(hour_successful) cnt   FROM  log_hourly_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR ) t;`
		} else if q["day"] != nil {
			sql = `select cnt from ( SELECT day,sum(day_successful) cnt   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day ) t ;`
		} else if q["month"] != nil {
			sql = `select cnt from ( SELECT month,sum(month_successful) cnt   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
		} else {
			sql = `select cnt from ( SELECT year,sum(year_successful) cnt   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
		}
		values=append(values,domain)
		_, err = o.Raw(sql, values).QueryRows(&calledSuccess)
	}


	//
	//logs.Info("--sql-->", sql, "--values--", values)
	//logs.Info("--calledSuccess--", calledSuccess)
	return

}

func timeHandler(q map[string]interface{}, values []interface{}, sql string) ([]interface{}, string) {
	//logs.Info("----------------",q["start_time"],q["end_time"])
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
		sql = sql + " and day >= ? "
	}

	if value, ok := q["end_time"]; ok {
		values = append(values, value)
		sql = sql + " and day <= ? "
	}
	return values, sql
}

func (this *LogHourlyQuotasStatisticModel) GetLabel(q map[string]interface{},domain string) (label []*string, err error) {
	label = make([]*string, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)
	//q["app_id"],q["api_id"],q["start_time"],q["end_time"]
	if value, ok := q["app_id"]; ok {
		values = append(values, value)
	}
	if value, ok := q["api_id"]; ok {
		values = append(values, value)
		//logs.Info("api_id------>",q["api_id"])
	}
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}

	//logs.Info("---->",q["app_id"],q["api_id"],q["start_time"],q["end_time"])
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	var sql string

	if q["app_id"] == 0 {
		if q["hour"] != nil {
			sql = `SELECT hour as count   FROM  log_hourly_quotas_statistic where  api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR;`
		} else if q["day"] != nil {
			sql = `SELECT day as count   FROM  log_daily_quotas_statistic where  api_id = ? and day >= ? and day <= ?   and domain= ?  GROUP BY day  ORDER BY day;`
		} else if q["month"] != nil {
			sql = `SELECT month as count   FROM  log_daily_quotas_statistic where api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month;`
		} else {
			sql = `SELECT year as count   FROM  log_daily_quotas_statistic where api_id = ? and day >= ? and day <= ?   and domain= ?  GROUP BY year  ORDER BY year;`
		}

		_, err = o.Raw(sql, q["api_id"], q["start_time"], q["end_time"],domain).QueryRows(&label)
	} else {
		if q["hour"] != nil {
			sql = `SELECT hour as count   FROM  log_hourly_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR;`
		} else if q["day"] != nil {
			sql = `SELECT day as count   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day;`
		} else if q["month"] != nil {
			sql = `SELECT month as count   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month;`
		} else {
			sql = `SELECT year as count   FROM  log_daily_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year;`
		}

		values=append(values,domain)

		_, err = o.Raw(sql, values).QueryRows(&label)

	}

	return

}

// 1.得到所有label
// 2.得到所有data [内]
//   2.1 得到 调用成功
//   2.2 得到 调用失败
//   2.3 得到 调用总数
func (this *LogHourlyQuotasStatisticModel) GetMonitorChartDimension(q map[string]interface{}) (monitors []*LogHourlyQuotasStatisticModel, totalCount int64, err error) {

	o := orm.NewOrm()
	monitors = make([]*LogHourlyQuotasStatisticModel, 0)

	var values = make([]interface{}, 0)
	//q["app_id"],q["api_id"],q["start_time"],q["end_time"]
	if value, ok := q["app_id"]; ok {
		values = append(values, value)
	}
	if value, ok := q["api_id"]; ok {
		values = append(values, value)
	}
	if value, ok := q["start_time"]; ok {
		values = append(values, value)
	}
	if value, ok := q["end_time"]; ok {
		values = append(values, value)
	}

	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	var sql string
	if q["month"] != nil {
		//monthSql := `SELECT month as dimension,sum(successful_called) as count FROM  hourly_quotas_statistic where app_id = ? and api_id = ? and created_time >= ? and created_time <= ? GROUP BY month`
		if monitorItem, ok := q["monitor_item"]; ok { //   monitor_item 1/调用成功 2/调用失败 3/QPS峰值
			if monitorItem == 1 {
				sql = `SELECT month as dimension,sum(successful_called) as count FROM  log_hourly_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ? GROUP BY month`
			} else {
				sql = `SELECT month as dimension,sum(failed_called) as count FROM  log_hourly_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ? GROUP BY month`
			}
		}
	} else if q["day"] != nil {
		//daySql:=`SELECT day as dimension,sum(successful_called) as count FROM  hourly_quotas_statistic where app_id = ? and api_id = ? and created_time >= ? and created_time <= ? GROUP BY day `
		if monitorItem, ok := q["monitor_item"]; ok { //   monitor_item 1/调用成功 2/调用失败 3/QPS峰值
			if monitorItem == 1 {
				sql = `SELECT month as dimension,sum(successful_called) as count FROM  log_hourly_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ? GROUP BY month`
			} else {
				sql = `SELECT month as dimension,sum(failed_called) as count FROM  log_hourly_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ? GROUP BY month`
			}
		}
	} else {
		//hourSql:=`SELECT hour as dimension,sum(successful_called) as count FROM  hourly_quotas_statistic where app_id = ? and api_id = ? and created_time >= ? and created_time <= ? GROUP BY hour `
		if monitorItem, ok := q["monitor_item"]; ok { //
			if monitorItem == 1 {
				sql = `SELECT month as dimension,sum(successful_called) as count FROM  log_hourly_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ? GROUP BY month`
			} else {
				sql = `SELECT month as dimension,sum(failed_called) as count FROM  log_hourly_quotas_statistic where app_id = ? and api_id = ? and day >= ? and day <= ? GROUP BY month`
			}
		}
	}
	totalCount, err = o.Raw(sql, values).QueryRows(&monitors)

	return

}

/////////////////// api ///////////////////////////////////////////////////

func (this *LogHourlyQuotasStatisticModel) GetApiCalledSuccessAndFailed(q map[string]interface{},domain string) (calledSuccessAndFailed []*uint64, err error) {
	calledSuccessAndFailed = make([]*uint64, 0)
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

	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值

	var sql string
	sql = `select cnt from (SELECT day,sum(day_successful)+sum(day_failed) cnt   FROM  log_hourly_quotas_statistic where  day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day) t;`
	_, err = o.Raw(sql, values).QueryRows(&calledSuccessAndFailed)

	return

}

////////////////////////////////////   AppCalledTotalTrend    ////////////////////////////////////

// 此处对时间正确的处理逻辑：如果 startTimeOk && endTimeOK 不传，默认给7天前和当前时间 。前端带来的不确定因素，在controller解决，不要留在model中处理
// 调用总量趋势统计
func (this *LogHourlyQuotasStatisticModel) AppCalledTotalTrendLabel(q map[string]interface{},domain string) (label []*string, err error) {
	label = make([]*string, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)

	//logs.Info("---->",q["app_id"],q["api_id"],q["start_time"],q["end_time"])
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	var sql string
	if q["api_id"] == nil {

		if value, ok := q["app_id"]; ok {
			values = append(values, value)
		}

		//logs.Info(" h: ",q["hour"]," d: ", q["day"], " m: ",q["month"])

		startTime, _ := q["start_time"]
		endTime, _ := q["end_time"]


		if q["hour"] != nil {

			sql = `SELECT hour as count   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?   and domain= ?   GROUP BY HOUR  ORDER BY HOUR;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `SELECT hour as count   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY HOUR  ORDER BY HOUR;`
		} else if q["day"] != nil {

			sql = `SELECT day as count   FROM  log_daily_quotas_statistic where  app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `SELECT day as count   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY day  ORDER BY day;`
		} else if q["month"] != nil {

			sql = `SELECT month as count   FROM  log_daily_quotas_statistic where app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month;`
			values = append(values, startTime)
			values = append(values, endTime)

		} else {

			sql = `SELECT year as count   FROM  log_daily_quotas_statistic where app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `SELECT year as count   FROM  log_hourly_quotas_statistic where app_id = ? and day >= ? and day <= ?  GROUP BY year  ORDER BY year;`
		}

		values=append(values,domain)

		_, err = o.Raw(sql, values).QueryRows(&label)

	}

	return

}

//   调用总量趋势统计 GetCalledSuccess
func (this *LogHourlyQuotasStatisticModel) AppCalledTotalTrendSuccess(q map[string]interface{},domain string) (calledSuccess []*uint64, err error) {
	calledSuccess = make([]*uint64, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)

	var sql string
	if q["api_id"] == nil {

		if value, ok := q["app_id"]; ok {
			values = append(values, value)
		}

		//logs.Info(" h: ",q["hour"]," d: ", q["day"], " m: ",q["month"])
		startTime, _ := q["start_time"]
		endTime, _ := q["end_time"]

		if q["hour"] != nil {

			sql = `select cnt from ( SELECT HOUR,sum(hour_successful) cnt   FROM  log_hourly_quotas_statistic where app_id = ? and day >= ? and day <= ?  and domain= ?    GROUP BY HOUR  ORDER BY HOUR ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

		} else if q["day"] != nil {

			sql = `select cnt from ( SELECT day,sum(day_successful) cnt   FROM  log_daily_quotas_statistic where  app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day ) t ;`
			values = append(values, startTime)
			values = append(values, endTime)

		} else if q["month"] != nil {

			sql = `select cnt from ( SELECT month,sum(month_successful) cnt   FROM  log_daily_quotas_statistic where app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT month,sum(month_successful) cnt   FROM  log_hourly_quotas_statistic where app_id = ? and day >= ? and day <= ?  GROUP BY month  ORDER BY month ) t;`

		} else if q["year"] != nil {

			sql = `select cnt from ( SELECT year,sum(year_successful) cnt   FROM  log_daily_quotas_statistic where app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

		}
		values=append(values,domain)

		_, err = o.Raw(sql, values).QueryRows(&calledSuccess)

	}


	//logs.Info("--sql-->", sql, "--values--", values)
	//logs.Info("--calledSuccess--", calledSuccess)
	return

}

func (this *LogHourlyQuotasStatisticModel) AppCalledTotalTrendFailed(q map[string]interface{},domain string) (calledFailed []*uint64, err error) {
	calledFailed = make([]*uint64, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)

	var sql string

	if q["api_id"] == nil {

		if value, ok := q["app_id"]; ok {
			values = append(values, value)
		}

		//logs.Info(" h: ",q["hour"]," d: ", q["day"], " m: ",q["month"])
		startTime, _ := q["start_time"]
		endTime, _ := q["end_time"]

		if q["hour"] != nil {

			//logs.Info("startTimeOk",startTimeOk,endTimeOK,endTimeOK)
			sql = `select cnt from ( SELECT HOUR,sum(hour_failed) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT HOUR,sum(hour_failed) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY HOUR  ORDER BY HOUR ) t;`
		} else if q["day"] != nil {

			sql = `select cnt from ( SELECT day,sum(day_failed) cnt   FROM  log_daily_quotas_statistic where  app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day ) t ;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT day,sum(day_failed) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY day  ORDER BY day ) t ;`
		} else if q["month"] != nil {

			sql = `select cnt from ( SELECT month,sum(month_failed) cnt   FROM  log_daily_quotas_statistic where app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT month,sum(month_failed) cnt   FROM  log_hourly_quotas_statistic where app_id = ? and day >= ? and day <= ?  GROUP BY month  ORDER BY month ) t;`
		} else {

			sql = `select cnt from ( SELECT year,sum(year_failed) cnt   FROM  log_daily_quotas_statistic where app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT year,sum(year_failed) cnt   FROM  log_hourly_quotas_statistic where app_id = ? and day >= ? and day <= ?  GROUP BY year  ORDER BY year ) t;`
		}
		values=append(values,domain)

		_, err = o.Raw(sql, values).QueryRows(&calledFailed)

	}

	return

}

func (this *LogHourlyQuotasStatisticModel) AppCalledTotalTrendSuccessAndFailed(q map[string]interface{},domain string) (calledSuccessAndFailed []*uint64, err error) {
	calledSuccessAndFailed = make([]*uint64, 0)
	o := orm.NewOrm()
	//label = make([]*Label, 0)

	var values = make([]interface{}, 0)

	var sql string
	if q["api_id"] == nil {

		if value, ok := q["app_id"]; ok {
			values = append(values, value)
		}

		//logs.Info(" h: ",q["hour"]," d: ", q["day"], " m: ",q["month"])
		startTime, _ := q["start_time"]
		endTime, _ := q["end_time"]

		if q["hour"] != nil {

			sql = `select cnt from ( SELECT HOUR,sum(hour_failed)+SUM(hour_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY HOUR  ORDER BY HOUR ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT HOUR,sum(hour_failed)+SUM(hour_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY HOUR  ORDER BY HOUR ) t;`
		} else if q["day"] != nil {

			sql = `select cnt from ( SELECT day,sum(day_failed)+SUM(day_successful) cnt   FROM  log_daily_quotas_statistic where  app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY day  ORDER BY day ) t ;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT day,sum(day_failed)+SUM(day_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY day  ORDER BY day ) t ;`
		} else if q["month"] != nil {

			sql = `select cnt from ( SELECT month,sum(month_failed)+SUM(month_successful) cnt   FROM  log_daily_quotas_statistic where  app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY month  ORDER BY month ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT month,sum(month_failed)+SUM(month_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY month  ORDER BY month ) t;`
		} else {

			sql = `select cnt from ( SELECT year,sum(year_failed)+SUM(year_successful) cnt   FROM  log_daily_quotas_statistic where  app_id = ? and day >= ? and day <= ?  and domain= ?   GROUP BY year  ORDER BY year ) t;`
			values = append(values, startTime)
			values = append(values, endTime)

			//sql = `select cnt from ( SELECT year,sum(year_failed)+SUM(year_successful) cnt   FROM  log_hourly_quotas_statistic where  app_id = ? and day >= ? and day <= ?  GROUP BY year  ORDER BY year ) t;`
		}
		values=append(values,domain)

		_, err = o.Raw(sql, values).QueryRows(&calledSuccessAndFailed)
	}

	return

}

func (this *LogHourlyQuotasStatisticModel) UpdateByIdAndHour(apiId, appId uint64, hour string, domain string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id", apiId).Filter("app_id", appId).Filter("hour", hour).Filter("domain",domain).Update(data)

	return err
}

func (this *LogHourlyQuotasStatisticModel) GetLogAppDailyStatistic(startTime, endTime string) (list []*LogAppDailyStatisticModel, err error) {

	o := orm.NewOrm()
	var values = make([]interface{}, 0)

	values = append(values, startTime)
	values = append(values, endTime)

	//  ApiIdName       this.TableName()
	sql := `
SELECT app_id, app_name ,
,year,month,day,successful_called as year_successful ,failed_called as year_failed,
month,successful_called as month_successful ,failed_called as month_failed,
DAY,successful_called as day_successful ,failed_called as day_failed
FROM (
SELECT app_id, app_name ,year,month,day,SUM(hour_successful) as successful_called ,SUM(hour_failed) as failed_called
FROM (
SELECT app_id ,app_name,year,month,day,hour_successful,hour_failed
FROM log_hourly_quotas_statistic
where created_time>= ? and created_time< ?
) t  GROUP BY app_id,app_name,year,month,DAY
)tt `
	//logs.Info("------2------->",sql,values)

	_, err = o.Raw(sql, values).QueryRows(&list)

	//for i, i2 := range list {
	//	logs.Info(i,i2)
	//
	//}
	//logs.Info(list)
	return

}

func (this *LogHourlyQuotasStatisticModel) GetLogApiDailyStatistic(startTime, endTime string) (list []*LogApiDailyStatistic, err error) {

	o := orm.NewOrm()
	var values = make([]interface{}, 0)

	values = append(values, startTime)
	values = append(values, endTime)

	//  ApiIdName       this.TableName()
	sql := `
SELECT app_id, app_name ,
year,successful_called as year_successful ,failed_called as year_failed,
month,successful_called as month_successful ,failed_called as month_failed,
DAY,successful_called as day_successful ,failed_called as day_failed
FROM (
SELECT app_id, app_name ,year,month,day,SUM(hour_successful) as successful_called ,SUM(hour_failed) as failed_called
FROM (
SELECT app_id ,app_name,year,month,day,hour_successful,hour_failed
FROM log_hourly_quotas_statistic
where created_time>= ? and created_time< ?
) t  GROUP BY app_id,app_name,year,month,DAY
)tt `
	//logs.Info("------2------->",sql,values)

	_, err = o.Raw(sql, values).QueryRows(&list)

	//for i, i2 := range list {
	//	logs.Info(i,i2)
	//
	//}
	//logs.Info(list)
	return

}
