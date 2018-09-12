package main

import (
	_ "SghenApi/routers"
	"SghenApi/models"
	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

func main() {
	models.InitGorm()
	db := models.GetDb()
	defer db.Close()

	peotry, err := models.QueryPeotry(1536735430035598)
	fmt.Println(err)
	fmt.Println(peotry)
	
	beego.Run()	
}
