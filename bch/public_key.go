package bch

import "fmt"

type PubKey struct {
    k []byte
}

func NewPubKey(b []byte) (*PubKey, error) {
    if len(b) == 33 || len(b) == 65 {
        return &PubKey{k: b}, nil
    }
    return nil, fmt.Errorf("invalid length")
}

func (pub *PubKey) Bytes() []byte {
    return pub.k
}

func (pub *PubKey) Address() *Address {
    return NewAddress(pub.k)
}
