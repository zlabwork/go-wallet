package btc

import "fmt"

type PubKey struct {
	k          []byte
	compressed bool
}

func NewPubKey(b []byte) (*PubKey, error) {
	l := len(b)
	if l == 33 || l == 65 {
		c := true
		if l != 33 {
			c = false
		}
		return &PubKey{k: b, compressed: c}, nil
	}
	return nil, fmt.Errorf("invalid length")
}

func (pub *PubKey) Bytes() []byte {
	return pub.k
}

func (pub *PubKey) Address() *Address {
	return NewAddress(pub.k)
}
