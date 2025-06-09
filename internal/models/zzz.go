package models

const (
	AttrIce      = "ice"
	AttrElectric = "electric"
	AttrPhysical = "physical"
	AttrFire     = "fire"
	AttrEther    = "ether"
)

type Damage struct {
	Attack               float64 `yaml:"attack,omitempty"`                 // 攻击力
	DamageMultiplier     float64 `yaml:"damage_multiplier,omitempty"`      // 伤害倍率
	CriticalDamage       float64 `yaml:"critical_damage,omitempty"`        // 暴击伤害
	CriticalRate         float64 `yaml:"critical_rate,omitempty"`          // 暴击率
	DamageResistance     float64 `yaml:"damage_resistance,omitempty"`      // 抗性
	CommonDamageBonus    float64 `yaml:"common_damage_bonus,omitempty"`    // 伤害加成
	DefenseReduction     float64 `yaml:"defense_reduction,omitempty"`      // 防御减伤
	StunDamageMultiplier float64 `yaml:"stun_damage_multiplier,omitempty"` // 失衡易伤

	CommonSheerDamageBonus float64 `yaml:"common_sheer_damage_bonus,omitempty"` // 贯穿伤害加成
	SheerForce             float64 `yaml:"sheer_force,omitempty"`               // 贯穿力

	HP float64 `yaml:"hp,omitempty"` // 生命值

}
type AgentAttribute struct {
	HP                 float64 `yaml:"hp,omitempty"`                  // 生命值
	Attack             float64 `yaml:"attack,omitempty"`              // 攻击力
	Defense            float64 `yaml:"defense,omitempty"`             // 防御值
	Impact             float64 `yaml:"impact,omitempty"`              // 冲击力
	CriticalRate       float64 `yaml:"critical_rate,omitempty"`       // 暴击率
	CriticalDamage     float64 `yaml:"critical_damage,omitempty"`     // 暴击伤害
	AnomalyMastery     float64 `yaml:"anomaly_mastery,omitempty"`     // 异常掌控
	AnomalyProficiency float64 `yaml:"anomaly_proficiency,omitempty"` // 异常精通
	PenetrationRadio   float64 `yaml:"penetration_radio,omitempty"`   // 穿透率
	Penetration        float64 `yaml:"penetration,omitempty"`         // 穿透值
	EnergyRegen        float64 `yaml:"energy_regen,omitempty"`        // 能量回复
	SheerForce         float64 `yaml:"sheer_force,omitempty"`         // 贯穿力

	// 偏局内加成
	HPBonus      float64 `yaml:"hp_bonus,omitempty"`      // 生命值加成
	DefenseBonus float64 `yaml:"defense_bonus,omitempty"` // 防御力加成
	AttackBonus  float64 `yaml:"attack_bonus,omitempty"`  // 局内攻击加成

	CommonDamageBonus   float64 `yaml:"common_damage_bonus,omitempty"`   // 通用伤害加成
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

	DefenseReduction     float64 `yaml:"defense_reduction,omitempty"`      // 防御减伤
	DamageResistance     float64 `yaml:"damage_resistance,omitempty"`      // 抗性
	StunDamageMultiplier float64 `yaml:"stun_damage_multiplier,omitempty"` // 失衡易伤

}

func (a *AgentAttribute) Add(other AgentAttribute) {
	a.AddWithoutBonus(other)

	a.HPBonus += other.HPBonus
	a.AttackBonus += other.AttackBonus
	a.DefenseBonus += other.DefenseBonus
	//a.SheerForce += other.SheerForce
}

func (a *AgentAttribute) AddWithoutBonus(other AgentAttribute) {
	a.HP += other.HP
	a.Attack += other.Attack
	a.Defense += other.Defense

	a.Impact += other.Impact
	a.CriticalRate += other.CriticalRate
	a.CriticalDamage += other.CriticalDamage
	a.AnomalyMastery += other.AnomalyMastery
	a.AnomalyProficiency += other.AnomalyProficiency
	a.PenetrationRadio += other.PenetrationRadio
	a.Penetration += other.Penetration
	a.EnergyRegen += other.EnergyRegen

	a.CommonDamageBonus += other.CommonDamageBonus
	a.IceDamageBonus += other.IceDamageBonus
	a.ElectricDamageBonus += other.ElectricDamageBonus
	a.PhysicalDamageBonus += other.PhysicalDamageBonus
	a.FireDamageBonus += other.FireDamageBonus
	a.EtherDamageBonus += other.EtherDamageBonus

	a.CommonSheerDamageBonus += other.CommonSheerDamageBonus
	a.IceSheerDamageBonus += other.IceSheerDamageBonus
	a.ElectricSheerDamageBonus += other.ElectricSheerDamageBonus
	a.PhysicalSheerDamageBonus += other.PhysicalSheerDamageBonus
	a.FireSheerDamageBonus += other.FireSheerDamageBonus
	a.EtherSheerDamageBonus += other.EtherSheerDamageBonus

	a.DefenseReduction += other.DefenseReduction
	a.DamageResistance += other.DamageResistance
	a.StunDamageMultiplier += other.StunDamageMultiplier

}

