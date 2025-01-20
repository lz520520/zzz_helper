package data

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"zzz_helper/internal/models"
	"zzz_helper/internal/res"
)

var (
	AgentInfos SAgentInfo
)

type SAgentInfo []models.AgentInfo

func (this SAgentInfo) GetInfo(name string) (info models.AgentInfo, err error) {
	for _, agentInfo := range this {
		if agentInfo.Name == name {
			info = agentInfo
			return
		}
	}
	err = fmt.Errorf("agent info not found")
	return
}

func init() {
	err := yaml.Unmarshal(res.Agents, &AgentInfos)
	if err != nil {
		panic(err)
	}
}
