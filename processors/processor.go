package processors

type Processor interface {
	Init()
	Process()
}

type ProcessorNet map[string]Processor
