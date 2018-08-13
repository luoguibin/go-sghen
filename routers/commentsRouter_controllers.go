package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotryimageController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:PeotrysetController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:UserController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:UserController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:UserController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:UserController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["SghenApi/controllers:UserController"] = append(beego.GlobalControllerRouter["SghenApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

}
