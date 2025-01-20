package models

type CriticalParam struct {
	CriticalDamage float64 `yaml:"critical_damage,omitempty"` // 暴击伤害
	CriticalRate   float64 `yaml:"critical_rate,omitempty"`   // 暴击率
}
