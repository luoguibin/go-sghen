package game

import (
	"go-sghen/helper"
	"go-sghen/models"
	"math/rand"
)

const (
	// skill count which is under 1000, because of its id is base on `OT_Skill`
	OT_SkillSingle0 = OT_SkillSingle + 1
	OT_SkillSingle1 = OT_SkillSingle + 2
	OT_SkillSingle2 = OT_SkillSingle + 3
	OT_SkillSingle3 = OT_SkillSingle + 4
	OT_SkillSingle4 = OT_SkillSingle + 5

	OT_SkillSingleK0 = OT_SkillSingleK + 1

	OT_SkillNear0 = OT_SkillNear + 1

	OT_SkillNearK0 = OT_SkillNearK + 1
)

func getSkillSingleDamage(id int, data0 *models.GameData, data1 *models.GameData, limit int) int {
	distance := helper.GClientDistance(data0.X, data0.Y, data1.X, data1.Y)
	if int(distance) > limit {
		return -limit
	}
	var damage int
	switch id {
	case OT_SkillSingle0:
		manaAdd := int(0.5 * float32(data0.Spear.SMana))
		data0.Spear.SMana += manaAdd
		damage = getSpearValue(data0.Spear, data1.Shield)
		data0.Spear.SMana -= manaAdd
	case OT_SkillSingle1:
		metalAdd := int(1.5 * float32(data0.Spear.SMetal))
		woodAdd := int(1.5 * float32(data0.Spear.SWood))
		waterAdd := int(1.5 * float32(data0.Spear.SWater))
		fireAdd := int(1.5 * float32(data0.Spear.SFire))
		earthAdd := int(1.5 * float32(data0.Spear.SEarth))

		data0.Spear.SMetal += metalAdd
		data0.Spear.SWood += woodAdd
		data0.Spear.SWater += waterAdd
		data0.Spear.SFire += fireAdd
		data0.Spear.SEarth += earthAdd

		damage = getSpearValue(data0.Spear, data1.Shield)

		data0.Spear.SMetal -= metalAdd
		data0.Spear.SWood -= woodAdd
		data0.Spear.SWater -= waterAdd
		data0.Spear.SFire -= fireAdd
		data0.Spear.SEarth -= earthAdd
	default:
		damage = getSpearValue(data0.Spear, data1.Shield)
	}

	ran := rand.Intn(data0.Spear.SStrength / 10)
	if rand.Intn(10) < 5 {
		damage += ran
	} else {
		damage -= ran
	}
	if damage < 0 {
		damage = 0
	}
	return damage
}

func getSpearValue(spear *models.GameSpear, shield *models.GameShield) int {
	power := helper.Max(spear.SStrength-shield.SStrength, 0)
	power += helper.Max(spear.SMana-shield.SMana, 0) * 10

	fiveVal := getAbsFiveEleVal(spear, shield)
	power += fiveVal * 100
	return power
}

func getAbsFiveEleVal(spear *models.GameSpear, shield *models.GameShield) int {
	val := 0
	val += helper.Max(spear.SMetal-shield.SMetal, 0)
	val += helper.Max(spear.SWood-shield.SWood, 0)
	val += helper.Max(spear.SWater-shield.SWater, 0)
	val += helper.Max(spear.SFire-shield.SFire, 0)
	val += helper.Max(spear.SEarth-shield.SEarth, 0)
	return val
}
