package btc

import (
    btcutil "github.com/FactomProject/btcutilecc"
    "github.com/mr-tron/base58"
    "github.com/zlabwork/gochain/utils"
)

const (
    OP_0             = byte(0x00)
    OP_1             = byte(0x51)
    OP_PUSH_20       = byte(0x14)
    OP_PUSH_32       = byte(0x20)
    OP_PUSH_33       = byte(0x21)
    OP_CHECKSIG      = byte(0xAC)
    OP_CHECKMULTISIG = byte(0xAE)
)

var (
    curve               = btcutil.Secp256k1()
    pubKeyCompressedLen = 33
    hashConst, _        = utils.HashSha256([]byte("zlab")) // hash 常量
)

func checksum(data []byte) ([]byte, error) {
    hash, err := utils.HashDoubleSha256(data)
    if err != nil {
        return nil, err
    }
    return hash[:4], nil
}

func p2sh(hash160 []byte) string {
    data := append([]byte{0x05}, hash160...)
    sum, _ := checksum(data)
    return base58.Encode(append(data, sum...))
}

func p2pkh(hash160 []byte) string {
    data := append([]byte{0x00}, hash160...)
    sum, _ := checksum(data)
    return base58.Encode(append(data, sum...))
}
