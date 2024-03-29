package btc

import (
	"fmt"
	"github.com/sipa/bech32/ref/go/src/bech32"
	"github.com/zlabwork/go-wallet/utils"
)

type Address struct {
	pub  []byte
	hash []byte
}

func NewAddress(b []byte) *Address {
	h, _ := utils.Hash160(b)
	return &Address{
		pub:  b,
		hash: h,
	}
}

// P2pkh
// e.g. 19c4pkCL2jvTFYkZXDyUHi4ceoNze44mXE
// @docs https://learnmeabitcoin.com/technical/public-key-hash
// @docs https://learnmeabitcoin.com/technical/address
// format: m/44'/0'/0' support: imToken, bitPay
func (ad *Address) P2pkh() string {
	return p2pkh(ad.hash)
}

// P2sh
// e.g. 3AJ5kHgmaeEqLiSzeKe4iLRYoKfiCH5Y1C
func (ad *Address) P2sh() string {
	return p2sh(ad.hash)
}

// P2wpkh
// e.g. bc1qte3vrg28lxm5yjh4579cjqzgadgmrghjm6hvjt
// Native Address
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
// format: m/44'/0'/0' support: bitPay
func (ad *Address) P2wpkh() string {
	var program []int
	for _, i := range ad.hash {
		program = append(program, int(i))
	}
	hrp := getHrpFlag()
	addr, err := bech32.SegwitAddrEncode(hrp, 0, program)
	if err != nil {
		return ""
	}
	return addr
}

// P2wsh
// e.g. bc1qpvw9q3u9yx9ga452yr2q4hypgnp8kqxfku9lcvxutlldqqcl06fs8pdyj8
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
// TODO :: 测试
func (ad *Address) P2wsh() string {
	d := append([]byte{OP_PUSH_33}, ad.pub...)
	d = append(d, OP_CHECKSIG)
	ha, _ := utils.HashSha256(d)
	var program []int
	for _, i := range ha {
		program = append(program, int(i))
	}
	hrp := getHrpFlag()
	addr, err := bech32.SegwitAddrEncode(hrp, 0, program)
	if err != nil {
		return ""
	}
	return addr
}

// P2wpkhInP2sh - P2WPKH nested in P2SH
// e.g. 3J8VzKMkGwzneEs6imQrGX2jgNe8gwdyNn
// p2sh-segwit
// format: m/49'/0'/0' support: imToken
func (ad *Address) P2wpkhInP2sh() string {

	// public key used in P2SH-P2WPKH MUST be compressed, i.e. 33 bytes in size
	if len(ad.pub) != 33 {
		return ""
	}

	// OP_0 size hash160
	pre := []byte{OP_0, OP_PUSH_20}
	redeem := append(pre, ad.hash...)
	// P2SH
	ha, _ := utils.Hash160(redeem)
	return p2sh(ha)
}

// MultiP2wshInP2sh
// e.g. 3Ly7sZXyv9zNKKV35ntbgraLy9zaykzKQL
// P2WSH nested in P2SH (1-of-1 multisig)
func (ad *Address) MultiP2wshInP2sh() string {
	r, err := MultiP2wshInP2sh([][]byte{ad.pub}, 1, 1)
	if err != nil {
		return ""
	}
	return r
}

// MultiP2wsh
// e.g. bc1q0fawq3lvmhq47443f5xsp7l95qq4xpz5gjjkljwzw933vnectsjs2ty768
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
func (ad *Address) MultiP2wsh() string {
	r, err := MultiP2wsh([][]byte{ad.pub}, 1, 1)
	if err != nil {
		return ""
	}
	return r
}

// MultiP2wshInP2sh -  P2WSH nested in P2SH (1-of-1 multisig)
// https://bitcoincore.org/en/segwit_wallet_dev/
func MultiP2wshInP2sh(pubKey [][]byte, m, n int) (string, error) {
	if m <= 0 || n <= 0 || m > n {
		return "", fmt.Errorf("error OP_M OP_N")
	}
	OP_M := byte(0x50 + m)
	OP_N := byte(0x50 + n)

	// redeem
	redeem := []byte{OP_M}
	for i := 0; i < len(pubKey); i++ {
		if len(pubKey[i]) != 33 {
			return "", fmt.Errorf("public key inside P2SH-P2WSH scripts MUST be compressed key")
		}
		redeem = append(redeem, OP_PUSH_33)
		redeem = append(redeem, pubKey[i]...)
	}
	redeem = append(redeem, OP_N)
	redeem = append(redeem, OP_CHECKMULTISIG)

	ha, _ := utils.HashSha256(redeem)
	hash160, _ := utils.Hash160(append([]byte{OP_0, OP_PUSH_32}, ha...))

	// P2SH
	return p2sh(hash160), nil
}

// MultiP2wsh - P2WSH (1-of1 multisig)
func MultiP2wsh(pubKey [][]byte, m, n int) (string, error) {
	if m <= 0 || n <= 0 || m > n {
		return "", fmt.Errorf("error OP_M OP_N")
	}
	OP_M := byte(0x50 + m)
	OP_N := byte(0x50 + n)

	// redeem
	redeem := []byte{OP_M}
	for i := 0; i < len(pubKey); i++ {
		if len(pubKey[i]) != 33 {
			return "", fmt.Errorf("public key inside P2SH-P2WSH scripts MUST be compressed key")
		}
		redeem = append(redeem, OP_PUSH_33)
		redeem = append(redeem, pubKey[i]...)
	}
	redeem = append(redeem, OP_N)
	redeem = append(redeem, OP_CHECKMULTISIG)

	ha, _ := utils.HashSha256(redeem)

	var program []int
	for _, i := range ha {
		program = append(program, int(i))
	}
	hrp := getHrpFlag()
	addr, err := bech32.SegwitAddrEncode(hrp, 0, program)
	if err != nil {
		return "", nil
	}
	return addr, nil
}
