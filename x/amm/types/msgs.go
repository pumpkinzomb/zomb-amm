package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgAddLiquidity) (nil)
	_ sdk.Msg = (*MsgRemoveLiquidity) (nil)
)

const (
	TypeMsgAddLiquidity    = "add_liquidity"
	TypeMsgRemoveLiquidity = "remove_liquidity"
	TypeMsgSwapExactIn     = "swap_exact_in"
	TypeMsgSwapExactOut    = "swap_exact_out"
)

func NewMsgAddLiquidity(sender sdk.AccAddress, coins sdk.Coins) *MsgAddLiquidity {
	return &MsgAddLiquidity{
		Sender: sender.String(),
		Coins: coins,
	}
}

func (msg MsgAddLiquidity) Route() string { return RouterKey }
func (msg MsgAddLiquidity) Type() string { return TypeMsgAddLiquidity }

func (msg MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgAddLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %v", err)
	}
	err = msg.Coins.Validate()
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err.Error())
	}
	if len(msg.Coins) != 2 {
		return ErrWrongCoinNumber
	}
	return nil
}

func NewMsgRemoveLiquidity(sender sdk.AccAddress, share sdk.Coin) *MsgRemoveLiquidity {
	return &MsgRemoveLiquidity{
		Sender: sender.String(),
		Share: share,
	}
}

func (msg MsgRemoveLiquidity) Route() string { return RouterKey }
func (msg MsgRemoveLiquidity) Type() string { return TypeMsgRemoveLiquidity }

func (msg MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgRemoveLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %v", err)
	}
	_, err = ParseShareDenom(msg.Share.Denom)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func NewMsgSwapExactIn(sender sdk.AccAddress, coinIn, minCoinOut sdk.Coin) *MsgSwapExactIn {
	return &MsgSwapExactIn{
		Sender:     sender.String(),
		CoinIn:     coinIn,
		MinCoinOut: minCoinOut,
	}
}

func (msg MsgSwapExactIn) Route() string { return RouterKey }
func (msg MsgSwapExactIn) Type() string  { return TypeMsgSwapExactIn }

func (msg MsgSwapExactIn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgSwapExactIn) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSwapExactIn) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %v", err)
	}
	if err := msg.CoinIn.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin in: %v", err)
	}
	if err := msg.MinCoinOut.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid min coin out: %v", err)
	}
	return nil
}

func NewMsgSwapExactOut(sender sdk.AccAddress, coinOut, maxCoinIn sdk.Coin) *MsgSwapExactOut {
	return &MsgSwapExactOut{
		Sender:    sender.String(),
		CoinOut:   coinOut,
		MaxCoinIn: maxCoinIn,
	}
}

func (msg MsgSwapExactOut) Route() string { return RouterKey }
func (msg MsgSwapExactOut) Type() string  { return TypeMsgSwapExactOut }

func (msg MsgSwapExactOut) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgSwapExactOut) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSwapExactOut) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %v", err)
	}
	if err := msg.CoinOut.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coin out: %v", err)
	}
	if err := msg.MaxCoinIn.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid max coin in: %v", err)
	}
	return nil
}