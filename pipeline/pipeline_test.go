package pipeline

import (
	"context"
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

func add1(block iface.Block) (iface.Block, error) {
	num = num + 1
	return block, nil
}

func add2(block iface.Block) (iface.Block, error) {
	num = num + 2
	return block, nil
}

var _ = Describe("Stage", func() {
	var (
		stage *Stage
		ctx   context.Background
	)

	BeforeEach(func() {
		ctx = context.Background()
		in = make(chan iface.Block)
		stage = new(Stage)
	})

	It("should run all steps in a stage", func() {
		num := 0

		stage.addStep(add1)
		stage.addStep(add2)
		out, errc, err := stage.Run(ctx, in)

		Expect(err).NotToHaveOccured()
		Expect(num).To(Equal(3))
	})
})
