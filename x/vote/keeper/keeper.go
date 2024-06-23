package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/samricotta/vote/x/vote"

	storetypes "cosmossdk.io/core/store"

	expectedkeepers "github.com/samricotta/vote/x/crs/expected_keepers"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	authority string

	// state management
	Schema      collections.Schema
	VoteID      collections.Sequence
	Vote        collections.Map[uint64, vote.Vote]
	ResolveVote collections.Map[uint64, vote.ResolveVote]

	bankKeeper expectedkeepers.BankKeeper
}

func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, bk expectedkeepers.BankKeeper, authority string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	// build collections schema
	sb := collections.NewSchemaBuilder(storeService)

	// create keeper
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Vote:         collections.NewMap(sb, VoteKey, "vote", collections.Uint64Key, codec.CollValue[vote.Vote](cdc)),
		ResolveVote:  collections.NewMap(sb, ResolveVoteKey, "resolve_vote", collections.Uint64Key, codec.CollValue[vote.ResolveVote](cdc)),
	}

	//build schema

	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build schema: %w", err))
	}

	k.Schema = schema

	return k
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) GetGenesisHandler() appmodule.HasGenesis {
	return nil
}

func (k Keeper) Endblocker(ctx context.Context) {
	// check if the vote has expired
	// check if the vote has been resolved
	// create a crs.NewDecision object
	// store the new decision
	// resolve the vote
	// delete the vote
}
