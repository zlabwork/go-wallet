package bitcoin

import (
    "github.com/mr-tron/base58"
)

func (pri *priKeyData) TestNetWIF() string {
    version := []byte{0xef}
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

func (addr *addrData) TestNetP2PKH() string {

    prefix := []byte{0x6f}
    preData := append(prefix, addr.hash...)
    sum, _ := checksum(preData)
    address := append(preData, sum...)

    return base58.Encode(address)
}

func (addr *addrData) TestNetP2SH() string {

    prefix := []byte{0xc4}
    preData := append(prefix, addr.hash...)
    sum, _ := checksum(preData)
    address := append(preData, sum...)

    return base58.Encode(address)
}
