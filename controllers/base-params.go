package controllers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/validation"
)


/*****************************/
type getCreateUserParams struct {
	Id     	int64 	`form:"id" valid:"Required"`
	Pw     	string 	`form:"pw" valid:"Required"`
	Name   	string 	`form:"name"`
}

func (params *getCreateUserParams) Valid(v *validation.Validation) {
	if params.Id <= 0 {
		v.SetError("id", "不能为空")
	} else if len(strings.TrimSpace(params.Pw)) == 0 {
		v.SetError("pw", "不能为空")
	} 
}


/*****************************/
type getUpdateUserParams struct {
	Id			int64		`form:"uId" valid:"Required"`
	Pw     	string 	`form:"pw"`
	Name		string	`form:"name"`
}

func (params *getUpdateUserParams) Valid(v *validation.Validation) {
    if params.Id <= 0 {
			v.SetError("user id", "不能为空")
		}
}



/*****************************/
type getQueryUserParams struct {
	Id			int64		`form:"uId" valid:"Required"`
	QueryId int64 	`form:"queryId" valid:"Required"`
	Level		int			`form:"level" valid:"Required"`
}

func (params *getQueryUserParams) Valid(v *validation.Validation) {
    if params.Id <= 0 {
			v.SetError("user id", "不能为空")
		} else if params.QueryId <= 0 {
			v.SetError("user queryId", "不能为空")
		} else if params.Level <= 0 {
			v.SetError("user level", "不能为空")
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
		v.SetError("token", "不能为空")
    } else if len(strings.TrimSpace(params.SetName)) == 0 {
		v.SetError("setName", "不能为空")
    } 
}


/*****************************/
type GetPeotrysetUpdateParams struct {
	GetUserParams
	SetName     string `form:"setName" valid:"Required"`
	SId     string `form:"sId" valid:"Required"`
}

func (params *GetPeotrysetUpdateParams) Valid(v *validation.Validation) {
	fmt.Println("GetPeotrysetUpdateParams Valid")
    if len(strings.TrimSpace(params.Token)) == 0 {
		v.SetError("token", "不能为空")
    } else if len(strings.TrimSpace(params.SetName)) == 0 {
		v.SetError("setName", "不能为空")
    }  else if len(strings.TrimSpace(params.SId)) == 0 {
		v.SetError("setId", "不能为空")
    } 
}


/*****************************/
type GetPeotryUpdateParams struct {
	GetUserParams
	PId			   string	`form:"pId" valid:"Required"`
}

func (params *GetPeotryUpdateParams) Valid(v *validation.Validation) {
	fmt.Println("GetPeotryUpdateParams Valid")
    if len(strings.TrimSpace(params.Token)) == 0 {
		v.SetError("token", "不能为空")
	} else if len(strings.TrimSpace(params.PId)) == 0 {
		v.SetError("PId", "不能为空")
    } 
}


/*****************************/
type GetPeotryimageParams struct {
	GetPeotryUpdateParams
	ImageDatas     []string `form:"imageDatas" valid:"Required"`
}

func (params *GetPeotryimageParams) Valid(v *validation.Validation) {
	fmt.Println("GetPeotryimageParams Valid")
    if len(strings.TrimSpace(params.Token)) == 0 {
		v.SetError("token", "不能为空")
	} else if len(strings.TrimSpace(params.PId)) == 0 {
		v.SetError("PId", "不能为空")
    } 
}