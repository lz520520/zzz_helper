package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v3"
)

type DriverDiskStat struct {
	Name     string             `yaml:"name,omitempty"`
	Position int                `yaml:"position,omitempty"`
	Main     DriverDiskMainStat `yaml:"main,omitempty"` // 主词条
	Sub      DriverDiskSubStat  `yaml:"sub,omitempty"`  // 副词条
}

func (d *DriverDiskStat) Hash() string {
	b, _ := yaml.Marshal(d)
	hash := md5.Sum(b)
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

func (d *DriverDiskStat) Add(other DriverDiskStat) {
	d.Main.Add(other.Main)
	d.Sub.Add(other.Sub)
}

type DriverDiskSet struct {
	Disks []DriverDiskStat

	AttackBonus  float64 // 攻击加成
	HPBonus      float64
	SetAttribute CommonAttribute
}

func (d *DriverDiskSet) Parse() error {
	if len(d.Disks) > 6 || len(d.Disks) == 0 {
		return fmt.Errorf("disks count is error")
	}
	setAttr := SetAttribute{}
	for _, disk := range d.Disks {
		switch disk.Name {
		case "极地重金属":
			setAttr.PolarMetal += 1
		case "河豚电音":
			setAttr.PufferElectro += 1
		case "啄木鸟电音":
			setAttr.WoodpeckerElectro += 1
		case "折枝剑歌":
			setAttr.BranchBladeSong += 1
		case "云岿如我":
			setAttr.YKRW += 1
		}
	}
	err := setAttr.Parse(d)
	if err != nil {
		return err
	}
	for _, disk := range d.Disks {
		d.SetAttribute.OutGame.AddDisk(disk)
	}
	return nil
}

type SetAttribute struct {
	PolarMetal        int // 极地重金属
	PufferElectro     int // 河豚电音
	WoodpeckerElectro int // 啄木鸟电音
	BranchBladeSong   int // 折枝剑歌
	YKRW              int // 云岿如我
}

func (this *SetAttribute) Parse(set *DriverDiskSet) error {
	if this.PolarMetal >= 2 {
		info, _ := DriversInfos.GetInfo("极地重金属")
		set.SetAttribute.OutGame.Add(info.Piece2.OutGame)
		set.SetAttribute.InGame.Add(info.Piece2.InGame)
	}
	if this.PolarMetal >= 4 {
		info, _ := DriversInfos.GetInfo("极地重金属")
		set.SetAttribute.OutGame.Add(info.Piece4.OutGame)
		set.SetAttribute.InGame.Add(info.Piece4.InGame)
	}

	if this.PufferElectro >= 2 {
		info, _ := DriversInfos.GetInfo("河豚电音")
		set.SetAttribute.OutGame.Add(info.Piece2.OutGame)
		set.SetAttribute.InGame.Add(info.Piece2.InGame)
	}
	if this.PufferElectro >= 4 {
		info, _ := DriversInfos.GetInfo("河豚电音")
		set.SetAttribute.OutGame.Add(info.Piece4.OutGame)
		set.SetAttribute.InGame.Add(info.Piece4.InGame)
	}

	if this.BranchBladeSong >= 2 {
		info, _ := DriversInfos.GetInfo("折枝剑歌")
		set.SetAttribute.OutGame.Add(info.Piece2.OutGame)
		set.SetAttribute.InGame.Add(info.Piece2.InGame)
	}
	if this.BranchBladeSong >= 4 {
		info, _ := DriversInfos.GetInfo("折枝剑歌")
		set.SetAttribute.OutGame.Add(info.Piece4.OutGame)
		set.SetAttribute.InGame.Add(info.Piece4.InGame)
	}

	if this.WoodpeckerElectro >= 2 {
		info, _ := DriversInfos.GetInfo("啄木鸟电音")
		set.SetAttribute.OutGame.Add(info.Piece2.OutGame)
		set.SetAttribute.InGame.Add(info.Piece2.InGame)
	}
	if this.WoodpeckerElectro >= 4 {
		info, _ := DriversInfos.GetInfo("啄木鸟电音")
		set.SetAttribute.OutGame.Add(info.Piece4.OutGame)
		set.SetAttribute.InGame.Add(info.Piece4.InGame)
	}

	if this.YKRW >= 2 {
		info, _ := DriversInfos.GetInfo("云岿如我")
		set.SetAttribute.OutGame.Add(info.Piece2.OutGame)
		set.SetAttribute.InGame.Add(info.Piece2.InGame)
	}
	if this.YKRW >= 4 {
		info, _ := DriversInfos.GetInfo("云岿如我")
		set.SetAttribute.OutGame.Add(info.Piece4.OutGame)
		set.SetAttribute.InGame.Add(info.Piece4.InGame)
	}
	return nil
}

type DriverSetAttribute struct {
	Name   string          `yaml:"name,omitempty"`
	Piece2 CommonAttribute `yaml:"piece2,omitempty"`
	Piece4 CommonAttribute `yaml:"piece4,omitempty"`
}

type DriverDiskMainStat struct {
	HP           float64 `yaml:"hp,omitempty"`            // 生命值
	HPBonus      float64 `yaml:"hp_bonus,omitempty"`      // 生命加成
	Attack       float64 `yaml:"attack,omitempty"`        // 攻击力
	Defense      float64 `yaml:"defense,omitempty"`       // 防御值
	DefenseBonus float64 `yaml:"defense_bonus,omitempty"` // 防御加成

	CriticalDamage     float64 `yaml:"critical_damage,omitempty"`     // 暴击伤害
	CriticalRate       float64 `yaml:"critical_rate,omitempty"`       // 暴击率
	AnomalyProficiency float64 `yaml:"anomaly_proficiency,omitempty"` // 异常精通

	CommonDamageBonus   float64 `yaml:"common_damage_bonus,omitempty"`   // 伤害加成
	IceDamageBonus      float64 `yaml:"ice_damage_bonus,omitempty"`      // 冰属性伤害加成
	ElectricDamageBonus float64 `yaml:"electric_damage_bonus,omitempty"` // 电属性伤害加成
	PhysicalDamageBonus float64 `yaml:"physical_damage_bonus,omitempty"` // 物理属性伤害加成
	FireDamageBonus     float64 `yaml:"fire_damage_bonus,omitempty"`     // 火属性伤害加成
	EtherDamageBonus    float64 `yaml:"ether_damage_bonus,omitempty"`    // 以太伤害加成

	CommonSheerDamageBonus   float64 `yaml:"common_sheer_damage_bonus,omitempty"`   // 贯穿伤害加成
	IceSheerDamageBonus      float64 `yaml:"ice_sheer_damage_bonus,omitempty"`      // 贯穿冰属性伤害加成
	ElectricSheerDamageBonus float64 `yaml:"electric_sheer_damage_bonus,omitempty"` // 贯穿电属性伤害加成
	PhysicalSheerDamageBonus float64 `yaml:"physical_sheer_damage_bonus,omitempty"` // 贯穿物理属性伤害加成
	FireSheerDamageBonus     float64 `yaml:"fire_sheer_damage_bonus,omitempty"`     // 贯穿火属性伤害加成
	EtherSheerDamageBonus    float64 `yaml:"ether_sheer_damage_bonus,omitempty"`    // 贯穿以太伤害加成

	AttackBonus      float64 `yaml:"attack_bonus,omitempty"`      // 攻击加成
	PenetrationRadio float64 `yaml:"penetration_radio,omitempty"` // 穿透率

	AnomalyMastery float64 `yaml:"anomaly_mastery,omitempty"` // 异常掌控
	EnergyRegen    float64 `yaml:"energy_regen,omitempty"`    // 能量回复
}

func (d *DriverDiskMainStat) Add(other DriverDiskMainStat) {
	d.HP += other.HP
	d.Attack += other.Attack
	d.Defense += other.Defense
	d.CriticalDamage += other.CriticalDamage
	d.CriticalRate += other.CriticalRate
	d.AnomalyProficiency += other.AnomalyProficiency

	d.CommonDamageBonus += other.CommonDamageBonus
	d.IceDamageBonus += other.IceDamageBonus
	d.ElectricDamageBonus += other.ElectricDamageBonus
	d.PhysicalDamageBonus += other.PhysicalDamageBonus
	d.FireDamageBonus += other.FireDamageBonus
	d.EtherDamageBonus += other.EtherDamageBonus

	d.CommonSheerDamageBonus += other.CommonSheerDamageBonus
	d.IceSheerDamageBonus += other.IceSheerDamageBonus
	d.ElectricSheerDamageBonus += other.ElectricSheerDamageBonus
	d.PhysicalSheerDamageBonus += other.PhysicalSheerDamageBonus
	d.FireSheerDamageBonus += other.FireSheerDamageBonus
	d.EtherSheerDamageBonus += other.EtherSheerDamageBonus

	d.AttackBonus += other.AttackBonus
	d.HPBonus += other.HPBonus

	d.PenetrationRadio += other.PenetrationRadio
	d.AnomalyMastery += other.AnomalyMastery
	d.EnergyRegen += other.EnergyRegen
}

type DriverDiskSubStat struct {
	Attack      float64 `yaml:"attack,omitempty"`       // 攻击力
	AttackBonus float64 `yaml:"attack_bonus,omitempty"` // 攻击加成
	Penetration float64 `yaml:"penetration,omitempty"`  // 穿透值

	CriticalDamage float64 `yaml:"critical_damage,omitempty"` // 暴击伤害
	CriticalRate   float64 `yaml:"critical_rate,omitempty"`   // 暴击率

	AnomalyProficiency float64 `yaml:"anomaly_proficiency,omitempty"` // 异常精通
	Defense            float64 `yaml:"defense,omitempty"`             // 防御值
	DefenseBonus       float64 `yaml:"defense_bonus,omitempty"`       // 防御加成

	HP      float64 `yaml:"hp,omitempty"`       // 生命值
	HPBonus float64 `yaml:"hp_bonus,omitempty"` // 生命值

}

func (d *DriverDiskSubStat) Add(other DriverDiskSubStat) {
	d.Attack += other.Attack
	d.AttackBonus += other.AttackBonus
	d.Penetration += other.Penetration
	d.CriticalDamage += other.CriticalDamage
	d.CriticalRate += other.CriticalRate
	d.AnomalyProficiency += other.AnomalyProficiency
	d.Defense += other.Defense
	d.DefenseBonus += other.DefenseBonus

	d.HP += other.HP
	d.HPBonus += other.HPBonus
}

type DiskCollection struct {
	Disk1 []DriverDiskStat
	Disk2 []DriverDiskStat
	Disk3 []DriverDiskStat
	Disk4 []DriverDiskStat
	Disk5 []DriverDiskStat
	Disk6 []DriverDiskStat
}
