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
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	params := &getCreateCommentParams{}

	if c.CheckFormParams(data, params) {
		comment, err := models.CreateComment(params.Type, params.TypeID, params.FromID, params.ToID, params.Content)

		if err == nil {
			data[models.STR_DATA] = comment
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "评论提交失败"
		}
	}

	c.respToJSON(data)
}

// QueryComments ...
func (c *CommentController) QueryComments() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getQueryCommentParams{}

	if c.CheckFormParams(data, params) {
		comments, err := models.QueryCommentByTypeID(params.TypeID)

		if err == nil {
			data[models.STR_DATA] = comments
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "评论查询失败"
		}
	}
	c.respToJSON(data)
}

// DeleteComment ...
func (c *CommentController) DeleteComment() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getDeleteCommentParams{}

	if c.CheckFormParams(data, params) {
		comment, err := models.QueryComment(params.ID)
		userID := c.Ctx.Input.GetData("userId").(int64)

		if err == nil {
			if userID == comment.FromID {
				err := models.DeleteComment(comment.ID)

				if err != nil {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "删除评论失败"
				}
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "只能删除自己的评论"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "评论查询失败"
		}
	}
	c.respToJSON(data)
}
