package calc

import (
	"testing"
	"zzz_helper/internal/data"
	"zzz_helper/internal/models"
)

func TestDefense(t *testing.T) {
	t.Log(DefenseCalc(models.DefenseParam{
		LevelBase:          592,
		MonsterBaseDefense: 36,
		DefenseReduction:   0,
		PenetrationRadio:   0,
		Penetration:        21,
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
		//DefenseReduction:   0.36,
		PenetrationRadio: 0.00,
		Penetration:      0,
	})
	out2 := DefenseCalc(models.DefenseParam{
		LevelBase:          794,
		MonsterBaseDefense: 60,
		//DefenseReduction:   0.36,
		PenetrationRadio: 0.32,
		Penetration:      0,
	})
	improve := out2/out1 - 1
	t.Logf("miyabi improve with 32%%: %.2f%%", improve*100)
}

func TestDamage(t *testing.T) {
	basicDamage := 1.0
	add := 0.3

	improve := (basicDamage+1+add)/(basicDamage+1) - 1
	t.Logf("damage improve: %.2f%%", improve*100)

}
