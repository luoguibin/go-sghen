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
type GetUserCreateParams struct {
	GetLoginParams
	Name     string `form:"name" valid:"Required"`
}

func (params *GetUserCreateParams) Valid(v *validation.Validation) {
	fmt.Println("GetUserCreateParams Valid")
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


/*****************************/
type GetUserUpdateParams struct {
	GetUserParams
	Name     string `form:"name"`
	Pw		 string `form:"pw"`
}

func (params *GetUserUpdateParams) Valid(v *validation.Validation) {
	fmt.Println("GetUserUpdateParams Valid")
    if len(strings.TrimSpace(params.Token)) == 0 {
        // 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("token", "不能为空")
    }
}


/*****************************/
type GetPeotrysetParams struct {
	GetUserParams
	SetName     string `form:"setName" valid:"Required"`
}

func (params *GetPeotrysetParams) Valid(v *validation.Validation) {
	fmt.Println("GetPeotrysetParams Valid")
    if len(strings.TrimSpace(params.Token)) == 0 {
        // 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("token", "不能为空")
    } else if len(strings.TrimSpace(params.SetName)) == 0 {
        // 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("setName", "不能为空")
    } 
}


/*****************************/
type GetPeotryimageParams struct {
	GetUserParams
	PId			   string	`form:"pId" valid:"Required"`
	ImageDatas     []string `form:"imageDatas" valid:"Required"`
}

func (params *GetPeotryimageParams) Valid(v *validation.Validation) {
	fmt.Println("GetPeotryimageParams Valid")
    if len(strings.TrimSpace(params.Token)) == 0 {
        // 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("token", "不能为空")
	} else if len(strings.TrimSpace(params.PId)) == 0 {
        // 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("PId", "不能为空")
    } 
}