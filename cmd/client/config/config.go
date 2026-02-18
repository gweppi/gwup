package config


import (
	"os"
	"fmt"
	"path/filepath"
	"encoding/json"

	"github.com/gweppi/gwup/internal/shared"
)

func GetConfig() (shared.Config, error) {
	// check if there is a file stored by the user that includes config values
	file, err := getConfigFile()
	if err != nil {
		return shared.Config{}, err
	}

	var config shared.Config

	stat, err := file.Stat()
	if err != nil {
		return shared.Config{}, err
	}

	bytes := make([]byte, stat.Size())
	file.Read(bytes)
	json.Unmarshal(bytes, &config)
	
	file.Close()
	return config, nil
}

func SetConfig(config shared.Config) error {
	file, err := getConfigFile()
	if err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.Encode(config)

	file.Close()

	return nil
}

func getConfigFile() (*os.File, error) {
	dir, err := os.UserHomeDir()
	// there was en error getting the config file, print error on screen
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	// path where the config file should be located
	configDir := filepath.Join(dir, ".config", "gwup")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("could not create config directory: %w", err)
	}
	
	configFilePath := filepath.Join(configDir, "config.json")
	configFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0600)
	
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	return configFile, nil
}

