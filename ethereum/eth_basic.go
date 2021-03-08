package ethereum

import (
    "errors"
    "github.com/ethereum/go-ethereum/crypto"
)

type priKeyData struct {
    key []byte
}

type pubKeyData struct {
    key []byte
}

// 创建账号
func NewPriKeyRandom() *priKeyData {
    key, err := crypto.GenerateKey()
    if err != nil {
        return nil
    }

    return &priKeyData{key: key.D.Bytes()}
}

func NewPriKey(priKey []byte) (*priKeyData, error) {
    if len(priKey) != 32 {
        return nil, errors.New("invalid length")
    }
    return &priKeyData{key: priKey}, nil
}

func (pri *priKeyData) Key() []byte {
    return pri.key
}

func (pri *priKeyData) Address() (string, error) {
    priKey, err := crypto.ToECDSA(pri.key)
    if err != nil {
        return "", err
    }

    addr := crypto.PubkeyToAddress(priKey.PublicKey)
    return addr.Hex(), nil
}
