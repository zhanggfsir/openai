package models

import (
	"openai-backend/utils"
	"openai-backend/utils/casbinUtil"
	"openai-backend/utils/conf"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

// 用户模型
type ChildUserModel struct {
	Id             int    `orm:"auto;pk" description:"ID" json:"id"`
	UserId         string `description:"用户ID" json:"user_id"`
	UserName       string `description:"登录的用户名" json:"user_name"`
	Name           string `description:"姓名" json:"name"`
	Sex            string `description:"性别(0:男;1:女)" json:"sex"`
	Password       string `description:"用户密码" json:"password"`
	Email          string `description:"邮箱" json:"email"`
	Phone          string `description:"手机" json:"phone"`
	Domain         string `description:"域名" json:"domain"`
	UserType       int64  `description:"0:编程访问(appKey/appSecret),1:控制台访问(可以登陆),2:编程访问(appKey/appSecret) + 控制台访问(可以登陆)" json:"user_type"`
	UserTypeStr    string `description:"多选0:编程访问;1:控制台密码访问(0,1)" json:"user_type_str"`
	RoleIds        string `description:"角色Id列表" json:"role_ids"`
	AppKey         string `description:"app_id" json:"app_key"`
	AppSecret      string `description:"app_id" json:"app_secret"`
	ResetPd        string `description:"是否强制修改密码" json:"reset_pd"`
	Remark         string `description:"备注" json:"remark"`
	DelFlag        string `description:"删除标记" json:"del_flag"`
	PublicModel
}

//type ChildType int

const (
	// 编程访问(appKey/appSecret)
	ChildTypeProgram int64 = iota
	// 控制台访问(可以登陆)
	ChildTypeConsole
	// 编程访问(appKey/appSecret) + 控制台访问(可以登陆)
	ChildTypeAll
)


// TableName 获取对应数据库表名.
func (this *ChildUserModel) TableName() string {
	return "sys_child_info"
}

// TableEngine 获取数据使用的引擎.
func (this *ChildUserModel) TableEngine() string {
	return "INNODB"
}

func (this *ChildUserModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + this.TableName()
}

func NewChildUserModel() *ChildUserModel {
	return &ChildUserModel{}
}

func (this *ChildUserModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}


func (this *ChildUserModel) FetchByUserId(userId string) (*ChildUserModel, error) {
	err := this.QueryTable().Filter("user_id", userId).One(this)
	return this, err
}

func (this *ChildUserModel) FetchByAppKey(appKey string) (*ChildUserModel, error) {
	err := this.QueryTable().Filter("app_key", appKey).One(this)
	return this, err
}

func (this *ChildUserModel) FetchByPhone(Pid, phone string) (*ChildUserModel, error) {
	err := this.QueryTable().Filter("p_user_id", Pid).Filter("phone", phone).One(this)
	return this, err
}

// 获得用户信息
func (this *ChildUserModel) FetchByUserName(domain, userName string) (*ChildUserModel, error) {
	err := this.QueryTable().Filter("domain", domain).Filter("user_name", userName).One(this)
	return this, err
}

// 获得用户信息
func (this *ChildUserModel) FetchByUserNameDomain(domain, userName string) (*ChildUserModel, error) {
	err := this.QueryTable().Filter("domain", domain).Filter("user_name", userName).One(this)
	return this, err
}

func (this *ChildUserModel) Insert() error {
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

func (this *ChildUserModel) SaveFilter(filterKey, filterValue string, value map[string]interface{}) error {
	value["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter(filterKey, filterValue).Update(value)
	return err
}

func (this *ChildUserModel) UpdateUserRole(domain, userId, roleListStr string, addList, delList []string) error {
	var err error
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		logs.Error("开启事物时出错 -> ", err)
		return err
	}

	updateValue := map[string]interface{}{
		"role_ids": roleListStr,
		"updated_time": time.Now(),
	}
	if _, err = o.QueryTable(this.TableNameWithPrefix()).Filter("user_id", userId).Update(updateValue); err != nil {
		err = errors.New(err.Error() + this.TableNameWithPrefix())
		goto ERR
	}

	for _, roleId := range addList {
		userRoleInfo := &UserRoleModel{
			Id:          domain + userId + roleId,
			UserId:      userId,
			RoleId:      roleId,
			Domain:      domain,
			PublicModel: PublicModel{},
		}
		if _, err := o.Insert(userRoleInfo); err != nil {
			err = errors.New(err.Error() + this.TableNameWithPrefix())
			goto ERR
		}
	}
	for _, roleId := range addList {
		if err := casbinUtil.AddRolesGroupPolicy(domain, userId, roleId); err != nil {
			err = errors.New(err.Error() + this.TableNameWithPrefix())
			goto ERR
		}
	}

	for _, roleId := range delList {
		id := domain + userId + roleId
		if _, err := o.QueryTable(NewUserRoleModel().TableNameWithPrefix()).Filter("id", id).Delete(); err != nil{
			err = errors.New(err.Error() + this.TableNameWithPrefix())
			goto ERR
		}
	}

	for _, roleId := range delList {
		if err := casbinUtil.DelRolesGroupPolicy(domain, userId, roleId); err != nil {
			err = errors.New(err.Error() + this.TableNameWithPrefix())
			goto ERR
		}
	}

	return o.Commit()

ERR:
	o.Rollback()
	logs.Error("UpdateUserRole=", utils.Fields{
		"domain": domain,
		"userId": userId,
		"roleListStr":roleListStr,
		"addList": addList,
		"delList": delList,
	})

	return err
}

func (this *ChildUserModel) DeleteById(userId string) error {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		logs.Error("开启事物时出错 -> ", err)
		return err
	}

	if _, err := o.QueryTable(this.TableNameWithPrefix()).Filter("user_id", userId).Delete(); err != nil{
		logs.Error("err", err)
		o.Rollback()
		return err
	}

	userRoleInfos, err := NewUserRoleModel().FetchByUserId(userId)
	if err != nil {
		logs.Error("err", err)
		o.Rollback()
		return err
	}

	for _, userRoleInfo := range userRoleInfos {
		id := userRoleInfo.Id
		roleId := userRoleInfo.RoleId
		if _, err := o.QueryTable(NewUserRoleModel().TableNameWithPrefix()).Filter("id", id).Delete(); err != nil{
			logs.Error("err", err)
			o.Rollback()
			return err
		}
		if _, err := casbinUtil.Enforcer.RemoveFilteredGroupingPolicy(0, userId, roleId, this.Domain); err != nil {
			logs.Error("err", err)
			o.Rollback()
			return err
		}
	}

	return o.Commit()
}


func (this *ChildUserModel) FindToPager(pageNo, pageSize int,
	queryFilter map[string]string) (list []*ChildUserModel, Count int, err error) {
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

func (this *ChildUserModel) filter(queryFilter map[string]string) (result orm.QuerySeter) {
	result = this.QueryTable()
	if len(queryFilter) > 0 {
		for key, value := range queryFilter {
			result = result.Filter(key, value)
		}
	}
	return
}