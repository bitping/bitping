package iface

import (
	"github.com/auser/bitping/types"
	"github.com/codegangsta/cli"
)

type Configurable interface {
	AddCLIFlags([]cli.Flag) []cli.Flag
	IsConfigured(c *cli.Context) bool
	Configure(c *cli.Context) error
}

type Watcher interface {
	Configurable

	Name() string
	Watch(blockCh chan types.Block, errCh chan error)
}
