package controllers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/validation"
)



/*****************************/
type GetLoginParams struct {
	Id     string `form:"id" valid:"Required"`
	Pw     string `form:"pw" valid:"Required"`
}

func (params *GetLoginParams) Valid(v *validation.Validation) {
	fmt.Println("GetLoginParams Valid")
    if len(strings.TrimSpace(params.Id)) == 0 {
		v.SetError("id", "不能为空")
    } else if len(strings.TrimSpace(params.Pw)) == 0 {
		v.SetError("pw", "不能为空")
    } 
}

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
