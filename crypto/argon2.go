package crypto

import (
	"golang.org/x/crypto/argon2"
)

// Argon2KDF produces a fixed length (32 bytes) key from the given password using the Argon2 Key derivation algorythm
func Argon2KDF(password []byte) []byte {
	salt := []byte{0x4F, 0xEB, 0x43, 0xDB, 0xBE, 0xB0, 0x43, 0x5C, 0x86, 0xC9, 0x7F, 0xA8, 0x9B, 0x4B, 0xDB, 0x0C}
	return argon2.Key(password, salt, 3, 32*1024, 4, 32)
}
