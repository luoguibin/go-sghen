package controllers

import(
	"SghenApi/models"
	"fmt"
	// "encoding/json"
)

type HtmlController struct {
	BaseController
}

type getHtmlInput struct {
	ID		int64		`form:"id"`
}

func (c *HtmlController) Get() {
	fmt.Println("html get")
	params := &getHtmlInput{}
	if err := c.ParseForm(params); err == nil {
		if params.ID == 0 {
			params.ID = -1
		}

		user, err := models.QueryUser(params.ID)
		if err == nil {
			c.Data["user"] = user
		} else {
			c.Data["user"] = models.User {
				ID:			-1,
				UName:		"游客",
			}
		}
	}
	
	c.TplName = "index.html"
}