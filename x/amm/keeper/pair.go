package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/pumpkinzomb/zomb-amm/x/amm/types"
)

func (k Keeper) AddLiquidity(ctx sdk.Context, fromAddr sdk.AccAddress, coins sdk.Coins) (mintedShare sdk.Coin, err error) {
	coin0, coin1 := coins[0], coins[1]
	denom0, denom1 := coin0.Denom, coin1.Denom

	pair, found := k.GetPairByDenoms(ctx, denom0, denom1)
	if !found {
		pairID := k.GetLastPairID(ctx)
		pairID++
		k.SetLastPairID(ctx, pairID)

		pair = types.NewPair(pairID, denom0, denom1)
		k.SetPair(ctx, pair)
		k.SetPairIndex(ctx, pair)
	}

	reserveAddr := types.PairReserveAddress(pair)
	shareDenom := types.ShareDenom(pair)

	reserveBalances := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx := reserveBalances.AmountOf(denom0)
	ry := reserveBalances.AmountOf(denom1)
	x := coin0.Amount
	y := coin1.Amount
	totalShare := k.bankKeeper.GetSupply(ctx, shareDenom).Amount

	var ax, ay, share sdk.Int
	if totalShare.IsZero() {
		var l sdk.Dec
		l, err = sdk.NewDecFromInt(x.Mul(y)).ApproxSqrt()
		if err != nil {
			return
		}
		share = l.TruncateInt()
		if share.LT(k.GetMinInitialLiquidity(ctx)) {
			err = sdkerrors.Wrapf(
				types.ErrInsufficientLiquidity, "insufficient initial liquidity: %s", share)
			return
		}
		ax = x
		ay = y
	} else {
		share = sdk.MinInt(totalShare.Mul(x).Quo(rx), totalShare.Mul(y).Quo(ry))
		ax = rx.Mul(share).Quo(totalShare)
		ay = ry.Mul(share).Quo(totalShare)
	}
	if !ax.IsPositive() || !ay.IsPositive() || !share.IsPositive() {
		err = types.ErrInsufficientLiquidity
		return
	}
	
	amt := sdk.NewCoins(sdk.NewCoin(denom0, ax), sdk.NewCoin(denom1, ay))
	err = k.bankKeeper.SendCoins(ctx, fromAddr, reserveAddr, amt)
	if err != nil {
		return 
	}
	mintedShare = sdk.NewCoin(shareDenom, share)
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedShare))
	if err != nil {
		return
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, fromAddr, sdk.NewCoins(mintedShare))
	if err != nil {
		return
	}
	return mintedShare, nil
}

func (k Keeper) RemoveLiquidity(ctx sdk.Context, fromAddr sdk.AccAddress, share sdk.Coin) (withdrawnCoins sdk.Coins, err error) {
	var pairID uint64
	pairID, err = types.ParseShareDenom(share.Denom)
	if err != nil {
		return
	}

	pair, found := k.GetPair(ctx, pairID)
	if !found {
		err = sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", pairID)
		return
	}

	reserveAddr := types.PairReserveAddress(pair)
	reserveBalances := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx := reserveBalances.AmountOf(pair.Denom0)
	ry := reserveBalances.AmountOf(pair.Denom1)
	totalShare := k.bankKeeper.GetSupply(ctx, share.Denom).Amount

	var wx, wy sdk.Int
	if share.Amount.Equal(totalShare){
		wx = rx
		wy = ry
		k.DeletePair(ctx, pair)
		k.DeletePairIndex(ctx, pair)
	} else {
		wx = rx.Mul(share.Amount).Quo(totalShare)
		wy = ry.Mul(share.Amount).Quo(totalShare)
	}
	if !wx.IsPositive() && !wy.IsPositive() {
		err = sdkerrors.Wrap(types.ErrInsufficientLiquidity, "too small share to remove")
		return
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, fromAddr, types.ModuleName, sdk.NewCoins(share))
	if err != nil {
		return
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(share))
	if err != nil {
		return
	}
	withdrawnCoins = sdk.NewCoins(sdk.NewCoin(pair.Denom0, wx), sdk.NewCoin(pair.Denom1, wy))
	err = k.bankKeeper.SendCoins(ctx, reserveAddr, fromAddr, withdrawnCoins)
	if err != nil {
		return 
	}
	return withdrawnCoins, nil
}

