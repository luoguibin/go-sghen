package main

import (
	_ "SghenApi/routers"
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	// _ "SghenApi/tests"
)

func init() {
	orm.RegisterDataBase(
		"default", 
		"mysql", 
		beego.AppConfig.String("mysqluser") + ":" + beego.AppConfig.String("mysqlpass") + "@tcp(" + beego.AppConfig.String("mysqlurls") + ")/" + beego.AppConfig.String("mysqldb")+"?charset=utf8&loc=Asia%2FShanghai")
	orm.DefaultTimeLoc = time.UTC
	orm.Debug = true
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()	
}
