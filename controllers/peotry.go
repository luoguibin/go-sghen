package controllers

import (
	"SghenApi/models"
	"SghenApi/helper"
)

// PeotryController operations for Peotry
type PeotryController struct {
	BaseController
}

func (c *PeotryController) QueryPeotry() {
	data := c.GetResponseData();
	params := &getQueryPeotryParams{}

	if c.CheckFormParams(data, params) {
		list, err, count, totalPage, currentPage, pageIsEnd := models.QueryPeotry(params.Id, params.SId, params.Page, params.Limit, params.Content)
		if err == nil {
			if params.Id > 0 {
				data[models.RESP_DATA] = list[0]
			} else {
				data[models.RESP_DATA] = list
				data["totalCount"] = count
				data["totalPage"] = totalPage
				data["currentPage"] = currentPage
				data["pageIsEnd"] = pageIsEnd
			}
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	}
	c.respToJSON(data)
}


func (c *PeotryController) CreatePeotry() {
	data := c.GetResponseData()
	params := &getCreatePeotryParams{}

	if c.CheckFormParams(data, params) {
		set, err := models.QueryPeotrySetByID(params.SId)
		if err == nil {
			if set.UID == params.UId {
				timeStr := helper.GetNowDateTime()
				pId, err := models.SavePeotry(params.UId, params.SId, params.Title, timeStr, params.Content, params.End, "[]")
				if err == nil {
					data[models.RESP_DATA] = pId
				} else {
					data[models.RESP_CODE] = models.RESP_ERR
					data[models.RESP_MSG] = err.Error()
				}
			} else {
				data[models.RESP_CODE] = models.RESP_ERR
				data[models.RESP_MSG] = "禁止在他人诗集中创建个人诗歌"
			}
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	}

	c.respToJSON(data)
}

func (c *PeotryController) UpdatePeotry() {
	data := c.GetResponseData()
	params := &getUpdatePeotryParams{}
	
	if c.CheckFormParams(data, params) {
		list, err, _, _, _, _ := models.QueryPeotry(params.PId, 0, 0, 0, "")
		qPeotry := list[0]
		if err == nil {
			if qPeotry.UID == params.UId {
				// 判断选集是否有更新
				if qPeotry.SID != params.SId {
					set, err := models.QueryPeotrySetByID(params.SId)
					if err == nil {
						if set.UID == params.UId {
							qPeotry.SID = params.SId
						} else {
							data[models.RESP_CODE] = models.RESP_ERR
							data[models.RESP_MSG] = "禁止在他人诗集中更新个人诗歌"
							c.respToJSON(data)
							return
						}
					} else {
						data[models.RESP_CODE] = models.RESP_ERR
						data[models.RESP_MSG] = err.Error()
						c.respToJSON(data)
						return
					}
				}
				qPeotry.PTitle  = params.Title
				qPeotry.PContent = params.Content
				qPeotry.PEnd = params.End
				// 更新时需要将这些附带的结构体置空
				qPeotry.UUser = nil
				qPeotry.SSet = nil
				qPeotry.PImage = nil

				err := models.UpdatePeotry(qPeotry)
				if err == nil {
					data[models.RESP_DATA] = qPeotry.ID
				} else {
					data[models.RESP_CODE] = models.RESP_ERR
					data[models.RESP_MSG] = err.Error()
				}
			} else {
				data[models.RESP_CODE] = models.RESP_ERR
				data[models.RESP_MSG] = "禁止更新他人诗歌"
			}
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	}

	c.respToJSON(data)
}

func (c *PeotryController) DeletePeotry() {
	data := c.GetResponseData()
	params := &getDeletePeotryParams{}

	if c.CheckFormParams(data, params) { 
		list, err, _, _, _, _ := models.QueryPeotry(params.PID, 0, 0, 0, "")
		if err == nil {
			peotry := list[0]
			if peotry.UID == params.UID {
				err := models.DeletePeotry(params.PID)
				if err == nil {
					data[models.RESP_DATA] = true
				} else {
					data[models.RESP_CODE] = models.RESP_ERR
					data[models.RESP_MSG] = err.Error()
				}
			} else {
				data[models.RESP_CODE] = models.RESP_ERR
				data[models.RESP_MSG] = "禁止删除他人诗歌"
			}
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	}

	c.respToJSON(data)
}