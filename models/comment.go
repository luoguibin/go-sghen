package models

import (
	"go-sghen/helper"
	"time"
)

// Comment ...
type Comment struct {
	ID         int64     `gorm:"primary_key,id" json:"id,omitempty"`
	Type       int       `gorm:"column:type" json:"type"`
	TypeID     int64     `gorm:"column:type_id" json:"typeId"`
	FromID     int64     `gorm:"column:from_id" json:"fromId"`
	ToID       int64     `gorm:"column:to_id" json:"toId"`
	Content    string    `gorm:"column:content" json:"content"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
}

// TableName ...
func (c Comment) TableName() string {
	return "comment"
}

func initSystemComment() {
	tx := dbOrmDefault.Model(&Comment{}).Begin()
	tx.Create(Comment{
		ID:         helper.GetMicrosecond(),
		Type:       0,
		TypeID:     0,
		FromID:     0,
		ToID:       0,
		Content:    "hello world",
		CreateTime: time.Now(),
	})
	tx.Commit()
}

// CreateComment ...
func CreateComment(Type int, typeID int64, fromID int64, toID int64, content string) (*Comment, error) {
	curTime := helper.GetMicrosecond()
	comment := &Comment{
		ID:         curTime,
		Type:       Type,
		TypeID:     typeID,
		FromID:     fromID,
		ToID:       toID,
		Content:    content,
		CreateTime: time.Now(),
	}

	err := dbOrmDefault.Model(&Comment{}).Create(comment).Error
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// QueryCommentByTypeID ...
func QueryCommentByTypeID(typeID int64) ([]*Comment, error) {
	list := make([]*Comment, 0)
	comment := &Comment{
		TypeID: typeID,
	}
	err := dbOrmDefault.Model(&Comment{}).Where(comment).Find(&list).Error
	if err == nil {
		return list, err
	}
	return nil, err
}
