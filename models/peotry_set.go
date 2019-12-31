package models

import (
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"
)

// PeotrySet 诗词选集
type PeotrySet struct {
	ID int `gorm:"column:id;primary_key;auto_increment" json:"id"`

	UserID int64 `gorm:"column:user_id" json:"userId"`
	User   *User `gorm:"foreignkey:user_id" json:"user,omitempty"`

	Name string `gorm:"column:name;size:100" json:"name"`
}

// initSystemPeotrySet ...
func initSystemPeotrySet() {
	setsJson, err := ioutil.ReadFile("data/sys-peotry-set.json")
	if err != nil {
		fmt.Println("read sys-peotry-set.json err")
		fmt.Println(err)
		return
	}

	re := gjson.ParseBytes(setsJson)
	re.ForEach(func(key, value gjson.Result) bool {
		id := value.Get("id").Int()
		userId := value.Get("userId").Int()
		sName := value.Get("name").String()
		savePeotrySet(int(id), userId, sName)
		return true
	})
}

// savePeotrySet ...
func savePeotrySet(id int, userId int64, name string) error {
	peotrySet := PeotrySet{
		ID:     id,
		UserID: userId,
		Name:   name,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Save(peotrySet).Error
	return err
}

// CreatePeotrySet ...
func CreatePeotrySet(userId int64, name string) error {
	set := PeotrySet{
		UserID: userId,
		Name:   name,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Create(&set).Error
	return err
}

// QueryPeotrySetByID ...
func QueryPeotrySetByID(id int) (*PeotrySet, error) {
	set := &PeotrySet{
		ID: id,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Find(set).Error
	return set, err
}

// QueryPeotrySetByUID ...
func QueryPeotrySetByUID(userId int64) ([]PeotrySet, error) {
	list := make([]PeotrySet, 0)
	err := dbOrmDefault.Model(&PeotrySet{}).Where("user_id = ? or user_id = 0", userId).Find(&list).Error
	return list, err
}

// DeletePeotrySet ...
func DeletePeotrySet(id int) error {
	set := &PeotrySet{
		ID: id,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Delete(set).Error
	return err
}
