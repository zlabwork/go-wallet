package btc

import (
	btcutil "github.com/FactomProject/btcutilecc"
	"github.com/mr-tron/base58"
	"github.com/sipa/bech32/ref/go/src/bech32"
	"github.com/zlabwork/go-wallet/utils"
)

var (
	curve               = btcutil.Secp256k1()
	pubKeyCompressedLen = 33
)

func checksum(data []byte) ([]byte, error) {
	hash, err := utils.HashDoubleSha256(data)
	if err != nil {
		return nil, err
	}
	return hash[:4], nil
}

func p2sh(hash160 []byte) string {
	data := append([]byte{getVer("P2SH")}, hash160...)
	sum, _ := checksum(data)
	return base58.Encode(append(data, sum...))
}

func p2pkh(hash160 []byte) string {
	data := append([]byte{getVer("P2PKH")}, hash160...)
	sum, _ := checksum(data)
	return base58.Encode(append(data, sum...))
}

// ParseHash160
// only support: P2pkh, P2sh, P2wpkh
// todo: support regtest testnet
func ParseHash160(address string) ([]byte, error) {
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
