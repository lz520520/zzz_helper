package calc

import (
	"testing"
)

func TestSheer(t *testing.T) {
	baseSheer := 2280.0 + 720

	hpResult := SheerBonusImprove(baseSheer, 0, 0, 8373, 1)

	attackResult := SheerBonusImprove(baseSheer, 1615, 1, 0, 0)

	expectCriticalDamage := CriticalDamageExpect2(hpResult, 1)
	expectDamageBonus := DamageExpect(hpResult*10, 0.3)

	t.Logf(`
生命词条贯穿力收益: %.2f%%
攻击词条贯穿力收益: %.2f%%
期望爆伤：%.2f%%
期望增伤: %.2f%%
`, hpResult*100, attackResult*100, expectCriticalDamage*100, expectDamageBonus*100)

}
