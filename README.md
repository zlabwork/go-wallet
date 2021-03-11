## 安装
```bash
go get github.com/zlabwork/go-chain

go get github.com/zlabwork/go-chain@v1.1.0
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

## 钱包工具
* http://webhdwallet.github.io/  
* https://iancoleman.io/bip39/  
* https://www.bitaddress.org/  


## docs
[椭圆曲线文档](http://www.secg.org/sec2-v2.pdf)  
[椭圆曲线图形](https://www.desmos.com/calculator/ialhd71we3?lang=zh-CN)  
[Graphical Address Generator](https://www.royalfork.org/2014/08/11/graphical-address-generator)  
[BCH地址规则](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md)  
[https://github.com/gcash/bchutil](https://github.com/gcash/bchutil)  
[bip-0044](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)  
