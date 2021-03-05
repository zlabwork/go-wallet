package bitcoincash

import (
    "crypto/sha256"
    "golang.org/x/crypto/ripemd160"
    "io"
)

func hash160(data []byte) ([]byte, error) {
    hash1, err := hashSha256(data)
    if err != nil {
        return nil, err
    }

    hash2, err := hashRipeMD160(hash1)
    if err != nil {
        return nil, err
    }

    return hash2, nil
}

func hashRipeMD160(data []byte) ([]byte, error) {
    h := ripemd160.New()
    _, err := io.WriteString(h, string(data))
    if err != nil {
        return nil, err
    }
    return h.Sum(nil), nil
}

func hashSha256(data []byte) ([]byte, error) {
    h := sha256.New()
    _, err := h.Write(data)
    if err != nil {
        return nil, err
    }
    return h.Sum(nil), nil
}
