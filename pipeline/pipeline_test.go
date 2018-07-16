package pipeline

import (
	"context"
	"errors"
	"testing"

	"github.com/auser/bitping/iface"
	. "github.com/auser/bitping/util/ginkgo"
)

func Test(t *testing.T) {
	Setup("pipeline", t)
}

func noopStep(block iface.Block) (iface.Block, error) {
	return block, nil
}

func errStep(block iface.Block) (iface.Block, error) {
	return nil, errors.New("blew up")
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

	It("should run all steps in a stage", func() {
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
		_, _, err := stage.Run(ctx, in)

		Expect(err).NotTo(HaveOccurred())
		Expect(num).To(Equal(3))
	})
})
