package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/lisijie/goblog/util"
	"github.com/nfnt/resize"
	"goblog/models"
	"image/jpeg"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	// "time"
)

type AdminController struct {
	beego.Controller
}

func (this *AdminController) Get() {

	/*
		获取session和cookie,管理员和用户的session和cookie都要获取，并进行分类处理输出对应模板
	*/
	//普通用户
	sessionuid := this.GetSession("sessionuid")
	cookieaccount := this.Ctx.GetCookie("cookieaccount")
	//管理员
	sessionaid := this.GetSession("sessionaid")
	cookieadmin := this.Ctx.GetCookie("cookieadmin")

	//判断用户类型
	if len(cookieadmin) > 0 {
		//管理员
		this.Data["admin"] = true
		this.TplNames = "index/index.html"
		return
	} else if sessionaid != nil {
		//管理员
		this.Data["admin"] = true
		this.TplNames = "index/index.html"
	} else if len(cookieaccount) > 0 {
		//普通用户,读取用户名
		account := DecodeCookie(cookieaccount)
		userinfo, _ := models.Userinfo(account)
		this.Data["username"] = userinfo.Username
		this.Data["admin"] = false
		this.TplNames = "index/index.html"
	} else if sessionuid != nil {
		//普通用户,读取用户uid
		userinfo, _ := models.Userinfowithuid(sessionuid)
		this.Data["username"] = userinfo.Username
		this.Data["admin"] = false
		this.TplNames = "index/index.html"
	} else {
		//返回登录界面
		this.Redirect("/", 301)
	}

}

func (this *AdminController) Login() {
	/*
		后台超级管理员登录
	*/

	//cookie存在时直接登录,session存在时直接登录
	testcookie := this.Ctx.GetCookie("cookieadmin")
	testsession := this.GetSession("sessionaid")
	if len(testcookie) > 0 {
		this.Redirect("/admin", 301)
		return
	} else if testsession != nil {
		this.Redirect("/admin", 301)
		return
	} else {
		this.TplNames = "index/adminlogin.html"
		if this.GetString("dosubmit") == "yes" {
			//获取表单传递数据
			account := strings.TrimSpace(this.Input().Get("account"))
			password := strings.TrimSpace(this.Input().Get("password"))
			//remember用于判断是否创建cookie
			remember := strings.TrimSpace(this.Input().Get("remember"))

			//加密用户密码，用于判断密码是否正确
			passwordinfo := md5.New()
			salt1 := beego.AppConfig.String("salt1")
			salt2 := beego.AppConfig.String("salt2")
			io.WriteString(passwordinfo, salt1)
			io.WriteString(passwordinfo, password)
			io.WriteString(passwordinfo, salt2)
			passwordinfofinal := fmt.Sprintf("%x", passwordinfo.Sum(nil))

			//查找管理员用户名是否存在！
			searchadminaccount := models.SearchAdminAccount(account)
			searchadminaccountandpwd := models.SearchAdminAccountAndPwd(account, passwordinfofinal)
			if len(account) == 0 {
				this.Data["errmsg"] = "账号不能为空！"
			} else if len(password) == 0 {
				this.Data["errmsg"] = "密码不能为空！"
			} else if searchadminaccount != nil {
				this.Data["errmsg"] = "账号不存在！"
			} else if searchadminaccountandpwd != nil {
				this.Data["errmsg"] = "密码不正确！"
			} else {
				/*
					判断是否记录密码一周
				*/
				if remember == "yes" {
					//记住密码一周
					account1 := beego.AppConfig.String("salt1") + " " + account + " " + beego.AppConfig.String("salt2")
					account2 := []byte(account1)
					//base64加密，需要使用byte类型
					cookieval := base64.StdEncoding.EncodeToString(account2)
					this.Ctx.SetCookie("cookieadmin", cookieval, 7*24*3600, "/")
					this.Redirect("/admin", 301)
				}
				/*
					只把当前账号存到session中,登录时不记住密码一周默认设置session，不设置cookie
				*/
				admininfo, _ := models.Admininfo(account)
				aid := admininfo.Id
				sessionaid := this.GetSession("sessionaid")
				if sessionaid == nil {
					//如果session不存在，就创建session
					this.SetSession("sessionaid", aid)
					//记录最新登录时间,根据userid更新
					ip := this.Ctx.Input.IP()
					err := models.Updateadminlogintime(aid, ip)
					if err != nil {
						return
					}
					this.Redirect("/admin", 301)
				} else {
					/*
						如果session存在，判断是否是本人博客，URL中获取用户博客名称，然后从
						数据库中读取用户账号,加密后判断是否和session数值一致，一致表示是本人
						博客，显示本人博客全部功能，不一致表示是其他人博客，只显示博客内容，不显示
						博客管理功能！
					*/
					ip := this.Ctx.Input.IP()
					err := models.Updateadminlogintime(aid, ip)
					if err != nil {
						return
					}
					this.Redirect("/admin", 301)
					return
				}
			}
		}
	}

}

