package types

type Block struct {
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
}

type Transaction struct {
	*EOSTransactionReceipt
	*EthereumTransaction

	BlockHash       string `json:"blockHash"`
	BlockNumber     int64  `json:"blockNumber"`
	TransactionHash string `json:"transactionHash"`
	Hash            string `json:"hash"`
	Nonce           int64  `json:"nonce"`

	// Is this a tx that's split form
	IsDerived    bool `json:"isSplit"`
	DerivedIndex int  `json:"derivedIndex"`

	Actions []Action `json:"actions"`
	Events  []Event  `json:"events"`
}

// Actions/Events/Function Calls

type Event struct {
	*EthereumEvent
}

// UTXO is address
type Action struct {
	*EOSAction
	*EthereumCall

	BlockHash       string `json:"blockHash"`
	BlockNumber     int64  `json:"blockNumber"`
	TransactionHash string `json:"transactionHash"`
	Address         string `json:"address"`
	Data            []byte `json:"data"`

	From      string  `json:"from"`
	To        string  `json:"to"`
	Value     *BigInt `json:"value"`
	Symbol    string  `json:"symbol"`
	Precision uint64  `json:"precision"`

	In  uint64 `json:"in"`
	Out uint64 `json:"out"`
}
