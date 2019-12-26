package controllers

import (
	"bytes"
	"encoding/json"
	"go-sghen/helper"
	"go-sghen/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SmsController 短信控制器
type SmsController struct {
	BaseController
}

// SmsResult 短信返回信息结构体
type SmsResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Fee    int    `json:"fee"`
	Sid    string `json:"sid"`
}

// SendSmsCode 发送短信验证码
func (c *SmsController) SendSmsCode() {
	data := c.GetResponseData()
	params := &getSmsSendParams{}

	if c.CheckFormParams(data, params) {
		if models.MConfig.SGHENENV != "prod" {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "非正式环境中已暂停验证码服务"
			c.respToJSON(data)
			return
		}

		user, _ := models.QueryUser(params.Phone)
		if user != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "该账号已注册"
			c.respToJSON(data)
			return
		}

		sdkAppID := models.MConfig.SmsSdkAppID
		appKey := models.MConfig.SmsAppKey
		if sdkAppID == 0 || len(appKey) == 0 {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "验证码服务尚未运行"
			models.SaveSmsCode(params.Phone, "1234", 0, 2*60*1000)
			c.respToJSON(data)
			return
		}

		smsCode, err := models.QuerySmsCode(params.Phone)
		if err != nil && !strings.Contains(err.Error(), "record not found") {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "验证码服务错误"
			c.respToJSON(data)
			return
		}

		if smsCode != nil {
			timeVal := helper.GetMillisecond() - smsCode.TimeCreate
			if timeVal > 0 && timeVal < 60*1000 {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "稍后再发送验证码"
				c.respToJSON(data)
				return
			}
		}
		models.SaveSmsCode(params.Phone, "", 0, 2*60*1000)

		random := helper.GetMicrosecond()
		time := time.Now().Unix()
		codeStr := helper.RandomNum4()
		text := "appkey=" + appKey + "&random=" + strconv.FormatInt(random, 10) + "&time=" + strconv.FormatInt(time, 10) + "&mobile=" + strconv.FormatInt(params.Phone, 10)
		url := "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=" + strconv.Itoa(sdkAppID) + "&random=" + strconv.FormatInt(random, 10)

		requestbody := make(map[string]interface{})
		tel := make(map[string]interface{})
		requestbody["params"] = []string{codeStr}
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
			models.DeleteSmsCode(params.Phone)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = err.Error()
		} else {
			var smsResult SmsResult
			json.Unmarshal(body, &smsResult)

			if smsResult.Result != 0 {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = smsResult.Errmsg
			} else {
				_, err := models.SaveSmsCode(params.Phone, codeStr, 0, 2*60*1000)
				if err != nil {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "验证码服务错误"
				}
			}
		}
	}
	c.respToJSON(data)
}