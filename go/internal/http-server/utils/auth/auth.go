package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

type Params struct {
	Time    uint32 // iterations
	Memory  uint32 // KiB
	Threads uint8
	KeyLen  uint32 // bytes
}

// Sensible defaults
var Default = Params{
	Time:    3,
	Memory:  64 * 1024, // 64 MiB
	Threads: 1,
	KeyLen:  32,
}

func GenerateSalt(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

// Hash(password, saltB64) -> hashB64
func Hash(password, saltB64 string, p Params) (string, error) {
	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLen)
	return base64.RawStdEncoding.EncodeToString(key), nil
}

// Verify compares provided password against stored salt+hash
func Verify(password, saltB64, expectedHashB64 string, p Params) (bool, error) {
	computed, err := Hash(password, saltB64, p)
	if err != nil {
		return false, err
	}
	// constant-time compare
	if subtle.ConstantTimeCompare([]byte(computed), []byte(expectedHashB64)) == 1 {
		return true, nil
	}
	return false, nil
}

// Convenience for signup: returns (saltB64, hashB64)
func NewSecret(password string, p Params) (string, string, error) {
	salt, err := GenerateSalt(16) // 128-bit
	if err != nil {
		return "", "", err
	}
	hash, err := Hash(password, salt, p)
	if err != nil {
		return "", "", err
	}
	if hash == "" {
		return "", "", errors.New("empty hash")
	}
	return salt, hash, nil
}
