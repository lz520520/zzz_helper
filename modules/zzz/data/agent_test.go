package data

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
	"zzz_helper/modules/zzz/models"
)

func TestGenerate(t *testing.T) {
	agent := models.AgentInfo{
		Name: "艾莲·乔",
		Stars: []models.AgentStar{
			{
				Level: 1,
				InGame: models.AgentAttribute{
					CommonDamageBonus: 0.25,
					CriticalRate:      0.12,
				},
			},
		},
		Attribute: models.AgentAttribute{
			HP:             7673,
			Attack:         863,
			Defense:        606,
			CriticalRate:   0.05,
			CriticalDamage: 0.5,
		},
		CorePassive: models.CommonAttribute{
			OutGame: models.AgentAttribute{
				Attack:       25 * 3,
				CriticalRate: 0.048 * 3,
			},
			InGame: models.AgentAttribute{
				CommonDamageBonus: 0.3,
			},
		},
	}
	data, _ := yaml.Marshal([]models.AgentInfo{agent})
	fmt.Println(string(data))

}
