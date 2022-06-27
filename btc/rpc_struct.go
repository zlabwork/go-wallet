package btc

type Response struct {
	Result interface{}
	Error  string
	Id     string
}

type Block struct {
	Hash              string        `json:"hash"`
	Confirmations     int           `json:"confirmations"`
	Strippedsize      int           `json:"strippedsize"`
	Size              int           `json:"size"`
	Weight            int           `json:"weight"`
	Height            int           `json:"height"`
	Version           int           `json:"version"`
	VersionHex        string        `json:"versionHex"`
	Merkleroot        string        `json:"merkleroot"`
	Tx                []Transaction `json:"tx"`
	Time              int           `json:"time"`
	Mediantime        int           `json:"mediantime"`
	Nonce             int           `json:"nonce"`
	Bits              string        `json:"bits"`
	Difficulty        float64       `json:"difficulty"`
	Chainwork         string        `json:"chainwork"`
	NTx               int           `json:"nTx"`
	Previousblockhash string        `json:"previousblockhash"`
	Nextblockhash     string        `json:"nextblockhash"`
}

type Transaction struct {
	Txid     string `json:"txid"`
	Hash     string `json:"hash"`
	Version  int    `json:"version"`
	Size     int    `json:"size"`
	Vsize    int    `json:"vsize"`
	Weight   int    `json:"weight"`
	Locktime int    `json:"locktime"`
	Vin      []Vin  `json:"vin"`
	Vout     []Vout `json:"vout"`
	Hex      string `json:"hex"`
}

type Vin struct {
	Coinbase  string `json:"coinbase"`
	Txid      string `json:"txid"`
	Vout      int    `json:"vout"`
	ScriptSig struct {
		Asm string `json:"asm"`
		Hex string `json:"hex"`
	} `json:"scriptSig"`
	Txinwitness []string `json:"txinwitness"`
	Sequence    int64    `json:"sequence"`
}

type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	Type      string   `json:"type"`
	ReqSigs   int      `json:"reqSigs"`
	Addresses []string `json:"addresses"`
}

type txOut struct {
	Result struct {
		Bestblock     string       `json:"bestblock"`
		Confirmations int          `json:"confirmations"`
		Value         float64      `json:"value"`
		ScriptPubKey  ScriptPubKey `json:"scriptPubKey"`
		Coinbase      bool         `json:"coinbase"`
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    string      `json:"id"`
}
