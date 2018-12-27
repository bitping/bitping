package iface

import (
	"github.com/auser/bitping/types"
	"github.com/codegangsta/cli"
)

// Configurable is a cli interface for adding and configuring using cli flags
type Configurable interface {
	// Takes a slice of cli.Flag and returns a slice of cli.Flag with new flags
	AddCLIFlags([]cli.Flag) []cli.Flag

	// Checks to see if the cli.Context has enough cli.Flags for configuring
	// this object
	IsConfigured(c *cli.Context) bool

	// Configure this object using the cli.Context
	Configure(c *cli.Context) error
}

// Watchers watch a blockchain and pipes blockchain blocks and errors
type Watcher interface {
	// Has to be configured
	Configurable

	// Name returns the name of the watcher
	Name() string

	// Watch should start running the blockchian watcher process
	// It pipes back the unified block type or errors
	Watch(blockCh chan types.Block, errCh chan error)
}
