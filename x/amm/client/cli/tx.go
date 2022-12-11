package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pumpkinzomb/zomb-amm/x/amm/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetAddLiquidityCmd(),
		GetRemovedLiquidityCmd(),
		GetSwapExactInCmd(),
		GetSwapExactOutCmd(),
	)

	return cmd
}

func GetAddLiquidityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidity [coins]",
		Args:  cobra.ExactArgs(1),
		Short: "Add liquidity to a pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return fmt.Errorf("invalid coins: %w", err)
			}
			msg := types.NewMsgAddLiquidity(clientCtx.GetFromAddress(), coins)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetRemovedLiquidityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity [share]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove liquidity from a pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			share, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("invalid share: %w", err)
			}
			msg := types.NewMsgRemoveLiquidity(clientCtx.GetFromAddress(), share)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetSwapExactInCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-exact-in [coin-in] [min-coin-out]",
		Args:  cobra.ExactArgs(2),
		Short: "Swap exact amount of input coin for output coin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coinIn, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("invalid coin in: %w", err)
			}
			minCoinOut, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid min coin out: %w", err)
			}
			msg := types.NewMsgSwapExactIn(clientCtx.GetFromAddress(), coinIn, minCoinOut)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetSwapExactOutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-exact-out [coin-out] [max-coin-in]",
		Args:  cobra.ExactArgs(2),
		Short: "Swap input coin for exact amount of output coin",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coinOut, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("invalid coin out: %w", err)
			}
			maxCoinIn, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid max coin in: %w", err)
			}
			msg := types.NewMsgSwapExactOut(clientCtx.GetFromAddress(), coinOut, maxCoinIn)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}