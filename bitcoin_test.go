package main

import (
    "encoding/hex"
    "github.com/zlabwork/go-chain/bitcoin"
    "testing"
)

const (
    testPriKey          = "DFB9E60F61CC1EE2CFCAAD7E9C7187121C9B0C21FD66E87D3BC32168AC14FFFF"
    testPriWif          = "L4ic4Xh9a7nJgvChwYLSBLL3guBDmJkJbZ72F5prD8TQkxATTMxk"
    testAddr1Compress   = "1HoYi6T28GBVU652SFfTCvrSf51wFQ9qvY"
    testAddr1UnCompress = "1PKX2i36ab9AAo2hCE6gkCVQQ28UYzeUdm"
    testAddr3Compress   = "3JVZddwTgAVsZFmTZML3dZDNobJenLjq3G"
    testAddr3UnCompress = "3Q1XxFXY8VTYFxj8KKmHAprLYYRC6VUswU"
)

func TestPriKeyWif(t *testing.T) {
    priKey, _ := bitcoin.ParseWIF(testPriWif)
    if priKey.WIF() != testPriWif {
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

    priKey, _ := bitcoin.ParseWIF(testPriWif)
    pub1 := priKey.PubKey()
    pub2 := priKey.PubKeyUncompressed()
    if pub1.Address().P2PKH() != testAddr1Compress {
        t.Error("error compress address when P2PKH")
    }
    if pub2.Address().P2PKH() != testAddr1UnCompress {
        t.Error("error uncompress address when P2PKH")
    }
}

func TestAddressP2SH(t *testing.T) {

    priKey, _ := bitcoin.ParseWIF(testPriWif)
    pub1 := priKey.PubKey()
    pub2 := priKey.PubKeyUncompressed()
    if pub1.Address().P2SH() != testAddr3Compress {
        t.Error("error compress address when P2SH")
    }
    if pub2.Address().P2SH() != testAddr3UnCompress {
        t.Error("error uncompress address when P2SH")
    }
}
