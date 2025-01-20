package data

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"zzz_helper/internal/models"
	"zzz_helper/internal/res"
)

var (
	EngineInfos SEngineInfo
)

type SEngineInfo []models.WeaponEngine

func (this SEngineInfo) GetInfo(name string) (info models.WeaponEngine, err error) {
	for _, agentInfo := range this {
		if agentInfo.Name == name {
			info = agentInfo
			return
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
