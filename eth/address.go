package eth

import (
	"encoding/hex"
	"github.com/zlabwork/go-chain"
)

type Address chain.EthAddress

func NewAddress(b []byte) *Address {
	var a Address
	if len(b) > len(a) {
		b = b[len(b)-20:]
	}
	copy(a[20-len(b):], b)
	return &a
}

func (addr *Address) String() string {
	return "0x" + hex.EncodeToString(addr[:])
}
