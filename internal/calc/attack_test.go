package calc

import (
	"testing"
)

func TestAttack(t *testing.T) {
	currentAttack := 2891.0
	basicAttack := 1651.0
	result := AttackImprove(currentAttack)
	t.Logf("attack: %.2f%%", result*100)

	result = AttackBonusImprove(basicAttack, currentAttack, 1)
	t.Logf("attack bonus: %.2f%%", result*100)

	result = AttackBonusImprove(basicAttack, currentAttack, 10)
	t.Logf("attack bonus with main: %.2f%%", result*100)
}

func TestAttack2(t *testing.T) {
	result := AttackBonusImprove(1651, 3187, 1)
	t.Logf("attack bonus: %.2f%%", result*100)

	criticalDamage := CriticalDamageExpect(0.922)
	t.Logf("CriticalDamage expect: %.2f%%", criticalDamage*100)

	criticalRate := CriticalRateExpect(2.0040)
	t.Logf("CriticalRate expect: %.2f%%", criticalRate*100)
}
