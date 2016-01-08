package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lisijie/goblog/util"
	"goblog/models"
	// "math/rand"
	"strconv"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get() {
	//获取cookie和session值
	cookieaccount := this.Ctx.GetCookie("cookieaccount")
	uid := this.GetSession("sessionuid")
	/*
		如果存在cookie，判断session是否存在，不存在就创建session
	*/
	if len(cookieaccount) > 0 {
		if uid == nil {
			//获取用户账号account
			account := DecodeCookie(cookieaccount)
			//获取用户uid
			userinfo, _ := models.Userinfo(account)
			uid = userinfo.Id
			this.SetSession("sessionuid", uid)
		}
		/*
			输出文章信息
		*/
		article, _ := models.Articleread(uid, 0)
		this.Data["article"] = article

		/*
			文章推荐信息
		*/
		articletuijian, _ := models.Articletuijian()
		this.Data["articletuijian"] = articletuijian

		userinfo, err := models.Userinfowithuid(uid)
		if err != nil {
			return
		}
		this.Data["userinfo"] = userinfo
		if userinfo.Job != "" {
			this.Data["job"] = userinfo.Job
		} else {
			this.Data["job"] = "暂未填写"
		}
		if userinfo.Brief != "" {
			this.Data["brief"] = userinfo.Brief
		} else {
			this.Data["brief"] = "暂未填写"
		}
		this.TplNames = "blog/index.html"
		return
	} else if uid != nil {

		/*
			输出文章信息
		*/
		article, _ := models.Articleread(uid, 0)
		this.Data["article"] = article

		/*
			文章推荐信息
		*/
		articletuijian, _ := models.Articletuijian()
		this.Data["articletuijian"] = articletuijian

		/*
			随机文章信息
		*/
		articlerand, _ := models.Articlerand()
		this.Data["articlerand"] = articlerand

		userinfo, err := models.Userinfowithuid(uid)
		if err != nil {
			return
		}
		this.Data["userinfo"] = userinfo
		if userinfo.Job != "" {
			this.Data["job"] = userinfo.Job
		} else {
			this.Data["job"] = "暂未填写"
		}
		if userinfo.Brief != "" {
			this.Data["brief"] = userinfo.Brief
		} else {
			this.Data["brief"] = "暂未填写"
		}
		this.TplNames = "blog/index.html"
		return
	} else {
		this.Redirect("/", 301)
	}

}

/*
关于我
*/

func (this *IndexController) About() {
	//获取user表中的introduce字段值
	uid := this.GetSession("sessionuid")
	cookieaccount := this.Ctx.GetCookie("cookieaccount")
	if len(cookieaccount) > 0 {
		if uid == nil {
			//获取用户账号account
			account := DecodeCookie(cookieaccount)
			//获取用户uid
			userinfo, _ := models.Userinfo(account)
			uid = userinfo.Id
			this.SetSession("sessionuid", uid)
		}

		/*
			输出文章信息
		*/
		article, _ := models.Articleread(uid, 0)
		this.Data["article"] = article

		/*
			文章推荐信息
		*/
		articletuijian, _ := models.Articletuijian()
		this.Data["articletuijian"] = articletuijian

		/*
			随机文章信息
		*/
		articlerand, _ := models.Articlerand()
		this.Data["articlerand"] = articlerand

		//读取数据
		userinfo, err := models.Userinfowithuid(uid)
		if err != nil {
			return
		}
		this.Data["userinfo"] = userinfo
		if userinfo.Job != "" {
			this.Data["job"] = userinfo.Job
		} else {
			this.Data["job"] = "暂未填写"
		}
		if userinfo.Brief != "" {
			this.Data["brief"] = userinfo.Brief
		} else {
			this.Data["brief"] = "暂未填写"
		}
		if userinfo.Introduce != "" {
			this.Data["introduce"] = userinfo.Introduce
			this.TplNames = "blog/about.html"
		} else {
			this.Data["introduce"] = "暂无内容！"
			this.TplNames = "blog/about.html"
		}
	} else if uid != nil {

		/*
			输出文章信息
		*/
		article, _ := models.Articleread(uid, 0)
		this.Data["article"] = article

		/*
			文章推荐信息
		*/
		articletuijian, _ := models.Articletuijian()
		this.Data["articletuijian"] = articletuijian

		/*
			随机文章信息
		*/
		articlerand, _ := models.Articlerand()
		this.Data["articlerand"] = articlerand

		userinfo, err := models.Userinfowithuid(uid)
		if err != nil {
			return
		}
		this.Data["userinfo"] = userinfo
		if userinfo.Job != "" {
			this.Data["job"] = userinfo.Job
		} else {
			this.Data["job"] = "暂未填写"
		}
		if userinfo.Brief != "" {
			this.Data["brief"] = userinfo.Brief
		} else {
			this.Data["brief"] = "暂未填写"
		}
		if userinfo.Introduce != "" {
			this.Data["introduce"] = userinfo.Introduce
			this.TplNames = "blog/about.html"
		} else {
			this.Data["introduce"] = "暂无内容！"
			this.TplNames = "blog/about.html"
		}
	} else {
		this.Redirect("/", 301)
	}

}

