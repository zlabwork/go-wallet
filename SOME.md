## 合约



## 合约相关

```bash
# 生成绑定文件 token.go
abigen --abi token.abi --pkg mytoken --type Token --out token.go

# 生成部署文件 token.go
abigen --abi token.abi --pkg mytoken --type Token --out token.go --bin token.bin
```
```bash
# 如果安装了 solc 可以直接绑定
abigen --sol token.sol --pkg mytoken --out token.go
```

[Reference Docs](https://geth.ethereum.org/docs/dapp/native-bindings)  
[Chainlink: LINK Token](https://cn.etherscan.com/address/0x514910771af9ca656af840dff83e8264ecf986ca#code)  
