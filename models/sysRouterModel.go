package models

import (
	"openai-backend/utils"
	"openai-backend/utils/conf"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

// 路由模型
type RouterModel struct {
	Id           int    `orm:"auto;pk" description:"id" json:"id"`
	Path         string `description:"请求路径" json:"path"`
	RouterGroup  string `description:"路由组" json:"routerGroup"`
	Method       string `description:"方法" json:"method"`
	Description  string `description:"路由描述" json:"description"`
	Domain       string `description:"域" json:"domain"`
	PublicModel
}

type ApiResponse struct {
	Path         string `description:"请求路径" json:"path"`
	RouterGroup  string `description:"路由组" json:"routerGroup"`
	Method       string `description:"方法" json:"method"`
	Description  string `description:"路由描述" json:"description"`
	Domain       string `description:"域" json:"domain"`
	Flag         string `description:"所有的falg" json:"flag"`
}

// TableName 获取对应数据库表名.
func (this *RouterModel) TableName() string {
	return "sys_router_info"
}

// TableEngine 获取数据使用的引擎.
func (this *RouterModel) TableEngine() string {
	return "INNODB"
}

func (this *RouterModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + this.TableName()
}

func NewRouterModel() *RouterModel {
	return &RouterModel{}
}

func (this *RouterModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

// FetchByUserId
func (this *RouterModel) FetchByPathAndMethod(path, method string) (*RouterModel, error) {
	err := this.QueryTable().Filter("path", path).Filter("method", method).One(this)
	return this, err
}

func (this *RouterModel) FetchByAppKey(appKey string) (*RouterModel, error) {
	err := this.QueryTable().Filter("app_key", appKey).One(this)
	return this, err
}

func (this *RouterModel) FetchByPhone(Pid, phone string) (*RouterModel, error) {
	err := this.QueryTable().Filter("p_user_id", Pid).Filter("phone", phone).One(this)
	return this, err
}

func (this *RouterModel) Insert() error {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		logs.Error("开启事物时出错 -> ", err)
		return err
	}

	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		o.Rollback()
		logs.Error("err=", err)
		return err
	}

	return o.Commit()
}

// SaveById 修改router信息
func (this *RouterModel) SaveById(id int, value map[string]interface{}) (err error) {
	value["updated_time"] = time.Now()
	if _, err = this.QueryTable().Filter("id", id).Update(value); err != nil {
		logs.Error(utils.Fields{
			"err": err.Error(),
			"id": id,
			"fileInfo": value,
		})
	}
	return err
}

// DeleteById 删除
func (this *RouterModel) DeleteById(ids []int) error {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		logs.Error("开启事物时出错 -> ", err)
		return err
	}

	for _, id := range ids {
		if _, err := o.QueryTable(this.TableNameWithPrefix()).Filter("id", id).Delete(); err != nil {
			logs.Error("delete router info err ", err, id)
			o.Rollback()
			return err
		}
	}

	return o.Commit()
}

func (this *RouterModel) FindToPager(pageNo, pageSize int, queryFilter map[string]string) (list []*RouterModel, Count int, err error) {
	var (
		c      int64
		offset int
	)
	// 获得分页的总条数
	if c, err = this.filter(queryFilter).Count(); err != nil {
		return
	}

	offset = getOffset(pageSize, pageNo, c)

	//查询sql语句
	if _, err = this.filter(queryFilter).OrderBy("-created_time").Offset(offset).Limit(pageSize).All(&list); err != nil {
		return
	}

	Count = int(c)
	return
}

func (this *RouterModel) filter(queryFilter map[string]string) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			result = result.Filter(key, value)
		}
	}
	return
}

// FetchAllRouters 获得所有数据
func (this *RouterModel) FetchAllRouters(queryFilter map[string]string) (list []*RouterModel, err error) {
	if _, err = this.filter(queryFilter).OrderBy("-created_time").All(&list); err != nil {
		logs.Error("err", err)
	}

	return
}