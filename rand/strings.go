package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// Bytes generates n random bytes or returns error
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	// Read return number of bytes - we're ignoring
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// NBytes returns the number of bytes in the base64 URL encoded string
// used for testing PasswordHash
func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

// String generates a byte slice of size n and return
// a base64 url encoded version of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper function used to generate
// remember tokens of a predetermined byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
