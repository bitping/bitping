package main

import (
	"fmt"
	"os"

	b "github.com/auser/bitping/blockchains"
	"github.com/codegangsta/cli"
)

func main() {
	Run(os.Args)
}

func Run(args []string) {
	app := cli.NewApp()
	app.Name = "bitping"
	app.Commands = []cli.Command{
		{
			Name:   "eth",
			Usage:  "Run ethereum server",
			Action: RunServer,
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
		},
	}
	app.Run(args)
}

// TODO: move to other file
func StartGeth(c *cli.Context) {
	fmt.Printf("Starting geth...\n")
	opts := &b.EthereumOptions{
		IpcPath: c.String("ipcPath"),
	}
	_ := b.NewClient(opts)
}

func RunServer(c *cli.Context) {
	app := NewApp(AppOptions{})
	app.Run()
}
