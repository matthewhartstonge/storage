package storage

import (
	"crypto/sha512"
	"encoding/hex"
)

// SignatureHash hashes the signature to prevent errors where the signature is
// longer than 128 characters.
func SignatureHash(signature string) string {
	sum := sha512.Sum384([]byte(signature))
	return hex.EncodeToString(sum[:])
}
