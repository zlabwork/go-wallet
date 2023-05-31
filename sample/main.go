package main

import (
	"fmt"
	"github.com/zlabwork/go-wallet/btc"
)

func main() {
	priKey := btc.NewPriKeyRandom()
	fmt.Println(priKey.PubKey().Address().P2pkh())
	fmt.Println(priKey.PubKey().Address().P2sh())
	fmt.Println(priKey.PubKey().Address().P2wpkh())
	fmt.Println(priKey.PubKey().Address().P2wsh())
	fmt.Println(priKey.PubKey().Address().P2wpkhInP2sh())
	fmt.Println(priKey.PubKey().Address().MultiP2wsh())
	fmt.Println(priKey.PubKey().Address().MultiP2wshInP2sh())
}
