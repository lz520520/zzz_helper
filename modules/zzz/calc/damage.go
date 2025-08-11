package calc

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	data2 "zzz_helper/modules/zzz/data"
	models2 "zzz_helper/modules/zzz/models"
)

func DamageFuzz(param models2.DamageFuzzParam, writer io.StringWriter) (*models2.DamageCalcResult, error) {
	agentInfo, err := data2.AgentInfos.GetInfo(param.Name)
	if err != nil {
		return nil, err
	}
	engineInfo, err := data2.EngineInfos.GetInfo(param.Engine, param.EngineStar)
	if err != nil {
		return nil, err
	}
	// 驱动盘

	agentInfo.Star = param.Star
	collection, count, err := data2.CollectDriverDisks(param.DriverType, param.DriverIds)
	if err != nil {
		return nil, err
	}
	writer.WriteString(fmt.Sprintf("[*] 共%d个驱动盘", count))
	agentInfo.Features = param.AgentFeatures

	base := models2.DamageParam{
		AgentInfo: agentInfo,
		//DriverDisks:  getEllenDisks(),
		WeaponEngine: engineInfo,
		TestData:     param.TestData,
	}

	maxOutput := 0.0
	result := &models2.DamageCalcResult{}
	for _, disk1 := range collection.Disk1 {
		for _, disk2 := range collection.Disk2 {
			for _, disk3 := range collection.Disk3 {
				for _, disk4 := range collection.Disk4 {
					for _, disk5 := range collection.Disk5 {
						for _, disk6 := range collection.Disk6 {
							disks := []models2.DriverDiskStat{disk1, disk2, disk3, disk4, disk5, disk6}
							if param.DriverFilter != nil && !param.DriverFilter(disks) {
								continue
							}
							set := models2.DriverDiskSet{Disks: disks}
							err = set.Parse()
							if err != nil {
								return nil, err
							}
							damageParam := base
							damageParam.DriverDisks = set
							damageParam.Attribute = param.Attribute
							damageParam.Improve = param.Improve
							damageParam.Stun = param.Stun
							tmp, err := DamageCalc(damageParam, param.InGameAttrFilter)
							if err == nil {
								writer.WriteString(fmt.Sprintf("[*] 伤害计算: %v", tmp.Output))

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
	writer.WriteString(fmt.Sprintf(`[+] 
【期望】
%s
`, string(s)))
	return result, nil
}

func DamageCalc(param models2.DamageParam, inGameAttrFilter func(attr models2.AgentAttribute) bool) (*models2.DamageCalcResult, error) {
	var star = make([]models2.AgentStar, 0)
	for _, s := range param.AgentInfo.Stars {
		if s.Level <= param.AgentInfo.Star {
			star = append(star, s)
		}
	}
	/*
	   局外
	*/
	// 局外白值计算
	// 代理人基本属性 + 武器属性 + 核心被动属性 + 影画加成
	baseAttribute := param.AgentInfo.Attribute
	baseAttribute.Add(param.WeaponEngine.OutGame)
	baseAttribute.Add(param.AgentInfo.CorePassive.OutGame)
	for _, s := range star {
		baseAttribute.Add(s.OutGame)
	}

	agentOutGameAttribute := models2.AgentAttribute{}
	agentOutGameAttribute.Add(baseAttribute)

	// 代理人黄值计算
	// 白值 + 驱动盘属性 + 驱动盘套装属性
	// 驱动盘套装
	agentOutGameAttribute.Add(param.DriverDisks.SetAttribute.OutGame)
	agentOutGameAttribute.Add(param.TestData.Attribute.OutGame)
	agentOutGameAttribute.AddDisk(param.TestData.Disk)

	// 计算局外黄字
	agentOutGameAttribute.Attack = baseAttribute.Attack*agentOutGameAttribute.AttackBonus +
		agentOutGameAttribute.Attack

	agentOutGameAttribute.HP = baseAttribute.HP*agentOutGameAttribute.HPBonus +
		agentOutGameAttribute.HP

	agentOutGameAttribute.Defense = baseAttribute.Defense*agentOutGameAttribute.DefenseBonus +
		agentOutGameAttribute.Defense
	if param.AgentInfo.Features.Attribute2Sheer != nil {
		agentOutGameAttribute.SheerForce = param.TestData.Attribute.OutGame.SheerForce + param.AgentInfo.Features.Attribute2Sheer(agentOutGameAttribute)
	}

	//if math.Round(agentOutGameAttribute.CriticalDamage*1000)/1000 == 1.604 &&
	//    math.Round(agentOutGameAttribute.CriticalRate*1000)/1000 == 0.698 &&
	//    math.Round(agentOutGameAttribute.Defense) == 951 {
	//    fmt.Println("critical damage")
	//}
	//data1, _ := json.MarshalIndent(agentOutGameAttribute, "", "  ")
	//fmt.Printf("out: %s\n", string(data1))
	/*
	   局内
	*/
	// 代理人局内面板 + 影画  + 武器属性 + 核心被动属性 + 驱动盘套装
	agentInGameAttribute := models2.AgentAttribute{}
	agentInGameAttribute.AddWithoutBonus(agentOutGameAttribute)
	agentInGameAttribute.Add(param.WeaponEngine.InGame)
	agentInGameAttribute.Add(param.AgentInfo.CorePassive.InGame)
	agentInGameAttribute.Add(param.DriverDisks.SetAttribute.InGame)

	for _, s := range star {
		agentInGameAttribute.Add(s.InGame)
	}
	agentInGameAttribute.Add(param.TestData.Attribute.InGame)

	agentInGameAttribute.Fix()
	if param.AgentInfo.Features.Attribute2Sheer != nil {
		agentInGameAttribute.SheerForce = param.TestData.Attribute.InGame.SheerForce + param.AgentInfo.Features.Attribute2Sheer(agentInGameAttribute)
	}

	if inGameAttrFilter != nil && !inGameAttrFilter(agentInGameAttribute) {
		return &models2.DamageCalcResult{}, nil
	}
	//data1, _ = json.MarshalIndent(agentInGameAttribute, "", "  ")

	//fmt.Printf("in: %s\n", string(data1))
	damage := models2.Damage{
		Attack:                 agentInGameAttribute.Attack + agentInGameAttribute.AttackBonus*agentOutGameAttribute.Attack,
		DamageMultiplier:       param.TestData.DamageMultiplier,
		CriticalDamage:         agentInGameAttribute.CriticalDamage,
		CriticalRate:           agentInGameAttribute.CriticalRate,
		DamageResistance:       agentInGameAttribute.DamageResistance,
		CommonDamageBonus:      agentInGameAttribute.CommonDamageBonus,
		CommonSheerDamageBonus: agentInGameAttribute.CommonSheerDamageBonus,
		//DefenseReduction:     1 - agentInGameAttribute.DefenseReduction,
		StunDamageMultiplier: 1.5 + agentInGameAttribute.StunDamageMultiplier,
		HP:                   agentInGameAttribute.HP + agentInGameAttribute.HPBonus*agentOutGameAttribute.HP,
	}
	// 贯穿力转模
	if param.AgentInfo.Features.Attribute2Sheer != nil {
		damage.SheerForce = param.AgentInfo.Features.Attribute2Sheer(models2.AgentAttribute{HP: damage.HP, Attack: damage.Attack})
	}

	//减防计算
	if param.AgentInfo.Features.LifeDestroy {
		damage.DefenseReduction = 1
	} else {
		damage.DefenseReduction = DefenseCalc(models2.DefenseParam{
			LevelBase:          param.TestData.LevelBase,
			MonsterBaseDefense: param.TestData.MonsterBaseDefense,
			DefenseReduction:   agentInGameAttribute.DefenseReduction,
			PenetrationRadio:   agentInGameAttribute.PenetrationRadio,
			Penetration:        agentInGameAttribute.Penetration,
		})
	}

	// 属性增伤计算
	damageBonus := 0.0
	sheerDamageBonus := 0.0
	switch param.Attribute {
	case models2.AttrIce:
		damageBonus = agentInGameAttribute.IceDamageBonus
		sheerDamageBonus = agentInGameAttribute.IceSheerDamageBonus
	case models2.AttrElectric:
		damageBonus = agentInGameAttribute.ElectricDamageBonus
		sheerDamageBonus = agentInGameAttribute.ElectricSheerDamageBonus
	case models2.AttrPhysical:
		damageBonus = agentInGameAttribute.PhysicalDamageBonus
		sheerDamageBonus = agentInGameAttribute.PhysicalSheerDamageBonus
	case models2.AttrFire:
		damageBonus = agentInGameAttribute.FireDamageBonus
		sheerDamageBonus = agentInGameAttribute.FireSheerDamageBonus
	case models2.AttrEther:
		damageBonus = agentInGameAttribute.EtherDamageBonus
		sheerDamageBonus = agentInGameAttribute.EtherSheerDamageBonus
	}
	damage.CommonDamageBonus += damageBonus
	damage.CommonSheerDamageBonus += sheerDamageBonus
	//damage.CriticalRate = 1
	// 输出计算
	// 爆伤 = (1 + CriticalDamage) * CriticalRate + 1 * (1-CriticalRate) =  CriticalDamage * CriticalRate + 1
	var output float64
	if param.AgentInfo.Features.LifeDestroy {
		output = damage.SheerForce *
			damage.DamageMultiplier *
			(damage.CriticalDamage*damage.CriticalRate + 1) *
			(1 - damage.DamageResistance) *
			(1 + damage.CommonDamageBonus) *
			(1 + damage.CommonSheerDamageBonus) *
			damage.DefenseReduction
	} else {
		output = damage.Attack *
			damage.DamageMultiplier *
			(damage.CriticalDamage*damage.CriticalRate + 1) *
			(1 - damage.DamageResistance) *
			(1 + damage.CommonDamageBonus) *
			damage.DefenseReduction
	}

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
		if param.AgentInfo.Features.LifeDestroy {

			attackImprove := AttackBonusImprove(baseAttribute.Attack, damage.Attack, 1)
			criticalRateImprove, criticalDamageImprove := CriticalImprove(models2.CriticalParam{
				CriticalDamage: damage.CriticalDamage,
				CriticalRate:   damage.CriticalRate,
			}, 1)
			penetrationImprove := DefensePenetrationImprove(models2.DefenseParam{
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
		} else {
			attackImprove := AttackBonusImprove(baseAttribute.Attack, damage.Attack, 1)
			criticalRateImprove, criticalDamageImprove := CriticalImprove(models2.CriticalParam{
				CriticalDamage: damage.CriticalDamage,
				CriticalRate:   damage.CriticalRate,
			}, 1)
			penetrationImprove := DefensePenetrationImprove(models2.DefenseParam{
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

	}

	return &models2.DamageCalcResult{
		Output:      output,
		OutGameAttr: agentOutGameAttribute,
		InGameAttr:  agentInGameAttribute,
		BaseAttr:    baseAttribute,
		Damage:      damage,
	}, nil
}

func DamageExpect(improve float64, damageBonus float64) (outDamageBonus float64) {
	outDamageBonus = damageBonus/improve - 1
	return
}
