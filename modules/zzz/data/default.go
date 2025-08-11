package data

import (
	"zzz_helper/modules/zzz/models"
)

var (
	BaseDriverDiskSubStat = models.DriverDiskSubStat{
		Attack:             19,
		AttackBonus:        0.03,
		Penetration:        9,
		CriticalDamage:     0.048,
		CriticalRate:       0.024,
		AnomalyProficiency: 9,
		Defense:            15,
		DefenseBonus:       0.048,

		HP:      112,
		HPBonus: 0.03,
	}

	BaseDriverDiskMainStat = models.DriverDiskMainStat{
		HP:                 2200,
		HPBonus:            0.3,
		Attack:             316,
		Defense:            184,
		DefenseBonus:       0.48,
		CriticalDamage:     0.48,
		CriticalRate:       0.24,
		AnomalyProficiency: 92,

		CommonDamageBonus:   0.3,
		IceDamageBonus:      0.3,
		ElectricDamageBonus: 0.3,
		PhysicalDamageBonus: 0.3,
		FireDamageBonus:     0.3,
		EtherDamageBonus:    0.3,

		AttackBonus:      0.3,
		PenetrationRadio: 0.24,
		AnomalyMastery:   0.3,
		EnergyRegen:      0.6,
	}
)
