package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/samricotta/crs"
	expectedkeepers "github.com/samricotta/crs/expected_keepers"
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
	Decisions  collections.Map[uint64, crs.Decision]                         // key: ID
	Commits    collections.Map[collections.Pair[uint64, []byte], crs.Commit] // key: (decision ID, voter)
	Reveals    collections.Map[collections.Pair[uint64, []byte], crs.Reveal] // key: (decision ID, voter)

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
		Decisions:    collections.NewMap(sb, crs.DecisionKey, "decision", collections.Uint64Key, codec.CollValue[crs.Decision](cdc)),
		Commits:      collections.NewMap(sb, crs.CommitKey, "commit", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey), codec.CollValue[crs.Commit](cdc)),
		Reveals:      collections.NewMap(sb, crs.RevealKey, "reveal", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey), codec.CollValue[crs.Reveal](cdc)),
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

func (k Keeper) EndBlocker(ctx context.Context, id uint64) { // Added id as a parameter
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	participants := [][]byte{}
	now := sdkCtx.BlockTime()
	reveals := []crs.Reveal{}

	err := k.Decisions.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint64, []byte](id),
		func(key uint64, decision crs.Decision) (bool, error) {
			if sdkCtx.BlockTime().After(decision.RevealTimeout) {
				err := k.Reveals.Walk(
					ctx,
					collections.NewPrefixedPairRange[uint64, []byte](id),
					func(key collections.Pair[uint64, []byte], reveal crs.Reveal) (bool, error) {
						participants = append(participants, key.K2())
						reveals = append(reveals, reveal)
						return false, nil
					},
				)
				if err != nil {
					return false, err
				}

			}
			return false, nil
		},
	)
	if err != nil {
		return
	}

	if len(participants) > 0 && now.After(decision.RevealTimeout) {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, crs.ModuleName, participants[0], sdk.NewCoins(decision.EntryFee)); err != nil {
			return
		}

		err = k.RefundAllParticipants(ctx, participants, sdk.NewCoins(decision.EntryFee))
		if err != nil {
			sdkCtx.Logger().Error("Error processing refunds:", "error", err)
			return
		}
	}
}

func (k Keeper) RefundAllParticipants(ctx context.Context, participants [][]byte, amount sdk.Coins) error {
	for _, addr := range participants {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, crs.ModuleName, addr, amount)
		if err != nil {
			return err
		}
	}
	return nil
}
