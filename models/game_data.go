package models


type GameData struct {
	ID         		int64     		`gorm:"primary_key,id" json:"id,omitempty"`
	Name      		string    		`gorm:"column:g_name;type:varchar(100)" json:"name,omitempty"`
	Level	   		int 			`gorm:"column:g_level" json:"level"`

	Condition		int   			`gorm:"column:g_condition" json:"condition"`

	BloodBase		int 			`gorm:"column:g_blood_base" json:"bloodBase"`
	Blood	   		int 			`gorm:"column:g_blood" json:"blood"`
	BloodAll	   	int 			`gorm:"-" json:"bloodAll"`

	Spear	   		*GameSpear 		`gorm:"foreignkey:id" json:"spear"`
	Shield	   		*GameShield 	`gorm:"foreignkey:id" json:"shield"`

	MapId			int				`gorm:"column:g_map_id" json:"mapId"`
	ScreenId		int 			`gorm:"-" json:"screenId"`
	SpeedBase		int 			`gorm:"column:g_speed_base" json:"speedBase"`
	Speed			int 			`gorm:"-" json:"speed"`
	X				int				`gorm:"column:g_x" json:"x"`
	Y				int 			`gorm:"column:g_y" json:"y"`

	Move			int 			`gorm:"-" json:"-"`
	X0				int				`gorm:"-" json:"-"`
	Y0				int				`gorm:"-" json:"-"`
	X1				int				`gorm:"-" json:"-"`
	Y1				int				`gorm:"-" json:"-"`
	MoveTime		int64			`gorm:"-" json:"-"`
	EndTime		int64			`gorm:"-" json:"-"`

	Orders			[]*interface{}		`gorm:"-" json:"orders"`
}

func initSystemGameData() {
	count := 0
	if dbOrmDefault.Model(&GameSpear{}).Count(&count); count == 0 {
		initSystemGameSpear()
	}
	if dbOrmDefault.Model(&GameShield{}).Count(&count); count == 0 {
		initSystemGameShield()
	}

	CreateGameData(15625045984, "Sghen", 103, 111000, 11100, 0, 50, 0, 0)
	CreateGameData(66666666, "Morge", 102, 100000, 10000, 0, 50, 0, 0)
	CreateGameData(88888888, "SghenMorge", 105, 110000, 11000, 0, 50, 0, 0)
}

func CreateGameData(id int64, name string, level int, bloodBase int, blood int, mapId int,speedBase int, x int, y int) {
	gameData := GameData {
		ID:			id,
		Name:		name,
		Level:		level,
		BloodBase:	bloodBase,
		Blood:		blood,
		SpeedBase:	speedBase,
		X:			x,
		Y:			y,
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
	err := dbOrmDefault.Model(&GameData{}).Preload("Spear").Preload("Shield").Find(gameData).Error
	if err != nil {
		return nil, err
	}
	return gameData, nil
}

func UpdateGameData(gameData *GameData) error{
	// update: nothing will be updated such as "", 0, false are blank values of their types
	err := dbOrmDefault.Model(&GameData{}).Save(gameData).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}