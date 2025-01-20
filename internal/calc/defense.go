package calc

import (
	"zzz_helper/internal/data"
	"zzz_helper/internal/models"
)

func DefenseCalc(param models.DefenseParam) float64 {
	if param.PenetrationRadio > 1 {
		param.PenetrationRadio = 1
	}
	if param.DefenseReduction > 1 {
		param.DefenseReduction = 1
	}
	monsterBasicDefense := param.MonsterBaseDefense / 50 * param.LevelBase // 怪物基础防御力
	monsterValidDefense := monsterBasicDefense*                            // 怪物防御力
		(1-param.PenetrationRadio)* // 穿透率
		(1-param.DefenseReduction) - // 减防
		param.Penetration // 穿透值
	if monsterBasicDefense < 0 {
		monsterBasicDefense = 0
	}

	return param.LevelBase / (monsterValidDefense + param.LevelBase)
}

func DefensePenetrationImprove(param models.DefenseParam) float64 {
	out1 := DefenseCalc(param)
	param.Penetration += data.BaseDriverDiskSubStat.Penetration
	out2 := DefenseCalc(param)
	return out2/out1 - 1
}
