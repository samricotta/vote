package keeper

import (
	"context"
	"fmt"

	"github.com/samricotta/vote/x/crs"
	"github.com/samricotta/vote/x/vote"
)

type msgServer struct {
	keeper Keeper
}

var _ vote.MsgServer = msgServer{}

func NewMsgServerImp(keeper Keeper) vote.MsgServer {
	return &msgServer{keeper: keeper}
}

func NewVote(ctx context.Context, msg *vote.MsgNewVote) (*vote.MsgNewVoteResponse, error) {
	// check if the vote is empty
	// check if the vote is valid
	// check if the vote is expired
	// check if the vote is resolved
	// create a crs.NewDecision object
	// store the new decision
	// resolve the vote
	// delete the vote

	if msg == nil {
		return nil, fmt.Errorf("vote is required")
	}

	for _, value := range msg.Options {
		if value == "" {
			return nil, fmt.Errorf("voting options are required")
		}
	}

	crs := crs.NewDecision{
		VotingOptions: msg.Options,
		
	}

}
