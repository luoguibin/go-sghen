// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"SghenApi/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.Router("/", &controllers.BaseController{}, "GET:BaseGetTest")

	beego.InsertFilter("/v1/user/update", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, true)
	})
	beego.InsertFilter("/v1/user/delete", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/user/query", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/peotry/create", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/peotry/update", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/peotry/delete", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/peotry-set/query", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/peotry-set/create", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/peotry-set/delete", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/upload", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, true)
	})

	//详见　https://beego.me/docs/mvc/controller/router.md
	nsv1 := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSRouter("/create", &controllers.UserController{}, "post:CreateUser"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:LoginUser"),
			beego.NSRouter("/update", &controllers.UserController{}, "post:UpdateUser"),
			beego.NSRouter("/delete", &controllers.UserController{}, "delete:DeleteUser"),
			beego.NSRouter("/query", &controllers.UserController{}, "get:QueryUser"),
		),
		beego.NSNamespace("/peotry",
			beego.NSRouter("/query", &controllers.PeotryController{}, "get:QueryPeotry"),
			beego.NSRouter("/create", &controllers.PeotryController{}, "get:CreatePeotry"),
			beego.NSRouter("/update", &controllers.PeotryController{}, "get:UpdatePeotry"),
			beego.NSRouter("/delete", &controllers.PeotryController{}, "delete:DeletePeotry"),
		),
		beego.NSNamespace("/peotry-set",
			beego.NSRouter("/query", &controllers.PeotrySetController{}, "get:QueryPeotrySet"),
			beego.NSRouter("/create", &controllers.PeotrySetController{}, "get:CreatePeotrySet"),
			beego.NSRouter("/delete", &controllers.PeotrySetController{}, "delete:DeletePeotrySet"),
		),
		beego.NSRouter("/upload", &controllers.FileUploaderController{}, "post:FileUpload"),
	)

	beego.AddNamespace(nsv1)
}
