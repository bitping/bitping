package cmd

import (
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"

	b "github.com/auser/bitping/blockchains"
	"github.com/auser/bitping/types"
	"github.com/auser/bitping/work"
	"github.com/codegangsta/cli"
)

var sharedFlags = append([]cli.Flag{}, cli.StringFlag{
	Name:   "eth",
	Usage:  "ethereum address",
	EnvVar: "ETH_PATH",
},
	cli.StringFlag{
		Name:  "eos",
		Usage: "eos address",
	},
)

var (
	PubsubClient *pubsub.Client
)

const PubsubTopicID = "ethereum-block-update"

var WatchCmd = cli.Command{
	Name:   "watch",
	Usage:  "Run watch server",
	Flags:  append([]cli.Flag{}, sharedFlags...),
	Action: StartListening,
}

func StartListening(c *cli.Context) {

	var ethAddr = c.String("eth")
	var eosAddr = c.String("eos")

	var workerPool = work.New(128)
	var in = make(chan types.Block)
	var errCh = make(chan error)
	// var done = make(chan struct{})

	// SETUP LISTENING PROCESS
	if ethAddr != "" {
		runEthereum(ethAddr, in, errCh)
	}
	if eosAddr != "" {
		runEos(eosAddr, in, errCh)
	}
	for {
		select {
		case block := <-in:
			fmt.Printf("Got a block: %d\n", block.BlockNumber)
		case err := <-errCh:
			fmt.Printf("Error listening: %#v\n", err)
			break
		}
	}
	// END SETUP

	fmt.Printf("Shutdown network\n")
	workerPool.Stop()
}

func runEthereum(addr string, in chan types.Block, errCh chan error) {
	opts := b.EthereumOptions{
		Node: addr,
	}
	ethClient, err := b.NewEthClient(opts)
	if err != nil {
		log.Fatal("Error: %s\n", err.Error())
	}
	go ethClient.Run(in, errCh)
}

func runEos(addr string, in chan types.Block, errCh chan error) {
	fmt.Printf("Do this next")
}