/*
主页面输出
*/
func (this *AdminController) Main() {
	/*
		获取session和cookie,管理员和用户的session和cookie都要获取，并进行分类处理输出对应模板
	*/
	//普通用户
	sessionuid := this.GetSession("sessionuid")
	cookieaccount := this.Ctx.GetCookie("cookieaccount")
	//管理员
	sessionaid := this.GetSession("sessionaid")
	cookieadmin := this.Ctx.GetCookie("cookieadmin")

	//判断用户类型
	if len(cookieadmin) > 0 {
		//管理员
		this.Data["hostname"], _ = os.Hostname()
		this.Data["goversion"] = runtime.Version()
		this.Data["os"] = runtime.GOOS
		this.Data["cpunum"] = runtime.NumCPU()
		this.Data["arch"] = runtime.GOARCH

		this.Data["admin"] = true
		this.TplNames = "layout.html"
		return
	} else if sessionaid != nil {
		//管理员
		this.Data["hostname"], _ = os.Hostname()
		this.Data["goversion"] = runtime.Version()
		this.Data["os"] = runtime.GOOS
		this.Data["cpunum"] = runtime.NumCPU()
		this.Data["arch"] = runtime.GOARCH

		//输出网站内容统计信息
		article0, _ := models.Articletypesumwithstatus(0)
		article1, _ := models.Articletypesumwithstatus(1)
		article2, _ := models.Articletypesumwithstatus(2)

		diary0, _ := models.Diarytypesumwithstatic(0)
		diary1, _ := models.Diarytypesumwithstatic(1)
		diary2, _ := models.Diarytypesumwithstatic(2)

		shuo, _ := models.Shuosumwithall()

		picture0, _ := models.Imagereadcountwithstatic(0)
		picture1, _ := models.Imagereadcountwithstatic(1)

		usersum, _ := models.Usersum()

		this.Data["article0"] = article0
		this.Data["article1"] = article1
		this.Data["article2"] = article2

		this.Data["diary0"] = diary0
		this.Data["diary1"] = diary1
		this.Data["diary2"] = diary2

		this.Data["shuo"] = shuo

		this.Data["picture0"] = picture0
		this.Data["picture1"] = picture1

		this.Data["usersum"] = usersum

		this.Data["admin"] = true
		this.TplNames = "layout.html"
	} else if len(cookieaccount) > 0 {
		//普通用户
		account := DecodeCookie(cookieaccount)
		userinfo, _ := models.Userinfo(account)

		//输出文章数量，日记数量，碎言碎语数量，照片数量
		userid := userinfo.Id
		article0, _ := models.Articletypesum(userid, 0)
		article1, _ := models.Articletypesum(userid, 1)
		article2, _ := models.Articletypesum(userid, 2)

		diary0, _ := models.Diarytypesum(userid, 0)
		diary1, _ := models.Diarytypesum(userid, 1)
		diary2, _ := models.Diarytypesum(userid, 2)

		shuo, _ := models.Shuosum(userid)

		picture0, _ := models.Imagereadcount(userid, 0)
		picture1, _ := models.Imagereadcount(userid, 1)

		this.Data["article0"] = article0
		this.Data["article1"] = article1
		this.Data["article2"] = article2

		this.Data["diary0"] = diary0
		this.Data["diary1"] = diary1
		this.Data["diary2"] = diary2

		this.Data["shuo"] = shuo

		this.Data["picture0"] = picture0
		this.Data["picture1"] = picture1
		this.Data["admin"] = false
		this.TplNames = "layout.html"
	} else if sessionuid != nil {
		//普通用户,读取用户uid
		userinfo, _ := models.Userinfowithuid(sessionuid)

		//输出文章数量，日记数量，碎言碎语数量，照片数量
		userid := userinfo.Id
		article0, _ := models.Articletypesum(userid, 0)
		article1, _ := models.Articletypesum(userid, 1)
		article2, _ := models.Articletypesum(userid, 2)

		diary0, _ := models.Diarytypesum(userid, 0)
		diary1, _ := models.Diarytypesum(userid, 1)
		diary2, _ := models.Diarytypesum(userid, 2)

		shuo, _ := models.Shuosum(userid)

		picture0, _ := models.Imagereadcount(userid, 0)
		picture1, _ := models.Imagereadcount(userid, 1)

		this.Data["article0"] = article0
		this.Data["article1"] = article1
		this.Data["article2"] = article2

		this.Data["diary0"] = diary0
		this.Data["diary1"] = diary1
		this.Data["diary2"] = diary2

		this.Data["shuo"] = shuo

		this.Data["picture0"] = picture0
		this.Data["picture1"] = picture1

		this.Data["admin"] = false
		this.TplNames = "layout.html"
	} else {
		//返回登录界面
		this.Redirect("/", 301)
	}
}

