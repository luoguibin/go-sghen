package models

import (
	"time"
)

type User struct {
	ID         		int64     	`gorm:"primary_key" json:"id"`
	UPassword  		string    	`gorm:"column:u_password;type:varchar(20)" json:"-"`
	UName      		string    	`gorm:"column:u_name;type:varchar(100)" json:"name"`
	UToken     		string    	`gorm:"-" json:"token"`
	UTimeCreate 	time.Time 	`gorm:"column:u_time_create" json:"timeCreate"`
	ULevel	   		int 		`gorm:"column:u_level" json:"-"`
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
			ID: 			15625045984, 
			UPassword: 		"123456",
			UName: 			"Sghen",
			ULevel: 		9,
			UTimeCreate:	time.Now(),
		})
	tx.Commit()
}


func CreateUser(id int64, password string, name string) (*User, error){
	user := &User{
		ID:				id,
		UPassword:		password,
		UName:			name,
		UTimeCreate:	time.Now(),
		ULevel:			1,
	}

	err := dbOrmDefault.Model(&User{}).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func QueryUser() {
	
}

func UpdateUser() {

}

func DeleteUser() {

}