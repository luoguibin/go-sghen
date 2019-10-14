package controllers

import (
	"errors"
	"fmt"
	"go-sghen/helper"
	"go-sghen/models"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserController operations for User
type UserController struct {
	BaseController
}

// CreateUser 创建user
func (c *UserController) CreateUser() {
	data := c.GetResponseData()
	params := &getCreateUserParams{}

	if c.CheckFormParams(data, params) {
		smsCode, err := models.QuerySmsCode(params.ID)
		if err != nil && !strings.Contains(err.Error(), "record not found") {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "验证码服务错误"
			return
		}
		if smsCode != nil {
			if smsCode.Code != params.Code {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "验证码错误"
				return
			}
			timeVal := helper.GetMillisecond() - smsCode.TimeCreate
			if timeVal < 0 || timeVal > smsCode.TimeLife {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "验证码已过有效期"
				return
			}
			models.DeleteSmsCode(params.ID)
		}

		user, err := models.CreateUser(params.ID, helper.MD5(params.Pw), params.Name, 1)
		if err == nil {
			createUserToken(user, data)
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
			compare := strings.Compare(user.Password, params.Pw)

			if compare == 0 {
				createUserToken(user, data)
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "用户账号或密码错误"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "用户账号或密码错误"
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

// createUserToken 创建用户token，基于json web token
func createUserToken(user *models.User, data ResponseData) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userId"] = strconv.FormatInt(user.ID, 10)
	claims["uLevel"] = strconv.Itoa(user.Level)

	token.Claims = claims

	tokenString, err := token.SignedString([]byte(models.MConfig.JwtSecretKey))
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户id签名失败"
		return
	}

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
	} else {
		return nil, errors.New("token get mapcliams err")
	}
}
