package core

import (
	"log"
	"os"
)

type Settings struct {
}

func ReadSettings(settingsPath string) Settings {
	settings := Settings{}

	settingsData, err := os.ReadFile(settingsPath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", settingsPath, err)
	}

	err = yaml.Unmarshal(settingsData, &settings)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", settingsPath, err)
	}

	return settings
}
