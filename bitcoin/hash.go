package bitcoin

import (
    "crypto/sha256"
    "golang.org/x/crypto/ripemd160"
    "io"
)

func checksum(data []byte) ([]byte, error) {
    hash, err := hashDoubleSha256(data)
    if err != nil {
        return nil, err
    }

    return hash[:4], nil
}

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
    hasher := ripemd160.New()
    _, err := io.WriteString(hasher, string(data))
    if err != nil {
        return nil, err
    }
    return hasher.Sum(nil), nil
}

func hashSha256(data []byte) ([]byte, error) {
    hasher := sha256.New()
    _, err := hasher.Write(data)
    if err != nil {
        return nil, err
    }
    return hasher.Sum(nil), nil
}

func hashDoubleSha256(data []byte) ([]byte, error) {
    hash1, err := hashSha256(data)
    if err != nil {
        return nil, err
    }

    hash2, err := hashSha256(hash1)
    if err != nil {
        return nil, err
    }
    return hash2, nil
}
