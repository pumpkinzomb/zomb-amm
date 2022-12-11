package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddLiquidity{}, "zomb/MsgAddLiquidity", nil)
	cdc.RegisterConcrete(&MsgRemoveLiquidity{}, "zomb/MsgRemoveLiquidity", nil)
	cdc.RegisterConcrete(&MsgSwapExactIn{}, "zomb/MsgSwapExactIn", nil)
	cdc.RegisterConcrete(&MsgSwapExactOut{}, "zomb/MsgSwapExactOut", nil)
}

var (
	amino = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry){
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
	
	RegisterLegacyAminoCodec(authzcodec.Amino)
}