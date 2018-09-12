package controllers

import (
	"strings"
	"github.com/astaxie/beego/validation"
)

// 
// 有token验证的，均采用query传递该参数；
// 跳转路由前，在路由InsertFilter中，手动设置解析token中信息，
// 并写入对应的输入结构体中。
// 

// user create和login共用输入结构体
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

// user update和delete共用输入结构体
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

// user query输入结构体
type getQueryUserParams struct {
	Id			int64		`form:"uId" valid:"Required"`
	QueryId 	int64 		`form:"queryId" valid:"Required"`
	Level		int			`form:"level" valid:"Required"`
}

func (params *getQueryUserParams) Valid(v *validation.Validation) {
    if params.Id <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.QueryId <= 0 {
		v.SetError("user queryId", "不能为空")
	} else if params.Level < 0 {
		v.SetError("user level", "不能为空")
	} 
}

// user query输入结构体
type getQueryPeotryParams struct {
	Id			int64		`form:"id"`
	SId			int			`form:"setId"`
	Page		int			`form:"page"`
	Limit		int			`form:"limit"`
	Content		string		`form:"content"`
}

func (params *getQueryPeotryParams) Valid(v *validation.Validation) {
}

// peotry create输入结构体
type getCreatePeotryParams struct {
	UId 		int64		`form:"uId" valid:"Required"`
	SId			int 		`form:"sId" valid:"Required"`
	Title		string		`form:"title"`
	Content		string		`form:"content" valid:"Required"`
	End			string		`form:"end"`
}

func (params *getCreatePeotryParams) Valid(v *validation.Validation) {
    if params.UId <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.SId < 0 {
		v.SetError("set id", "不能为空")
	} else if len(strings.TrimSpace(params.Content)) < 5 {
		v.SetError("peotry content", "不能少于5个字符")
	} 
}

// peotry update输入结构体
type getUpdatePeotryParams struct {
	PId			int64		`form:"pId" valid:"Required"`
	UId 		int64		`form:"uId" valid:"Required"`
	SId			int 		`form:"sId" valid:"Required"`
	Title		string		`form:"title"`
	Content		string		`form:"content" valid:"Required"`
	End			string		`form:"end"`
}

func (params *getUpdatePeotryParams) Valid(v *validation.Validation) {
	if params.PId <= 0 {
		v.SetError("peotry id", "不能为空")
	} else if params.UId <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.SId < 0 {
		v.SetError("set id", "不能为空")
	} else if len(strings.TrimSpace(params.Content)) < 5 {
		v.SetError("peotry content", "不能少于5个字符")
	} 
}
