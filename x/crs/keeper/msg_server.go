package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/samricotta/vote/x/crs"
)

type msgServer struct {
	k Keeper
}

var _ crs.MsgServer = msgServer{}

func NewMsgServerImp() crs.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) NewDecision(ctx context.Context, msg *crs.MsgNewDecision) (*crs.MsgNewDecisionResponse, error) {
	// check if the entry fee is zero
	// check if the voting options are empty
	// convert sender address from string to bytes
	// transfer coins from sender's account to module account:
	// retrieve the next decision ID
	// fetch parameters from the keeper
	// create and store new decision
	// create and store commit
	// return response

	if msg.EntryFee.IsZero() {
		return nil, fmt.Errorf("entry fee is required")
	}

	for _, value := range msg.VotingOptions {
		if value == nil {
			return nil, fmt.Errorf("voting options are required")
		}
	}
	votersAddr, err := ms.k.addressCodec.StringToBytes(msg.Voter)
	if err != nil {
		return nil, fmt.Errorf("invalid sender address: %w", err)
	}

	err = ms.k.bankKeeper.SendCoinsFromAccountToModule(ctx, votersAddr, ModuleName, sdk.NewCoins(msg.EntryFee))
	if err != nil {
		return nil, err
	}

	// retrieve the next decision id
	voteid, err := ms.k.NewDecisionID.Next(ctx)
	if err != nil {
		return nil, err
	}

	params, err := ms.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// store decision
	newDecision := crs.NewDecision{
		Id:            voteid,
		EntryFee:      msg.EntryFee,
		CommitTimeout: sdkCtx.BlockTime().Add(time.Second * time.Duration(params.CommitTimeout)),
	}

	err = ms.k.NewDecision.Set(ctx, voteid, newDecision)
	if err != nil {
		return nil, err
	}

	//store the commit
	commit := crs.Commit{
		Commit: msg.Commit,
	}

	err = ms.k.Commit.Set(ctx, collections.Join(voteid, votersAddr), commit)
	if err != nil {
		return nil, err
	}

	return &crs.MsgNewDecisionResponse{}, nil
}

func saveVoterInfo(decisionID uint64, msg *crs.MsgReveal) {

}

func (ms msgServer) Commit(ctx context.Context, msg *crs.MsgCommit) {

}
