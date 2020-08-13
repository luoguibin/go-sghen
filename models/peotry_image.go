package models

import (
	"go-sghen/helper"
	"time"
)

// PeotryImage ...
type PeotryImage struct {
	ID         int64     `gorm:"column(id);" json:"id,omitempty"`
	Images     string    `gorm:"column(images);type:mediumtext" json:"images"`
	Count      int       `gorm:"column(count);" json:"count"`
	TimeCreate time.Time `gorm:"column:time_create" json:"time"`
}

// SavePeotryImage ...
func SavePeotryImage(id int64, images string, count int) error {
	peotryImage := &PeotryImage{
		ID:         id,
		Images:     images,
		Count:      count,
		TimeCreate: helper.StrToTimeStamp(helper.GetNowDateTime()),
	}

	err := dbOrmDefault.Model(&PeotryImage{}).Save(peotryImage).Error
	return err
}
