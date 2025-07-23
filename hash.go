package storage

import (
	"crypto/sha512"
	"fmt"
)

// SignatureHash hashes the signature to prevent errors where the signature is
// longer than 128 characters.
func SignatureHash(signature string) string {
	if signature == "" {
		return ""
	}
	return fmt.Sprintf("%x", sha512.Sum384([]byte(signature)))
}
