package calc

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"math"
	"zzz_helper/internal/data"
	"zzz_helper/internal/models"
)

func DamageFuzz(param models.DamageFuzzParam) error {
	agentInfo, err := data.AgentInfos.GetInfo(param.Name)
	if err != nil {
		return err
	}
	engineInfo, err := data.EngineInfos.GetInfo(param.Engine)
	if err != nil {
		return err
	}
	// 驱动盘

	agentInfo.Star = param.Star
	collection, err := data.CollectDriverDisks(param.DriverPath)
	if err != nil {
		return err
	}

	base := models.DamageParam{
		AgentInfo: agentInfo,
		//DriverDisks:  getEllenDisks(),
		WeaponEngine: engineInfo,
		TestData:     param.TestData,
	}

	maxOutput := 0.0
	result := &models.DamageCalcResult{}
	for _, disk1 := range collection.Disk1 {
		for _, disk2 := range collection.Disk2 {
			for _, disk3 := range collection.Disk3 {
				for _, disk4 := range collection.Disk4 {
					for _, disk5 := range collection.Disk5 {
						for _, disk6 := range collection.Disk6 {
							disks := []models.DriverDiskStat{disk1, disk2, disk3, disk4, disk5, disk6}
							if param.DriverFilter != nil && !param.DriverFilter(disks) {
								continue
							}
							set := models.DriverDiskSet{Disks: disks}
							err = set.Parse()
							if err != nil {
								return err
							}
							damageParam := base
							damageParam.DriverDisks = set
							damageParam.Attribute = param.Attribute
							damageParam.Improve = param.Improve
							damageParam.Stun = param.Stun
							tmp, err := DamageCalc(damageParam, param.InGameAttrFilter)
							if err == nil {
								fmt.Printf("伤害计算: %v\n", tmp.Output)
								if tmp.Output > maxOutput {
									maxOutput = tmp.Output
									result = tmp
									result.Set = set
								}
							}
						}
					}
				}
			}
		}
	}
	s, _ := yaml.Marshal(result)
	fmt.Printf(`
【期望】
%s
`, string(s))
	return nil
}

