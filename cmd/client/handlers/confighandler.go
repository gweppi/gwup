package handlers

import (
	"fmt"
	"context"
	"net/url"
	"net/http"
	"encoding/json"
	
	"github.com/urfave/cli/v3"
	"github.com/gweppi/gwup/internal/shared"
	"github.com/gweppi/gwup/cmd/client/config"
)

func HandleConfig(ctx context.Context, cmd *cli.Command) error {
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
	} else {
		fmt.Println("This server does not require authentication to paste, continuing...")
	}

	if err := config.SetConfig(newConfig); err != nil {
		return err
	}
	// if server requires authcode ask for it to be provided, for now just print it
	return nil
}

