package btc

import (
	btcutil "github.com/FactomProject/btcutilecc"
	"github.com/mr-tron/base58"
	"github.com/zlabwork/go-chain/utils"
)

var (
	curve               = btcutil.Secp256k1()
	pubKeyCompressedLen = 33
	hashConst, _        = utils.HashSha256([]byte("zlab")) // hash 常量
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
