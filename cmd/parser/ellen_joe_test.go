package main

import (
	"os"
	"testing"
	"zzz_helper/modules/zzz/calc"
	models2 "zzz_helper/modules/zzz/models"
)

func TestEllenJoe(t *testing.T) {
	_, err := calc.DamageFuzz(models2.DamageFuzzParam{
		Name:      "艾莲·乔",
		Attribute: models2.AttrIce,
		Star:      2,
		Engine:    "深海访客",
		TestData: models2.TestData{
			LevelBase:          794,
			MonsterBaseDefense: 60,
			DamageMultiplier:   1.997, // 普攻第一段伤害
			Attribute: models2.CommonAttribute{
				OutGame: models2.AgentAttribute{},
				InGame: models2.AgentAttribute{
					//AttackBonus:      0.128,
					//DamageResistance: -0.2, // 怪物减抗
					//CommonDamageBonus:      -0.3,
					Attack:            1600,             // 嘉音buff
					CommonDamageBonus: 0.2 + 0.24 + 0.9, // 嘉音增伤+莱特
					CriticalDamage:    0.25 - 0.4,       // 嘉音爆伤
				},
			},
		},
	}, os.Stdout)
	if err != nil {
		t.Error(err)
	}
}

//
//func DamageCompare(
//    param1 models.DamageParam,
//    param2 models.DamageParam,
//    base models.DamageParam) error {
//    output1, err := calc.DamageCalc(param1, models.AttrIce, false, false)
//    if err != nil {
//        return err
//    }
//    output2, err := calc.DamageCalc(param2, models.AttrIce, false, false)
//    if err != nil {
//        return err
//    }
//    baseOutput, err := calc.DamageCalc(base, models.AttrIce, false, false)
//    if err != nil {
//        return err
//    }
//    relOutput1 := output1.Output - baseOutput.Output
//    relOutput2 := output2.Output - baseOutput.Output
//    if relOutput1 < 0 || relOutput2 < 0 {
//        return fmt.Errorf("invalid damage")
//    }
//
//    fmt.Printf(`
//base: %v, output1: %v, output2: %v
//output1 with base: %v
//output2 with base: %v
//compare: %v
//`, baseOutput, output1, output2,
//        relOutput1/baseOutput.Output, relOutput2/baseOutput.Output,
//        relOutput1/relOutput2)
//    return nil
//}
//
//func DamageImprove(param1 models.DamageParam, param2 models.DamageParam) error {
//    output1, err := calc.DamageCalc(param1, models.AttrIce, false, false)
//    if err != nil {
//        return err
//    }
//    output2, err := calc.DamageCalc(param2, models.AttrIce, false, false)
//    if err != nil {
//        return err
//    }
//    fmt.Printf(`
//output1: %v, output2: %v
//relative: %v
//`, output1, output2, (output1.Output-output2.Output)/output2.Output)
//    return nil
//}
