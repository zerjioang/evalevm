package uuid

import (
	"crypto/rand"
	"encoding/hex"
)

// UUIDv4 generates a random UUIDv4 string with no heap allocations.
func UUIDv4() string {
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		panic("failed to generate UUIDv4: " + err.Error())
	}

	// Set version (4) and variant (10)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variant 10

	var buf [36]byte // UUID string format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

	hex.Encode(buf[0:8], uuid[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:16])

	return string(buf[:])
}
