package btc

// GetBlockHash
// https://developer.bitcoin.org/reference/rpc/getblockhash.html
func (rc *RpcClient) GetBlockHash(blockHeight int64) (string, error) {
	b, err := rc.Request("getblockhash", []interface{}{blockHeight})
	if err != nil {
		return "", err
	}
	return string(b[11:75]), nil
}

// GetBlock
// https://developer.bitcoin.org/reference/rpc/getblock.html
func (rc *RpcClient) GetBlock(blockHeight int64) ([]byte, error) {
	// 1. 获取块hash
	h, err := rc.GetBlockHash(blockHeight)
	if err != nil {
		return nil, err
	}

	// 2. 获取块数据
	b, err := rc.Request("getblock", []interface{}{h, 2})
	if err != nil {
		return nil, err
	}

	return b, nil
}
