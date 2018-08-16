package main

import (
	"os"

	cmd "github.com/auser/bitping/cmd"
	"github.com/codegangsta/cli"
)

func main() {
	Run(os.Args)
}

func Run(args []string) {
	app := cli.NewApp()
	app.Name = "bitping"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "daemonize",
			Usage: "run as a daemon",
		},
	}
	app.Commands = []cli.Command{cmd.EthCmd}
	app.Run(args)
}

// // TODO: move to other file
// func RunServer(c *cli.Context) {
// 	app := NewApp(AppOptions{})
// 	app.Run()
// }