/*
碎言碎语
*/

func (this *IndexController) Shuo() {
	//获取uid，判断是否非法访问
	cookieaccount := this.Ctx.GetCookie("cookieaccount")
	uid := this.GetSession("sessionuid")
	if len(cookieaccount) > 0 {
		if uid == nil {
			//获取用户账号account
			account := DecodeCookie(cookieaccount)
			//获取用户uid
			userinfo, _ := models.Userinfo(account)
			uid = userinfo.Id
			this.SetSession("sessionuid", uid)
		}
		//读取数据
		// userinfo, _ := models.Userinfowithuid(uid)
		// chicken, err := models.Selectshuo(userinfo.Id)
		// if err != nil {
		// 	return
		// }
		/*
			定义分页方法
		*/
		var page int
		var pagesize int = 10
		var list []*models.Chicken_soup
		var chicken models.Chicken_soup

		if page, _ = this.GetInt("page"); page < 1 {
			page = 1
		}
		offset := (page - 1) * pagesize

		count, _ := chicken.Query(uid).Count()
		if count > 0 {
			chicken.Query(uid).OrderBy("-id").Limit(pagesize, offset).All(&list)
		}

		//判断count是否为0
		if count == 0 {
			this.Data["noshuo"] = true
		} else {
			this.Data["noshuo"] = false
		}
		this.Data["count"] = count
		this.Data["list"] = list
		this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/index/shuo", true).ToString()
		this.TplNames = "blog/shuo.html"

	} else if uid != nil {
		//读取数据
		// userinfo, _ := models.Userinfowithuid(uid)
		// chicken, err := models.Selectshuo(userinfo.Id)
		// if err != nil {
		// 	return
		// }
		// this.Data["chicken"] = chicken
		// this.TplNames = "blog/shuo.html"
		/*
			定义分页方法
		*/
		var page int
		var pagesize int = 10
		var list []*models.Chicken_soup
		var chicken models.Chicken_soup

		if page, _ = this.GetInt("page"); page < 1 {
			page = 1
		}
		offset := (page - 1) * pagesize

		count, _ := chicken.Query(uid).Count()
		if count > 0 {
			chicken.Query(uid).OrderBy("-id").Limit(pagesize, offset).All(&list)
		}
		//判断count是否为0
		if count == 0 {
			this.Data["noshuo"] = true
		} else {
			this.Data["noshuo"] = false
		}
		this.Data["count"] = count
		this.Data["list"] = list
		this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/index/shuo", true).ToString()
		this.TplNames = "blog/shuo.html"
	} else {
		this.Redirect("/", 301)
	}
}

/*
个人日记
*/

func (this *IndexController) Riji() {
	uid := this.GetSession("sessionuid")
	/*
		输出文章信息
	*/
	article, _ := models.Articleread(uid, 0)
	this.Data["article"] = article

	/*
		文章推荐信息
	*/
	articletuijian, _ := models.Articletuijian()
	this.Data["articletuijian"] = articletuijian

	/*
		随机文章信息
	*/
	articlerand, _ := models.Articlerand()
	this.Data["articlerand"] = articlerand

	//获取用户信息
	userinfo, _ := models.Userinfowithuid(uid)
	if userinfo.Job != "" {
		this.Data["job"] = userinfo.Job
	} else {
		this.Data["job"] = "暂未填写"
	}
	if userinfo.Brief != "" {
		this.Data["brief"] = userinfo.Brief
	} else {
		this.Data["brief"] = "暂未填写"
	}
	//获取日记信息
	diary, _, count := models.Dirayfindwithuserid(uid)
	if count == 0 {
		this.Data["nodiary"] = true
	} else {
		this.Data["nodiary"] = false
	}
	this.Data["diary"] = diary
	this.Data["userinfo"] = userinfo
	this.TplNames = "blog/riji.html"
}

