package calc

import (
	"testing"
	"zzz_helper/internal/models"
)

func TestCriticalImprove(t *testing.T) {
	rate, damage := CriticalImprove(models.CriticalParam{
		CriticalDamage: 1.47 + 0.8,
		CriticalRate:   0.698 + 0.27,
	}, 1)
	attack := AttackBonusExpected(1623, damage, 1)
	t.Logf("rate: %.2f%%, damage: %.2f%%", rate*100, damage*100)
	t.Logf("expect attack: %v", attack)

}
