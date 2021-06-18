package btc

import (
    "fmt"
    "github.com/sipa/bech32/ref/go/src/bech32"
    "github.com/zlabwork/go-chain/utils"
)

type Address struct {
    pub  []byte
    hash []byte
}

func NewAddress(b []byte) *Address {
    h, _ := utils.Hash160(b)
    return &Address{
        pub:  b,
        hash: h,
    }
}

// @docs https://learnmeabitcoin.com/technical/public-key-hash
// @docs https://learnmeabitcoin.com/technical/address
// format: m/44'/0'/0' support: imToken, bitPay
func (ad *Address) P2PKH() string {
    return p2pkh(ad.hash)
}

func (ad *Address) P2SH() string {
    return p2sh(ad.hash)
}

// Native Address
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
// format: m/44'/0'/0' support: bitPay
func (ad *Address) P2WPKH() string {
    var program []int
    for _, i := range ad.hash {
        program = append(program, int(i))
    }
    addr, err := bech32.SegwitAddrEncode("bc", 0, program)
    if err != nil {
        return ""
    }
    return addr
}

// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
// TODO :: 测试
func (ad *Address) P2WSH() string {
    d := append([]byte{OP_PUSH_33}, ad.pub...)
    d = append(d, OP_CHECKSIG)
    ha, _ := utils.HashSha256(d)
    var program []int
    for _, i := range ha {
        program = append(program, int(i))
    }
    addr, err := bech32.SegwitAddrEncode("bc", 0, program)
    if err != nil {
        return ""
    }
    return addr
}

// P2SH(P2WPKH)
// p2sh-segwit
// format: m/49'/0'/0' support: imToken
func (ad *Address) P2SHP2WPKH() string {
    // OP_0 size hash160
    pre := []byte{OP_0, OP_PUSH_20}
    redeem := append(pre, ad.hash...)
    // P2SH
    ha, _ := utils.Hash160(redeem)
    return p2sh(ha)
}

// P2SH-P2WSH
// https://bitcoincore.org/en/segwit_wallet_dev/
func P2SHP2WSH(pubKey [][]byte, m, n int) (string, error) {
    if m <= 0 || n <= 0 || m > n {
        return "", fmt.Errorf("error OP_M OP_N")
    }
    OP_M := byte(0x50 + m)
    OP_N := byte(0x50 + n)

    // redeem
    redeem := []byte{OP_M}
    for i := 0; i < len(pubKey); i++ {
        if len(pubKey[i]) != 33 {
            return "", fmt.Errorf("public key inside P2SH-P2WSH scripts MUST be compressed key")
        }
        redeem = append(redeem, OP_PUSH_33)
        redeem = append(redeem, pubKey[i]...)
    }
    redeem = append(redeem, OP_N)
    redeem = append(redeem, OP_CHECKMULTISIG)

    ha, _ := utils.HashSha256(redeem)
    hash160, _ := utils.Hash160(append([]byte{OP_0, OP_PUSH_32}, ha...))

    // P2SH
    return p2sh(hash160), nil
}
