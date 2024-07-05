package keeper

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/samricotta/vote/x/crs"
	"github.com/samricotta/vote/x/crs/utils"
)

type msgServer struct {
	k Keeper
}

var _ crs.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) crs.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) CreateDecision(ctx context.Context, msg *crs.MsgCreateDecision) (*crs.MsgCreateDecisionResponse, error) {

	// check sender address
	senderAddr, err := ms.k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, fmt.Errorf("invalid sender address: %w", err)
	}

	// check if commit and reveal times are in the future
	if msg.CommitDuration <= 0 {
		return nil, fmt.Errorf("commit duration must be greater than 0")
	}

	if msg.RevealDuration <= 0 {
		return nil, fmt.Errorf("reveal duration must be greater than 0")
	}

	// check if entry fee is nil
	if !msg.EntryFee.IsNil() {
		err = ms.k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, crs.ModuleName, sdk.NewCoins(msg.EntryFee))
		if err != nil {
			return nil, err
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// store decision
	commitTimeout := sdkCtx.BlockTime().Add(time.Second * time.Duration(msg.CommitDuration))
	revealTimeout := commitTimeout.Add(time.Second * time.Duration(msg.RevealDuration))
	newDecision := crs.Decision{
		EntryFee:      msg.EntryFee,
		Options:       msg.Options,
		CommitTimeout: commitTimeout,
		RevealTimeout: revealTimeout,
	}

	if err = ms.k.CreateDecision(ctx, newDecision); err != nil {
		return nil, err
	}

	return &crs.MsgCreateDecisionResponse{}, nil
}

// Commit checks if the decision is still open for committing, then stores the commit
func (ms msgServer) Commit(ctx context.Context, msg *crs.MsgCommit) (*crs.MsgCommitResponse, error) {
	// convert sender address from string to bytes
	senderAddr, err := ms.k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, fmt.Errorf("invalid sender address: %w", err)
	}

	decision, err := ms.k.Decisions.Get(ctx, msg.DecisionId)
	if err != nil {
		return nil, err
	}

	// check if the commit period has ended
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if sdkCtx.BlockTime().After(decision.CommitTimeout) {
		return nil, fmt.Errorf("commit period has ended")
	}

	// check if a commit already exists for this decision and sender
	_, err = ms.k.Commits.Get(ctx, collections.Join(msg.DecisionId, senderAddr))
	if err == nil {
		return nil, fmt.Errorf("commit already exists")
	}

	// store commit
	commit := crs.Commit{
		Commit: msg.Commit,
	}

	err = ms.k.Commits.Set(ctx, collections.Join(msg.DecisionId, senderAddr), commit)
	if err != nil {
		return nil, err
	}

	return &crs.MsgCommitResponse{}, nil
}

// Reveal checks if the decision is still open for revealing, checks the reveal against the commit, then stores the reveal
func (ms msgServer) Reveal(ctx context.Context, msg *crs.MsgReveal) (*crs.MsgRevealResponse, error) {
	// convert sender address from string to bytes
	senderAddr, err := ms.k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, fmt.Errorf("invalid sender address: %w", err)
	}

	decision, err := ms.k.Decisions.Get(ctx, msg.DecisionId)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if sdkCtx.BlockTime().Before(decision.CommitTimeout) {
		return nil, fmt.Errorf("commit period has not ended")
	}

	if sdkCtx.BlockTime().After(decision.RevealTimeout) {
		return nil, fmt.Errorf("reveal period has ended")
	}

	// check if a reveal already exists for this decision and sender
	_, err = ms.k.Reveals.Get(ctx, collections.Join(msg.DecisionId, senderAddr))
	if err == nil {
		return nil, fmt.Errorf("reveal already exists")
	}

	// fetch commit
	commit, err := ms.k.Commits.Get(ctx, collections.Join(msg.DecisionId, senderAddr))
	if err != nil {
		return nil, err
	}

	// check if the reveal matches the commit, by recalculating the sha256 hash
	hash, err := utils.CalculateCommit(msg.DecisionId, msg.OptionChosen, msg.Salt)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(hash, commit.Commit) {
		return nil, fmt.Errorf("reveal does not match commit")
	}

	// store reveal
	reveal := crs.Reveal{
		Option: msg.OptionChosen,
	}

	err = ms.k.Reveals.Set(ctx, collections.Join(msg.DecisionId, senderAddr), reveal)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateParams params is defining the handler for the MsgUpdateParams message.
func (ms msgServer) UpdateParams(ctx context.Context, msg *crs.MsgUpdateParams) (*crs.MsgUpdateParamsResponse, error) {
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

	return &crs.MsgUpdateParamsResponse{}, nil
}