func (k Keeper) SwapExactIn(ctx sdk.Context, fromAddr sdk.AccAddress, coinIn, minCoinOut sdk.Coin) (coinOut sdk.Coin, err error) {
	pair, found := k.GetPairByDenoms(ctx, coinIn.Denom, minCoinOut.Denom)
	if !found {
		err = sdkerrors.Wrap(sdkerrors.ErrNotFound, "pair not found")
		return
	}

	reserveAddr := types.PairReserveAddress(pair)
	reserveBalances := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx := reserveBalances.AmountOf(pair.Denom0)
	ry := reserveBalances.AmountOf(pair.Denom1)
	feeRate := k.GetFeeRate(ctx)

	var reserveIn, reserveOut sdk.Int
	amtInWithoutFee := sdk.NewDecFromInt(coinIn.Amount).MulTruncate(sdk.OneDec().Sub(feeRate)).TruncateInt()
	if coinIn.Denom == pair.Denom0 {
		reserveIn, reserveOut = rx, ry
		coinOut.Denom = pair.Denom1
	} else {
		reserveIn, reserveOut = ry, rx
		coinOut.Denom = pair.Denom0
	}
	coinOut.Amount = reserveOut.Mul(amtInWithoutFee).Quo(reserveIn.Add(amtInWithoutFee))
	if coinOut.Amount.LT(minCoinOut.Amount) {
		err = sdkerrors.Wrapf(
			types.ErrSmallOutCoin, "%s is smaller than %s", coinOut.Amount, minCoinOut.Amount)
		return
	}

	err = k.bankKeeper.SendCoins(ctx, fromAddr, reserveAddr, sdk.NewCoins(coinIn))
	if err != nil {
		return 
	}
	err = k.bankKeeper.SendCoins(ctx, reserveAddr, fromAddr, sdk.NewCoins(coinOut))
	if err != nil {
		return 
	}
	return coinOut, nil
}

func (k Keeper) SwapExactOut(ctx sdk.Context, fromAddr sdk.AccAddress, coinOut, maxCoinIn sdk.Coin) (coinIn sdk.Coin, err error) {
	pair, found := k.GetPairByDenoms(ctx, maxCoinIn.Denom, coinOut.Denom)
	if !found {
		err = sdkerrors.Wrap(sdkerrors.ErrNotFound, "pair not found")
		return
	}
	
	reserveAddr := types.PairReserveAddress(pair)
	reserveBalances := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx := reserveBalances.AmountOf(pair.Denom0)
	ry := reserveBalances.AmountOf(pair.Denom1)
	feeRate := k.GetFeeRate(ctx)

	var reserveIn, reserveOut sdk.Int

	if coinOut.Denom == pair.Denom1 { // x to y
		reserveIn, reserveOut = rx, ry
		coinIn.Denom = pair.Denom0
	} else { // y to x
		reserveIn, reserveOut = ry, rx
		coinIn.Denom = pair.Denom1
	}
	coinIn.Amount = sdk.NewDecFromInt(reserveIn.Mul(coinOut.Amount)).QuoInt(reserveOut.Sub(coinOut.Amount)).Mul(sdk.OneDec().Add(feeRate)).Ceil().TruncateInt()
	if coinIn.Amount.GT(maxCoinIn.Amount) {
		err = sdkerrors.Wrapf(
			types.ErrBigInCoin, "%s is bigger than %s", coinIn.Amount, maxCoinIn.Amount)
		return
	}
	if err = k.bankKeeper.SendCoins(ctx, fromAddr, reserveAddr, sdk.NewCoins(coinIn)); err != nil {
		return
	}
	if err = k.bankKeeper.SendCoins(ctx, reserveAddr, fromAddr, sdk.NewCoins(coinOut)); err != nil {
		return
	}
	return coinIn, nil
}