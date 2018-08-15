package controllers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/validation"
)

/*****************************/
type GetUserParams struct {
	Token     string `form:"token" valid:"Required"`
}

func (params *GetUserParams) Valid(v *validation.Validation) {
	fmt.Println("GetUserParams Valid")
    if len(strings.TrimSpace(params.Token)) == 0 {
        // 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("token", "不能为空")
    }
}
