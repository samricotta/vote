package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

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
	Schema     collections.Schema
	Params     collections.Item[crs.Params]
	DecisionID collections.Sequence
	Decision   collections.Map[uint64, crs.Decision]                         // key: ID
	Commit     collections.Map[collections.Pair[uint64, []byte], crs.Commit] // key: (decision ID, voter)
	Reveal     collections.Map[collections.Pair[uint64, []byte], crs.Reveal] // key: (decision ID, voter)

	bankKeeper expectedkeepers.BankKeeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, bankKeeper expectedkeepers.BankKeeper, authority string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Params:       collections.NewItem(sb, crs.ParamsKey, "params", codec.CollValue[crs.Params](cdc)),
		DecisionID:   collections.NewSequence(sb, crs.DecisionIDKey, "decision_id"),
		Decision:     collections.NewMap(sb, crs.DecisionKey, "decision", collections.Uint64Key, codec.CollValue[crs.Decision](cdc)),
		Commit:       collections.NewMap(sb, crs.CommitKey, "commit", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey), codec.CollValue[crs.Commit](cdc)),
		Reveal:       collections.NewMap(sb, crs.RevealKey, "reveal", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey), codec.CollValue[crs.Reveal](cdc)),
		bankKeeper:   bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
