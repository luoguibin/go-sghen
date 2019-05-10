package models

import (
	"github.com/jinzhu/gorm"
)

func InitGorm() {
	db, er := gorm.Open("mysql", MConfig.dBUsername+":"+MConfig.dBPassword+"@tcp("+MConfig.dBHost+")/"+MConfig.dBName+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if er != nil {
		MConfig.MLogger.Error(er.Error())
		return
	}

	db.DB().SetMaxIdleConns(MConfig.dBMaxIdle)
	db.DB().SetMaxOpenConns(MConfig.dBMaxConn)

	db.SingularTable(true) //禁用创建表名自动添加负数形式

	dbOrmDefault = db

	db.AutoMigrate(&User{}, &PeotrySet{}, &PeotryImage{}, &Peotry{}, &Comment{}, &GameData{}, &GameSpear{}, &GameShield{})

	count := 0
	if db.Model(&User{}).Count(&count); count == 0 {
		initSystemUser()
	}
	if db.Model(&PeotrySet{}).Count(&count); count == 0 {
		initSystemPeotrySet()
	}
	if db.Model(&Peotry{}).Count(&count); count == 0 {
		initSystemPeotry()
	}
	if db.Model(&Comment{}).Count(&count); count == 0 {
		initSystemComment()
	}
	if db.Model(&GameData{}).Count(&count); count == 0 {
		initSystemGameData()
	}
}
