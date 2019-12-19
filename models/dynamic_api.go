package models

import (
	"fmt"
	"go-sghen/helper"
	"time"
)

// DynamicAPI 自定义脚本制作接口
type DynamicAPI struct {
	ID int64 `gorm:"primary_key" json:"id,omitempty"`

	Name    string `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Comment string `gorm:"column:comment;type:varchar(200)" json:"comment,omitempty"`
	Content string `gorm:"column:content;type:mediumtext" json:"content"`

	Status int `gorm:"column:status" json:"status"`

	TimeCreate time.Time `gorm:"column:time_create" json:"timeCreate"`
	TimeUpdate time.Time `gorm:"column:time_update" json:"timeUpdate"`

	UserID int64 `gorm:"column:user_id" json:"-"`
	User   *User `gorm:"foreignkey:user_id;" json:"user"`
}

// TableName ...
func (u DynamicAPI) TableName() string {
	return "dynamic_api"
}

// CreateDynamicAPI 创建一个接口
func CreateDynamicAPI(name string, comment string, content string, status int, userID int64) (*DynamicAPI, error) {
	id := helper.GetMicrosecond()
	timeNow := time.Now()

	dynamicAPI := &DynamicAPI{
		ID:         id,
		Name:       name,
		Comment:    comment,
		Content:    content,
		Status:     status,
		UserID:     userID,
		TimeCreate: timeNow,
		TimeUpdate: timeNow,
	}

	err := dbOrmDefault.Model(&DynamicAPI{}).Create(dynamicAPI).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return dynamicAPI, nil
}

// UpdateDynamicAPI 更新一个接口
func UpdateDynamicAPI(id int64, name string, comment string, content string, status int) (*DynamicAPI, error) {
	dynamicAPI := &DynamicAPI{
		ID:         id,
		Name:       name,
		Comment:    comment,
		Content:    content,
		Status:     status,
		TimeUpdate: time.Now(),
	}

	err := dbOrmDefault.Model(&DynamicAPI{}).Update(dynamicAPI).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return dynamicAPI, nil
}

// QueryDynamicAPI 查询接口列表
func QueryDynamicAPI(id int64, name string, comment string, status int, userID int64, limit int, page int) ([]*DynamicAPI, int, int, int, int, error) {
	list := make([]*DynamicAPI, 0)
	totalPage := 0
	count := 0
	curPage := page
	pageIsEnd := 0

	if limit <= 0 {
		limit = 10
	}

	db := dbOrmDefault.Model(&DynamicAPI{})
	query := &DynamicAPI{}
	if id > 0 {
		query.ID = id
	}
	if status > 0 {
		query.Status = status
	}
	if userID > 0 {
		query.UserID = userID
	}
	db = db.Where(query)

	if len(name) > 1 {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if len(comment) > 1 {
		db = db.Where("comment LIKE ?", "%"+comment+"%")
	}

	db.Count(&count)
	db = db.Preload("User")
	err := db.Limit(limit).Offset(helper.PageOffset(limit, page)).Find(&list).Error

	if err == nil {
		totalPage, pageIsEnd = helper.PageTotal(limit, page, int64(count))
		return list, count, totalPage, curPage, pageIsEnd, nil
	}
	return nil, 0, 0, 0, 0, err
}

// DeleteDynamicAPI 删除接口
func DeleteDynamicAPI(id int64) error {
	dynamicAPI := &DynamicAPI{
		ID: id,
	}

	err := dbOrmDefault.Model(&DynamicAPI{}).Delete(&dynamicAPI).Error
	return err
}
