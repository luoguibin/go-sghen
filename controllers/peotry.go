package controllers

import (
	"SghenApi/models"
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"fmt"
)

// PeotryController operations for Peotry
type PeotryController struct {
	BaseController
}

// URLMapping ...
func (c *PeotryController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Peotry
// @Param	body		body 	models.Peotry	true		"body for Peotry content"
// @Success 201 {int} models.Peotry
// @Failure 403 body is empty
// @router / [post]
func (c *PeotryController) Post() {
	data := c.GetResponseData()
	params := &GetUserParams{}

	if c.CheckFormParams(data, params) {
		claims, err := c.ParseUserToken(params.Token)

		if err == nil {
			fmt.Println(claims)
			var v models.Peotry

			if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
				v.UId = &models.User{Id: 15625045984}
				v.PTime = time.Now()
				idStr := v.PTime.Format("20060102") + "01"
				v.Id, _ = strconv.ParseInt(idStr, 10, 64)
				fmt.Println(v.PTime)

				if _, err := models.AddPeotry(&v); err == nil {
					data[models.RESP_DATA] = v
				} else {
					data[models.RESP_CODE] = models.RESP_ERR
					data[models.RESP_MSG] = err.Error()
				}
			} else {
				data[models.RESP_CODE] = models.RESP_ERR
				data[models.RESP_MSG] = err.Error()
			}
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = "token invalid"
		}
	}
	
	c.respToJSON(data)
}

// GetOne ...
// @Title Get One
// @Description get Peotry by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Peotry
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PeotryController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	v, err := models.GetPeotryById(id)
	data := c.GetResponseData()

	if err != nil {
		data[models.RESP_CODE] = models.RESP_ERR
		data[models.RESP_MSG] = err.Error()
	} else {
		data[models.RESP_DATA] = v
	}
	c.respToJSON(data)
}

// GetAll ...
// @Title Get All
// @Description get Peotry
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Peotry
// @Failure 403
// @router / [get]
func (c *PeotryController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64
	data := c.GetResponseData()

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				data[models.RESP_CODE] = models.RESP_ERR
				data[models.RESP_MSG] = "Error: invalid query key/value pair"
				c.respToJSON(data)
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllPeotry(query, fields, sortby, order, offset, limit)
	
	if err != nil {
		data[models.RESP_CODE] = models.RESP_ERR
		data[models.RESP_MSG] = err.Error()
	} else {
		data[models.RESP_DATA] = l
	}
	c.respToJSON(data)
}

// Put ...
// @Title Put
// @Description update the Peotry
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Peotry	true		"body for Peotry content"
// @Success 200 {object} models.Peotry
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PeotryController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	v := models.Peotry{Id: id}
	data := c.GetResponseData()
	
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdatePeotryById(&v); err == nil {
			data[models.RESP_MSG] = "OK"
		} else {
			data[models.RESP_CODE] = models.RESP_ERR
			data[models.RESP_MSG] = err.Error()
		}
	} else {
		data[models.RESP_CODE] = models.RESP_ERR
		data[models.RESP_MSG] = err.Error()
	}
	c.respToJSON(data)
}

// Delete ...
// @Title Delete
// @Description delete the Peotry
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PeotryController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	data := c.GetResponseData()
	if err := models.DeletePeotry(id); err == nil {
		data[models.RESP_MSG] = "删除成功"
	} else {
		data[models.RESP_CODE] = models.RESP_ERR
		data[models.RESP_MSG] = err.Error()
	}
	c.respToJSON(data)
}
