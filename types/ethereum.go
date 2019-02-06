package types

import (
	t "github.com/ethereum/go-ethereum/core/types"
)

type GethHeader = t.Header
type GethBlock = t.Block
type GethHomesteadSigner = t.HomesteadSigner

// https://github.com/ethereum/wiki/wiki/JavaScript-API#web3ethgetblock
// {
//   "number": 3,
//   "hash": "0xef95f2f1ed3ca60b048b4bf67cde2195961e0bba6f70bcbea9a2c4e133e34b46",
//   "parentHash": "0x2302e1c0b972d00932deb5dab9eb2982f570597d9d42504c05d9c2147eaf9c88",
//   "nonce": "0xfb6e1a62d119228b",
//   "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
//   "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
//   "transactionsRoot": "0x3a1b03875115b79539e5bd33fb00d8f7b7cd61929d5a3c574f507b8acf415bee",
//   "stateRoot": "0xf1133199d44695dfa8fd1bcfe424d82854b5cebef75bddd7e40ea94cda515bcb",
//   "miner": "0x8888f1f195afa192cfee860698584c030f4c9db1",
//   "difficulty": BigNumber,
//   "totalDifficulty": BigNumber,
//   "size": 616,
//   "extraData": "0x",
//   "gasLimit": 3141592,
//   "gasUsed": 21662,
//   "timestamp": 1429287689,
//   "transactions": [
//     "0x9fc76417374aa880d4449a1f7f31ec597f00b1f6f3dd2d66f4c9c6c445836d8b"
//   ],
//   "uncles": []
// }

type EthereumBlock struct {
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	TotalDifficulty  *BigInt  `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         uint64   `json:"gasLimit"`
	GasUsed          uint64   `json:"gasUsed"`
	Uncles           []string `json:"uncles"`
}

// {
//   "hash": "0x9fc76417374aa880d4449a1f7f31ec597f00b1f6f3dd2d66f4c9c6c445836d8b",
//   "nonce": 2,
//   "blockHash": "0xef95f2f1ed3ca60b048b4bf67cde2195961e0bba6f70bcbea9a2c4e133e34b46",
//   "blockNumber": 3,
//   "transactionIndex": 0,
//   "from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
//   "to": "0x6295ee1b4f6dd65047762f924ecd367c17eabf8f",
//   "value": BigNumber,
//   "gas": 314159,
//   "gasPrice": BigNumber,
//   "input": "0x57cb2fc4"
// }

type EthereumCall struct {
	Input []byte `json:"input"`
}

type EthereumEvent struct {
	LogIndex         uint64   `json:"logIndex"`
	TransactionIndex uint64   `json:"transactionIndex"`
	Address          string   `json:"address"`
	Data             []byte   `json:"data"`
	Topics           []string `json:"topics"`
	Removed          bool     `json:"bool"`
}

type EthereumTransaction struct {
	TransactionIndex int64   `json:"transactionIndex"`
	GasPrice         *BigInt `json:"gasPrice"`
	Gas              uint64  `json:"gas"`
}
