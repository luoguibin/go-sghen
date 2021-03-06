package models

import (
	"fmt"
	"go-sghen/helper"
	"io/ioutil"
	"sort"
	"time"

	"github.com/tidwall/gjson"
)

// Peotry ...
type Peotry struct {
	ID int64 `gorm:"column:id;primary_key;" json:"id"`

	UserID int64 `gorm:"column:user_id" json:"-"`
	User   *User `gorm:"foreignkey:user_id;" json:"user"`

	SetID int        `gorm:"column:set_id" json:"-"`
	Set   *PeotrySet `gorm:"foreignkey:set_id" json:"set"`

	Title      string    `gorm:"column:title;type:varchar(20)" json:"title"`
	TimeCreate time.Time `gorm:"column:time_create" json:"time"`
	Content    string    `gorm:"column:content;type:mediumtext" json:"content"`
	End        string    `gorm:"column:end" json:"end"`

	Image *PeotryImage `gorm:"foreignkey:id" json:"image,omitempty"`

	Comments []*Comment `gorm:"-" json:"comments,omitempty"`
}

// initSystemPeotry 初始化诗词
func initSystemPeotry(isDefault bool) {
	var name string
	if isDefault {
		name = "data/sys-peotry.json"
	} else {
		name = "data/sys-peotry_temp.json"
	}
	peotriesJson, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("read sys-peotry.json err")
		fmt.Println(err)
		return
	}

	re := gjson.ParseBytes(peotriesJson)
	re.ForEach(func(key, value gjson.Result) bool {
		userId := value.Get("userId").Int()
		setId := value.Get("setId").Int()
		title := value.Get("title").String()
		time := value.Get("timeCreate").String()
		content := value.Get("content").String()
		end := value.Get("end").String()
		image := value.Get("image").String()
		// todo: check the image json string if valid
		CreatePeotry(userId, int(setId), title, time, content, end, image)
		return true
	})
}

func AddTempPeotry() {
	initSystemPeotry(false)
}

// CreatePeotry 创建诗词
func CreatePeotry(userId int64, setId int, title string, time string, content string, end string, images string) (int64, error) {
	id := helper.NewUinqueID()
	peotry := Peotry{
		ID:         id,
		UserID:     userId,
		SetID:      setId,
		Title:      title,
		TimeCreate: helper.StrToTimeStamp(time),
		Content:    content,
		End:        end,
	}

	err := dbOrmDefault.Model(&Peotry{}).Save(peotry).Error
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	res := gjson.Parse(images)
	imgs := res.Array()
	l := len(imgs)
	if l > 0 {
		SavePeotryImage(id, images, l)
	}
	return id, nil
}

// UpdatePeotry ...
func UpdatePeotry(peotry *Peotry) error {
	err := dbOrmDefault.Model(&Peotry{}).Update(peotry).Error
	return err
}

// QueryPeotry ...
func QueryPeotry(userID int64, setId int, page int, limit int, content string) ([]*Peotry, int, int, int, int, error) {
	list := make([]*Peotry, 0)
	totalPage := 0
	count := 0
	curPage := page
	pageIsEnd := 0

	if limit == 0 {
		limit = 10
	}

	db := dbOrmDefault.Model(&Peotry{})
	if userID > 0 {
		query := &Peotry{
			UserID: userID,
		}
		db = db.Where(query)
	}
	if setId > 0 {
		query := &Peotry{
			SetID: setId,
		}
		db = db.Where(query)
	}

	if len(content) > 1 {
		// todo
		db = db.Where("content LIKE ?", "%"+content+"%")
	}
	db = db.Order("time_create desc")

	db.Count(&count)
	db = db.Preload("User").Preload("Set").Preload("Image")
	err := db.Limit(limit).Offset(helper.PageOffset(limit, page)).Find(&list).Error

	if err == nil {
		totalPage, pageIsEnd = helper.PageTotal(limit, page, int64(count))
		return list, count, totalPage, curPage, pageIsEnd, nil
	}
	return nil, 0, 0, 0, 0, err
}

// QueryPopularPeotry ...
func QueryPopularPeotry(limit int) ([]*Peotry, error) {
	comments := make([]*Comment, 0)

	db := dbOrmDefault.Model(&Comment{})
	db = db.Select("type_id, count(*) as repeat_count")
	db = db.Where("to_id=? AND content=?", -1, "praise")
	db = db.Group("type_id").Having("repeat_count > 1").Order("repeat_count DESC")
	err := db.Limit(limit).Find(&comments).Error

	if err != nil {
		return nil, err
	}

	var ids []int64
	idMap := make(map[int64]int)
	for i, comment := range comments {
		ids = append(ids, comment.TypeID)
		idMap[comment.TypeID] = i
	}

	peotrys := make([]*Peotry, 0)
	db = dbOrmDefault.Model(&Peotry{})
	db = db.Preload("User").Preload("Set").Preload("Image")
	err = db.Where("id in (?)", ids).Find(&peotrys).Error

	if err == nil {
		sort.Slice(peotrys, func(i, j int) bool {
			return idMap[peotrys[i].ID] < idMap[peotrys[j].ID]
		})
	}

	return peotrys, err
}

// QueryPeotryByID ...
func QueryPeotryByID(id int64) (*Peotry, error) {
	peotry := &Peotry{
		ID: id,
	}

	err := dbOrmDefault.Model(&Peotry{}).Preload("User").Preload("Set").Preload("Image").Find(peotry).Error
	return peotry, err
}

// DeletePeotry ...
func DeletePeotry(id int64) error {
	peotry := &Peotry{
		ID: id,
	}

	err := dbOrmDefault.Model(&Peotry{}).Delete(peotry).Error
	return err
}
