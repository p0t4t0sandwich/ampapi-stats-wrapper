package stats

import (
	"encoding/json"
	"fmt"
	"log"
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
	if settings.ADDRESS == "" {
		if settings.USE_UDS {
			settings.ADDRESS = "/tmp/stats.sock"
		} else {
			settings.ADDRESS = "0.0.0.0:3021"
		}
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
	log.Printf("Creating new settings object with the following properties:")
	log.Printf("ADDRESS: %s", settings.ADDRESS)
	log.Printf("USE_UDS: %t", settings.USE_UDS)
	log.Printf("AMP_API_URL: %s", settings.AMP_API_URL)
	log.Printf("AMP_API_USERNAME: %s", settings.AMP_API_USERNAME)
	if len(settings.AMP_API_PASSWORD) > 6 {
		log.Printf("AMP_API_PASSWORD: %s", settings.AMP_API_PASSWORD[0:3]+"***"+settings.AMP_API_PASSWORD[len(settings.AMP_API_PASSWORD)-3:])
	} else {
		log.Printf("AMP_API_PASSWORD: %s", settings.AMP_API_PASSWORD)
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
