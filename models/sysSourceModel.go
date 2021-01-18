package models

import (
	"openai-backend/utils/conf"
	"github.com/astaxie/beego/orm"
	"time"
)

type SourceModel struct {
	Id          string  `orm:"pk" description:"domain + source_id" json:"id"`
	SourceId    string  `description:"资源Id" json:"source_id"`
	SourceName  string  `description:"资源Id" json:"source_name"`
	Remark      string  `description:"备注" json:"remark"`
	Domain      string  `description:"域名" json:"domain"`
	PublicModel
}

// TableName 获取对应数据库表名.
func (this *SourceModel) TableName() string {
	return "sys_source_info"
}

// TableEngine 获取数据使用的引擎.
func (this *SourceModel) TableEngine() string {
	return "INNODB"
}

func (this *SourceModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + this.TableName()
}

func NewSourceModel() *SourceModel {
	return &SourceModel{}
}

func (this *SourceModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

func (this *SourceModel) FetchAllInfo() (list []*SourceModel, err error) {
	_, err = this.QueryTable().All(&list)
	return
}

func (this *SourceModel) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(this)
	return err
}

func (d *SourceModel) SaveFilter(domain, id string, value map[string]interface{}) error {
	value["updated_time"] = time.Now()
	_, err := d.QueryTable().Filter("domain", domain).Filter("id", id).Update(value)
	return err
}

func (this *SourceModel) FetchById(domain, id string) (*SourceModel, error) {
	err := this.QueryTable().Filter("domain", domain).Filter("id", id).One(this)
	return this, err
}
