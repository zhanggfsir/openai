package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type LogApiDailyHistoryAddUpTemp struct {

	Id          		 			uint64 `orm:"auto" description:"ID" json:"id"`
	Day             	 			string `description:"day" json:"day"`
	DayTotal        	 			uint64 `description:"day_total" json:"day_total"`
	DaySuccessful   	 			uint64 `description:"day" json:"day_successful"`
	DayFailed       	 			uint64 `description:"day" json:"day_failed"`
	DailyHistoryAddUp             	uint64 `description:"day_total" json:"daily_history_add_up"`
	DailySuccessfulHistoryAddUp   	uint64 `description:"day" json:"daily_successful_history_add_up"`
	DailyFailedHistoryAddUp       	uint64 `description:"day" json:"daily_failed_history_add_up"`
	Domain      string  `description:"域名" json:"domain"`

	//Domain     string `description:"所属域" json:"domain"`
	//CreateTime string `json:"createdTime"`
	//ModifyTime string `json:"updatedTime"`
}

type LogApiDailyHistoryAddUp struct {
	Id          		 			uint64 `orm:"auto" description:"ID" json:"id"`
	Day             	 			string `description:"day" json:"day"`
	DayTotal        	 			uint64 `description:"day_total" json:"day_total"`
	DaySuccessful   	 			uint64 `description:"day" json:"day_successful"`
	DayFailed       	 			uint64 `description:"day" json:"day_failed"`
	DailyHistoryAddUp             	uint64 `description:"day_total" json:"daily_history_add_up"`
	DailySuccessfulHistoryAddUp   	uint64 `description:"day" json:"daily_successful_history_add_up"`
	DailyFailedHistoryAddUp       	uint64 `description:"day" json:"daily_failed_history_add_up"`
	Domain      string  `description:"域名" json:"domain"`

	PublicModel
}



func (this *LogApiDailyHistoryAddUp) TableName() string {
	return "log_api_daily_history_addup"
}

// TableEngine 获取数据使用的引擎.
func (this *LogApiDailyHistoryAddUp) TableEngine() string {
	return "INNODB"
}

func (this *LogApiDailyHistoryAddUp) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewLogApiDailyHistoryAddUp() *LogApiDailyHistoryAddUp {
	return &LogApiDailyHistoryAddUp{}
}

func (this *LogApiDailyHistoryAddUp) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix()) //.Filter("domain", domain)
}

//添加设备
func (this *LogApiDailyHistoryAddUp) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//修改设备
func (this *LogApiDailyHistoryAddUp) Update(id string, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *LogApiDailyHistoryAddUp) UpdateById(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

//多个过滤条件
func (this *LogApiDailyHistoryAddUp) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}



func (this *LogApiDailyHistoryAddUp) GetApiLabel(q map[string]interface{}) (label []*string, err error) {
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
	//logs.Info("---->",q["app_id"],q["api_id"],q["start_time"],q["end_time"])
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	var sql string
	sql = `SELECT day  FROM  log_api_daily_history_addup where  day >= ? and day <= ?  ORDER BY day;`

	_, err = o.Raw(sql, values).QueryRows(&label)

	return

}




func (this *LogApiDailyHistoryAddUp) GetApiCalledSuccess(q map[string]interface{}  ) (calledSuccess []*uint64, err error) {
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

	//logs.Info("---->",q["app_id"],q["api_id"],q["start_time"],q["end_time"])
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值
	var sql string
	sql = `select daily_successful_history_add_up   FROM  log_api_daily_history_addup  where  day >= ? and day <= ?    ORDER BY day ;`
	_, err = o.Raw(sql, values).QueryRows(&calledSuccess)

	return

}


func (this *LogApiDailyHistoryAddUp) GetApiCalledFailed(q map[string]interface{}) (calledFailed []*uint64, err error) {
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

	//logs.Info("---->",q["app_id"],q["api_id"],q["start_time"],q["end_time"])
	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值

	var sql string
	//sql = `select cnt from ( SELECT day,sum(day_failed) cnt   FROM  log_hourly_quotas_statistic where  day >= ? and day <= ?  GROUP BY day  ORDER BY day  ) t  ;`
	sql = `select daily_failed_history_add_up   FROM  log_api_daily_history_addup  where  day >= ? and day <= ?    ORDER BY day ;`
	_, err = o.Raw(sql, values).QueryRows(&calledFailed)

	//}

	return

}


func (this *LogApiDailyHistoryAddUp) GetApiCalledSuccessAndFailed(q map[string]interface{} ) (calledSuccessAndFailed []*uint64, err error) {
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


	//monitor_item 1/调用成功 2/调用失败 3/QPS峰值

	var sql string
	sql = `SELECT daily_history_add_up cnt   FROM  log_api_daily_history_addup where  day >= ? and day <= ?   ORDER BY DAY `
	_, err = o.Raw(sql, values).QueryRows(&calledSuccessAndFailed)

	return

}




func (this *LogApiDailyHistoryAddUp) UpdateByIdAndDay(apiId uint64,day string, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id",apiId).Filter("day",day).Filter("domain",domain).Update(data)

	return err
}






