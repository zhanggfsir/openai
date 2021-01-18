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
@Description:APImodel
*/

type AbilityCategoryModelTemp struct {
	CategoryId  uint64 `orm:"auto" description:"能力种类ID" json:"categoryId"`
	CategoryName string `description:"能力种类名称" json:"categoryName"`
	AbilityList []*AbilityModelTemp `description:"AbilityList" json:"abilityList"`
}

type AbilityModelTemp struct {
	AbilityId uint64 `description:"能力ID" json:"abilityId"`
	AbilityName  string `description:"能力名称" json:"abilityName"`
	ApiList []*ApiModelTemp `description:"apiList" json:"apiList"`
}

type ApiModelTemp struct {
	Id           uint64 `orm:"auto" description:"ID" json:"id"`
	ApiName      string `description:"api名称" json:"apiName"`
	Url          string `description:"url" json:"url"`
	MaxQps       uint64 `description:"最大Qps" json:"maxQps"`
	Quotas       uint64 `description:"配额" json:"quotas"`
	QuotasPeriod string `description:"配额单位 h d m y" json:"quotasPeriod"`
	CreateTime   string `json:"createdTime"`
	ModifyTime   string `json:"updatedTime"`
}


type ApiModel struct {
	Id        uint64 `orm:"auto;pk" description:"ID" json:"id"`
	AbilityId uint64 `description:"能力ID" json:"abilityId"`
	ApiName   string `description:"api名称" json:"apiName"`
	Url       string `description:"url" json:"url"`
	UrlDesc       string `description:"url_desc" json:"urlDesc"`
	PublicModel
}

func (this *ApiModel) TableName() string {
	return "api_info"
}

// TableEngine 获取数据使用的引擎.
func (this *ApiModel) TableEngine() string {
	return "INNODB"
}

func (this *ApiModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func (this *ApiModel) QueryTable() orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	//logs.Info("-------------",this.TableNameWithPrefix())
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}



func NewApiModel() *ApiModel {
	return &ApiModel{}
}

//添加
func (this *ApiModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//根据ID来删除
func (this *ApiModel) DeleteById(o orm.Ormer, Ids []uint64) error {

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

		//删除能力api
		if _, err := o.QueryTable(this.TableName()).Filter("id", this.Id).Delete(); err != nil {
			logs.Error("err=", err)
			return err
		}

		//
		//TODO 删除default_quotas_info 中 api
		//删除quotas_info 中 api
	}

	return nil
}

//修改
func (this *ApiModel) Update(id string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *ApiModel) UpdateById(id uint64, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

func (this *ApiModel) FindToPager(pageNo, pageSize int,
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
func (this *ApiModel) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}
