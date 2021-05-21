package bitcoin

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "github.com/mr-tron/base58"
    "math/big"
)

type priKeyData struct {
    key []byte
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
