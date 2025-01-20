package data

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
	"zzz_helper/internal/models"
)

func TestEngineGenearte(t *testing.T) {
	agent := models.WeaponEngine{
		Name: "深海访客",
		OutGame: models.AgentAttribute{
			Attack:            713,
			CriticalRate:      0.24,
			CommonDamageBonus: 0.25,
		},
		InGame: models.AgentAttribute{
			CriticalRate: 0.20,
		},
	}
	data, _ := yaml.Marshal([]models.WeaponEngine{agent})
	fmt.Println(string(data))
}
