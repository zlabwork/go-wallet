package bitcoin

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "github.com/FactomProject/btcutilecc"
    "github.com/mr-tron/base58"
    "github.com/sipa/bech32/ref/go/src/bech32"
    "math/big"
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

type addrData struct {
    hash   []byte
    pubKey []byte
}

type priKeyData struct {
    key []byte
}

type pubKeyData struct {
    key []byte
}

// TODO :: 私钥值约束，最大不能大于 fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140
func NewPriKeyRandom() *priKeyData {
    b := make([]byte, 32)
    rand.Read(b)
    return &priKeyData{
        key: b,
    }
}

// https://learnmeabitcoin.com/technical/private-key
func NewPriKey(priKey []byte) (*priKeyData, error) {
    if len(priKey) != 32 {
        return nil, errors.New("invalid length")
    }
    return &priKeyData{
        key: priKey,
    }, nil
}

// 脑钱包
func NewBrainWallet(words, salt string) (*priKeyData, error) {
    priKey, err := hashSha256([]byte(words + salt))
    if err != nil {
        return nil, err
    }
    return &priKeyData{
        key: priKey,
    }, nil
}

// 脑钱包 - 非通用方法
func NewBrainWalletSpecial(words, salt string) (*priKeyData, error) {
    wh, _ := hashSha256([]byte(words))
    sh, _ := hashSha256([]byte(salt))
    bs := append(append(wh, sh...), hashConst...)
    priKey, err := hashSha256(bs)
    if err != nil {
        return nil, err
    }
    return &priKeyData{
        key: priKey,
    }, nil
}

// parse WIF
func ParseWIF(wif string) (*priKeyData, error) {
    b, err := base58.Decode(wif)
    if err != nil {
        return nil, err
    }
    if len(b) < 33 {
        return nil, errors.New("invalid WIF data")
    }
    return NewPriKey(b[1:33])
}

func (pri *priKeyData) Key() []byte {
    return pri.key
}

func (pri *priKeyData) PubKey() *pubKeyData {
    return &pubKeyData{
        key: pri.compressPublicKey(curve.ScalarBaseMult(pri.key)),
    }
}

func (pri *priKeyData) PubKeyUncompressed() *pubKeyData {
    return &pubKeyData{
        key: pri.uncompressedPublicKey(curve.ScalarBaseMult(pri.key)),
    }
}

// https://learnmeabitcoin.com/technical/wif
func (pri *priKeyData) WIF() string {
    version := []byte{0x80}
    compression := byte(0x01)
    key := append(version, pri.key...)
    key = append(key, compression)
    sum, err := checksum(key)
    if err != nil {
        return ""
    }
    key = append(key, sum...)
    return base58.Encode(key)
}

func (pri *priKeyData) Base64() string {
    return base64.StdEncoding.EncodeToString(pri.key)
}

func (pri *priKeyData) Hex() string {
    return hex.EncodeToString(pri.key)
}

func (pri *priKeyData) compressPublicKey(x *big.Int, y *big.Int) []byte {
    var key bytes.Buffer

    // Write header; 0x2 for even y value; 0x3 for odd
    key.WriteByte(byte(0x2) + byte(y.Bit(0)))

    // Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
    xBytes := x.Bytes()
    for i := 0; i < (publicKeyCompressedLength - 1 - len(xBytes)); i++ {
        key.WriteByte(0x0)
    }
    key.Write(xBytes)

    return key.Bytes()
}

func (pri *priKeyData) uncompressedPublicKey(x *big.Int, y *big.Int) []byte {
    var key bytes.Buffer
    key.WriteByte(byte(0x4))
    key.Write(x.Bytes())
    key.Write(y.Bytes())
    return key.Bytes()
}

// https://learnmeabitcoin.com/technical/public-key
func NewPubKey(pubKey []byte) (*pubKeyData, error) {
    if len(pubKey) == 33 || len(pubKey) == 65 {
        return &pubKeyData{
            key: pubKey,
        }, nil
    }
    return nil, errors.New("invalid length")
}

func (pub *pubKeyData) Key() []byte {
    return pub.key
}

func (pub *pubKeyData) Address() *addrData {
    h, _ := hash160(pub.key)
    return &addrData{
        hash:   h,
        pubKey: pub.key,
    }
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
func (addr *addrData) P2PKH() string {
    return p2pkh(addr.hash)
}

func (addr *addrData) P2SH() string {
    return p2sh(addr.hash)
}

// Native Address
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
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
func (addr *addrData) P2SHP2WPKH() string {
    // OP_0 size hash160
    pre := []byte{OP_0, OP_PUSH_20}
    redeem := append(pre, addr.hash...)

    // P2SH
    ha, _ := hash160(redeem)
    return p2sh(ha)
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
