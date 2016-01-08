package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Article struct {
	Id         int
	UserId     int
	Author     string
	Title      string
	Color      string
	Content    string
	Tags       string
	PostTime   int
	Views      int
	Status     int8
	UpdateTime int
	IsTop      int8
	Picture    string
	Brief      string
}

//插入数据
// func Articleadd(title string, color string, content string) error {
// 	o := orm.NewOrm()
// 	article := &Article{
// 		// Author:  author,
// 		Title:   title,
// 		Color:   color,
// 		Content: content,
// 	}
// 	_, err := o.Insert(article)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

//添加文章
func Articleadd(author, title, content, tags, color, picture, brief string, status, istop int8) (articleid int, err error) {
	if len(picture) == 0 {
		picture = ""
	}
	curtime := int(time.Now().Unix())
	o := orm.NewOrm()
	//根据author来获取uid
	userinfo, _ := Userinfo(author)
	userid := userinfo.Id
	article := &Article{
		UserId:   userid,
		Author:   author,
		Title:    title,
		Content:  content,
		Tags:     tags,
		Color:    color,
		Status:   status,
		IsTop:    istop,
		PostTime: curtime,
		Picture:  picture,
		Brief:    brief,
	}
	id, err := o.Insert(article)

	if err != nil {
		return 0, err
	}
	return int(id), nil

}

//文章编辑操作
func Articleedit(title, color, tags, content, brief string, status, istop int8, posttime, articleid int) error {
	//当前更新时间
	curtime := int(time.Now().Unix())
	article := &Article{
		Title:      title,
		Color:      color,
		Tags:       tags,
		Content:    content,
		Status:     status,
		IsTop:      istop,
		PostTime:   posttime,
		UpdateTime: curtime,
		Id:         articleid,
		Brief:      brief,
	}
	o := orm.NewOrm()
	_, err := o.Update(article, "Title", "Color", "Tags", "Content", "Status", "IsTop", "PostTime", "UpdateTime", "Id", "Brief")
	// qs := o.QueryTable("article")
	// _, err := qs.Filter("Id", articleid).Update(article)
	if err != nil {
		return err
	}
	return nil
}

//读取数据,根据userid和状态来获取
func Articleread(userid, status interface{}) ([]*Article, error) {
	o := orm.NewOrm()
	article := make([]*Article, 0)
	qs := o.QueryTable("article")
	_, err := qs.Filter("UserId", userid).Filter("Status", status).OrderBy("-post_time").All(&article)
	if err != nil {
		return nil, err
	}
	return article, nil
}

//精心推荐文章，根据浏览量来查询前五条

//读取数据,根据userid,状态,查询关键字来获取,keywordtype表示查询类型，0是标题，1是作者，2是标签
func Articlereadwithkeyword(userid, status interface{}, keyword string, keywordtype int) ([]*Article, error) {
	o := orm.NewOrm()
	article := make([]*Article, 0)
	qs := o.QueryTable("article")
	//判断keyword类型
	var err error
	if keywordtype == 0 {
		_, err = qs.Filter("UserId", userid).Filter("Status", status).Filter("Title", keyword).OrderBy("-post_time").All(&article)
	} else if keywordtype == 1 {
		_, err = qs.Filter("UserId", userid).Filter("Status", status).Filter("Author", keyword).OrderBy("-post_time").All(&article)
	} else {
		//标签还需要处理
		_, err = qs.Filter("UserId", userid).Filter("Status", status).Filter("Tags", keyword).OrderBy("-post_time").All(&article)
	}
	// _, err := qs.Filter("UserId", userid).Filter("Status", status).Filter("Title", keyword).All(&article)
	if err != nil {
		return nil, err
	}
	return article, nil
}

//读取文章数据，根据文章主码来获取
func Articlefindwithaid(articleid interface{}) (*Article, error) {
	o := orm.NewOrm()
	article1 := new(Article)
	qs := o.QueryTable("article")
	err := qs.Filter("Id", articleid).One(article1)
	if err != nil {
		return nil, err
	}
	return article1, nil
}

//查询已发布，草稿箱，回收站文章总数
func Articletypesum(userid, status interface{}) (int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	count, err := qs.Filter("Userid", userid).Filter("Status", status).Count()
	if err != nil {
		return 0, err
	}
	countresult := int(count)
	return countresult, nil
}

//查询网站中已发布，草稿箱，回收站文章总数
func Articletypesumwithstatus(status interface{}) (int, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	count, err := qs.Filter("Status", status).Count()
	if err != nil {
		return 0, err
	}
	countresult := int(count)
	return countresult, nil
}

//删除一篇文章
func Articledelete(articleid int) error {
	o := orm.NewOrm()
	article := &Article{
		Id: articleid,
	}
	_, err := o.Delete(article)
	if err != nil {
		return err
	}
	return nil
}

//批量操作
func Articlebatch(idarr []int, status int) error {
	o := orm.NewOrm()
	qs := o.QueryTable("article")
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

//批量删除文章操作
func Articlebatchdelete(idarr []int) error {
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	_, err := qs.Filter("id__in", idarr).Delete()
	if err != nil {
		return err
	}
	return nil
}

//文章搜索功能，标题，作者，标签进行搜索，实现精确查找
func Articleselect(keyword string, searchtype string) ([]*Article, error) {
	o := orm.NewOrm()
	article := make([]*Article, 0)
	if searchtype == "title" {
		//按标题关键字搜索
		qs := o.QueryTable("article")
		_, err := qs.Filter("Title", keyword).All(&article)
		if err != nil {
			return nil, err
		}
		return article, nil
	}
	return nil, nil
}

//文章封面图片更新
func Articleimageupdate(articleid int, picture string) error {

	article := &Article{
		Picture: picture,
		Id:      articleid,
	}
	o := orm.NewOrm()
	_, err := o.Update(article, "Picture")
	// qs := o.QueryTable("article")
	// _, err := qs.Filter("Id", articleid).Update(article)
	if err != nil {
		return err
	}
	return nil
}

//文章浏览量加1
func Articleviewsadd(articleid int) error {
	o := orm.NewOrm()
	article := &Article{
		Id: articleid,
	}
	qs := o.QueryTable("article")
	err := qs.Filter("Id", articleid).One(article)
	if err != nil {
		return err
	}
	views := article.Views + 1
	_, err = qs.Filter("Id", articleid).Update(orm.Params{"Views": views})
	if err != nil {
		return err
	}
	return nil
}

/*
精心推荐博客功能，主要根据点击量进行排序输出
*/
func Articletuijian() ([]*Article, error) {
	o := orm.NewOrm()
	article := make([]*Article, 0)
	qs := o.QueryTable("article")
	_, err := qs.OrderBy("-views").Limit(5, 0).All(&article)
	if err != nil {
		return nil, err
	}
	return article, nil
}

/*
随机文章
*/
func Articlerand() (result []Article, err error) {
	o := orm.NewOrm()
	// var r orm.RawSeter
	_, err = o.Raw("SELECT * FROM article ORDER BY rand() LIMIT 5").QueryRows(&result)

	return result, err
}

/*
分类文章
*/
func Articlecategory(category string) ([]*Article, error) {
	o := orm.NewOrm()
	article := make([]*Article, 0)
	qs := o.QueryTable("article")
	_, err := qs.Filter("Tags", category).OrderBy("-post_time").All(&article)
	if err != nil {
		return nil, err
	}
	return article, nil
}
