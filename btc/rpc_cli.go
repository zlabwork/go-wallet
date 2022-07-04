package btc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zlabwork/go-zlibs"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	minTxAmount = 546 // satoshis
)

type RpcClient struct {
	req      *zlibs.HttpLib
	auth     string
	endpoint string
}

type HandleConfigs struct {
	Host string // http://127.0.0.1:18443
	User string
	Pass string
}

func NewRpcClient(c *HandleConfigs) *RpcClient {
	return &RpcClient{
		req:      zlibs.NewHttpLib(),
		auth:     "Basic " + base64.StdEncoding.EncodeToString([]byte(c.User+":"+c.Pass)),
		endpoint: c.Host,
	}
}

func (rc *RpcClient) Request(method string, params []interface{}) ([]byte, error) {

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
	header["Authorization"] = rc.auth
	rc.req.SetHeaders(header)
	resp, err := rc.req.RequestRaw("POST", rc.endpoint, []byte(data))
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
