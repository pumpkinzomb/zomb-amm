package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrWrongCoinNumber       = sdkerrors.Register(ModuleName, 2, "wrong number of coins")
	ErrInsufficientLiquidity = sdkerrors.Register(ModuleName, 3, "insufficient liquidity")
	ErrSmallOutCoin          = sdkerrors.Register(ModuleName, 4, "calculated out coin is smaller than the minimum")
	ErrBigInCoin             = sdkerrors.Register(ModuleName, 5, "calculated in coin is bigger than the maximum")
)