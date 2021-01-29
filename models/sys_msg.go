package models

import (
	"go-sghen/helper"
	"time"
)

// SysMsg ...
type SysMsg struct {
	ID     int64 `gorm:"primary_key,id" json:"id,omitempty"`
	UserID int64 `gorm:"index;column:user_id" json:"userId"`

	MsgType int    `gorm:"index;column:msg_type" json:"msgType"`
	Status  int    `gorm:"index;column:status" json:"status"`
	Content string `gorm:"column:content;type:mediumtext" json:"content"`

	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName ...
func (c SysMsg) TableName() string {
	return "sys_msg"
}

// CreateSysMsg ...
func CreateSysMsg(userID int64, msgType int, content string) (*SysMsg, error) {
	createTime := time.Now()
	sysMsg := &SysMsg{
		ID:         helper.NewUinqueID(),
		UserID:     userID,
		MsgType:    msgType,
		Status:     -1,
		Content:    content,
		CreateTime: createTime,
		UpdateTime: createTime,
	}

	err := dbOrmDefault.Model(&SysMsg{}).Create(sysMsg).Error
	if err != nil {
		return nil, err
	}

	return sysMsg, nil
}

// GetSysMsg ...
func GetSysMsg(id, userId int64) (*SysMsg, error) {
	sysMsg := &SysMsg{
		ID:     id,
		UserID: userId,
	}
	err := dbOrmDefault.Model(&SysMsg{}).Where(sysMsg).Find(sysMsg).Error
	if err == nil {
		return sysMsg, nil
	}
	return nil, err
}

// GetSysMsgs ...
func GetSysMsgs(userId int64, status int, page int, limit int) ([]*SysMsg, int, error) {
	list := make([]*SysMsg, 0)
	count := 0

	if limit == 0 {
		limit = 10
	}

	db := dbOrmDefault.Model(&SysMsg{})
	if userId > 0 {
		query := &SysMsg{
			UserID: userId,
		}
		db = db.Where(query)
	}

	db = db.Order("time_create desc")
	db.Count(&count)

	err := db.Limit(limit).Offset(helper.PageOffset(limit, page)).Find(&list).Error

	if err == nil {
		return list, count, nil
	}
	return nil, 0, err
}

// UpdateSysMsgStatus ...
func UpdateSysMsgStatus(id int64, userId int64, status int) error {
	sysMsg := &SysMsg{
		ID:         id,
		UserID:     userId,
		Status:     status,
		UpdateTime: time.Now(),
	}

	err := dbOrmDefault.Model(&SysMsg{}).Update(sysMsg).Error
	return err
}

// DeleteSysMsg ...
func DeleteSysMsg(id int64) error {
	sysMsg := &SysMsg{
		ID: id,
	}
	err := dbOrmDefault.Model(&SysMsg{}).Where(sysMsg).Delete(sysMsg).Error
	return err
}
