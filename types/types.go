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

type Block struct {
	Map

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
	// Height
	// StrippedSize
	// Weight
	// Version
	// VersionHex
	// Merkleroot
	// Time
	// MedianTime
	// Bits
	// Chainwork
	// PreviousBlockHash

	// Ethereum Data
	// Sha3Uncles       string   `json:"Sha3Uncles"`
	// LogsBloom        string   `json:"LogsBloom"`
	// TransactionsRoot string   `json:"TransactionsRoot"`
	// StateRoot        string   `json:"StateRoot"`
	// Miner            string   `json:"Miner"`
	// TotalDifficulty  int64    `json:"TotalDifficulty"`
	// ExtraData        string   `json:"ExtraData"`
	// GasLimit         uint64   `json:"GasLimit"`
	// GasUsed          uint64   `json:"GasUsed"`
	// TimeStamp        int64    `json:"Timestamp"`
	// Uncles           []string `json:"Uncles"`
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
