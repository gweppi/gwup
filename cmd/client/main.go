package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gweppi/gwup/internal/shared"
	"github.com/gweppi/gwup/cmd/client/config"
	"github.com/urfave/cli/v3"
)

type contextKey string
const configKey contextKey = "config"

// Command should include the following subcommands:
// - join (join a copy paste server)
// - paste (paste a file or dir to a server)
func main() {
	ctx := context.Background()

	config, err := config.GetConfig()
	if err != nil {
		// config could not be loaded, return error
		fmt.Println(err)
		return
	}
	ctx = context.WithValue(ctx, configKey, config)

	cmd := &cli.Command {
		Name: "gwup",
		Version: "1.0.0",
		Action: handleUpload,
		Commands: []*cli.Command {
			{
				Name: "config",
				Usage: "Configurate gwup (set sever and authcode)",
				Action: handleConfig,
			},
			{
				Name: "paste",
				Usage: "Paste a file or directory",
				Action: handlePaste,
			},
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		fmt.Println(err)
	}
}

func handleUpload(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("This is the paste command")
	return nil
}


func handleConfig(ctx context.Context, cmd *cli.Command) error {
	newConfig := shared.Config{}

	var serverUrl string;
	fmt.Print("Please enter your server url...\n > ")
	fmt.Scan(&serverUrl)
	// do some checks on the server url
	if _, err := url.Parse(serverUrl); err != nil {
		return err
	}
	newConfig.ServerUrl = serverUrl
	
	res, err := http.Get(serverUrl + "/health")
	if err != nil {
		return err
	}

	defer res.Body.Close()
	var status shared.ServerInfo
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&status); err != nil {
		return err
	}

	// check if server status is ok
	if status.Status != "ok" {
		return fmt.Errorf("The server is not healthy")
	}

	// check if server requires authcode
	if status.RequiresAuth {
	}


	if err := config.SetConfig(newConfig); err != nil {
		return err
	}
	// if server requires authcode ask for it to be provided, for now just print it
	return nil
}

func handlePaste(ctx context.Context, cmd *cli.Command) error {
	return nil
}

