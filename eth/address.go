package eth

import (
	"github.com/ethereum/go-ethereum/common"
)

type Address struct {
	addr common.Address
}

func NewAddress(b []byte) *Address {
	return &Address{
		addr: common.BytesToAddress(b),
	}
}

func (ad *Address) String() string {
	return ad.addr.Hex()
}
