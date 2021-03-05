package bitcoincash

import (
    "bytes"
    "crypto/rand"
    "errors"
    btcutil "github.com/FactomProject/btcutilecc"
    "github.com/zlabwork/go-chain/utils"
    "golang.org/x/crypto/ripemd160"
    "math/big"
)

const (
    addrTypeP2PKH             addrType = 0
    addrTypeP2SH              addrType = 1
    publicKeyCompressedLength          = 33
    charset                            = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
)

var (
    curve = btcutil.Secp256k1()
)

type addrType int

type addressData struct {
    hash   []byte
    prefix string
}

type priKeyData struct {
    key []byte
}

type pubKeyData struct {
    key []byte
}

func NewPriKeyRandom() *priKeyData {
    b := make([]byte, 32)
    rand.Read(b)
    return &priKeyData{
        key: b,
    }
}

func NewPriKey(priKey []byte) *priKeyData {
    return &priKeyData{
        key: priKey,
    }
}

func NewPubKey(priKey []byte) *pubKeyData {
    curve.ScalarBaseMult(priKey)
    return &pubKeyData{
        key: compressPublicKey(curve.ScalarBaseMult(priKey)),
    }
}

func NewAddress(pubKey []byte) *addressData {
    h, _ := hash160(pubKey)
    return &addressData{
        hash:   h[:ripemd160.Size],
        prefix: "bitcoincash", // bchreg、bchtest、bchsim
    }
}

func (pri *priKeyData) Key() []byte {
    return pri.key
}

func (pri *priKeyData) PubKey() *pubKeyData {
    return NewPubKey(pri.key)
}

func (pub *pubKeyData) Key() []byte {
    return pub.key
}

func (pub *pubKeyData) Address() *addressData {
    return NewAddress(pub.key)
}

func (addr *addressData) String() string {
    return addr.P2PKH()
}

func (addr *addressData) Hash160() []byte {
    return addr.hash
}

func (addr *addressData) P2PKH() string {
    return newCashAddress(addr.hash, addr.prefix, addrTypeP2PKH)
}

func (addr *addressData) P2SH() string {
    return newCashAddress(addr.hash, addr.prefix, addrTypeP2SH)
}

func newCashAddress(input []byte, prefix string, t addrType) string {
    k, err := packAddressData(t, input)
    if err != nil {
        return ""
    }
    return encode(prefix, k)
}

func packAddressData(addrType addrType, addrHash []byte) ([]byte, error) {
    // Pack addr data with version byte.
    if addrType != addrTypeP2PKH && addrType != addrTypeP2SH {
        return nil, errors.New("invalid AddressType")
    }
    versionByte := uint(addrType) << 3
    encodedSize := (uint(len(addrHash)) - 20) / 4
    if (len(addrHash)-20)%4 != 0 {
        return nil, errors.New("invalid address hash size")
    }
    if encodedSize < 0 || encodedSize > 8 {
        return nil, errors.New("encoded size out of valid range")
    }
    versionByte |= encodedSize
    var addrHashUint []byte
    addrHashUint = append(addrHashUint, addrHash...)
    data := append([]byte{byte(versionByte)}, addrHashUint...)
    packedData, err := utils.ConvertBits(data, 8, 5, true)
    if err != nil {
        return []byte{}, err
    }
    return packedData, nil
}

func encode(prefix string, payload []byte) string {
    sum := checksum(prefix, payload)
    combined := cat(payload, sum)
    ret := ""

    for _, c := range combined {
        ret += string(charset[c])
    }
    return ret
}

// @see https://github.com/gcash/bchutil/blob/master/address.go
func polyMod(v []byte) uint64 {

    c := uint64(1)
    for _, d := range v {

        c0 := byte(c >> 35)
        c = ((c & 0x07ffffffff) << 5) ^ uint64(d)

        if c0&0x01 > 0 {
            c ^= 0x98f2bc8e61
        }
        if c0&0x02 > 0 {
            c ^= 0x79b76d99e2
        }
        if c0&0x04 > 0 {
            c ^= 0xf33e5fb3c4
        }
        if c0&0x08 > 0 {
            c ^= 0xae2eabe2a8
        }
        if c0&0x10 > 0 {
            c ^= 0x1e4f43e470
        }
    }

    return c ^ 1
}

func checksum(prefix string, payload []byte) []byte {
    enc := cat(expandPrefix(prefix), payload)
    // Append 8 zeroes.
    enc = cat(enc, []byte{0, 0, 0, 0, 0, 0, 0, 0})
    // Determine what to XOR into those 8 zeroes.
    mod := polyMod(enc)
    ret := make([]byte, 8)
    for i := 0; i < 8; i++ {
        // Convert the 5-bit groups in mod to checksum values.
        ret[i] = byte((mod >> uint(5*(7-i))) & 0x1f)
    }
    return ret
}

func expandPrefix(prefix string) []byte {
    ret := make([]byte, len(prefix)+1)
    for i := 0; i < len(prefix); i++ {
        ret[i] = prefix[i] & 0x1f
    }

    ret[len(prefix)] = 0
    return ret
}

func cat(x, y []byte) []byte {
    return append(x, y...)
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
