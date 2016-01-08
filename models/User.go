package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

/*
根据account字段获取用户表信息,返回一个user结构体，主要用于登录注册判断
*/

func Userinfo(account interface{}) (*User, error) {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("account", account).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/*
根据account字段获取管理员表信息，返回一个admin结构体
*/
func Admininfo(account interface{}) (*Admin, error) {
	o := orm.NewOrm()
	admin := new(Admin)
	qs := o.QueryTable("admin")
	err := qs.Filter("account", account).One(admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

/*
根据aid字段获取用户表信息，返回一个admin结构体
*/
func Admininfowithaid(aid interface{}) (*Admin, error) {
	o := orm.NewOrm()
	admin := new(Admin)
	qs := o.QueryTable("admin")
	err := qs.Filter("id", aid).One(admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

/*
根据uid字段获取用户表信息，返回一个user结构体
*/
func Userinfowithuid(uid interface{}) (*User, error) {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("id", uid).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/*
修改个人资料
*/
func Edituserinfo(user User) error {
	o := orm.NewOrm()
	//要新建一个User类型的struct，不然数据会被写入
	userinfo := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("Id", user.Id).One(userinfo)
	if err == nil {
		userinfo.Account = user.Account
		userinfo.Username = user.Username
		userinfo.Sex = user.Sex
		if user.Job != "" {
			userinfo.Job = user.Job
		}
		if user.Brief != "" {
			userinfo.Brief = user.Brief
		}
		if user.Introduce != "" {
			userinfo.Introduce = user.Introduce
		}

		_, err := o.Update(userinfo) //要写入地址,切记
		if err != nil {
			return err
		}
	}

	return nil
}

//更新普通用户登录时间信息和IP信息
func Updateuserlogintime(userid interface{}, ip string) error {
	curtime1 := time.Now().Unix()
	curtime := int(curtime1)
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	_, err := qs.Filter("Id", userid).Update(orm.Params{"lastlogintime": curtime, "lastloginip": ip})
	if err != nil {
		return err
	}
	return nil
}

//更新管理员登录时间信息和IP信息
func Updateadminlogintime(adminid interface{}, ip string) error {
	curtime1 := time.Now().Unix()
	curtime := int(curtime1)
	o := orm.NewOrm()
	qs := o.QueryTable("admin")
	_, err := qs.Filter("Id", adminid).Update(orm.Params{"lastlogintime": curtime, "lastloginip": ip})
	if err != nil {
		return err
	}
	return nil
}

//

//更新用户密码
func Updateuserpassword(userid interface{}, password string) error {
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	_, err := qs.Filter("Id", userid).Update(orm.Params{"password": password})
	if err != nil {
		return err
	}
	return nil
}

//更新管理员密码
func Updateadminpassword(adminid interface{}, password string) error {
	o := orm.NewOrm()
	qs := o.QueryTable("admin")
	_, err := qs.Filter("Id", adminid).Update(orm.Params{"password": password})
	if err != nil {
		return err
	}
	return nil
}

//更新普通用户头像信息
func Updateuserface(userid interface{}, facebig, facesmall string) error {
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	_, err := qs.Filter("Id", userid).Update(orm.Params{"facebig": facebig, "facesmall": facesmall})
	if err != nil {
		return err
	}
	return nil
}

//查询整个网站用户总数
func Usersum() (int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	count, err := qs.Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

/*
根据lock字段查询用户
*/

func SearchUserWithLock(lock int) ([]*User, error) {
	lock1 := int8(lock)
	o := orm.NewOrm()
	user := make([]*User, 0)
	qs := o.QueryTable("user")
	_, err := qs.Filter("Lock", lock1).All(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/*
解锁，锁定用户
*/
func Lock(userid, lock int) error {
	lock1 := int8(lock)
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	_, err := qs.Filter("Id", userid).Update(orm.Params{"lock": lock1})
	if err != nil {
		return err
	}
	return nil
}

/*
锁定和正常用户总数
*/
func Locksum(lock int) (int, error) {
	lock1 := int8(lock)
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	count, err := qs.Filter("Lock", lock1).Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
