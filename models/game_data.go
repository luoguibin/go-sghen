package models


type GameData struct {
	ID         		int64     		`gorm:"primary_key,id" json:"id,omitempty"`
	GName      		string    		`gorm:"column:g_name;type:varchar(100)" json:"name,omitempty"`
	GLevel	   		int 			`gorm:"column:g_level" json:"level"`

	GCondition		int   			`gorm:"column:g_condition" json:"condition"`

	GBloodBase		int 			`gorm:"column:g_blood_base" json:"bloodBase"`
	GBlood	   		int 			`gorm:"column:g_blood" json:"blood"`
	GBloodAll	   	int 			`gorm:"column:g_blood_all" json:"bloodAll"`

	GSpear	   		*GameSpear 		`gorm:"foreignkey:id" json:"spear"`
	GShield	   		*GameShield 	`gorm:"foreignkey:id" json:"shield"`

	GX				int				`gorm:"column:g_x" json:"x"`
	GY				int 			`gorm:"column:g_y" json:"y"`

	GOrders			[]*GameOrder		`gorm:"-" json:"orders"`
}

func initSystemGameData() {
	count := 0
	if dbOrmDefault.Model(&GameSpear{}).Count(&count); count == 0 {
		initSystemGameSpear()
	}
	if dbOrmDefault.Model(&GameShield{}).Count(&count); count == 0 {
		initSystemGameShield()
	}

	CreateGameData(15625045984, "Sghen", 103, 11100, 300000, 350000, 5000, 0, 0)
	CreateGameData(66666666, "Morge", 102, 10000, 350000, 500000, 4900, 0, 0)
	CreateGameData(88888888, "SghenMorge", 105, 11000, 320000, 400000, 5500, 0, 0)
}

func CreateGameData(id int64, gName string, gLevel int, gBloodBase int, gBlood int, gBloodAll int, gPower int, gX int, gY int) {
	gameData := GameData {
		ID:			id,
		GName:		gName,
		GLevel:		gLevel,
		GBloodBase:	gBloodBase,
		GBlood: 	gBlood,
		GBloodAll: 	gBloodAll,
		GX:			gX,
		GY:			gY,
	}
	err := dbOrmDefault.Model(&GameData{}).Save(gameData).Error
	if err != nil {
		MConfig.MLogger.Error(err.Error())
	}
}

func QueryGameData(id int64) (*GameData, error){
	gameData := &GameData {
		ID: 	id,
	}
	err := dbOrmDefault.Model(&GameData{}).Preload("GSpear").Preload("GShield").Find(gameData).Error
	if err != nil {
		return nil, err
	}
	return gameData, nil
}

func UpdateGameData(gameData *GameData) error{
	err := dbOrmDefault.Model(&GameData{}).Update(gameData).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}