func (this *AgentAttribute) AddDisk(set DriverDiskStat) {
	this.CommonDamageBonus += set.Main.CommonDamageBonus
	this.IceDamageBonus += set.Main.IceDamageBonus
	this.FireDamageBonus += set.Main.FireDamageBonus
	this.PhysicalDamageBonus += set.Main.PhysicalDamageBonus
	this.ElectricDamageBonus += set.Main.ElectricDamageBonus
	this.EtherDamageBonus += set.Main.EtherDamageBonus

	this.CommonSheerDamageBonus += set.Main.CommonSheerDamageBonus
	this.IceSheerDamageBonus += set.Main.IceSheerDamageBonus
	this.ElectricSheerDamageBonus += set.Main.ElectricSheerDamageBonus
	this.PhysicalSheerDamageBonus += set.Main.PhysicalSheerDamageBonus
	this.FireSheerDamageBonus += set.Main.FireSheerDamageBonus
	this.EtherSheerDamageBonus += set.Main.EtherSheerDamageBonus

	//agentOutGameAttribute.EtherSheerDamageBonus += set.Main.EtherSheerDamageBonus

	this.AttackBonus += set.Main.AttackBonus + set.Sub.AttackBonus
	this.HPBonus += set.Main.HPBonus + set.Sub.HPBonus
	this.DefenseBonus += set.Main.DefenseBonus + set.Sub.DefenseBonus

	this.Attack += set.Main.Attack + set.Sub.Attack
	this.HP += set.Main.HP + set.Sub.HP
	this.Defense += set.Main.Defense + set.Sub.Defense

	this.CriticalRate += set.Main.CriticalRate + set.Sub.CriticalRate
	this.CriticalDamage += set.Main.CriticalDamage + set.Sub.CriticalDamage

	this.Penetration += set.Sub.Penetration
	this.PenetrationRadio += set.Main.PenetrationRadio

	this.AnomalyProficiency += set.Main.AnomalyProficiency + set.Sub.AnomalyProficiency
	this.AnomalyMastery += set.Main.AnomalyMastery

	this.EnergyRegen += set.Main.EnergyRegen
}

func (a *AgentAttribute) Fix() {
	if a.CriticalRate > 1 {
		a.CriticalRate = 1
	}
	if a.PenetrationRadio > 1 {
		a.PenetrationRadio = 1
	}
}

func (a *AgentAttribute) AddEngine(other WeaponEngine) {
	a.Add(other.OutGame)
}

type CommonAttribute struct {
	OutGame AgentAttribute `yaml:"out_game,omitempty"`
	InGame  AgentAttribute `yaml:"in_game,omitempty"`
}

type AgentInfo struct {
	Name        string          `yaml:"name"` // 角色名称
	Star        int             `yaml:"-"`
	Stars       []AgentStar     `yaml:"stars"`
	Attribute   AgentAttribute  `yaml:"attribute,omitempty"`
	CorePassive CommonAttribute `yaml:"core_passive,omitempty"`

	Features AgentFeatures `yaml:"-"` // 代理人特性

}

type AgentStar struct {
	Level   int            `yaml:"level"`
	OutGame AgentAttribute `yaml:"out_game,omitempty"`
	InGame  AgentAttribute `yaml:"in_game,omitempty"`
}
type DamageParam struct {
	AgentInfo    AgentInfo
	DriverDisks  DriverDiskSet
	WeaponEngine WeaponEngine

	Attribute string // 属性
	Improve   bool   // 提升
	Stun      bool   // 失衡状态

	TestData TestData
}
type TestData struct {
	LevelBase          float64 // 等级基数
	MonsterBaseDefense float64 // 怪物基础防御

	DamageMultiplier float64 // 技能伤害倍率
	Disk             DriverDiskStat
	Attribute        CommonAttribute
}

type WeaponEngine struct {
	Name    string         `yaml:"name"`
	OutGame AgentAttribute `yaml:"out_game,omitempty"`
	InGame  AgentAttribute `yaml:"in_game,omitempty"`
}

type AgentFeatures struct {
	LifeDestroy     bool                              `yaml:"life_destroy"`
	Attribute2Sheer func(attr AgentAttribute) float64 `yaml:"-"`
}
