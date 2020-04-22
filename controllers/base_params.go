package controllers

import (
	"go-sghen/helper"
	"strconv"
	"strings"

	"github.com/astaxie/beego/validation"
)

// getSmsSendParams 短信验证码
type getSmsSendParams struct {
	Phone int64 `form:"phone" valid:"Required"`
	CaptchaID string `form:"captchaId" valid:"Required"`
	CaptchaValue string `form:"captchaValue" valid:"Required"`
}

func (params *getSmsSendParams) Valid(v *validation.Validation) {
	phoneStr := strconv.FormatInt(params.Phone, 10)
	if !helper.IsPhone(phoneStr) {
		v.SetError("phone", "请输入正确的手机号码")
	}
	if len(params.CaptchaID) == 0 {
		v.SetError("captchaId", "校验码ID错误")
	}
	if len(params.CaptchaValue) == 0 {
		v.SetError("cap", "校验码错误")
	}
}

// getCommonParams ...
type getCommonParams struct {
	MType string `form:"type" valid:"Required"`
	Data  string `form:"data"`
	Field bool   `form:"field"`
}

func (params *getCommonParams) Valid(v *validation.Validation) {
	if len(strings.TrimSpace(params.MType)) == 0 {
		v.SetError("type", "type类型错误")
	}
}

// user create和login共用输入结构体
type getCreateUserParams struct {
	Account string `form:"account" valid:"Required"`
	Pw      string `form:"pw" valid:"Required"`
	Name    string `form:"name"`
	Type    int    `form:"type"`
	Code    string `form:"code"`
}

func (params *getCreateUserParams) Valid(v *validation.Validation) {
	if len(strings.TrimSpace(params.Account)) == 0 {
		v.SetError("account", "不能为空")
	} else if len(strings.TrimSpace(params.Pw)) == 0 {
		v.SetError("pw", "不能为空")
	} else if params.Type < 0 || params.Type > 1 {
		v.SetError("type", "账号类型错误")
	}
}

// user update和delete共用输入结构体
type getUpdateUserParams struct {
	Pw     string `form:"pw" json:"pw"`
	Name   string `form:"name" json:"name"`
	Avatar string `form:"avatar" json:"avatar"`
	Mobile string `form:"mobile" json:"mobile"`
	Mood   string `form:"mood" json:"mood"`
}

func (params *getUpdateUserParams) Valid(v *validation.Validation) {
}

// user query输入结构体
type getQueryUserParams struct {
	ID       int64 `form:"id" valid:"Required"`
	QueryUID int64 `form:"queryId" valid:"Required"`
	Level    int   `form:"level" valid:"Required"`
}

func (params *getQueryUserParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.QueryUID <= 0 {
		v.SetError("user queryId", "不能为空")
	} else if params.Level < 0 {
		v.SetError("user level", "不能为空")
	}
}

// users query输入结构体
type getQueryUsersParams struct {
	IDStrs string `form:"idStrs" valid:"Required"`
}

func (params *getQueryUsersParams) Valid(v *validation.Validation) {
	if len(params.IDStrs) == 0 {
		v.SetError("user ids", "不能为空")
	}
}

