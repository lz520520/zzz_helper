package calc

import "zzz_helper/internal/data"

func SheerBonusImprove(sheer float64, basicAttack float64, attackBonusCount float64, basicHP float64, hpBonusCount float64) float64 {
	if attackBonusCount > 0 {
		out := sheer + data.BaseDriverDiskSubStat.AttackBonus*attackBonusCount*0.3*basicAttack
		return out/sheer - 1
	}
	if hpBonusCount > 0 {
		out := sheer + data.BaseDriverDiskSubStat.HPBonus*hpBonusCount*0.1*basicHP
		return out/sheer - 1
	}
	return 0
}
