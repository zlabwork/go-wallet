## 安装
```bash
go get github.com/zlabwork/go-chain
# or use this lib with go.mod
```
## bitcoin
```golang
import "github.com/zlabwork/go-chain/bitcoin"

priKey := bitcoin.NewPriKeyRandom()
priKey.WIF() // L211iZmidtxLQ2s7hzM9BYacPUu2asT1KkCkyrTbNbDib2N85ai5
pubKey := priKey.PubKey()
address1 := pubKey.Address().P2PKH() // 19c4pkCL2jvTFYkZXDyUHi4ceoNze44mXE
address2 := pubKey.Address().P2SH()  // 3AJ5kHgmaeEqLiSzeKe4iLRYoKfiCH5Y1C
```


## bitcoincash
```golang
priKey := bitcoincash.NewPriKeyRandom()
address := priKey.PubKey().Address().P2PKH()
log.Println(address) // qzz6eq5we2qdxg29jkzxkxafc34xhduk7vhayz3z06
```


## ethereum
```golang
import "github.com/zlabwork/go-chain/ethereum"

lib := ethereum.NewEthLib()
priKey, _ := lib.GenPriKey()
address, _ := lib.GetAddress(priKey)

// 测试网络
lib.Connect("http://127.0.0.1:8545")

// 正式网络 
// https://infura.io/ 需要申请一个KEY
// lib.Connect("https://mainnet.infura.io/v3/xxxxxxxx")
// lib.Connect("wss://mainnet.infura.io/ws/v3/xxxxxxxx")
```

## docs
[BCH地址规则](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md)
[https://github.com/gcash/bchutil](https://github.com/gcash/bchutil)