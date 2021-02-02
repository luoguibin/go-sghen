package controllers

import "go-sghen/models"

// SysMsgContronller ...
type SysMsgContronller struct {
	BaseController
}

// QuerySysMsgs ...
func (c *SysMsgContronller) QuerySysMsgs() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getQuerySysMsgParams{}
	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	list, count, err := models.GetSysMsgs(userID, params.Status, params.Page, params.Limit)

	if err == nil {
		data[models.STR_DATA] = list
		data["count"] = count
	} else {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "未查询到消息"
	}

	c.respToJSON(data)
}

// ReadSysMsg ...
func (c *SysMsgContronller) ReadSysMsg() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getUpdateSysMsgParams{}
	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	err := models.UpdateSysMsgStatus(params.ID, userID, 1)

	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "更新消息失败"
	}

	c.respToJSON(data)
}
