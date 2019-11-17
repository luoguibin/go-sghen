package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// InitGorm ...
func InitGorm() {
	db, err := gorm.Open("mysql", MConfig.dBUsername+":"+MConfig.dBPassword+"@tcp("+MConfig.dBHost+")/"+MConfig.dBName+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		MConfig.MLogger.Error(err.Error())
		fmt.Println(err)
		return
	}

	db.DB().SetMaxIdleConns(MConfig.dBMaxIdle)
	db.DB().SetMaxOpenConns(MConfig.dBMaxConn)
	db.SingularTable(true) //禁用创建表名自动添加负数形式
	dbOrmDefault = db

	db.AutoMigrate(&User{}, &Peotry{}, &PeotrySet{}, PeotryImage{}, Comment{}, SmsCode{})

	count := 0
	if db.Model(&User{}).Count(&count); count == 0 {
		initSystemUser()
	}
	if db.Model(&Peotry{}).Count(&count); count == 0 {
		initSystemPeotry(true)
	}
	initSystemPeotry(false)
	if db.Model(&PeotrySet{}).Count(&count); count == 0 {
		initSystemPeotrySet()
	}

	db0, err0 := gorm.Open("mysql", MConfig.dBUsername+":"+MConfig.dBPassword+"@tcp("+MConfig.dBHost+")/"+MConfig.dBName0+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if err0 != nil {
		MConfig.MLogger.Error(err0.Error())
		fmt.Println(err0)
		return
	}

	db0.DB().SetMaxIdleConns(MConfig.dBMaxIdle)
	db0.DB().SetMaxOpenConns(MConfig.dBMaxConn)
	db0.SingularTable(true) //禁用创建表名自动添加负数形式
	db0.AutoMigrate()
	dbOrmDynamic = db0
}
