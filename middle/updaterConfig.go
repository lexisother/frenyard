package middle

import (
	"os"
	"fmt"
	"path/filepath"
	"encoding/json"
)

// UpdaterConfig is the configuration structure for the application.
type UpdaterConfig struct {
	GamePath string `json:"gamePath"`
	DevMode bool `json:"devMode"`
}

func getUpdaterConfigPath() string {
	cfg, err := os.UserConfigDir()
	if err != nil {
		cfg = ""
	}
	return filepath.Join(cfg, "ccmodupdater.json")
}

// ReadUpdaterConfig returns the current configuration for the application.
func ReadUpdaterConfig() UpdaterConfig {
	cfg := &UpdaterConfig{}
	cfgFile, err := os.Open(getUpdaterConfigPath())
	if err != nil {
		fmt.Printf("Failed to open updater config: %s\n", err.Error())
		return *cfg
	}
	defer cfgFile.Close()
	if json.NewDecoder(cfgFile).Decode(cfg) != nil {
		fmt.Printf("Failed to decode updater config\n")
		return *cfg
	}
	WriteUpdaterConfig(*cfg)
	return *cfg
}

// WriteUpdaterConfig saves a configuration as the current configuration.
func WriteUpdaterConfig(cfg UpdaterConfig) {
	cfgPath := getUpdaterConfigPath()
	// success/failure doesn't matter if the OpenFile works
	os.MkdirAll(filepath.Dir(cfgPath), os.ModePerm)
	cfgFile, err := os.OpenFile(cfgPath, os.O_WRONLY | os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to save updater config: %s\n", err.Error())
		return
	}
	defer cfgFile.Close()
	if json.NewEncoder(cfgFile).Encode(cfg) != nil {
		fmt.Printf("Failed to encode updater config\n")
	}
}
