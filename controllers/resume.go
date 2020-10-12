package controllers

import "go-sghen/models"

// ResumeController ...
type ResumeController struct {
	BaseController
}

// CreateResume ...
func (c *ResumeController) CreateResume() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getCreateResumeParams{}

	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	resume, err := models.CreateResume(userID, params.PersonalInfos, params.SkillJob, params.Educations, params.Experiences, params.Projects, params.Descriptions, params.Hobby)

	if err == nil {
		data[models.STR_DATA] = resume
	} else {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "创建简历失败"
	}

	c.respToJSON(data)
}

// GetResumeDetail ...
func (c *ResumeController) GetResumeDetail() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	resume, err := models.GetResume(userID)

	if err == nil {
		data[models.STR_DATA] = resume
	} else {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "未查询到简历"
	}

	c.respToJSON(data)
}

// UpdateResume ...
func (c *ResumeController) UpdateResume() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}
	params := &getCreateResumeParams{}

	if !c.CheckFormParams(data, params) {
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	err := models.UpdateResume(userID, params.PersonalInfos, params.SkillJob, params.Educations, params.Experiences, params.Projects, params.Descriptions, params.Hobby)

	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "更新简历失败"
	}

	c.respToJSON(data)
}

// DeleteResume ...
func (c *ResumeController) DeleteResume() {
	data, isOk := c.GetResponseData()
	if !isOk {
		c.respToJSON(data)
		return
	}

	userID := c.Ctx.Input.GetData("userId").(int64)
	err := models.DeleteResume(userID)

	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "删除简历失败"
	}

	c.respToJSON(data)
}
