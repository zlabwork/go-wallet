## 安装
```bash
go get github.com/zlabwork/go-chain
# or use this lib with go.mod
```
## bitcoin
```golang
import "github.com/zlabwork/go-chain/bitcoin"

priKey := bitcoin.GeneratePrivateKey()
pubKey := bitcoin.GeneratePublicKey(priKey)
address1 := bitcoin.P2PKH(pubKey)
address2 := bitcoin.P2SH(pubKey)
```

## ethereum
```golang
import "github.com/zlabwork/go-chain/ethereum"

lib := ethereum.NewEthLib()
priKey, _ := lib.GeneratePrivateKey()
address, _ := lib.GetAddress(priKey)

// 测试网络
lib.Connect("http://127.0.0.1:8545")

// 正式网络 
// https://infura.io/ 需要申请一个KEY
// lib.Connect("https://mainnet.infura.io/v3/xxxxxxxx")
// lib.Connect("wss://mainnet.infura.io/ws/v3/xxxxxxxx")
```
