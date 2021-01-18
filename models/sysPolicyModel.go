package models

import (
	"openai-backend/utils/casbinUtil"
	"openai-backend/utils/conf"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type PolicyModel struct {
	Id             string  `orm:"pk" description:"policy + action" json:"id"`
	PolicyGroup    string  `description:"策略组" json:"policy_group"`
	Action         string  `description:"策略action(add,read,update,delete)" json:"action"`
	Url            string  `description:"请求的路径" json:"url"`
	Remark         string  `description:"备注" json:"remark"`
	PublicModel
}

// TableName 获取对应数据库表名.
func (this *PolicyModel) TableName() string {
	return "sys_policy_info"
}

// TableEngine 获取数据使用的引擎.
func (this *PolicyModel) TableEngine() string {
	return "INNODB"
}

func (this *PolicyModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + this.TableName()
}

func NewPolicyModel() *PolicyModel {
	return &PolicyModel{}
}

func (this *PolicyModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

func (this *PolicyModel) FetchAllInfo() (list []*PolicyModel, err error) {
	_, err = this.QueryTable().All(&list)
	return
}

func (this *PolicyModel) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(this)
	return err
}

func (this *PolicyModel) SaveFilter(id string, value map[string]interface{}) error {
	value["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("id", id).Update(value)
	return err
}

func (this *PolicyModel) FetchById(id string) (*PolicyModel, error) {
	err := this.QueryTable().Filter("id", id).One(this)
	return this, err
}

// Delete
func (this *PolicyModel) DeleteById(ids []string) error {
	o := orm.NewOrm()

	if err := o.Begin(); err != nil {
		logs.Error("开启事物时出错 -> ", err)
		return err
	}

	for _, id := range ids {
		_, err := o.QueryTable(this.TableNameWithPrefix()).Filter("id", id).Delete()
		if err != nil {
			o.Rollback()
			return err
		}
	}

	return o.Commit()
}

func (this *PolicyModel) FindToPager(pageNo, pageSize int,
	queryFilter map[string]string) (list []*PolicyModel, Count int, err error) {
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

//多个过滤条件
func (this *PolicyModel) filter(queryFilter map[string]string) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			result = result.Filter(key, value)
		}
	}
	return
}

// 注册角色模型 - 初始化
func RegisterPolicy() {
	for _, policy := range casbinUtil.DefaultPolicyList {
		for _, action := range  casbinUtil.DefaultActionList {
			policyTab := &PolicyModel{
				Id: policy + action,
				PolicyGroup: policy,
				Action: action,
				Url: "/:ver/" + policy,
			}

			if _, err := NewPolicyModel().FetchById(policyTab.Id); err == orm.ErrNoRows {
				policyTab.Insert()
			}
		}
	}
}