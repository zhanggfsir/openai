package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

/*
@Author:
@Time: 2020-11-05 15:02
@Description:quotasmodel
*/

type QuotasModel struct {
	Id           uint64 `orm:"auto" description:"ID" json:"id"`
	AppId        uint64 `description:"appId" json:"appId"`
	ApiId        uint64 `description:"apiId" json:"apiId"`
	MaxQps       uint64 `description:"最大Qps" json:"maxQps"`
	Quotas       uint64 `description:"配额" json:"maxQps"`
	QuotasPeriod string `description:"单位 h d m y" json:"quotasPeriod"`
	Status       string `description:"api状态 0 免费 1 付费" json:"status"`
	PublicModel
}

func (this *QuotasModel) TableName() string {
	return "quotas_info"
}

// TableEngine 获取数据使用的引擎.
func (this *QuotasModel) TableEngine() string {
	return "INNODB"
}

func (this *QuotasModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewQuotasModel() *QuotasModel {
	return &QuotasModel{}
}

func (this *QuotasModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

//添加
func (this *QuotasModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//根据ID来删除
func (this *QuotasModel) DeleteById(o orm.Ormer, Ids []uint64) error {

	for _, id := range Ids {

		_, err := o.QueryTable(this.TableName()).Filter("id", id).All(this)
		if err != nil {
			if err := o.Rollback(); err != nil {
				logs.Error("err=", err)
				return err
			}
			logs.Error("err=", err)
			return err
		}

		//删除
		if _, err := o.QueryTable(this.TableName()).Filter("id", this.Id).Delete(); err != nil {
			logs.Error("err=", err)
			return err
		}

	}

	return nil
}


//修改
func (this *QuotasModel) Update(id string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *QuotasModel) UpdateById(id uint64, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

func (this *QuotasModel) FindToPager(pageNo, pageSize int,
	queryFilter map[string]interface{}) (list []*AppModel, Count int, err error) {
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
func (this *QuotasModel) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}

//
////多个过滤条件
//func (this *QuotasModel) Filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
//	return this.filter(queryFilter)
//	//result = this.QueryTable()
//	//if len(queryFilter) > 0 {
//	//	for key, value := range queryFilter {
//	//		//result.Filter(key, value).Filter(key, value).Filter(key, value)
//	//		result = result.Filter(key, value)
//	//	}
//	//}
//	//return
//}
