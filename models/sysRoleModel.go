package models

import (
	"openai-backend/utils/conf"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type RoleModel struct {
	Id        string  `orm:"pk" description:"domain + roleId" json:"id"`
	RoleId    string  `description:"角色id" json:"role_id"`
	RoleName  string  `description:"角色名" json:"role_name"`
	Remark    string  `description:"备注" json:"remark"`
	Domain    string  `description:"域名"   json:"domain"`
	PublicModel
}

type UserRoleModel struct {
	Id        string  `orm:"pk" description:"ID" json:"id"`   // domain + userId + roleId
	UserId    string  `description:"用户Id" json:"user_id"`
	RoleId    string  `description:"角色id" json:"role_id"`
	Domain    string  `description:"域名"   json:"domain"`
	PublicModel
}

// ---------------角色表----------------------
//====================================================================
func (this *RoleModel) TableName() string {
	return "sys_role_info"
}

// TableEngine 获取数据使用的引擎.
func (this *RoleModel) TableEngine() string {
	return "INNODB"
}

func (this *RoleModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + this.TableName()
}

func NewRoleModel() *RoleModel {
	return &RoleModel{}
}


func (this *RoleModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

func (this *RoleModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

func (this *RoleModel) SaveFilter(domain, id string, value map[string]interface{}) error {
	value["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("domain", domain).Filter("id", id).Update(value)
	return err
}

func (this *RoleModel) FetchById(id string) (*RoleModel, error) {
	err := this.QueryTable().Filter("id", id).One(this)
	return this, err
}

func (this *RoleModel) FindToPager(pageNo, pageSize int,
	queryFilter map[string]string) (list []*RoleModel, Count int, err error) {
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
func (this *RoleModel) filter(queryFilter map[string]string) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			result = result.Filter(key, value)
		}
	}
	return
}

// Delete
func (this *RoleModel) DeleteById(id string) error {
	o := orm.NewOrm()
	_, err := o.QueryTable(this.TableNameWithPrefix()).Filter("id", id).Delete()

	return err
}



// ---------------人员和角色关联表----------------------
//====================================================================
func (this *UserRoleModel) TableName() string {
	return "sys_user_role"
}

// TableEngine 获取数据使用的引擎.
func (this *UserRoleModel) TableEngine() string {
	return "INNODB"
}

func (this *UserRoleModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + this.TableName()
}

func NewUserRoleModel() *UserRoleModel {
	return &UserRoleModel{}
}

func (this *UserRoleModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

func (this *UserRoleModel) Insert() error {
	o := orm.NewOrm()
	//插入到数据库
	if _, err := o.Insert(this); err != nil {
		logs.Error("err=", err)
		return err
	}
	return nil
}

func (this *UserRoleModel) FetchByUserId(userId string) (list []*UserRoleModel, err error) {
	_, err =this.QueryTable().Filter("user_id", userId).All(&list)
	return
}