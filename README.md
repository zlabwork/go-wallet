## 安装
```bash
go get github.com/zlabwork/go-chain

go get github.com/zlabwork/go-chain@v1.3.0
```
## BTC
```golang
import "github.com/zlabwork/go-chain/btc"

priKey := btc.NewPriKeyRandom()
priKey.WIF() // L211iZmidtxLQ2s7hzM9BYacPUu2asT1KkCkyrTbNbDib2N85ai5
pubKey := priKey.PubKey()
address1 := pubKey.Address().P2PKH() // 19c4pkCL2jvTFYkZXDyUHi4ceoNze44mXE
address2 := pubKey.Address().P2SH()  // 3AJ5kHgmaeEqLiSzeKe4iLRYoKfiCH5Y1C
```


## BCH
```golang
import "github.com/zlabwork/go-chain/bch"

priKey := bch.NewPriKeyRandom()
address := priKey.PubKey().Address().P2PKH()
log.Println(address) // qzz6eq5we2qdxg29jkzxkxafc34xhduk7vhayz3z06
```


## ETH
```golang
import "github.com/zlabwork/go-chain/eth"

priKey := eth.NewPriKeyRandom()
priKey.Address().String()

// 测试网络
cs, _ := eth.NewConnectService("http://127.0.0.1:8545")
cs.GetBalance("0x06****11")

// 正式网络 
// https://infura.io/ 需要申请一个KEY
// eth.NewConnectService("https://mainnet.infura.io/v3/xxxxxxxx")
// eth.NewConnectService("wss://mainnet.infura.io/ws/v3/xxxxxxxx")
```


## HD Wallet Example
[HD Wallet Example](docs/HDWallet.md)  


## 钱包工具
* http://webhdwallet.github.io/  
* https://iancoleman.io/bip39/  
* https://www.bitaddress.org/  


## Docs
##### 1. articles
[椭圆曲线文档](http://www.secg.org/sec2-v2.pdf)  
[Graphical Address Generator](https://www.royalfork.org/2014/08/11/graphical-address-generator)  

##### 2. bch
[BCH地址规则](https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md)  
[https://github.com/gcash/bchutil](https://github.com/gcash/bchutil)


##### 3. wallet
[bip-0044](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)  
[钱包兼容性汇总](https://bitcoinops.org/en/compatibility/)  

##### 4. transactions
[btc transactions](https://developer.bitcoin.org/reference/transactions.html)  
[btc Script](https://en.bitcoin.it/wiki/Script)  


## Tools
[椭圆曲线图形](https://www.desmos.com/calculator/ialhd71we3?lang=zh-CN)  
[SegWit address format](http://bitcoin.sipa.be/bech32/demo/demo.html)  


## Similar Tools 
[python bitcoin-utils](https://pypi.org/project/bitcoin-utils/)  


## Bitcoin Transaction
![bitcoin-tx](docs/assets/bitcoin-tx.png)
