package types

const (
	// ModuleName defines the module name
	ModuleName = "twin"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_twin"

	// TrainingStateKey is the key for TrainingState
	TrainingStateKey = "TrainingState/value/"

	// TwinKeyPrefix is the prefix to retrieve all Twin
	TwinKeyPrefix = "Twin/value/"
)

// Convert string key in bytes
func KeyPrefix(p string) []byte {
	return []byte(p)
}

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
