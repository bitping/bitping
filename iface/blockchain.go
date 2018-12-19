package iface

import (
	"github.com/auser/bitping/types"
)

type BlockListener interface {
	Listen(ch chan types.Block, err chan error)
}

type TransactionListener interface {
	Listen(ch chan types.Transaction, err chan error)
}

type BlockchainListener interface {
	Block() BlockListener
	Transaction() TransactionListener
}
