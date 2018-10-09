package crypto

import (
	"encoding/hex"
	"loki/log"
	"testing"
)

func TestArgon2hashes(t *testing.T) {
	// []byte("Here is a string....")

	passwords := []string{"x", "m", "mm", "Matthias"}

	for _, password := range passwords {
		key := Argon2KDF([]byte(password))
		log.Debug("Key: " + hexdump(key) + ", password: " + password)
	}

	// t.Fail()

}

// Hexdump provides a string with the hex-representation of the byte-array given in data.
func hexdump(data []byte) string {
	return hex.EncodeToString(data)
}
