package processors

import (
	"fmt"

	"github.com/auser/bitping/types"
)

type Printer struct {
	In   <-chan types.Block
	Done chan<- struct{}
}

func (p *Printer) Init() {
	fmt.Printf("Initializing Printer\n")
}

func (p *Printer) Process() {
	go func() {
		for {
			c, ok := <-p.In
			if !ok {
				fmt.Println("Printer finished")
				close(p.Done)
				return
			}
			fmt.Printf("Block number: %#v\n", c.BlockNumber)
		}
	}()
}
