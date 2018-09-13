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

// // URLMapping ...
// func (c *PeotryController) URLMapping() {
// 	c.Mapping("Post", c.Post)
// 	c.Mapping("GetOne", c.GetOne)
// 	c.Mapping("GetAll", c.GetAll)
// 	c.Mapping("Put", c.Put)
// 	c.Mapping("Delete", c.Delete)
// }

// // Post ...
// // @Title Post
// // @Description create Peotry
// // @Param	body		body 	models.Peotry	true		"body for Peotry content"
// // @Success 201 {int} models.Peotry
// // @Failure 403 body is empty
// // @router / [post]
// func (c *PeotryController) Post() {
// 	data := c.GetResponseData()
// 	params := &GetUserParams{}

// 	if c.CheckFormParams(data, params) {
// 		claims, errToken := c.ParseUserToken(params.Token)
// 		if errToken == nil {
// 			var v models.Peotry

// 			if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
// 				uIdStr := claims["uid"].(string)
// 				uId, _ := strconv.ParseInt(uIdStr, 10, 64)
				
// 				v.UId = &models.User{Id: uId}
// 				v.PTime = time.Now()
// 				v.Id = v.PTime.UnixNano() / 1e3
// 				v.IId = &models.Peotryimage{Id: 0}
// 				fmt.Println(v)

// 				if _, err := models.AddPeotry(&v); err == nil {
// 					data[models.RESP_DATA] = v
// 				} else {
// 					data[models.RESP_CODE] = models.RESP_ERR
// 					data[models.RESP_MSG] = err.Error()
// 				}
// 			} else {
// 				data[models.RESP_CODE] = models.RESP_ERR
// 				data[models.RESP_MSG] = err.Error()
// 			}
// 		} else {
// 			data[models.RESP_CODE] = models.RESP_ERR
// 			data[models.RESP_MSG] = errToken.Error()
// 		}
// 	}
	
// 	c.respToJSON(data)
// }

// // GetOne ...
// // @Title Get One
// // @Description get Peotry by id
// // @Param	id		path 	string	true		"The key for staticblock"
// // @Success 200 {object} models.Peotry
// // @Failure 403 :id is empty
// // @router /:id [get]
// func (c *PeotryController) GetOne() {
// 	idStr := c.Ctx.Input.Param(":id")
// 	id, _ := strconv.ParseInt(idStr, 10, 64)
// 	v, err := models.GetPeotryById(id)
// 	data := c.GetResponseData()

// 	if err != nil {
// 		data[models.RESP_CODE] = models.RESP_ERR
// 		data[models.RESP_MSG] = err.Error()
// 	} else {
// 		data[models.RESP_DATA] = v
// 	}
// 	c.respToJSON(data)
// }

// // GetAll ...
// // @Title Get All
// // @Description get Peotry
// // @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// // @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// // @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// // @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// // @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// // @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// // @Success 200 {object} models.Peotry
// // @Failure 403
// // @router / [get]
// func (c *PeotryController) GetAll() {
// 	var fields []string
// 	var sortby []string
// 	var order []string
// 	var query = make(map[string]string)
// 	var limit int64 = 10
// 	var offset int64
// 	data := c.GetResponseData()

// 	// fields: col1,col2,entity.col3
// 	if v := c.GetString("fields"); v != "" {
// 		fields = strings.Split(v, ",")
// 	}
// 	// limit: 10 (default is 10)
// 	if v, err := c.GetInt64("limit"); err == nil {
// 		limit = v
// 	}
// 	// offset: 0 (default is 0)
// 	if v, err := c.GetInt64("offset"); err == nil {
// 		offset = v
// 	}
// 	// sortby: col1,col2
// 	if v := c.GetString("sortby"); v != "" {
// 		sortby = strings.Split(v, ",")
// 	}
// 	// order: desc,asc
// 	if v := c.GetString("order"); v != "" {
// 		order = strings.Split(v, ",")
// 	}
// 	// query: k:v,k:v
// 	if v := c.GetString("query"); v != "" {
// 		for _, cond := range strings.Split(v, ",") {
// 			kv := strings.SplitN(cond, ":", 2)
// 			if len(kv) != 2 {
// 				data[models.RESP_CODE] = models.RESP_ERR
// 				data[models.RESP_MSG] = "Error: invalid query key/value pair"
// 				c.respToJSON(data)
// 				return
// 			}
// 			k, v := kv[0], kv[1]
// 			query[k] = v
// 		}
// 	}

