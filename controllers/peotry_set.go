package controllers

import (
	"go-sghen/models"
)

// PeotrysetController operations for Peotryset
type PeotrySetController struct {
	BaseController
}

// 查询用户的选集和系统默认选集
func (c *PeotrySetController) QueryPeotrySet() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getQueryPoetrySetParams{}

	if c.CheckFormParams(data, params) {
		list, err := models.QueryPeotrySetByUID(params.UserID)

		if err == nil {
			data[models.STR_DATA] = list
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "获取用户诗歌选集失败"
		}
	}

	c.respToJSON(data)
}

// CreatePeotrySet 创建选集
func (c *PeotrySetController) CreatePeotrySet() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getCreatePoetrySetParams{}

	if c.CheckFormParams(data, params) {
		userID := c.Ctx.Input.GetData("userId").(int64)
		err := models.CreatePeotrySet(userID, params.Name)

		if err == nil {
			data[models.STR_DATA] = true
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "创建选集失败"
		}
	}

	c.respToJSON(data)
}

// DeletePeotrySet 删除选集
func (c *PeotrySetController) DeletePeotrySet() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getDeletePoetrySetParams{}

	if c.CheckFormParams(data, params) {
		set, err := models.QueryPeotrySetByID(params.ID)

		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "未获取到选集id"
			c.respToJSON(data)
			return
		}

		userID := c.Ctx.Input.GetData("userId").(int64)
		if set.UserID != userID {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "禁止删除非本人创建的诗集"
			c.respToJSON(data)
			return
		}

		_, count, _, _, _, err := models.QueryPeotry(0, set.ID, 1, 10, "")
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "尝试查询该选集下的诗词失败"
			c.respToJSON(data)
			return
		}

		if count > 0 {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "该选集下还有诗词，禁止删除"
			c.respToJSON(data)
			return
		}

		err = models.DeletePeotrySet(params.ID)
		if err == nil {
			data[models.STR_DATA] = true
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "删除选集失败"
		}
	}

	c.respToJSON(data)
}
