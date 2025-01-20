package calc

import "zzz_helper/internal/data"

func AttackBonusImprove(basicAttack, attack float64, count int) float64 {
	out := attack + data.BaseDriverDiskSubStat.AttackBonus*float64(count)*basicAttack
	return out/attack - 1
}

func AttackBonusExpected(basicAttack, improve float64, count int) (attack float64) {
	attack = data.BaseDriverDiskSubStat.AttackBonus * float64(count) * basicAttack / improve
	return
}

func AttackImprove(attack float64) float64 {
	out := attack + data.BaseDriverDiskSubStat.Attack
	return out/attack - 1
}
