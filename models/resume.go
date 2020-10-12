package models

import (
	"go-sghen/helper"
	"time"
)

// Resume ...
type Resume struct {
	ID     int64 `gorm:"primary_key,id" json:"id,omitempty"`
	UserID int64 `gorm:"index;column:user_id" json:"userId"`

	PersonalInfos string `gorm:"column:personal_infos;type:mediumtext" json:"personalInfos"`
	SkillJob      string `gorm:"column:skill_job;type:mediumtext" json:"skillJob"`
	Educations    string `gorm:"column:educations;type:mediumtext" json:"educations"`
	Experiences   string `gorm:"column:experiences;type:mediumtext" json:"experiences"`
	Projects      string `gorm:"column:projects;type:mediumtext" json:"projects"`
	Descriptions  string `gorm:"column:descriptions;type:mediumtext" json:"descriptions"`
	Hobby         string `gorm:"column:hobby;type:mediumtext" json:"hobby"`

	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName ...
func (c Resume) TableName() string {
	return "resume"
}

// CreateResume ...
func CreateResume(userID int64, personalInfos, skillJob, educations, experiences, projects, descriptions, hobby string) (*Resume, error) {
	createTime := time.Now()
	resume := &Resume{
		ID:            helper.NewUinqueID(),
		UserID:        userID,
		PersonalInfos: personalInfos,
		SkillJob:      skillJob,
		Educations:    educations,
		Experiences:   experiences,
		Projects:      projects,
		Descriptions:  descriptions,
		Hobby:         hobby,
		CreateTime:    createTime,
		UpdateTime:    createTime,
	}

	err := dbOrmDefault.Model(&Resume{}).Create(resume).Error
	if err != nil {
		return nil, err
	}

	return resume, nil
}

// GetResume ...
func GetResume(userID int64) (*Resume, error) {
	resume := &Resume{
		UserID: userID,
	}
	err := dbOrmDefault.Model(&Resume{}).Find(resume).Error
	if err == nil {
		return resume, nil
	}
	return nil, err
}

// UpdateResume ...
func UpdateResume(userID int64, personalInfos, skillJob, educations, experiences, projects, descriptions, hobby string) error {
	resume := &Resume{
		UserID:     userID,
		UpdateTime: time.Now(),
	}
	if len(personalInfos) != 0 {
		resume.PersonalInfos = personalInfos
	}
	if len(skillJob) != 0 {
		resume.SkillJob = skillJob
	}
	if len(educations) != 0 {
		resume.Educations = educations
	}
	if len(experiences) != 0 {
		resume.Experiences = experiences
	}
	if len(projects) != 0 {
		resume.Projects = projects
	}
	if len(descriptions) != 0 {
		resume.Descriptions = descriptions
	}
	if len(hobby) != 0 {
		resume.Hobby = hobby
	}

	err := dbOrmDefault.Model(&Resume{}).Update(resume).Error
	return err
}

// DeleteResume ...
func DeleteResume(userID int64) error {
	resume := &Resume{
		UserID: userID,
	}
	err := dbOrmDefault.Model(&Resume{}).Delete(resume).Error
	return err
}
