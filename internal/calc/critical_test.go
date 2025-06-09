package calc

import (
	"testing"
	"zzz_helper/internal/models"
)

func TestCriticalImprove(t *testing.T) {
	criticalDamage := 1.2 + 0.8 + 0.3 + 0.72
	criticalRate := 0.73 + 0.27
	attack := 2920
	rate, damage := CriticalImprove(models.CriticalParam{
		CriticalDamage: criticalDamage,
		CriticalRate:   criticalRate,
	}, 1)
	attackImprove := AttackBonusImprove(1623, float64(attack), 1)

	t.Logf("局内暴击: %.2f%%, 局内爆伤: %.2f%%, 局内攻击: %v", criticalRate*100, criticalDamage*100, attack)
	t.Logf("暴击词条收益: %.2f%%, 爆伤词条收益: %.2f%%, 大攻击词条收益: %.2f%%", rate*100, damage*100, attackImprove*100)

	expectAttack := AttackBonusExpected(1623, damage, 1)
	t.Logf("期望攻击力: %v", expectAttack)

}

func TestYifuCritical(t *testing.T) {
	// 基础 + 专武 + 4号盘 + 嘉音buff
	criticalDamage := 0.5 + 0.5 + 0.48 + 0.25
	// 基础 + 被动 + 专武 + 被动*3 + 啄木鸟2
	criticalRate := 0.05 + 0.25 + 0.24 + 0.048*3 + 0.08

	// 每个驱动盘双爆副词条 5个，共30个副词条
	criticalRate += 10 * 0.024
	criticalDamage += 20 * 0.048

	t.Logf("暴击率: %.2f%%, 爆伤: %.2f%%", criticalRate*100, criticalDamage*100)

	improve := (1+criticalDamage)/(1+criticalDamage-0.16) - 1

	t.Logf("鸟2收益: %.2f%%", improve*100)

}
