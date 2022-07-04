package btc

import (
	"encoding/json"
)

// 获取缓存池描述
// @doc https://developer.bitcoin.org/reference/rpc/getmempoolinfo.html
// [root@space ~]# bitcoin-cli getmempoolinfo
// {
//   "loaded": true,
//   "size": 188,
//   "bytes": 65981,
//   "usage": 297456,
//   "maxmempool": 300000000,
//   "mempoolminfee": 0.00001000,
//   "minrelaytxfee": 0.00001000,
//   "unbroadcastcount": 0
// }

// GetRawMemPool
// @doc https://developer.bitcoin.org/reference/rpc/getrawmempool.html
// @ref https://developer.bitcoin.org/reference/rpc/getmempoolentry.html
// @ref https://developer.bitcoin.org/reference/rpc/getmempoolinfo.html
func (rc *RpcClient) GetRawMemPool() ([]string, error) {
	b, err := rc.Request("getrawmempool", []interface{}{false, true})
	if err != nil {
		return nil, err
	}

	type typeTxIds struct {
		TxIds      []string `json:"txids"`
		MemPoolSeq int64    `json:"mempool_sequence"`
	}
	type response struct {
		Response
		Result typeTxIds
	}
	var resp response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp.Result.TxIds, err
}