/*
系统设置页面
*/
func (this *AdminController) Setting() {
	this.TplNames = "system/setting.html"
}

/*
添加文章页面
*/
func (this *AdminController) Addarticle() {
	this.TplNames = "article/add.html"
}

/*
文章列表页面
*/
func (this *AdminController) Articlelist() {
	keyword := this.GetString("keyword")
	searchtype := this.GetString("searchtype")
	uid := this.GetSession("sessionuid")
	//判断点击状态，分别是已发布，草稿箱，回收站
	var status int = 0
	status, _ = this.GetInt("status")
	articlesum1, _ := models.Articletypesum(uid, 0)
	articlesum2, _ := models.Articletypesum(uid, 1)
	articlesum3, _ := models.Articletypesum(uid, 2)
	//用于显示关键字类型
	this.Data["searchtype"] = "false"
	//判断状态
	if status == 0 {
		statusid := 0
		var err error
		var articlelist []*models.Article
		//判断是否关键字搜索,并分类处理
		if len(keyword) > 0 {
			if searchtype == "title" {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 0)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
			} else if searchtype == "author" {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 1)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "author"
			} else {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 2)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "tag"

			}
		} else {
			articlelist, err = models.Articleread(uid, statusid)
		}
		if err != nil {
			return
		}

		this.TplNames = "article/list.html"
		this.Data["article"] = articlelist
		this.Data["yifabu"] = true
		this.Data["articlesum1"] = articlesum1
		this.Data["articlesum2"] = articlesum2
		this.Data["articlesum3"] = articlesum3
		return
	} else if status == 1 {
		//草稿箱
		statusid := 1
		var err error
		var articlelist []*models.Article
		//判断是否关键字搜索,并分类处理
		if len(keyword) > 0 {
			if searchtype == "title" {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 0)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
			} else if searchtype == "author" {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 1)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "author"
			} else {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 2)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "tag"
			}
		} else {
			articlelist, err = models.Articleread(uid, statusid)
		}
		if err != nil {
			return
		}
		this.TplNames = "article/list.html"
		this.Data["article"] = articlelist
		this.Data["caogaoxiang"] = true
		this.Data["articlesum1"] = articlesum1
		this.Data["articlesum2"] = articlesum2
		this.Data["articlesum3"] = articlesum3
		return
	} else {
		statusid := 2
		var err error
		var articlelist []*models.Article
		//判断是否关键字搜索,并分类处理
		if len(keyword) > 0 {
			if searchtype == "title" {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 0)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
			} else if searchtype == "author" {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 1)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "author"
			} else {
				articlelist, err = models.Articlereadwithkeyword(uid, statusid, keyword, 2)
				this.Data["guanjianzi"] = true
				this.Data["guanjianzivalue"] = keyword
				this.Data["searchtype"] = "tag"
			}
		} else {
			articlelist, err = models.Articleread(uid, statusid)
		}
		if err != nil {
			return
		}
		this.TplNames = "article/list.html"
		this.Data["article"] = articlelist
		this.Data["huishouzhan"] = true
		this.Data["articlesum1"] = articlesum1
		this.Data["articlesum2"] = articlesum2
		this.Data["articlesum3"] = articlesum3
		return
	}

}

