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

var (
    curve                     = btcutil.Secp256k1()
    publicKeyCompressedLength = 33
    hashConst, _              = hashSha256([]byte("zlab")) // hash 常量
)

type addrData struct {
    hash []byte
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
    if len(priKey) < 32 {
        return nil, errors.New("invalid length")
    }
    return &priKeyData{
        key: priKey,
    }, nil
}

// https://learnmeabitcoin.com/technical/public-key
func NewPubKey(priKey []byte) (*pubKeyData, error) {
    if len(priKey) < 32 {
        return nil, errors.New("invalid length")
    }
    return &pubKeyData{
        key: compressPublicKey(curve.ScalarBaseMult(priKey)),
    }, nil
}

func NewPubKeyUncompressed(priKey []byte) (*pubKeyData, error) {
    if len(priKey) < 32 {
        return nil, errors.New("invalid length")
    }
    return &pubKeyData{
        key: uncompressedPublicKey(curve.ScalarBaseMult(priKey)),
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

func (pri *priKeyData) Key() []byte {
    return pri.key
}

func (pri *priKeyData) PubKey() *pubKeyData {
    pub, err := NewPubKey(pri.key)
    if err != nil {
        return nil
    }
    return pub
}

func (pri *priKeyData) PubKeyUncompressed() *pubKeyData {
    pub, err := NewPubKeyUncompressed(pri.key)
    if err != nil {
        return nil
    }
    return pub
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

func (pub *pubKeyData) Key() []byte {
    return pub.key
}

func (pub *pubKeyData) Address() *addrData {
    h, _ := hash160(pub.key)
    return &addrData{
        hash: h,
    }
}

func (addr *addrData) Hash160() []byte {
    return addr.hash
}

// @docs https://learnmeabitcoin.com/technical/public-key-hash
// @docs https://learnmeabitcoin.com/technical/address
func (addr *addrData) P2PKH() string {
    prefix := []byte{0x00}
    preData := append(prefix, addr.hash...)
    sum, _ := checksum(preData)
    address := append(preData, sum...)
    return base58.Encode(address)
}

func (addr *addrData) P2SH() string {
    prefix := []byte{0x05}
    preData := append(prefix, addr.hash...)
    sum, _ := checksum(preData)
    address := append(preData, sum...)
    return base58.Encode(address)
}

// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
func (addr *addrData) Segwit() string {
    var program []int
    for _, i := range addr.hash {
        program = append(program, int(i))
    }
    address, err := bech32.SegwitAddrEncode("bc", 0, program)
    if err != nil {
        return ""
    }
    return address
}

func compressPublicKey(x *big.Int, y *big.Int) []byte {
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

func uncompressedPublicKey(x *big.Int, y *big.Int) []byte {
    var key bytes.Buffer
    key.WriteByte(byte(0x4))
    key.Write(x.Bytes())
    key.Write(y.Bytes())
    return key.Bytes()
}
