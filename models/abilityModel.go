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
@Description:AbilityModel
*/


type AbilityModel struct {
	Id          uint64 `orm:"auto;pk" description:"ID" json:"id"`
	CategoryId  uint64  `description:"能力种类ID" json:"categoryId"`
	AbilityName string `description:"能力名称" json:"abilityName"`
	PublicModel
}

func (this *AbilityModel) TableName() string {
	return "ability_info"
}

// TableEngine 获取数据使用的引擎.
func (this *AbilityModel) TableEngine() string {
	return "INNODB"
}

func (this *AbilityModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewAbilityModel() *AbilityModel {
	return &AbilityModel{}
}

func (this *AbilityModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

//添加
func (this *AbilityModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//根据ID来删除
func (this *AbilityModel) DeleteById(o orm.Ormer, Ids []uint64) error {

	for _, id := range Ids {

		_, _ = o.QueryTable(this.TableName()).Filter("id", id).All(this)

		//删除能力
		if _, err := o.QueryTable(this.TableName()).Filter("id", this.Id).Delete(); err != nil {
			logs.Error("err=", err)
			return err
		}
		//TODO 删除能力
		//删除app_ability_info 中ability
		//删除能力api
		//删除default_quotas_info 中 api
		//删除quotas_info 中 api

	}

	return nil
}

//修改
func (this *AbilityModel) Update(id string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *AbilityModel) UpdateById(id uint64, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

func (this *AbilityModel) FindToPager(pageNo, pageSize int,
	queryFilter map[string]interface{}) (list []*AbilityModel, Count int, err error) {
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
func (this *AbilityModel) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}
