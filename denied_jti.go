package storage

import (
	// Standard Library Imports
	"crypto/sha256"
	"fmt"
	"time"
)

// DeniedJTI provides the structure for a Denied JSON Web Token (JWT) Token
// Identifier.
type DeniedJTI struct {
	JTI       string `bson:"-" json:"-" xml:"-"`
	Signature string `bson:"signature" json:"signature" xml:"signature"`
	Expiry    int64  `bson:"exp" json:"exp" xml:"exp"`
}

// NewDeniedJTI returns a new jti to be denied.
func NewDeniedJTI(jti string, exp time.Time) DeniedJTI {
	return DeniedJTI{
		JTI:       jti,
		Signature: SignatureFromJTI(jti),
		Expiry:    exp.Unix(),
	}
}

// SignatureFromJTI creates a JTI signature from the JWT Token ID.
func SignatureFromJTI(jti string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(jti)))
}
