package bitping

import (
	"context"
	"fmt"
	"log"
	"math/big"
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

func NewClient(opts EthereumOptions) *EthereumApp {
	ipcPath := opts.IpcPath
	client, err := ethclient.Dial(ipcPath)
	if err != nil {
		log.Fatal(err)
	}

	app := &EthereumApp{
		Client:  client,
		Options: opts,
	}

	return app
}

func (app *EthereumApp) Run() {
	fmt.Printf("Running ethereum\n")
	networkId := app.GetNetwork()
	fmt.Printf("Network id: %v\n", networkId)

	// test
	var headsCh = make(chan *types.Header, 16)
	var errCh = make(chan error, 16)
	go app.SubscribeToAddressEvents(headsCh, errCh)

	for {
		select {
		case err := <-errCh:
			fmt.Printf("Got an error: %v\n", err)
			log.Fatal(err)
		case head := <-headsCh:
			block, err := app.makeBlockFromHeader(head)
			if err == nil {
				fmt.Printf("Got a block: %#v\n", block)
			}
			transactions, err := app.makeTransactionsFromHeader(head)
			if err == nil {
				fmt.Printf("Txs: %#v\n", transactions)
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

func (app *EthereumApp) SubscribeToAddressEvents(
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

	for _, uncle := range block.Uncles() {
		fmt.Printf("uncle: %#v\n", uncle)
	}

	for _, tx := range block.Transactions() {
		fmt.Printf("tx: %#v\n", tx)
	}

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
) ([]types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	miner := head.Hash()
	count, err := app.Client.TransactionCount(ctx, miner)

	if err != nil {
		return nil, err
	}

	fmt.Printf("count: %d\n", count)

	return nil, nil
}
