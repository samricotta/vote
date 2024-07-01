package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/samricotta/vote/x/crs"
	expectedkeepers "github.com/samricotta/vote/x/crs/expected_keepers"
	"github.com/samricotta/vote/x/crs/keeper"
)

type testFixture struct {
	ctx         sdk.Context
	k           keeper.Keeper
	msgServer   crs.MsgServer
	queryServer crs.QueryServer

	addrs []sdk.AccAddress
}

func initFixture(t *testing.T, bankkeeper expectedkeepers.BankKeeper) *testFixture {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	key := storetypes.NewKVStoreKey(crs.ModuleName)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	storeService := runtime.NewKVStoreService(key)
	addrs := simtestutil.CreateIncrementalAccounts(3)

	k := keeper.NewKeeper(encCfg.Codec, addresscodec.NewBech32Codec("cosmos"), storeService, bankkeeper, addrs[0].String())
	err := k.InitGenesis(testCtx.Ctx, crs.NewGenesisState())
	if err != nil {
		panic(err)
	}

	return &testFixture{
		ctx:         testCtx.Ctx,
		k:           k,
		msgServer:   keeper.NewMsgServerImpl(k),
		queryServer: keeper.NewQueryServerImpl(k),
		addrs:       addrs,
	}
}
