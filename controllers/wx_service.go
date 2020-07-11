package controllers

import (
	"encoding/json"
	"errors"
	"go-sghen/models"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
)

// WxServiceController 微信服务控制器
type WxServiceController struct {
	BaseController
}

// WxLoginResult 登录凭证校验返回实体
type WxLoginResult struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

// LoginWxUser 微信登录凭证校验入口
func (c *WxServiceController) LoginWxUser() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getWxLoginParams{}
	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	if models.MConfig.SGHENENV != "prod" {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "非正式环境中不支持小程序登陆服务"
		c.respToJSON(data)
		return
	}

	result, err := verifyLoginCode(params.Code)
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = err.Error()
		c.respToJSON(data)
		return
	}

	if params.Type == "bind" {
		// 绑定微信用户
		userID := c.Ctx.Input.GetData("userId").(int64)
		_, err := models.UpdateUserAccount(userID, result.OpenID, "")
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = err.Error()
		}
		c.respToJSON(data)
		return
	}

	// 获取到OpenID，查询本系统user
	user, err := models.QueryUser(result.OpenID, "")
	if err != nil {
		// 账号不存在时应提示是否关联
		if gorm.IsRecordNotFoundError(err) {
			data[models.STR_CODE] = models.CODE_NOT_FOUND
			data[models.STR_MSG] = "用户不存在"
			// user, err := models.CreateUser(result.OpenID, "", "", params.Name, "", "", 1)
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = err.Error()
		}
		c.respToJSON(data)
		return
	}

	createUserToken(c.Ctx, user, data)
	c.respToJSON(data)
}

// BindWxUser 绑定微信用户
func (c *WxServiceController) BindWxUser() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getWxLoginParams{}
	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	if models.MConfig.SGHENENV != "prod" {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "非正式环境中不支持小程序登陆服务"
		c.respToJSON(data)
		return
	}

	result, err := verifyLoginCode(params.Code)
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = err.Error()
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	_, err = models.UpdateUserAccount(userID, result.OpenID, "")
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = err.Error()
	}
	c.respToJSON(data)
}

// 微信登录凭证校验
func verifyLoginCode(code string) (WxLoginResult, error) {
	wxAppID := models.MConfig.WxAppID
	wxSecret := models.MConfig.WxSecret
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + wxAppID + "&secret=" + wxSecret + "&js_code=" + code + "&grant_type=authorization_code"

	var result WxLoginResult
	resp, err := http.Get(url)
	if err != nil {
		models.MConfig.MLogger.Error(err.Error())
		return result, errors.New("微信登录凭证校验错误")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	json.Unmarshal(body, &result)
	if result.ErrCode != 0 {
		// -1	系统繁忙，此时请开发者稍候再试
		// 0	请求成功
		// 40029	code 无效
		// 45011	频率限制，每个用户每分钟100次
		return result, errors.New(result.ErrMsg)
	}
	return result, nil
}