// 	l, err := models.GetAllPeotry(query, fields, sortby, order, offset, limit)
	
// 	if err != nil {
// 		data[models.RESP_CODE] = models.RESP_ERR
// 		data[models.RESP_MSG] = err.Error()
// 	} else {
// 		data[models.RESP_DATA] = l
// 	}
// 	c.respToJSON(data)
// }

// // Put ...
// // @Title Put
// // @Description update the Peotry
// // @Param	id		path 	string	true		"The id you want to update"
// // @Param	body		body 	models.Peotry	true		"body for Peotry content"
// // @Success 200 {object} models.Peotry
// // @Failure 403 :id is not int
// // @router /:id [put]
// func (c *PeotryController) Put() {
// 	data := c.GetResponseData()
// 	params := &GetPeotryUpdateParams{}

// 	if c.CheckFormParams(data, params) {
// 		claims, errToken := c.ParseUserToken(params.Token)

// 		if errToken == nil {
// 			pId, _ := strconv.ParseInt(params.PId, 10, 64)
// 			rV, rErr := models.GetPeotryById(pId)
// 			if rErr != nil {
// 				data[models.RESP_CODE] = models.RESP_ERR
// 				data[models.RESP_MSG] = rErr.Error()
// 			} else {
// 				uIdStr := claims["uid"].(string)
// 				uId, _ := strconv.ParseInt(uIdStr, 10, 64)

// 				if rV.UId.Id == uId {
// 					var v models.Peotry
// 					if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
// 						if err := models.UpdatePeotryById(&v); err == nil {
// 							data[models.RESP_MSG] = "更新成功"
// 						} else {
// 							data[models.RESP_CODE] = models.RESP_ERR
// 							data[models.RESP_MSG] = err.Error()
// 						}
// 					} else {
// 						data[models.RESP_CODE] = models.RESP_ERR
// 						data[models.RESP_MSG] = err.Error()
// 					}
// 				} else {
// 					data[models.RESP_CODE] = models.RESP_ERR
// 					data[models.RESP_MSG] = "该诗不属于当前用户"
// 				}
// 			}
// 		} else {
// 			data[models.RESP_CODE] = models.RESP_ERR
// 			data[models.RESP_MSG] = errToken.Error()
// 		}
// 	}
	
// 	c.respToJSON(data)
// }

// // Delete ...
// // @Title Delete
// // @Description delete the Peotry
// // @Param	id		path 	string	true		"The id you want to delete"
// // @Success 200 {string} delete success!
// // @Failure 403 id is empty
// // @router /:id [delete]
// func (c *PeotryController) Delete() {
// 	data := c.GetResponseData()
// 	params := &GetPeotryUpdateParams{}

// 	if c.CheckFormParams(data, params) {
// 		claims, errToken := c.ParseUserToken(params.Token)

// 		if errToken == nil {
// 			pId, _ := strconv.ParseInt(params.PId, 10, 64)
// 			rV, rErr := models.GetPeotryById(pId)
// 			if rErr != nil {
// 				data[models.RESP_CODE] = models.RESP_ERR
// 				data[models.RESP_MSG] = rErr.Error()
// 			} else {
// 				uIdStr := claims["uid"].(string)
// 				uId, _ := strconv.ParseInt(uIdStr, 10, 64)

// 				if rV.UId.Id == uId {
// 					if err := models.DeletePeotry(pId); err == nil {
// 						data[models.RESP_MSG] = "删除成功"
// 					} else {
// 						data[models.RESP_CODE] = models.RESP_ERR
// 						data[models.RESP_MSG] = err.Error()
// 					}
// 				} else {
// 					data[models.RESP_CODE] = models.RESP_ERR
// 					data[models.RESP_MSG] = "该诗不属于当前用户"
// 				}
// 			}
// 		} else {
// 			data[models.RESP_CODE] = models.RESP_ERR
// 			data[models.RESP_MSG] = errToken.Error()
// 		}
// 	}
	
// 	c.respToJSON(data)
// }
