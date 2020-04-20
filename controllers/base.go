package controllers

import (
	"fmt"
	"go-sghen/helper"
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

// respToJSON JSON统一接口返回
func (c *BaseController) respToJSON(data ResponseData) {
	respMsg, ok := data[models.STR_MSG]
	if !ok || (ok && len(respMsg.(string)) <= 0) {
		data[models.STR_MSG] = models.MConfig.CodeMsgMap[data[models.STR_CODE].(int)]
	}
	// c.Ctx.Output.SetStatus(201)
	c.Data["json"] = data
	c.ServeJSON()
}

// TestGet 基础测试调用
func (c *BaseController) TestGet() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	data["ip"] = helper.GetRequestIP(c.Ctx)

	c.respToJSON(data)
}

func (c *BaseController) GetPageConfig() {
	data, _ := c.GetResponseData()

	c.respToJSON(data)
}

// CheckFormParams 检测表单信息
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

// ResponseData 请求返回体
type ResponseData map[string]interface{}

// GetResponseData 获取请求返回体
func (c *BaseController) GetResponseData() (ResponseData, bool) {
	// return ResponseData{models.STR_CODE: models.CODE_OK}, true
	return ResponseData{
		models.STR_CODE: models.CODE_MAINTENANCE,
		models.STR_DATA: "服务器维护中",
	}, false
}

// CheckAccessToken 检测用户身份，通过后将相关信息写入Input对象
func CheckAccessToken(ctx *context.Context) {
	datas := ResponseData{}
	// tokenCookie, err := ctx.Request.Cookie(models.STR_SGHEN_SESSION)
	// if err != nil {
	// 	datas[models.STR_CODE] = models.CODE_ERR
	// 	datas[models.STR_MSG] = "会话ID为空"
	// 	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	// 	ctx.Output.JSON(datas, false, true)
	// 	return
	// }

	token := ctx.Request.Header.Get("Authorization")
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

	userID, _ := strconv.ParseInt(claims["userId"].(string), 10, 64)
	userLevel, _ := strconv.Atoi(claims["uLevel"].(string))
	ctx.Input.SetData("userId", userID)
	ctx.Input.SetData("level", userLevel)
	return
}
