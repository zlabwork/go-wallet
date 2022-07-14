package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"strings"
)

type TxIn struct {
	Tx string // txId
	N  uint32 // vout
}

// GetTxOut
// https://developer.bitcoin.org/reference/rpc/gettxout.html
func (rc *RpcClient) GetTxOut(tx string, index uint32) (*TxOut, error) {

	b, err := rc.Request("gettxout", []interface{}{tx, index})
	if err != nil {
		return nil, err
	}

	type response struct {
		Response
		Result TxOut
	}
	var out response
	if err = json.Unmarshal(b, &out); err != nil {
		return nil, err
	}

	//if out.Value == 0 {
	//	return nil, fmt.Errorf("%s, error txId or tx has been spent", tx)
	//}
	return &out.Result, nil
}

func (rc *RpcClient) CreateTransferAll(ins []string, addr string, feeSat int64) (hex string, error error) {

	var t int64
	var inT []TxIn
	for _, str := range ins {
		s := strings.Split(str, ":")
		tx := strings.TrimSpace(s[0])
		n, err := strconv.ParseUint(strings.TrimSpace(s[1]), 10, 64)
		// GetTxOut
		ou, err := rc.GetTxOut(tx, uint32(n))
		if err != nil {
			return "", err
		}

		v, err := strconv.ParseInt(decimal.NewFromFloat(ou.Value).Mul(decimal.NewFromInt(100000000)).String(), 10, 64)
		if err != nil {
			return "", err
		}

		if v < minTxAmount {
			return "", fmt.Errorf("error, [tx:%s vout:%d], %d satoshis, less than minimum number (%d satoshis)", tx, n, v, minTxAmount)
		}
		t += v
		inT = append(inT, TxIn{Tx: tx, N: uint32(n)})
	}

	// fees
	size := 148*len(ins) + 34 + 10
	fee := int64(size) * feeSat

	outs := map[string]int64{
		addr: t - fee,
	}
	return rc.createRawTX(inT, outs, "")
}

// CreateTXAlias - Alias for CreateTX
func (rc *RpcClient) CreateTXAlias(ins []string, outs map[string]int64, hexData string, feeSat int64, changeAddress string) (hex string, error error) {
	var inT []TxIn
	for _, str := range ins {
		s := strings.Split(str, ":")
		n, err := strconv.ParseUint(strings.TrimSpace(s[1]), 10, 64)
		if err != nil {
			return "", errors.New(str + ", error vout format")
		}
		inT = append(inT, TxIn{Tx: strings.TrimSpace(s[0]), N: uint32(n)})
	}
	return rc.CreateTX(inT, outs, hexData, feeSat, changeAddress)
}

func (rc *RpcClient) CreateTX(ins []TxIn, outs map[string]int64, hexData string, feeSat int64, changeAddress string) (hex string, error error) {

	if len(changeAddress) < 10 {
		return "", errors.New("error change address")
	}

	// 1. total in
	var totalIn int64
	for _, in := range ins {
		txOut, err := rc.GetTxOut(in.Tx, in.N)
		if err != nil {
			return "", err
		}
		amt := int64(txOut.Value * math.Pow10(8))
		if amt < minTxAmount {
			return "", fmt.Errorf("error, [tx:%s vout:%d], %d satoshis, less than minimum number (%d satoshis)", in.Tx, in.N, amt, minTxAmount)
		}
		totalIn += amt
	}

	// 2. total out
	var totalOut int64
	isInOuts := false
	for addr, sa := range outs {
		if addr == changeAddress {
			isInOuts = true
		}
		if sa < minTxAmount {
			return "", fmt.Errorf("transfer %d satoshis, less than minimum amount %d satoshis", sa, minTxAmount)
		}
		totalOut += sa
	}

	// TODO: 1. confirm the calculate process
	// 3. fee sat
	size := 148*len(ins) + 34*len(outs) + 10
	fee := int64(size) * feeSat
	left := totalIn - totalOut - fee
	if left < 0 {
		return "", fmt.Errorf("fee is not enough")
	}

	// 4. charge back
	if isInOuts {
		for addr, _ := range outs {
			if addr == changeAddress {
				outs[addr] += left
			}
		}
	} else {
		if left-34*feeSat > minTxAmount {
			outs[changeAddress] = left - 34*feeSat
			fee += 34 * feeSat
		}
	}

	return rc.createRawTX(ins, outs, hexData)
}

