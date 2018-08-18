package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"

	b "github.com/auser/bitping/blockchains"
	"github.com/auser/bitping/types"
	"github.com/auser/bitping/work"
	"github.com/codegangsta/cli"
	bolt "github.com/coreos/bbolt"
)

var sharedFlags = append([]cli.Flag{}, cli.StringFlag{
	Name:   "node",
	Value:  "http://localhost:8545",
	Usage:  "node path for geth",
	EnvVar: "NODE_PATH",
},
	cli.StringFlag{
		Name:   "googleProjectId",
		Usage:  "google project id",
		EnvVar: "GOOGLE_PROJECT_ID",
	}, cli.StringFlag{
		Name:   "googleKeyfile",
		Usage:  "google auth client",
		EnvVar: "GOOGLE_AUTH_FILE",
	})

var (
	PubsubClient *pubsub.Client
	Database     *bolt.DB
)

const PubsubTopicID = "ethereum-block-update"

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
				}, cli.BoolFlag{
					Name:  "background",
					Usage: "run in the background",
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

func blockHandler(in <-chan types.Block, errCh chan error) {
	logger := func(block types.Block) {
		fmt.Printf("Got a block in the block handler: %#v\n", block)
	}

	for block := range in {
		logger(block)
		// 		writeToDb(block, errCh)
	}
}

func writeToDb(block types.Block, errCh chan error) {
	ctx := context.Background()

	// err := Database.Update(func(tx *bolt.Tx) error {
	// 	b, err := tx.CreateBucketIfNotExists([]byte("blocks"))
	// 	if err != nil {
	// 		return err
	// 	}

	id := make([]byte, 8)
	binary.BigEndian.PutUint64(id, uint64(block.BlockNumber))

	// here
	buf, err := json.Marshal(block)
	if err != nil {
		errCh <- err
	}
	// 	if err := b.Put(id, buf); err != nil {
	// 		return err
	// 	}

	// 	if PubsubClient == nil {
	// 		return nil
	// 	}

	// 	return nil
	// })
	topic := PubsubClient.Topic(PubsubTopicID)
	_, err = topic.Publish(ctx, &pubsub.Message{Data: []byte(buf)}).Get(ctx)

	if err != nil {
		fmt.Printf("Error publishing: %s\n", err)
		errCh <- err
	}
}

// ----------------------------------------------------------------------
// TODO: Abstract this part
// ----------------------------------------------------------------------

// ----------------------------------------------------------------------
// TODO: end abstraction
// ----------------------------------------------------------------------
func StartListening(c *cli.Context) {
	var workerPool = work.New(128)
	var in = make(chan types.Block)
	var errCh = make(chan error)
	// var done = make(chan struct{})

	// SETUP LISTENING PROCESS
	client := makeClient(c)
	go client.Run(in, errCh)
	for {
		select {
		case block := <-in:
			workerPool.Submit(func() {

				b, err := json.Marshal(block)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Printf("%s\n", b)
				}

				// 	// TODO: Use alice to make this a workflow
				// 	txs, err := client.GetTransactionsFromBlock(o)
				// 	if err != nil {
				// 		log.Printf("Error: %s\n", err.Error())
				// 		return
				// 	}
				// fmt.Printf("%d) txs: %#v\n", block.BlockNumber, len(txs))
				// 	// TODO: push to google pub/sub
			})

		case err := <-errCh:
			fmt.Printf("Error listening: %#v\n", err)
			// close(in)
			// <-done
		}
	}
	// END SETUP

	fmt.Printf("Shutdown network\n")
	workerPool.Stop()
	// if c.Bool("background") == true {
	// go startListeningForeground(c)
	// } else {
	// startListeningForeground(c)
	// }
}

func startListeningForeground(c *cli.Context) {
	client := makeClient(c)

	var blockCh = make(chan types.Block, 16)
	// var transactionCh = make(chan []types.Transaction, 16)
	var errCh = make(chan error, 16)

	// open db
	go client.Run(blockCh, errCh)
	go blockHandler(blockCh, errCh)

	for {
		select {

		case err := <-errCh:
			fmt.Printf("Error listening: %#v\n", err)
			//case block := <-blockCh:
			//fmt.Printf("Got a block: %#v\n", block)
			// case txs := <-transactionCh:
			// fmt.Printf("\nGot some transactions: %#v\n", txs)
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
		Node: c.String("node"),
	}
	client, err := b.NewClient(opts)

	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// pubsubClient, err := configurePubsub(c)
	// if err != nil {
	// 	log.Fatalf("error: %v\n", err)
	// }

	// PubsubClient = pubsubClient

	configureDB(c)
	return client

}

func configureDB(c *cli.Context) {
	db, err := bolt.Open("database/eth.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	Database = db
}

func configurePubsub(c *cli.Context) (*pubsub.Client, error) {
	googleProjectId := c.String("googleProjectId")

	ctx := context.Background()

	// clientOption := option.WithServiceAccountFile(c.String("googleKeyfile"))
	client, err := pubsub.NewClient(ctx, googleProjectId)
	if err != nil {
		return nil, err
	}

	// Create the topic if it doesn't exist.
	if exists, err := client.Topic(PubsubTopicID).Exists(ctx); err != nil {
		return nil, err
	} else if !exists {
		if _, err := client.CreateTopic(ctx, PubsubTopicID); err != nil {
			return nil, err
		}
	}
	return client, nil
}
