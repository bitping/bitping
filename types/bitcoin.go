package types

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
