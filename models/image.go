package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Picture struct {
	Id       int
	Picture  string
	Height   int
	Uid      int
	PostTime int
	Status   int8
}

/*
插入数据
*/
func Insertpicture(userid, height int, image string) error {
	curtime := time.Now().Unix()
	curtime1 := int(curtime)
	o := orm.NewOrm()
	picture := &Picture{
		Picture:  image,
		Height:   height,
		Uid:      userid,
		PostTime: curtime1,
	}
	_, err := o.Insert(picture)
	if err != nil {
		return err
	}
	return nil
}

/*
读取图片表数据,根据uid来读取
*/
func Pictureread(userid, status int) ([]*Picture, error) {
	o := orm.NewOrm()
	picture := make([]*Picture, 0)
	qs := o.QueryTable("picture")
	_, err := qs.Filter("Uid", userid).Filter("Status", status).OrderBy("-post_time").All(&picture)
	if err != nil {
		return nil, err
	}
	return picture, nil
}

/*
删除图片
*/

func Picturedelete(pictureid int) error {
	o := orm.NewOrm()
	picture := &Picture{
		Id: pictureid,
	}
	_, err := o.Delete(picture)
	if err != nil {
		return err
	}
	return nil
}

/*
查找图片，根据pictureid来查找
*/
func Pictureselect(pictureid int) (*Picture, error) {
	o := orm.NewOrm()
	picture := new(Picture)
	qs := o.QueryTable("picture")
	err := qs.Filter("Id", pictureid).One(picture)
	if err != nil {
		return nil, err
	}
	return picture, nil
}

//批量操作
func Imagebatch(idarr []int, status int) error {
	o := orm.NewOrm()
	qs := o.QueryTable("picture")
	var err error
	//根据status值批量操作
	if status == 0 {
		//批量移到已发布
		_, err = qs.Filter("id__in", idarr).Update(orm.Params{"status": 0})
	} else {
		//批量移到回收站
		_, err = qs.Filter("id__in", idarr).Update(orm.Params{"status": 1})
	}
	if err != nil {
		return err
	}
	return nil
}

//批量删除图片操作
func Imagebatchdelete(idarr []int) error {
	o := orm.NewOrm()
	qs := o.QueryTable("picture")
	_, err := qs.Filter("id__in", idarr).Delete()
	if err != nil {
		return err
	}
	return nil
}

//批量获取图片
func Imagebatchselect(idarr []int) ([]*Picture, error) {
	o := orm.NewOrm()
	picture := make([]*Picture, 0)
	qs := o.QueryTable("picture")
	_, err := qs.Filter("id__in", idarr).All(&picture)
	if err != nil {
		return nil, err
	}
	return picture, nil
}

/*
获取图片数量，根据status来区分
*/

func Imagereadcount(userid interface{}, status1 int) (int, error) {
	status := int8(status1)
	o := orm.NewOrm()
	qs := o.QueryTable("picture")
	count, err := qs.Filter("Status", status).Filter("Uid", userid).Count()
	count1 := int(count)
	if err != nil {
		return 0, err
	}
	return count1, nil
}

/*
获取整个网站图片数量，根据status来区分
*/

func Imagereadcountwithstatic(status1 int) (int, error) {
	status := int8(status1)
	o := orm.NewOrm()
	qs := o.QueryTable("picture")
	count, err := qs.Filter("Status", status).Count()
	count1 := int(count)
	if err != nil {
		return 0, err
	}
	return count1, nil
}
