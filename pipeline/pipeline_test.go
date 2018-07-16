package pipeline

import (
	// "context"
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/auser/bitping/iface"
	"github.com/auser/bitping/types"
	. "github.com/auser/bitping/util/ginkgo"
)

func Test(t *testing.T) {
	Setup("pipeline", t)
}

type block struct {
}

func (b *block) Listen(ch chan types.Block, err chan error) {
	fmt.Printf("Listen...\n")
}

var _ = Describe("Stage", func() {
	var (
		stage *Stage
		ctx   context.Context
		in    chan iface.Block
	)

	BeforeEach(func() {
		log.Printf("beforeeach")
		ctx = context.Background()
		in = make(chan iface.Block)
		stage = new(Stage)
	})

	AfterEach(func() {

	})

	It("should run all steps in a stage", func() {
		defer close(in)
		num := 0

		add1 := func(block iface.Block) (iface.Block, error) {
			num = num + 1
			return block, nil
		}

		add2 := func(block iface.Block) (iface.Block, error) {
			num = num + 2
			return block, nil
		}

		stage.addStep(add1)
		stage.addStep(add2)

		out, _, err := stage.Run(ctx, in)

		go func() {
			in <- new(block)
			in <- new(block)
		}()

		// wait for blocks to return
		<-out
		<-out

		ctx.Done()

		Expect(err).NotTo(HaveOccurred())
		Expect(num).To(Equal(6))
		log.Print("test done")
	})
})
