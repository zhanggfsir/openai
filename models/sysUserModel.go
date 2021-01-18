package models

import (
	"github.com/astaxie/beego/orm"
	"openai-backend/utils/conf"
	"time"
)

// 用户模型
type UserModel struct {
	Id             int    `orm:"auto;pk" description:"ID" json:"id"`
	UserId         string `description:"用户ID" json:"user_id"`
	UserName       string `description:"登录的用户名" json:"user_name"`
	Name           string `description:"姓名" json:"name"`
	Sex            string `description:"性别" json:"sex"`
	Password       string `description:"用户密码" json:"password"`
	Email          string `description:"邮箱" json:"email"`
	Phone          string `description:"手机" json:"phone"`
	Remark         string `description:"备注" json:"remark"`
	DelFlag        string `description:"删除标记" json:"del_flag"`
	PublicModel
}

// TableName 获取对应数据库表名.
func (this *UserModel) TableName() string {
	return "sys_user_info"
}

// TableEngine 获取数据使用的引擎.
func (this *UserModel) TableEngine() string {
	return "INNODB"
}

func (this *UserModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + this.TableName()
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (this *UserModel) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(this.TableNameWithPrefix())
}

func (this *UserModel) UpdateByPhone(phone string, data map[string]interface{}) error {
	data["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter("phone", phone).Update(data)

	return err
}


func (this *UserModel) FetchByUserId(userId string) (*UserModel, error) {
	err := this.QueryTable().Filter("user_id", userId).One(this)
	return this, err
}

func (this *UserModel) FetchByPhone(phone string) (*UserModel, error) {
	err := this.QueryTable().Filter("phone", phone).One(this)
	return this, err
}

// 获得用户信息
func (this *UserModel) FetchByUserName(userName string) (*UserModel, error) {
	err := this.QueryTable().Filter("username", userName).One(this)
	return this, err
}

func (this *UserModel) FindToPager(pageNo, pageSize int) (list []*UserModel, totalCount int, err error) {
	var (
		c int64
		offset int
	)

	c, err = this.QueryTable().Count()
	if err != nil {
		return
	}

	offset = getOffset(pageSize, pageNo, c)
	_, err = this.QueryTable().OrderBy("-created_time").Offset(offset).Limit(pageSize).All(&list)
	if err != nil {
		return
	}

	totalCount = int(c)
	return
}

func (this *UserModel) SaveFilter(filterKey, filterValue string, value map[string]interface{}) error {
	value["updated_time"] = time.Now()
	_, err := this.QueryTable().Filter(filterKey, filterValue).Update(value)
	return err
}
