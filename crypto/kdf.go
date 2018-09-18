package crypto

// KeyDerivator defines a type which turns a password given as byte-array into a fixed-sized key returned as byte-array as well.
type KeyDerivator func(password []byte) []byte

// NewKeyDerivator produces the systems default KeyDerivator: Argon2KDF
func NewKeyDerivator() KeyDerivator {
	return Argon2KDF
}
