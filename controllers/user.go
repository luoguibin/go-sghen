package controllers

import (
	"SghenApi/models"
	"time"
	"strconv"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// UserController operations for User
type UserController struct {
	BaseController
}

func (c *UserController)CreateUser() {
	data := c.GetResponseData()
	params := &getCreateUserParams{}
	if (c.CheckPostParams(data, params)) {
		user, err := models.CreateUser(params.Id, params.Pw, params.Name)
		if err == nil {
			createUserToken(user, data)
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	}
	
	c.respToJSON(data)
}

func (c *UserController)LoginUser() {
	data := c.GetResponseData()
	params := &getCreateUserParams{}
	if (c.CheckPostParams(data, params)) {
		user, err := models.LoginUser(params.Id, params.Pw)
		if err == nil {
			createUserToken(user, data)
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	}
	c.respToJSON(data)
}

func (c *UserController)QueryUser() {
	data := c.GetResponseData()
	params := &getQueryUserParams{}

	if c.CheckFormParams(data, params) {
		if params.Level >= 5 {
			user, err := models.QueryUser(params.QueryId)
			if err == nil {
				data[models.RESP_DATA] = user
			} else {
				data[models.RESP_CODE] = models.RESP_ERR
				data[models.RESP_MSG] = err.Error()
			}		
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = "用户等级低，限制查询"
		}
	}
	c.respToJSON(data)
}

func (c *UserController)UpdateUser() {
	data := c.GetResponseData()
	params := &getUpdateUserParams{}

	if c.CheckPostParams(data, params) {
		fmt.Println(params)
		_, err := models.UpdateUser(params.Id, params.Pw, params.Name)
		if err != nil {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	}
	
	c.respToJSON(data)
}

func (c *UserController)DeleteUser() {
	data := c.GetResponseData()
	params := &getUpdateUserParams{}

	if c.CheckFormParams(data, params) {
		fmt.Println(params)
		err := models.DeleteUser(params.Id)
		if err != nil {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	}
	
	c.respToJSON(data)
}


func createUserToken(user *models.User, data ResponseData) {
	token := jwt.New(jwt.SigningMethodHS256)
    claims := make(jwt.MapClaims)
    claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["uid"] = strconv.FormatInt(user.ID, 10)
	claims["ulevel"] = strconv.Itoa(user.ULevel)

    token.Claims = claims

    tokenString, err := token.SignedString([]byte(models.MConfig.JwtSecretKey))
    if err != nil {
		data[models.RESP_CODE] = models.RESP_ERR
		data[models.RESP_MSG] = "Error while signing the token"
		return
	}
	
	data[models.RESP_TOKEN] = tokenString
}

func CheckUserToken(tokenString string) (map[string]interface{}) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		return []byte(models.MConfig.JwtSecretKey), nil
	})

	if (err != nil) {
		return nil
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	} else {
		return nil
	}
}