// CreateRawTX - TODO: Fix the address duplicated in VOut
// https://developer.bitcoin.org/reference/rpc/createrawtransaction.html
func (rc *RpcClient) createRawTX(ins []TxIn, outs map[string]int64, hexData string) (hex string, error error) {

	// 1. check
	if len(ins) < 1 {
		return "", fmt.Errorf("no inputs")
	}
	if len(outs) < 1 {
		return "", fmt.Errorf("no outputs")
	}
	for _, sat := range outs {
		if sat < minTxAmount {
			return "", fmt.Errorf("min allow amount %d satoshis", minTxAmount)
		}
	}

	// 2. in
	type inType struct {
		TxId string `json:"txid"`
		VOut uint32 `json:"vout"`
		// Seq  uint32 `json:"sequence"`
	}
	var inData []inType
	for _, item := range ins {
		if len(item.Tx) != 64 {
			return "", fmt.Errorf("error transaction: %s", item.Tx)
		}
		inData = append(inData, inType{TxId: item.Tx, VOut: item.N})
	}

	// 3. out
	outData := make(map[string]interface{})
	for addr, sat := range outs {
		n := decimal.NewFromInt(sat).Div(decimal.NewFromInt(100000000)) // ou.Amt / 10^8
		outData[addr] = n.String()
	}
	if len(hexData) > 0 {
		outData["data"] = hexData
	}

	// 4.
	param := []interface{}{inData, outData}
	b, err := rc.Request("createrawtransaction", param)
	if err != nil {
		return "", err
	}

	// 5. parse
	type desc struct {
		Result string      `json:"result"`
		Error  interface{} `json:"error"`
		ID     string      `json:"id"`
	}
	var resp desc
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("error return, reqId: %s", resp.ID)
	}

	return resp.Result, nil
}

// SignRawTX
// https://developer.bitcoin.org/reference/rpc/signrawtransactionwithkey.html
// priKeys: base58-encoded private keys
func (rc *RpcClient) SignRawTX(hex string, priKeys []string) (string, error) {

	// 1.
	param := []interface{}{hex, priKeys}
	b, err := rc.Request("signrawtransactionwithkey", param)
	if err != nil {
		return "", err
	}

	// 2.
	type desc struct {
		Result struct {
			Hex      string      `json:"hex"`
			Complete bool        `json:"complete"`
			Errors   interface{} `json:"errors"`
		} `json:"result"`
		Error interface{} `json:"error"`
		ID    string      `json:"id"`
	}
	var resp desc
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return "", err
	}
	if resp.Result.Complete != true {
		return "", fmt.Errorf("error response: %s", string(b))
	}

	return resp.Result.Hex, err
}

// SendRawTX
// https://developer.bitcoin.org/reference/rpc/sendrawtransaction.html
// 失败: 报错500, 可使用命令行测试问题
// 格式: bitcoin-cli sendrawtransaction <signedHex>
// ---------------
// 报错1: Fee exceeds maximum configured by user (e.g. -maxtxfee, maxfeerate)
// 报错2: min relay fee not met, 100 < 141
// 原因: 手续费用太高或太低
func (rc *RpcClient) SendRawTX(signedHex string) (tx string, error error) {

	param := []interface{}{signedHex}
	b, err := rc.Request("sendrawtransaction", param)
	if err != nil {
		return "", fmt.Errorf("%s; it maybe invalid fee or txn-mempool-conflict; try test command `bitcoin-cli sendrawtransaction <signedHex>`", err.Error())
	}

	type desc struct {
		Result string      `json:"result"`
		Error  interface{} `json:"error"`
		ID     string      `json:"id"`
	}
	var resp desc
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return "", err
	}

	return resp.Result, nil
}

// GetRawTransaction
// @docs https://developer.bitcoin.org/reference/rpc/getrawtransaction.html
func (rc *RpcClient) GetRawTransaction(txId string, blockHash string) (*RawTransaction, error) {

	// bitcoin-cli getrawtransaction "mytxid" true "myblockhash"
	// bitcoin-cli getrawtransaction "mytxid" false "myblockhash" // 返回 hexString 可使用 decoderawtransaction 解析
	args := []interface{}{txId, true, blockHash}
	if blockHash == "" {
		args = []interface{}{txId, true}
	}
	b, err := rc.Request("getrawtransaction", args)
	if err != nil {
		return nil, err
	}

	type response struct {
		Response
		Result RawTransaction
	}
	var resp response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &resp.Result, nil
}
