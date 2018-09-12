package controllers

import (
	"fmt"
	"strconv"
	"encoding/json"
	"SghenApi/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/context"
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

func GatewayAccessUser(ctx *context.Context, setInPost bool) {
	datas := make(map[string]interface{})
	// userId := ctx.Input.Query("userId")
	token := ctx.Input.Query("token")

	if len(token) <= 0 {
		datas[models.RESP_CODE] = models.RESP_ERR
		datas[models.RESP_MSG] = "token is empty"
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Output.JSON(datas, false, false)
		return
	}

	claims := CheckUserToken(token)
	if  claims == nil{
		datas[models.RESP_CODE] = models.RESP_ERR
		datas[models.RESP_MSG] = "token is invalid"
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Output.JSON(datas, false, false)
		return
	}
	
	if setInPost {
		uId, _ := strconv.ParseInt(claims["uid"].(string), 10, 64)
		ctx.Input.SetData("uId", uId)
		ctx.Input.SetData("level", claims["uid"].(int))
	} else {
		ctx.Input.Context.Request.Form.Add("uId", claims["uid"].(string))
		ctx.Input.Context.Request.Form.Add("level", claims["uid"].(string))
	}
	
	return
}


func (c *BaseController) CheckFormParams(data ResponseData, params interface{}) bool {
	//验证参数是否异常
	if err := c.ParseForm(params); err != nil {
		data[models.RESP_CODE] = models.RESP_ERR
		return false
	}
	fmt.Println("CheckFormParams")
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

func (c *BaseController) CheckPostParams(data ResponseData, params interface{}) bool {
	//验证参数是否异常
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &params); err != nil {
		data[models.RESP_CODE] = models.RESP_ERR
		data[models.RESP_MSG] = err.Error()
		return false
	}

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
