package common

type PriKey interface {
    Bytes() []byte
    Base64() string
    Hex() string
    // WIF() string
}

type PubKey interface {
    Bytes() []byte
}

type BtcAddress interface {
    P2PKH() string
    P2SH() string
    P2WPKH() string
    P2WSH() string
    P2SHP2WPKH() string
}

type BchAddress interface {
    P2PKH() string
    P2SH() string
}
