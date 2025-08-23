package auth

import (
	"encoding/hex"
	"crypto/rand"
)

func MakeRefreshToken() (string, error) {
	// Generate a random 256-bit (32-byte) hex-encoded string
	key := make([]byte, 32)


	/* 
		Don't need to check error because it will panic if it errors out:

		"Read fills b with cryptographically secure random bytes. It never
		returns an error, and always fills b entirely." - docs

	*/
	rand.Read(key)

	hexKey := hex.EncodeToString(key)

	return hexKey, nil

}