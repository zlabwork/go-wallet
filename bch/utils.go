package bch

import (
    "fmt"
    btcutil "github.com/FactomProject/btcutilecc"
    "github.com/zlabwork/gochain/utils"
)

var (
    curve = btcutil.Secp256k1()
)

const (
    addrTypeP2PKH       int = 0
    addrTypeP2SH        int = 1
    pubKeyCompressedLen     = 33
    charset                 = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
)

func createCashAddress(b []byte, pre string, t int) string {
    k, err := packAddrData(t, b)
    if err != nil {
        return ""
    }
    return encode(pre, k)
}

func packAddrData(addrType int, addrHash []byte) ([]byte, error) {
    // Pack addr data with version byte.
    if addrType != addrTypeP2PKH && addrType != addrTypeP2SH {
        return nil, fmt.Errorf("invalid AddressType")
    }
    versionByte := uint(addrType) << 3
    encodedSize := (uint(len(addrHash)) - 20) / 4
    if (len(addrHash)-20)%4 != 0 {
        return nil, fmt.Errorf("invalid address hash size")
    }
    if encodedSize < 0 || encodedSize > 8 {
        return nil, fmt.Errorf("encoded size out of valid range")
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

func encode(pre string, payload []byte) string {
    sum := checksum(pre, payload)
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
