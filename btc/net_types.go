package btc

type txOut struct {
	Result struct {
		Bestblock     string  `json:"bestblock"`
		Confirmations int     `json:"confirmations"`
		Value         float64 `json:"value"`
		ScriptPubKey  struct {
			Asm       string   `json:"asm"`
			Hex       string   `json:"hex"`
			ReqSigs   int      `json:"reqSigs"`
			Type      string   `json:"type"`
			Addresses []string `json:"addresses"`
		} `json:"scriptPubKey"`
		Coinbase bool `json:"coinbase"`
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    string      `json:"id"`
}
