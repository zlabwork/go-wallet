package bitcoin

import (
    "bytes"
    "crypto/rand"
    "errors"
    "github.com/FactomProject/btcutilecc"
    "github.com/mr-tron/base58"
    "github.com/sipa/bech32/ref/go/src/bech32"
    "math/big"
)

var (
    curve                     = btcutil.Secp256k1()
    PublicKeyCompressedLength = 33
)

// https://learnmeabitcoin.com/technical/private-key
func GenPriKey() []byte {
    // TODO :: 私钥值约束，最大不能大于 fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140
    b := make([]byte, 32)
    rand.Read(b)
    return b
}

// https://learnmeabitcoin.com/technical/public-key
func GenPubKey(priKey []byte) []byte {
    curve.ScalarBaseMult(priKey)
    return compressPublicKey(curve.ScalarBaseMult(priKey))
}

func GenPubKeyUncompressed(priKey []byte) []byte {
    curve.ScalarBaseMult(priKey)
    return uncompressedPublicKey(curve.ScalarBaseMult(priKey))
}

// https://learnmeabitcoin.com/technical/wif
func WIF(priKey []byte) string {
    version := []byte{0x80}
    compression := byte(0x01)
    key := append(version, priKey...)
    key = append(key, compression)
    sum, err := checksum(key)
    if err != nil {
        return ""
    }
    key = append(key, sum...)
    return base58.Encode(key)
}

// parse WIF
func ParseWIF(wif string) ([]byte, error) {
    b, err := base58.Decode(wif)
    if err != nil {
        return nil, err
    }
    if len(b) < 33 {
        return nil, errors.New("error WIF data")
    }
    return b[1:33], nil
}

// @docs https://learnmeabitcoin.com/technical/public-key-hash
// @docs https://learnmeabitcoin.com/technical/address
func P2PKH(pubKey []byte) string {
    h, _ := hash160(pubKey)

    prefix := []byte{0x00}
    preData := append(prefix, h...)
    sum, _ := checksum(preData)
    addr := append(preData, sum...)

    return base58.Encode(addr)
}

func P2SH(pubKey []byte) string {
    h, _ := hash160(pubKey)

    prefix := []byte{0x05}
    preData := append(prefix, h...)
    sum, _ := checksum(preData)
    addr := append(preData, sum...)

    return base58.Encode(addr)
}

// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
func Segwit(pubKey []byte) string {
    h, _ := hash160(pubKey)
    var program []int
    for _, i := range h {
        program = append(program, int(i))
    }
    addr, err := bech32.SegwitAddrEncode("bc", 0, program)
    if err != nil {
        return ""
    }
    return addr
}

func Addr2Hash160(address string) ([]byte, error) {
    if address[0:2] == "bc" {
        _, n, err := bech32.SegwitAddrDecode("bc", address)
        if err != nil {
            return nil, err
        }
        var bs []byte
        for _, d := range n {
            bs = append(bs, byte(d))
        }
        return bs, err
    }

    b, err := base58.Decode(address)
    if err != nil {
        return nil, err
    }
    return b[1:21], nil
}

func compressPublicKey(x *big.Int, y *big.Int) []byte {
    var key bytes.Buffer

    // Write header; 0x2 for even y value; 0x3 for odd
    key.WriteByte(byte(0x2) + byte(y.Bit(0)))

    // Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
    xBytes := x.Bytes()
    for i := 0; i < (PublicKeyCompressedLength - 1 - len(xBytes)); i++ {
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
