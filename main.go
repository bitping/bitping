package main

import (
	"io/ioutil"
	"os"

	cmd "github.com/auser/bitping/cmd"
	"github.com/codegangsta/cli"
)

func main() {
	Run(os.Args)
}

func Run(args []string) {
	app := cli.NewApp()
	var version = readVersion()
	app.Name = "bitping"
	app.Version = version
	app.Usage = "bitpingger"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "daemonize",
			Usage: "run as a app",
		},
	}
	app.Commands = []cli.Command{cmd.WatchCmd}
	app.Run(args)
}

func readVersion() string {
	var version = "0.0.0"
	da, err := ioutil.ReadFile("./Version")
	if err == nil {
		version = string(da)
	}
	return version
}
