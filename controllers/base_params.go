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
	ID   int64  `form:"uId" json:"uId" valid:"Required"`
	Pw   string `form:"pw" json:"pw" valid:"Required"`
	Name string `form:"name"`
}

func (params *getCreateUserParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("uId", "不能为空")
	} else if len(strings.TrimSpace(params.Pw)) == 0 {
		v.SetError("pw", "不能为空")
	}
}

// user update和delete共用输入结构体
type getUpdateUserParams struct {
	ID      int64  `form:"uId" json:"uId" valid:"Required"`
	Pw      string `form:"pw" json:"pw"`
	Name    string `form:"name" json:"name"`
	IconURL string `form:"iconUrl" json:"iconUrl"`
}

func (params *getUpdateUserParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("user id", "不能为空")
	}
}

// user query输入结构体
type getQueryUserParams struct {
	ID       int64 `form:"uId" valid:"Required"`
	QueryUID int64 `form:"queryUId" valid:"Required"`
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
	ID          int64  `form:"pId"`
	SID         int    `form:"setId"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
	Content     string `form:"qContent"`
	NeedComment bool   `form:"needComment"`
}

func (params *getQueryPeotryParams) Valid(v *validation.Validation) {
}

// peotry create输入结构体
type getCreatePeotryParams struct {
	UID     int64  `form:"uId" valid:"Required"`
	SID     int    `form:"sId" valid:"Required"`
	Title   string `form:"title"`
	Content string `form:"content" valid:"Required"`
	End     string `form:"end"`
}

func (params *getCreatePeotryParams) Valid(v *validation.Validation) {
	if params.UID <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.SID < 0 {
		v.SetError("set id", "不能为空")
	} else if len(strings.TrimSpace(params.Content)) < 5 {
		v.SetError("peotry content", "不能少于5个字符")
	}
}

// peotry update输入结构体
type getUpdatePeotryParams struct {
	PID     int64  `form:"pId" valid:"Required"`
	UID     int64  `form:"uId" valid:"Required"`
	SID     int    `form:"sId" valid:"Required"`
	Title   string `form:"title"`
	Content string `form:"content" valid:"Required"`
	End     string `form:"end"`
}

func (params *getUpdatePeotryParams) Valid(v *validation.Validation) {
	if params.PID <= 0 {
		v.SetError("peotry id", "不能为空")
	} else if params.UID <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.SID < 0 {
		v.SetError("set id", "不能为空")
	} else if len(strings.TrimSpace(params.Content)) < 5 {
		v.SetError("peotry content", "不能少于5个字符")
	}
}

// peotry delete输入结构体
type getDeletePeotryParams struct {
	UID int64 `form:"uId" valid:"Required"`
	PID int64 `form:"pId" valid:"Required"`
}

func (params *getDeletePeotryParams) Valid(v *validation.Validation) {
	if params.UID <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.PID <= 0 {
		v.SetError("peotry id", "不能为空")
	}
}

// peotryset query输入结构体
type getQueryPoetrySetParams struct {
	UID int64 `form:"uId" valid:"Required"`
}

func (params *getQueryPoetrySetParams) Valid(v *validation.Validation) {
	if params.UID == 0 {
		v.SetError("set id", "不能为空")
	}
}

// peotryset create
type getCreatePoetrySetParams struct {
	UID  int64  `form:"uId" valid:"Required"`
	Name string `form:"name" valid:"Required"`
}

func (params *getCreatePoetrySetParams) Valid(v *validation.Validation) {
	if params.UID <= 0 {
		v.SetError("user id", "不能为空")
	} else if len(strings.TrimSpace(params.Name)) == 0 {
		v.SetError("set name", "不能为空")
	}
}

// peotryset delete
type getDeletePoetrySetParams struct {
	UID int64 `form:"uId" valid:"Required"`
	SID int   `form:"sId" valid:"Required"`
}

func (params *getDeletePoetrySetParams) Valid(v *validation.Validation) {
	if params.UID <= 0 {
		v.SetError("user id", "不能为空")
	} else if params.SID <= 0 {
		v.SetError("set id", "不能为空")
	}
}

// comment create输入结构体
type getCreateCommentParams struct {
	Type    int    `form:"type" json:"type" valid:"Required"`
	TypeID  int64  `form:"typeId" json:"typeId" valid:"Required"`
	FromID  int64  `form:"fromId" json:"fromId" valid:"Required"`
	ToID    int64  `form:"toId" json:"toId" valid:"Required"`
	Content string `form:"Content" json:"Content" valid:"Required"`
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
