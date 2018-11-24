package models

type GameSpear struct {
	ID			int64	`gorm:"primary_key,id" json:"id,omitempty"`

	SStrength	int 	`gorm:"column:s_strength" json:"strength"`
	SMana		int  	`gorm:"column:s_mana" json:"mana"`

	//  Five elements 
	SMetal		int 	`gorm:"column:s_metal" json:"metal"`
	SWood		int 	`gorm:"column:s_wood" json:"wood"`
	SWater		int 	`gorm:"column:s_water" json:"water"`
	SFire		int 	`gorm:"column:s_fire" json:"fire"`
	SEarth		int 	`gorm:"column:s_earth" json:"earth"` 
}

func initSystemGameSpear() {
	CreateGameSpear(15625045984, 900, 90, 1, 2, 5, 1, 2)
	CreateGameSpear(66666666, 880, 93, 1, 2, 3, 1, 2)
	CreateGameSpear(88888888, 910, 92, 2, 2, 2, 1, 3)
}

func CreateGameSpear(id int64, sStrength int, sMana int, sMetal int, sWood int, sWater int, sFire int, sEarth int) {
	gameSpear := GameSpear {
		ID:			id,
		SStrength: 	sStrength,
		SMana:		sMana,
		SMetal:		sMana,
		SWood: 		sWood,
		SWater:		sWater,
		SFire:		sFire,
		SEarth: 	sEarth,
	}
	err := dbOrmDefault.Model(&GameSpear{}).Save(gameSpear).Error
	if err != nil {
		MConfig.MLogger.Error(err.Error())
	}
}

func QueryGameSpear(id int64) (*GameSpear, error){
	gameSpear := &GameSpear {
		ID: 	id,
	}
	err := dbOrmDefault.Model(&GameSpear{}).Find(gameSpear).Error
	if err != nil {
		return nil, err
	}
	return gameSpear, nil
}

func UpdateGameSpear(gameSpear *GameSpear) error{
	err := dbOrmDefault.Model(&GameSpear{}).Update(gameSpear).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}