package routers

import (
	"github.com/astaxie/beego"
	"goblog/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/index", &controllers.IndexController{})
	beego.Router("/admin", &controllers.AdminController{})
	beego.Router("/loginout", &controllers.LoginOutController{})
	beego.Router("/article", &controllers.ArticleController{})
	beego.Router("/diary", &controllers.DiaryController{})
	beego.Router("/image", &controllers.ImageController{})
	//Login控制器自动路由
	beego.AutoRouter(&controllers.LoginController{})
	beego.AutoRouter(&controllers.AdminController{})
	beego.AutoRouter(&controllers.IndexController{})
	beego.AutoRouter(&controllers.ArticleController{})
	beego.AutoRouter(&controllers.DiaryController{})
	beego.AutoRouter(&controllers.ImageController{})
}
