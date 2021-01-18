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
@Description:AbilityCategorymodel
*/

type AbilityCategoryModel struct {
	Id           uint64 `orm:"auto;pk" description:"ID" json:"id"`
	CategoryName string `description:"能力种类名称" json:"categoryName"`
	PublicModel
}
//SELECT category_id,category_name,SUM(successful_called) as successful_called ,SUM(failed_called)  as failed_called

type CategoryModelTemp struct {
	CategoryId           uint64 ` description:"ID" json:"category_id"`
	CategoryName 	   	 string `description:"能力种类名称" json:"category_name"`
	SuccessfulCalled     uint64 `description:"apiId" json:"successful_called"`
	FailedCalled         uint64 `description:"apiId" json:"failed_called"`
}



func (this *AbilityCategoryModel) TableName() string {
	return "ability_category_info"
}

// TableEngine 获取数据使用的引擎.
func (this *AbilityCategoryModel) TableEngine() string {
	return "INNODB"
}

func (this *AbilityCategoryModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func (this *AbilityCategoryModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

func NewAbilityCategoryModel() *AbilityCategoryModel {
	return &AbilityCategoryModel{}
}

//添加
func (this *AbilityCategoryModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//根据ID来删除
func (this *AbilityCategoryModel) DeleteById(o orm.Ormer, Ids []uint64) error {

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

		//删除能力类型
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
func (this *AbilityCategoryModel) Update(id string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *AbilityCategoryModel) UpdateById(id uint64, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

func (this *AbilityCategoryModel) AbilityCategoryList(
	queryFilter map[string]interface{}) (list []*AbilityCategoryModel, Count int, err error) {
	var c int64

	c, err = this.filter(queryFilter).Count()

	//查询sql语句
	_, err = this.filter(queryFilter).OrderBy("-created_time").All(&list)
	if err != nil {
		return
	}
	Count = int(c)
	return
}

func (this *AbilityCategoryModel) FindToPager(pageNo, pageSize int,
	queryFilter map[string]interface{}) (list []*AbilityCategoryModel, Count int, err error) {
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
func (this *AbilityCategoryModel) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}





