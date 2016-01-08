package controllers

import (
	"github.com/astaxie/beego"
	"goblog/models"
	"strconv"
	"strings"
)

type DiaryController struct {
	beego.Controller
}

func (this *DiaryController) Get() {

}

/*
发表日记页面
*/
func (this *DiaryController) Diary() {
	this.Data["status"] = 0
	this.TplNames = "user/adddiary.html"
}

/*
保存日记内容，分为添加和编辑
*/
func (this *DiaryController) Save() {
	var (
		id      int    = 0
		title   string = strings.TrimSpace(this.Input().Get("title"))
		content string = this.Input().Get("content")
		tags    string = strings.TrimSpace(this.Input().Get("tags"))
		status1 int    = 0
	)
	status1, _ = this.GetInt("status")
	if title == "" {
		this.Data["errmsg"] = "标题不能为空！"
		this.Data["content"] = content
		this.Data["tags"] = tags
		this.Data["status"] = status1
		this.TplNames = "user/adddiary.html"
		return
	}
	if content == "" {
		this.Data["errmsg"] = "日记内容不能为空！"
		this.Data["title"] = title
		this.Data["tags"] = tags
		this.Data["status"] = status1
		this.TplNames = "user/adddiary.html"
		return
	}

	status := int8(status1)
	id1 := this.GetString("id")
	id2, _ := strconv.ParseInt(id1, 10, 64)
	id = int(id2)

	if id < 1 {
		//id 小于1为添加文章操作
		userid := this.GetSession("sessionuid")
		//类型断言
		uid := userid.(int)
		//发表日记
		err := models.Adddiary(title, content, tags, status, uid)
		if err != nil {
			return
		}
		this.Redirect("/diary/diarylist", 301)
	} else {
		//id大于等于1时为编辑日记操作
		//表单传递的时间格式转换为时间戳

		timevalue := this.GetString("posttime")
		timestampvalue := int(Datetotimestamp(timevalue))

		//文章编辑操作
		err := models.Diaryedit(title, tags, content, status, timestampvalue, id)
		if err != nil {
			return
		}
		this.Redirect("/diary/diarylist", 301)
	}

}

// /*
// 个人日记页面
// */
// func (this *DiaryController) Diarylist() {
// 	this.TplNames = "user/list.html"

// }

/*
日记列表页面
*/
func (this *DiaryController) Diarylist() {
	keyword := this.GetString("keyword")
	searchtype := this.GetString("searchtype")
	uid := this.GetSession("sessionuid")
	//判断点击状态，分别是已发布，草稿箱，回收站
	var status int = 0
	status, _ = this.GetInt("status")
	diarysum1, _ := models.Diarytypesum(uid, 0)
	diarysum2, _ := models.Diarytypesum(uid, 1)
	diarysum3, _ := models.Diarytypesum(uid, 2)
	//用于显示关键字类型
	this.Data["searchtype"] = "false"
	//判断状态
	if status == 0 {
		statusid := 0
		var err error
		var diarylist []*models.Diary
		//判断是否关键字搜索,并分类处理
		if len(keyword) > 0 {
			if searchtype == "title" {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 0)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
			} else if searchtype == "author" {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 1)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "author"
			} else {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 2)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "tag"

			}
		} else {
			diarylist, err = models.Diaryread(uid, statusid)
		}
		if err != nil {
			return
		}

		this.TplNames = "user/list.html"
		this.Data["diary"] = diarylist
		this.Data["yifabu"] = true
		this.Data["diarysum1"] = diarysum1
		this.Data["diarysum2"] = diarysum2
		this.Data["diarysum3"] = diarysum3
		return
	} else if status == 1 {
		//草稿箱
		statusid := 1
		var err error
		var diarylist []*models.Diary
		//判断是否关键字搜索,并分类处理
		if len(keyword) > 0 {
			if searchtype == "title" {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 0)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
			} else if searchtype == "author" {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 1)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "author"
			} else {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 2)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "tag"
			}
		} else {
			diarylist, err = models.Diaryread(uid, statusid)
		}
		if err != nil {
			return
		}
		this.TplNames = "user/list.html"
		this.Data["diary"] = diarylist
		this.Data["caogaoxiang"] = true
		this.Data["diarysum1"] = diarysum1
		this.Data["diarysum2"] = diarysum2
		this.Data["diarysum3"] = diarysum3
		return
	} else {
		statusid := 2
		var err error
		var diarylist []*models.Diary
		//判断是否关键字搜索,并分类处理
		if len(keyword) > 0 {
			if searchtype == "title" {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 0)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
			} else if searchtype == "author" {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 1)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "author"
			} else {
				diarylist, err = models.Diaryreadwithkeyword(uid, statusid, keyword, 2)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "tag"
			}
		} else {
			diarylist, err = models.Diaryread(uid, statusid)
		}
		if err != nil {
			return
		}
		this.TplNames = "user/list.html"
		this.Data["diary"] = diarylist
		this.Data["huishouzhan"] = true
		this.Data["diarysum1"] = diarysum1
		this.Data["diarysum2"] = diarysum2
		this.Data["diarysum3"] = diarysum3
		return
	}

}

//批处理日记
func (this *DiaryController) Batch() {
	ids := this.GetStrings("ids[]")
	op := this.GetString("op")
	idarr := make([]int, 0)
	for _, v := range ids {
		if id, _ := strconv.Atoi(v); id > 0 {
			idarr = append(idarr, id)
		}
	}

	switch op {
	case "topub":
		//移到已发布
		if len(idarr) == 0 {
			break
		}
		err := models.Diarybatch(idarr, 0)
		if err != nil {
			return
		}
	case "todrafts":
		//移到草稿箱
		if len(idarr) == 0 {
			break
		}
		err := models.Diarybatch(idarr, 1)
		if err != nil {
			return
		}
	case "totrash":
		//移到回收站
		if len(idarr) == 0 {
			break
		}
		err := models.Diarybatch(idarr, 2)
		if err != nil {
			return
		}
	case "delete":
		//批量删除
		if len(idarr) == 0 {
			break
		}
		err := models.Diarybatchdelete(idarr)
		if err != nil {
			return
		}
	default:
		if len(idarr) == 0 {
			break
		}
	}
	this.Redirect(this.Ctx.Request.Referer(), 302)
}

func (this *DiaryController) Editdiary() {
	diaryid, _ := this.GetInt32("did")
	//如果日记id存在，进入编辑页面，不存在返回后台默认页面
	if diaryid > 0 {
		//判断是否存在
		diaryinfo, err := models.Dirayfindwithdid(diaryid)
		if err != nil {
			//没有找到相应文章,返回
			this.Redirect(this.Ctx.Request.Referer(), 302)
			return
		}
		if diaryinfo.UpdateTime > 0 {
			//有更新时间就显示
			this.Data["updatetime"] = true
		}
		this.Data["diary"] = diaryinfo
		this.TplNames = "user/edit.html"
	} else {
		this.Redirect("/admin", 301)
	}
	// this.TplNames = "article/edit.html"
}

//删除日记
func (this *DiaryController) Delete() {
	//获取日记id
	diaryid1 := this.Input().Get("diaryid")
	status := this.Input().Get("status")
	diaryid2, err := strconv.ParseInt(diaryid1, 10, 32)
	diaryid := int(diaryid2)
	// articleid := int(articleid1)
	err = models.Diarydelete(diaryid)
	if err != nil {
		return
	}
	this.Redirect("/diary/diarylist?status="+status, 301)
}
