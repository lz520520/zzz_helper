package calc

import (
	"testing"
	"zzz_helper/modules/zzz/data"
	"zzz_helper/modules/zzz/models"
)

func TestDefense(t *testing.T) {
	t.Log(DefenseCalc(models.DefenseParam{
		LevelBase:          794,
		MonsterBaseDefense: 60,
		DefenseReduction:   0,
		PenetrationRadio:   1,
		Penetration:        0,
	}))
}

func TestDefensePenetrationImprove(t *testing.T) {
	for i := 0; i < 10; i++ {
		pen := data.BaseDriverDiskSubStat.Penetration * float64(i)
		result := DefensePenetrationImprove(models.DefenseParam{
			LevelBase:          794,
			MonsterBaseDefense: 58,
			DefenseReduction:   0.36,
			PenetrationRadio:   0.32 + 0.39,
			Penetration:        pen,
		})
		t.Logf("improve with %v: %.2f%%", pen, result*100)
	}

	result := DefensePenetrationImprove(models.DefenseParam{
		LevelBase:          794,
		MonsterBaseDefense: 58,
		DefenseReduction:   0,
		PenetrationRadio:   0.24,
		Penetration:        0,
	})
	t.Logf("improve with 24%%: %.2f%%", result*100)

	result = DefensePenetrationImprove(models.DefenseParam{
		LevelBase:          794,
		MonsterBaseDefense: 58,
		DefenseReduction:   0,
		PenetrationRadio:   0.32,
		Penetration:        0,
	})
	t.Logf("improve with 32%%: %.2f%%", result*100)
}

func TestMiyabi(t *testing.T) {
	out1 := DefenseCalc(models.DefenseParam{
		LevelBase:          794,
		MonsterBaseDefense: 60,
		DefenseReduction:   0.0,
		PenetrationRadio:   0.24,
		Penetration:        0,
	})
	out2 := DefenseCalc(models.DefenseParam{
		LevelBase:          794,
		MonsterBaseDefense: 60,
		DefenseReduction:   0.0,
		PenetrationRadio:   0.32,
		Penetration:        0,
	})
	improve := out2/out1 - 1
	t.Logf("穿透盘提升: %.2f%%", improve*100)

	expectDamage := 0.3 / improve
	t.Logf("属性盘期望: %.2f%%", (expectDamage-1)*100)

	expectAttack := AttackBonusExpected(1651, improve, 10)
	t.Logf("攻击盘期望: %v", expectAttack)

}

func TestDamage(t *testing.T) {
	basicDamage := 1.0
	add := 0.3

	improve := (basicDamage+1+add)/(basicDamage+1) - 1
	t.Logf("damage improve: %.2f%%", improve*100)

}
