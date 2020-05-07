package controllers

import (
	"go-sghen/helper"
	"go-sghen/models"

	"encoding/json"
	"strings"
)

// PeotryController operations for Peotry
type PeotryController struct {
	BaseController
}

// AddTempPeotry 添加系统临时诗词
func (c *PeotryController) AddTempPeotry() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	userLevel := c.Ctx.Input.GetData("level").(int)

	if userLevel < 9 {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户权限不够"
	} else {
		models.AddTempPeotry()
	}
	c.respToJSON(data)
}

func (c *PeotryController) QueryPeotry() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getQueryPeotryParams{}

	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	if params.ID > 0 {
		peotry, err := models.QueryPeotryByID(params.ID)

		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "未查询到对应id的诗歌"
			c.respToJSON(data)
			return
		}

		if params.NeedComment {
			comments, err := models.QueryCommentByTypeID(peotry.ID)
			if err != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "诗歌评论查询失败"
				c.respToJSON(data)
				return
			}
			peotry.Comments = comments
		}
		data[models.STR_DATA] = peotry
		c.respToJSON(data)
		return
	} 

	// 诗词列表查询
	list, count, totalPage, curPage, pageIsEnd, err := models.QueryPeotry(params.UserID, params.SetID, params.Page, params.Limit, params.Content)

	if err == nil {
		if params.NeedComment {
			for _, peotry := range list {
				comments, e := models.QueryCommentByTypeID(peotry.ID)
				if e == nil {
					peotry.Comments = comments
				}
			}
		}

		data[models.STR_DATA] = list
		data["totalCount"] = count
		data["totalPage"] = totalPage
		data["curPage"] = curPage
		data["pageIsEnd"] = pageIsEnd
	} else {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "未查询到对应的诗歌"
	}
	
	c.respToJSON(data)
}

func (c *PeotryController) QueryPopularPeotry() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	limit, err := c.GetInt("limit", 5)
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "参数错误"
		c.respToJSON(data)
		return
	}
	if limit > 20 {
		limit = 20
	}
	list, err := models.QueryPopularPeotry(limit)
	if err == nil {
		for _, peotry := range list {
			comments, e := models.QueryCommentByTypeID(peotry.ID)
			if e == nil {
				peotry.Comments = comments
			}
		}
		data[models.STR_DATA] = list
	}
	c.respToJSON(data)
}

// CreatePeotry ...
func (c *PeotryController) CreatePeotry() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getCreatePeotryParams{}

	if c.CheckFormParams(data, params) {
		set, err := models.QueryPeotrySetByID(params.SetID)

		if err == nil {
			if set.UserID == 0 || set.UserID == params.UserID {
				imageNames := make([]string, 0)

				// 判断是否有图片
				if len(strings.TrimSpace(params.ImageNames)) > 0 {
					err := json.Unmarshal([]byte(params.ImageNames), &imageNames)
					if err == nil {
						if len(imageNames) > 10 {
							data[models.STR_MSG] = "诗歌图片超过10张，只保存前10张"
							imageNames = imageNames[0:10]
						}
					} else {
						data[models.STR_MSG] = "诗词图片列表解析失败"
					}
				}

				timeStr := helper.GetNowDateTime()
				pId, err := models.CreatePeotry(params.UserID, params.SetID, params.Title, timeStr, params.Content, params.End, params.ImageNames)

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

// UpdatePeotry ...
func (c *PeotryController) UpdatePeotry() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getUpdatePeotryParams{}

	if c.CheckFormParams(data, params) {
		qPeotry, err := models.QueryPeotryByID(params.ID)

		if err == nil {
			if qPeotry.UserID == params.UserID {
				// 判断选集是否有更新
				if qPeotry.SetID != params.SetID {
					set, err := models.QueryPeotrySetByID(params.SetID)
					if err == nil {
						if set.UserID == 0 || set.UserID == params.UserID {
							qPeotry.SetID = params.SetID
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
				qPeotry.Title = params.Title
				qPeotry.Content = params.Content
				qPeotry.End = params.End
				// 更新时需要将这些附带的结构体置空
				qPeotry.User = nil
				qPeotry.Set = nil
				qPeotry.Image = nil

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

// DeletePeotry ...
func (c *PeotryController) DeletePeotry() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getDeletePeotryParams{}

	if c.CheckFormParams(data, params) {
		peotry, err := models.QueryPeotryByID(params.ID)
		userID := c.Ctx.Input.GetData("userId").(int64)

		if err == nil {
			if peotry.UserID == userID {
				err := models.DeletePeotry(peotry.ID)
				models.DeleteComments(peotry.ID)

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
