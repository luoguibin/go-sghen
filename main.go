package main

import (
	_ "SghenApi/routers"
	"SghenApi/models"
	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	models.InitGorm()
	db := models.GetDb()
	defer db.Close()

	beego.Run()	
}
