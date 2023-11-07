package hash

import (
	"crypto/sha256"
	"fmt"
)

// PasswordHasher provides hashing logic for securely passwords storing
type PasswordHasher interface {
	Hash(password string) (string, error)
}

// SHA256Hasher uses SHA256 to hash passwords with provided salt.
type SHA256Hasher struct {
	salt string
}

func NewSHA256Hasher(salt string) *SHA256Hasher {
	return &SHA256Hasher{salt: salt}
}

// Hash creates SHA256 hash of given password with salt
func (h *SHA256Hasher) Hash(password string) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
