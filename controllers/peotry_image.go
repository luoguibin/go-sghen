package controllers

// import (
// 	"SghenApi/models"
// 	"encoding/json"
// 	"strconv"
// 	"strings"
// 	"fmt"
// )

// // PeotryimageController operations for Peotryimage
// type PeotryimageController struct {
// 	BaseController
// }

// // URLMapping ...
// func (c *PeotryimageController) URLMapping() {
// 	c.Mapping("Post", c.Post)
// 	c.Mapping("GetOne", c.GetOne)
// 	c.Mapping("GetAll", c.GetAll)
// 	c.Mapping("Put", c.Put)
// 	c.Mapping("Delete", c.Delete)
// }

// // Post ...
// // @Title Post
// // @Description create Peotryimage
// // @Param	body		body 	models.Peotryimage	true		"body for Peotryimage content"
// // @Success 201 {int} models.Peotryimage
// // @Failure 403 body is empty
// // @router / [post]
// func (c *PeotryimageController) Post() {
// 	// for update test
// 	// data := c.GetResponseData()
// 	// var v models.Peotryimage
// 	// if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
// 	// 	if _, err := models.AddPeotryimage(&v); err == nil {
// 	// 		data[models.STR_DATA] = v
// 	// 	} else {
// 	// 		data[models.STR_CODE] = models.CODE_ERR
// 	// 				data[models.STR_MSG] = err.Error()
// 	// 	}
// 	// } else {
// 	// 	data[models.STR_CODE] = models.CODE_ERR
// 	// 	data[models.STR_MSG] = err.Error()
// 	// }
// 	// c.respToJSON(data)

// 	data := c.GetResponseData()
// 	params := &GetPeotryimageParams{}

// 	if c.CheckPostParams(data, params) {
// 		claims, errToken := c.ParseUserToken(params.Token)

// 		if errToken == nil {
// 			pId, _ := strconv.ParseInt(params.PId, 10, 64)
// 			rV, rErr := models.GetPeotryById(pId)

// 			if rErr != nil {
// 				data[models.STR_CODE] = models.CODE_ERR
// 				data[models.STR_MSG] = rErr.Error()
// 			} else {
// 				uIdStr := claims["uid"].(string)
// 				uId, _ := strconv.ParseInt(uIdStr, 10, 64)

// 				if rV.UId.Id == uId {
// 					// params.ImageDatas
// 					imagesStr := "["
// 					count := 0
// 					length := len(params.ImageDatas)
// 					if length > 0 {
// 						for index, imageData := range params.ImageDatas {
// 							rename := params.PId + "-" + strconv.Itoa(index + 1)
// 							format, err := models.SavePeotryimage(imageData, rename)
// 							if err == nil {
// 								imagesStr += "\"" + rename + "." + format + "\""
// 								if index != length - 1 {
// 									imagesStr += ","
// 								}
// 								count++
// 							} else {
// 								count = -1
// 								break
// 							}
// 						}

// 						if count >= 0 {
// 							imagesStr += "]"
// 							pId, _ := strconv.ParseInt(params.PId, 10, 64)
// 							v := models.Peotryimage {
// 								Id: pId,
// 								IImages: imagesStr,
// 								ICount: count,
// 							}

// 							if _, err := models.AddPeotryimage(&v); err == nil {
// 								data[models.STR_DATA] = v
// 								rV.IId.Id = pId
// 								vErr := models.UpdatePeotryById(rV)
// 								fmt.Println(vErr)
// 								if vErr == nil {
// 									data[models.STR_MSG] = "创建成功"
// 								} else {
// 									data[models.STR_CODE] = models.CODE_ERR
// 									data[models.STR_MSG] = vErr.Error()
// 								}
// 							} else {
// 								data[models.STR_CODE] = models.CODE_ERR
// 								data[models.STR_MSG] = err.Error()
// 							}
// 						} else {
// 							data[models.STR_CODE] = models.CODE_ERR
// 							data[models.STR_MSG] = "保存图片出错：" + imagesStr
// 						}
// 					} else {
// 						data[models.STR_MSG] = "无图片数据"
// 					}
// 				} else {
// 					data[models.STR_CODE] = models.CODE_ERR
// 					data[models.STR_MSG] = "该诗不属于当前用户"
// 				}
// 			}
// 		} else {
// 			data[models.STR_CODE] = models.CODE_ERR
// 			data[models.STR_MSG] = errToken.Error()
// 		}	
// 	}
// 	c.respToJSON(data)
// }

