package eth

import (
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "github.com/ethereum/go-ethereum/crypto"
)

type PriKey struct {
    k []byte
}

func NewPriKeyRandom() *PriKey {
    key, err := crypto.GenerateKey()
    if err != nil {
        return nil
    }
    return &PriKey{k: key.D.Bytes()}
}

func NewPriKey(b []byte) (*PriKey, error) {
    if len(b) != 32 {
        return nil, fmt.Errorf("%s", "invalid length")
    }
    return &PriKey{k: b}, nil
}

func (pri *PriKey) Bytes() []byte {
    return pri.k
}

func (pri *PriKey) Base64() string {
    return base64.StdEncoding.EncodeToString(pri.Bytes())
}

func (pri *PriKey) Hex() string {
    return hex.EncodeToString(pri.Bytes())
}

func (pri *PriKey) Address() *Address {
    priKey, err := crypto.ToECDSA(pri.k)
    if err != nil {
        return nil
    }
    addr := crypto.PubkeyToAddress(priKey.PublicKey)
    return &Address{addr: addr}
}
