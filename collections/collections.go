// This package provides collections for data manipulation on the blockchain.
package collections

import (
	"errors"
	"fmt"

	"github.com/vlmoon99/near-sdk-go/borsh"
	"github.com/vlmoon99/near-sdk-go/env"
)

const (
	CollectionErrIndexOutOfBounds    = "index out of bounds"
	CollectionErrVectorEmpty         = "vector is empty"
	CollectionErrMapEmpty            = "map is empty"
	CollectionErrNoKeyLessOrEqual    = "no key less than or equal to given key"
	CollectionErrNoKeyGreaterOrEqual = "no key greater than or equal to given key"
	CollectionErrUnsupportedKeyType  = "unsupported key type"
)

const (
	KeySeparator = ":"
)

type Collection interface {
	GetPrefix() string
}

type BaseCollection struct {
	prefix string
}

func NewBaseCollection(prefix string) BaseCollection {
	return BaseCollection{
		prefix: prefix,
	}
}

func (b BaseCollection) GetPrefix() string {
	return b.prefix
}

func (b BaseCollection) createKey(key interface{}) string {
	switch k := key.(type) {
	case string:
		return b.prefix + KeySeparator + k
	case uint64:
		return b.prefix + KeySeparator + fmt.Sprintf("%d", k)
	case int:
		return b.createKey(uint64(k))
	default:
		env.PanicStr(CollectionErrUnsupportedKeyType)
		return ""
	}
}

type Storage interface {
	Read(key string) ([]byte, error)
	Write(key string, value []byte) (bool, error)
	Delete(key string) (bool, error)
	HasKey(key string) (bool, error)
}

type DefaultStorage struct{}

func (s DefaultStorage) Read(key string) ([]byte, error) {
	return env.StorageRead([]byte(key))
}

func (s DefaultStorage) Write(key string, value []byte) (bool, error) {
	return env.StorageWrite([]byte(key), value)
}

func (s DefaultStorage) Delete(key string) (bool, error) {
	return env.StorageRemove([]byte(key))
}

func (s DefaultStorage) HasKey(key string) (bool, error) {
	return env.StorageHasKey([]byte(key))
}

type Serializer interface {
	Serialize(value interface{}) ([]byte, error)
	Deserialize(data []byte, value interface{}) error
}

type DefaultSerializer struct{}

func (s DefaultSerializer) Serialize(value interface{}) ([]byte, error) {
	return borsh.Serialize(value)
}

func (s DefaultSerializer) Deserialize(data []byte, value interface{}) error {
	return borsh.Deserialize(data, value)
}

type Vector struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	length     uint64
}

func NewVector(prefix string) *Vector {
	return &Vector{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		length:         0,
	}
}

func (v *Vector) Length() uint64 {
	return v.length
}

func (v *Vector) Push(value interface{}) error {
	data, err := v.serializer.Serialize(value)
	if err != nil {
		return err
	}

	key := v.createKey(v.length)
	_, err = v.storage.Write(key, data)
	if err != nil {
		return err
	}

	v.length++

	// Verify the write
	_, err = v.storage.Read(key)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vector) Get(index uint64, value interface{}) error {
	if index >= v.length {
		return errors.New(CollectionErrIndexOutOfBounds)
	}

	key := v.createKey(index)
	data, err := v.storage.Read(key)
	if err != nil {
		return err
	}

	return v.serializer.Deserialize(data, value)
}

func (v *Vector) Set(index uint64, value interface{}) error {
	if index >= v.length {
		return errors.New(CollectionErrIndexOutOfBounds)
	}

	data, err := v.serializer.Serialize(value)
	if err != nil {
		return err
	}

	key := v.createKey(index)
	_, err = v.storage.Write(key, data)
	return err
}

func (v *Vector) Pop(value interface{}) error {
	if v.length == 0 {
		return errors.New(CollectionErrVectorEmpty)
	}

	v.length--
	key := v.createKey(v.length)
	data, err := v.storage.Read(key)
	if err != nil {
		return err
	}

	err = v.serializer.Deserialize(data, value)
	if err != nil {
		return err
	}

	_, err = v.storage.Delete(key)
	return err
}

func (v *Vector) Clear() error {
	for i := uint64(0); i < v.length; i++ {
		key := v.createKey(i)
		_, err := v.storage.Delete(key)
		if err != nil {
			return err
		}
	}
	v.length = 0
	return nil
}

