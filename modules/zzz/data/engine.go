package data

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"zzz_helper/internal/config"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/modules/zzz/models"
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

func GetEngineInfos() (SEngineInfo, error) {
	b, err := file2.ReadFileBytes(filepath.Join(config.CurrentPath, "conf/engines.yml"))
	if err != nil {
		return nil, err
	}
	var infos SEngineInfo

	err = yaml.Unmarshal(b, &infos)
	return infos, err
}
