package main

import (
    "encoding/hex"
    "github.com/zlabwork/go-chain/bitcoin"
    "testing"
)

const (
    btcPriKey              = "DFB9E60F61CC1EE2CFCAAD7E9C7187121C9B0C21FD66E87D3BC32168AC14FFFF"
    btcPriWif              = "L4ic4Xh9a7nJgvChwYLSBLL3guBDmJkJbZ72F5prD8TQkxATTMxk"
    btcAddrUnCompressP2PKH = "1PKX2i36ab9AAo2hCE6gkCVQQ28UYzeUdm"
    btcAddrUnCompressP2SH  = "3Q1XxFXY8VTYFxj8KKmHAprLYYRC6VUswU"
    btcAddrP2PKH           = "1HoYi6T28GBVU652SFfTCvrSf51wFQ9qvY"
    btcAddrP2SH            = "3JVZddwTgAVsZFmTZML3dZDNobJenLjq3G"
    btcAddrP2WPKH          = "bc1qhp8em87yulx7dyuzy4f76svtc0h89n2vpkx5px"
    btcAddrP2WSH           = "bc1qr0rsjgmc92uhymwr6eddhv8u7y006jg4f8aa5tr5pnhc2resjn7q5kf3l2"
    btcAddrP2SHP2WPKH      = "3Gnnj3oozxwWs4FHpbnoy8BVJFvjB42Lj3"
)

func TestBTC_PriKeyWif(t *testing.T) {
    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    if priKey.WIF() != btcPriWif {
        t.Errorf("priKey.WIF() %s", priKey.WIF())
    }
}

func TestBTC_Compressed(t *testing.T) {

    priKey := bitcoin.NewPriKeyRandom()
    p1 := priKey.PubKey()
    p2, _ := bitcoin.NewPubKey(p1.Key())

    if hex.EncodeToString(p1.Key()) != hex.EncodeToString(p2.Key()) {
        t.Error("method priKey.PubKey() and NewPubKey() not matched ")
    }
}

func TestBTC_Uncompressed(t *testing.T) {

    priKey := bitcoin.NewPriKeyRandom()
    p1 := priKey.PubKeyUncompressed()
    p2, _ := bitcoin.NewPubKey(p1.Key())

    if hex.EncodeToString(p1.Key()) != hex.EncodeToString(p2.Key()) {
        t.Error("method priKey.PubKeyUncompressed() and NewPubKey() not matched ")
    }
}

func TestBTC_P2PKH(t *testing.T) {

    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    pub1 := priKey.PubKey()
    pub2 := priKey.PubKeyUncompressed()
    if pub1.Address().P2PKH() != btcAddrP2PKH {
        t.Error("error compress address when P2PKH")
    }
    if pub2.Address().P2PKH() != btcAddrUnCompressP2PKH {
        t.Error("error uncompress address when P2PKH")
    }
}

func TestBTC_P2SH(t *testing.T) {

    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    pub1 := priKey.PubKey()
    pub2 := priKey.PubKeyUncompressed()
    if pub1.Address().P2SH() != btcAddrP2SH {
        t.Error("error compress address when P2SH")
    }
    if pub2.Address().P2SH() != btcAddrUnCompressP2SH {
        t.Error("error uncompress address when P2SH")
    }
}

func TestBTC_P2WPKH(t *testing.T) {

    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    pub := priKey.PubKey()
    addr, _ := pub.Address().P2WPKH()
    if addr != btcAddrP2WPKH {
        t.Error("error address P2WPKH")
    }
}

func TestBTC_P2WSH(t *testing.T) {

    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    pub := priKey.PubKey()
    addr, _ := pub.Address().P2WSH()
    if addr != btcAddrP2WSH {
        t.Error("error address P2WSH")
    }
}

func TestBTC_P2SHP2WPKH(t *testing.T) {

    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    pub := priKey.PubKey()
    if pub.Address().P2SHP2WPKH() != btcAddrP2SHP2WPKH {
        t.Error("error address P2SHP2WPKH")
    }
}
