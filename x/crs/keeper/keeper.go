package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/samricotta/vote/x/crs"

	expectedkeepers "github.com/samricotta/vote/x/crs/expected_keepers"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Schema        collections.Schema
	Params        collections.Item[crs.Params]
	NewDecisionID collections.Sequence
	NewDecision   collections.Map[uint64, crs.NewDecision]
	Commit        collections.Map[collections.Pair[uint64, []byte], crs.Commit]
	Reveal        collections.Map[collections.Pair[uint64, []byte], crs.Reveal]

	bankKeeper expectedkeepers.BankKeeper
}

func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, bk expectedkeepers.BankKeeper, authority string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Params:       collections.NewItem(sb, ParamsKey, "params", codec.CollValue[crs.Params](cdc)),
		NewDecision:  collections.NewMap(sb, NewDecisionKey, "new_decision", collections.Uint64Key, codec.CollValue[crs.NewDecision](cdc)),
		Commit:       collections.NewMap(sb, CommitKey, "commit", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey), codec.CollValue[crs.Commit](cdc)),
		Reveal:       collections.NewMap(sb, RevealKey, "reveal", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey), codec.CollValue[crs.Reveal](cdc)),
	}

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

func (k Keeper) GenesisHandler() appmodule.HasGenesis {
	return k.Schema
}

func (k Keeper) Endblocker(ctx context.Context) error {
	// - goes through all of the NewDecisions and see if the timeout has passed for the vote
	// - if the timeout has passed, delete the NewDecision and refund the entry fee
	// - if both reveals are available, delete the NewDecision and announce the winning option
	// - if the reveal timeout has passed, delete the NewDecision and announce the winning option

	newDecisionToDelete := []uint64{}



	return nil
}
