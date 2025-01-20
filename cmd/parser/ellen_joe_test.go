package main

import (
	"testing"
	"zzz_helper/internal/calc"
	"zzz_helper/internal/models"
)

func TestEllenJoe(t *testing.T) {
	err := calc.DamageFuzz(models.DamageFuzzParam{
		Name:       "艾莲·乔",
		Attribute:  models.AttrIce,
		Star:       1,
		Engine:     "深海访客",
		DriverPath: "conf/drivers_all.yml",
		TestData: models.TestData{
			LevelBase:          794,
			MonsterBaseDefense: 60,
			DamageMultiplier:   1.997, // 普攻第一段伤害
			Attribute: models.CommonAttribute{
				OutGame: models.AgentAttribute{},
				InGame: models.AgentAttribute{
					//AttackBonus:      0.128,
					DamageResistance: -0.2, // 怪物减抗
					//CommonDamageBonus:      -0.3,
					Attack:            1000, // 苍角buff
					CommonDamageBonus: 0.2,  // 苍角核心被动额外能力增加冰伤20%
				},
			},
		},
	})
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
