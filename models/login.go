package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//定义struct类型数据（数据库数据字段）

type User struct {
	Id            int
	Account       string
	Username      string
	Password      string
	Sex           int8
	Job           string
	Brief         string
	Facebig       string
	Facesmall     string
	Lock          int8
	Introduce     string
	Lastlogintime int
	Lastloginip   string
}
type Admin struct {
	Id            int
	Account       string
	Password      string
	Lastlogintime int
	Lastloginip   string
}

func RegisterDB() {
	//注册模型、驱动
	orm.RegisterModel(new(User), new(Chicken_soup), new(Admin), new(Article), new(Diary), new(Picture))
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/go_blog?charset=utf8&loc=Local", 10)
}

/*
注册账号
*/

func Register(account, password, username, ip string) error {
	//数据都不为空时
	if len(account) > 0 && len(password) > 0 && len(username) > 0 {
		//获取当前时间
		curtime64 := time.Now().Unix()
		curtime := int(curtime64)
		o := orm.NewOrm()

		user := &User{
			Account:       account,
			Password:      password,
			Username:      username,
			Lastlogintime: curtime,
			Lastloginip:   ip,
		}
		_, err := o.Insert(user)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

/*
根据account返回用户信息
*/
func Searchuser(account string) (*User, error) {
	o := orm.NewOrm()
	user := &User{}
	qs := o.QueryTable("user")
	err := qs.Filter("Account", account).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/*
查找数据库中是否存在该用户账号
*/
func SearchAccount(account string) error {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("account", account).One(user)
	//如果err不是nil,说明数据库中没有该账号！
	if err != nil {
		return err
	}
	return nil
}

/*
查找数据库中是否存在该管理员账号
*/
func SearchAdminAccount(account string) error {
	o := orm.NewOrm()
	admin := new(Admin)
	qs := o.QueryTable("admin")
	err := qs.Filter("account", account).One(admin)
	//如果err不是nil，说明数据库中没有改管理员账号
	if err != nil {
		return err
	}
	return nil
}

/*
查找数据库中是否存在该博客名称
*/

func SearchUsername(username string) error {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("username", username).One(user)
	//如果err不是nil,说明数据库中没有该博客名称！
	if err != nil {
		return err
	}
	return nil
}

/*
查找用户账号和密码是否正确,返回nil表示账号密码正确
*/

func SearchAccountAndPwd(account, password string) error {
	o := orm.NewOrm()
	user := &User{
		Account:  account,
		Password: password,
	}
	err := o.Read(user, "account", "password")
	if err != nil {
		return err
	}
	return nil
}

/*
查找管理员账号和密码是否正确，返回nil表示账号密码正确
*/
func SearchAdminAccountAndPwd(account, password string) error {
	o := orm.NewOrm()
	admin := &Admin{
		Account:  account,
		Password: password,
	}
	err := o.Read(admin, "account", "password")
	if err != nil {
		return err
	}
	return nil
}