// peotry query输入结构体
type getQueryPeotryParams struct {
	ID          int64  `form:"id"`
	UserID      int64  `form:"userId"`
	SetID       int    `form:"setId"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
	Content     string `form:"content"`
	NeedComment bool   `form:"needComment"`
}

func (params *getQueryPeotryParams) Valid(v *validation.Validation) {
}

// peotry create输入结构体
type getCreatePeotryParams struct {
	UserID     int64  `form:"userId" valid:"Required"`
	SetID      int    `form:"setId" valid:"Required"`
	Title      string `form:"title"`
	Content    string `form:"content" valid:"Required"`
	ImageNames string `form:"imageNames"`
	End        string `form:"end"`
}

func (params *getCreatePeotryParams) Valid(v *validation.Validation) {
	if params.UserID <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.SetID < 0 {
		v.SetError("set id", "不能为空")
	} else if len(strings.TrimSpace(params.Content)) < 5 {
		v.SetError("peotry content", "不能少于5个字符")
	}
}

// peotry update输入结构体
type getUpdatePeotryParams struct {
	ID      int64  `form:"id" valid:"Required"`
	UserID  int64  `form:"userId" valid:"Required"`
	SetID   int    `form:"setId" valid:"Required"`
	Title   string `form:"title"`
	Content string `form:"content" valid:"Required"`
	End     string `form:"end"`
}

func (params *getUpdatePeotryParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("peotry id", "不能为空")
	} else if params.UserID <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.SetID < 0 {
		v.SetError("set id", "不能为空")
	} else if len(strings.TrimSpace(params.Content)) < 5 {
		v.SetError("peotry content", "不能少于5个字符")
	}
}

// peotry delete输入结构体
type getDeletePeotryParams struct {
	UserID int64 `form:"userId" valid:"Required"`
	ID     int64 `form:"id" valid:"Required"`
}

func (params *getDeletePeotryParams) Valid(v *validation.Validation) {
	if params.UserID <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.ID <= 0 {
		v.SetError("peotry id", "不能为空")
	}
}

// peotryset query输入结构体
type getQueryPoetrySetParams struct {
	UserID int64 `form:"userId" valid:"Required"`
}

func (params *getQueryPoetrySetParams) Valid(v *validation.Validation) {
	if params.UserID == 0 {
		v.SetError("set id", "不能为空")
	}
}

// peotryset create
type getCreatePoetrySetParams struct {
	Name string `form:"name" valid:"Required"`
}

func (params *getCreatePoetrySetParams) Valid(v *validation.Validation) {
	if len(strings.TrimSpace(params.Name)) == 0 {
		v.SetError("set name", "不能为空")
	}
}

// peotryset delete
type getDeletePoetrySetParams struct {
	ID int `form:"id" valid:"Required"`
}

func (params *getDeletePoetrySetParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("set id", "不能为空")
	}
}

// comment create输入结构体
type getCreateCommentParams struct {
	Type    int    `form:"type" json:"type" valid:"Required"`
	TypeID  int64  `form:"typeId" json:"typeId" valid:"Required"`
	FromID  int64  `form:"fromId" json:"fromId" valid:"Required"`
	ToID    int64  `form:"toId" json:"toId" valid:"Required"`
	Content string `form:"content" json:"content" valid:"Required"`
}

func (params *getCreateCommentParams) Valid(v *validation.Validation) {
	if params.Type <= 0 {
		v.SetError("type", "不能为空")
	} else if params.TypeID <= 0 {
		v.SetError("typeId", "不能为空")
	} else if params.FromID <= 0 {
		v.SetError("fromId", "不能为空")
	} else if len(strings.TrimSpace(params.Content)) == 0 {
		v.SetError("comment", "不能为空")
	}
}

// comment query,delete输入结构体
type getQueryCommentParams struct {
	Type   int   `form:"type" json:"type" valid:"Required"`
	TypeID int64 `form:"typeId" json:"typeId" valid:"Required"`
}

func (params *getQueryCommentParams) Valid(v *validation.Validation) {
	if params.Type <= 0 {
		v.SetError("type", "不能为空")
	} else if params.TypeID <= 0 {
		v.SetError("typeId", "不能为空")
	}
}

// comment query,delete输入结构体
type getDeleteCommentParams struct {
	ID     int64 `form:"id" json:"id" valid:"Required"`
	FromID int64 `form:"fromId" json:"fromId" valid:"Required"`
}

func (params *getDeleteCommentParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("id", "不能为空")
	} else if params.FromID <= 0 {
		v.SetError("fromId", "不能为空")
	}
}

// DynamicAPI create
type getCreateDynamicAPIParams struct {
	SuffixPath string `form:"suffixPath" valid:"Required"`
	Name       string `form:"name" valid:"Required"`
	Comment    string `form:"comment"`
	Content    string `form:"content" valid:"Required"`
	Status     int    `form:"status"`
}

func (params *getCreateDynamicAPIParams) Valid(v *validation.Validation) {
	if len(strings.TrimSpace(params.SuffixPath)) == 0 {
		v.SetError("set suffixPath", "不能为空")
	}
	if len(strings.TrimSpace(params.Name)) == 0 {
		v.SetError("set name", "不能为空")
	}
	if len(strings.TrimSpace(params.Content)) == 0 {
		v.SetError("set content", "不能为空")
	}
}

// DynamicAPI create
type getUpdateDynamicAPIParams struct {
	ID         int64  `form:"id" valid:"Required"`
	SuffixPath string `form:"suffixPath" valid:"Required"`
	Name       string `form:"name" valid:"Required"`
	Comment    string `form:"comment"`
	Content    string `form:"content" valid:"Required"`
	Status     int    `form:"status"`
}

func (params *getUpdateDynamicAPIParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("set id", "不能为空")
	}
	if len(strings.TrimSpace(params.SuffixPath)) == 0 {
		v.SetError("set suffixPath", "不能为空")
	}
	if len(strings.TrimSpace(params.Content)) == 0 {
		v.SetError("set content", "不能为空")
	}
}

// getQueryDynamicAPIParams qeury
type getQueryDynamicAPIParams struct {
	ID         int64  `form:"id"`
	SuffixPath string `form:"suffixPath"`
	Name       string `form:"name"`
	Comment    string `form:"comment"`
	Status     int    `form:"status"`
	UserID     int64  `form:"userId"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
}

func (params *getQueryDynamicAPIParams) Valid(v *validation.Validation) {
}

// getDeleteDynamicAPIParams delete
type getDeleteDynamicAPIParams struct {
	ID int64 `form:"id" valid:"Required"`
}

func (params *getDeleteDynamicAPIParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("set id", "不能为空")
	}
}
