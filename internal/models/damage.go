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
	Name       string   `yaml:"name"`        // 代理人
	Attribute  string   `yaml:"attribute"`   // 属性
	Star       int      `yaml:"star"`        // 影画等级
	Engine     string   `yaml:"engine"`      // 音擎
	Improve    bool     `yaml:"improve"`     // 是否输出提升率
	Stun       bool     `yaml:"stun"`        // 是否计算失衡易伤
	DriverPath string   `yaml:"driver_path"` // 驱动盘配置文件路径
	TestData   TestData `yaml:"test_data"`   // 测试数据

	AgentFeatures AgentFeatures `yaml:"agent_features"`

	DriverFilter     func(disks []DriverDiskStat) bool
	InGameAttrFilter func(attr AgentAttribute) bool
}
