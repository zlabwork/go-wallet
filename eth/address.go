package eth

import "encoding/hex"

type Address [20]byte

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