func DamageCalc(param models.DamageParam, inGameAttrFilter func(attr models.AgentAttribute) bool) (*models.DamageCalcResult, error) {
	var star = make([]models.AgentStar, 0)
	for _, s := range param.AgentInfo.Stars {
		if s.Level <= param.AgentInfo.Star {
			star = append(star, s)
		}
	}
	// 局外白值计算
	// 代理人基本属性 + 武器属性 + 核心被动属性 + 影画加成
	baseAttribute := param.AgentInfo.Attribute
	baseAttribute.Add(param.WeaponEngine.OutGame)
	baseAttribute.Add(param.AgentInfo.CorePassive.OutGame)
	for _, s := range star {
		baseAttribute.Add(s.OutGame)
	}

	agentOutGameAttribute := models.AgentAttribute{}
	agentOutGameAttribute.Add(baseAttribute)

	// 代理人黄值计算
	// 白值 + 驱动盘属性 + 驱动盘套装属性
	driverDiskFinalStat := models.DriverDiskStat{}
	for _, disk := range param.DriverDisks.Disks {
		driverDiskFinalStat.Add(disk)
	}
	driverDiskFinalStat.Add(param.TestData.Disk)

	agentOutGameAttribute.Attack = baseAttribute.Attack*
		(1+driverDiskFinalStat.Main.AttackBonus+driverDiskFinalStat.Sub.AttackBonus+param.DriverDisks.AttackBonus) +
		driverDiskFinalStat.Main.Attack + driverDiskFinalStat.Sub.Attack

	agentOutGameAttribute.CommonDamageBonus += driverDiskFinalStat.Main.CommonDamageBonus
	agentOutGameAttribute.IceDamageBonus += driverDiskFinalStat.Main.IceDamageBonus
	agentOutGameAttribute.FireDamageBonus += driverDiskFinalStat.Main.FireDamageBonus
	agentOutGameAttribute.PhysicalDamageBonus += driverDiskFinalStat.Main.PhysicalDamageBonus
	agentOutGameAttribute.ElectricDamageBonus += driverDiskFinalStat.Main.ElectricDamageBonus
	agentOutGameAttribute.EtherDamageBonus += driverDiskFinalStat.Main.EtherDamageBonus

	agentOutGameAttribute.CriticalRate += driverDiskFinalStat.Main.CriticalRate + driverDiskFinalStat.Sub.CriticalRate
	agentOutGameAttribute.CriticalDamage += driverDiskFinalStat.Main.CriticalDamage + driverDiskFinalStat.Sub.CriticalDamage

	agentOutGameAttribute.Penetration += driverDiskFinalStat.Sub.Penetration
	agentOutGameAttribute.PenetrationRadio += driverDiskFinalStat.Main.PenetrationRadio

	agentOutGameAttribute.AnomalyProficiency += driverDiskFinalStat.Main.AnomalyProficiency + driverDiskFinalStat.Sub.AnomalyProficiency
	agentOutGameAttribute.AnomalyMastery += driverDiskFinalStat.Main.AnomalyMastery

	agentOutGameAttribute.EnergyRegen += driverDiskFinalStat.Main.EnergyRegen

	agentOutGameAttribute.HP = baseAttribute.HP*
		(1+driverDiskFinalStat.Sub.HPBonus) +
		driverDiskFinalStat.Main.HP + driverDiskFinalStat.Sub.HP
	agentOutGameAttribute.Defense = baseAttribute.Defense*
		(1+driverDiskFinalStat.Sub.DefenseBonus) +
		driverDiskFinalStat.Main.Defense + driverDiskFinalStat.Sub.Defense
	// 驱动盘套装
	agentOutGameAttribute.Add(param.DriverDisks.SetAttribute.OutGame)
	agentOutGameAttribute.Add(param.TestData.Attribute.OutGame)

	if math.Round(agentOutGameAttribute.CriticalDamage*1000)/1000 == 1.604 &&
		math.Round(agentOutGameAttribute.CriticalRate*1000)/1000 == 0.698 &&
		math.Round(agentOutGameAttribute.Defense) == 951 {
		fmt.Println("critical damage")
	}
	//data1, _ := json.MarshalIndent(agentOutGameAttribute, "", "  ")
	//fmt.Printf("out: %s\n", string(data1))
	// 代理人局内面板 + 影画  + 武器属性 + 核心被动属性 + 驱动盘套装
	agentInGameAttribute := models.AgentAttribute{}
	agentInGameAttribute.Add(agentOutGameAttribute)
	agentInGameAttribute.Add(param.WeaponEngine.InGame)
	agentInGameAttribute.Add(param.AgentInfo.CorePassive.InGame)
	agentInGameAttribute.Add(param.DriverDisks.SetAttribute.InGame)

	for _, s := range star {
		agentInGameAttribute.Add(s.InGame)
	}
	agentInGameAttribute.Add(param.TestData.Attribute.InGame)

	agentInGameAttribute.Fix()
	if inGameAttrFilter != nil && !inGameAttrFilter(agentInGameAttribute) {
		return &models.DamageCalcResult{}, nil
	}
	//data1, _ = json.MarshalIndent(agentInGameAttribute, "", "  ")

	//fmt.Printf("in: %s\n", string(data1))
	damage := models.Damage{
		Attack:            agentInGameAttribute.Attack + agentInGameAttribute.AttackBonus*agentOutGameAttribute.Attack,
		DamageMultiplier:  param.TestData.DamageMultiplier,
		CriticalDamage:    agentInGameAttribute.CriticalDamage,
		CriticalRate:      agentInGameAttribute.CriticalRate,
		DamageResistance:  1 - agentInGameAttribute.DamageResistance,
		CommonDamageBonus: 1 + agentInGameAttribute.CommonDamageBonus,
		//DefenseReduction:     1 - agentInGameAttribute.DefenseReduction,
		StunDamageMultiplier: 1.5 + agentInGameAttribute.StunDamageMultiplier,
	}
	//减防计算
	damage.DefenseReduction = DefenseCalc(models.DefenseParam{
		LevelBase:          param.TestData.LevelBase,
		MonsterBaseDefense: param.TestData.MonsterBaseDefense,
		DefenseReduction:   agentInGameAttribute.DefenseReduction,
		PenetrationRadio:   agentInGameAttribute.PenetrationRadio,
		Penetration:        agentInGameAttribute.Penetration,
	})
	// 属性增伤计算
	damageBonus := 0.0
	switch param.Attribute {
	case models.AttrIce:
		damageBonus = agentInGameAttribute.IceDamageBonus
	case models.AttrElectric:
		damageBonus = agentInGameAttribute.ElectricDamageBonus
	case models.AttrPhysical:
		damageBonus = agentInGameAttribute.PhysicalDamageBonus
	case models.AttrFire:
		damageBonus = agentInGameAttribute.FireDamageBonus
	case models.AttrEther:
		damageBonus = agentInGameAttribute.EtherDamageBonus
	}
	damage.CommonDamageBonus += damageBonus
	//damage.CriticalRate = 1
	// 输出计算
	// 爆伤 = (1 + CriticalDamage) * CriticalRate + 1 * (1-CriticalRate) =  CriticalDamage * CriticalRate + 1
	output := damage.Attack *
		damage.DamageMultiplier *
		(damage.CriticalDamage*damage.CriticalRate + 1) *
		damage.DamageResistance *
		damage.CommonDamageBonus *
		damage.DefenseReduction
	if param.Stun {
		output = output * damage.StunDamageMultiplier
	}
	//data1, _ = json.MarshalIndent(damage, "", "  ")

	// fmt.Printf("damage: \n%s\noutput: %v\n", string(data1), output)
	//
	// 收益计算
	/*
	   x: 爆伤 y: 暴击率
	   1 + (x + 0.048)*y >= 1+ (x) * (y + 0.024)
	   y >= x/2

	   x: 爆伤 y: 暴击率 z: 攻击 base: 基础攻击力
	   期望爆伤: xy = x*y+1
	   z * ((x + 0.048)*y+1) = (z + base *0.03) * (x*y+1)
	   z = (0.03 * base) (1+x * y) / 0.048y
	*/
	//attackExpected := baseAttribute.Attack * data.BaseDriverDiskSubStat.AttackBonus *
	//    (1 + damage.CriticalDamage*damage.CriticalRate) / (data.BaseDriverDiskSubStat.CriticalDamage * damage.CriticalRate)
	//criticalRateExpected := damage.CriticalDamage / 2
	//if criticalRateExpected > 1 {
	//    criticalRateExpected = 1
	//}
	if param.Improve {
		attackImprove := AttackBonusImprove(baseAttribute.Attack, damage.Attack, 1)
		criticalRateImprove, criticalDamageImprove := CriticalImprove(models.CriticalParam{
			CriticalDamage: damage.CriticalDamage,
			CriticalRate:   damage.CriticalRate,
		}, 1)
		penetrationImprove := DefensePenetrationImprove(models.DefenseParam{
			LevelBase:          param.TestData.LevelBase,
			MonsterBaseDefense: param.TestData.MonsterBaseDefense,
			DefenseReduction:   agentInGameAttribute.DefenseReduction,
			PenetrationRadio:   agentInGameAttribute.PenetrationRadio,
			Penetration:        agentInGameAttribute.Penetration,
		})

		criticalRateExpected := CriticalRateExpect(damage.CriticalDamage)
		attackExpected := AttackBonusExpected(baseAttribute.Attack, criticalDamageImprove, 1)

		fmt.Printf(`
当前局内
    暴伤: %.2f%%
    暴击: %.2f%%
    攻击: %v

爆伤词条收益:   %.2f%%
暴击词条收益:   %.2f%%
大攻击词条收益: %.2f%%
穿透值词条收益: %.2f%%

爆伤词条收益大于暴击率，需局内暴击率: %.2f%%
爆伤词条收益大于攻击力，需局内攻击力: %v
`,
			damage.CriticalDamage*100,
			damage.CriticalRate*100,
			damage.Attack,

			criticalDamageImprove*100,
			criticalRateImprove*100,
			attackImprove*100,
			penetrationImprove*100,

			criticalRateExpected*100,
			attackExpected,
		)
	}

	return &models.DamageCalcResult{
		Output:      output,
		OutGameAttr: agentOutGameAttribute,
		InGameAttr:  agentInGameAttribute,
		Damage:      damage,
	}, nil
}
