package models


type GameData struct {
	ID         		int64     	`gorm:"primary_key" json:"id,omitempty"`
	GName      		string    	`gorm:"column:g_name;type:varchar(100)" json:"name,omitempty"`
	GLevel	   		int 		`gorm:"column:g_level" json:"level"`

	GBlood	   		int 		`gorm:"column:g_blood" json:"blood"`
	GBloodAll	   	int 		`gorm:"column:g_blood_all" json:"blood_all"`
	GPower	   		int 		`gorm:"column:g_power" json:"power"`
}

func initSystemGameData() {
	CreateGameData(15625045984, "Sghen", 103, 300000, 350000, 5000)
	CreateGameData(66666666, "Morge", 102, 350000, 500000, 4900)
	CreateGameData(88888888, "SghenMorge", 105, 320000, 400000, 5500)
}

func CreateGameData(id int64, gName string, gLevel int, gBlood int, gBloodAll int, gPower int) {
	gameData := GameData {
		ID:			id,
		GName:		gName,
		GLevel:		gLevel,
		GBlood: 	gBlood,
		GBloodAll: 	gBloodAll,
		GPower:		gPower,
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
	err := dbOrmDefault.Model(&GameData{}).Find(gameData).Error
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