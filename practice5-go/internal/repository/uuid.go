package repository

import (
	"crypto/rand"
	"encoding/hex"
)

// NewUUIDv4 generates a random UUID v4 string using only the standard library.
func NewUUIDv4() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	// Set version (4) and variant bits.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	hexStr := hex.EncodeToString(b) // 32 chars
	return hexStr[0:8] + "-" + hexStr[8:12] + "-" + hexStr[12:16] + "-" + hexStr[16:20] + "-" + hexStr[20:32]
}
