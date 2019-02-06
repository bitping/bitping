package iface

import (
	"github.com/codegangsta/cli"
)

// Configurable is a cli interface for adding and configuring using cli flags
type Configurable interface {
	// Takes a slice of cli.Flag and returns a slice of cli.Flag with new flags
	AddCLIFlags([]cli.Flag) []cli.Flag

	// Checks to see if the cli.Context has enough cli.Flags for configuring
	// this object
	CanConfigure(c *cli.Context) bool

	// Configure this object using the cli.Context
	Configure(c *cli.Context) error
}
