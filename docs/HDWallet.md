## HD Wallet Example

```golang
package main

import (
    "encoding/hex"
    "fmt"
    "github.com/tyler-smith/go-bip39"
    "github.com/zlabwork/go-bip32"
    "github.com/zlabwork/go-wallet/btc"
    "github.com/zlabwork/go-wallet/docs"
)

// Legacy: m/44'/0'/0'
// SegWit: m/49'/0'/0'
func main() {

    // 1. 生成助记词
    bip39.SetWordList(docs.WordList().ChineseSimple())
    entropy, _ := bip39.NewEntropy(128)
    mnemonic, _ := bip39.NewMnemonic(entropy)
    fmt.Println("助记词:", mnemonic)

    // 2. 验证助记词
    if !bip39.IsMnemonicValid(mnemonic) {
        panic("Valid Mnemonic")
    }

    // 3. 种子
    // a> 随机种子
    //seed, err := bip32.NewSeed()
    //if err != nil {
    //    log.Fatalln("Error generating seed:", err)
    //}

    // b> 助记词种子
    // 注: 多数HD钱包使用空密码
    seed := bip39.NewSeed(mnemonic, "") // custom passphrase
    fmt.Println("种子:", hex.EncodeToString(seed))
    fmt.Println("\n----------------------")

    // 4. 主秘钥
    masterKey, _ := bip32.NewMasterKey(seed)
    publicKey := masterKey.PublicKey()
    fmt.Println("root xprv:", masterKey)
    fmt.Println("root xpub:", publicKey)

    // 4. 子秘钥
    key, _ := masterKey.NewChildKey(uint32(0x80000000) + 44) // m/44'
    key, _ = key.NewChildKey(uint32(0x80000000))             // m/44'/0'
    key, _ = key.NewChildKey(uint32(0x80000000))             // m/44'/0'/0'

    // @see https://iancoleman.io/bip39/
    fmt.Println("\n--------- m/44'/0'/0' ---------")
    fmt.Println("account xprv:", key.String())
    fmt.Println("account xpub:", key.PublicKey().String())

    key, _ = key.NewChildKey(uint32(0)) // m/44'/0'/0'/0
    key, _ = key.NewChildKey(uint32(0)) // m/44'/0'/0'/0/0

    fmt.Println("\n--------- m/44'/0'/0'/0/0 ---------")
    fmt.Println("私钥:", hex.EncodeToString(key.Key))
    fmt.Println("公钥:", hex.EncodeToString(key.PublicKey().Key))
    pubKey, _ := btc.NewPubKey(key.PublicKey().Key)
    fmt.Println("地址:", pubKey.Address().P2pkh())
}
```
