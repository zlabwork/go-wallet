package bitcoin

import (
    "github.com/FactomProject/btcutilecc"
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
    curve                     = btcutil.Secp256k1()
    publicKeyCompressedLength = 33
    hashConst, _              = hashSha256([]byte("zlab")) // hash 常量
)
