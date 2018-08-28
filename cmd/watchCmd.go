package cmd

import (
	"encoding/json"
	"log"

	"github.com/auser/bitping/blockchains"
	"github.com/auser/bitping/types"
	"github.com/auser/bitping/work"
	"github.com/codegangsta/cli"
	"github.com/thedevsaddam/gojsonq"

	. "github.com/auser/bitping/iface"
)

var watchables []Watchable
var WatchCmd cli.Command

func init() {
	watchables = []Watchable{
		&blockchains.EthereumApp{},
		&blockchains.EosApp{},
	}

	fs := []cli.Flag{}

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
					// jq := gojsonq.New().JSONString(jsonString).From("transactions").Where("gas", ">", 100000)
					jq := gojsonq.New().JSONString(jsonString).From("singletonTransactions").Where("value", ">", 0)
					log.Printf("JQ %#v\n", jq.Get())

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
