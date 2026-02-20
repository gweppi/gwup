package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"io"
	"path/filepath"

	"github.com/gweppi/gwup/cmd/client/config"
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
		Action: handlePaste,
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
			{
				Name: "upload",
				Usage: "Upload a file or directory to the server",
				Action: handleUpload,
			},
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		fmt.Println(err)
	}
}

func handlePaste(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("This is the paste command")
	return nil
}

func handleUpload(ctx context.Context, cmd *cli.Command) error {
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

	readPipe, writePipe := io.Pipe()
	writer := multipart.NewWriter(writePipe)

	go func() {
		defer file.Close()
		defer writePipe.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile(shared.FileName, filepath.Base(fileName))
		if err != nil {
			writePipe.CloseWithError(err)
		}

		if _, err := io.Copy(part, file); err != nil {
			writePipe.CloseWithError(err)
		}
	}()
	
	request, err := http.NewRequest("POST", config.ServerUrl + "/upload", readPipe)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
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
	} else {
		fmt.Println("This server does not require authentication to paste, continuing...")
	}

	if err := config.SetConfig(newConfig); err != nil {
		return err
	}
	// if server requires authcode ask for it to be provided, for now just print it
	return nil
}


