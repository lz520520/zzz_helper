package models

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"zzz_helper/res"
)

type SDriversInfo []DriverSetAttribute

var (
	DriversInfos SDriversInfo
)

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

func init() {
	err := yaml.Unmarshal(res.Drivers, &DriversInfos)
	if err != nil {
		panic(err)
	}
}
