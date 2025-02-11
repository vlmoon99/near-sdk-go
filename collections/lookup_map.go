package collections

import (
	"errors"

	"github.com/vlmoon99/near-sdk-go/borsh"
	"github.com/vlmoon99/near-sdk-go/env"
)

type LookupMap struct {
	KeyPrefix []byte
}

func NewLookupMap(keyPrefix []byte) *LookupMap {
	return &LookupMap{KeyPrefix: keyPrefix}
}

func (m *LookupMap) rawKeyToStorageKey(rawKey []byte) []byte {
	combined := make([]byte, len(m.KeyPrefix)+len(rawKey))

	copy(combined, m.KeyPrefix)
	copy(combined[len(m.KeyPrefix):], rawKey)

	return combined
}

func (m *LookupMap) ContainsKey(key []byte) (bool, error) {
	storageKey := m.rawKeyToStorageKey(key)

	return env.StorageHasKey(storageKey)
}

func (m *LookupMap) Get(key []byte) (interface{}, error) {

	storageKey := m.rawKeyToStorageKey(key)
	valueBytes, err := env.StorageRead(storageKey)

	if err != nil {
		return nil, errors.New("storage read error")
	}

	if valueBytes == nil {
		return nil, nil
	}

	var value string

	err = borsh.Deserialize(valueBytes, &value)

	if err != nil {
		return nil, errors.New("deserialization error")
	}

	return value, nil
}

func (m *LookupMap) Insert(key []byte, value interface{}) error {

	storageKey := m.rawKeyToStorageKey(key)
	valueBytes, err := borsh.Serialize(value)

	if err != nil {
		return errors.New("serialization error")
	}
	_, err = env.StorageWrite(storageKey, valueBytes)

	return err
}

func (m *LookupMap) Remove(key []byte) error {

	storageKey := m.rawKeyToStorageKey(key)
	_, err := env.StorageRemove(storageKey)

	if err != nil {
		return err
	}

	return nil
}

func (m *LookupMap) Serialize() ([]byte, error) {
	return borsh.Serialize(m)
}

func DeserializeLookupMap(data []byte) (*LookupMap, error) {
	var lm LookupMap
	err := borsh.Deserialize(data, &lm)
	if err != nil {
		return nil, err
	}
	return &lm, nil
}
