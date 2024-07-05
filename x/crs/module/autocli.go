package module

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/samricotta/vote/x/crs"
	crsv1 "github.com/samricotta/vote/x/crs/api/v1"
	"github.com/samricotta/vote/x/crs/utils"
	"github.com/spf13/cobra"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: crsv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Decision",
					Use:       "decision [decision_id]",
					Short:     "Get a decision by ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: crsv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CreateDecision",
					Use:       "create-decision [commit_duration] [reveal_duration] [entry_fee] [refund] [send_funds_to] [options]",
					Short:     "Create a new decision",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "commit_duration"},
						{ProtoField: "reveal_duration"},
						{ProtoField: "entry_fee"},
						{ProtoField: "refund"},
						{ProtoField: "send_funds_to"},
						{ProtoField: "options", Varargs: true},
					},
				},
				{
					RpcMethod: "Commit",
					Use:       "commit [decision_id] [commit_hash]",
					Short:     "Commit to a decision",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "decision_id"},
						{ProtoField: "commit"},
					},
				},
				{
					RpcMethod: "Reveal",
					Use:       "reveal [decision_id] [option_chosen] [salt]",
					Short:     "Reveal the choice after a commit",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "decision_id"},
						{ProtoField: "option_chosen"},
						{ProtoField: "salt"},
					},
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // This is a authority gated tx, so we skip it.
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}

// For the custom "commit" command
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   crs.ModuleName,
		Short: "Transactions commands for the crs module",
		RunE:  client.ValidateCmd,
	}

	//   minid tx crs commit [decision_id] [commit_hash] [flags]
	commitCmd := &cobra.Command{
		Use:   "commit [decision_id] [option]",
		Short: "Commit to a decision",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var option []byte

			// accept options in hex and base64
			if data, err := hex.DecodeString(args[1]); err == nil {
				option = data
			} else if data, err := base64.StdEncoding.DecodeString(args[1]); err == nil {
				option = data
			}

			if len(option) == 0 {
				return errors.New("option must be a valid hex or base64 string")
			}

			decisionID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			// salt is 32 random bytes
			salt := make([]byte, 32)
			_, err = rand.Read(salt)
			if err != nil {
				return err
			}

			cmd.Printf("Salt: %x \n", salt)

			hash, err := utils.CalculateCommit(decisionID, option, salt)
			if err != nil {
				return err
			}
			msg := &crs.MsgCommit{
				Sender:     clientCtx.GetFromAddress().String(),
				DecisionId: decisionID,
				Commit:     hash,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(commitCmd)
	txCmd.AddCommand(commitCmd)
	return txCmd
}
