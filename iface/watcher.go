package iface

import (
	"github.com/auser/bitping/types"
)

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
