package btc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/mr-tron/base58"
)

type addrDesc struct {
	Ver      uint8
	AddrType string
	Addr     string
	Hash160  []byte
}

type msgTx struct {
	Version  int32
	TxIn     []*VIn
	TxOut    []*VOut
	LockTime uint32
}

type transaction struct {
}

func NewTransaction() *transaction {
	return &transaction{}
}

// CreateRawTx
// @link https://www.royalfork.org/2014/11/20/txn-demo/
// @link https://developer.bitcoin.org/reference/transactions.html#raw-transaction-format
func (tx *transaction) CreateRawTx(ins []VIn, outs []VOut, lockTime uint32) ([]byte, error) {

	// 格式: 01000000 NUM01 INPUT NUM02 OUTPUTS 00000000
	ver := []byte{0x01, 0x00, 0x00, 0x00} // 版本
	t1 := byte(len(ins))                  // 输入数量
	t2 := byte(len(outs))                 // 输出数量
	lt := []byte{0x00, 0x00, 0x00, 0x00}  // 锁定时间

	// inputs
	var inputs []byte
	for _, i := range ins {
		in, err := tx.txIn(i)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, in...)
	}

	// outputs
	var outputs []byte
	for _, o := range outs {
		ou, err := tx.txOut(o)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, ou...)
	}

	r := append(ver, t1)
	r = append(r, inputs...)
	r = append(r, t2)
	r = append(r, outputs...)
	r = append(r, lt...)
	return r, nil
}

// FIXME:: 未完待续
// @docs https://www.royalfork.org/2014/11/20/txn-demo/
// @link https://developer.bitcoin.org/reference/transactions.html#txin-a-transaction-input-non-coinbase
func (tx *transaction) txIn(in VIn) ([]byte, error) {
	// ======================================================================
	// 格式: tx的反转 + INDEX (uint32) + length of script + script + 0xFFFFFFFF
	// ======================================================================

	// 1. txId - little-endian 反转
	txId1, err := hex.DecodeString(in.Tx)
	l := len(txId1)
	txId := make([]byte, l)
	for i := 0; i < l; i++ {
		txId[i] = txId1[l-i-1]
	}

	// 2. vOut index
	idx := make([]byte, 4)
	binary.LittleEndian.PutUint32(idx, in.N)

	// TODO :: 签名格式
	// Signature contains 5 items:
	// 1 byte length of the following 2 fields,
	// the DER encoded signature,
	// the hash type (this is usually 1, but there are several hash types),
	// 1 byte length of public key,
	// then the public key.
	sig, err := hex.DecodeString("47304402204e572c0587b2147efaa5685b470350bad9561c359056ecb2abb0eca05bc612f502203aae1b45aa24215b2575a26871f18c95fb1b911eaed7705eaf53cb3a2b031ea0012103c13dca192f1ba64265d8efca97d43b822ff24db357c13b0e6e0395cf91e9efae")
	if err != nil {
		return nil, err
	}
	len := byte(len(sig))
	end := []byte{0xFF, 0xFF, 0xFF, 0xFF}

	var r []byte
	r = append(r, txId...)
	r = append(r, idx...)
	r = append(r, len)
	r = append(r, sig...)
	r = append(r, end...)

	return r, nil
}

// @docs https://www.royalfork.org/2014/11/20/txn-demo/
// @link https://developer.bitcoin.org/reference/transactions.html#txout-a-transaction-output
func (tx *transaction) txOut(out VOut) ([]byte, error) {

	// ============================================
	// 格式: sat value + pk_script bytes + pk_script
	// pk_script 的最大长度 10,000 bytes
	// ============================================

	addr := out.Addr
	sat := out.Amt

	// 1. amount in satoshis
	val := make([]byte, 8)
	binary.LittleEndian.PutUint64(val, uint64(sat))

	// 2. 锁定脚本
	pks, err := tx.pkScript(addr)
	if err != nil {
		return nil, err
	}

	// 3. length
	bf := bytes.NewBuffer(nil)
	binary.Write(bf, binary.BigEndian, uint8(len(pks)))
	l := bf.Bytes()

	var r []byte
	r = append(val, l...)
	r = append(r, pks...)
	return r, nil
}

// 锁定脚本 - Lock Script
func (tx *transaction) pkScript(addr string) ([]byte, error) {

	desc, err := tx.parseAddr(addr)
	if err != nil {
		return nil, err
	}
	var r []byte

	switch desc.Ver {
	case 0x00, 0x6F: // P2PKH
		r = append(r, OP_DUP)
		r = append(r, OP_HASH160)
		r = append(r, byte(len(desc.Hash160))) // 0x14, Push 20 bytes as data TODO :: 是否进一步确认？
		r = append(r, desc.Hash160...)
		r = append(r, OP_EQUALVERIFY)
		r = append(r, OP_CHECKSIG)

	case 0x04, 0x05: // P2SH TODO :: 补充

	}

	return r, nil
}

func (tx *transaction) parseAddr(addr string) (*addrDesc, error) {

	// TODO :: 验证 checksum
	b, err := base58.Decode(addr)
	if err != nil {
		return nil, err
	}

	return &addrDesc{
		AddrType: "p2pkh",
		Addr:     addr,
		Ver:      b[0],
		Hash160:  b[1:21],
	}, nil
}
