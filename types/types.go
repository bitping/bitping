package types

import (
	"fmt"
	"math/big"

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

type BigInt big.Int

func NewBigInt(i *big.Int) *BigInt {
	if i == nil {
		i = new(big.Int).SetInt64(0)
	}
	bi := BigInt(*i)
	return &bi
}

func BigIntFromInt(i int64) *BigInt {
	bi := new(big.Int).SetInt64(i)
	return NewBigInt(bi)
}

func BigIntFromString(s string) (*BigInt, bool) {
	i, worked := new(big.Int).SetString(s, 10)

	if i != nil {
		bi := BigInt(*i)
		return &bi, worked
	}

	return nil, worked
}

func (i BigInt) MarshalJSON() ([]byte, error) {
	i2 := big.Int(i)
	return []byte(fmt.Sprintf(`"%s"`, i2.String())), nil
}

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
	MedianTime        uint64 `json:"medianTime"`
	Bits              string `json:"bits"`
	Chainwork         string `json:"chainWork"`
	PreviousBlockHash string `json:"previousBlockHash"`
	NextBlockHash     string `json:"nextBlockHash"`
}

// https://github.com/EOSIO/eosjs/blob/master/src/schema/chain_types.json
// {
//   "timestamp": "2018-08-23T17:39:23.500",
//   "producer": "acryptolions",
//   "confirmed": 0,
//   "previous": "00aab1dbad726a64d7a034073c63135e04ef8673bd5e9ef09e91d164bf2935ed",
//   "transaction_mroot": "c9e3c0fa1d236d0cd5ec6527b3f54e954cf6e76595640f1bfd7e8705418ea643",
//   "action_mroot": "207a06bd344056e456ab84ca4f271ca14f1617e18dde9060110f4290159d2d9c",
//   "schedule_version": 196,
//   "new_producers": null,
//   "header_extensions": [],
//   "producer_signature": "SIG_K1_KcCBD54ZyX1JiQqpjMDuN99XZi9kR8wPprqUSgHwUwfdHuMrYauNSh47bb9gHmTZeBz5YLc8qYMWHywBt4qNE5sqUNYLHC",
//   "id": "00aab1dce0ee6cf0457088907ac017f32356dc34790b9a03a4f5147dbce3b3b8",
//   "block_num": 11186652,
//   "ref_block_prefix": 2424860741,
//   "transactions": [],
// }

type EOSBlock struct {
	Producer              string `json:"producer"`
	Confirmed             uint64 `json:"confirmed"`
	TransactionMerkleRoot string `json:"transactionMroot"`
	ActionMerkleRoot      string `json:"actionMroot"`
	ScheduleVersion       uint64 `json:"scheduleVersion"`
	// NewProducers       [] `json:"newProducers"`
	// HeaderExtensions   [] `json:"headerExtensions"`
	ProducerSignature string `json:"producerSignature"`
	RefBlockPrefix    uint64 `json:"refBlockPrefix"`
	ChainID           string `json:"chainID"`
}

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

type Block struct {
	*BitcoinBlock
	*EOSBlock
	*EthereumBlock

	Hash       string  `json:"hash"`
	HeaderHash string  `json:"headerHash"`
	Network    string  `json:"network"`
	NetworkID  int64   `json:"networkID"`
	Number     int64   `json:"number"`
	Size       float64 `json:"size"`
	Time       int64   `json:"time"`

	// Only used for PoW coins
	Nonce      string  `json:"nonce"`
	Difficulty *BigInt `json:"difficulty"`

	// BTC PreviousBlockHash
	// EOS Previous
	ParentHash string `json:"parentHash"`

	Transactions []Transaction `json:"transactions"`

	// Batched Transactions are broken up
	SingletonTransactions []Transaction `json:"singletonTransactions"`
}

