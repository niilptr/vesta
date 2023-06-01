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
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TrainingStateKey = "TrainingState/value/"
)

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
