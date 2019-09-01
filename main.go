package main

import (
	"go-sghen/models"
	_ "go-sghen/routers"
	"os"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	beego.SetStaticPath("file", "./file")
	os.Setenv("ZONEINFO", "./lib/time/zoneinfo.zip")

	models.InitGorm()
	db := models.GetDb()
	defer db.Close()

	beego.Run()
}
