package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Diary struct {
	Id         int
	Title      string
	Content    string
	Category   string
	Status     int8
	PostTime   int
	UpdateTime int
	Uid        int
	Author     string
}

/*
发表日记
*/

func Adddiary(title, content, category string, status int8, userid int) error {
	curtime1 := time.Now().Unix()
	curtime := int(curtime1)
	o := orm.NewOrm()
	userinfo, err := Userinfowithuid(userid)
	author := userinfo.Username
	diary := &Diary{
		Title:    title,
		Content:  content,
		Category: category,
		Status:   status,
		Uid:      userid,
		PostTime: curtime,
		Author:   author,
	}
	_, err = o.Insert(diary)
	if err != nil {
		return err
	}
	return nil

}

//日记编辑操作
func Diaryedit(title, category, content string, status int8, posttime, diaryid int) error {
	//当前更新时间
	curtime := int(time.Now().Unix())
	diary := &Diary{
		Title:      title,
		Category:   category,
		Content:    content,
		Status:     status,
		PostTime:   posttime,
		UpdateTime: curtime,
		Id:         diaryid,
	}
	o := orm.NewOrm()
	_, err := o.Update(diary, "Title", "Category", "Content", "Status", "PostTime", "UpdateTime", "Id")
	if err != nil {
		return err
	}
	return nil
}

//查询已发布，草稿箱，回收站日记总数
func Diarytypesum(userid, status interface{}) (int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("diary")
	count, err := qs.Filter("Uid", userid).Filter("Status", status).Count()
	if err != nil {
		return 0, err
	}
	countresult := int(count)
	return countresult, nil
}

//查询整个网站已发布，草稿箱，回收站日记总数
func Diarytypesumwithstatic(status interface{}) (int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("diary")
	count, err := qs.Filter("Status", status).Count()
	if err != nil {
		return 0, err
	}
	countresult := int(count)
	return countresult, nil
}

//读取数据,根据userid,状态,查询关键字来获取,keywordtype表示查询类型，0是标题，1是作者，2是标签
func Diaryreadwithkeyword(userid, status interface{}, keyword string, keywordtype int) ([]*Diary, error) {
	o := orm.NewOrm()
	diary := make([]*Diary, 0)
	qs := o.QueryTable("diary")
	//判断keyword类型
	var err error
	if keywordtype == 0 {
		_, err = qs.Filter("Uid", userid).Filter("Status", status).Filter("Title", keyword).OrderBy("-post_time").All(&diary)
	} else if keywordtype == 1 {
		_, err = qs.Filter("Uid", userid).Filter("Status", status).Filter("Author", keyword).OrderBy("-post_time").All(&diary)
	} else {
		//标签还需要处理
		_, err = qs.Filter("Uid", userid).Filter("Status", status).Filter("Category", keyword).OrderBy("-post_time").All(&diary)
	}
	// _, err := qs.Filter("UserId", userid).Filter("Status", status).Filter("Title", keyword).All(&article)
	if err != nil {
		return nil, err
	}
	return diary, nil
}

//读取数据,根据userid和状态来获取
func Diaryread(userid, status interface{}) ([]*Diary, error) {
	o := orm.NewOrm()
	diary := make([]*Diary, 0)
	qs := o.QueryTable("diary")
	_, err := qs.Filter("Uid", userid).Filter("Status", status).OrderBy("-post_time").All(&diary)
	if err != nil {
		return nil, err
	}
	return diary, nil
}

//批量操作
func Diarybatch(idarr []int, status int) error {
	o := orm.NewOrm()
	qs := o.QueryTable("diary")
	var err error
	//根据status值批量操作
	if status == 0 {
		//批量移到已发布
		_, err = qs.Filter("id__in", idarr).Update(orm.Params{"status": 0})
	} else if status == 1 {
		//批量移到草稿箱
		_, err = qs.Filter("id__in", idarr).Update(orm.Params{"status": 1})
	} else {
		//批量移到回收站
		_, err = qs.Filter("id__in", idarr).Update(orm.Params{"status": 2})
	}
	if err != nil {
		return err
	}
	return nil
}

//批量删除日记操作
func Diarybatchdelete(idarr []int) error {
	o := orm.NewOrm()
	qs := o.QueryTable("diary")
	_, err := qs.Filter("id__in", idarr).Delete()
	if err != nil {
		return err
	}
	return nil
}

//读取日记数据，根据日记主码来获取
func Dirayfindwithdid(diaryid interface{}) (*Diary, error) {
	o := orm.NewOrm()
	diary1 := new(Diary)
	qs := o.QueryTable("diary")
	err := qs.Filter("Id", diaryid).One(diary1)
	if err != nil {
		return nil, err
	}
	return diary1, nil
}

//删除一篇日记
func Diarydelete(diaryid int) error {
	o := orm.NewOrm()
	diary := &Diary{
		Id: diaryid,
	}
	_, err := o.Delete(diary)
	if err != nil {
		return err
	}
	return nil
}

//读取日记数据，根据日记外码uid来获取
func Dirayfindwithuserid(userid interface{}) ([]*Diary, error, int) {
	status := 0
	o := orm.NewOrm()
	diary1 := make([]*Diary, 0)
	qs := o.QueryTable("diary")
	_, err := qs.Filter("Uid", userid).Filter("Status", status).OrderBy("-post_time").All(&diary1)
	count, _ := qs.Filter("Uid", userid).Count()
	if err != nil {
		return nil, err, 0
	}

	return diary1, nil, int(count)
}
