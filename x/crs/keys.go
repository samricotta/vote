package crs

import "cosmossdk.io/collections"

const ModuleName = "crs"

var (
	ParamsKey     = collections.NewPrefix(0)
	DecisionIDKey = collections.NewPrefix(1)
	DecisionKey   = collections.NewPrefix(2)
	CommitKey     = collections.NewPrefix(3)
	RevealKey     = collections.NewPrefix(4)
)