/*
编辑文章页面
*/

func (this *AdminController) Editarticle() {
	articleid, _ := this.GetInt32("aid")
	//如果文章id存在，进入编辑页面，不存在返回后台默认页面
	if articleid > 0 {
		//判断是否存在
		articleinfo, err := models.Articlefindwithaid(articleid)
		if err != nil {
			//没有找到相应文章,返回
			this.Redirect(this.Ctx.Request.Referer(), 302)
			return
		}
		if articleinfo.UpdateTime > 0 {
			//有更新时间就显示
			this.Data["updatetime"] = true
		}
		this.Data["article"] = articleinfo
		this.TplNames = "article/edit.html"
	} else {
		this.Redirect("/admin", 301)
	}
	// this.TplNames = "article/edit.html"
}

/*
标签管理页面
*/
func (this *AdminController) Tag() {
	this.TplNames = "tag/list.html"
}

/*
添加用户页面
*/
func (this *AdminController) Adduser() {
	this.TplNames = "user/add.html"
	if this.GetString("dosubmit") == "yes" {
		account := strings.TrimSpace(this.Input().Get("account"))
		password := strings.TrimSpace(this.Input().Get("password"))
		password1 := strings.TrimSpace(this.Input().Get("password1"))
		username := strings.TrimSpace(this.Input().Get("username"))

		/*
			查找数据库是否存在该账号和博客名称
		*/
		searchaccount := models.SearchAccount(account)
		searchusername := models.SearchUsername(username)

		/*
			判断输入是否合法
		*/
		if len(account) == 0 {
			this.Data["errmsg"] = "请输入账号！"
		} else if len(password) == 0 {
			this.Data["errmsg"] = "请输入密码！"
		} else if len(password1) == 0 {
			this.Data["errmsg"] = "请确认密码！"
		} else if len(username) == 0 {
			this.Data["errmsg"] = "请输入博客名称！"
		} else if searchaccount == nil {
			this.Data["errmsg"] = "该账号已被注册！"
		} else if searchusername == nil {
			this.Data["errmsg"] = "该博客名称已被注册！"
		} else if password1 != password {
			this.Data["errmsg"] = "两次输入的密码不一致！"
		} else {
			/*
				加密注册密码
			*/
			this.Data["zhuce"] = true
			hpassword := md5.New()
			io.WriteString(hpassword, beego.AppConfig.String("salt1"))
			io.WriteString(hpassword, password)
			io.WriteString(hpassword, beego.AppConfig.String("salt2"))
			hpasswordfinal := fmt.Sprintf("%x", hpassword.Sum(nil))
			//记录ip地址
			ip := this.Ctx.Input.IP()
			err := models.Register(account, hpasswordfinal, username, ip)
			if err != nil {
				this.Ctx.WriteString("注册失败！")
			}
		}

	}
}

/*
用户列表页面
*/

func (this *AdminController) Userlist() {

	//默认status为1
	var status int
	status1 := this.GetString("status")
	if status1 == "" {
		status = 0
	} else {
		status, _ = strconv.Atoi(status1)
	}

	//获取用户信息,根据lock字段查询
	userinfo, err := models.SearchUserWithLock(status)
	if err != nil {
		return
	}

	//获取锁定和正常用户的总数
	count1, _ := models.Locksum(0)
	count2, _ := models.Locksum(1)

	this.Data["count1"] = count1
	this.Data["count2"] = count2
	this.Data["userinfo"] = userinfo
	if status == 0 {
		this.Data["lock"] = false
	} else {
		this.Data["lock"] = true
	}
	this.TplNames = "user/userlist.html"
}

