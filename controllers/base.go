package controllers

import (
	"fmt"
	"SghenApi/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/plugins/cors"
)



/*****************************/
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

func (c *BaseController) respToJSON(data ResponseData) {
    respMsg, ok := data[models.RESP_MSG]
	if !ok || (ok && len(respMsg.(string)) <= 0) {
		data[models.RESP_MSG] = models.MConfig.CodeMsgMap[data[models.RESP_CODE].(int)]
	}
	// c.Ctx.Output.SetStatus(201)
	c.Data["json"] = data
    c.ServeJSON()
}

func (c *BaseController) BaseGetTest() {
	data := c.GetResponseData()
	c.respToJSON(data)
}

func (c *BaseController) CheckUserParams(data ResponseData, params interface{}) bool {
	//验证参数是否异常
	if err := c.ParseForm(params); err != nil {
		data[models.RESP_CODE] = models.RESP_ERR
		return false
	}
	fmt.Println("CheckUserParams")
	fmt.Println(params)

	//验证参数
	valid := validation.Validation{}
	if ok, _ := valid.Valid(params); ok {
		return true
	}

	data[models.RESP_CODE] = models.RESP_ERR
	data[models.RESP_MSG] = fmt.Sprint(valid.ErrorsMap)
	return false
}



/*****************************/
type ResponseData map[string]interface{}

func (self *BaseController) GetResponseData() ResponseData {
	return ResponseData{ models.RESP_CODE: models.RESP_OK }
}

