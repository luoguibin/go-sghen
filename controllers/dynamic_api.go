package controllers

import (
	"errors"
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
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getCreateDynamicAPIParams{}

	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	userLevel := c.Ctx.Input.GetData("level").(int)

	if userLevel < 9 {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户权限不够，禁止创建接口"
		c.respToJSON(data)
		return
	}

	// 禁止保存状态0
	if params.Status == 0 {
		params.Status = -1
	}
	dynamicAPI, err := models.CreateDynamicAPI(params.SuffixPath, params.Name, params.Comment, params.Content, params.Status, userID)
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "创建接口失败"
	} else {
		data[models.STR_DATA] = dynamicAPI
		if dynamicAPI.Status > 0 {
			dynamicAPICacheTask()
		}
	}

	c.respToJSON(data)
}

// UpdateDynamicAPI 更新
func (c *DynamicAPIController) UpdateDynamicAPI() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getUpdateDynamicAPIParams{}

	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	originAPI, err := models.QueryOneDynamicAPI(params.ID)
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "接口ID校验失败"
		c.respToJSON(data)
		return
	}
	if userID != originAPI.UserID {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "禁止更新非本人创建的接口"
		c.respToJSON(data)
		return
	}

	userLevel := c.Ctx.Input.GetData("level").(int)
	if userLevel < 9 {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户权限不够，禁止更新接口"
		c.respToJSON(data)
		return
	}

	// 检测是否已加载该api
	cacheAPI0, ok := models.MConfig.DynamicAPIMap.Load(originAPI.SuffixPath)
	if ok {
		cacheAPI := cacheAPI0.(*models.DynamicAPI)
		originAPI.Count = cacheAPI.Count
	}

	dynamicAPI, err := models.UpdateDynamicAPI(params.ID, params.SuffixPath, params.Name, params.Comment, params.Content, params.Status, originAPI.Count)
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "更新接口失败"
	} else {
		data[models.STR_DATA] = dynamicAPI
		dynamicAPICacheTask()
	}

	c.respToJSON(data)
}

// QueryDynamicAPI 查询
func (c *DynamicAPIController) QueryDynamicAPI() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getQueryDynamicAPIParams{}
	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

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

	c.respToJSON(data)
}

// DeleteDynamicAPI 删除
func (c *DynamicAPIController) DeleteDynamicAPI() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getDeleteDynamicAPIParams{}
	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	// userID := c.Ctx.Input.GetData("userId").(int64)
	// userLevel := c.Ctx.Input.GetData("level").(int)
	userLevel := c.Ctx.Input.GetData("level").(int)
	if userLevel < 9 {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户权限不够，禁止更新接口"
		c.respToJSON(data)
		return
	}

	err := models.DeleteDynamicAPI(params.ID)
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "删除接口失败"
	} else {
		dynamicAPICacheTask()
	}

	c.respToJSON(data)
}

// GetDynamicDataByPath 获取数据
func (c *DynamicAPIController) GetDynamicDataByPath() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	dynamicData, err := c.getDynamicData()
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = err.Error()
		c.respToJSON(data)
		return
	}
	data[models.STR_DATA] = dynamicData
	c.respToJSON(data)
}

// PostDynamicData 更改数据
func (c *DynamicAPIController) PostDynamicData() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	c.respToJSON(data)
}

