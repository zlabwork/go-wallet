package btc

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/zlabwork/go-zlibs"
    "io/ioutil"
    "math"
    "strconv"
    "strings"
    "time"
)

const (
    minTxAmount = 546 // satoshis
)

type ServiceClient struct {
    req      *zlibs.HttpLib
    auth     string
    endpoint string
}

func NewServiceClient(handle *serviceHandle) *ServiceClient {
    return &ServiceClient{
        req:      handle.req,
        auth:     "Basic " + base64.StdEncoding.EncodeToString([]byte(handle.cfg.User+":"+handle.cfg.Pass)),
        endpoint: handle.cfg.Host,
    }
}

func (sc *ServiceClient) Request(method string, params []interface{}) ([]byte, error) {

    reqId := strconv.FormatInt(time.Now().UnixNano(), 10)
    type rd struct {
        Ver    string        `json:"jsonrpc"`
        Id     string        `json:"id"`
        Method string        `json:"method"`
        Params []interface{} `json:"params"`
    }
    b, err := json.Marshal(rd{
        Ver:    "1.0",
        Id:     reqId,
        Method: method,
        Params: params,
    })
    if err != nil {
        return nil, err
    }

    // setting header & body
    data := string(b)
    header := make(map[string]string)
    header["Authorization"] = sc.auth
    sc.req.SetHeaders(header)
    resp, err := sc.req.RequestRaw("POST", sc.endpoint, []byte(data))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // http body
    body, err := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("code: %d, message: %s", resp.StatusCode, strings.Trim(string(body), "\n"))
    }
    return body, err
}

// https://developer.bitcoin.org/reference/rpc/getblockhash.html
func (sc *ServiceClient) GetBlockHash(blockHeight int64) (string, error) {
    b, err := sc.Request("getblockhash", []interface{}{blockHeight})
    if err != nil {
        return "", err
    }
    return string(b[11:75]), nil
}

// https://developer.bitcoin.org/reference/rpc/getblock.html
func (sc *ServiceClient) GetBlock(blockHeight int64) ([]byte, error) {
    // 1. 获取块hash
    h, err := sc.GetBlockHash(blockHeight)
    if err != nil {
        return nil, err
    }

    // 2. 获取块数据
    b, err := sc.Request("getblock", []interface{}{h, 2})
    if err != nil {
        return nil, err
    }

    return b, nil
}

// https://developer.bitcoin.org/reference/rpc/gettxout.html
func (sc *ServiceClient) GetTxOut(tx string, index int) ([]byte, error) {

    b, err := sc.Request("gettxout", []interface{}{tx, index})
    if err != nil {
        return nil, err
    }
    return b, nil
}

func (sc *ServiceClient) CreateTX(ins []VIn, outs []VOut, hexData string, sat int64, chargeBack string) (hex string, error error) {

    // 1. total in
    var totalIn int64
    for _, in := range ins {
        if in.Amt < minTxAmount {
            return "", fmt.Errorf("min allow amount %d satoshis", minTxAmount)
        }
        totalIn += in.Amt
    }

    // 2. total out
    var totalOut int64
    for _, out := range outs {
        if out.Amt < minTxAmount {
            return "", fmt.Errorf("min allow amount %d satoshis", minTxAmount)
        }
        totalOut += out.Amt
    }

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

    return sc.CreateRawTX(ins, outs, hexData)
}

// https://developer.bitcoin.org/reference/rpc/createrawtransaction.html
func (sc *ServiceClient) CreateRawTX(ins []VIn, outs []VOut, hexData string) (hex string, error error) {

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
        outData[ou.Addr] = float64(ou.Amt) * math.Pow10(-8)
    }
    if len(hexData) > 0 {
        outData["data"] = hexData
    }
    param := []interface{}{inData, outData}

    // 4.
    b, err := sc.Request("createrawtransaction", param)
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

// https://developer.bitcoin.org/reference/rpc/signrawtransactionwithkey.html
// priKeys: base58-encoded private keys
func (sc *ServiceClient) SignRawTX(hex string, priKeys []string) (string, error) {

    // 1.
    param := []interface{}{hex, priKeys}
    b, err := sc.Request("signrawtransactionwithkey", param)
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

// https://developer.bitcoin.org/reference/rpc/sendrawtransaction.html
// 失败: 报错500, 可使用命令行测试问题
// 格式: bitcoin-cli sendrawtransaction <signedHex>
// ---------------
// 报错1: Fee exceeds maximum configured by user (e.g. -maxtxfee, maxfeerate)
// 报错2: min relay fee not met, 100 < 141
// 原因: 手续费用太高或太低
func (sc *ServiceClient) SendRawTX(signedHex string) (tx string, error error) {

    param := []interface{}{signedHex}
    b, err := sc.Request("sendrawtransaction", param)
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
