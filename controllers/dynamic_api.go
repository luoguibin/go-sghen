package controllers

import (
	"go-sghen/models"
)

// DynamicAPIController 自定义接口控制器
type DynamicAPIController struct {
	BaseController
}

// CreateDynamicAPI 新增
func (c *DynamicAPIController) CreateDynamicAPI() {
	data := c.GetResponseData()
	params := &getCreateDynamicAPIParams{}

	if c.CheckFormParams(data, params) {
		userID := c.Ctx.Input.GetData("userId").(int64)
		userLevel := c.Ctx.Input.GetData("level").(int)

		if userLevel < 9 {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "用户权限不够，禁止创建接口"
		} else {
			dynamicAPI, err := models.CreateDynamicAPI(params.Name, params.Comment, params.Content, params.Status, userID)
			if err != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "创建接口失败"
			} else {
				data[models.STR_DATA] = dynamicAPI
			}
		}
	}

	c.respToJSON(data)
}

// UpdateDynamicAPI 更新
func (c *DynamicAPIController) UpdateDynamicAPI() {
	data := c.GetResponseData()
	params := &getUpdateDynamicAPIParams{}

	if c.CheckFormParams(data, params) {
		// userID := c.Ctx.Input.GetData("userId").(int64)
		userLevel := c.Ctx.Input.GetData("level").(int)

		if userLevel < 9 {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "用户权限不够，禁止更新接口"
		} else {
			dynamicAPI, err := models.UpdateDynamicAPI(params.ID, params.Name, params.Comment, params.Content, params.Status)
			if err != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "更新接口失败"
			} else {
				data[models.STR_DATA] = dynamicAPI
			}
		}
	}

	c.respToJSON(data)
}

// QueryDynamicAPI 查询
func (c *DynamicAPIController) QueryDynamicAPI() {
	data := c.GetResponseData()
	params := &getQueryDynamicAPIParams{}

	if c.CheckFormParams(data, params) {
		list, count, totalPage, curPage, pageIsEnd, err := models.QueryDynamicAPI(params.ID, params.Name, params.Comment, params.Status, params.UserID, params.Limit, params.Page)
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "查询列表失败"
		} else {
			data[models.STR_DATA] = list
			data["totalCount"] = count
			data["totalPage"] = totalPage
			data["curPage"] = curPage
			data["pageIsEnd"] = pageIsEnd
		}
	}

	c.respToJSON(data)
}

// DeleteDynamicAPI 删除
func (c *DynamicAPIController) DeleteDynamicAPI() {
	data := c.GetResponseData()
	params := &getDeleteDynamicAPIParams{}

	if c.CheckFormParams(data, params) {
		// userID := c.Ctx.Input.GetData("userId").(int64)
		// userLevel := c.Ctx.Input.GetData("level").(int)

		err := models.DeleteDynamicAPI(params.ID)
		if err != nil {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "删除接口失败"
		}
	}

	c.respToJSON(data)
}
