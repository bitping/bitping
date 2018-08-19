package bitping

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	types "github.com/auser/bitping/types"
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
	Client  *ethclient.Client
	Options EthereumOptions
	types.BlockChainRunner
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

	return app, nil
}

func (app *EthereumApp) Run(
	blockChan chan types.Block,
	// transChan chan []types.Transaction,
	errChan chan error,
) {
	fmt.Printf("Running ethereum\n")
	networkId := app.GetNetwork()
	fmt.Printf("Network id: %v\n", networkId)

	// test
	var headsCh = make(chan *types.Header)
	var errCh = make(chan error)
	go app.SubscribeToNewBlocks(headsCh, errCh)

	for {
		select {
		case err := <-errCh:
			fmt.Printf("Got an error in client.Run(): %v\n", err)
			// TODO: Reconnect here
			// go app.SubscribeToNewBlocks(headsCh, errCh)
		case head := <-headsCh:
			block, err := app.GetBlockFromHeader(head)
			if err != nil {
				fmt.Printf("Error happened: %s\n", err.Error())
				errChan <- err
			} else {
				blockChan <- block
			}
			// transactions, err := app.makeTransactionsFromBlock(block)
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

func (app *EthereumApp) SubscribeToNewBlocks(
	heads chan *types.Header,
	errCh chan error,
) {
	var ch = make(chan *types.Header)
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

func (app *EthereumApp) getBlockByNumWithBackoff(num *big.Int) (*types.BlockchainBlock, error) {
	ctx := context.Background()
	var (
		block *types.BlockchainBlock
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
				return &types.BlockchainBlock{}, err
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

func (app *EthereumApp) GetBlockFromHeader(
	head *types.Header,
) (types.Block, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	// right now, this is blocking... do we want it to block?
	block, err := app.getBlockByNumWithBackoff(head.Number)
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
		if msg, err := tx.AsMessage(types.HomesteadSigner{}); err != nil {
			txFromStr = msg.From().Hex()
			if msg.To() != nil {
				txToStr = msg.To().Hex()
			}
		}

		transaction := types.Transaction{
			BlockHash:        block.Hash().Hex(),
			BlockNumber:      block.Number().Int64(),
			Hash:             tx.Hash().String(),
			Nonce:            int64(tx.Nonce()),
			TransactionIndex: int64(i),
			From:             txFromStr,
			To:               txToStr,
			Value:            types.BigNumber(fmt.Sprint(tx.Value())),
			GasPrice:         types.BigNumber(fmt.Sprint(tx.Cost())),
			Gas:              types.BigNumber(fmt.Sprint(tx.Gas())),
		}
		transactions = append(transactions, transaction)
	}

	blockObj := types.Block{
		Network:               "ethereum",
		HeaderHash:            head.Hash().Hex(),
		BlockHash:             block.HashNoNonce().Hex(),
		BlockNumber:           block.Number().Int64(),
		BlockDifficulty:       block.Difficulty().Int64(),
		BlockTotalDifficulty:  block.Difficulty().Int64(), // make sense?
		BlockNonce:            fmt.Sprint(block.Nonce()),
		BlockSize:             float64(block.Size()),
		BlockGasUsed:          block.GasUsed(),
		BlockGasLimit:         block.GasLimit(),
		BlockExtraData:        fmt.Sprint(block.Extra()),
		BlockParentHash:       block.ParentHash().String(),
		BlockSha3Uncles:       head.UncleHash.String(),
		BlockMiner:            block.Hash().Hex(),
		BlockTransactionsRoot: head.TxHash.String(),
		BlockStateRoot:        head.Root.String(),
		Transactions:          transactions,
	}

	return blockObj, nil
}
