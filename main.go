package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	"goblog/models"
	_ "goblog/routers"
	"time"
)

func init() {
	models.RegisterDB()
}

/*
时间戳转换模板函数
*/
func timeformat(time1 int) (out string) {
	curtime := int64(time1)
	out = time.Unix(curtime, 0).Format("2006-01-02 15:04:05")
	return
}
func timeformat1(time1 int) (out string) {
	curtime := int64(time1)
	out = time.Unix(curtime, 0).Format("2006-01-02")
	return
}

// func add(value int) (out int) {
// 	out = value + 1
// 	return
// }

// func dec(value int) (out int) {
// 	out = value - 5
// 	return
// }

func main() {
	orm.Debug = true
	orm.RunSyncdb("default", false, true)
	beego.SessionOn = true

	//模板函数
	beego.AddFuncMap("time", timeformat)
	beego.AddFuncMap("time1", timeformat1)
	beego.Run()
}
