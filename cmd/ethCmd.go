package cmd

import (
	"fmt"
	"log"

	b "github.com/auser/bitping/blockchains"
	"github.com/auser/bitping/types"
	"github.com/codegangsta/cli"
)

var sharedFlags = append([]cli.Flag{}, cli.StringFlag{
	Name:   "ipcPath",
	Value:  "/root/.ethereum/geth.ipc",
	Usage:  "ipcPath for geth",
	EnvVar: "IPC_PATH",
},
)

var EthCmd = cli.Command{
	Name:  "eth",
	Usage: "Run ethereum server",
	Flags: append([]cli.Flag{}, sharedFlags...),
	Subcommands: []cli.Command{
		{
			Name:   "run",
			Usage:  "Run geth",
			Action: StartListening,
			Flags: append([]cli.Flag{
				cli.BoolFlag{
					Name:  "stdout",
					Usage: "print to stdout",
				}}, sharedFlags...),
		},
		{
			Name:   "info",
			Usage:  "Get information",
			Action: GetInfo,
			Flags:  append([]cli.Flag{}, sharedFlags...),
		},
	},
}

//
// TODO: Break this out and make configurable
// with config file
func StartListening(c *cli.Context) {
	client := makeClient(c)

	var blockCh = make(chan *types.Block, 16)
	var transactionCh = make(chan *[]types.Transaction, 16)
	var errCh = make(chan error, 16)

	go client.Run(blockCh, transactionCh, errCh)

	for {
		select {
		case err := <-errCh:
			fmt.Printf("Error listening: %#v\n", err)
		case block := <-blockCh:
			fmt.Printf("Got a block: %#v\n", block)
		case txs := <-transactionCh:
			fmt.Printf("\nGot some transactions: %#v\n", txs)
		}
	}
}

func GetInfo(c *cli.Context) {
	client := makeClient(c)
	networkId := client.GetNetwork()
	fmt.Printf("Running on network: %v\n", networkId)
}

func makeClient(c *cli.Context) *b.EthereumApp {
	opts := b.EthereumOptions{
		IpcPath: c.String("ipcPath"),
	}
	client, err := b.NewClient(opts)

	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	return client
}
