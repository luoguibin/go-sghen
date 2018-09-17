package controllers

import (
	"SghenApi/models"
	"fmt"
	"time"
	"strconv"
	"errors"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

// UserController operations for User
type UserController struct {
	BaseController
}

// 创建user
func (c *UserController)CreateUser() {
	data := c.GetResponseData()
	params := &getCreateUserParams{}
	if (c.CheckPostParams(data, params)) {
		user, err := models.CreateUser(params.ID, params.Pw, params.Name)
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

// user登录
func (c *UserController)LoginUser() {
	data := c.GetResponseData()
	params := &getCreateUserParams{}
	if (c.CheckPostParams(data, params)) {
		user, err := models.QueryUser(params.ID)
		if err == nil {
			if (user.UPassword == params.Pw) {
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

// 查询user，限制level等级为５以下的user
func (c *UserController)QueryUser() {
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

// 更新user
func (c *UserController)UpdateUser() {
	data := c.GetResponseData()
	params := &getUpdateUserParams{}

	if c.CheckPostParams(data, params) {
		fmt.Println(params)
		_, err := models.UpdateUser(params.ID, params.Pw, params.Name)
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "更新用户信息失败"
		}
	}
	
	c.respToJSON(data)
}

// 删除user
func (c *UserController)DeleteUser() {
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

// 创建用户token，基于json web token
func createUserToken(user *models.User, data ResponseData) {
	token := jwt.New(jwt.SigningMethodHS256)
    claims := make(jwt.MapClaims)
    claims["exp"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["uId"] = strconv.FormatInt(user.ID, 10)
	claims["uLevel"] = strconv.Itoa(user.ULevel)

    token.Claims = claims

    tokenString, err := token.SignedString([]byte(models.MConfig.JwtSecretKey))
    if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户id签名失败"
		return
	}
	
	user.UToken = tokenString
	data[models.STR_DATA] = user
}

// 检测解析token
func CheckUserToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		return []byte(models.MConfig.JwtSecretKey), nil
	})

	if (err != nil) {
		return nil, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token get mapcliams err")
	}
}