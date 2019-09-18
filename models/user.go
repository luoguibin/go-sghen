package models

import (
	"strings"
	"time"
)

// User ...
type User struct {
	ID         int64     `gorm:"primary_key" json:"id,omitempty"`
	Password   string    `gorm:"column:password;type:varchar(20)" json:"-"`
	Name       string    `gorm:"column:name;type:varchar(100)" json:"name,omitempty"`
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
	tx := dbOrmDefault.Model(&User{}).Begin()
	tx.Create(User{
		ID:         15625045984,
		Password:   "123456",
		Name:       "乂末",
		Level:      9,
		TimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:         15688888888,
		Password:   "123456",
		Name:       "Sghen",
		Level:      9,
		TimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:         15622222222,
		Password:   "123456",
		Name:       "Morge",
		Level:      9,
		TimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:         15666666666,
		Password:   "123456",
		Name:       "SghenMorge",
		Level:      9,
		TimeCreate: time.Now(),
	})
	tx.Commit()
}

// CreateUser ...
func CreateUser(ID int64, Password string, Name string) (*User, error) {
	user := &User{
		ID:         ID,
		Password:   Password,
		Name:       Name,
		TimeCreate: time.Now(),
		Level:      1,
	}

	err := dbOrmDefault.Model(&User{}).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// QueryUser ...
func QueryUser(ID int64) (*User, error) {
	user := &User{
		ID: ID,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err == nil {
		return user, nil
	} else {
		return nil, err
	}
}

// QueryUsers ...
func QueryUsers(IDs []int64) ([]*User, error) {
	list := make([]*User, 0)
	err := dbOrmDefault.Model(&User{}).Select("id, u_name, u_icon_url").Where("id in (?)", IDs).Find(&list).Error
	if err == nil {
		return list, nil
	}
	return nil, err
}

// UpdateUser ...
func UpdateUser(ID int64, Password string, Name string, IconURL string) (*User, error) {
	user := &User{
		ID: ID,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err == nil {
		if len(strings.TrimSpace(Password)) > 0 {
			user.Password = Password
		}
		if len(strings.TrimSpace(Name)) > 0 {
			user.Name = Name
		}
		if len(strings.TrimSpace(IconURL)) > 0 {
			user.IconURL = IconURL
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

// DeleteUser ...
func DeleteUser(ID int64) error {
	user := &User{
		ID: ID,
	}

	err := dbOrmDefault.Model(&User{}).Delete(user).Error
	return err
}
