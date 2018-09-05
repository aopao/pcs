package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"uugu.org/pcs/common"
	"time"
)

type User struct {
	Id         int64     `json:"id"`
	UserName   string    `json:"userName"`
	Name       string    `json:"name"`
	Password   string    `json:"-"`
	Status     int       `json:"status"`
	Tel        string    `json:"tel"`
	Email      string    `json:"email"`
	Qq         string    `json:"qq"`
	Weixin     string    `json:"weixin"`    //`json:"description,omitempty"`
	ParentTel  string    `json:"parentTel"` //`json:"description,omitempty"`
	Address    string    `json:"address"`   //`json:"description,omitempty"`
	School     string    `json:"school"`    //`json:"description,omitempty"`
	Province   string    `json:"province"`  //`json:"description,omitempty"`
	City       string    `json:"city"`      //`json:"description,omitempty"`
	Country    string    `json:"country"`   //`json:"description,omitempty"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Valid      int       `json:"valid"`
}

func (m *User) TableName() string {
	return TableName("pcs_users")
}

func GetUserById(id int64) (*User, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM User WHERE id IN (?, ?, ?)", ids)
	user := User{Id: id}

	//err := o.QueryTable("user").Filter("name", "slene").One(&user)

	ormerr := o.Read(&user)

	if ormerr == orm.ErrNoRows {
		logs.Error("查询不到,主键:", user.Id)
	} else if ormerr == orm.ErrMissPK {
		logs.Error("找不到主键:", user.Id)
	} else {
		return &user, nil
	}
	return nil, ormerr
}

func GetUserByUserName(userName string) (*User, error) {
	var user User
	err := o.QueryTable(new(User).TableName()).Filter("user_name", userName).One(&user)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		logs.Error("Returned Multi Rows Not One")
		return nil, err
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		logs.Error("Not row found")
		return nil, err
	}
	return &user, nil
}

func GetUsers(user *User) ([]User, error) {
	//o := orm.NewOrm()
	// o.Raw("aa")
	// ids := []int{1, 2, 3}
	// o.Raw("SELECT * FROM user WHERE id IN (?, ?, ?)", ids)
	var users []User
	//num, err := o.Raw("SELECT * FROM ? as bc WHERE bc.key = ?", param.TableName(), param.Key).QueryRows(&users)
	num, err := o.QueryTable(user.TableName()).Filter("user_name", user.UserName).Filter("password", user.Password).All(&users)
	if err != nil {
		return nil, err
	}
	logs.Info("User nums: ", num)
	//users := []*User{}
	//o.QueryTable(new(User).TableName()).All(&users)
	return users, nil
}

func GetAllUsers() ([]User, error) {
	//o := orm.NewOrm()
	var users []User
	//num, err := o.Raw("SELECT * FROM ? as bc WHERE bc.key = ?", param.TableName(), param.Key).QueryRows(&users)
	num, err := o.QueryTable(new(User).TableName()).All(&users)
	if err != nil {
		return nil, err
	}
	logs.Info("Get User nums: ", num)
	//users := []*User{}
	//o.QueryTable(new(User).TableName()).All(&users)
	return users, nil
}

func AddUser(param *User) (*common.CommonResponse) {
	commonResponse := common.CommonResponse{}
	o := orm.NewOrm()
	result, ormerr := o.Insert(param)
	if nil != ormerr {
		logs.Error("Add User Error:", ormerr)
		return commonResponse.ToFail(ormerr.Error())
	} else {
		return commonResponse.ToSuccess("Add User Success! Id:" + string(result))
	}
}
