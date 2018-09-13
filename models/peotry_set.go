package models

import (
	"fmt"
	"io/ioutil"
	"github.com/tidwall/gjson"
)

type PeotrySet struct {
	ID    	int  	`gorm:"column:id;primary_key;auto_increment" json:"id"`

	UID   	int64  	`gorm:"column:u_id" json:"-"`
	UUser 	*User  	`gorm:"foreignkey:u_id" json:"user,omitempty"`

	SName 	string 	`gorm:"column(s_name);size(100)" json:"name"`
}

func initSystemPeotrySet() {
	setsJson, err := ioutil.ReadFile("data/sys-peotry-set.json")
	if err != nil {
		fmt.Println("read sys-peotry-set.json err")
		fmt.Println(err)
		return
	}

	re := gjson.ParseBytes(setsJson)
	re.ForEach(func (key, value gjson.Result) bool {
		sId := value.Get("s_id").Int()
		uId	:= value.Get("u_id").Int()
		sName := value.Get("s_name").String()
		savePeotrySet(int(sId), uId, sName)
		return true
	})
}

func savePeotrySet(id int, uId int64, name string) {
	peotrySet := PeotrySet {
		ID:		id,
		UID:	uId,
		SName:	name,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Save(peotrySet).Error
	if err != nil {
		fmt.Println(err)
	}
}

func CreatePeotrySet(uId int64, name string) error {
	set := PeotrySet {
		UID:	uId,
		SName: 	name,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Create(&set).Error
	return err
}

func QueryPeotrySetByID(id int) (*PeotrySet, error){
	set := &PeotrySet{
		ID: 	id,
	}

	err := dbOrmDefault.Model(&PeotrySet{}).Find(set).Error
	if err != nil {
		return nil, err
	}
	return set, nil
}

func QueryPeotrySetByUID(uId int64) ([]PeotrySet, error) {
	list := make([]PeotrySet, 0)
	set := &PeotrySet{
		UID:	uId,
	}
	err := dbOrmDefault.Model(&PeotrySet{}).Where(set).Find(&list).Error
	if err == nil {
		return list, err
	} else {
		return nil, err
	}
}

func DeletePeotrySet(id int) error{
	set := &PeotrySet{
		ID: 	id,
	}
	err := dbOrmDefault.Model(&PeotrySet{}).Delete(set).Error
	return err
}