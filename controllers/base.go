package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

type BaseController struct {
	beego.Controller
}

func init() {
	fmt.Println("basecontroller init");
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
        AllowAllOrigins:  true,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
        AllowCredentials: true,
    }))
}

// func (c *BaseController) ServeJSON() {
    
//     c.ServeJSON()
// }

func (c *BaseController) BaseGetTest() {
	var jsonObj = make(map[string]interface{})
	jsonObj["code"] = 1000
	jsonObj["message"] = "ok"

	c.Data["json"] = jsonObj
	c.ServeJSON()
}