package cmd

import (
	"fmt"

	b "github.com/auser/bitping/blockchains"
	"github.com/codegangsta/cli"
)

var EthCmd = cli.Command{
	Name:  "eth",
	Usage: "Run ethereum server",
	Subcommands: []cli.Command{
		{
			Name:   "run",
			Usage:  "Run geth",
			Action: StartGeth,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ipcPath",
					Value: "/root/.ethereum/geth.ipc",
					Usage: "ipcPath for geth",
				},
			},
		},
	},
}

func StartGeth(c *cli.Context) {
	fmt.Printf("Starting geth...\n")
	opts := b.EthereumOptions{
		IpcPath: c.String("ipcPath"),
	}
	client := b.NewClient(opts)
	client.Run()
}
