package handlers

import (
	"os"
	"fmt"
	"context"
	"net/http"
	"path/filepath"


	"github.com/urfave/cli/v3"
	"github.com/gweppi/gwup/cmd/client/config"
)


func HandleUpload(ctx context.Context, cmd *cli.Command) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	// check if server is configured
	if config.IsUndefined() {
		return fmt.Errorf("You have not set up a server, please do so by running the config command")
	}

	// get the name of the file that has to be uploaded
	fileName := cmd.Args().First()

	// check if the file was actually provided
	if fileName == "" {
		return fmt.Errorf("There was no file provided to upload")
	}

	// check if the file does actually exist
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("The provided file does not exist")
		} else {
			return fmt.Errorf("Something went wrong locating the file")
		}
	}

	// now that we know that the file exits we can stream it to the server
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", config.ServerUrl + "/upload", file)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/octet-stream")
	request.Header.Add("Content-Disposition", "attachment; filename=\"" + filepath.Base(fileName) + "\"")
	if config.AuthCode != "" {
		request.Header.Add("Authorization", config.AuthCode)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Println("File uploaded with code", response.StatusCode)

	return nil
}
