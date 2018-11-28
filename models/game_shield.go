package models

type GameShield struct {
	ID			int64	`gorm:"primary_key,id" json:"id,omitempty"`

	SStrength	int 	`gorm:"column:s_strength" json:"antiStrength"`
	SMana		int  	`gorm:"column:s_mana" json:"antiMana"`

	//  Five elements 
	SMetal		int 	`gorm:"column:s_metal" json:"antiMetal"`
	SWood		int 	`gorm:"column:s_wood" json:"antiWood"`
	SWater		int 	`gorm:"column:s_water" json:"antiWater"`
	SFire		int 	`gorm:"column:s_fire" json:"antiFire"`
	SEarth		int 	`gorm:"column:s_earth" json:"antiEarth"` 
}

func initSystemGameShield() {
	CreateGameShield(15625045984, 800, 10, 0, 0, 2, 0, 0)
	CreateGameShield(66666666, 780, 13, 0, 0, 1, 0, 0)
	CreateGameShield(88888888, 810, 12, 0, 0, 0, 0, 1)
}

func CreateGameShield(id int64, sStrength int, sMana int, sMetal int, sWood int, sWater int, sFire int, sEarth int) {
	gameShield := GameShield {
		ID:			id,
		SStrength: 	sStrength,
		SMana:		sMana,
		SMetal:		sMetal,
		SWood: 		sWood,
		SWater:		sWater,
		SFire:		sFire,
		SEarth: 	sEarth,
	}
	err := dbOrmDefault.Model(&GameShield{}).Save(gameShield).Error
	if err != nil {
		MConfig.MLogger.Error(err.Error())
	}
}

func QueryGameShield(id int64) (*GameShield, error){
	gameShield := &GameShield {
		ID: 	id,
	}
	err := dbOrmDefault.Model(&GameShield{}).Find(gameShield).Error
	if err != nil {
		return nil, err
	}
	return gameShield, nil
}

func UpdateGameShield(gameShield *GameShield) error{
	err := dbOrmDefault.Model(&GameShield{}).Update(gameShield).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}