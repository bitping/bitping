package iface

type Step interface {
	Execute()
}

type Stage interface {
	AddChain(blockchain Blockchain)
	RemoveChain(blockchain Blockchain)
	AddStep(step Step)
	RemoveStep(step Step)
	AddStage(stage ...Stage)
	RemoveStage(stage Stage)
}