// // GetOne ...
// // @Title Get One
// // @Description get Peotryimage by id
// // @Param	id		path 	string	true		"The key for staticblock"
// // @Success 200 {object} models.Peotryimage
// // @Failure 403 :id is empty
// // @router /:id [get]
// func (c *PeotryimageController) GetOne() {
// 	idStr := c.Ctx.Input.Param(":id")
// 	id, _ := strconv.ParseInt(idStr, 10, 64)
// 	v, err := models.GetPeotryimageById(id)
// 	data := c.GetResponseData()
			
// 	if err != nil {
// 		data[models.STR_CODE] = models.CODE_ERR
// 		data[models.STR_MSG] = err.Error()
// 	} else {
// 		data[models.STR_DATA] = v
// 	}	
// 	c.respToJSON(data)
// }

// // GetAll ...
// // @Title Get All
// // @Description get Peotryimage
// // @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// // @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// // @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// // @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// // @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// // @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// // @Success 200 {object} models.Peotryimage
// // @Failure 403
// // @router / [get]
// func (c *PeotryimageController) GetAll() {
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
// 				data[models.STR_CODE] = models.CODE_ERR
// 				data[models.STR_MSG] = "Error: invalid query key/value pair"
// 				c.respToJSON(data)
// 				return
// 			}
// 			k, v := kv[0], kv[1]
// 			query[k] = v
// 		}
// 	}

// 	l, err := models.GetAllPeotryimage(query, fields, sortby, order, offset, limit)
// 	if err != nil {
// 		data[models.STR_CODE] = models.CODE_ERR
// 		data[models.STR_MSG] = err.Error()
// 	} else {
// 		data[models.STR_DATA] = l
// 	}
// 	c.respToJSON(data)
// }

// // Put ...
// // @Title Put
// // @Description update the Peotryimage
// // @Param	id		path 	string	true		"The id you want to update"
// // @Param	body		body 	models.Peotryimage	true		"body for Peotryimage content"
// // @Success 200 {object} models.Peotryimage
// // @Failure 403 :id is not int
// // @router /:id [put]
// func (c *PeotryimageController) Put() {
// 	idStr := c.Ctx.Input.Param(":id")
// 	id, _ := strconv.ParseInt(idStr, 10, 64)
// 	v := models.Peotryimage{Id: id}
// 	data := c.GetResponseData()

// 	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
// 		if err := models.UpdatePeotryimageById(&v); err == nil {
// 			data[models.STR_MSG] = "OK"
// 		} else {
// 			data[models.STR_CODE] = models.CODE_ERR
// 			data[models.STR_MSG] = err.Error()
// 		}
// 	} else {
// 		data[models.STR_CODE] = models.CODE_ERR
// 		data[models.STR_MSG] = err.Error()
// 	}
// 	c.respToJSON(data)
// }

// // Delete ...
// // @Title Delete
// // @Description delete the Peotryimage
// // @Param	id		path 	string	true		"The id you want to delete"
// // @Success 200 {string} delete success!
// // @Failure 403 id is empty
// // @router /:id [delete]
// func (c *PeotryimageController) Delete() {
// 	idStr := c.Ctx.Input.Param(":id")
// 	id, _ := strconv.ParseInt(idStr, 10, 64)
// 	data := c.GetResponseData()
	
// 	if err := models.DeletePeotryimage(id); err == nil {
// 		data[models.STR_MSG] = "OK"
// 	} else {
// 		data[models.STR_CODE] = models.CODE_ERR
// 		data[models.STR_MSG] = err.Error()
// 	}
// 	c.respToJSON(data)
// }
