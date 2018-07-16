package pipeline

import (
	// "context"
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
		// ctx   context.Context
		// in    chan iface.Block
	)

	BeforeEach(func() {
		log.Printf("beforeeach")
		// ctx = context.Background()
		// in = make(chan iface.Block)
		stage = new(Stage)
	})

	AfterEach(func() {

	})

	It("should run all steps in a stage", func() {
		log.Printf("it")
		num := 0

		add1 := func(block iface.Block) (iface.Block, error) {
			log.Printf("add1")
			num = num + 1
			log.Printf("%v", num)
			return block, nil
		}

		add2 := func(block iface.Block) (iface.Block, error) {
			log.Printf("add2")
			num = num + 2
			log.Printf("%v", num)
			return block, nil
		}

		log.Printf("waiting...")
		stage.addStep(add1)
		stage.addStep(add2)

		// _, _, err := stage.Run(ctx, in)
		// log.Printf("add block")
		// in <- new(block)
		// log.Printf("add block")
		// in <- new(block)
		// log.Printf("done")
		// close(in)
		// ctx.Done()

		// Expect(err).NotTo(HaveOccurred())
		Expect(num).To(Equal(3))
	})
})