/*
锁定用户和解锁用户
*/
func (this *AdminController) Lockuser() {
	lock, _ := this.GetInt("lock")
	userid := this.GetString("userid")
	userid1, _ := strconv.Atoi(userid)
	if lock == 0 {
		//解锁
		err := models.Lock(userid1, lock)
		if err != nil {
			return
		}
	} else {
		//锁定
		err := models.Lock(userid1, lock)
		if err != nil {
			return
		}
	}
	this.Redirect(this.Ctx.Request.Referer(), 301)
}

/*
编辑用户页面
*/
func (this *AdminController) Edituser() {
	this.TplNames = "user/edit.html"
}

/*
个人信息页面
*/
func (this *AdminController) Profile() {
	//判断是管理员还是普通用户
	sessionuid := this.GetSession("sessionuid")
	// cookieaccount := this.Ctx.GetCookie("cookieaccount")
	sessionaid := this.GetSession("sessionaid")
	cookieadmin := this.Ctx.GetCookie("cookieadmin")
	//管理员
	if sessionaid != nil {
		userinfo, err := models.Admininfowithaid(sessionaid)
		if err != nil {
			return
		}
		this.Data["userinfo"] = userinfo
		this.Data["user"] = false
	} else if len(cookieadmin) > 0 {

	} else if sessionuid != nil {
		//普通用户
		userinfo, err := models.Userinfowithuid(sessionuid)
		if err != nil {
			return
		}
		this.Data["userinfo"] = userinfo
		this.Data["user"] = true
	} else {

	}
	this.TplNames = "account/profile.html"
}

/*
修改个人密码
*/
func (this *AdminController) Modifypassword() {
	style, _ := this.GetInt32("type")
	oldpassword := this.Input().Get("password")
	newpassword := this.Input().Get("newpassword")
	newpassword2 := this.Input().Get("newpassword2")

	if style == 0 {
		//普通用户修改密码
		userid := this.GetSession("sessionuid")
		password := md5.New()
		salt1 := beego.AppConfig.String("salt1")
		salt2 := beego.AppConfig.String("salt2")
		io.WriteString(password, salt1)
		io.WriteString(password, oldpassword)
		io.WriteString(password, salt2)
		testpassword := fmt.Sprintf("%x", password.Sum(nil))
		//比较加密后的旧密码和数据库中密码是否一致
		userinfo, _ := models.Userinfowithuid(userid)
		if oldpassword == "" {
			this.Data["errmsg"] = "当前密码不能为空！"
			this.Profile()
		} else if newpassword == "" {
			this.Data["errmsg"] = "新密码不能为空！"
			this.Profile()
		} else if newpassword2 == "" {
			this.Data["errmsg"] = "请确认密码！"
			this.Profile()
		} else if userinfo.Password != testpassword {
			this.Data["errmsg"] = "原密码不正确！"
			this.Profile()
		} else if newpassword != newpassword2 {
			this.Data["errmsg"] = "两次输入的新密码不相同！"
			this.Profile()

		} else {
			//加密
			password1 := md5.New()
			io.WriteString(password1, salt1)
			io.WriteString(password1, newpassword)
			io.WriteString(password1, salt2)
			password1final := fmt.Sprintf("%x", password1.Sum(nil))
			err := models.Updateuserpassword(userid, password1final)
			if err != nil {
				return
			}
			this.Data["errmsg"] = "修改成功！"
		}
		this.Profile()
		// this.Redirect("/admin/modifypassword", 301)
	} else {
		//管理员修改密码
		adminid := this.GetSession("sessionaid")
		password := md5.New()
		salt1 := beego.AppConfig.String("salt1")
		salt2 := beego.AppConfig.String("salt2")
		io.WriteString(password, salt1)
		io.WriteString(password, oldpassword)
		io.WriteString(password, salt2)
		testpassword := fmt.Sprintf("%x", password.Sum(nil))
		//比较加密后的旧密码和数据库中密码是否一致
		userinfo, _ := models.Admininfowithaid(adminid)
		if oldpassword == "" {
			this.Data["errmsg"] = "当前密码不能为空！"
			this.Profile()
		} else if newpassword == "" {
			this.Data["errmsg"] = "新密码不能为空！"
			this.Profile()
		} else if newpassword2 == "" {
			this.Data["errmsg"] = "请确认密码！"
			this.Profile()
		} else if userinfo.Password != testpassword {
			this.Data["errmsg"] = "原密码不正确！"
			this.Profile()
		} else if newpassword != newpassword2 {
			this.Data["errmsg"] = "两次输入的新密码不相同！"
			this.Profile()

		} else {
			//加密
			password1 := md5.New()
			io.WriteString(password1, salt1)
			io.WriteString(password1, newpassword)
			io.WriteString(password1, salt2)
			password1final := fmt.Sprintf("%x", password1.Sum(nil))
			err := models.Updateadminpassword(adminid, password1final)
			if err != nil {
				return
			}
			this.Data["errmsg"] = "修改成功！"
		}
		this.Profile()

	}
}

