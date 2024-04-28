package stats

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/p0t4t0sandwich/ampapi-go"
	"github.com/p0t4t0sandwich/ampapi-go/modules"
)

// Settings - Settings struct
type Settings struct {
	ADDRESS          string `json:"ADDRESS"`
	USE_UDS          bool   `json:"USE_UDS"`
	AMP_API_URL      string `json:"AMP_API_URL"`
	AMP_API_USERNAME string `json:"AMP_API_USERNAME"`
	AMP_API_PASSWORD string `json:"AMP_API_PASSWORD"`
}

// NewSettings - Create a new settings instance
func NewSettings() *Settings {
	// Get settings from settings.json
	settingsFile, err := os.ReadFile("./settings.json")
	if err != nil {
		fmt.Println("Error reading settings.json")
	}
	var settings Settings
	_ = json.Unmarshal(settingsFile, &settings)

	// Override settings from env
	ENV_ADDRESS := os.Getenv("ADDRESS")
	if ENV_ADDRESS != "" {
		settings.ADDRESS = ENV_ADDRESS
	}
	ENV_USE_UDS := os.Getenv("USE_UDS") == "true"
	if ENV_USE_UDS {
		settings.USE_UDS = ENV_USE_UDS
	}
	ENV_AMP_API_URL := os.Getenv("AMP_API_URL")
	if ENV_AMP_API_URL != "" {
		settings.AMP_API_URL = ENV_AMP_API_URL
	}
	ENV_AMP_API_USERNAME := os.Getenv("AMP_API_USERNAME")
	if ENV_AMP_API_USERNAME != "" {
		settings.AMP_API_USERNAME = ENV_AMP_API_USERNAME
	}
	ENV_AMP_API_PASSWORD := os.Getenv("AMP_API_PASSWORD")
	if ENV_AMP_API_PASSWORD != "" {
		settings.AMP_API_PASSWORD = ENV_AMP_API_PASSWORD
	}
	return &settings
}

// TargetData - TargetData struct
type TargetData struct {
	TargetID ampapi.UUID
	Target   *modules.ADS
}

// InstanceData - InstanceData struct
type InstanceData struct {
	InstanceID ampapi.UUID
	Instance   *modules.CommonAPI
}
