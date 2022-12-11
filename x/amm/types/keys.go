package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pumpkinzomb/zomb-amm/utils"
)

const (
	ModuleName = "amm"
	StoreKey = ModuleName
	RouterKey = ModuleName
	MemStoreKey = "meme_amm"
)

var (
	LastPairIDKey = []byte{0x01}
	PairKeyPrefix = []byte{0x02}
	PairIndexKeyPrefix = []byte{0x03}
)

func GetPairKey(pairID uint64) []byte {
	return append(PairKeyPrefix, sdk.Uint64ToBigEndian(pairID)...)
}

func GetPairIndexKey(denom0, denom1 string) []byte {
	return append(
		append(PairIndexKeyPrefix, utils.LengthPrefix([]byte(denom0))...,), denom1...
	)
}