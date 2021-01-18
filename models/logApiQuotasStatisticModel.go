package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type LogApiQuotasStatisticModelTemp struct {
	ApiId            uint64 `description:"apiId" json:"api_id"`
	ApiName        string `description:"成功调用次数" json:"api_name"`

	CalledTotal      uint64 `description:"成功调用次数" json:"total_called"`
	SuccessfulCalled uint64 `description:"成功调用次数" json:"successful_called"`
	FailedCalled     uint64 `description:"失败调用次数" json:"failed_called"`
	Domain      string  `description:"域名" json:"domain"`

	//Domain     string `description:"所属域" json:"domain"`
	//CreateTime string `json:"createdTime"`
	//ModifyTime string `json:"updatedTime"`
}

type LogApiQuotasStatisticModel struct {
	Id               uint64 `orm:"auto" description:"ID" json:"id"`
	ApiId            uint64 `description:"apiId" json:"api_id"`
	ApiName        string `description:"成功调用次数" json:"api_name"`
	SuccessfulCalled uint64 `description:"成功调用次数" json:"successful_called"`
	FailedCalled     uint64 `description:"失败调用次数" json:"failed_called"`
	Domain      string  `description:"域名" json:"domain"`

	PublicModel
}

func (this *LogApiQuotasStatisticModel) TableName() string {
	return "log_api_quotas_statistic"
}

// TableEngine 获取数据使用的引擎.
func (this *LogApiQuotasStatisticModel) TableEngine() string {
	return "INNODB"
}

func (this *LogApiQuotasStatisticModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewLogApiQuotasStatisticModel() *LogApiQuotasStatisticModel {
	return &LogApiQuotasStatisticModel{}
}

func (this *LogApiQuotasStatisticModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix()) //.Filter("domain", domain)
}

//添加设备
func (this *LogApiQuotasStatisticModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//修改设备
func (this *LogApiQuotasStatisticModel) Update(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id", id).Update(data)
	return err
}

func (this *LogApiQuotasStatisticModel) UpdateByApiId(id uint64, data map[string]interface{}, domain string) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id", id).Update(data)

	return err
}

//多个过滤条件
func (this *LogApiQuotasStatisticModel) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}




func (this *LogApiQuotasStatisticModel) GetAppQuotasStatisticModel(
	q map[string]interface{}) (list []*LogApiQuotasStatisticModel, Count int, err error) {

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

func (this *LogApiQuotasStatisticModel) GetCalledTotal(q map[string]interface{},domain string) (int64, error) {
	var totalSuccessFailed int64
	o := orm.NewOrm()
	//////////////////////
	var values = make([]interface{}, 0)
	sql := `select sum(successful_called) + sum(failed_called) as cnt  FROM  log_api_quotas_statistic WHERE 1=1 and domain = ?  `

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

type SuccessAndFailed struct {
	Success uint64 `json:"success"`
	Failed  uint64 `json:"failed"`
}

func (this *LogApiQuotasStatisticModel) GetSuccessAndFailed(q map[string]interface{},domain string) (successAndFailed SuccessAndFailed, err error) {

	o := orm.NewOrm()

	var values = make([]interface{}, 0)
	sql := `select sum(successful_called) as success , sum(failed_called) as failed FROM  log_api_quotas_statistic WHERE 1=1 and domain = ?   `

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



func (this *LogApiQuotasStatisticModel) UpdateByIdAndDay(apiId uint64, domain string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("api_id",apiId).Filter("domain",domain).Update(data)

	return err
}




func (this *AbilityCategoryModel) CategoryPieChart( domain string ) (list []*CategoryModelTemp, err error) {
	o := orm.NewOrm()
	var values = make([]interface{}, 0)
	//values=append(values, startTime)
	//values=append(values, endTime)

	//logs.Info("---- this.TableName()-->", this.TableName())
	//  ApiIdName       this.TableName()  //  openai_detail_log_20201221
	sql :=`
SELECT category_id,category_name,SUM(successful_called) as successful_called ,SUM(failed_called)  as failed_called
FROM 
(SELECT api_id,successful_called, failed_called FROM  log_api_quotas_statistic where domain= ? ) t1 
INNER JOIN 
(SELECT id ,ability_id FROM api_info )t2
ON t1.api_id=t2.id
INNER JOIN
(SELECT id,category_id from ability_info)  t3
ON  t2.ability_id=t3.id
INNER JOIN 
(SELECT id ,category_name FROM  ability_category_info ) t4
ON t3.category_id=t4.id
GROUP BY category_id,category_name; `

	values=append(values, domain)

	_, err = o.Raw(sql, values).QueryRows(&list)

	//for i, i2 := range list {
	//	logs.Info(i,"------>",i2)
	//}
	//logs.Info(list)
	return

}

