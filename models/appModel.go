package models

/*
@Author:
@Time: 2020-11-05 15:02
@Description:APPmodel
*/
import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type AppsAbilityModel struct {
	AbilityId   uint64          `description:"能力ID" json:"abilityId"`
	AbilityName string          `description:"能力名称" json:"abilityName"`
	ApiList     []*AppsApiModel `description:"apiList" json:"apiList"`
}

type AppsApiModel struct {
	ApiId        uint64 `orm:"auto" description:"ID" json:"appId"`
	ApiName      string `description:"api名称" json:"apiName"`
	Url          string `description:"url" json:"url"`
	MaxQps       uint64 `description:"最大Qps" json:"maxQps"`
	Quotas       uint64 `description:"配额" json:"quotas"`
	QuotasPeriod string `description:"配额单位 h d m y" json:"quotasPeriod"`
}

type AppModelTemp struct {
	Id uint64 `orm:"auto" description:"ID" json:"id"`
	//AppId       string   `description:"设备ID" json:"appId"`
	AppName     string   `description:"App名称" json:"appName"`
	AppType     string   `description:"类型" json:"appType"`
	ApiKey      string   `description:"apiKey" json:"apiKey"`
	SecretKey   string   `description:"SecretKey" json:"secretKey"`
	AbilityList []uint64 `description:"api列表" json:"abilityList"`
	Desc        string   `description:"备注" json:"desc"`
	Domain      string   `description:"所属域" json:"domain"`
	CreateTime  string   `json:"createdTime"`
	ModifyTime  string   `json:"updatedTime"`
}

type AppModel struct {
	Id        uint64 `orm:"auto" description:"ID" json:"id"`
	//AppId     string `description:"AppID" json:"appId"`
	AppName   string `description:"App名称" json:"appName"`
	ApiKey    string `description:"apiKey 24位" json:"apiKey"`
	SecretKey string `description:"SecretKey 32位" json:"secretKey"`
	AppType   string `description:"类型" json:"appType"`
	Desc      string `description:"备注" json:"desc"`
	Domain    string `description:"所属域" json:"domain"`
	PublicModel
}

func (this *AppModel) TableName() string {
	return "app_info"
}

// TableEngine 获取数据使用的引擎.
func (this *AppModel) TableEngine() string {
	return "INNODB"
}

func (this *AppModel) TableNameWithPrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "") + this.TableName()
}

func NewAppModel() *AppModel {
	return &AppModel{}
}

func (this *AppModel) QueryTable()  orm.QuerySeter {
	//return orm.NewOrm().QueryTable(this.TableNameWithPrefix()).Filter("del_flag", "0")
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

//添加
func (this *AppModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

//根据ID来删除
func (this *AppModel) DeleteById(o orm.Ormer, Ids []uint64) error {

	for _, id := range Ids {

		_, err := o.QueryTable(this.TableName()).Filter("id", id).All(this)
		if err != nil {
			_ = o.Rollback()
			logs.Error("err=", err)
			return err
		}
		//删除quotas_info 中app
		if _, err := o.QueryTable(NewQuotasModel().TableName()).Filter("app_id", this.Id).Delete(); err != nil {
			_ = o.Rollback()
			logs.Error("err=", err)
			return err
		}
		//删除app_ability_info 中 app
		if _, err := o.QueryTable(NewAppAbilityModel().TableName()).Filter("app_id", this.Id).Delete(); err != nil {
			_ = o.Rollback()
			logs.Error("err=", err)
			return err
		}

		//删除App
		if _, err := o.QueryTable(this.TableName()).Filter("id", this.Id).Delete(); err != nil {
			_ = o.Rollback()
			logs.Error("err=", err)
			return err
		}
		//TODO 删除redis缓存



	}

	return nil
}

//修改
func (this *AppModel) Update(id string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)
	return err
}

func (this *AppModel) UpdateById(id uint64, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(data)

	return err
}

func (this *AppModel) FindToPager(domain string,pageNo, pageSize int,
	queryFilter map[string]interface{}) (list []*AppModel, Count int, err error) {
	var (
		c      int64
		offset int
	)
	// 获得分页的总条数
	c, err = this.filter(queryFilter).Filter("domain",domain).Count()
	if err != nil {
		return
	}
	offset = getOffset(pageSize, pageNo, c)

	//查询sql语句
	_, err = this.filter(queryFilter).Filter("domain",domain).OrderBy("-created_time").Offset(offset).Limit(pageSize).All(&list)
	if err != nil {
		return
	}
	Count = int(c)
	return
}

//多个过滤条件
func (this *AppModel) filter(queryFilter map[string]interface{}) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			//result.Filter(key, value).Filter(key, value).Filter(key, value)
			result = result.Filter(key, value)
		}
	}
	return
}
