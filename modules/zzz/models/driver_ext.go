package models

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"zzz_helper/internal/config"
	"zzz_helper/internal/utils/file2"
)

type SDriversInfo []DriverSetAttribute

func (this SDriversInfo) GetInfo(name string) (info DriverSetAttribute, err error) {
	for _, agentInfo := range this {
		if agentInfo.Name == name {
			info = agentInfo
			return
		}
	}
	err = fmt.Errorf("driver info not found")
	return
}

func GetDriversInfos() (SDriversInfo, error) {
	b, err := file2.ReadFileBytes(filepath.Join(config.CurrentPath, "conf/driver.yml"))
	if err != nil {
		return nil, err
	}
	var infos SDriversInfo

	err = yaml.Unmarshal(b, &infos)
	return infos, err
}
