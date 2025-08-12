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
func GetAgentInfos() (SAgentInfo, error) {
	b, err := file2.ReadFileBytes(filepath.Join(config.CurrentPath, "conf/agents.yml"))
	if err != nil {
		return nil, err
	}
	var infos SAgentInfo

	err = yaml.Unmarshal(b, &infos)
	return infos, err
}
