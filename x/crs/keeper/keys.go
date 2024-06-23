package keeper

import "cosmossdk.io/collections"

const ModuleName = "crs"

var (
	ParamsKey      = collections.NewPrefix(0)
	NewDecisionKey = collections.NewPrefix(1)
	CommitKey      = collections.NewPrefix(2)
	RevealKey      = collections.NewPrefix(4)
)

