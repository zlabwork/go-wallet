package bitcoin

import (
    "github.com/mr-tron/base58"
)

func TestNetWIF(priKey []byte) string {
    version := []byte{0xef}
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

func TestNetP2PKH(pubKey []byte) string {
    h, _ := hash160(pubKey)

    prefix := []byte{0x6f}
    preData := append(prefix, h...)
    sum, _ := checksum(preData)
    addr := append(preData, sum...)

    return base58.Encode(addr)
}

func TestNetP2SH(pubKey []byte) string {
    h, _ := hash160(pubKey)

    prefix := []byte{0xc4}
    preData := append(prefix, h...)
    sum, _ := checksum(preData)
    addr := append(preData, sum...)

    return base58.Encode(addr)
}
