package controllers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

// APIController ...
type APIController struct {
	beego.Controller
}

// APIEntity ...
type APIEntity struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

var (
	// UserListEntity ...
	UserListEntity = APIEntity{
		Name: "user-list",
		Data: "",
	}

	// UserDataEntity ...
	UserDataEntity = APIEntity{
		Name: "user-data",
		Data: "id",
	}

	// APIRefreshCount ...
	APIRefreshCount = 0

	r = rand.New(rand.NewSource(time.Now().Unix()))
)

// RandString 生成随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 97 // 65大写  97小写
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// APIList 获取具体api列表
func (c *APIController) APIList() {
	data := make([]APIEntity, 0)
	UserListEntity.Name = RandString(10) + strconv.Itoa(APIRefreshCount)
	UserDataEntity.Name = RandString(10) + strconv.Itoa(APIRefreshCount)
	APIRefreshCount = APIRefreshCount + 1
	data = append(data, UserListEntity)
	data = append(data, UserDataEntity)
	c.Data["json"] = data
	c.ServeJSON()
}

// APICenter api处理入口
func (c *APIController) APICenter() {
	params := c.Ctx.Input.Params()
	delete(params, ":splat")

	if len(params) == 0 {
		c.APIList()
		return
	}

	data := ResponseData{}
	switch params["0"] {
	case UserListEntity.Name:
		apiUserList(params, data)
	case UserDataEntity.Name:
		apiUserData(params, data)
	default:
		apiError(data)
	}
	if (data["code"].(int)) == 404 {
		fmt.Println(404)
		c.Ctx.ResponseWriter.WriteHeader(404)
	} else {
		c.Data["json"] = data
		c.ServeJSON()
	}
}

func apiUserList(params map[string]string, data ResponseData) {
	length := len(params)
	if length == 1 {
		data["code"] = 200
		data["msg"] = "success"
		data["data"] = make([]APIEntity, 0)
	} else {
		apiError(data)
	}
}

func apiUserData(params map[string]string, data ResponseData) {
	length := len(params)

	if length == 2 {
		data["code"] = 200
		data["msg"] = "success"
		data["data"] = &APIEntity{
			Name: "liming",
			Data: params["1"],
		}
	} else {
		apiError(data)
	}
}

func apiError(data ResponseData) {
	data["code"] = 404
	data["msg"] = "api url format error"
}
