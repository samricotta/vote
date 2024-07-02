package election

import "cosmossdk.io/collections"

const ModuleName = "election"

var (
	ParamsKey  = collections.NewPrefix(0)
	CounterKey = collections.NewPrefix(1)
)
