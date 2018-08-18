package controllers

import (
	"fmt"
	"time"
	"strconv"
	"encoding/json"
	"SghenApi/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/plugins/cors"
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
		fmt.Println(err)
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
func (c *BaseController)CreateUserToken (user *models.User, data ResponseData) {
	token := jwt.New(jwt.SigningMethodHS256)
    claims := make(jwt.MapClaims)
    claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["uid"] = strconv.FormatInt(user.Id, 10)
	fmt.Println(claims)

    token.Claims = claims

    tokenString, err := token.SignedString([]byte(models.JWT_SECRET_KEY))
    if err != nil {
		data[models.RESP_CODE] = models.RESP_ERR
		data[models.RESP_MSG] = "Error while signing the token"
		return
	}
	
	data[models.RESP_TOKEN] = tokenString
}

func (c *BaseController)ParseUserToken (tokenString string) (map[string]interface{}, error) {
	fmt.Println("ParseUserToken()", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		return []byte(models.JWT_SECRET_KEY), nil
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



/*****************************/
type ResponseData map[string]interface{}

func (self *BaseController) GetResponseData() ResponseData {
	return ResponseData{ models.RESP_CODE: models.RESP_OK }
}

