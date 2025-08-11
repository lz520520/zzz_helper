package main

import (
	"testing"
	"zzz_helper/modules/zzz/calc"
	models2 "zzz_helper/modules/zzz/models"
)

func TestMiyabi(t *testing.T) {
	err := calc.DamageFuzz(models2.DamageFuzzParam{
		Name:      "星见雅",
		Attribute: models2.AttrIce,
		Star:      6,
		Engine:    "霰落星殿",
		//Improve:   true,
		//DriverPath: "conf/drivers_miyabi.yml",
		DriverType: []string{"折枝剑歌", "河豚电音"},
		TestData: models2.TestData{
			LevelBase:          794,
			MonsterBaseDefense: 60,
			DamageMultiplier:   1, // 普攻第一段伤害
			Attribute: models2.CommonAttribute{
				OutGame: models2.AgentAttribute{},
				InGame: models2.AgentAttribute{
					//AttackBonus:      0.128,
					//DamageResistance: -0.2, // 怪物减抗
					IceDamageBonus: 0.3,
					Attack:         1200 + 1000, // 柚叶buff+苍角buff
					//CriticalDamage:    0.25,
					//CommonDamageBonus: 0.2, // 佳音
					CommonDamageBonus: 0.2, // 苍角核心被动额外能力增加冰伤20%
				},
			},
		},
		//DriverFilter: func(disks []models.DriverDiskStat) bool {
		//    if disks[4].Main.PenetrationRadio > 0 {
		//        return true
		//    }
		//    return false
		//},
		//InGameAttrFilter: func(attr models.AgentAttribute) bool {
		//    if attr.CriticalRate >= 0.80+0.12 {
		//        return true
		//    }
		//    return false
		//},
	})
	if err != nil {
		t.Error(err)
	}
}
