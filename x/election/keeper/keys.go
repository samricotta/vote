package keeper

import "cosmossdk.io/collections"

const ModuleName = "election"

var (
	ElectionKey = collections.NewPrefix(0)
)
