package utils

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"io"
)

func Hash160(data []byte) ([]byte, error) {
	h1, err := HashSha256(data)
	if err != nil {
		return nil, err
	}
	h2, err := HashRipeMD160(h1)
	if err != nil {
		return nil, err
	}
	return h2, nil
}

func HashRipeMD160(data []byte) ([]byte, error) {
	h := ripemd160.New()
	_, err := io.WriteString(h, string(data))
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func HashSha256(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func HashDoubleSha256(data []byte) ([]byte, error) {
	h1, err := HashSha256(data)
	if err != nil {
		return nil, err
	}
	h2, err := HashSha256(h1)
	if err != nil {
		return nil, err
	}
	return h2, nil
}
