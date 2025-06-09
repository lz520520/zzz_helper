package calc

import "testing"

func TestVvaAnomalyImprove(t *testing.T) {
	result := AnomalyImprove(456+145, 1)
	t.Logf("Anomaly: %.2f%%", result*100)

	currentAttack := 2453.0
	basicAttack := 880 + 713.0
	result = AttackBonusImprove(basicAttack, currentAttack, 2)
	t.Logf("attack bonus with main: %.2f%%", result*100)

}

func TestAnomalyExpected(t *testing.T) {
	result := AnomalyImprove(460, 1)
	attack := AttackBonusExpected(1651, result, 1)
	t.Logf("Anomaly: %.2f%%, expected attack: %v", result*100, attack)
}
