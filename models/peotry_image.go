package models

type PeotryImage struct {
	ID    		int64  	`orm:"column(id);pk" json:"id,omitempty"`
	IImages   	string  `orm:"column(i_images);type:mediumtext" json:"images"`
	ICount 		int 	`orm:"column(i_count);" json:"count"`
}

func SavePeotryImage(id int64, images string, count int) error {
	peotryImage := &PeotryImage {
		ID:			id,
		IImages:	images,
		ICount:		count,
	}

	err := dbOrmDefault.Model(&PeotryImage{}).Save(peotryImage).Error
	return err
}
