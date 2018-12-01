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


func getSkillSingleDamage(id int, data0 *models.GameData, data1 *models.GameData, limit int) int {
	distance := helper.GClientDistance(data0.X, data0.Y, data1.X, data1.Y)
	if int(distance) > limit {
		return -limit
	}
	spear0 := getSpearValue(data0.Spear)
	switch id {
	case OT_SkillSingle0:
		spear0 += int(0.5 * float32(data0.Spear.SMana)) * 10
	case OT_SkillSingle1:
		spear0 += int(1.3 * float32(getSpearFiveEleValue(data0.Spear))) * 100
	}
	shield1 := getShieldValue(data1.Shield)
	damage := spear0 - shield1

	ran := rand.Intn(data0.Spear.SStrength / 10)
	if rand.Intn(10) < 5 {
		damage += ran
	} else {
		damage -= ran
	}
	if (damage < 0) {
		damage = 0
	}
	return damage
}


func getSpearValue(spear *models.GameSpear) int {
	power := spear.SStrength;
	power += spear.SMana * 10;

	fiveVal := getSpearFiveEleValue(spear)
	power += fiveVal * 100
	return power
}

func getShieldValue(shield *models.GameShield) int {
	power := shield.SStrength;
	power += shield.SMana * 10;

	fiveVal := getShieldFiveEleValue(shield)
	power += fiveVal * 100
	return power
}

func getSpearFiveEleValue(spear *models.GameSpear) int {
	return spear.SMetal + spear.SWood + spear.SWater + spear.SFire + spear.SEarth
}

func getShieldFiveEleValue(shield *models.GameShield) int {
	return shield.SMetal + shield.SWood + shield.SWater + shield.SFire + shield.SEarth
}