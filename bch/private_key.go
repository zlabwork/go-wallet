package bch

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "github.com/mr-tron/base58"
)

type PriKey struct {
    k []byte
}

// TODO :: 私钥值约束，最大不能大于 fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140
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

// TODO
func (pri *PriKey) WIF() string {
    return ""
}

func (pri *PriKey) PubKey() *PubKey {

    x, y := curve.ScalarBaseMult(pri.Bytes())

    var key bytes.Buffer

    // Write header; 0x2 for even y value; 0x3 for odd
    key.WriteByte(byte(0x2) + byte(y.Bit(0)))

    // Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
    xBytes := x.Bytes()
    for i := 0; i < (pubKeyCompressedLen - 1 - len(xBytes)); i++ {
        key.WriteByte(0x0)
    }
    key.Write(xBytes)

    pub, _ := NewPubKey(key.Bytes())
    return pub
}
