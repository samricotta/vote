package keeper

import "cosmossdk.io/collections"

const ModuleName = "vote"

var (
	VoteKey      = collections.NewPrefix(0)
	ResolveVoteKey = collections.NewPrefix(1)
)

