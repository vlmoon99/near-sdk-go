package collections

import (
	"github.com/vlmoon99/near-sdk-go/borsh"
	"github.com/vlmoon99/near-sdk-go/env"
)

// The LookupMap type represents a map for storing and retrieving key-value pairs.
type LookupMap struct {
	keyPrefix []byte
}

// Creates and returns a new LookupMap instance.
//
// Parameters:
//
//	keyPrefix: The prefix to be used for the keys in the map.
func NewLookupMap(keyPrefix []byte) *LookupMap {
	return &LookupMap{keyPrefix: keyPrefix}
}

// Combines the key prefix and the raw key to create a storage key.
//
// Parameters:
//
//	rawKey: The raw key to be combined with the key prefix.
func (m *LookupMap) rawKeyToStorageKey(rawKey []byte) []byte {
	combined := make([]byte, len(m.keyPrefix)+len(rawKey))

	copy(combined, m.keyPrefix)
	copy(combined[len(m.keyPrefix):], rawKey)

	return combined
}

// Checks if a key exists in the map.
//
// Parameters:
//
//	key: The key to check for existence.
//
// Returns:
//
//	bool: True if the key exists, false otherwise.
//	error: An error if the key checking fails.
func (m *LookupMap) ContainsKey(key []byte) (bool, error) {
	storageKey := m.rawKeyToStorageKey(key)

	return env.StorageHasKey(storageKey)
}

// Retrieves the value associated with the specified key.
//
// Parameters:
//
//	key: The key to retrieve the value for.
//
// Returns:
//
//	interface{}: The value associated with the key, or nil if not found.
//	error: An error if the key retrieval or deserialization fails.
func (m *LookupMap) Get(key []byte) (interface{}, error) {
	storageKey := m.rawKeyToStorageKey(key)
	valueBytes, err := env.StorageRead(storageKey)

	if err != nil {
		return nil, err
	}

	if valueBytes == nil {
		return nil, nil
	}

	var value string

	err = borsh.Deserialize(valueBytes, &value)

	if err != nil {
		return nil, err
	}

	return value, nil
}

// Inserts a key-value pair into the map.
//
// Parameters:
//
//	key: The key to insert.
//	value: The value to insert.
//
// Returns:
//
//	error: An error if the key insertion or serialization fails.
func (m *LookupMap) Insert(key []byte, value interface{}) error {
	storageKey := m.rawKeyToStorageKey(key)
	valueBytes, err := borsh.Serialize(value)
	if err != nil {
		return err
	}
	_, err = env.StorageWrite(storageKey, valueBytes)

	return err
}

// Removes a key-value pair from the map.
//
// Parameters:
//
//	key: The key to remove.
//
// Returns:
//
//	error: An error if the key removal fails.
func (m *LookupMap) Remove(key []byte) error {
	storageKey := m.rawKeyToStorageKey(key)
	_, err := env.StorageRemove(storageKey)

	if err != nil {
		return err
	}

	return nil
}

// Serializes the LookupMap instance to a byte slice.
//
// Returns:
//
//	[]byte: The serialized byte slice.
//	error: An error if the serialization fails.
func (m *LookupMap) Serialize() ([]byte, error) {
	return m.keyPrefix, nil
}

// Deserializes a byte slice into a LookupMap instance.
//
// Parameters:
//
//	data: The byte slice to deserialize.
//
// Returns:
//
//	*LookupMap: The deserialized LookupMap instance.
//	error: An error if the deserialization fails.
func DeserializeLookupMap(data []byte) (*LookupMap, error) {
	return &LookupMap{keyPrefix: data}, nil
}
