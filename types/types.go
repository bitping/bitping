package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Header = types.Header
type BlockchainBlock = types.Block

type BigNumber string

type Block struct {
	HeaderHash            common.Hash `json:"headerHash"`
	Network               string      `json:"network"`
	BlockNumber           int64       `json:"blockNumber"`
	BlockHash             string      `json:"blockHash"`
	BlockParentHash       string      `json:"blockParentHash"`
	BlockNonce            string      `json:"blockNonce"`
	BlockSha3Uncles       string      `json:"blockSha3Uncles"`
	BlockLogsBloom        string      `json:"blockLogsBloom"`
	BlockTransactionsRoot string      `json:"blockTransactionsRoot"`
	BlockStateRoot        string      `json:"blockStateRoot"`
	BlockMiner            string      `json:"blockMiner"`
	BlockDifficulty       int64       `json:"blockDifficulty"`
	BlockTotalDifficulty  int64       `json:"blockTotalDifficulty"`
	BlockExtraData        string      `json:"blockExtraData"`
	BlockSize             float64     `json:"blockSize"`
	BlockGasLimit         uint64      `json:"blockGasLimit"`
	BlockGasUsed          uint64      `json:"blockGasUsed"`
	BlockTimeStamp        int64       `json:"blockTimestamp"`
	BlockUncles           []string    `json:"blockUncles"`
}

type Transaction struct {
	BlockHash        string    `json:"blockHash"`
	BlockNumber      int64     `json:"blockNumber"`
	Hash             string    `json:"hash"`
	Nonce            int64     `json:"nonce"`
	TransactionIndex int64     `json:"transactionIndex"`
	From             string    `json:"from"`
	To               string    `json:"to"`
	Value            BigNumber `json:"value"`
	GasPrice         BigNumber `json:"gasPrice"`
	Gas              BigNumber `json:"gas"`
	Input            string    `json:"input"`
}

type Log struct {
	LogIndex         int64  `json:"logIndex"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      int64  `json:"blockNumber"`
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex int64  `json:"transactionIndex"`
	Address          string `json:"address"`
	Data             string `json:"data"`
	Topics           string `json:"topics"`
}

type Receipt struct {
	ReceiptBlockHash         string `json:"receiptBlockHash"`
	ReceiptBlockNumber       int64  `json:"receiptBlockNumber"`
	ReceiptTransactionHash   string `json:"receiptTransactionHash"`
	ReceiptTransactionIndex  int64  `json:"receiptTransactionIndex"`
	ReceiptFrom              string `json:"receiptFrom"`
	ReceiptTo                string `json:"receiptTo"`
	ReceiptCumulativeGasUsed int64  `json:"receiptCumulativeGasUsed"`
	ReceiptGasUsed           int64  `json:"receiptGasUsed"`
	ReceiptContractAddress   string `json:"receiptContractAddress"`
	Logs                     []Log  `json:"logs,omitempty"`
}

type PercentageCalculations struct {
	Percentage           float64
	BlocksToGo           uint64
	Bps                  int
	EstimatedMinutesLeft int
	CurrentBlock         uint64
}
