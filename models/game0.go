package models

type Game0 struct {
	ID         		int64     	`gorm:"primary_key" json:"id,omitempty"`
	GName      		string    	`gorm:"column:g_name;type:varchar(100)" json:"name,omitempty"`
	GLevel	   		int 		`gorm:"column:g_level" json:"level"`

	GBlood	   		int 		`gorm:"column:g_blood" json:"blood"`
	GBloodAll	   	int 		`gorm:"column:g_blood_all" json:"blood_all"`
	GPower	   		int 		`gorm:"column:g_power" json:"power"`
}

type GameAction0 struct {
	Action 		string		`json:"action"`
	Target		int64		`json:"target"`
	Msg         string		`json:"msg"`
}