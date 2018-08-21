package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	b "github.com/auser/bitping/blockchains"
	"github.com/auser/bitping/types"
	"github.com/auser/bitping/work"
	"github.com/codegangsta/cli"
	"github.com/thedevsaddam/gojsonq"
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
	cli.StringFlag{
		Name:  "eos-p2p",
		Usage: "eos p2p address",
		Value: "peering.mainnet.eoscanada.com:9876",
	},
	cli.Int64Flag{
		Name:  "eos-version",
		Usage: "eos network version",
		Value: int64(1206),
	},
)

var WatchCmd = cli.Command{
	Name:   "watch",
	Usage:  "Run watch server",
	Flags:  append([]cli.Flag{}, sharedFlags...),
	Action: StartListening,
}

func StartListening(c *cli.Context) {

	var ethAddr = c.String("eth")
	var eosAddr = c.String("eos")
	var eosp2pAddr = c.String("eos-p2p")
	var eosNetworkVersion = c.Int64("eos-version")

	var workerPool = work.New(128)
	var in = make(chan types.Block)
	var errCh = make(chan error)
	// var done = make(chan struct{})

	// SETUP LISTENING PROCESS
	if ethAddr != "" {
		runEthereum(ethAddr, in, errCh)
	}
	if eosAddr != "" {
		runEos(eosAddr, eosp2pAddr, eosNetworkVersion, in, errCh)
	}
	for {
		select {
		case block := <-in:
			// Handle querying here
			workerPool.Submit(func() {
				dat, err := json.Marshal(block)
				if err != nil {
					fmt.Printf("Error marshaling block: %s\n", err.Error())
					return
				} else {
					// fmt.Printf("Running submitted block to worker pool %s\n", dat)
					jsonString := string(dat[:])

					// SELECT * FROM ethereum transactions WHERE address = "0xdeadbeef" AND gas > 1000000 confirmed;
					// SELECT * transactions WHERE to = "0xcoffeeshop" WHERE UTXO is complete; -> transactions*n
					jq := gojsonq.New().JSONString(jsonString).From("transactions").Where("gas", ">", 100000)
					fmt.Printf("%#v\n", jq.Get())

					// fire event
					// for i, matched := range(jq.Get()) {
					// 	matchedEndpoint := db.GetEndpoint(i)
					// 	// build JSON
					// 	httpClient.POST(matchedEndpoint)
					// }

				}
			})
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

func runEos(addr string, p2pAddr string, eosNetworkVersion int64, in chan types.Block, errCh chan error) {
	opts := b.EosOptions{
		P2PAddr:        p2pAddr,
		Node:           addr,
		NetworkVersion: eosNetworkVersion,
	}

	fmt.Printf("%#v\n", opts)

	eosClient, err := b.NewEosClient(opts)
	if err != nil {
		log.Fatal("Error loading Eos client: %s\n", err.Error())
	}
	go eosClient.Run(in, errCh)
}
