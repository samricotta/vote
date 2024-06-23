package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
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
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// iterate through all the new decisions
	err := k.NewDecision.Walk(ctx, nil, func(key uint64, game crs.NewDecision) (bool, error) {
		now := sdkCtx.BlockTime()
		
		//voters that committed
		votersCommitted := [][]byte{}
		commits := []crs.Commit{}
		// iterate through all the voters that committed
		err := k.Commit.Walk(
			ctx, 
			collections.NewPrefixedPairRange[uint64, []byte](id),
			func(key uint64, value []byte) error {
				votersCommitted = append(votersCommitted, value)
				return nil
			},
		)
		if err != nil {
			return false, err
		}

		// if the game has less than 0 voters and the commit timeout has passed, delete the game and refund the entry fee
		if  len(votersCommitted) == 0 && now.After(newdecision.CommitTimeout) {
			newDecisionToDelete = append(newDecisionToDelete, key)
			return false, nil
		}

		// no voter revealed so there's no reveal timeout set yet
		if newDecision.timeout.isZero() {
			return false, nil
		}

		votersRevealed := [][]byte{}
		reveals := []crs.Reveal{}
		// iterate through all the reveals
		err = k.Reveal.Walk(
			ctx,
			collections.NewPrefixedPairRange[uint64, []byte](id),
			func(key collections.Pair[uint64, []byte], value crs.Reveal) (bool, error) {
				votersRevealed = append(votersRevealed, value)
				return false, nil
			},
		)
		if err != nil {
			return false, err
		}

		// if the reveal timeout hasnt passed and less than 2 players are revealed, lets wait
		if len(votersRevealed) < 2 && now.Before(newDecision.timeout) {
			return false, nil
		}

		// this game is over, lets delete it
		newDecisionToDelete :=  append(newDecisionToDelete, id)
	
		//calculate the winning option
		optionCounts := make(map[string]int)
		var winningOption string
		maxCount := 0

		for _, reveal := range reveals {
			optionCounts[reveal.Options]++
			if optionCounts[reveal.Options] > maxCount {
				maxCount = optionCounts[reveal.Options]
				winningOption = reveal.Options
			}
		} 

		for _, id := range newDecisionToDelete {
			if err := k.NewDecision.Remove(ctx, id); err != nil {
				return false, err
			}
		}

	return nil
}


