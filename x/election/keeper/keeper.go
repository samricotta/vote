package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/samricotta/vote/x/election"

	storetypes "cosmossdk.io/core/store"

	expectedkeepers "github.com/samricotta/vote/x/election/expected_keepers"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	authority string

	// state management
	Schema     collections.Schema
	ElectionID collections.Sequence
	Election   collections.Map[uint64, election.Election] // key: ID
	Params     collections.Item[election.Params]

	bankKeeper expectedkeepers.BankKeeper
	crsKeeper  expectedkeepers.CrsKeeper
}

func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, bk expectedkeepers.BankKeeper, crs expectedkeepers.CrsKeeper, authority string) Keeper {
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
		Election:     collections.NewMap(sb, ElectionKey, "election", collections.Uint64Key, codec.CollValue[election.Election](cdc)),
		bankKeeper:   bk,
		crsKeeper:    crs,
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
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.crsKeeper.
	// check if the election has expired
	// check if the election has been resolved
	// create a crs.NewDecision object
	// store the new decision
	// resolve the election
	// delete the election

}
