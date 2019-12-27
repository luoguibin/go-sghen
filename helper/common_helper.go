package helper

import (
	"regexp"
	"strings"

	"github.com/astaxie/beego/context"
)

//IsPhone 判断是否为手机号码
func IsPhone(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

// GetRequestIP 获取请求的真实IP
func GetRequestIP(ctx *context.Context) string {
	ip := ctx.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = ctx.Request.Header.Get("X-real-ip")
	}

	if ip == "" {
		ip = "127.0.0.1"
	}
	return ip
}
