package main

import (
	"os"
	"fmt"
	"context"

	"github.com/urfave/cli/v3"
	"github.com/gweppi/gwup/cmd/client/config"
	"github.com/gweppi/gwup/cmd/client/handlers"
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
		Action: handlers.HandlePaste,
		Commands: []*cli.Command {
			{
				Name: "config",
				Usage: "Configurate gwup (set sever and authcode)",
				Action: handlers.HandleConfig,
			},
			{
				Name: "paste",
				Aliases: []string{"d"},
				Usage: "Paste a file or directory",
				Action: handlers.HandlePaste,
			},
			{
				Name: "upload",
				Aliases: []string{"u"},
				Usage: "Upload a file or directory to the server",
				Action: handlers.HandleUpload,
			},
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		fmt.Println(err)
	}
}