/*
文章显示，每浏览一次浏览量增加1
*/
func (this *IndexController) New() {
	uid := this.GetSession("sessionuid")

	/*
		输出文章信息
	*/
	article1, _ := models.Articleread(uid, 0)
	this.Data["article1"] = article1

	/*
		文章推荐信息
	*/
	articletuijian, _ := models.Articletuijian()
	this.Data["articletuijian"] = articletuijian

	/*
		随机文章信息
	*/
	articlerand, _ := models.Articlerand()
	this.Data["articlerand"] = articlerand

	articleid1 := this.Input().Get("articleid")
	articleid, _ := strconv.Atoi(articleid1)
	// this.Data["test"] = articleid
	// this.TplNames = "article/testarticle.html"
	// return
	//浏览量加1
	err := models.Articleviewsadd(articleid)
	if err != nil {
		return
	}
	article, _ := models.Articlefindwithaid(articleid1)
	userid := this.GetSession("sessionuid")
	var userinfo *models.User
	userinfo, err = models.Userinfowithuid(userid)
	if err != nil {
		return
	}
	if userinfo.Job != "" {
		this.Data["job"] = userinfo.Job
	} else {
		this.Data["job"] = "暂未填写"
	}
	if userinfo.Brief != "" {
		this.Data["brief"] = userinfo.Brief
	} else {
		this.Data["brief"] = "暂未填写"
	}
	this.Data["userinfo"] = userinfo
	this.Data["article"] = article
	this.TplNames = "blog/news.html"
}

/*
学无止境
*/
func (this *IndexController) Xuewuzhijing() {
	/*
		输出文章信息
	*/
	uid := this.GetSession("sessionuid")
	article, _ := models.Articleread(uid, 0)
	this.Data["article"] = article

	/*
		文章推荐信息
	*/
	articletuijian, _ := models.Articletuijian()
	this.Data["articletuijian"] = articletuijian

	/*
		随机文章信息
	*/
	articlerand, _ := models.Articlerand()
	this.Data["articlerand"] = articlerand
	this.Data["status"] = true
	userinfo, err := models.Userinfowithuid(uid)
	if err != nil {
		return
	}
	this.Data["userinfo"] = userinfo
	if userinfo.Job != "" {
		this.Data["job"] = userinfo.Job
	} else {
		this.Data["job"] = "暂未填写"
	}
	if userinfo.Brief != "" {
		this.Data["brief"] = userinfo.Brief
	} else {
		this.Data["brief"] = "暂未填写"
	}
	this.TplNames = "blog/learn.html"
}

/*
相册显示
*/
func (this *IndexController) Show() {
	userid := this.GetSession("sessionuid")
	userid1 := userid.(int)
	pictureinfo, err := models.Pictureread(userid1, 0)
	if err != nil {
		return
	}
	this.Data["pictureinfo"] = pictureinfo
	this.TplNames = "blog/xc.html"
}

/*
分类显示
*/
func (this *IndexController) Category() {
	//获取分类
	category := this.GetString("tags")
	//输出文章，显示在左侧
	article1, err := models.Articlecategory(category)
	if err != nil {
		return
	}
	this.Data["article1"] = article1
	this.Data["tags"] = category

	this.Data["status"] = false
	/*
		输出文章信息,用于右侧显示
	*/
	uid := this.GetSession("sessionuid")
	article, _ := models.Articleread(uid, 0)
	this.Data["article"] = article

	/*
		文章推荐信息
	*/
	articletuijian, _ := models.Articletuijian()
	this.Data["articletuijian"] = articletuijian

	/*
		随机文章信息
	*/
	articlerand, _ := models.Articlerand()
	this.Data["articlerand"] = articlerand

	userinfo, err := models.Userinfowithuid(uid)
	if err != nil {
		return
	}
	this.Data["userinfo"] = userinfo
	if userinfo.Job != "" {
		this.Data["job"] = userinfo.Job
	} else {
		this.Data["job"] = "暂未填写"
	}
	if userinfo.Brief != "" {
		this.Data["brief"] = userinfo.Brief
	} else {
		this.Data["brief"] = "暂未填写"
	}
	this.TplNames = "blog/learn.html"

}
