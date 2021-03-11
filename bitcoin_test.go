package main

import (
    "encoding/hex"
    "github.com/zlabwork/go-chain/bitcoin"
    "testing"
)

const (
    btcPriKey          = "DFB9E60F61CC1EE2CFCAAD7E9C7187121C9B0C21FD66E87D3BC32168AC14FFFF"
    btcPriWif          = "L4ic4Xh9a7nJgvChwYLSBLL3guBDmJkJbZ72F5prD8TQkxATTMxk"
    btcAddr1Compress   = "1HoYi6T28GBVU652SFfTCvrSf51wFQ9qvY"
    btcAddr1UnCompress = "1PKX2i36ab9AAo2hCE6gkCVQQ28UYzeUdm"
    btcAddr3Compress   = "3JVZddwTgAVsZFmTZML3dZDNobJenLjq3G"
    btcAddr3UnCompress = "3Q1XxFXY8VTYFxj8KKmHAprLYYRC6VUswU"
)

func TestPriKeyWif(t *testing.T) {
    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    if priKey.WIF() != btcPriWif {
        t.Errorf("priKey.WIF() %s", priKey.WIF())
    }
}

func TestPubKeyCompressed(t *testing.T) {

    priKey := bitcoin.NewPriKeyRandom()
    p1 := priKey.PubKey()
    p2, _ := bitcoin.NewPubKey(p1.Key())

    if hex.EncodeToString(p1.Key()) != hex.EncodeToString(p2.Key()) {
        t.Error("method priKey.PubKey() and NewPubKey() not matched ")
    }
}

func TestPubKeyUncompressed(t *testing.T) {

    priKey := bitcoin.NewPriKeyRandom()
    p1 := priKey.PubKeyUncompressed()
    p2, _ := bitcoin.NewPubKey(p1.Key())

    if hex.EncodeToString(p1.Key()) != hex.EncodeToString(p2.Key()) {
        t.Error("method priKey.PubKeyUncompressed() and NewPubKey() not matched ")
    }
}

func TestAddressP2PKH(t *testing.T) {

    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    pub1 := priKey.PubKey()
    pub2 := priKey.PubKeyUncompressed()
    if pub1.Address().P2PKH() != btcAddr1Compress {
        t.Error("error compress address when P2PKH")
    }
    if pub2.Address().P2PKH() != btcAddr1UnCompress {
        t.Error("error uncompress address when P2PKH")
    }
}

func TestAddressP2SH(t *testing.T) {

    priKey, _ := bitcoin.ParseWIF(btcPriWif)
    pub1 := priKey.PubKey()
    pub2 := priKey.PubKeyUncompressed()
    if pub1.Address().P2SH() != btcAddr3Compress {
        t.Error("error compress address when P2SH")
    }
    if pub2.Address().P2SH() != btcAddr3UnCompress {
        t.Error("error uncompress address when P2SH")
    }
}
