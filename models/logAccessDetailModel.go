package models

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

/*
@Author:
@Time: 2020-11-05 15:02
@Description:APImodel
*/

type AccessLogTemp struct {
	Id           uint64 `orm:"auto" description:"ID" json:"id"`
	AppId        uint64 `description:"AppId" json:"AppId"`
	ApiId        uint64 `description:"ApiId" json:"ApiId"`
	AppName      string `description:"app名称" json:"appName"`
	ApiName      string `description:"api名称" json:"apiName"`
	Domain      string  `description:"域名" json:"domain"`
	Year         string `description:"2020-11" json:"year"`
	month        string `description:"2020-11" json:"month"`
	day          string `description:"2020-11-16" json:"day"`
	hour         string `description:"2020-11-16 18" json:"hour"`
	CalledStatus string `description:"calledStatus" json:"calledStatus"`
	Url     	 string `description:"url" json:"url"`

	CreateTime   string `json:"createdTime"`
	ModifyTime   string `json:"updatedTime"`
}

type LogAccessDetail struct {
	Id           uint64 `orm:"auto" description:"ID" json:"id"`
	AppId        uint64 `description:"AppId" json:"AppId"`
	ApiId        uint64 `description:"ApiId" json:"ApiId"`
	AppName      string `description:"app名称" json:"appName"`
	ApiName      string `description:"api名称" json:"apiName"`
	Domain      string  `description:"域名" json:"domain"`
	Year         string `description:"2020-11" json:"year"`
	Month        string `description:"2020-11" json:"month"`
	Day          string `description:"2020-11-16" json:"day"`
	Hour         string `description:"2020-11-16 18" json:"hour"`
	CalledStatus string `description:"calledStatus" json:"calledStatus"`
	Url     	 string `description:"url" json:"url"`

	CreatedTime  time.Time  `description:"created_time" json:"created_time"`
	//CreateTime time.Time  `orm:"column(created_time);type(datetime);auto_now_add" json:"created_time"`
	UpdateTime   time.Time  `orm:"column(updated_time);type(datetime);auto_now" json:"updated_time"`
}

func (this *LogAccessDetail) TableName() string {
	return "openai_detail_log_"+time.Now().Format("20060102")
	//return  "openai_detail_log_20201230"
}

// TableEngine 获取数据使用的引擎.
func (this *LogAccessDetail) TableEngine() string {
	return "INNODB"
}

func (this *LogAccessDetail) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func (this *LogAccessDetail) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

func NewLogAccessDetail() *LogAccessDetail {
	return &LogAccessDetail{}
}

//添加
func (this *LogAccessDetail) Insert() error {
	var oo = orm.NewOrm()
	//插入到数据库
	if _, err := oo.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}


var OpenaiDetailLog ="openai_detail_log"

//插入es
func (info *LogAccessDetail) InsertElasticFromKafka(values []interface{}) error {

	err:=info.InsertMySql(values)
	if err!=nil{
		logs.Info(err)
		return err
	}
	resp, err := EsClient.Index().Index(esIndex(OpenaiDetailLog)).BodyJson(
		info).Do(context.Background())
	if err != nil {
		logs.Error("Elastic err", err, "operation=", resp)
		return err
	}
	return err
}

//插入es
func (info *LogAccessDetail) InsertMySql(values []interface{}) error {

	err:=Insert(info,values)
	if err!=nil{
		logs.Info(err)
		return err
	}

	return nil
}

//  data["updated_time"] = time.Now()
func Insert(this *LogAccessDetail,values []interface{}) error{

	o := orm.NewOrm()
	_,err:=o.Raw(` insert into `+ this.TableName()+` (app_id,api_id,app_name,api_name,year,month,day,hour,called_status,url,domain,created_time,updated_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?)`,values...).
		Exec()
	if err!=nil{
		logs.Error(err)
		return err

	}
	return nil

}




//修改
func (this *LogAccessDetail) Update(id string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *LogAccessDetail) UpdateById(id uint64, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

func (this *LogAccessDetail) FindToPager(pageNo, pageSize int,
	queryFilter map[string]interface{}) (list []*ApiModel, Count int, err error) {
	var (
		c      int64
		offset int
	)
	// 获得分页的总条数
	c, err = this.filter(queryFilter).Count()
	if err != nil {
		return
	}
	offset = getOffset(pageSize, pageNo, c)

	//查询sql语句
	_, err = this.filter(queryFilter).OrderBy("-created_time").Offset(offset).Limit(pageSize).All(&list)
	if err != nil {
		return
	}
	Count = int(c)
	return
}

//多个过滤条件
func (this *LogAccessDetail) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}




type WindowStatisticModel struct {
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
	Domain      	string  `description:"domain" json:"domain"`

	PublicModel
}


func (this *LogAccessDetail) GetWindowStatistic(startTime,endTime string) (list []*WindowStatisticModel, err error) {
	o := orm.NewOrm()
	var values = make([]interface{}, 0)

	values=append(values, startTime)
	values=append(values, endTime)

	//logs.Info("---- this.TableName()-->", this.TableName())
	//  ApiIdName       this.TableName()  //  openai_detail_log_20201221
	sql :=`
SELECT  app_id,api_id ,app_name, api_name ,
year,successful_called as year_successful ,failed_called as year_failed,
month,successful_called as month_successful ,failed_called as month_failed,
day,successful_called as day_successful ,failed_called as day_failed,
hour,successful_called as hour_successful ,failed_called as hour_failed,domain
FROM (
SELECT app_id,api_id ,app_name, api_name ,year,month,day,hour,SUM(successful) as successful_called ,SUM(failed) as failed_called,domain
FROM (
SELECT app_id,api_id ,app_name, api_name ,year,month,day,HOUR,domain,
case when called_status ="200" then 1 ELSE 0 END  AS successful,
case when called_status!="200"  then 1 ELSE 0 END AS failed
FROM `+ this.TableName()+` where updated_time >= ? and updated_time < ?
) t  GROUP BY app_id,api_id ,app_name, api_name ,year,month,day,HOUR,domain
)tt `
	//logs.Info("------2------->",sql,values)

	_, err = o.Raw(sql, values).QueryRows(&list)

	if err!=nil{
		logs.Error(err)
		return nil, err
	}
	//for i, i2 := range list {
	//	logs.Info(i,"------>",i2)
	//}
	//logs.Info(list)
	return

}




