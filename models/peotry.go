package models

import (
	"fmt"
	"go-sghen/helper"
	"io/ioutil"
	"time"

	"github.com/tidwall/gjson"
)

// Peotry ...
type Peotry struct {
	ID int64 `gorm:"column:id;primary_key;" json:"id"`

	UID   int64 `gorm:"column:u_id" json:"-"`
	UUser *User `gorm:"foreignkey:u_id;" json:"user"`

	SID  int        `gorm:"column:s_id" json:"-"`
	SSet *PeotrySet `gorm:"foreignkey:s_id" json:"set"`

	PTitle   string    `gorm:"column:p_title;type:varchar(20)" json:"title"`
	PTime    time.Time `gorm:"column:p_time" json:"time"`
	PContent string    `gorm:"column:p_content;type:mediumtext" json:"content"`
	PEnd     string    `gorm:"column:p_end" json:"end"`

	PImage *PeotryImage `gorm:"foreignkey:id" json:"image,omitempty"`

	Comments []*Comment `gorm:"-" json:"comments,omitempty"`
}

func initSystemPeotry() {
	peotriesJson, err := ioutil.ReadFile("data/sys-peotry.json")
	if err != nil {
		fmt.Println("read sys-peotry.json err")
		fmt.Println(err)
		return
	}

	re := gjson.ParseBytes(peotriesJson)
	re.ForEach(func(key, value gjson.Result) bool {
		uId := value.Get("u_id").Int()
		sId := value.Get("s_id").Int()
		pTitle := value.Get("p_title").String()
		pTime := value.Get("p_time").String()
		pContent := value.Get("p_content").String()
		pEnd := value.Get("p_end").String()
		pImages := value.Get("p_images").String()
		CreatePeotry(uId, int(sId), pTitle, pTime, pContent, pEnd, pImages)
		return true
	})
}

func CreatePeotry(userId int64, setId int, title string, pTime string, content string, end string, images string) (int64, error) {
	curTime := helper.GetMicrosecond()
	peotry := Peotry{
		ID:       curTime,
		UID:      userId,
		SID:      setId,
		PTitle:   title,
		PTime:    helper.StrToTimeStamp(pTime),
		PContent: content,
		PEnd:     end,
	}

	err := dbOrmDefault.Model(&Peotry{}).Save(peotry).Error
	if err != nil {
		fmt.Println(err)
		return 0, err
	} else {
		res := gjson.Parse(images)
		imgs := res.Array()
		l := len(imgs)
		if l > 0 {
			SavePeotryImage(curTime, images, l)
		}
	}
	return curTime, nil
}

func UpdatePeotry(peotry *Peotry) error {
	err := dbOrmDefault.Model(&Peotry{}).Update(peotry).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func QueryPeotry(setId int, page int, limit int, content string) ([]*Peotry, error, int, int, int, int) {
	list := make([]*Peotry, 0)
	totalPage := 0
	count := 0
	curPage := page
	pageIsEnd := 0

	if limit == 0 {
		limit = 10
	}

	db := dbOrmDefault.Model(&Peotry{})
	if setId > 0 {
		query := &Peotry{
			SID: setId,
		}
		db = db.Where(query)
	}

	if len(content) > 1 {
		db = db.Where("p_content LIKE ?", "%"+content+"%")
	}

	db.Count(&count)
	db = db.Preload("UUser").Preload("SSet").Preload("PImage")
	err := db.Limit(limit).Offset(helper.PageOffset(limit, page)).Find(&list).Error

	if err == nil {
		totalPage, pageIsEnd = helper.PageTotal(limit, page, int64(count))
	} else {
		return nil, err, 0, 0, 0, 0
	}

	return list, nil, count, totalPage, curPage, pageIsEnd
}

func QueryPeotryByID(id int64) (*Peotry, error) {
	peotry := &Peotry{
		ID: id,
	}

	err := dbOrmDefault.Model(&Peotry{}).Preload("UUser").Preload("SSet").Preload("PImage").Find(peotry).Error
	if err == nil {
		return peotry, nil
	} else {
		return nil, err
	}
}

func DeletePeotry(id int64) error {
	set := &Peotry{
		ID: id,
	}

	err := dbOrmDefault.Model(&Peotry{}).Delete(&set).Error
	return err
}
