package btc

import (
	"encoding/hex"
	"testing"
)

const (
	btcPriKey               = "DFB9E60F61CC1EE2CFCAAD7E9C7187121C9B0C21FD66E87D3BC32168AC14FFFF"
	btcPriWif               = "L4ic4Xh9a7nJgvChwYLSBLL3guBDmJkJbZ72F5prD8TQkxATTMxk"
	btcAddrUnCompressP2pkh  = "1PKX2i36ab9AAo2hCE6gkCVQQ28UYzeUdm"
	btcAddrUnCompressP2sh   = "3Q1XxFXY8VTYFxj8KKmHAprLYYRC6VUswU"
	btcAddrP2pkh            = "1HoYi6T28GBVU652SFfTCvrSf51wFQ9qvY"
	btcAddrP2sh             = "3JVZddwTgAVsZFmTZML3dZDNobJenLjq3G"
	btcAddrP2wpkh           = "bc1qhp8em87yulx7dyuzy4f76svtc0h89n2vpkx5px"
	btcAddrP2wsh            = "bc1qr0rsjgmc92uhymwr6eddhv8u7y006jg4f8aa5tr5pnhc2resjn7q5kf3l2"
	btcAddrP2wpkhInP2sh     = "3Gnnj3oozxwWs4FHpbnoy8BVJFvjB42Lj3"
	btcAddrMultiP2wshInP2sh = "39pvWsE9ctd6C9huSYHq4afLdhQjzQCJkD"
	btcAddrMultiP2wsh       = "bc1qru76nhzv58h4c347r9xtlks97gm3d0a7eej39rys2dfy3z9ckjmqa25hz2"
)

func TestBTC_PriKeyWif(t *testing.T) {
	priKey, _ := ParseWIF(btcPriWif)
	if priKey.WIF() != btcPriWif {
		t.Errorf("priKey.WIF() %s", priKey.WIF())
	}
}

func TestBTC_Compressed(t *testing.T) {

	priKey := NewPriKeyRandom()
	p1 := priKey.PubKey()
	p2, _ := NewPubKey(p1.Bytes())

	if hex.EncodeToString(p1.Bytes()) != hex.EncodeToString(p2.Bytes()) {
		t.Error("method priKey.PubKey() and NewPubKey() not matched ")
	}
}

func TestBTC_Uncompressed(t *testing.T) {

	priKey := NewPriKeyRandom()
	p1 := priKey.PubKeyUnCompressed()
	p2, _ := NewPubKey(p1.Bytes())

	if hex.EncodeToString(p1.Bytes()) != hex.EncodeToString(p2.Bytes()) {
		t.Error("method priKey.PubKeyUncompressed() and NewPubKey() not matched ")
	}
}

func TestBTC_P2pkh(t *testing.T) {

	priKey, _ := ParseWIF(btcPriWif)
	pub1 := priKey.PubKey()
	pub2 := priKey.PubKeyUnCompressed()
	if pub1.Address().P2pkh() != btcAddrP2pkh {
		t.Error("error compress address when P2pkh")
	}
	if pub2.Address().P2pkh() != btcAddrUnCompressP2pkh {
		t.Error("error uncompress address when P2pkh")
	}
}

func TestBTC_P2sh(t *testing.T) {

	priKey, _ := ParseWIF(btcPriWif)
	pub1 := priKey.PubKey()
	pub2 := priKey.PubKeyUnCompressed()
	if pub1.Address().P2sh() != btcAddrP2sh {
		t.Error("error compress address when P2sh")
	}
	if pub2.Address().P2sh() != btcAddrUnCompressP2sh {
		t.Error("error uncompress address when P2sh")
	}
}

func TestBTC_P2wpkh(t *testing.T) {

	priKey, _ := ParseWIF(btcPriWif)
	pub := priKey.PubKey()
	addr := pub.Address().P2wpkh()
	if addr != btcAddrP2wpkh {
		t.Error("error address P2wpkh")
	}
}

func TestBTC_P2wsh(t *testing.T) {

	priKey, _ := ParseWIF(btcPriWif)
	pub := priKey.PubKey()
	addr := pub.Address().P2wsh()
	if addr != btcAddrP2wsh {
		t.Error("error address P2wsh")
	}
}

func TestBTC_P2wpkhInP2sh(t *testing.T) {

	priKey, _ := ParseWIF(btcPriWif)
	pub := priKey.PubKey()
	if pub.Address().P2wpkhInP2sh() != btcAddrP2wpkhInP2sh {
		t.Error("error address P2wpkhInP2sh")
	}
}

func TestBTC_MultiP2wshInP2sh(t *testing.T) {
	priKey, _ := ParseWIF(btcPriWif)
	pub := priKey.PubKey()
	if pub.Address().MultiP2wshInP2sh() != btcAddrMultiP2wshInP2sh {
		t.Error("error address MultiP2wshInP2sh")
	}
}

func TestBTC_MultiP2wsh(t *testing.T) {
	priKey, _ := ParseWIF(btcPriWif)
	pub := priKey.PubKey()
	if pub.Address().MultiP2wsh() != btcAddrMultiP2wsh {
		t.Error("error address MultiP2wsh")
	}
}
