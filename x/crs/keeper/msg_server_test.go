package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/collections"
	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/samricotta/vote/x/crs"
	"github.com/stretchr/testify/require"

	"github.com/samricotta/vote/x/crs/keeper"
	"github.com/samricotta/vote/x/crs/mocks"
)

func TestUpdateParams(t *testing.T) {
	f := initFixture(t, nil)
	require := require.New(t)

	testCases := []struct {
		name         string
		request      *crs.MsgUpdateParams
		expectErrMsg string
	}{
		{
			name: "set invalid authority (not an address)",
			request: &crs.MsgUpdateParams{
				Authority: "foo",
			},
			expectErrMsg: "invalid authority address",
		},
		{
			name: "set invalid authority (not defined authority)",
			request: &crs.MsgUpdateParams{
				Authority: f.addrs[1].String(),
			},
			expectErrMsg: fmt.Sprintf("unauthorized, authority does not match the module's authority: got %s, want %s", f.addrs[1].String(), f.k.GetAuthority()),
		},
		{
			name: "set valid params",
			request: &crs.MsgUpdateParams{
				Authority: f.k.GetAuthority(),
				Params:    crs.Params{},
			},
			expectErrMsg: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := f.msgServer.UpdateParams(f.ctx, tc.request)
			if tc.expectErrMsg != "" {
				require.Error(err)
				require.ErrorContains(err, tc.expectErrMsg)
			} else {
				require.NoError(err)
			}
		})
	}
}

func TestCreateDecision(t *testing.T) {
	mockctrl := gomock.NewController(t)
	defer mockctrl.Finish()

	bankkeeper := mocks.NewMockBankKeeper(mockctrl)

	f := initFixture(t, bankkeeper)
	require := require.New(t)

	testCases := []struct {
		name         string
		request      *crs.MsgCreateDecision
		expectErrMsg string
		malleate     func()
	}{
		{
			name: "set invalid sender (not an address)",
			request: &crs.MsgCreateDecision{
				Sender: "foo",
			},
			expectErrMsg: "invalid sender address",
		},
		{
			name: "set invalid decision (commit duration is 0)",
			request: &crs.MsgCreateDecision{
				Sender:         f.addrs[0].String(),
				CommitDuration: 0,
				RevealDuration: 100,
			},
			expectErrMsg: "commit duration must be greater than 0",
		},
		{
			name: "set invalid decision (reveal duration is 0)",
			request: &crs.MsgCreateDecision{
				Sender:         f.addrs[0].String(),
				CommitDuration: 100,
				RevealDuration: 0,
			},
			expectErrMsg: "reveal duration must be greater than 0",
		},
		{
			name: "set valid decision",
			request: &crs.MsgCreateDecision{
				Sender:         f.addrs[0].String(),
				CommitDuration: 100,
				RevealDuration: 100,
				EntryFee:       sdk.NewInt64Coin("stake", 100),
			},
			expectErrMsg: "",
			malleate: func() {
				bankkeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), f.addrs[0], crs.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("stake", 100))).Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.malleate != nil {
				tc.malleate()
			}
			_, err := f.msgServer.CreateDecision(f.ctx, tc.request)
			if tc.expectErrMsg != "" {
				require.Error(err)
				require.ErrorContains(err, tc.expectErrMsg)
			} else {
				require.NoError(err)
			}
		})
	}

}

func TestCommit(t *testing.T) {
	f := initFixture(t, nil)
	require := require.New(t)

	testCases := []struct {
		name         string
		request      *crs.MsgCommit
		expectErrMsg string
		malleate     func()
	}{
		{
			name: "set invalid sender (not an address)",
			request: &crs.MsgCommit{
				Sender: "foo",
			},
			expectErrMsg: "invalid sender address",
		},
		{
			name: "set invalid decision ID (does not exist)",
			request: &crs.MsgCommit{
				Sender:     f.addrs[0].String(),
				DecisionId: 123,
			},
			expectErrMsg: "collections: not found: key '123' of type github.com/cosmos/gogoproto/samricotta.crs.v1.Decision",
		},
		{
			name: "commit period is over",
			request: &crs.MsgCommit{
				Sender:     f.addrs[0].String(),
				DecisionId: 2,
				Commit:     []byte("commit"),
			},
			expectErrMsg: "commit period has ended",
			malleate: func() {
				f.k.Decision.Set(f.ctx, 2, crs.Decision{
					CommitTimeout: f.ctx.BlockTime().Add(-2), // expired
				})
			},
		},
		{
			name: "set \"valid\" commit",
			request: &crs.MsgCommit{
				Sender:     f.addrs[0].String(),
				DecisionId: 1,
				Commit:     []byte("commit"),
			},
			expectErrMsg: "",
			malleate: func() {
				f.k.Decision.Set(f.ctx, 1, crs.Decision{
					CommitTimeout: f.ctx.BlockTime().Add(10),
				})
			},
		},
		{
			name: "commit already exists",
			request: &crs.MsgCommit{
				Sender:     f.addrs[0].String(),
				DecisionId: 1,
				Commit:     []byte("commit"),
			},
			expectErrMsg: "commit already exists",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.malleate != nil {
				tc.malleate()
			}
			_, err := f.msgServer.Commit(f.ctx, tc.request)
			if tc.expectErrMsg != "" {
				require.Error(err)
				require.ErrorContains(err, tc.expectErrMsg)
			} else {
				require.NoError(err)
			}
		})
	}

}

