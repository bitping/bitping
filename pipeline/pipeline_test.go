package pipeline

import (
	// "context"
	"context"
	"errors"
	"fmt"
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

		in <- new(block)
		in <- new(block)

		// wait for blocks to return
		<-out
		<-out

		ctx.Done()

		Expect(err).NotTo(HaveOccurred())
		Expect(num).To(Equal(6))
	})

	It("should clean up channels on failure", func() {

		defer close(in)
		num := 0

		add1 := func(block iface.Block) (iface.Block, error) {
			num = num + 1
			return block, nil
		}

		err1 := func(block iface.Block) (iface.Block, error) {
			return nil, errors.New("oops")
		}

		stage.addStep(add1)
		stage.addStep(err1)

		_, errc, err := stage.Run(ctx, in)
		Expect(err).NotTo(HaveOccurred())

		in <- new(block)

		// wait for blocks to return
		err = <-errc
		Expect(err).To(HaveOccurred())
	})
})