// getDynamicData ...
func (c *DynamicAPIController) getDynamicData() (*interface{}, error) {
	suffixPath := c.Ctx.Input.Param(":splat")
	dynamicAPI0, ok := models.MConfig.DynamicAPIMap.Load(suffixPath)
	if !ok {
		return nil, errors.New("接口未加载或未定义")
	}
	dynamicAPI := (dynamicAPI0).(*models.DynamicAPI)
	if dynamicAPI.Status < 1 {
		return nil, errors.New("接口未加载")
	}

	cacheData, cachedOk := models.MConfig.DynamicCachedMap[suffixPath]
	if dynamicAPI.Status == 2 && cachedOk {
		// 读取缓存数据
		dynamicAPI.Count = dynamicAPI.Count + 1
		return cacheData, nil
	}

	var datas []interface{}
	sqlStrs := strings.Split(dynamicAPI.Content, "$$$")
	r, _ := regexp.Compile("\\$\\{[0-9a-zA-Z_]{1,}\\}")
	for _, str := range sqlStrs {
		sqlStr := str
		keyNames := r.FindAllStringSubmatch(sqlStr, -1)
		if len(keyNames) > 0 {
			for _, keyName0 := range keyNames {
				keyName := keyName0[0]
				orderName := keyName[2 : len(keyName)-1]

				switch orderName {
				case "limit":
					limit := c.GetString("limit", "20")
					r, _ := regexp.Compile("^[0-9]+$")
					if limit != "0" && r.MatchString(limit) {
						sqlStr = strings.Replace(sqlStr, "${limit}", limit, -1)
					} else {
						return nil, errors.New("操作失败")
					}
				case "offset":
					offset := c.GetString("offset", "0")
					r, _ := regexp.Compile("^[0-9]+$")
					if r.MatchString(offset) {
						sqlStr = strings.Replace(sqlStr, "${offset}", offset, -1)
					} else {
						return nil, errors.New("操作失败")
					}
				case "id":
					id := c.GetString("id", "0")
					r, _ := regexp.Compile("^[0-9]+$")
					if id != "0" && r.MatchString(id) {
						sqlStr = strings.Replace(sqlStr, "${id}", id, -1)
					} else {
						return nil, errors.New("操作失败")
					}
				case "datas":
					datas := c.GetString("datas", "")
					r, _ := regexp.Compile("^[0-9,]+$")
					if len(datas) > 0 && r.MatchString(datas) {
						// datasStr := strings.Join(datas,  ",")
						sqlStr = strings.Replace(sqlStr, "${datas}", datas, -1)
					} else {
						return nil, errors.New("操作失败")
					}
				case "date0":
					date0 := c.GetString("date0", "")
					r, _ := regexp.Compile("^[0-9\\-:\\s]+$")
					if len(date0) > 0 && r.MatchString(date0) {
						// datasStr := strings.Join(datas,  ",")
						sqlStr = strings.Replace(sqlStr, "${date0}", date0, -1)
					} else {
						return nil, errors.New("操作失败")
					}
				case "date1":
					date1 := c.GetString("date1", "")
					r, _ := regexp.Compile("^[0-9\\-:\\s]+$")
					if len(date1) > 0 && r.MatchString(date1) {
						// datasStr := strings.Join(datas,  ",")
						sqlStr = strings.Replace(sqlStr, "${date1}", date1, -1)
					} else {
						return nil, errors.New("操作失败")
					}
				default:
					return nil, errors.New("操作失败")
				}
			}
		}

		tempData, err := models.GetDynamicData(sqlStr)
		if err != nil {
			models.MConfig.MLogger.Error(err.Error())
			return nil, errors.New("操作失败")
		}
		datas = append(datas, tempData)
	}

	dynamicAPI.Count = dynamicAPI.Count + 1
	var dynamicData interface{}
	if len(datas) == 1 {
		dynamicData = datas[0]
	} else {
		dynamicData = datas
	}

	if dynamicAPI.Status == 2 && !cachedOk {
		models.MConfig.DynamicCachedMap[suffixPath] = &dynamicData
	}
	return &dynamicData, nil
}

// dynamicAPICacheTask ...
func dynamicAPICacheTask() {
	models.MConfig.MLogger.Info("动态API缓存数据重置")

	models.MConfig.DynamicAPIMap.Range(func(k, v interface{}) bool {
		dynamicAPI, ok := v.(*models.DynamicAPI)
		if ok {
			models.UpdateDynamicAPI(dynamicAPI.ID, "", "", "", "", 0, dynamicAPI.Count)
		}
		return true
	})

	models.InitDynamicAPIMap()
}
