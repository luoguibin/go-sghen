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

func createUserToken(user *models.User, data ResponseData) {
	token := jwt.New(jwt.SigningMethodHS256)
    claims := make(jwt.MapClaims)
    claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["uid"] = strconv.FormatInt(user.ID, 10)
	fmt.Println(claims)

    token.Claims = claims

    tokenString, err := token.SignedString([]byte(models.MConfig.JwtSecretKey))
    if err != nil {
		data[models.RESP_CODE] = models.RESP_ERR
		data[models.RESP_MSG] = "Error while signing the token"
		return
	}
	
	data[models.RESP_TOKEN] = tokenString
}

func ParseUserToken(tokenString string) (map[string]interface{}, error) {
	fmt.Println("ParseUserToken()", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		return []byte(models.MConfig.JwtSecretKey), nil
	})

	if (err != nil) {
		fmt.Println(err)
		return nil, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		return claims, nil
	} else {
		fmt.Println("Claims parse error", err)
		return nil, err
	}
}

func (c *UserController)QueryUser() {
	data := c.GetResponseData()
	models.QueryUser()
	c.respToJSON(data)
}

func (c *UserController)UpdateUser() {
	data := c.GetResponseData()
	models.UpdateUser()
	c.respToJSON(data)
}

func (c *UserController)DeleteUser() {
	data := c.GetResponseData()
	models.DeleteUser()
	c.respToJSON(data)
}