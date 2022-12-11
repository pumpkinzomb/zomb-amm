package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pumpkinzomb/zomb-amm/x/amm/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) AddLiquidity(c context.Context, msg *types.MsgAddLiquidity) (* types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	mintedShare, err := k.Keeper.AddLiquidity(
		ctx, sdk.MustAccAddressFromBech32(msg.Sender), msg.Coins)

	if err != nil {
		return nil, err
	}
	return &types.MsgAddLiquidityResponse{
		MintedShare: mintedShare,
	}, nil
}

func (k msgServer) RemoveLiquidity(c context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	withdrawnCoins, err := k.Keeper.RemoveLiquidity(
		ctx, sdk.MustAccAddressFromBech32(msg.Sender), msg.Share)
	if err != nil {
		return nil, err
	}
	return &types.MsgRemoveLiquidityResponse{
		WithdrawnCoins: withdrawnCoins,
	}, nil
}

func (k msgServer) SwapExactIn(c context.Context, msg *types.MsgSwapExactIn) (*types.MsgSwapExactInResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	coinOut, err := k.Keeper.SwapExactIn(
		ctx, sdk.MustAccAddressFromBech32(msg.Sender), msg.CoinIn, msg.MinCoinOut)
	if err != nil {
		return nil, err
	}
	return &types.MsgSwapExactInResponse{
		CoinOut: coinOut,
	}, nil
}

func (k msgServer) SwapExactOut(c context.Context, msg *types.MsgSwapExactOut) (*types.MsgSwapExactOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	coinIn, err := k.Keeper.SwapExactOut(
		ctx, sdk.MustAccAddressFromBech32(msg.Sender), msg.CoinOut, msg.MaxCoinIn)
	if err != nil {
		return nil, err
	}
	return &types.MsgSwapExactOutResponse{
		CoinIn: coinIn,
	}, nil
}