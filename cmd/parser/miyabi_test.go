package main

import (
	"testing"
	"zzz_helper/internal/calc"
	"zzz_helper/internal/models"
)

func TestMiyabi(t *testing.T) {
	err := calc.DamageFuzz(models.DamageFuzzParam{
		Name:      "星见雅",
		Attribute: models.AttrIce,
		Star:      6,
		Engine:    "霰落星殿",
		//Improve:   true,
		//DriverPath: "conf/drivers_miyabi.yml",
		DriverPath: "conf/drivers-2.yml",
		TestData: models.TestData{
			LevelBase:          794,
			MonsterBaseDefense: 60,
			DamageMultiplier:   1, // 普攻第一段伤害
			Attribute: models.CommonAttribute{
				OutGame: models.AgentAttribute{},
				InGame: models.AgentAttribute{
					//AttackBonus:      0.128,
					//DamageResistance: -0.2, // 怪物减抗
					IceDamageBonus: 0.3,
					Attack:         1200, // 耀佳音buff
					//CommonDamageBonus: 0.2,  // 苍角核心被动额外能力增加冰伤20%
				},
			},
		},
		DriverFilter: func(disks []models.DriverDiskStat) bool {
			if disks[4].Main.PenetrationRadio > 0 {
				return true
			}
			return false
		},
		InGameAttrFilter: func(attr models.AgentAttribute) bool {
			if attr.CriticalRate >= 0.80+0.12 {
				return true
			}
			return false
		},
	})
	if err != nil {
		t.Error(err)
	}
}
