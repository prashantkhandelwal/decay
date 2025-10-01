package utils

import (
	"crypto/rand"
	"math/big"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// EncodeBase62 converts a big.Int to a Base62 string.
func EncodeBase62(n *big.Int) string {
	if n.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	base := big.NewInt(62)
	result := ""
	for n.Cmp(big.NewInt(0)) > 0 {
		mod := new(big.Int)
		n.DivMod(n, base, mod)
		result = string(base62Chars[mod.Int64()]) + result
	}
	return result
}

// GenerateBase62ID creates a random 128-bit Base62 ID.
func GenerateBase62ID() (string, error) {
	// Generate 128-bit random number
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	n := new(big.Int).SetBytes(b)
	return EncodeBase62(n), nil
}
