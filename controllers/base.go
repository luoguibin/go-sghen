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

// func (c *BaseController) AllowCross() {
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")       //允许访问源
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
//     c.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
// }