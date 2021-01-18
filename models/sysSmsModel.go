package models

import (
	"openai-backend/utils/conf"
	"github.com/astaxie/beego/orm"
)
// 用户模型
type SmsModel struct {
	Id          int `orm:"auto;pk" description:"ID" json:"id"`
	Phone       string `description:"发送手机号码" json:"phone"`
	Captcha     string `description:"验证码"  json:"captcha"`
	Type        string `description:"1:注册; 2修改密码"  json:"type"`
	Status      string `description:"发送的状态(0:fail;1:success)" json:"status"`
	PublicModel
}

// TableName 获取对应数据库表名.
func (d *SmsModel) TableName() string {
	return "sys_sms_info"
}

// TableEngine 获取数据使用的引擎.
func (d *SmsModel) TableEngine() string {
	return "INNODB"
}

func (d *SmsModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + d.TableName()
}

func NewSmsModel() *SmsModel {
	return &SmsModel{}
}

func (d *SmsModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(d.TableNameWithPrefix())
}

func (d *SmsModel) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(d)
	return err
}