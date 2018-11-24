package models


type GameData struct {
	ID         		int64     		`gorm:"primary_key,id" json:"id,omitempty"`
	GName      		string    		`gorm:"column:g_name;type:varchar(100)" json:"name,omitempty"`
	GLevel	   		int 			`gorm:"column:g_level" json:"level"`

	GCondition		int   			`gorm:"column:g_condition" json:"condition"`

	GBloodBase		int 			`gorm:"column:g_blood_base" json:"bloodBase"`
	GBlood	   		int 			`gorm:"column:g_blood" json:"blood"`
	GBloodAll	   	int 			`gorm:"-" json:"bloodAll"`

	GSpear	   		*GameSpear 		`gorm:"foreignkey:id" json:"spear"`
	GShield	   		*GameShield 	`gorm:"foreignkey:id" json:"shield"`

	GMapId			int				`gorm:"column:g_map_id" json:"mapId"`
	GSpeedBase		int 			`gorm:"column:g_speed_base" json:"speedBase"`
	GSpeed			int 			`gorm:"-" json:"speed"`
	GX				int				`gorm:"column:g_x" json:"x"`
	GY				int 			`gorm:"column:g_y" json:"y"`

	GOrders			[]*interface{}		`gorm:"-" json:"orders"`
}

func initSystemGameData() {
	count := 0
	if dbOrmDefault.Model(&GameSpear{}).Count(&count); count == 0 {
		initSystemGameSpear()
	}
	if dbOrmDefault.Model(&GameShield{}).Count(&count); count == 0 {
		initSystemGameShield()
	}

	CreateGameData(15625045984, "Sghen", 103, 11100, 11100, 0, 50, 0, 0)
	CreateGameData(66666666, "Morge", 102, 10000, 10000, 0, 50, 0, 0)
	CreateGameData(88888888, "SghenMorge", 105, 11000, 11000, 0, 50, 0, 0)
}

func CreateGameData(id int64, gName string, gLevel int, gBloodBase int, gBlood int, gMapId int, gSpeedBase int, gX int, gY int) {
	gameData := GameData {
		ID:			id,
		GName:		gName,
		GLevel:		gLevel,
		GBloodBase:	gBloodBase,
		GBlood:		gBlood,
		GSpeedBase:	gSpeedBase,
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