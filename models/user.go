package models

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"go-sghen/helper"

	"github.com/tidwall/gjson"
)

// User ...
type User struct {
	ID         int64     `gorm:"primary_key" json:"id,omitempty"`
	Password   string    `gorm:"column:password;type:varchar(300)" json:"-"`
	Name       string    `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Token      string    `gorm:"-" json:"token,omitempty"`
	IconURL    string    `gorm:"column:icon_url" json:"iconUrl"`
	TimeCreate time.Time `gorm:"column:time_create" json:"timeCreate"`
	Level      int       `gorm:"column:level" json:"-"`
}

// TableName ...
func (u User) TableName() string {
	return "user"
}

func initSystemUser() {
	usersJson, err := ioutil.ReadFile("data/sys-account.json")
	if err != nil {
		fmt.Println("read sys-account.json err")
		fmt.Println(err)
		return
	}

	re := gjson.ParseBytes(usersJson)
	re.ForEach(func(key, value gjson.Result) bool {
		id := value.Get("id").Int()
		password := value.Get("password").String()
		name := value.Get("name").String()
		level := value.Get("level").Int()

		CreateUser(id, helper.MD5(password), name, int(level))
		return true
	})
}

// CreateUser ...
func CreateUser(id int64, password string, name string, level int) (*User, error) {
	user := &User{
		ID:         id,
		Password:   password,
		Name:       name,
		TimeCreate: time.Now(),
		Level:      level,
	}

	err := dbOrmDefault.Model(&User{}).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// QueryUser ...
func QueryUser(id int64) (*User, error) {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err == nil {
		return user, nil
	}
	return nil, err
}

// QueryUsers ...
func QueryUsers(ids []int64) ([]*User, error) {
	list := make([]*User, 0)
	err := dbOrmDefault.Model(&User{}).Select("id, name, icon_url").Where("id in (?)", ids).Find(&list).Error
	if err == nil {
		return list, nil
	}
	return nil, err
}

// UpdateUser ...
func UpdateUser(id int64, password string, name string, iconUrl string) (*User, error) {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err == nil {
		if len(strings.TrimSpace(password)) > 0 {
			user.Password = password
		}
		if len(strings.TrimSpace(name)) > 0 {
			user.Name = name
		}
		if len(strings.TrimSpace(iconUrl)) > 0 {
			user.IconURL = iconUrl
		}

		err = dbOrmDefault.Model(&User{}).Save(user).Error
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, err
}

// DeleteUser ...
func DeleteUser(id int64) error {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Delete(user).Error
	return err
}
