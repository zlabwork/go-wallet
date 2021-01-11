package chain

import "crypto/rand"

func CreatePrivateKey() []byte {
	b := make([]byte, 32)
	rand.Read(b)
	return b
}
