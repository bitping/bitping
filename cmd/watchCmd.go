package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	contracts "github.com/auser/bitping/contracts"

	"github.com/ethereum/go-ethereum/common"

	"github.com/auser/bitping/blockchains"
	"github.com/auser/bitping/types"
	"github.com/auser/bitping/work"
	"github.com/codegangsta/cli"
	"github.com/thedevsaddam/gojsonq"

	. "github.com/auser/bitping/iface"
)

var watchables []Watchable
var WatchCmd cli.Command
var contract *contracts.USDToken

func init() {
	watchables = []Watchable{
		&blockchains.EthereumApp{},
		&blockchains.EosApp{},
	}

	fs := append([]cli.Flag{},
		cli.StringFlag{
			Name: "contractAddress",
		},
	)

	for _, w := range watchables {
		fs = w.AddCLIFlags(fs)
	}

	WatchCmd = cli.Command{
		Name:   "watch",
		Usage:  "Run watch server",
		Flags:  fs,
		Action: StartListening,
	}
}

func StartListening(c *cli.Context) {
	// Open the contract
	contractAddrStr := c.String("contractAddress")
	ethClient := watchables[0].(*blockchains.EthereumApp).Client

	if contractAddrStr != "" {
		contractAddr := common.HexToAddress(contractAddrStr)
		contract, err := contracts.NewUSDToken(contractAddr, ethClient)

		if err != nil {
			panic(err)
		}

		fmt.Println("contract is loaded")
		_ = contract
	}

	var workerPool = work.New(128)
	defer workerPool.Stop()
	var blockCh = make(chan types.Block)
	var errCh = make(chan error)

	// CONFIGURE WATCHABLES
	for _, w := range watchables {
		if w.IsConfigured(c) {
			log.Printf("Configuring %v", w.Name())
		} else {
			log.Printf("Not Starting %v", w.Name())
			continue
		}

		if err := w.Configure(c); err != nil {
			log.Printf("Failed to Start %v: %v", w.Name(), err)
			continue
		}

		log.Printf("Starting %v", w.Name())

		go w.Watch(blockCh, errCh)
	}

	// SETUP LISTENING PROCESS
	for {
		select {
		case block := <-blockCh:
			// Handle querying here
			workerPool.Submit(func() {
				dat, err := json.Marshal(block)
				if err != nil {
					log.Printf("Error marshaling block: %s\n", err.Error())
					return
				} else {
					// fmt.Printf("Running submitted block to worker pool %s\n", dat)
					jsonString := string(dat[:])
					// fmt.Printf("%s\n", jsonString)

					// SELECT * FROM ethereum transactions WHERE address = "0xdeadbeef" AND gas > 1000000 confirmed;
					// SELECT * transactions WHERE to = "0xcoffeeshop" WHERE UTXO is complete; -> transactions*n
					jq := gojsonq.New().JSONString(jsonString).From("transactions").Where("from", "=", "0xf8f59f0269c4f6d7b5c5ab98d70180eaa0c7507e").OrWhere("to", "=", "0xf8f59f0269c4f6d7b5c5ab98d70180eaa0c7507e")
					// jq := gojsonq.New().JSONString(jsonString).From("singletonTransactions").Where("value", ">", 0)
					// log.Printf("%#v\n", jsonString)
					if jq.Count() > 0 {
						fmt.Printf("%s\n", jsonString)
						log.Printf("An event occurred on the address")
					}

					// fire event
					// for i, matched := range(jq.Get()) {
					// 	matchedEndpoint := db.GetEndpoint(i)
					// 	// build JSON
					// 	httpClient.POST(matchedEndpoint)
					// }

				}
			})
		case err := <-errCh:
			log.Printf("Error listening: %#v\n", err)
			break
		}
	}
}
