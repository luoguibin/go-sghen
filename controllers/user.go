package controllers

import (
	"errors"
	"fmt"
	"go-sghen/helper"
	"go-sghen/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/context"
	jwt "github.com/dgrijalva/jwt-go"
)

// UserController 用户控制器
type UserController struct {
	BaseController
}

// CreateUser 创建用户
func (c *UserController) CreateUser() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getCreateUserParams{}
	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	// 默认以手机注册账号
	mobile := params.Account
	if params.Type == 0 {
		err := checkSmsCode(mobile, params.Code)
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = err.Error()
			c.respToJSON(data)
			return
		}
	} else {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "暂只支持手机账号注册"
		c.respToJSON(data)
		return
	}

	_, err := models.QueryUser(params.Account, mobile)
	if err == nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "该账号已被注册"
		c.respToJSON(data)
		return
	}

	user, err := models.CreateUser(params.Account, mobile, params.Pw, params.Name, "", "", 1)
	if err == nil {
		createUserToken(c.Ctx, user, data)

		// 首次创建用户，写入消息
		models.CreateSysMsg(user.ID, models.MODULE_USER_CREATE, "")
	} else {
		data[models.STR_CODE] = models.CODE_ERR
		errStr := err.Error()

		if strings.Contains(errStr, "PRIMARY") {
			data[models.STR_MSG] = "已存在该用户"
		} else {
			data[models.STR_MSG] = "用户注册失败"
		}
	}

	c.respToJSON(data)
}

// LoginUser 登录
func (c *UserController) LoginUser() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getCreateUserParams{}

	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	account := params.Account
	mobile := params.Account
	if params.Type == 0 {
		account = ""
	} else {
		mobile = ""
	}
	user, err := models.QueryUser(account, mobile)
	if err == nil {
		compare := -1
		// 优先判断是否启动验证码登陆
		if len(strings.TrimSpace(params.Code)) > 0 {
			smsErr := checkSmsCode(params.Account, params.Code)
			if smsErr != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = smsErr.Error()
				c.respToJSON(data)
				return
			}
			compare = 0
		} else {
			if len(strings.TrimSpace(params.Random)) > 0 {
				user.UserPWD = helper.MD5(user.UserPWD + params.Random)
			}
			compare = strings.Compare(user.UserPWD, params.Pw)
		}

		if compare == 0 {
			createUserToken(c.Ctx, user, data)
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "用户账号或密码错误"
		}
	} else {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户账号错误或用户不存在"
	}
	c.respToJSON(data)
}

// UpdateUser 更新user
func (c *UserController) UpdateUser() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getUpdateUserParams{}

	if c.CheckFormParams(data, params) {
		userID := c.Ctx.Input.GetData("userId").(int64)

		_, err := models.UpdateUser(userID, params.Pw, params.Name, params.Avatar, params.Mood)
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "更新用户信息失败"
		}
	}

	c.respToJSON(data)
}

// DeleteUser 删除user
func (c *UserController) DeleteUser() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getUpdateUserParams{}
	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	data[models.STR_CODE] = models.CODE_ERR
	data[models.STR_MSG] = "接口维护中"

	c.respToJSON(data)
}

// checkSmsCode prod模式校验验证码
func checkSmsCode(mobile, Code string) error {
	if models.MConfig.SGHENENV != "prod" {
		return nil
	}
	id, err := strconv.ParseInt(mobile, 10, 64)
	if err != nil {
		return errors.New("手机号码格式错误")
	}
	smsCode, err := models.QuerySmsCode(id)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return errors.New("验证码服务错误")
	}
	if smsCode == nil {
		return errors.New("请重新发送验证码")
	}
	if smsCode.Code != Code {
		if smsCode.CountRead >= 2 {
			models.DeleteSmsCode(smsCode.ID)
			return errors.New("验证码错误，请重新发送")
		}
		smsCode.CountRead = smsCode.CountRead + 1
		models.SaveSmsCode(smsCode.ID, smsCode.Code, smsCode.CountRead, smsCode.TimeLife)
		return errors.New("验证码错误")
	}
	timeVal := helper.GetMillisecond() - smsCode.TimeCreate
	if timeVal < 0 || timeVal > smsCode.TimeLife {
		return errors.New("验证码已过有效期")
	}
	models.DeleteSmsCode(id)
	return nil
}

// createUserToken 创建用户token，基于json web token
func createUserToken(c *context.Context, user *models.User, data ResponseData) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	nowTime := time.Now().Unix()
	claims["exp"] = nowTime + models.MConfig.JwtExpireDuration
	claims["iat"] = nowTime
	claims["userId"] = strconv.FormatInt(user.ID, 10)
	claims["userName"] = user.UserName
	claims["uLevel"] = strconv.Itoa(user.Level)

	token.Claims = claims

	tokenString, err := token.SignedString([]byte(models.MConfig.JwtSecretKey))
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户id签名失败"
		return
	}

	// uidCookie := &http.Cookie{
	// 	Name:     models.STR_SGHEN_SESSION,
	// 	Value:    tokenString,
	// 	HttpOnly: true,
	// }
	// http.SetCookie(c.ResponseWriter, uidCookie)

	user.Token = tokenString
	data["expireDuration"] = models.MConfig.JwtExpireDuration
	data[models.STR_DATA] = user
}

// CheckUserToken 检测解析token
func CheckUserToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(models.MConfig.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token get mapcliams err")
}
