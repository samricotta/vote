package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	electionv1 "github.com/samricotta/vote/x/election/api/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: electionv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Counter",
					Use:       "counter [address]",
					Short:     "Get the current value of the counter for an address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Get the current module parameters",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: electionv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "IncrementCounter",
					Use:       "counter [sender]",
					Short:     "Increments the counter by 1 for the sender",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "sender"},
					},
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // This is a authority gated tx, so we skip it.
				},
			},
		},
	}
}
