package processors

import (
	"fmt"

	"github.com/auser/bitping/types"
)

type Listener struct {
	In  <-chan types.Block
	Out chan<- types.Block
}

func (l *Listener) Init() {
	fmt.Printf("Initializing Listener\n")
}

func (l *Listener) Process() {
	go func() {
		for {
			s, ok := <-l.In
			if !ok {
				fmt.Println("Listener finished")
				close(l.Out)
				return
			}
			l.Out <- s
		}
	}()
}
