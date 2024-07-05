package keeper

import (
	"context"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/samricotta/vote/x/crs"
	"github.com/samricotta/vote/x/election"
)

type msgServer struct {
	k Keeper
}

var _ election.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) election.MsgServer {
	return &msgServer{k: keeper}
}

// UpdateParams params is defining the handler for the MsgUpdateParams message.
func (ms msgServer) UpdateParams(ctx context.Context, msg *election.MsgUpdateParams) (*election.MsgUpdateParamsResponse, error) {
	if _, err := ms.k.addressCodec.StringToBytes(msg.Authority); err != nil {
		return nil, fmt.Errorf("invalid authority address: %w", err)
	}

	if authority := ms.k.GetAuthority(); !strings.EqualFold(msg.Authority, authority) {
		return nil, fmt.Errorf("unauthorized, authority does not match the module's authority: got %s, want %s", msg.Authority, authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := ms.k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &election.MsgUpdateParamsResponse{}, nil
}

func (ms msgServer) NewElection(ctx context.Context, msg *election.MsgNewElection) (*election.MsgNewElectionResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("election is required")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	commitTimeout := sdkCtx.BlockTime().Add(time.Minute * 2)
	revealTimeout := commitTimeout.Add(time.Minute * 2)
	options := [][]byte{}
	for _, v := range msg.Options {
		options = append(options, []byte(v))
	}

	decision := crs.Decision{
		EntryFee:      sdk.NewInt64Coin("stake", 1000),
		Options:       options,
		CommitTimeout: commitTimeout,
		RevealTimeout: revealTimeout,
		Refund:        true,
		PayoutAddress: "",
	}

	if err := ms.k.crsKeeper.CreateDecision(ctx, decision); err != nil {
		return nil, err
	}

	return &election.MsgNewElectionResponse{}, nil
}
