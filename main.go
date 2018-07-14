package main

import (
	"os"

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
			Name:   "run",
			Usage:  "Run server",
			Action: RunServer,
		},
	}
	app.Run(args)
}

func RunServer(c *cli.Context) {
	app := NewApp(&AppOptions{})
	app.Run()
}
