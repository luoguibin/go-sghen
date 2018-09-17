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
		if params.ID > 0 {
			peotry, err := models.QueryPeotryByID(params.ID)
			if err == nil {
				data[models.STR_DATA] = peotry
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "未查询到对应id的诗歌"
			}
		} else {
			list, err, count, totalPage, curPage, pageIsEnd := models.QueryPeotry(params.SID, params.Page, params.Limit, params.Content)
			if err == nil {
				data[models.STR_DATA] = list
				data["totalCount"] = count
				data["totalPage"] = totalPage
				data["curPage"] = curPage
				data["pageIsEnd"] = pageIsEnd
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "未查询到对应的诗歌"
			}
		}
	}
	c.respToJSON(data)
}


func (c *PeotryController) CreatePeotry() {
	data := c.GetResponseData()
	params := &getCreatePeotryParams{}

	if c.CheckFormParams(data, params) {
		set, err := models.QueryPeotrySetByID(params.SID)
		if err == nil {
			if set.UID == params.UID {
				timeStr := helper.GetNowDateTime()
				pId, err := models.CreatePeotry(params.UID, params.SID, params.Title, timeStr, params.Content, params.End, "[]")
				if err == nil {
					data[models.STR_DATA] = pId
				} else {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "创建诗歌失败"
				}
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "禁止在他人选集中创建个人诗歌"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "非本人创建的选集下不能创建新的诗歌"
		}
	}

	c.respToJSON(data)
}

func (c *PeotryController) UpdatePeotry() {
	data := c.GetResponseData()
	params := &getUpdatePeotryParams{}
	
	if c.CheckFormParams(data, params) {
		qPeotry, err := models.QueryPeotryByID(params.PID)
		if err == nil {
			if qPeotry.UID == params.UID {
				// 判断选集是否有更新
				if qPeotry.SID != params.SID {
					set, err := models.QueryPeotrySetByID(params.SID)
					if err == nil {
						if set.UID == params.UID {
							qPeotry.SID = params.SID
						} else {
							data[models.STR_CODE] = models.CODE_ERR
							data[models.STR_MSG] = "禁止在他人选集中更新个人诗歌"
							c.respToJSON(data)
							return
						}
					} else {
						data[models.STR_CODE] = models.CODE_ERR
						data[models.STR_MSG] = "未获取到相应选集id"
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
					data[models.STR_DATA] = qPeotry.ID
				} else {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "更新诗歌失败"
				}
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "禁止更新他人诗歌"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "未获取到相应诗歌id"
		}
	}

	c.respToJSON(data)
}

func (c *PeotryController) DeletePeotry() {
	data := c.GetResponseData()
	params := &getDeletePeotryParams{}

	if c.CheckFormParams(data, params) { 
		peotry, err := models.QueryPeotryByID(params.PID)
		if err == nil {
			if peotry.UID == params.UID {
				err := models.DeletePeotry(params.PID)
				if err == nil {
					data[models.STR_DATA] = true
				} else {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "删除诗歌失败"
				}
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "禁止删除他人诗歌"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "未获取到相应诗歌id"
		}
	}

	c.respToJSON(data)
}