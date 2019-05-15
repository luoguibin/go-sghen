package models

import (
	"strings"
	"time"
)

type User struct {
	ID          int64     `gorm:"primary_key" json:"id,omitempty"`
	UPassword   string    `gorm:"column:u_password;type:varchar(20)" json:"-"`
	UName       string    `gorm:"column:u_name;type:varchar(100)" json:"name,omitempty"`
	UToken      string    `gorm:"-" json:"token,omitempty"`
	UIconURL    string    `gorm:"column:u_icon_url" json:"iconUrl"`
	UTimeCreate time.Time `gorm:"column:u_time_create" json:"timeCreate"`
	ULevel      int       `gorm:"column:u_level" json:"-"`
}

//  `json:"-"` 把struct编码成json字符串时，会忽略这个字段
//	`json:"id,omitempty"` //如果这个字段是空值，则不编码到JSON里面，否则用id为名字编码
//	`json:",omitempty"`   //如果这个字段是空值，则不编码到JSON里面，否则用属性名为名字编码

func (u User) TableName() string {
	return "user"
}

func initSystemUser() {
	tx := dbOrmDefault.Model(&User{}).Begin()
	tx.Create(User{
		ID:          15625045984,
		UPassword:   "123456",
		UName:       "Sghen",
		ULevel:      9,
		UTimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:          88888888,
		UPassword:   "123456",
		UName:       "Sghen",
		ULevel:      9,
		UTimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:          22222222,
		UPassword:   "123456",
		UName:       "Sghen",
		ULevel:      9,
		UTimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:          66666666,
		UPassword:   "123456",
		UName:       "Sghen",
		ULevel:      9,
		UTimeCreate: time.Now(),
	})
	tx.Commit()
}

func CreateUser(id int64, password string, name string) (*User, error) {
	user := &User{
		ID:          id,
		UPassword:   password,
		UName:       name,
		UTimeCreate: time.Now(),
		ULevel:      1,
	}

	err := dbOrmDefault.Model(&User{}).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func QueryUser(id int64) (*User, error) {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err == nil {
		return user, nil
	} else {
		return nil, err
	}
}

// QueryUsers ...
func QueryUsers(ids []int64) ([]*User, error) {
	list := make([]*User, 0)
	err := dbOrmDefault.Model(&User{}).Select("id, u_name, u_icon_url").Where("id in (?)", ids).Find(&list).Error
	if err == nil {
		return list, nil
	}
	return nil, err
}

func UpdateUser(id int64, password string, name string, iconURL string) (*User, error) {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err == nil {
		if len(strings.TrimSpace(password)) > 0 {
			user.UPassword = password
		}
		if len(strings.TrimSpace(name)) > 0 {
			user.UName = name
		}
		if len(strings.TrimSpace(iconURL)) > 0 {
			user.UIconURL = iconURL
		}

		err = dbOrmDefault.Model(&User{}).Save(user).Error
		if err != nil {
			return nil, err
		} else {
			return user, nil
		}
	} else {
		return nil, err
	}
}

func DeleteUser(id int64) error {
	user := &User{
		ID: id,
	}

	err := dbOrmDefault.Model(&User{}).Delete(user).Error
	return err
}
