package btc

import (
    "crypto/rand"
    "fmt"
    "github.com/mr-tron/base58"
    "github.com/zlabwork/gochain/utils"
)

type TestNetPriKey struct {
    k []byte
}

func TestNewPriKeyRandom() *TestNetPriKey {
    b := make([]byte, 32)
    rand.Read(b)
    return &TestNetPriKey{k: b}
}

func TestNewPriKey(b []byte) (*TestNetPriKey, error) {
    if len(b) != 32 {
        return nil, fmt.Errorf("%s", "invalid length")
    }
    return &TestNetPriKey{k: b}, nil
}

type TestNetAddress struct {
    pub  []byte
    hash []byte
}

func TestNewAddress(b []byte) *TestNetAddress {
    h, _ := utils.Hash160(b)
    return &TestNetAddress{
        pub:  b,
        hash: h,
    }
}

func (pri *TestNetPriKey) WIF() string {
    ver := []byte{0xef}
    compression := byte(0x01)
    k := append(ver, pri.k...)
    k = append(k, compression)
    sum, err := checksum(k)
    if err != nil {
        return ""
    }
    k = append(k, sum...)
    return base58.Encode(k)
}

func (addr *TestNetAddress) P2PKH() string {
    prefix := []byte{0x6f}
    preData := append(prefix, addr.hash...)
    sum, _ := checksum(preData)
    address := append(preData, sum...)
    return base58.Encode(address)
}

func (addr *TestNetAddress) P2SH() string {
    prefix := []byte{0xc4}
    preData := append(prefix, addr.hash...)
    sum, _ := checksum(preData)
    address := append(preData, sum...)
    return base58.Encode(address)
}
