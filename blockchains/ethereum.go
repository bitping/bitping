package bitping

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/auser/bitping/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumOptions struct {
	IpcPath string
}

// TODO: make interface for blockchains
type EthereumApp struct {
	Client  *ethclient.Client
	Options EthereumOptions
}

func NewClient(opts EthereumOptions) (*EthereumApp, error) {
	ipcPath := opts.IpcPath
	client, err := ethclient.Dial(ipcPath)
	if err != nil {
		return nil, err
	}

	app := &EthereumApp{
		Client:  client,
		Options: opts,
	}

	return app, nil
}

func (app *EthereumApp) Run(
	blockChan chan *types.Block,
	transChan chan *[]types.Transaction,
	errChan chan error,
) {
	fmt.Printf("Running ethereum\n")
	networkId := app.GetNetwork()
	fmt.Printf("Network id: %v\n", networkId)

	// test
	var headsCh = make(chan *types.Header, 16)
	var errCh = make(chan error, 16)
	go app.SubscribeToNewBlocks(headsCh, errCh)

	for {
		select {
		case err := <-errCh:
			fmt.Printf("Got an error: %v\n", err)
			log.Fatal(err)
		case head := <-headsCh:
			block, err := app.makeBlockFromHeader(head)
			if err != nil {
				errChan <- err
			} else {
				blockChan <- block
			}
			transactions, err := app.makeTransactionsFromHeader(head)
			if err != nil {
				errChan <- err
			} else {
				transChan <- transactions
			}
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

func (app *EthereumApp) SubscribeToNewBlocks(
	heads chan *types.Header,
	errCh chan error,
) {
	var ch = make(chan *types.Header)
	ctx := context.Background()
	sub, err := app.Client.SubscribeNewHead(ctx, ch)
	if err != nil {
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

func (app *EthereumApp) makeBlockFromHeader(
	head *types.Header,
) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	miner := head.Hash()
	block, err := app.Client.BlockByHash(ctx, miner)

	if err != nil {
		return nil, err
	}

	difficulty := types.BigNumber(block.Difficulty().String())
	totalDifficulty := types.BigNumber(head.Difficulty.String())
	cancel()

	blockObj := &types.Block{
		BlockHash:             block.Hash().Hex(),
		BlockNumber:           block.Number().Int64(),
		BlockDifficulty:       difficulty,
		BlockTotalDifficulty:  totalDifficulty,
		BlockNonce:            fmt.Sprint(block.Nonce),
		BlockSize:             float64(block.Size()),
		BlockGasUsed:          block.GasUsed(),
		BlockGasLimit:         block.GasLimit(),
		BlockExtraData:        fmt.Sprint(block.Extra()),
		BlockParentHash:       block.ParentHash().String(),
		BlockSha3Uncles:       head.UncleHash.String(),
		BlockMiner:            miner.String(),
		BlockTransactionsRoot: head.TxHash.String(),
		BlockStateRoot:        head.Root.String(),
	}

	return blockObj, nil
}

func (app *EthereumApp) makeTransactionsFromHeader(
	head *types.Header,
) (*[]types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	miner := head.Hash()
	count, err := app.Client.TransactionCount(ctx, miner)
	block, err := app.Client.BlockByHash(ctx, miner)

	if err != nil {
		return nil, err
	}

	// Is this the right approach?
	var transactions *[]types.Transaction
	var wg sync.WaitGroup

	queue := make(chan *types.Transaction, 1)
	wg.Add(int(count))

	for i := 0; i < int(count); i++ {
		go func(i int) {
			tx, err := app.Client.TransactionInBlock(ctx, miner, uint(i))
			if err == nil {
				sender, err := app.Client.TransactionSender(ctx, tx, miner, uint(i))
				if err == nil {
					fmt.Printf("%#v\n", tx)
					transaction := &types.Transaction{
						BlockHash:        miner.String(),
						BlockNumber:      block.Number().Int64(),
						Hash:             tx.Hash().String(),
						Nonce:            int64(tx.Nonce()),
						TransactionIndex: int64(i),
						From:             sender.String(),
						To:               tx.To().String(),
						Value:            types.BigNumber(fmt.Sprint(tx.Value())),
						GasPrice:         types.BigNumber(fmt.Sprint(tx.Cost())),
						Gas:              types.BigNumber(fmt.Sprint(tx.Gas())),
					}
					queue <- transaction
				}
			}
		}(i)
	}

	go func() {
		for t := range queue {
			*transactions = append(*transactions, *t)
			wg.Done()
		}
	}()

	wg.Wait()

	cancel()

	return transactions, nil
}

/**
* TODO: Add watcher for contract addresses
**/
