package models

type DefenseParam struct {
	LevelBase          float64
	MonsterBaseDefense float64
	DefenseReduction   float64 // 减防
	PenetrationRadio   float64 // 穿透率
	Penetration        float64 // 穿透值
}
