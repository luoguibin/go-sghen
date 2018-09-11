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
		controllers.GatewayAccessUser(ctx)
	})

	//详见　https://beego.me/docs/mvc/controller/router.md
	nsv1 := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSRouter("/create", &controllers.UserController{}, "post:CreateUser"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:LoginUser"),
			beego.NSRouter("/update", &controllers.UserController{}, "post:UpdateUser"),
			// TODO
			beego.NSRouter("/query", &controllers.UserController{}, "post:QueryUser"),
			beego.NSRouter("/delete", &controllers.UserController{}, "delete:DeleteUser"),
		),
	)

	beego.AddNamespace(nsv1)
}
