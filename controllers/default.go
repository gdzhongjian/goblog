package controllers

import (
	"github.com/astaxie/beego"
	// "goblog/models"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	/*
		判断存不存在cookie信息和session信息
	*/
	cookieaccount := this.Ctx.GetCookie("cookieaccount")
	uid := this.GetSession("sessionuid")
	//存在cookie
	if len(cookieaccount) > 0 {
		this.Redirect("/index", 301)
	} else if uid != nil {
		//存在session
		this.Redirect("/index", 301)

	} else {
		this.TplNames = "account/prologin.html"
	}

}