func TestReveal(t *testing.T) {
	f := initFixture(t, nil)
	require := require.New(t)

	salt := rand.Bytes(32)

	testCases := []struct {
		name         string
		request      *crs.MsgReveal
		expectErrMsg string
		malleate     func()
	}{
		{
			name: "set invalid sender (not an address)",
			request: &crs.MsgReveal{
				Sender: "foo",
			},
			expectErrMsg: "invalid sender address",
		},
		{
			name: "set invalid decision ID (does not exist)",
			request: &crs.MsgReveal{
				Sender:       f.addrs[0].String(),
				DecisionId:   123,
				OptionChosen: []byte("reveal"),
			},
			expectErrMsg: "collections: not found: key '123' of type github.com/cosmos/gogoproto/samricotta.crs.v1.Decision",
		},
		{
			name: "reveal period is over",
			request: &crs.MsgReveal{
				Sender:       f.addrs[0].String(),
				DecisionId:   2,
				OptionChosen: []byte("reveal"),
			},
			expectErrMsg: "reveal period has ended",
			malleate: func() {
				f.k.Decision.Set(f.ctx, 2, crs.Decision{
					RevealTimeout: f.ctx.BlockTime().Add(-2), // expired
				})
			},
		},
		{
			name: "reveal does not match commit",
			request: &crs.MsgReveal{
				Sender:       f.addrs[0].String(),
				DecisionId:   1,
				OptionChosen: []byte("DIFFERENT"),
				Salt:         salt,
			},
			expectErrMsg: "reveal does not match commit",
			malleate: func() {
				f.k.Decision.Set(f.ctx, 1, crs.Decision{
					RevealTimeout: f.ctx.BlockTime().Add(10),
				})

				commit, err := keeper.CalculateCommit(
					1,
					[]byte("ABC"),
					salt,
				)
				require.NoError(err)

				f.k.Commit.Set(f.ctx, collections.Join(uint64(1), f.addrs[0].Bytes()), crs.Commit{
					Commit: commit,
				})

			},
		},
		{
			name: "set \"valid\" reveal",
			request: &crs.MsgReveal{
				Sender:       f.addrs[0].String(),
				DecisionId:   1,
				OptionChosen: []byte("AAA"),
				Salt:         salt,
			},
			expectErrMsg: "",
			malleate: func() {
				f.k.Decision.Set(f.ctx, 1, crs.Decision{
					RevealTimeout: f.ctx.BlockTime().Add(10),
				})

				commit, err := keeper.CalculateCommit(
					1,
					[]byte("AAA"),
					salt,
				)
				require.NoError(err)

				f.k.Commit.Set(f.ctx, collections.Join(uint64(1), f.addrs[0].Bytes()), crs.Commit{
					Commit: commit,
				})

			},
		},
		{
			name: "reveal already exists",
			request: &crs.MsgReveal{
				Sender:       f.addrs[0].String(),
				DecisionId:   1,
				OptionChosen: []byte("reveal"),
			},
			expectErrMsg: "reveal already exists",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.malleate != nil {
				tc.malleate()
			}
			_, err := f.msgServer.Reveal(f.ctx, tc.request)
			if tc.expectErrMsg != "" {
				require.Error(err)
				require.ErrorContains(err, tc.expectErrMsg)
			} else {
				require.NoError(err)
			}
		})
	}

	// f.k.Decision.Set(f.ctx, 1, crs.Decision{
	// 	RevealTimeout: f.ctx.BlockTime().Add(10),
	// })

	// f.k.Decision.Set(f.ctx, 2, crs.Decision{
	// 	RevealTimeout: f.ctx.BlockTime().Add(-2), // expired
	// })

}

// func TestIncrementCounter(t *testing.T) {
// 	f := initFixture(t)
// 	require := require.New(t)

// 	testCases := []struct {
// 		name            string
// 		request         *crs.MsgIncrementCounter
// 		expectErrMsg    string
// 		expectedCounter uint64
// 	}{
// 		{
// 			name: "set invalid sender (not an address)",
// 			request: &crs.MsgIncrementCounter{
// 				Sender: "foo",
// 			},
// 			expectErrMsg: "invalid sender address",
// 		},
// 		{
// 			name: "set valid sender",
// 			request: &crs.MsgIncrementCounter{
// 				Sender: "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
// 			},
// 			expectErrMsg:    "",
// 			expectedCounter: 1,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			_, err := f.msgServer.IncrementCounter(f.ctx, tc.request)
// 			if tc.expectErrMsg != "" {
// 				require.Error(err)
// 				require.ErrorContains(err, tc.expectErrMsg)
// 			} else {
// 				require.NoError(err)

// 				counter, err := f.k.Counter.Get(f.ctx, tc.request.Sender)
// 				require.NoError(err)
// 				require.Equal(tc.expectedCounter, counter)
// 			}
// 		})
// 	}

// }
