package utils

import (
	"bytes"
	"errors"
)

const (
	hex1F = 0x1F // 00011111
)

// base32
var charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

func Base32Encode(data []byte, padding bool) (string, error) {
	bits, err := ConvertBits(data, 8, 5, padding)
	if err != nil {
		return "", err
	}
	ret := ""
	for _, c := range bits {
		ret += string(charset[c])
	}
	return ret, nil
}

// @deprecated
// Deprecated: 自写函数实现暂时废弃
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
func ConvertBits(data []byte, fromBits uint, tobits uint, pad bool) ([]byte, error) {
	// General power-of-2 base conversion.
	var uintArr []uint
	for _, i := range data {
		uintArr = append(uintArr, uint(i))
	}
	acc := uint(0)
	bits := uint(0)
	var ret []uint
	maxv := uint((1 << tobits) - 1)
	maxAcc := uint((1 << (fromBits + tobits - 1)) - 1)
	for _, value := range uintArr {
		acc = ((acc << fromBits) | value) & maxAcc
		bits += fromBits
		for bits >= tobits {
			bits -= tobits
			ret = append(ret, (acc>>bits)&maxv)
		}
	}
	if pad {
		if bits > 0 {
			ret = append(ret, (acc<<(tobits-bits))&maxv)
		}
	} else if bits >= fromBits || ((acc<<(tobits-bits))&maxv) != 0 {
		return []byte{}, errors.New("encoding padding error")
	}
	var dataArr []byte
	for _, i := range ret {
		dataArr = append(dataArr, byte(i))
	}
	return dataArr, nil
}
