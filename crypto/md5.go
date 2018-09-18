package crypto

import (
	"bytes"
	"crypto/md5"
)

// GetStringMD5 returns the md5 checksum of the given string as byte-array.
func GetStringMD5(payload string) []byte {
	return ComputeMD5checksum([]byte(payload))
}

// ComputeMD5checksum computes the md5 checksum for the given byte-array. Returns a byte-array.
func ComputeMD5checksum(payload []byte) []byte {
	computedMD5 := md5.Sum(payload)
	return computedMD5[0:]
}

// VerifyMD5 verifies that the md5 checksum of the data given in the payload parameter matches the one
// provided in the givenMD5 param.
func VerifyMD5(payload []byte, givenMD5 []byte) bool {
	return bytes.Equal(ComputeMD5checksum(payload), givenMD5)
}
