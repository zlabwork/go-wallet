package btc

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "github.com/mr-tron/base58"
    "github.com/zlabwork/gochain/utils"
)

type PriKey struct {
    k []byte
}

func NewPriKeyRandom() *PriKey {
    b := make([]byte, 32)
    rand.Read(b)
    return &PriKey{k: b}
}

func NewPriKey(b []byte) (*PriKey, error) {
    if len(b) != 32 {
        return nil, fmt.Errorf("%s", "invalid length")
    }
    return &PriKey{k: b}, nil
}

// 脑钱包
func NewPriKeyBrain(words, salt string) (*PriKey, error) {
    k, err := utils.HashSha256([]byte(words + salt))
    if err != nil {
        return nil, err
    }
    return &PriKey{
        k: k,
    }, nil
}

// 脑钱包 - 非通用方法
func NewPriKeyBrainSP(words, salt string) (*PriKey, error) {
    h1, _ := utils.HashSha256([]byte(words))
    h2, _ := utils.HashSha256([]byte(salt))
    b := append(append(h1, h2...), hashConst...)
    k, err := utils.HashSha256(b)
    if err != nil {
        return nil, err
    }
    return &PriKey{
        k: k,
    }, nil
}

func ParseWIF(wif string) (*PriKey, error) {
    b, err := base58.Decode(wif)
    if err != nil {
        return nil, err
    }
    if len(b) < 33 {
        return nil, fmt.Errorf("invalid WIF data")
    }
    return NewPriKey(b[1:33])
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

func (pri *PriKey) WIF() string {
    ver := []byte{0x80}
    compression := byte(0x01)
    k := append(ver, pri.k...)
    k = append(k, compression)
    s, err := checksum(k)
    if err != nil {
        return ""
    }
    k = append(k, s...)
    return base58.Encode(k)
}

func (pri *PriKey) PubKey() *PubKey {
    x, y := curve.ScalarBaseMult(pri.Bytes())
    var k bytes.Buffer
    // Write header; 0x2 for even y value; 0x3 for odd
    k.WriteByte(byte(0x2) + byte(y.Bit(0)))
    // Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
    xBytes := x.Bytes()
    for i := 0; i < (pubKeyCompressedLen - 1 - len(xBytes)); i++ {
        k.WriteByte(0x0)
    }
    k.Write(xBytes)
    p, _ := NewPubKey(k.Bytes())
    return p
}

func (pri *PriKey) PubKeyUnCompressed() *PubKey {
    x, y := curve.ScalarBaseMult(pri.Bytes())
    var k bytes.Buffer
    k.WriteByte(byte(0x4))
    k.Write(x.Bytes())
    k.Write(y.Bytes())
    p, _ := NewPubKey(k.Bytes())
    return p
}
