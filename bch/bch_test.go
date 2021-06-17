package bch

import (
    "encoding/hex"
    "strings"
    "testing"
)

const (
    bchPri    = "9F17A0F210D1030508B1A803D24A439B5784AE5F28897F162DAF380D691E52A2"
    bchPriWif = "L2YxxGfLPCqSXVPS4Gh9AGwkVRxjMoRPtCaSSXYQxkaa8taEHwMJ"
    bchAddr1  = "qrmfkdmj45n0y7qy8n767gsnwajz9pksmyme2sd8gv"
    bchAddr2  = "prmfkdmj45n0y7qy8n767gsnwajz9pksmyvuhl2yn3"
)

func TestBCHPriKey(t *testing.T) {

    bs, _ := hex.DecodeString(bchPri)

    priKey, _ := NewPriKey(bs)
    if hex.EncodeToString(priKey.Bytes()) != strings.ToLower(bchPri) {
        t.Error("private key is not matched")
    }

    // test address
    if priKey.PubKey().Address().P2PKH() != bchAddr1 {
        t.Error("P2PKH address error")
    }
    if priKey.PubKey().Address().P2SH() != bchAddr2 {
        t.Error("P2SH address error")
    }
}
