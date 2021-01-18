package models

import (
	"encoding/base64"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"openai-backend/utils"
	"openai-backend/utils/casbinUtil"
	"openai-backend/utils/conf"
	"time"
)

// 用户模型
type DomainModel struct {
	Id     int    `orm:"auto;pk" description:"ID" json:"id"`
	UserId string `description:"用户名" json:"user_id"`
	Domain string `description:"domain" json:"domain"`
	Name   string `description:"固件ID" json:"name"`
	Phone  string `description:"电话号码" json:"phone"`
	Email  string `description:"邮箱" json:"email"`
	PublicModel
}

// TableName 获取对应数据库表名.
func (this *DomainModel) TableName() string {
	return "sys_domain_info"
}

// TableEngine 获取数据使用的引擎.
func (this *DomainModel) TableEngine() string {
	return "INNODB"
}

func (this *DomainModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + this.TableName()
}

func NewDomainModel() *DomainModel {
	return &DomainModel{}
}

func (this *DomainModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

func (this *DomainModel) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(this)
	return err
}

// 获得域名信息
func (this *DomainModel) FetchByDomain(domain string) (*DomainModel, error) {
	err := this.QueryTable().Filter("domain", domain).One(this)
	return this, err
}

func (this *DomainModel) FetchByCompanyName(companyName string) (*DomainModel, error) {
	err := this.QueryTable().Filter("name", companyName).One(this)
	return this, err
}

func (this *DomainModel) FetchByMail(mail string) (*DomainModel, error) {
	err := this.QueryTable().Filter("mail", mail).One(this)
	return this, err
}

func (this *DomainModel) FetchByPhone(phone string) (*DomainModel, error) {
	err := this.QueryTable().Filter("phone", phone).One(this)
	return this, err
}

func (this *DomainModel) FetchByUserId(userId string) (*DomainModel, error) {
	err := this.QueryTable().Filter("user_id", userId).One(this)
	return this, err
}

func (this *DomainModel) SaveFilter(domain string, value map[string]interface{}) error {
	value["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("domain", domain).Update(value)
	return err
}

//域名注册,使用事务
func (this *DomainModel) DomainRegister(password, roleId, userName string) error {
	var err error
	o := orm.NewOrm()
	domain := this.Domain
	phone := this.Phone
	userId := this.UserId

	if err = o.Begin(); err != nil {
		logs.Error("开启事物时出错 -> ", err)
		return err
	}

	//1.首先插入域表
	_, err = o.Insert(this)
	if err != nil {
		o.Rollback()
		goto ERR
	}

	//2.首先在user表中添加用户
	if err = addUser(o, userId, userName, password, phone); err != nil {
		o.Rollback()
		goto ERR
	}

	//3.添加admin角色表
	if err = addRole(o, domain, roleId); err != nil {
		o.Rollback()
		goto ERR
	}

	//4.添加user admin角色表
	if err = addUserRole(o, domain, userId, roleId); err != nil {
		o.Rollback()
		goto ERR
	}

	//5.添加casbin
	if err = AddPolicyFromController(domain, roleId); err != nil {
		o.Rollback()
		goto ERR
	}

	//6.在casbin添加权限和用户
	if err = casbinUtil.AddRolesGroupPolicy(domain, userId, roleId); err != nil {
		o.Rollback()
		goto ERR
	}

	return o.Commit()

ERR:
	logs.Error("register domain err", utils.Fields{
		"phone":  phone,
		"domain": domain,
		"userId": userId,
		"err":    err.Error(),
	})

	return err
}

func addUser(o orm.Ormer, userId, userName, password, phone string) error {
	decodePd, _ := base64.StdEncoding.DecodeString(password)
	userStruct := &UserModel{
		UserId:   userId,
		UserName: userName,
		Name:     userName,
		Phone:    phone,
		Password: utils.GenPassword(string(decodePd)),
		DelFlag:  "0",
	}
	_, err := o.Insert(userStruct)
	return err
}

func addRole(o orm.Ormer, domain, roleId string) error {
	roleModel := &RoleModel{
		Id:       domain + roleId,
		RoleId:   roleId,
		RoleName: roleId,
		Domain:   domain,
	}

	_, err := o.Insert(roleModel)
	return err
}

func addUserRole(o orm.Ormer, domain, userId, roleId string) error {
	userRoleModel := &UserRoleModel{
		Id:     domain + userId + roleId,
		UserId: userId,
		RoleId: roleId,
		Domain: domain,
	}

	_, err := o.Insert(userRoleModel)
	return err
}

func AddPolicyFromController(domain, roleId string) error {
	var (
		err  error
		list []*RouterModel
	)

	if list, err = NewRouterModel().FetchAllRouters(map[string]string{}); err != nil {
		logs.Error("err=", err)
		return err
	}

	for _, routerInfo := range list {
		path := routerInfo.Path
		method := routerInfo.Method
		if flag, err := casbinUtil.Enforcer.AddPolicy(roleId, domain, path, method); err != nil && !flag {
			logs.Error("fetch all policy", utils.Fields{
				"err":    err,
				"flag":   flag,
				"domain": domain,
				"role":   roleId,
				"path":   path,
			})

			return errors.New("add policy err")
		}
	}
	return nil
}
