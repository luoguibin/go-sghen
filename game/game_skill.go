package game

import (
	"SghenApi/models"
	"SghenApi/helper"
	"math/rand"
)

const (
	// skill count which is under 1000, because of its id is base on `OT_Skill`
	OT_SkillSingle0     =   OT_SkillSingle + 1
	OT_SkillSingle1     =   OT_SkillSingle + 2
	OT_SkillSingle2     =   OT_SkillSingle + 3
	OT_SkillSingle3     =   OT_SkillSingle + 4
	OT_SkillSingle4     =   OT_SkillSingle + 5

	OT_SkillSingleK0    =   OT_SkillSingleK + 1
	   
	OT_SkillNear0       =   OT_SkillNear + 1
	   
    OT_SkillNearK0      =   OT_SkillNearK + 1   
)


func getSkillSingleDamage(id int, data0 *models.GameData, data1 *models.GameData) int {
	d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
	if d > 50 {
		return int(-d)
	}

	ran := rand.Intn(100)
	if rand.Intn(10) < 5 {
		ran = data0.GSpear.SStrength + ran
	} else {
		ran = data0.GSpear.SStrength - ran
	}
	return ran
}
