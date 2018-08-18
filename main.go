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
	orm.DefaultTimeLoc = time.UTC
	orm.Debug = true
	orm.RegisterDataBase("default", "mysql", "root:root123@tcp(127.0.0.1:3306)/sghen?charset=utf8&loc=Local")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