func (v *Vector) ToSlice(valueType interface{}) ([]interface{}, error) {
	result := make([]interface{}, v.length)
	for i := uint64(0); i < v.length; i++ {
		err := v.Get(i, &result[i])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type UnorderedMap[K comparable, V any] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	keys       Vector
}

func NewUnorderedMap[K comparable, V any](prefix string) *UnorderedMap[K, V] {
	return &UnorderedMap[K, V]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		keys:           *NewVector(prefix + KeySeparator + "keys"),
	}
}

func (m *UnorderedMap[K, V]) Insert(key K, value V) error {
	data, err := m.serializer.Serialize(value)
	if err != nil {
		return err
	}

	storageKey := m.createKey(key)
	_, err = m.storage.Write(storageKey, data)
	if err != nil {
		return err
	}

	contains, err := m.Contains(key)
	if err != nil {
		return err
	}
	if !contains {
		// Serialize the key for storage
		keyData, err := m.serializer.Serialize(key)
		if err != nil {
			return err
		}

		// Store the key in the Vector
		err = m.keys.Push(keyData)
		if err != nil {
			return err
		}

		// Verify the key was stored
		var storedKeyData []byte
		err = m.keys.Get(m.keys.Length()-1, &storedKeyData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *UnorderedMap[K, V]) Get(key K) (V, error) {
	var value V
	storageKey := m.createKey(key)
	data, err := m.storage.Read(storageKey)
	if err != nil {
		return value, err
	}
	err = m.serializer.Deserialize(data, &value)
	return value, err
}

func (m *UnorderedMap[K, V]) Remove(key K) error {
	storageKey := m.createKey(key)
	_, err := m.storage.Delete(storageKey)
	if err != nil {
		return err
	}

	length := m.keys.Length()
	for i := uint64(0); i < length; i++ {
		var keyData []byte
		err := m.keys.Get(i, &keyData)
		if err != nil {
			return err
		}

		var storedKey K
		err = m.serializer.Deserialize(keyData, &storedKey)
		if err != nil {
			return err
		}

		if storedKey == key {
			if i < length-1 {
				var lastKey []byte
				err := m.keys.Get(length-1, &lastKey)
				if err != nil {
					return err
				}
				err = m.keys.Set(i, lastKey)
				if err != nil {
					return err
				}
			}
			var unused interface{}
			return m.keys.Pop(&unused)
		}
	}

	return nil
}

func (m *UnorderedMap[K, V]) Contains(key K) (bool, error) {
	storageKey := m.createKey(key)
	return m.storage.HasKey(storageKey)
}

func (m *UnorderedMap[K, V]) Keys() ([]K, error) {
	length := m.keys.Length()
	keys := make([]K, length)

	for i := uint64(0); i < length; i++ {
		var keyData []byte
		err := m.keys.Get(i, &keyData)
		if err != nil {
			return nil, err
		}

		err = m.serializer.Deserialize(keyData, &keys[i])
		if err != nil {
			return nil, err
		}
	}

	return keys, nil
}

func (m *UnorderedMap[K, V]) Values() ([]V, error) {
	keys, err := m.Keys()
	if err != nil {
		return nil, err
	}

	values := make([]V, len(keys))
	for i, key := range keys {
		value, err := m.Get(key)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	return values, nil
}

func (m *UnorderedMap[K, V]) Items() ([]struct {
	Key   K
	Value V
}, error) {
	keys, err := m.Keys()
	if err != nil {
		return nil, err
	}

	items := make([]struct {
		Key   K
		Value V
	}, len(keys))
	for i, key := range keys {
		value, err := m.Get(key)
		if err != nil {
			return nil, err
		}
		items[i].Key = key
		items[i].Value = value
	}

	return items, nil
}

func (m *UnorderedMap[K, V]) Clear() error {
	keys, err := m.Keys()
	if err != nil {
		return err
	}

	for _, key := range keys {
		err := m.Remove(key)
		if err != nil {
			return err
		}
	}

	return m.keys.Clear()
}

type LookupSet struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
}

func NewLookupSet(prefix string) *LookupSet {
	return &LookupSet{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
	}
}

func (s *LookupSet) Insert(value interface{}) error {
	data, err := s.serializer.Serialize(value)
	if err != nil {
		return err
	}

	storageKey := s.createKey(value)
	_, err = s.storage.Write(storageKey, data)
	return err
}

func (s *LookupSet) Remove(value interface{}) error {
	storageKey := s.createKey(value)
	_, err := s.storage.Delete(storageKey)
	return err
}

func (s *LookupSet) Contains(value interface{}) (bool, error) {
	storageKey := s.createKey(value)
	return s.storage.HasKey(storageKey)
}

type UnorderedSet struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	values     Vector
}

func NewUnorderedSet(prefix string) *UnorderedSet {
	return &UnorderedSet{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		values:         *NewVector(prefix + KeySeparator + "values"),
	}
}

func (s *UnorderedSet) Insert(value interface{}) error {
	contains, err := s.Contains(value)
	if err != nil {
		return err
	}
	if contains {
		return nil
	}

	data, err := s.serializer.Serialize(value)
	if err != nil {
		return err
	}

	storageKey := s.createKey(value)
	_, err = s.storage.Write(storageKey, data)
	if err != nil {
		return err
	}

	return s.values.Push(data)
}

func (s *UnorderedSet) Remove(value interface{}) error {
	storageKey := s.createKey(value)
	_, err := s.storage.Delete(storageKey)
	if err != nil {
		return err
	}

	length := s.values.Length()
	for i := uint64(0); i < length; i++ {
		var valueData []byte
		err := s.values.Get(i, &valueData)
		if err != nil {
			return err
		}

		var storedValue interface{}
		err = s.serializer.Deserialize(valueData, &storedValue)
		if err != nil {
			return err
		}

		if storedValue == value {
			if i < length-1 {
				var lastValue []byte
				err := s.values.Get(length-1, &lastValue)
				if err != nil {
					return err
				}
				err = s.values.Set(i, lastValue)
				if err != nil {
					return err
				}
			}
			var unused interface{}
			return s.values.Pop(&unused)
		}
	}

	return nil
}

func (s *UnorderedSet) Contains(value interface{}) (bool, error) {
	storageKey := s.createKey(value)
	return s.storage.HasKey(storageKey)
}

func (s *UnorderedSet) Values() ([]interface{}, error) {
	length := s.values.Length()
	values := make([]interface{}, length)

	for i := uint64(0); i < length; i++ {
		var valueData []byte
		err := s.values.Get(i, &valueData)
		if err != nil {
			return nil, err
		}

		err = s.serializer.Deserialize(valueData, &values[i])
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func (s *UnorderedSet) Clear() error {
	values, err := s.Values()
	if err != nil {
		return err
	}

	for _, value := range values {
		err := s.Remove(value)
		if err != nil {
			return err
		}
	}

	return s.values.Clear()
}

type TreeMap struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	keys       Vector
}

func NewTreeMap(prefix string) *TreeMap {
	return &TreeMap{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		keys:           *NewVector(prefix + KeySeparator + "keys"),
	}
}

func (m *TreeMap) Insert(key interface{}, value interface{}) error {
	data, err := m.serializer.Serialize(value)
	if err != nil {
		return err
	}

	storageKey := m.createKey(key)
	_, err = m.storage.Write(storageKey, data)
	if err != nil {
		return err
	}

	contains, err := m.Contains(key)
	if err != nil {
		return err
	}
	if !contains {
		keyData, err := m.serializer.Serialize(key)
		if err != nil {
			return err
		}

		length := m.keys.Length()
		insertIndex := uint64(0)

		for i := uint64(0); i < length; i++ {
			var existingKeyData []byte
			err := m.keys.Get(i, &existingKeyData)
			if err != nil {
				return err
			}

			var existingKey interface{}
			err = m.serializer.Deserialize(existingKeyData, &existingKey)
			if err != nil {
				return err
			}

			if compareKeys(key, existingKey) < 0 {
				break
			}
			insertIndex = i + 1
		}

		if insertIndex == length {
			return m.keys.Push(keyData)
		}

		for i := length; i > insertIndex; i-- {
			var prevKey []byte
			err := m.keys.Get(i-1, &prevKey)
			if err != nil {
				return err
			}
			err = m.keys.Set(i, prevKey)
			if err != nil {
				return err
			}
		}

		return m.keys.Set(insertIndex, keyData)
	}

	return nil
}

func (m *TreeMap) Get(key interface{}, value interface{}) error {
	storageKey := m.createKey(key)
	data, err := m.storage.Read(storageKey)
	if err != nil {
		return err
	}

	return m.serializer.Deserialize(data, value)
}

func (m *TreeMap) Remove(key interface{}) error {
	storageKey := m.createKey(key)
	_, err := m.storage.Delete(storageKey)
	if err != nil {
		return err
	}

	length := m.keys.Length()
	for i := uint64(0); i < length; i++ {
		var keyData []byte
		err := m.keys.Get(i, &keyData)
		if err != nil {
			return err
		}

		var storedKey interface{}
		err = m.serializer.Deserialize(keyData, &storedKey)
		if err != nil {
			return err
		}

		if storedKey == key {
			for j := i; j < length-1; j++ {
				var nextKey []byte
				err := m.keys.Get(j+1, &nextKey)
				if err != nil {
					return err
				}
				err = m.keys.Set(j, nextKey)
				if err != nil {
					return err
				}
			}
			var unused interface{}
			return m.keys.Pop(&unused)
		}
	}

	return nil
}

func (m *TreeMap) Contains(key interface{}) (bool, error) {
	storageKey := m.createKey(key)
	return m.storage.HasKey(storageKey)
}

func (m *TreeMap) Keys() ([]interface{}, error) {
	length := m.keys.Length()
	keys := make([]interface{}, length)

	for i := uint64(0); i < length; i++ {
		var keyData []byte
		err := m.keys.Get(i, &keyData)
		if err != nil {
			return nil, err
		}

		err = m.serializer.Deserialize(keyData, &keys[i])
		if err != nil {
			return nil, err
		}
	}

	return keys, nil
}

func (m *TreeMap) Values(valueType interface{}) ([]interface{}, error) {
	keys, err := m.Keys()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(keys))
	for i, key := range keys {
		err := m.Get(key, &values[i])
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func (m *TreeMap) Items(valueType interface{}) ([]struct{ Key, Value interface{} }, error) {
	keys, err := m.Keys()
	if err != nil {
		return nil, err
	}

	items := make([]struct{ Key, Value interface{} }, len(keys))
	for i, key := range keys {
		items[i].Key = key
		err := m.Get(key, &items[i].Value)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

func (m *TreeMap) MinKey() (interface{}, error) {
	if m.keys.Length() == 0 {
		return nil, errors.New(CollectionErrMapEmpty)
	}

	var keyData []byte
	err := m.keys.Get(0, &keyData)
	if err != nil {
		return nil, err
	}

	var key interface{}
	err = m.serializer.Deserialize(keyData, &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (m *TreeMap) MaxKey() (interface{}, error) {
	length := m.keys.Length()
	if length == 0 {
		return nil, errors.New(CollectionErrMapEmpty)
	}

	var keyData []byte
	err := m.keys.Get(length-1, &keyData)
	if err != nil {
		return nil, err
	}

	var key interface{}
	err = m.serializer.Deserialize(keyData, &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (m *TreeMap) FloorKey(key interface{}) (interface{}, error) {
	keys, err := m.Keys()
	if err != nil {
		return nil, err
	}

	for i := len(keys) - 1; i >= 0; i-- {
		if compareKeys(keys[i], key) <= 0 {
			return keys[i], nil
		}
	}

	return nil, errors.New(CollectionErrNoKeyLessOrEqual)
}

func (m *TreeMap) CeilingKey(key interface{}) (interface{}, error) {
	keys, err := m.Keys()
	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		if compareKeys(k, key) >= 0 {
			return k, nil
		}
	}

	return nil, errors.New(CollectionErrNoKeyGreaterOrEqual)
}

func (m *TreeMap) Clear() error {
	keys, err := m.Keys()
	if err != nil {
		return err
	}

	for _, key := range keys {
		err := m.Remove(key)
		if err != nil {
			return err
		}
	}

	return m.keys.Clear()
}

func compareKeys(a, b interface{}) int {
	switch a := a.(type) {
	case string:
		if b, ok := b.(string); ok {
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		}
	case int:
		if b, ok := b.(int); ok {
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		}
	case uint64:
		if b, ok := b.(uint64); ok {
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		}
	}
	env.PanicStr(CollectionErrUnsupportedKeyType)
	return 0
}

type LookupMap[K comparable, V any] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
}

func NewLookupMap[K comparable, V any](prefix string) *LookupMap[K, V] {
	return &LookupMap[K, V]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
	}
}

func (m *LookupMap[K, V]) Insert(key K, value V) error {
	data, err := m.serializer.Serialize(value)
	if err != nil {
		return err
	}

	storageKey := m.createKey(key)
	_, err = m.storage.Write(storageKey, data)
	return err
}

func (m *LookupMap[K, V]) Get(key K) (V, error) {
	var value V
	storageKey := m.createKey(key)
	data, err := m.storage.Read(storageKey)
	if err != nil {
		return value, err
	}
	err = m.serializer.Deserialize(data, &value)
	return value, err
}

func (m *LookupMap[K, V]) Remove(key K) error {
	storageKey := m.createKey(key)
	_, err := m.storage.Delete(storageKey)
	return err
}

func (m *LookupMap[K, V]) Contains(key K) (bool, error) {
	storageKey := m.createKey(key)
	return m.storage.HasKey(storageKey)
}
