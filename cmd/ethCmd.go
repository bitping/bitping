package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	b "github.com/auser/bitping/blockchains"
	"github.com/auser/bitping/types"
	"github.com/codegangsta/cli"
	bolt "github.com/coreos/bbolt"
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

// TODO: Break this out and make configurable
// with config file

func blockHandler(db *bolt.DB, in <-chan types.Block, errCh chan error) {
	logger := func(block types.Block) {
		fmt.Printf("Got a block in the block handler: %#v\n", block)
	}

	writeToDb := func(block types.Block) {
		err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("blocks"))
			if err != nil {
				return err
			}

			id := make([]byte, 8)
			binary.BigEndian.PutUint64(id, uint64(block.BlockNumber))

			buf, err := json.Marshal(block)
			if err != nil {
				return err

			}
			if err := b.Put(id, buf); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			errCh <- err
		}
	}

	for block := range in {
		logger(block)
		writeToDb(block)
	}
}

func StartListening(c *cli.Context) {
	client := makeClient(c)

	var blockCh = make(chan types.Block, 16)
	var transactionCh = make(chan []types.Transaction, 16)
	var errCh = make(chan error, 16)

	// open db
	db, err := bolt.Open("database/eth.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	go client.Run(blockCh, transactionCh, errCh)
	go blockHandler(db, blockCh, errCh)

	for {
		select {

		case err := <-errCh:
			fmt.Printf("Error listening: %#v\n", err)
		//case block := <-blockCh:
		//fmt.Printf("Got a block: %#v\n", block)
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
