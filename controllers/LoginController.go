package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"goblog/models"
	"io"
	"strings"
)

type LoginController struct {
	beego.Controller
}

/*
注册页面
*/
func (this *LoginController) Register() {
	this.TplNames = "account/proregister.html"
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
			hpassword := md5.New()
			io.WriteString(hpassword, beego.AppConfig.String("salt1"))
			io.WriteString(hpassword, password)
			io.WriteString(hpassword, beego.AppConfig.String("salt2"))
			hpasswordfinal := fmt.Sprintf("%x", hpassword.Sum(nil))
			//记录ip地址
			ip := this.GetClientIp()
			err := models.Register(account, hpasswordfinal, username, ip)
			if err != nil {
				this.Ctx.WriteString("注册失败！")
			} else {
				/*
					注册时默认设置session，不设置cookie,根据account字段查找user表主码id
				*/
				userinfo, _ := models.Userinfo(account)
				uid := userinfo.Id
				//普通用户sessionuid
				sessionuid := this.GetSession("sessionuid")
				if sessionuid == nil {
					//如果session不存在，就创建session
					this.SetSession("sessionuid", uid)
					this.Redirect("/index", 301)
				} else {
					/*
						如果session存在，判断是否是本人博客，URL中获取用户博客名称，然后从
						数据库中读取用户账号,加密后判断是否和session数值一致，一致表示是本人
						博客，显示本人博客全部功能，不一致表示是其他人博客，只显示博客内容，不显示
						博客管理功能！
					*/
					this.Redirect("/index", 301)
					return
				}
			}

		}

	}
}

/*
登录页面
*/
func (this *LoginController) Post() {
	this.TplNames = "account/prologin.html"
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

		//查找用户名是否存在和用户密码是否正确
		searchaccount := models.SearchAccount(account)
		searchaccountandpwd := models.SearchAccountAndPwd(account, passwordinfofinal)

		userinfo1, _ := models.Searchuser(account)
		if len(account) == 0 {
			this.Data["errmsg"] = "账号不能为空！"
		} else if len(password) == 0 {
			this.Data["errmsg"] = "密码不能为空！"
		} else if searchaccount != nil {
			this.Data["errmsg"] = "账号不存在！"
		} else if searchaccountandpwd != nil {
			this.Data["errmsg"] = "密码不正确！"
		} else if userinfo1.Lock == 1 {
			this.Data["errmsg"] = "该用户已锁定！"
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
				this.Ctx.SetCookie("cookieaccount", cookieval, 7*24*3600, "/")
			}

			/*
				只把当前账号存到session中，登录时不记住密码一周默认设置session，不设置cookie
			*/

			userinfo, _ := models.Userinfo(account)
			uid := userinfo.Id
			sessionuid := this.GetSession("sessionuid")
			if sessionuid == nil {
				//如果session不存在，就创建session
				this.SetSession("sessionuid", uid)
				//记录最新登录时间,根据userid更新
				ip := this.GetClientIp()
				err := models.Updateuserlogintime(uid, ip)
				if err != nil {
					return
				}
				this.Redirect("/index", 301)
			} else {
				/*
					如果session存在，判断是否是本人博客，URL中获取用户博客名称，然后从
					数据库中读取用户账号,加密后判断是否和session数值一致，一致表示是本人
					博客，显示本人博客全部功能，不一致表示是其他人博客，只显示博客内容，不显示
					博客管理功能！
				*/
				//记录最新登录时间,根据userid更新
				ip := this.GetClientIp()
				err := models.Updateuserlogintime(uid, ip)
				if err != nil {
					return
				}
				this.Redirect("/index", 301)
				return
			}
		}
	}
}

/*
解析cookie内容
*/
func DecodeCookie(cookie string) string {
	//解密cookie
	decodecookie, _ := base64.StdEncoding.DecodeString(string(cookie))
	var value []string
	value = strings.Split(string(decodecookie), " ")
	return value[1]
}

/*
获取用户IP
*/

func (this *LoginController) GetClientIp() string {
	// s := this.Ctx.Request.Header.Get("X-Forwarded-For")
	s := this.Ctx.Input.IP()
	return s
}
