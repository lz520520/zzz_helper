package data

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"zzz_helper/modules/zzz/models"
	"zzz_helper/res"
)

var (
	EngineInfos SEngineInfo
)

type SEngineInfo []models.WeaponEngineInfo

func (this SEngineInfo) GetInfo(name string, star int) (info models.WeaponEngine, err error) {
	for _, agentInfo := range this {
		if agentInfo.Name == name {
			for _, weaponStar := range agentInfo.Stars {
				if weaponStar.Level == star {
					return weaponStar, nil
				}
			}
		}
	}
	err = fmt.Errorf("engine info not found")
	return
}

func init() {
	err := yaml.Unmarshal(res.Engines, &EngineInfos)
	if err != nil {
		panic(err)
	}
}
