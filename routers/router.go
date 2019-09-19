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

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	// only GET, POST method
	beego.Router("/", &controllers.BaseController{}, "GET:BaseGetTest")

	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		if ctx.Request.Method == "POST" {
			controllers.GatewayAccessUser(ctx)
		}
	})

	//详见　https://beego.me/docs/mvc/controller/router.md
	nsv1 := beego.NewNamespace("/v1",
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
			beego.NSRouter("/create", &controllers.PeotryController{}, "post:CreatePeotry"),
			beego.NSRouter("/update", &controllers.PeotryController{}, "post:UpdatePeotry"),
			beego.NSRouter("/delete", &controllers.PeotryController{}, "post:DeletePeotry"),
		),
		beego.NSNamespace("/peotry-set",
			beego.NSRouter("/query", &controllers.PeotrySetController{}, "get:QueryPeotrySet"),
			beego.NSRouter("/create", &controllers.PeotrySetController{}, "get:CreatePeotrySet"),
			beego.NSRouter("/delete", &controllers.PeotrySetController{}, "post:DeletePeotrySet"),
		),
		beego.NSNamespace("/comment",
			beego.NSRouter("/query", &controllers.CommentController{}, "get:QueryComments"),
			beego.NSRouter("/create", &controllers.CommentController{}, "post:CreateComment"),
			beego.NSRouter("/delete", &controllers.CommentController{}, "post:DeleteComment"),
		),
		beego.NSRouter("/upload", &controllers.FileUploaderController{}, "post:FileUpload"),
	)

	beego.AddNamespace(nsv1)
}
