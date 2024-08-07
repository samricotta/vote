package keeper_test

import (
	"testing"

	"github.com/samricotta/vote/x/crs"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	fixture := initFixture(t, nil)

	data := &crs.GenesisState{
		Params: crs.DefaultParams(),
	}
	err := fixture.k.InitGenesis(fixture.ctx, data)
	require.NoError(t, err)

	params, err := fixture.k.Params.Get(fixture.ctx)
	require.NoError(t, err)

	require.Equal(t, crs.DefaultParams(), params)

	// count, err := fixture.k.Counter.Get(fixture.ctx, fixture.addrs[0].String())
	// require.NoError(t, err)
	// require.Equal(t, uint64(5), count)
}

func TestExportGenesis(t *testing.T) {
	fixture := initFixture(t, nil)

	// _, err := fixture.msgServer.IncrementCounter(fixture.ctx, &crs.MsgIncrementCounter{
	// 	Sender: fixture.addrs[0].String(),
	// })
	// require.NoError(t, err)

	out, err := fixture.k.ExportGenesis(fixture.ctx)
	require.NoError(t, err)

	require.Equal(t, crs.DefaultParams(), out.Params)
	// require.Equal(t, uint64(1), out.Counters[0].Count)
}
