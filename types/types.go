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

// {
// hash: "0000000000000000002662f1c485796dfd585c63daf606231e548a78c78dcaab",
// ver: 536870912,
// prev_block: "00000000000000000025dd1f000e816c0947d056beb498ee308057b33f98057c",
// mrkl_root: "3494036329589c46408d5b2a078b04bcb0d1118759a2e2a65ded13a7d122b1b8",
// time: 1534846055,
// bits: 388763047,
// fee: 15664446,
// nonce: 835124619,
// n_tx: 1371,
// size: 750983,
// block_index: 1719065,
// main_chain: true,
// height: 537787,
// received_time: 1534846055,
// relayed_by: "0.0.0.0",
// tx: []
// }

type BitcoinBlock struct {
	Height            uint64 `json:"height"`
	Confirmations     uint64 `json:"confirmations"`
	StrippedSize      uint64 `json:"strippedSize"`
	Weight            uint64 `json:"weight"`
	Version           string `json:"version"`
	VersionHex        string `json:"versionHex"`
	Merkleroot        string `json:"merkleroot"`
	Time              uint64 `json:"time"`
	MedianTime        uint64 `json:"medianTime"`
	Bits              string `json:"bits"`
	Chainwork         string `json:"chainWork"`
	PreviousBlockHash string `json:"previousBlockHash"`
	NextBlockHash     string `json:"nextBlockHash"`
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

	Difficulty float64 `json:"difficulty"`
	Hash       string  `json:"hash"`
	HeaderHash string  `json:"headerHash"`
	Network    string  `json:"network"`
	NetworkID  int64   `json:"networkID"`
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
