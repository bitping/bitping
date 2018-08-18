package bitping

import (
	"fmt"

	types "github.com/auser/bitping/types"
)

type EosOptions struct {
	Node string
}

// TODO: make interface for blockchains
type EosApp struct {
	Client  interface{}
	Options EosOptions
}

func NewEosClient(opts EosOptions) (*EosApp, error) {
	app := &EosApp{
		Options: opts,
	}

	return app, nil
}

func (app *EthereumApp) RunEos(
	blockChan chan types.Block,
	// transChan chan []types.Transaction,
	errChan chan error,
) {
	fmt.Printf("Run EOS...\n")
}
