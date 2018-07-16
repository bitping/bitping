package pipeline

import (
	"context"

	"github.com/auser/bitping/iface"
)

type Source interface{}
type Sink interface{}
type Step func(block iface.Block) (iface.Block, error)

type Stage struct {
	steps []Step
	sinks []Sink
}

func (p *Stage) addStep(step Step) {
	p.steps = append(p.steps, step)
}

func (p *Stage) addSink(sink Sink) {
	p.sinks = append(p.sinks, sink)
}

func (p *Stage) Run(ctx context.Context, in <-chan iface.Block) (<-chan iface.Block, <-chan error, error) {
	out := make(chan iface.Block)
	errc := make(chan error, 1)

	go func() {
		// defer close(out)
		// defer close(errc)

		for block := range in {
			var err error
			// For each step in stage, modify block / do work
			for _, step := range p.steps {
				// Use block returned by step
				block, err = step(block)
				if err != nil {
					errc <- err
					return
				}
			}
			select {
			case out <- block:
			case <-ctx.Done():
				return
			}

			// Feed final block output to next stage and/or wait for context to
			// close
		}
	}()

	return out, errc, nil
}
