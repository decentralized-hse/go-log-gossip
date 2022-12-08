package main

import (
	initCmd "github.com/decentralized-hse/go-log-gossip/cmd/init"
	"github.com/decentralized-hse/go-log-gossip/cmd/remove"
	"github.com/decentralized-hse/go-log-gossip/cmd/start"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "go-log-gossip",
		Usage: "Multilog system working via gossip",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "init config files",
				Action: func(_ *cli.Context) error {
					return initCmd.CommandInitConfig()
				},
			},
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "starts node",
				Action: func(_ *cli.Context) error {
					return start.CommandStartNode()
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "removes all configs, including folders, keys and etc",
				Action: func(_ *cli.Context) error {
					return remove.CommandRemoveConfig()
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
