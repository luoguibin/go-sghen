package controllers

import (
	"go-sghen/models"
)

// CommentController ...
type CommentController struct {
	BaseController
}

// CreateComment ...
func (c *CommentController) CreateComment() {
	data := c.GetResponseData()
	params := &getCreateCommentParams{}
	if c.CheckPostParams(data, params) {
		comment, err := models.CreateComment(params.Type, params.TypeID, params.FromID, params.ToID, params.Comment)
		if err == nil {
			data[models.STR_DATA] = comment.ID
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "评论提交失败"
		}
	}

	c.respToJSON(data)
}

// QueryComments ...
func (c *CommentController) QueryComments() {
	data := c.GetResponseData()
	params := &getQueryCommentParams{}
	if c.CheckFormParams(data, params) {
		comments, err := models.QueryCommentByTypeID(params.TypeID)
		if err == nil {
			data[models.STR_DATA] = comments
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "评论提交失败"
		}
	}
}
