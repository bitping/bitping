package types

import (
	t "github.com/ethereum/go-ethereum/core/types"
)

type Header = t.Header
type BlockchainBlock = t.Block
type HomesteadSigner = t.HomesteadSigner

type BlockChainRunner interface {
	NewClient()
	Run()
}

type Map map[string]interface{}

type BitcoinBlock struct {
	Height            uint64
	StrippedSize      uint64
	Weight            uint64
	Version           string
	VersionHex        string
	Merkleroot        string
	Time              uint64
	MedianTime        uint64
	Bits              string
	Chainwork         string
	PreviousBlockHash string
}

type EthereumBlock struct {
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	TotalDifficulty  int64    `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         uint64   `json:"gasLimit"`
	GasUsed          uint64   `json:"gasUsed"`
	TimeStamp        int64    `json:"timestamp"`
	Uncles           []string `json:"uncles"`
}

type Block struct {
	*BitcoinBlock
	*EthereumBlock

	Difficulty int64   `json:"difficulty"`
	Hash       string  `json:"hash"`
	HeaderHash string  `json:"headerHash"`
	Network    string  `json:"network"`
	Nonce      string  `json:"nonce"`
	Number     int64   `json:"number"`
	Size       float64 `json:"size"`

	// BTC PreviousBlockHash
	ParentHash string `json:"parentHash"`

	Transactions []Transaction `json:"transactions"`

	// Bitcoin Data
	// Ethereum Data
}

// UTX is address
type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      int64  `json:"blockNumber"`
	Hash             string `json:"hash"`
	Nonce            int64  `json:"nonce"`
	TransactionIndex int64  `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            int64  `json:"value"`
	GasPrice         int64  `json:"gasPrice"`
	Gas              uint64 `json:"gas"`
	Input            string `json:"input"`
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
