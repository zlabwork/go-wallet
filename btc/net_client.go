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
