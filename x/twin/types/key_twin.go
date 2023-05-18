package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// TwinKeyPrefix is the prefix to retrieve all Twin
	TwinKeyPrefix = "Twin/value/"
)

// TwinKey returns the store key to retrieve a Twin from the index fields
func TwinKey(
	name string,
) []byte {
	var key []byte

	nameBytes := []byte(name)
	key = append(key, nameBytes...)
	key = append(key, []byte("/")...)

	return key
}
