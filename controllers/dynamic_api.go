package controllers

import (
	"go-sghen/models"
	"regexp"
	"strings"
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
			if params.Status == 0 {
				params.Status = -1
			}
			dynamicAPI, err := models.CreateDynamicAPI(params.SuffixPath, params.Name, params.Comment, params.Content, params.Status, userID)
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
			dynamicAPI, err := models.UpdateDynamicAPI(params.ID, params.SuffixPath, params.Name, params.Comment, params.Content, params.Status)
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
		list, count, totalPage, curPage, pageIsEnd, err := models.QueryDynamicAPI(params.ID, params.SuffixPath, params.Name, params.Comment, params.Status, params.UserID, params.Limit, params.Page)
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
		userLevel := c.Ctx.Input.GetData("level").(int)

		if userLevel < 9 {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "用户权限不够，禁止更新接口"
		} else {
			err := models.DeleteDynamicAPI(params.ID)
			if err != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "删除接口失败"
			}
		}
	}

	c.respToJSON(data)
}

// GetDynamicDataByPath 获取数据
func (c *DynamicAPIController) GetDynamicDataByPath() {
	data := c.GetResponseData()

	suffixPath := c.Ctx.Input.Param(":splat")
	// pathKeys := strings.Split(splat, "/")
	dynamicAPI, ok := models.MConfig.DynamicAPIMap[suffixPath]
	// fmt.Println("DynamicAPI", suffixPath, ok, dynamicAPI, c.Ctx.Request.URL)
	if ok {
		if dynamicAPI.Status == 1 {
			sqlStr := dynamicAPI.Content
			r, _ := regexp.Compile("\\$\\{[0-9a-zA-Z_]{1,}\\}")
			keyNames := r.FindAllStringSubmatch(sqlStr, -1)
			if len(keyNames) > 0 {
				for _, keyName0 := range keyNames {
					keyName := keyName0[0]
					orderName := keyName[2 : len(keyName)-1]

					switch orderName {
					case "limit":
						limit := c.GetString("limit", "20")
						sqlStr = strings.Replace(sqlStr, "${limit}", limit, -1)
					case "offset":
						offset := c.GetString("offset", "0")
						sqlStr = strings.Replace(sqlStr, "${offset}", offset, -1)
					case "id":
						id := c.GetString("id", "0")
						if id != "0" {
							sqlStr = strings.Replace(sqlStr, "${id}", id, -1)
						}
					case "datas":
						datas := c.GetString("datas", "")
						r, _ := regexp.Compile("[0-9,]+")
						if len(datas) > 0 && r.MatchString(datas) {
							// datasStr := strings.Join(datas,  ",")
							sqlStr = strings.Replace(sqlStr, "${datas}", datas, -1)
						} else {
							data[models.STR_CODE] = models.CODE_ERR
							data[models.STR_MSG] = "操作失败"
							c.respToJSON(data)
							return
						}
					}
				}
			}

			// fmt.Println(sqlStr)

			list, err := models.GetDynamicData(sqlStr)
			if err != nil {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "操作失败"
				data[models.STR_DETAIL] = err
			} else {
				data[models.STR_DATA] = list
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "接口未加载"
		}
	} else {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "接口未加载或未定义"
	}

	c.respToJSON(data)
}

// PostDynamicData 更改数据
func (c *DynamicAPIController) PostDynamicData() {
	data := c.GetResponseData()

	c.respToJSON(data)
}
