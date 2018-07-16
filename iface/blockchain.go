package iface

import (
	"github.com/auser/bitping/types"
)

type Block interface {
	Listen(ch chan types.Block, err chan error)
}

type Transaction interface {
	Listen(ch chan types.Transaction, err chan error)
}

type Blockchain interface {
	Block() Block
	Transaction() Transaction
}
