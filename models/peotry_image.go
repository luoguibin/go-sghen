package models

// PeotryImage ...
type PeotryImage struct {
	ID     int64  `gorm:"column(id);" json:"id,omitempty"`
	Images string `gorm:"column(images);type:mediumtext" json:"images"`
	Count  int    `gorm:"column(count);" json:"count"`
}

// SavePeotryImage ...
func SavePeotryImage(id int64, images string, count int) error {
	peotryImage := &PeotryImage{
		ID:     id,
		Images: images,
		Count:  count,
	}

	err := dbOrmDefault.Model(&PeotryImage{}).Save(peotryImage).Error
	return err
}
