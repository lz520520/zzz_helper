package models

type DamageCalcResult struct {
	Output float64 `yaml:"output"`
	Damage Damage  `yaml:"damage"`

	OutGameAttr AgentAttribute `yaml:"out_game_attr"`
	InGameAttr  AgentAttribute `yaml:"in_game_attr"`
	BaseAttr    AgentAttribute `yaml:"base_attr"`
	Set         DriverDiskSet  `yaml:"set"`
}
type DamageFuzzDriver struct {
	Include map[int]DriverDiskMainStat
}

type DamageFuzzParam struct {
	Name       string `yaml:"name" json:"name"`     // 代理人
	Attribute  string `yaml:"attribute" json:"-"`   // 属性
	Star       int    `yaml:"star" json:"star"`     // 影画等级
	Engine     string `yaml:"engine" json:"engine"` // 音擎
	EngineStar int    `yaml:"engine_star" json:"engine_star"`

	Improve  bool     `yaml:"improve" json:"improve"`     // 是否输出提升率
	Stun     bool     `yaml:"stun" json:"stun"`           // 是否计算失衡易伤
	TestData TestData `yaml:"test_data" json:"test_data"` // 测试数据

	AgentFeatures AgentFeatures `yaml:"agent_features" json:"-"`

	DriverType []string `yaml:"driver_type" json:"driver_type"`
	DriverIds  []string `yaml:"driver_ids" json:"driver_ids"`

	DriverFilter     func(disks []DriverDiskStat) bool `yaml:"-" json:"-"`
	InGameAttrFilter func(attr AgentAttribute) bool    `yaml:"-" json:"-"`
}
