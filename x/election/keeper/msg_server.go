package keeper

import (
	"context"
	"fmt"
	"strings"

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

// IncrementCounter defines the handler for the MsgIncrementCounter message.
// func (ms msgServer) IncrementCounter(ctx context.Context, msg *election.MsgIncrementCounter) (*election.MsgIncrementCounterResponse, error) {
// 	if _, err := ms.k.addressCodec.StringToBytes(msg.Sender); err != nil {
// 		return nil, fmt.Errorf("invalid sender address: %w", err)
// 	}

// 	counter, err := ms.k.Counters.Get(ctx, msg.Sender)
// 	if err != nil && !errors.Is(err, collections.ErrNotFound) {
// 		return nil, err
// 	}

// 	counter++

// 	if err := ms.k.Counter.Set(ctx, msg.Sender, counter); err != nil {
// 		return nil, err
// 	}

// 	return &election.MsgIncrementCounterResponse{}, nil
// }

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

func NewElection(ctx context.Context, msg *election.MsgNewElection) (*election.MsgNewElectionResponse, error) {

	// check if the election is expired
	// create a crs.NewDecision object
	// store the new decision
	// resolve the election
	// delete the election

	if msg == nil {
		return nil, fmt.Errorf("election is required")
	}

	if msg.

	// for _, value := range msg.Options {
	// 	if value == "" {
	// 		return nil, fmt.Errorf("voting options are required")
	// 	}
	// }

	// election := election.NewDecision{
	// 	VotingOptions: msg.Options,
	// }

}
