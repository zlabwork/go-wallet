package wallet

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
	P2pkh() string
	P2sh() string
	P2wpkh() string
	P2wsh() string
	P2wpkhInP2sh() string
}

type BchAddress interface {
	P2pkh() string
	P2sh() string
}
