package iface

type Step interface {
	Execute()
}

type Stage interface {
	AddChain(blockchain BlockchainListener)
	RemoveChain(blockchain BlockchainListener)
	AddStep(step Step)
	RemoveStep(step Step)
	AddStage(stage ...Stage)
	RemoveStage(stage Stage)
}
