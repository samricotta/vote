package keeper_test

import (
	"testing"

	"github.com/samricotta/vote/x/election"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	fixture := initFixture(t)

	data := &election.GenesisState{
		// Counters: []election.Counter{
		// 	{
		// 		Address: fixture.addrs[0].String(),
		// 		Count:   5,
		// 	},
		// },
		Params: election.DefaultParams(),
	}
	err := fixture.k.InitGenesis(fixture.ctx, data)
	require.NoError(t, err)

	params, err := fixture.k.Params.Get(fixture.ctx)
	require.NoError(t, err)

	require.Equal(t, election.DefaultParams(), params)

	// count, err := fixture.k.Counter.Get(fixture.ctx, fixture.addrs[0].String())
	require.NoError(t, err)
	// require.Equal(t, uint64(5), count)
}

// func TestExportGenesis(t *testing.T) {
// 	fixture := initFixture(t)

// 	// _, err := fixture.msgServer.IncrementCounter(fixture.ctx, &election.MsgIncrementCounter{
// 	// 	Sender: fixture.addrs[0].String(),
// 	// })
// 	// require.NoError(t, err)

// 	// out, err := fixture.k.ExportGenesis(fixture.ctx)
// 	require.NoError(t, err)

// 	require.Equal(t, election.DefaultParams(), out.Params)
// 	require.Equal(t, uint64(1), out.Counters[0].Count)
// }
