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
	data := c.GetResponseData()
	params := &getCreateUserParams{}

	if c.CheckFormParams(data, params) {
		if params.Code != "test" {
			err := checkSmsCode(params.ID, params.Code)
			if err != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = err.Error()
				c.respToJSON(data)
				return
			}
		}

		user, err := models.CreateUser(params.ID, params.Pw, params.Name, 1)
		if err == nil {
			createUserToken(c.Ctx, user, data)
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			errStr := err.Error()

			if strings.Contains(errStr, "PRIMARY") {
				data[models.STR_MSG] = "已存在该用户"
			} else {
				data[models.STR_MSG] = "用户注册失败"
			}
		}
	}

	c.respToJSON(data)
}

// LoginUser 登录
func (c *UserController) LoginUser() {
	data := c.GetResponseData()
	params := &getCreateUserParams{}

	if c.CheckFormParams(data, params) {
		user, err := models.QueryUser(params.ID)

		if err == nil {
			compare := -1
			// 优先判断是否启动验证码登陆
			if len(strings.TrimSpace(params.Code)) > 0 {
				smsErr := checkSmsCode(params.ID, params.Code)
				if smsErr != nil {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = smsErr.Error()
					c.respToJSON(data)
					return
				}
				compare = 0
			} else {
				compare = strings.Compare(user.Password, params.Pw)
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
	}
	c.respToJSON(data)
}

// QueryUser 查询user，限制level等级为5以下的user
func (c *UserController) QueryUser() {
	data := c.GetResponseData()
	params := &getQueryUserParams{}

	if c.CheckFormParams(data, params) {
		if params.Level >= 5 {
			user, err := models.QueryUser(params.QueryUID)
			if err == nil {
				data[models.STR_DATA] = user
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "未查询到对应用户"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "用户等级低，限制查询"
		}
	}
	c.respToJSON(data)
}

// QueryUsers 查询users，输入ids的json
func (c *UserController) QueryUsers() {
	data := c.GetResponseData()
	params := &getQueryUsersParams{}

	if c.CheckFormParams(data, params) {
		IDStrs := strings.Split(params.IDStrs, ",")
		if len(IDStrs) > 0 && len(IDStrs) < 100 {
			ids := make([]int64, 0)
			for _, IDStr := range IDStrs {
				id, err := strconv.ParseInt(IDStr, 10, 64)
				if err == nil {
					ids = append(ids, id)
				}
			}
			users, err := models.QueryUsers(ids)
			if err == nil {
				data[models.STR_DATA] = users
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "未查询到对应用户列表"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			if len(IDStrs) > 100 {
				data[models.STR_MSG] = "查询用户个数不能超过100个"
			} else {
				data[models.STR_MSG] = "请输入用户id列表"
			}
		}
	}
	c.respToJSON(data)
}

// UpdateUser 更新user
func (c *UserController) UpdateUser() {
	data := c.GetResponseData()
	params := &getUpdateUserParams{}

	if c.CheckFormParams(data, params) {
		userID := c.Ctx.Input.GetData("userId").(int64)

		if params.ID == userID {
			_, err := models.UpdateUser(params.ID, params.Pw, params.Name, params.IconURL)
			if err != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "更新用户信息失败"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "禁止更新他人用户信息"
		}
	}

	c.respToJSON(data)
}

// DeleteUser 删除user
func (c *UserController) DeleteUser() {
	data := c.GetResponseData()
	params := &getUpdateUserParams{}

	if c.CheckFormParams(data, params) {
		err := models.DeleteUser(params.ID)
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "删除用户失败"
		}
	}

	c.respToJSON(data)
}

// checkSmsCode prod模式校验验证码
func checkSmsCode(ID int64, Code string) error {
	if models.MConfig.SGHENENV != "prod" {
		return nil
	}
	smsCode, err := models.QuerySmsCode(ID)
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
	models.DeleteSmsCode(ID)
	return nil
}

// createUserToken 创建用户token，基于json web token
func createUserToken(c *context.Context, user *models.User, data ResponseData) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userId"] = strconv.FormatInt(user.ID, 10)
	claims["userName"] = user.Name
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
