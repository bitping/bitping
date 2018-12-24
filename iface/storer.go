package iface

import "github.com/codegangsta/cli"

// Storage interface is for each storage solution
type Storer interface {
	Name() string
	AddCLIFlags([]cli.Flag) []cli.Flag
	Configure(c *cli.Context) error
	IsConfigured(c *cli.Context) bool

	Push(interface{}) bool
}