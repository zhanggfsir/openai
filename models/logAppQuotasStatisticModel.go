package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type LogAppQuotasStatisticModelTemp struct {
	AppId            uint64 `description:"apiId" json:"app_id"`
	AppName        string `description:"成功调用次数" json:"app_name"`

	CalledTotal      uint64 `description:"成功调用次数" json:"total_called"`
	SuccessfulCalled uint64 `description:"成功调用次数" json:"successful_called"`
	FailedCalled     uint64 `description:"失败调用次数" json:"failed_called"`
	Domain      string  `description:"域名" json:"domain"`

	//Domain     string `description:"所属域" json:"domain"`
	//CreateTime string `json:"createdTime"`
	//ModifyTime string `json:"updatedTime"`
}

type LogAppQuotasStatisticModel struct {
	Id               uint64 `orm:"auto" description:"ID" json:"id"`
	AppId            uint64 `description:"apiId" json:"app_id"`
	AppName        string `description:"成功调用次数" json:"app_name"`
	SuccessfulCalled uint64 `description:"成功调用次数" json:"successful_called"`
	FailedCalled     uint64 `description:"失败调用次数" json:"failed_called"`
	Domain      string  `description:"域名" json:"domain"`

	PublicModel
}

func (this *LogAppQuotasStatisticModel) TableName() string {
	return "log_app_quotas_statistic"
}

// TableEngine 获取数据使用的引擎.
func (this *LogAppQuotasStatisticModel) TableEngine() string {
	return "INNODB"
}

func (this *LogAppQuotasStatisticModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewLogAppQuotasStatisticModel() *LogAppQuotasStatisticModel {
	return &LogAppQuotasStatisticModel{}
}

func (this *LogAppQuotasStatisticModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix()) //.Filter("domain", domain)
}

//添加设备
func (this *LogAppQuotasStatisticModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//修改设备
func (this *LogAppQuotasStatisticModel) Update(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id", id).Update(data)
	return err
}

func (this *LogAppQuotasStatisticModel) UpdateByApiId(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id", id).Update(data)

	return err
}

//多个过滤条件
func (this *LogAppQuotasStatisticModel) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}




func (this *LogAppQuotasStatisticModel) GetAppQuotasStatisticModel(
	q map[string]interface{}) (list []*LogAppQuotasStatisticModel, Count int, err error) {

	o := orm.NewOrm()
	var values = make([]interface{}, 0)
	// todo sql 拼接
	sql := `SELECT * from log_quotas_statistic WHERE api_id IN (SELECT id AS api_id from api_info WHERE ability_id IN(SELECT id FROM  ability_info WHERE category_id = ? )   )   `
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

func (this *LogAppQuotasStatisticModel) GetCalledTotal(q map[string]interface{}) (int64, error) {
	var totalSuccessFailed int64
	o := orm.NewOrm()
	//////////////////////
	var values = make([]interface{}, 0)
	sql := `select sum(successful_called) + sum(failed_called) as cnt  FROM  log_quotas_statistic WHERE 1=1   `

	if q["start_time"]!=""{
		values = append(values, q["start_time"])
		sql = sql + " and created_time >= ? "
	}
	if q["end_time"]!=""{
		values = append(values, q["end_time"])
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


	err := o.Raw(sql,values).QueryRow(&totalSuccessFailed)

	return totalSuccessFailed, err
}


func (this *LogAppQuotasStatisticModel) GetSuccessAndFailed(q map[string]interface{}) (successAndFailed SuccessAndFailed, err error) {

	o := orm.NewOrm()

	var values = make([]interface{}, 0)
	sql := `select sum(successful_called) as success , sum(failed_called) as failed FROM  log_quotas_statistic WHERE 1=1   `

	if q["start_time"]!=""{
		values = append(values, q["start_time"])
		sql = sql + " and created_time >= ? "
	}
	if q["end_time"]!=""{
		values = append(values, q["end_time"])
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

	err = o.Raw(sql,values).QueryRow(&successAndFailed)
	return successAndFailed, err
}



func (this *LogAppQuotasStatisticModel) UpdateByIdAndDay(appId uint64, domain string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("app_id",appId).Filter("domain",domain).Update(data)

	return err
}




func (this *LogAppQuotasStatisticModel) GetAppCalledTotal(q map[string]interface{},domain string) (int64, error) {
	var totalSuccessFailed int64
	o := orm.NewOrm()
	//////////////////////
	var values = make([]interface{}, 0)
	sql := `select sum(successful_called) + sum(failed_called) as cnt  FROM  log_quotas_statistic_app WHERE 1=1  and domain= ?      `

	values=append(values,domain)

	if q["start_time"]!=""{
		values = append(values, q["start_time"])
		sql = sql + " and created_time >= ? "
	}
	if q["end_time"]!=""{
		values = append(values, q["end_time"])
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


	err := o.Raw(sql,values).QueryRow(&totalSuccessFailed)

	return totalSuccessFailed, err
}





	func (this *LogAppQuotasStatisticModel) GetAppSuccessAndFailed(q map[string]interface{},domain string) (successAndFailed SuccessAndFailed, err error) {

		o := orm.NewOrm()

		var values = make([]interface{}, 0)
		sql := `select sum(successful_called) as success , sum(failed_called) as failed FROM  log_quotas_statistic_app WHERE 1=1   and domain= ?   `

		values=append(values,domain)


		if q["start_time"]!=""{
			values = append(values, q["start_time"])
			sql = sql + " and created_time >= ? "
		}
		if q["end_time"]!=""{
			values = append(values, q["end_time"])
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

		err = o.Raw(sql,values).QueryRow(&successAndFailed)
		return successAndFailed, err
	}
