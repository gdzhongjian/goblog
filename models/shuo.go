package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Chicken_soup struct {
	Id      int
	Content string
	Time    int
	Uid     int
}

/**
 * 添加心灵鸡汤数据
 * @param {[type]} content string [内容]
 * @param {[type]} uid     int    [外码]
 */
func Addshuo(content string, uid int) error {
	o := orm.NewOrm()
	curtime := int(time.Now().Unix())
	chickensoup := &Chicken_soup{
		Content: content,
		Time:    curtime,
		Uid:     uid,
	}
	_, err := o.Insert(chickensoup)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 更新心灵鸡汤数据
 */
func Updateshuo(content string, uid int) error {
	//获取主码
	o := orm.NewOrm()
	curtime := int(time.Now().Unix())
	chicken_soup := new(Chicken_soup)
	qs := o.QueryTable("chicken_soup")
	err := qs.Filter("Uid", uid).One(chicken_soup)
	id := chicken_soup.Id
	newchicken := &Chicken_soup{
		Id:      id,
		Content: content,
		Time:    curtime,
		Uid:     uid,
	}
	_, err = o.Update(newchicken)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 根据外码uid来查询
 */
func Selectshuo(uid int) ([]*Chicken_soup, error) {
	o := orm.NewOrm()
	chicken := make([]*Chicken_soup, 0)
	qs := o.QueryTable("chicken_soup")
	_, err := qs.Filter("uid", uid).OrderBy("-time").All(&chicken)
	if err != nil {
		return nil, err
	}
	return chicken, nil

}

/*
查询碎言碎语数据
*/
func Shuofind(id int) (*Chicken_soup, error) {
	o := orm.NewOrm()
	chicken := new(Chicken_soup)
	qs := o.QueryTable("chicken_soup")
	err := qs.Filter("Id", id).One(chicken)
	if err != nil {
		return nil, err
	}
	return chicken, nil
}

/*
查询数据总量
*/
func (this *Chicken_soup) Query(userid interface{}) orm.QuerySeter {
	return orm.NewOrm().QueryTable(this).Filter("Uid", userid)

}

/*
查询userid用户碎言碎语总数
*/
func Shuosum(userid interface{}) (int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("chicken_soup")
	count, err := qs.Filter("Uid", userid).Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

/*
查询整个网站碎言碎语数据总数
*/
func Shuosumwithall() (int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("chicken_soup")
	count, err := qs.Count()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
