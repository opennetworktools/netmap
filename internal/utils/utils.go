package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveStructAsJson(filename string, structData any, timestamp string) error {
	jsonRes, err := json.Marshal(structData)
	if err != nil {
		return err
	}

	appDir, err := CreateDirectoryToSaveOutput()
	if err != nil {
		return err
	}

	jsonFilename := fmt.Sprintf("%s_%s.json", "graph", timestamp)
	jsonPilePath := filepath.Join(appDir, jsonFilename)
	err = os.WriteFile(jsonPilePath, jsonRes, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	return nil
}

func CreateDirectoryToSaveOutput() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(configDir, "netmap")
	os.MkdirAll(appDir, 0755) // Ensure the directory exists

	return appDir, nil
}

func TruncateString(s string, n int) string {
	if n >= len(s) {
		return s
	}
	return s[:n] + "..."
}
