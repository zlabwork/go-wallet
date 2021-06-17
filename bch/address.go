package bch

import (
    "github.com/zlabwork/gochain/utils"
    "golang.org/x/crypto/ripemd160"
)

type Address struct {
    pub  []byte
    hash []byte
    pre  string // bchreg、bchtest、bchsim
}

func NewAddress(b []byte) *Address {
    h, _ := utils.Hash160(b)
    return &Address{
        pub:  b,
        hash: h[:ripemd160.Size],
        pre:  "bitcoincash",
    }
}

func (ad *Address) P2PKH() string {
    return createCashAddress(ad.hash, ad.pre, addrTypeP2PKH)
}

func (ad *Address) P2SH() string {
    return createCashAddress(ad.hash, ad.pre, addrTypeP2SH)
}
