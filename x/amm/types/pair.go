package types

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	ReserveAddressNamePrefix = "reserve/"
	ShareDenomPrefix = "amm/share/"
)

var MinInitialLiquidity = sdk.NewInt(1000)

func NewPair(id uint64, denom0, denom1 string) Pair {
	return Pair{
		Id: id,
		Denom0: denom0,
		Denom1: denom1,
	}
}

func PairReserveAddress(pair Pair) sdk.AccAddress {
	return address.Module(
		ModuleName, []byte(ReserveAddressNamePrefix + strconv.FormatUint(pair.Id,10)))
}

func ShareDenom(pair Pair) string {
	return ShareDenomPrefix + strconv.FormatUint(pair.Id, 10)
}

func ParseShareDenom(denom string) (pairID uint64, err error) {
	if !strings.HasPrefix(denom, ShareDenomPrefix) {
		return 0, fmt.Errorf("share denom must have %s as prefix", ShareDenomPrefix)
	}
	return strconv.ParseUint(strings.TrimPrefix(denom, ShareDenomPrefix), 10, 64)
}

func SortDenoms(denomA, denomB string) (denom0, denom1 string) {
	if denomA < denomB {
		return denomA, denomB
	}
	return denomB, denomA
}