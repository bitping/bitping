package bitping

import (
	"log"

	ethclient "github.com/ethereum/go-ethereum"
)

type EthereumOptions struct {
	IpcPath string
}

func NewClient(opts EthereumOptions) *ethclient.Client {
	ipcPath := opts.ipcPath
	client, err := ethclient.Client.Dial(ipcPath)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
