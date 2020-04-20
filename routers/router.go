// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"go-sghen/controllers"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.Router("/", &controllers.BaseController{}, "GET:TestGet")

	// 设置路由过滤器，校验身份
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		flag := ctx.Request.Method == "POST"
		flag = flag && strings.Index(ctx.Request.URL.Path, "login") == -1
		flag = flag && strings.Index(ctx.Request.URL.Path, "/user/create") == -1
		flag = flag && strings.Index(ctx.Request.URL.Path, "/sms/send") == -1
		if flag {
			controllers.CheckAccessToken(ctx)
		}
	})

	// 路由定义
	nsv1 := beego.NewNamespace("/v1",
		beego.NSNamespace("/common",
			beego.NSRouter("/page-config", &controllers.BaseController{}, "GET:GetPageConfig"),
		),
		beego.NSNamespace("/user",
			beego.NSRouter("/create", &controllers.UserController{}, "post:CreateUser"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:LoginUser"),
			beego.NSRouter("/update", &controllers.UserController{}, "post:UpdateUser"),
			beego.NSRouter("/delete", &controllers.UserController{}, "post:DeleteUser"),
			beego.NSRouter("/query", &controllers.UserController{}, "get:QueryUser"),
			beego.NSRouter("/query-list", &controllers.UserController{}, "get:QueryUsers"),
		),
		beego.NSNamespace("/peotry",
			beego.NSRouter("/query", &controllers.PeotryController{}, "get:QueryPeotry"),
			beego.NSRouter("/query-popular", &controllers.PeotryController{}, "get:QueryPopularPeotry"),
			beego.NSRouter("/create", &controllers.PeotryController{}, "post:CreatePeotry"),
			beego.NSRouter("/update", &controllers.PeotryController{}, "post:UpdatePeotry"),
			beego.NSRouter("/delete", &controllers.PeotryController{}, "post:DeletePeotry"),
			beego.NSRouter("/add-temp", &controllers.PeotryController{}, "post:AddTempPeotry"),
		),
		beego.NSNamespace("/peotry-set",
			beego.NSRouter("/query", &controllers.PeotrySetController{}, "get:QueryPeotrySet"),
			beego.NSRouter("/create", &controllers.PeotrySetController{}, "post:CreatePeotrySet"),
			beego.NSRouter("/delete", &controllers.PeotrySetController{}, "post:DeletePeotrySet"),
		),
		beego.NSNamespace("/comment",
			beego.NSRouter("/query", &controllers.CommentController{}, "get:QueryComments"),
			beego.NSRouter("/create", &controllers.CommentController{}, "post:CreateComment"),
			beego.NSRouter("/delete", &controllers.CommentController{}, "post:DeleteComment"),
		),
		beego.NSNamespace("/api",
			beego.NSRouter("/create", &controllers.DynamicAPIController{}, "post:CreateDynamicAPI"),
			beego.NSRouter("/update", &controllers.DynamicAPIController{}, "post:UpdateDynamicAPI"),
			beego.NSRouter("/query", &controllers.DynamicAPIController{}, "get:QueryDynamicAPI"),
			beego.NSRouter("/delete", &controllers.DynamicAPIController{}, "post:DeleteDynamicAPI"),
			beego.NSRouter("/get/*", &controllers.DynamicAPIController{}, "get:GetDynamicDataByPath"),
			beego.NSRouter("/post", &controllers.DynamicAPIController{}, "post:PostDynamicData"),
		),
		beego.NSNamespace("/sms/",
			beego.NSRouter("/send", &controllers.SmsController{}, "post:SendSmsCode"),
		),
		beego.NSRouter("/upload", &controllers.FileUploaderController{}, "post:FileUpload"),
	)

	beego.AddNamespace(nsv1)
}
