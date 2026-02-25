package handlers

import (
	"os"
	"io"
	"fmt"
	"mime"
	"net/http"
	"context"

	"github.com/gweppi/gwup/cmd/client/config"
	"github.com/urfave/cli/v3"
)


func HandlePaste(ctx context.Context, cmd *cli.Command) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("GET", config.ServerUrl + "/download", nil)
	if err != nil {
		return err
	}

	request.Header.Add("X-File-Id", cmd.Args().Get(1))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()


	_, params, err := mime.ParseMediaType(response.Header.Get("Content-Disposition"))
	if err != nil {
		return err
	}

	fileName := params["filename"]
	if fileName == "" {
		return fmt.Errorf("No file name provided in Content-Disposition header")
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, response.Body)

	fmt.Printf("Written file named %s to disk\n", fileName)
	return nil
}
