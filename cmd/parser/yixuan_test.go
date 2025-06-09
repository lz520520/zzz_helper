package main

import (
	"testing"
	"zzz_helper/internal/calc"
	"zzz_helper/internal/models"
)

func TestYixuan(t *testing.T) {
	err := calc.DamageFuzz(models.DamageFuzzParam{
		Name:      "仪玄",
		Attribute: models.AttrEther,
		Star:      2,
		Engine:    "青溟笼舍",
		//Improve:   true,
		//DriverPath: "conf/drivers_miyabi.yml",
		DriverPath: "../../conf/drivers_yixuan.yml",
		AgentFeatures: models.AgentFeatures{
			LifeDestroy: true,
			Attribute2Sheer: func(attr models.AgentAttribute) float64 {
				return attr.Attack*0.3 + attr.HP*0.1 + 720
			},
		},
		TestData: models.TestData{
			LevelBase:          794,
			MonsterBaseDefense: 60,
			DamageMultiplier:   1, // 普攻第一段伤害
			Attribute: models.CommonAttribute{
				OutGame: models.AgentAttribute{},
				InGame: models.AgentAttribute{
					CommonDamageBonus: 0.6, // 被动增伤
					CriticalDamage:    0.2, // 额外能力爆伤
					//SheerForce:        720, // 潘子提供
					//SheerForce:    720,
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
