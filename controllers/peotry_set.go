package controllers

import (
	"go-sghen/models"
)

// PeotrysetController operations for Peotryset
type PeotrySetController struct {
	BaseController
}

func (c *PeotrySetController) QueryPeotrySet() {
	data := c.GetResponseData()
	params := &getQueryPoetrySetParams{}

	if c.CheckFormParams(data, params) {
		list, err := models.QueryPeotrySetByUID(params.UID)
		if err == nil {
			data[models.STR_DATA] = list
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "获取用户诗歌选集失败"
		}
	}

	c.respToJSON(data)
}

func (c *PeotrySetController) CreatePeotrySet() {
	data := c.GetResponseData()
	params := &getCreatePoetrySetParams{}

	if c.CheckFormParams(data, params) {
		err := models.CreatePeotrySet(params.UID, params.Name)
		if err == nil {
			data[models.STR_DATA] = true
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "创建选集失败"
		}
	}

	c.respToJSON(data)
}

func (c *PeotrySetController) DeletePeotrySet() {
	data := c.GetResponseData()
	params := &getDeletePoetrySetParams{}

	if c.CheckFormParams(data, params) {
		set, err := models.QueryPeotrySetByID(params.SID)
		if err == nil {
			if set.UID == params.UID {
				err := models.DeletePeotrySet(params.SID)
				if err == nil {
					data[models.STR_DATA] = true
				} else {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "删除选集失败"
				}
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "禁止删除非本人创建的诗集"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "未获取到选集id"
		}
	}

	c.respToJSON(data)
}
