package main

import (
	"testing"
	"zzz_helper/modules/zzz/calc"
	models2 "zzz_helper/modules/zzz/models"
)

func TestYixuan(t *testing.T) {
	err := calc.DamageFuzz(models2.DamageFuzzParam{
		Name:      "仪玄",
		Attribute: models2.AttrEther,
		Star:      2,
		Engine:    "青溟笼舍",
		//Improve:   true,
		//DriverPath: "conf/drivers_miyabi.yml",
		DriverType: []string{"啄木鸟电音", "折枝剑歌", "云岿如我"},
		AgentFeatures: models2.AgentFeatures{
			LifeDestroy: true,
			Attribute2Sheer: func(attr models2.AgentAttribute) float64 {
				return attr.Attack*0.3 + attr.HP*0.1
			},
		},
		TestData: models2.TestData{
			LevelBase:          794,
			MonsterBaseDefense: 60,
			DamageMultiplier:   1, // 普攻第一段伤害
			Attribute: models2.CommonAttribute{
				OutGame: models2.AgentAttribute{},
				InGame: models2.AgentAttribute{
					CommonDamageBonus: 0.6,             // 被动增伤
					CriticalDamage:    0.2 + 0.5 + 0.3, // 额外能力爆伤+橘福福50%+30%
					//SheerForce:        720, // 潘子提供
					//SheerForce:    720,
				},
			},
		},
		DriverFilter: func(disks []models2.DriverDiskStat) bool {
			count := 0
			count2 := 0
			for _, disk := range disks {
				if disk.Name == "云岿如我" {
					count++
				} else if disk.Name == "啄木鸟电音" {
					count2++
				}
			}
			//if count2 == 2 && count == 4 {
			//    fmt.Println("zmn111111111111111")
			//}
			return count == 4
		},
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
