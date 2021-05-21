package bitcoin

import "errors"

type pubKeyData struct {
    key []byte
}

// https://learnmeabitcoin.com/technical/public-key
func NewPubKey(pubKey []byte) (*pubKeyData, error) {
    if len(pubKey) == 33 || len(pubKey) == 65 {
        return &pubKeyData{
            key: pubKey,
        }, nil
    }
    return nil, errors.New("invalid length")
}

func (pub *pubKeyData) Key() []byte {
    return pub.key
}

func (pub *pubKeyData) Address() *addrData {
    h, _ := hash160(pub.key)
    return &addrData{
        hash:   h,
        pubKey: pub.key,
    }
}

func (pub *pubKeyData) isCompressed() bool {
    return len(pub.key) == 33
}
