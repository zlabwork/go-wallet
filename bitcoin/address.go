package bitcoin

import (
    "errors"
    "github.com/mr-tron/base58"
    "github.com/sipa/bech32/ref/go/src/bech32"
)

type addrData struct {
    hash   []byte
    pubKey []byte
}

func ParseAddress(addr string) (*addrData, error) {
    if addr[0:2] == "bc" {
        _, n, err := bech32.SegwitAddrDecode("bc", addr)
        if err != nil {
            return nil, err
        }
        var bs []byte
        for _, d := range n {
            bs = append(bs, byte(d))
        }
        return &addrData{hash: bs}, err
    }

    b, err := base58.Decode(addr)
    if err != nil {
        return nil, err
    }
    return &addrData{hash: b[1:21]}, nil
}

func (addr *addrData) Hash160() []byte {
    return addr.hash
}

// @docs https://learnmeabitcoin.com/technical/public-key-hash
// @docs https://learnmeabitcoin.com/technical/address
// format: m/44'/0'/0' support: imToken, bitPay
func (addr *addrData) P2PKH() string {
    return p2pkh(addr.hash)
}

func (addr *addrData) P2SH() string {
    return p2sh(addr.hash)
}

// Native Address
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
// format: m/44'/0'/0' support: bitPay
func (addr *addrData) P2WPKH() (string, error) {
    var program []int
    for _, i := range addr.hash {
        program = append(program, int(i))
    }
    address, err := bech32.SegwitAddrEncode("bc", 0, program)
    if err != nil {
        return "", err
    }
    return address, nil
}

// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
// TODO :: 测试
func (addr *addrData) P2WSH() (string, error) {
    if addr.pubKey == nil {
        return "", errors.New("pubKey is not specified")
    }
    data := append([]byte{OP_PUSH_33}, addr.pubKey...)
    data = append(data, OP_CHECKSIG)
    ha, _ := hashSha256(data)

    var program []int
    for _, i := range ha {
        program = append(program, int(i))
    }
    address, err := bech32.SegwitAddrEncode("bc", 0, program)
    if err != nil {
        return "", err
    }
    return address, nil
}

// P2SH(P2WPKH)
// p2sh-segwit
// format: m/49'/0'/0' support: imToken
func (addr *addrData) P2SHP2WPKH() string {
    // OP_0 size hash160
    pre := []byte{OP_0, OP_PUSH_20}
    redeem := append(pre, addr.hash...)

    // P2SH
    ha, _ := hash160(redeem)
    return p2sh(ha)
}

// hash160
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

// P2SH-P2WSH
// https://bitcoincore.org/en/segwit_wallet_dev/
func P2SHP2WSH(pubKey [][]byte, m, n int) (string, error) {
    if m <= 0 || n <= 0 || m > n {
        return "", errors.New("error OP_M OP_N")
    }
    OP_M := byte(0x50 + m)
    OP_N := byte(0x50 + n)

    // redeem
    redeem := []byte{OP_M}
    for i := 0; i < len(pubKey); i++ {
        if len(pubKey[i]) != 33 {
            return "", errors.New("public key inside P2SH-P2WSH scripts MUST be compressed key")
        }
        redeem = append(redeem, OP_PUSH_33)
        redeem = append(redeem, pubKey[i]...)
    }
    redeem = append(redeem, OP_N)
    redeem = append(redeem, OP_CHECKMULTISIG)

    ha, _ := hashSha256(redeem)
    hash160, _ := hash160(append([]byte{OP_0, OP_PUSH_32}, ha...))

    // P2SH
    return p2sh(hash160), nil
}
