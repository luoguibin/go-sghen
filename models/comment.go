package models

import (
	"go-sghen/helper"
	"time"
)

// Comment ...
type Comment struct {
	ID         int64     `gorm:"primary_key,id" json:"id,omitempty"`
	Type       int       `gorm:"column:type" json:"type"` // 1为诗歌
	TypeID     int64     `gorm:"index;column:type_id" json:"typeId"`
	TypeUserID int64     `gorm:"index;column:type_user_id" json:"typeUserId"`
	FromID     int64     `gorm:"index;column:from_id" json:"fromId"`
	ToID       int64     `gorm:"index;column:to_id" json:"toId"`
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
		ID:         helper.NewUinqueID(),
		Type:       0,
		TypeID:     0,
		TypeUserID: 0,
		FromID:     0,
		ToID:       0,
		Content:    "hello world",
		CreateTime: time.Now(),
	})
	tx.Commit()
}

// CreateComment ...
func CreateComment(Type int, typeID int64, typeUserID int64, fromID int64, toID int64, content string) (*Comment, error) {
	comment := &Comment{
		ID:         helper.NewUinqueID(),
		Type:       Type,
		TypeID:     typeID,
		TypeUserID: typeUserID,
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

// QueryComment ...
func QueryComment(ID int64) (*Comment, error) {
	comment := &Comment{
		ID: ID,
	}
	err := dbOrmDefault.Model(&Comment{}).Find(comment).Error
	if err == nil {
		return comment, nil
	}
	return nil, err
}

// QueryCommentByTypeIDFromID ...
func QueryCommentByTypeIDFromID(typeID int64, fromID int64, toID int64) (*Comment, error) {
	comment := &Comment{
		TypeID: typeID,
		FromID: fromID,
		ToID:   toID,
	}
	err := dbOrmDefault.Model(&Comment{}).Where(comment).Find(comment).Error
	if err == nil {
		return comment, nil
	}
	return nil, err
}

// QueryCommentByTypeID ...
func QueryCommentByTypeID(typeID int64) ([]*Comment, error) {
	list := make([]*Comment, 0)
	comment := &Comment{
		TypeID: typeID,
	}
	err := dbOrmDefault.Model(&Comment{}).Where(comment).Find(&list).Error
	if err == nil {
		return list, nil
	}
	return nil, err
}

// SaveComment ...
func SaveComment(comment *Comment) error {
	comment.CreateTime = time.Now()
	err := dbOrmDefault.Model(&Comment{}).Save(comment).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteComment ...
func DeleteComment(ID int64) error {
	comment := &Comment{
		ID: ID,
	}
	err := dbOrmDefault.Model(&Comment{}).Delete(comment).Error
	return err
}

// DeleteComments ...
func DeleteComments(typeID int64) error {
	err := dbOrmDefault.Delete(&Comment{}, "type_id = ?", typeID).Error
	return err
}
