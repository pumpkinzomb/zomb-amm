package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddLiquidity{},
		&MsgRemoveLiquidity{},
		&MsgSwapExactIn{},
		&MsgSwapExactOut{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddLiquidity{}, "zomb/MsgAddLiquidity", nil)
	cdc.RegisterConcrete(&MsgRemoveLiquidity{}, "zomb/MsgRemoveLiquidity", nil)
	cdc.RegisterConcrete(&MsgSwapExactIn{}, "zomb/MsgSwapExactIn", nil)
	cdc.RegisterConcrete(&MsgSwapExactOut{}, "zomb/MsgSwapExactOut", nil)
}