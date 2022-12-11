package utils

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func LengthPrefix(b []byte) []byte {
	l := len(b)
	if l == 0 {
		return b
	}
	return append([]byte{byte(l)}, b...)
}

func SampleAddress(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}