/*
个人资料页面
*/
func (this *AdminController) Userinfo() {
	/*
		读取数据库信息显示到页面(读取两张表数据，分别是user表和about表)
	*/
	uid := this.GetSession("sessionuid")
	userinfo, err := models.Userinfowithuid(uid)
	if err != nil {
		this.TplNames = "message/userinfo.html"
		return
	}
	this.TplNames = "message/userinfo.html"
	this.Data["userinfo"] = userinfo
}

/*
心灵鸡汤页面
*/
func (this *AdminController) Shuo() {
	op := this.Input().Get("op")
	uid := this.GetSession("sessionuid")
	if op == "add" {
		this.TplNames = "message/editshuo.html"
		this.Data["op"] = "add"
	} else if op == "edit" {
		shuoid, _ := this.GetInt("shuoid")

		shuo, err := models.Shuofind(shuoid)
		if err != nil {
			return
		}
		this.Data["content"] = shuo.Content
		this.TplNames = "message/editshuo.html"
		this.Data["op"] = "edit"
	} else {

		// this.TplNames = "message/shuo.html"
		// //查询数据，输出
		// uid := this.GetSession("sessionuid")
		// userinfo, err := models.Userinfowithuid(uid)
		// if err != nil {
		// 	return
		// }
		// chicken, err := models.Selectshuo(userinfo.Id)
		// if err != nil {
		// 	return
		// }
		// this.Data["chicken"] = chicken
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

		this.Data["count"] = count
		this.Data["list"] = list
		this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/admin/shuo", true).ToString()
		this.TplNames = "message/shuo.html"

	}
}

/*
编辑心灵鸡汤内容
*/
func (this *AdminController) Editshuo() {
	this.TplNames = "message/editshuo.html"
	//判断是新增内容还是修改内容
	op := this.Input().Get("op")
	content := this.Input().Get("content")
	if op == "add" {
		//添加操作
		uid := this.GetSession("sessionuid")
		userinfo, err := models.Userinfowithuid(uid)
		if err != nil {
			return
		}
		//执行添加
		err = models.Addshuo(content, userinfo.Id)
		if err == nil {
			this.Redirect("/admin/shuo", 301)
		}
	} else {
		//更新操作
		uid := this.GetSession("sessionuid")
		userinfo, err := models.Userinfowithuid(uid)
		if err != nil {
			return
		}
		//执行更新
		err = models.Updateshuo(content, userinfo.Id)
		if err != nil {
			this.Redirect("/admin/shuo", 301)
		}
	}
}

/*
个人资料修改
*/

