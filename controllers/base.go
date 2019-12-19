package controllers

import (
	"fmt"
	"go-sghen/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/validation"
)

// BaseController ...
type BaseController struct {
	beego.Controller
}

func init() {
	fmt.Println("basecontroller init")
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
}

// respToJSON 统一接口返回
func (c *BaseController) respToJSON(data ResponseData) {
	respMsg, ok := data[models.STR_MSG]
	if !ok || (ok && len(respMsg.(string)) <= 0) {
		data[models.STR_MSG] = models.MConfig.CodeMsgMap[data[models.STR_CODE].(int)]
	}
	// c.Ctx.Output.SetStatus(201)
	c.Data["json"] = data
	c.ServeJSON()
}

// BaseGetTest 基础测试调用
func (c *BaseController) BaseGetTest() {
	data := c.GetResponseData()

	ip := c.Ctx.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = c.Ctx.Request.Header.Get("X-real-ip")
	}

	if ip == "" {
		ip = "127.0.0.1"
	}
	data["ip"] = ip

	c.respToJSON(data)
}

// GatewayAccessUser ...
func GatewayAccessUser(ctx *context.Context) {
	datas := ResponseData{}
	token := ctx.Input.Query("token")

	if len(token) <= 0 {
		datas[models.STR_CODE] = models.CODE_ERR
		datas[models.STR_MSG] = "token不能为空"
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Output.JSON(datas, false, true)
		return
	}

	claims, err := CheckUserToken(token)
	if err != nil {
		datas[models.STR_CODE] = models.CODE_ERR_TOKEN
		errStr := err.Error()

		if strings.Contains(errStr, "expired") {
			datas[models.STR_MSG] = "token失效，请重新登录"
		} else {
			datas[models.STR_MSG] = "token参数错误"
		}

		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Output.JSON(datas, false, false)
		return
	}

	userId, _ := strconv.ParseInt(claims["userId"].(string), 10, 64)
	userLevel, _ := strconv.Atoi(claims["uLevel"].(string))
	ctx.Input.SetData("userId", userId)
	ctx.Input.SetData("level", userLevel)
	// ctx.Input.Context.Request.Form.Add("userId", claims["userId"].(string))
	return
}

// CheckFormParams ...
func (c *BaseController) CheckFormParams(data ResponseData, params interface{}) bool {
	// 验证参数是否异常
	if err := c.ParseForm(params); err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		return false
	}

	//验证参数
	valid := validation.Validation{}
	if ok, _ := valid.Valid(params); ok {
		return true
	}
	data[models.STR_CODE] = models.CODE_ERR
	data[models.STR_MSG] = fmt.Sprint(valid.ErrorsMap)
	return false
}

// CheckPostParams ...
// func (c *BaseController) CheckPostParams(data ResponseData, params interface{}) bool {
// 	// 验证参数是否异常
// 	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &params); err != nil {
// 		data[models.STR_CODE] = models.CODE_ERR
// 		data[models.STR_MSG] = err.Error()
// 		return false
// 	}
// 	// 验证参数
// 	valid := validation.Validation{}
// 	if ok, _ := valid.Valid(params); ok {
// 		return true
// 	}

// 	data[models.STR_CODE] = models.CODE_ERR
// 	data[models.STR_MSG] = fmt.Sprint(valid.ErrorsMap)
// 	return false
// }

// ResponseData ...
type ResponseData map[string]interface{}

// GetResponseData ...
func (c *BaseController) GetResponseData() ResponseData {
	return ResponseData{models.STR_CODE: models.CODE_OK}
}
