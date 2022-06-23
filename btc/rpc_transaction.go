package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// GetTxOut
// https://developer.bitcoin.org/reference/rpc/gettxout.html
func (rc *RpcClient) GetTxOut(tx string, index int) (*txOut, error) {

	b, err := rc.Request("gettxout", []interface{}{tx, index})
	if err != nil {
		return nil, err
	}
	var out txOut
	if err = json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	if out.Result.Value == 0 {
		return nil, fmt.Errorf("%s, error txId or tx has been spent", tx)
	}
	return &out, nil
}

func (rc *RpcClient) CreateTransferAll(ins map[string]uint32, addr string, sat int64) (hex string, error error) {

	var t int64
	var inT []VIn
	for tx, n := range ins {
		ou, err := rc.GetTxOut(tx, int(n))
		if err != nil {
			return "", err
		}
		v := int64(ou.Result.Value * math.Pow10(8))
		if v < minTxAmount {
			return "", fmt.Errorf("%s current %d satoshis, less than minimum amount %d satoshis", tx, v, minTxAmount)
		}
		t += v
		inT = append(inT, VIn{Tx: tx, N: n})
	}

	// fees
	size := 148*len(ins) + 34 + 10
	fee := int64(size) * sat

	return rc.createRawTX(inT, []VOut{{Addr: addr, Amt: t - fee}}, "")
}

// CreateTXAlias - Alias for CreateTX
func (rc *RpcClient) CreateTXAlias(ins []string, outs map[string]int64, hexData string, sat int64, chargeBack string) (hex string, error error) {
	var inT []VIn
	var outT []VOut
	for _, str := range ins {
		s := strings.Split(str, ":")
		n, err := strconv.ParseUint(strings.TrimSpace(s[1]), 10, 64)
		if err != nil {
			return "", errors.New(str + ", error vout format")
		}
		inT = append(inT, VIn{Tx: strings.TrimSpace(s[0]), N: uint32(n)})
	}
	for ad, n := range outs {
		outT = append(outT, VOut{Addr: ad, Amt: n})
	}
	return rc.CreateTX(inT, outT, hexData, sat, chargeBack)
}

// CreateTX
// outs := []VOut{{Addr: "btc address 2", Amt: 1000}, {Addr: "btc address 1", Amt: 2000}}
func (rc *RpcClient) CreateTX(ins []VIn, outs []VOut, hexData string, sat int64, chargeBack string) (hex string, error error) {

	// 1. total in
	var totalIn int64
	for _, in := range ins {
		txOut, err := rc.GetTxOut(in.Tx, int(in.N))
		if err != nil {
			return "", err
		}
		amt := int64(txOut.Result.Value * math.Pow10(8))
		if amt < minTxAmount {
			return "", fmt.Errorf("current %d satoshis, less than minimum amount %d satoshis", amt, minTxAmount)
		}
		totalIn += amt
	}

	// 2. total out
	var totalOut int64
	for _, out := range outs {
		if out.Amt < minTxAmount {
			return "", fmt.Errorf("transfer %d satoshis, less than minimum amount %d satoshis", out.Amt, minTxAmount)
		}
		totalOut += out.Amt
	}

	// TODO: 1. how to calculate if chargeBack in the outs list 2. confirm the calculate process
	// 3. fee sat
	size := 148*len(ins) + 34*len(outs) + 10
	fee := int64(size) * sat
	left := totalIn - totalOut - fee
	if left < 0 {
		return "", fmt.Errorf("fee is not enough")
	}

	// 4. charge back
	if left-34*sat > minTxAmount {
		outs = append(outs, VOut{Addr: chargeBack, Amt: left - 34*sat})
		fee += 34 * sat
	}

	return rc.createRawTX(ins, outs, hexData)
}

// CreateRawTX - TODO: Fix the address duplicated in VOut
// https://developer.bitcoin.org/reference/rpc/createrawtransaction.html
func (rc *RpcClient) createRawTX(ins []VIn, outs []VOut, hexData string) (hex string, error error) {

	// 1. check
	if len(ins) < 1 {
		return "", fmt.Errorf("no inputs")
	}
	if len(outs) < 1 {
		return "", fmt.Errorf("no outputs")
	}
	for _, t := range outs {
		if t.Amt < minTxAmount {
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
	for _, ou := range outs {
		outData[ou.Addr] = float64(ou.Amt) / 100000000 // FIXME: 浮点精度问题
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