func (this *AdminController) Edituserinfo() {
	//获取表单数据

	account := strings.TrimSpace(this.Input().Get("account"))
	username := strings.TrimSpace(this.Input().Get("username"))
	sex := strings.TrimSpace(this.Input().Get("sex"))
	job := strings.TrimSpace(this.Input().Get("job"))
	brief := strings.TrimSpace(this.Input().Get("brief"))
	uid := strings.TrimSpace(this.Input().Get("uid"))
	introduce := strings.TrimSpace(this.Input().Get("introduce"))
	var sex1 int8
	if sex == "男" {
		sex1 = 0
	} else {
		sex1 = 1
	}
	//要转换成int类型
	uid1, err := strconv.ParseInt(uid, 10, 32)
	//64转32
	uid2 := int(uid1)
	user := models.User{
		Id:        uid2,
		Account:   account,
		Username:  username,
		Sex:       sex1,
		Job:       job,
		Brief:     brief,
		Introduce: introduce,
	}
	err = models.Edituserinfo(user)
	if err != nil {
		this.Redirect("/admin/userinfo", 301)
		return
	}
	this.Redirect("/admin/userinfo", 301)

}

func (this *AdminController) Headimage() {
	this.TplNames = "upload/image.html"
}

/*
上传头像处理
*/
func (this *AdminController) Postimage() {
	_, f, err := this.GetFile("uploadimage")
	if err != nil {
		this.Data["shangchuan"] = "请上传文件！"
		this.Data["errmsg"] = true
		this.TplNames = "upload/image.html"
		return
	}
	if f == nil {
		this.Data["shangchuan"] = "请上传文件！"
		this.Data["errmsg"] = true
		this.TplNames = "upload/image.html"
		return
	}
	filename := f.Filename
	houzhui := strings.LastIndex(filename, ".")
	//后缀名
	realhouzhui := filename[houzhui:]
	// this.Data["test"] = realhouzhui
	// this.TplNames = "article/testarticle.html"
	// return
	// curtime := time.Now().Unix()
	// formatcurtime := this.Timetoformat(curtime)
	//获取当前用户名
	userid := this.GetSession("sessionuid")
	if userid == nil {
		//session不存在
		return
	}
	var userinfo *models.User
	userinfo, err = models.Userinfowithuid(userid)
	if err != nil {
		return
	}
	username := userinfo.Username
	//头像保存路径
	path := "static/uploadheadimage/" + username + "/"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return
	}
	err = this.SaveToFile("uploadimage", path+"original"+realhouzhui)
	if err != nil {
		return
	}
	//图片保存成功后对图片生成缩略图，传递图片路径,图片名称，和缩略图大小和缩略图类型
	name := "original" + realhouzhui
	var size1 uint = 230
	var size2 uint = 120
	//生成大头像
	err = this.Thumbnail(path, name, size1, true)
	if err != nil {
		return
	}
	//生成小头像
	err = this.Thumbnail(path, name, size2, false)
	if err != nil {
		return
	}

	//获取上传后头像路径
	facebig := path + "facebig" + realhouzhui
	facesmall := path + "facesmall" + realhouzhui
	//先判断数据库是否存在头像路径，存在就删除头像文件，之后再更新数据库
	// if len(userinfo.Facebig) > 0 && len(userinfo.Facesmall) > 0 {
	// 	oldpath := userinfo.Facebig
	// 	oldpathmulu := oldpath[:strings.LastIndex(oldpath, "/")]
	// 	err := os.Remove(oldpathmulu)
	// 	if err != nil {
	// 		return
	// 	}
	// }
	err = models.Updateuserface(userid, facebig, facesmall)
	if err != nil {
		return
	}
	this.Data["success"] = true
	this.TplNames = "upload/image.html"
}

/*
生成图片缩略图函数
*/
func (this *AdminController) Thumbnail(path, name string, size uint, style bool) error {
	//获取图片后缀和图片完整地址
	houzhui := name[strings.LastIndex(name, "."):]
	realpath := path + name
	// 打开图片
	file, err := os.Open(realpath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(size, size, img, resize.Lanczos3)
	//如果是大头像
	if style {
		out, err := os.Create(path + "facebig" + houzhui)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		// write new image to file
		err = jpeg.Encode(out, m, nil)
		if err != nil {
			return err
		}
		return nil
	} else {
		out, err := os.Create(path + "facesmall" + houzhui)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		// write new image to file
		err = jpeg.Encode(out, m, nil)
		if err != nil {
			return err
		}
		return nil
	}

}