// {
//       "status": "executed",
//       "cpu_usage_us": 1696,
//       "net_usage_words": 28,
//       "trx": {
//         "id": "4528f0f169031bd02402ef096095b4801cb6766b7462f8cfc24626cb2dbd8e92",
//         "signatures": [
//           "SIG_K1_KerL5zixYTk7ojcosM2iG8v2vpu1bbocQXoipmQNevjvRXGxhT6AcND9Ts7MTd1rFMsmn6W2LT3BqZqUYss853YfuhtJUf"
//         ],
//         "compression": "none",
//         "packed_context_free_data": "",
//         "context_free_data": [],
//         "packed_trx": "c0f17e5b8cb0b810b9890000000001808ec958e5ab983b000080f1a86c52d501808ec958e5ab983b00000000a8ed32328101808ec958e5ab983b366ade6765010000d0b1aa000000000017323031382d30382d32335431373a33393a31372e3530302a090000000000004030306161623164303562636364373364323830323736346232633935366165613736343736353065336230336463363238306134393534646462353161326530010000000000000000",
//         "transaction": {
//           "expiration": "2018-08-23T17:41:20",
//           "ref_block_num": 45196,
//           "ref_block_prefix": 2310607032,
//           "max_net_usage_words": 0,
//           "max_cpu_usage_ms": 0,
//           "delay_sec": 0,
//           "context_free_actions": [],
//           "actions": [
//             {
//               "account": "bigertestabc",
//               "name": "updateblk",
//               "authorization": [
//                 {
//                   "actor": "bigertestabc",
//                   "permission": "active"
//                 }
//               ],
//               "data": {
//                 "account": "bigertestabc",
//                 "dealId": "1535045954102",
//                 "eosRealblock": 11186640,
//                 "blktime": "2018-08-23T17:39:17.500",
//                 "blockSize": 2346,
//                 "blockId": "00aab1d05bccd73d2802764b2c956aea7647650e3b03dc6280a4954ddb51a2e0",
//                 "txCounts": 1
//               },
//               "hex_data": "808ec958e5ab983b366ade6765010000d0b1aa000000000017323031382d30382d32335431373a33393a31372e3530302a0900000000000040303061616231643035626363643733643238303237363462326339353661656137363437363530653362303364633632383061343935346464623531613265300100000000000000"
//             }
//           ],
//           "transaction_extensions": []
//         }
//       }
//     }

type EOSExtension struct {
	Type uint64 `json:"type"`
	Data string `json:"data"`
}

type EOSPermissionLevel struct {
	Actor       string `json:"actor"`
	Permisssion string `json:"permission"`
}

type EOSAction struct {
	Account       string               `json:"account"`
	Name          string               `json:"name"`
	Authorization []EOSPermissionLevel `json:"authorization"`
	HexData       string               `json:"hexData"`
	Data          interface{}          `json:"data"`
}

type EOSUnpackedTransaction struct {
	Expiration              int64          `json:"expiration"`
	RefBlockNum             uint64         `json:"refBlockNum"`
	RefBlockPrefix          uint64         `json:"refBlockPrefix"`
	MaxNetUsageWords        uint64         `json:"maxNetUsageWords"`
	MaxCPUUsageMicroSeconds uint64         `json:"maxCPUUsageMS"`
	DelaySec                uint64         `json:"delaySec"`
	Actions                 []EOSAction    `json:"actions"`
	ContextFreeActions      []EOSAction    `json:"contextFreeActions"`
	TransactionExtensions   []EOSExtension `json:"transactionExtensions"`
}

type EOSTransactionWithID struct {
	ID                    string                 `json:"id"`
	Signatures            []string               `json:"signatures"`
	Compression           string                 `json:"compression"`
	PackedTRX             string                 `json:"packedTRX"`
	PackedContextFreeData string                 `json:"packedContextFreeData"`
	ContextFreeData       []string               `json:"contextFreeData"`
	Transaction           EOSUnpackedTransaction `json:"transaction"`
}

type EOSTransactionReceipt struct {
	Status               string               `json:"status"`
	CPUUsageMicroSeconds uint64               `json:"cpuUsageUS"`
	NetUsageWords        uint64               `json:"netUsageWords"`
	TRX                  EOSTransactionWithID `json:"trx"`
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

type EthereumTransaction struct {
	TransactionIndex int64   `json:"transactionIndex"`
	GasPrice         *BigInt `json:"gasPrice"`
	Gas              uint64  `json:"gas"`
}

// UTX is address
type Transaction struct {
	*EOSTransactionReceipt
	*EthereumTransaction

	BlockHash   string  `json:"blockHash"`
	BlockNumber int64   `json:"blockNumber"`
	Hash        string  `json:"hash"`
	Nonce       int64   `json:"nonce"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Value       *BigInt `json:"value"`
	Symbol      string  `json:"symbol"`
	Precision   uint64  `json:"precision"`
	Data        []byte  `json:"data"`

	// Is this a tx that's split form
	IsSingleton    bool   `json:"isSplit"`
	SingletonIndex uint64 `json:"singletonIndex"`
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
