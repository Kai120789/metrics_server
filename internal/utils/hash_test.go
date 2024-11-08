package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestGenerateHash(t *testing.T) {
	// test values
	testKey := "example_key"

	// get hash with SHA-256
	hash := sha256.New()
	hash.Write([]byte(testKey))
	expectedHash := hex.EncodeToString(hash.Sum(nil))

	// gen hash by GenerateHash
	generatedHash := GenerateHash(testKey)

	// cimpare generated and expected hashes
	if generatedHash != expectedHash {
		t.Errorf("Ожидаемый хеш %s, но получили %s", expectedHash, generatedHash)
	}
}
