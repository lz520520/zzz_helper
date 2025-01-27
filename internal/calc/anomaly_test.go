package calc

import "testing"

func TestAnomalyImprove(t *testing.T) {
	result := AnomalyImprove(500, 1)
	t.Logf("Anomaly: %.2f%%", result*100)

}

func TestAnomalyExpected(t *testing.T) {
	result := AnomalyImprove(556, 1)
	attack := AttackBonusExpected(1651, result, 1)
	t.Logf("Anomaly: %.2f%%, expected attack: %v", result*100, attack)
}
