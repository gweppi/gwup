package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"

	//"net/url"

	"github.com/gweppi/gwup/internal/shared"
	"github.com/urfave/cli/v3"
)

type contextKey string
const configKey contextKey = "config"

// Command should include the following subcommands:
// - join (join a copy paste server)
// - paste (paste a file or dir to a server)
func main() {
	ctx := context.Background()

	config, err := getConfig()
	if err != nil {
		// config could not be loaded, return error
		fmt.Println(err)
		return
	}
	ctx = context.WithValue(ctx, configKey, config)

	cmd := &cli.Command {
		Name: "gwup",
		Version: "1.0.0",
		Commands: []*cli.Command {
			{
				Name: "config",
				Usage: "Configurate gwup (set sever and authcode)",
				Action: configCommand,
			},
			{
				Name: "paste",
				Usage: "Paste a file or directory",
				Action: pasteCommand,
			},
		},
	}

	cmd.Run(ctx, os.Args)
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

func getConfig() (shared.Config, error) {
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
	
	return config, nil
}

func setConfig(config shared.Config) error {
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

	return nil
}

func configCommand(ctx context.Context, cmd *cli.Command) error {
	newConfig := shared.Config{}

	var serverUrl string;
	fmt.Print("Please enter your server url...\n > ")
	fmt.Scan(&serverUrl)
	// do some checks on the server url
	newConfig.ServerUrl = serverUrl


	if err := setConfig(newConfig); err != nil {
		return err
	}
	// if server requires authcode ask for it to be provided, for now just print it
	return nil
}

func pasteCommand(ctx context.Context, cmd *cli.Command) error {
	return nil
}

