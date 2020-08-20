## 安装
```bash
go get github.com/zlabwork/libschain
# or use this lib with go.mod
```


## 使用
```golang
// 创建账户
lib := libschain.NewEthLib()
account, _ := lib.CreateAccount()

// 测试网络
lib.Connect("http://127.0.0.1:8545")

// 正式网络 
// https://infura.io/ 需要申请一个KEY
// lib.Connect("https://mainnet.infura.io/v3/xxxxxxxx")
// lib.Connect("wss://mainnet.infura.io/ws/v3/xxxxxxxx")
```
