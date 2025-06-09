package calc

import (
	"zzz_helper/internal/data"
	"zzz_helper/internal/models"
)

func CriticalCalc(param models.CriticalParam) float64 {
	if param.CriticalRate > 1 {
		param.CriticalRate = 1
	}
	return 1 + param.CriticalDamage*param.CriticalRate
}

// 暴击/爆伤
func CriticalImprove(param models.CriticalParam, count int) (float64, float64) {
	base := CriticalCalc(param)

	rateParam := param
	rateParam.CriticalRate += data.BaseDriverDiskSubStat.CriticalRate * float64(count)
	rateOut := CriticalCalc(rateParam)

	damageParam := param
	damageParam.CriticalDamage += data.BaseDriverDiskSubStat.CriticalDamage * float64(count)
	damageOut := CriticalCalc(damageParam)
	return rateOut/base - 1, damageOut/base - 1
}

func CriticalRateExpect(criticalDamage float64) (criticalRate float64) {
	criticalRate = data.BaseDriverDiskSubStat.CriticalRate / data.BaseDriverDiskSubStat.CriticalDamage * criticalDamage
	if criticalRate > 1 {
		criticalRate = 1
	}
	return
}

func CriticalDamageExpect(criticalRate float64) (criticalDamage float64) {
	criticalDamage = data.BaseDriverDiskSubStat.CriticalDamage / data.BaseDriverDiskSubStat.CriticalRate * criticalRate
	return
}

/*
(1 + criticalDamage + 0.048) * criticalRate  / (1 + criticalDamage) * criticalRate  = improve +1
criticalDamage =  (0.048 - improve) / improve
*/
func CriticalDamageExpect2(improve float64, count float64) (criticalDamage float64) {
	criticalDamage = (data.BaseDriverDiskSubStat.CriticalDamage - improve*count) / improve * count
	return
}
