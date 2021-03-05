package utils

import (
    "bytes"
    "errors"
    "fmt"
)

const (
    hex1F = 0x1F // 00011111
)

// base32
var charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

func Base32Encode(data []byte, padding bool) (string, error) {
    var input []int
    for _, item := range data {
        input = append(input, int(item))
    }
    bits, err := ConvertBits(input, 8, 5, padding)
    if err != nil {
        return "", err
    }
    ret := ""
    for _, i := range bits {
        ret += string(charset[i])
    }
    return ret, nil
}

func Base32Encode2(data []byte, padding bool) (string, error) {
    bits := 5

    // TODO::填充
    if len(data)%5 != 0 {
        return "", errors.New("no padding")
    }

    var ret bytes.Buffer
    l := len(data) * 8 / bits
    for i := 0; i < l; i++ {
        s := i * bits
        n := s / 8
        m := s % 8
        var num uint8
        if m <= 3 {
            num = data[n] >> (3 - m) & hex1F
        } else {
            num = (data[n] << (m - 3) & hex1F) | (data[n+1] >> (11 - m) & hex1F)
        }
        ret.WriteByte(charset[num])
    }

    return ret.String(), nil
}

// @see https://github.com/sipa/bech32/blob/master/ref/go/src/bech32/bech32.go
func ConvertBits(data []int, frombits, tobits uint, pad bool) ([]int, error) {
    acc := 0
    bits := uint(0)
    ret := []int{}
    maxv := (1 << tobits) - 1
    for idx, value := range data {
        if value < 0 || (value>>frombits) != 0 {
            return nil, fmt.Errorf("invalid data range : data[%d]=%d (frombits=%d)", idx, value, frombits)
        }
        acc = (acc << frombits) | value
        bits += frombits
        for bits >= tobits {
            bits -= tobits
            ret = append(ret, (acc>>bits)&maxv)
        }
    }
    if pad {
        if bits > 0 {
            ret = append(ret, (acc<<(tobits-bits))&maxv)
        }
    } else if bits >= frombits {
        return nil, fmt.Errorf("illegal zero padding")
    } else if ((acc << (tobits - bits)) & maxv) != 0 {
        return nil, fmt.Errorf("non-zero padding")
    }
    return ret, nil
}
