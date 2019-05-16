package controllers

import (
	"fmt"
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
		if params.ToID < 0 {
			fmt.Println(params)
			comment, _ := models.QueryCommentByTypeIDFromID(params.TypeID, params.FromID, params.ToID)
			fmt.Println(comment)
			if comment == nil {
				comment = &models.Comment{
					Type:    params.Type,
					TypeID:  params.TypeID,
					FromID:  params.FromID,
					ToID:    params.ToID,
					Content: params.Comment,
				}
				comment, err := models.CreateComment(params.Type, params.TypeID, params.FromID, params.ToID, params.Comment)
				if err == nil {
					data[models.STR_DATA] = comment.ID
				} else {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "操作失败"
				}
			} else {
				comment.Content = params.Comment
				err := models.SaveComment(comment)
				if err != nil {
					data[models.STR_CODE] = models.CODE_ERR
					data[models.STR_MSG] = "操作失败"
				}
			}
		} else {
			comment, err := models.CreateComment(params.Type, params.TypeID, params.FromID, params.ToID, params.Comment)
			if err == nil {
				data[models.STR_DATA] = comment.ID
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "评论提交失败"
			}
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
	c.respToJSON(data)
}

// DeleteComment ...
func (c *CommentController) DeleteComment() {
	data := c.GetResponseData()
	params := &getDeleteCommentParams{}
	if c.CheckFormParams(data, params) {
		comment, err := models.QueryComment(params.ID)
		userID := c.Ctx.Input.GetData("uId").(int64)
		if err == nil {
			if userID == comment.FromID {
				err := models.DeleteComemnt(comment.ID)
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
