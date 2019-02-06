package blockchains

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	types "github.com/auser/bitping/types"
	"github.com/codegangsta/cli"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	backoff "github.com/jpillora/backoff"
)

// Sometimes it takes a minute for the blockchain to catch up
// so we need to provide a small delay before we look up the block
var blockLookupDelay int = 1000
var b = &backoff.Backoff{
	Max: 5 * time.Minute,
}

type EthereumOptions struct {
	Node string
}

// TODO: make interface for blockchains
type EthereumApp struct {
	Client    *ethclient.Client
	Options   EthereumOptions
	NetworkId big.Int
}

func NewEthClient(opts EthereumOptions) (*EthereumApp, error) {
	nodePath := opts.Node
	client, err := ethclient.Dial(nodePath)
	if err != nil {
		return nil, err
	}

	app := &EthereumApp{
		Client:  client,
		Options: opts,
	}

	networkId := app.GetNetwork()
	fmt.Printf("Network id: %v\n", networkId)
	app.NetworkId = *networkId

	return app, nil
}

func (app EthereumApp) Name() string {
	return "Ethereum Watcher"
}

func (app EthereumApp) AddCLIFlags(fs []cli.Flag) []cli.Flag {
	return append(fs, cli.StringFlag{
		Name:   "eth",
		Usage:  "ethereum address",
		EnvVar: "ETH_PATH",
	})
}

func (app EthereumApp) IsConfigured(c *cli.Context) bool {
	return c.String("eth") != ""
}

func (app *EthereumApp) Configure(c *cli.Context) error {
	nodePath := c.String("eth")
	client, err := ethclient.Dial(nodePath)
	if err != nil {
		return err
	}

	app.Client = client
	app.Options = EthereumOptions{
		Node: nodePath,
	}

	networkId := app.GetNetwork()
	fmt.Printf("Network id: %v\n", networkId)
	app.NetworkId = *networkId

	return nil
}

func (app *EthereumApp) Watch(
	blockChan chan types.Block,
	// transChan chan []types.Transaction,
	errChan chan error,
) {
	fmt.Printf("Running Ethereum\n")

	// test
	var headsCh = make(chan *types.GethHeader)
	var errCh = make(chan error)
	go app.SubscribeToNews(headsCh, errCh)

	for {
		select {
		case err := <-errCh:
			fmt.Printf("Got an error in client.Run(): %v\n", err)
			// TODO: Reconnect here
			// go app.SubscribeToNews(headsCh, errCh)
		case head := <-headsCh:
			block, err := app.GetFromHeader(head)
			if err != nil {
				fmt.Printf("Error happened: %s\n", err.Error())
				errChan <- err
			} else {
				blockChan <- block
			}
			// transactions, err := app.makeTransactionsFrom(block)
			// if err != nil {
			// 	errChan <- err
			// } else {
			// 	transChan <- transactions
			// }
		}
	}
}

func (app *EthereumApp) GetNetwork() *big.Int {
	ctx := context.Background()
	networkId, err := app.Client.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return networkId
}

func (app *EthereumApp) SubscribeToNews(
	heads chan *types.GethHeader,
	errCh chan error,
) {
	var ch = make(chan *types.GethHeader)
	ctx := context.Background()
	sub, err := app.Client.SubscribeNewHead(ctx, ch)
	if err != nil {
		fmt.Printf("Error Subscribing to neew head: %s\n", err.Error())
		return
	}

	for {
		select {
		case err := <-sub.Err():
			fmt.Printf("Err: %v\n", err)
			errCh <- err
		case head := <-ch:
			heads <- head
		}
	}
}

func (app *EthereumApp) getByNumWithBackoff(num *big.Int) (*types.GethBlock, error) {
	ctx := context.Background()
	var (
		block *types.GethBlock
		err   error
	)
	var b = &backoff.Backoff{
		Max: 5 * time.Minute,
	}
	b.Reset()

	for {
		block, err = app.Client.BlockByNumber(ctx, num)

		if err != nil {
			d := b.Duration()
			if d > b.Max {
				return &types.GethBlock{}, err
			}
			time.Sleep(d)
			continue
		} else {
			break
		}
	}
	return block, nil
}

func (app *EthereumApp) getTransactionCountWithBackoff(hsh common.Hash) (uint, error) {
	var (
		count uint
		err   error
	)

	var b = &backoff.Backoff{
		Max: 5 * time.Minute,
	}
	// hsh := common.HexToHash(miner)
	ctx := context.Background()

	b.Reset()
	for {
		count, err = app.Client.TransactionCount(ctx, hsh)

		if count == 0 || err != nil {
			d := b.Duration()
			if d > b.Max {
				return 0, err
			}
			time.Sleep(d)
			continue
		} else {
			break
		}
	}
	return count, nil
}

func (app *EthereumApp) GetFromHeader(
	head *types.GethHeader,
) (types.Block, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	// head.Number = new(big.Int).SetInt64(3948621)

	// right now, this is blocking... do we want it to block?
	log.Printf("ETH Getting Block: %v", head.Number)
	block, err := app.getByNumWithBackoff(head.Number)
	if err != nil {
		return types.Block{}, err
	}

	// difficulty := types.BigNumber(block.Difficulty().String())
	// totalDifficulty := types.BigNumber(head.Difficulty.String())
	// cancel()

	var transactions []types.Transaction
	for i, tx := range block.Transactions() {
		var txFromStr string = "unknown"
		var txToStr string = "unknown"
		if msg, err := tx.AsMessage(types.GethHomesteadSigner{}); err != nil {

			txFromStr = msg.From().Hex()
			if msg.To() != nil {
				txToStr = msg.To().Hex()
			}
		}

		transaction := types.Transaction{
			BlockHash:   block.Hash().Hex(),
			BlockNumber: block.Number().Int64(),
			Hash:        tx.Hash().String(),
			Nonce:       int64(tx.Nonce()),

			EthereumTransaction: &types.EthereumTransaction{
				GasPrice:         types.NewBigInt(tx.Cost()),
				Gas:              tx.Gas(),
				TransactionIndex: int64(i),
			},

			Actions: []types.Action{
				types.Action{
					EthereumCall: &types.EthereumCall{
						Input: tx.Data(),
					},
					From:  strings.ToLower(txFromStr),
					To:    strings.ToLower(txToStr),
					Value: types.NewBigInt(tx.Value()),
					Data:  tx.Data(),
				},
			},
		}
		transactions = append(transactions, transaction)
	}

	blockObj := types.Block{
		Difficulty: types.NewBigInt(block.Difficulty()),
		Hash:       block.Hash().Hex(),
		HeaderHash: head.Hash().Hex(),
		Network:    "ethereum",
		NetworkID:  app.NetworkId.Int64(),
		Nonce:      fmt.Sprint(block.Nonce()),
		Number:     block.Number().Int64(),
		Size:       float64(block.Size()),
		ParentHash: block.ParentHash().String(),
		Time:       block.Time().Int64(),

		EthereumBlock: &types.EthereumBlock{
			TotalDifficulty:  types.NewBigInt(block.Difficulty()), // make sense?
			GasUsed:          block.GasUsed(),
			GasLimit:         block.GasLimit(),
			ExtraData:        fmt.Sprint(block.Extra()),
			Sha3Uncles:       head.UncleHash.String(),
			Miner:            block.Coinbase().Hex(),
			TransactionsRoot: head.TxHash.String(),
			StateRoot:        head.Root.String(),
		},

		Transactions: transactions,
	}

	return blockObj, nil
}
