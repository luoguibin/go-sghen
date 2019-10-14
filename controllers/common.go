package controllers

import (
	"bytes"
	"encoding/json"
	"go-sghen/helper"
	"go-sghen/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// CommonController ...
type CommonController struct {
	BaseController
}

// GetFlag ...
func (c *CommonController) GetFlag() {
	data := c.GetResponseData()
	c.respToJSON(data)
}

// SendSmsCode ...
func (c *CommonController) SendSmsCode() {
	data := c.GetResponseData()
	params := &getSmsSendParams{}

	if c.CheckFormParams(data, params) {
		random := helper.GetMicrosecond()
		sdkAppID := 0 // todo
		appKey := ""  // todo
		time := time.Now().Unix()

		text := "appKey=" + appKey + "&random=" + strconv.FormatInt(random, 10) + "&time=" + strconv.FormatInt(time, 10) + "&mobile=" + strconv.FormatInt(params.Phone, 10)

		url := "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=" + strconv.Itoa(sdkAppID) + "&random=" + strconv.FormatInt(random, 10)

		requestbody := make(map[string]interface{})
		tel := make(map[string]interface{})
		requestbody["params"] = []string{"8080"}
		requestbody["sig"] = helper.Sha256(text)
		requestbody["sign"] = "Sghen三行"
		tel["mobile"] = strconv.FormatInt(params.Phone, 10)
		tel["nationcode"] = "86"
		requestbody["tel"] = tel
		requestbody["time"] = time
		requestbody["tpl_id"] = 442430

		bytesData, err := json.Marshal(requestbody)

		resp, err := http.Post(url, "application/json", bytes.NewReader(bytesData))
		if err != nil {
			models.MConfig.MLogger.Error(err.Error())

			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "验证码服务掉线"
			c.respToJSON(data)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = err.Error()
		} else {
			data[models.STR_DATA] = string(body)
		}
	}
	c.respToJSON(data)
}
