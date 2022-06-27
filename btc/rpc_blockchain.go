package btc

import (
	"encoding/json"
	"errors"
)

// GetBlockHash
// https://developer.bitcoin.org/reference/rpc/getblockhash.html
func (rc *RpcClient) GetBlockHash(blockHeight int64) (string, error) {
	b, err := rc.Request("getblockhash", []interface{}{blockHeight})
	if err != nil {
		return "", err
	}
	// {"result":"3290a295cbb8e27c845f24699ba20f161743a57aa3960006113b16fbdf6e5b73","error":null,"id":"1656317255927154700"}
	return string(b[11:75]), nil
}

// GetBlock
// https://developer.bitcoin.org/reference/rpc/getblock.html
func (rc *RpcClient) GetBlock(hashString string) (*Block, error) {

	// 1. request rpc
	b, err := rc.Request("getblock", []interface{}{hashString, 2})
	if err != nil {
		return nil, err
	}

	// 2. json decode
	type response struct {
		Response
		Result Block
	}
	var resp response
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	// 3. return
	return &resp.Result, nil
}
