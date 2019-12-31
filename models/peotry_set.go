package models

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/tidwall/gjson"
)

// PeotrySet 诗词选集
type PeotrySet struct {
	ID int `gorm:"column:id;primary_key;auto_increment" json:"id"`

	UserID int64 `gorm:"column:user_id" json:"userId"`
	User   *User `gorm:"foreignkey:user_id" json:"user,omitempty"`

	Name       string    `gorm:"column:name;size:100" json:"name"`
	TimeCreate time.Time `gorm:"column:time_create" json:"timeCreate"`
}

// initSystemPeotrySet 初始化系统选集数据
func initSystemPeotrySet() {
	dataJSON, err := ioutil.ReadFile("data/sys-peotry-set.json")
	if err != nil {
		fmt.Println("read sys-peotry-set.json err")
		fmt.Println(err)
		return
	}

	re := gjson.ParseBytes(dataJSON)
	re.ForEach(func(key, value gjson.Result) bool {
		id := value.Get("id").Int()
		userID := value.Get("userId").Int()
		sName := value.Get("name").String()
		savePeotrySet(int(id), userID, sName)
		return true
	})
}

// savePeotrySet 本地保存选集
func savePeotrySet(id int, userID int64, name string) error {
	peotrySet := PeotrySet{
		ID:         id,
		UserID:     userID,
		Name:       name,
		TimeCreate: time.Now(),
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Save(peotrySet).Error
	return err
}

// CreatePeotrySet 创建选集，id自增
func CreatePeotrySet(userID int64, name string) error {
	set := PeotrySet{
		UserID:     userID,
		Name:       name,
		TimeCreate: time.Now(),
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Create(&set).Error
	return err
}

// QueryPeotrySetByID 查询某个选集
func QueryPeotrySetByID(id int) (*PeotrySet, error) {
	set := &PeotrySet{
		ID: id,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Find(set).Error
	return set, err
}

// QueryPeotrySetByUID 查询某个用户的选集和系统默认选集
func QueryPeotrySetByUID(userID int64) ([]PeotrySet, error) {
	list := make([]PeotrySet, 0)
	err := dbOrmDefault.Model(&PeotrySet{}).Where("user_id = ? or user_id = 0", userID).Find(&list).Error
	return list, err
}

// DeletePeotrySet 删除选集
func DeletePeotrySet(id int) error {
	set := &PeotrySet{
		ID: id,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Delete(set).Error
	return err
}
