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
```
