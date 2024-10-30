package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	return hex.EncodeToString(hash.Sum(nil))
}
