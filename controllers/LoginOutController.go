package controllers

import (
	"github.com/astaxie/beego"
	// "goblog/models"
)

type LoginOutController struct {
	beego.Controller
}

func (this *LoginOutController) Get() {
	op := this.Input().Get("op")
	//管理员和普通用户退出页面不相同，需要判断要退出的用户类型
	if op == "admin" {
		//管理员退出操作,获取session和cookie值
		sessionaid := this.GetSession("sessionaid")
		cookieadmin := this.Ctx.GetCookie("cookieadmin")
		if len(cookieadmin) > 0 {
			this.Ctx.SetCookie("cookieadmin", "", -1, "/")
		}
		if sessionaid != nil {
			this.DestroySession()
		}
		this.Redirect("/admin/login", 301)
		return
	} else {
		//普通用户退出操作,获取session和cookie值
		sessionuid := this.GetSession("sessionuid")
		cookieaccount := this.Ctx.GetCookie("cookieaccount")
		if len(cookieaccount) > 0 {
			this.Ctx.SetCookie("cookieaccount", "", -1, "/")
		}
		if sessionuid != nil {
			this.DestroySession()
		}
		this.Redirect("/", 301)
		return
	}
}
