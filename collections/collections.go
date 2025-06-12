// This package provides collections for data manipulation on the blockchain.
package collections

import (
	"errors"
	"fmt"

	"github.com/vlmoon99/near-sdk-go/borsh"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/types"
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
		return b.prefix + KeySeparator + types.Uint64ToString(k)
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

type Vector[T any] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	length     uint64
}

func NewVector[T any](prefix string) *Vector[T] {
	return &Vector[T]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		length:         0,
	}
}

func (v *Vector[T]) Length() uint64 {
	return v.length
}

func (v *Vector[T]) Push(value T) error {
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

	_, err = v.storage.Read(key)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vector[T]) Get(index uint64) (T, error) {
	var value T
	if index >= v.length {
		return value, errors.New(CollectionErrIndexOutOfBounds)
	}

	key := v.createKey(index)
	data, err := v.storage.Read(key)
	if err != nil {
		return value, err
	}

	err = v.serializer.Deserialize(data, &value)
	return value, err
}

func (v *Vector[T]) Set(index uint64, value T) error {
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

func (v *Vector[T]) Pop() (T, error) {
	var value T
	if v.length == 0 {
		return value, errors.New(CollectionErrVectorEmpty)
	}

	v.length--
	key := v.createKey(v.length)
	data, err := v.storage.Read(key)
	if err != nil {
		return value, err
	}

	err = v.serializer.Deserialize(data, &value)
	if err != nil {
		return value, err
	}

	_, err = v.storage.Delete(key)
	return value, err
}

func (v *Vector[T]) Clear() error {
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

func (v *Vector[T]) ToSlice() ([]T, error) {
	result := make([]T, v.length)
	for i := uint64(0); i < v.length; i++ {
		value, err := v.Get(i)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}
	return result, nil
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

type UnorderedMap[K comparable, V any] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	keys       Vector[[]byte]
}

func NewUnorderedMap[K comparable, V any](prefix string) *UnorderedMap[K, V] {
	return &UnorderedMap[K, V]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		keys:           *NewVector[[]byte](prefix + KeySeparator + "keys"),
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

	keyData, err := m.serializer.Serialize(key)
	if err != nil {
		return err
	}

	err = m.keys.Push(keyData)
	if err != nil {
		return err
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
		keyData, err := m.keys.Get(i)
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
				lastKey, err := m.keys.Get(length - 1)
				if err != nil {
					return err
				}
				err = m.keys.Set(i, lastKey)
				if err != nil {
					return err
				}
			}
			_, err := m.keys.Pop()
			return err
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
		keyData, err := m.keys.Get(i)
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

type LookupSet[T comparable] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
}

func NewLookupSet[T comparable](prefix string) *LookupSet[T] {
	return &LookupSet[T]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
	}
}

func (s *LookupSet[T]) Insert(value T) error {
	data, err := s.serializer.Serialize(value)
	if err != nil {
		return err
	}

	storageKey := s.createKey(value)
	_, err = s.storage.Write(storageKey, data)
	return err
}

func (s *LookupSet[T]) Remove(value T) error {
	storageKey := s.createKey(value)
	_, err := s.storage.Delete(storageKey)
	return err
}

func (s *LookupSet[T]) Contains(value T) (bool, error) {
	storageKey := s.createKey(value)
	return s.storage.HasKey(storageKey)
}

type UnorderedSet[T comparable] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	values     Vector[[]byte]
}

func NewUnorderedSet[T comparable](prefix string) *UnorderedSet[T] {
	return &UnorderedSet[T]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		values:         *NewVector[[]byte](prefix + KeySeparator + "values"),
	}
}

func (s *UnorderedSet[T]) Insert(value T) error {
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
	err = s.values.Push(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *UnorderedSet[T]) Remove(value T) error {
	storageKey := s.createKey(value)
	_, err := s.storage.Delete(storageKey)
	if err != nil {
		return err
	}

	length := s.values.Length()
	for i := uint64(0); i < length; i++ {
		valueData, err := s.values.Get(i)
		if err != nil {
			return err
		}

		var storedValue T
		err = s.serializer.Deserialize(valueData, &storedValue)
		if err != nil {
			return err
		}

		if storedValue == value {
			if i < length-1 {
				lastValue, err := s.values.Get(length - 1)
				if err != nil {
					return err
				}
				err = s.values.Set(i, lastValue)
				if err != nil {
					return err
				}
			}
			_, err := s.values.Pop()
			return err
		}
	}

	return nil
}

func (s *UnorderedSet[T]) Contains(value T) (bool, error) {
	storageKey := s.createKey(value)
	return s.storage.HasKey(storageKey)
}

func (s *UnorderedSet[T]) Values() ([]T, error) {
	length := s.values.Length()
	values := make([]T, length)

	for i := uint64(0); i < length; i++ {
		valueData, err := s.values.Get(i)
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

func (s *UnorderedSet[T]) Clear() error {
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

type TreeMap[K comparable, V any] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	keys       Vector[[]byte]
}

func NewTreeMap[K comparable, V any](prefix string) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		keys:           *NewVector[[]byte](prefix + KeySeparator + "keys"),
	}
}

func (m *TreeMap[K, V]) Insert(key K, value V) error {
	data, err := m.serializer.Serialize(value)
	if err != nil {
		return err
	}

	storageKey := m.createKey(key)
	fmt.Printf("storageKey: %v\n", storageKey)
	isSuccess, err := m.storage.Write(storageKey, data)
	fmt.Printf("isSuccess: %v\n", isSuccess)
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
		return err
	}

	keyData, err := m.serializer.Serialize(key)
	if err != nil {
		return err
	}

	keys, err := m.Keys()
	if err != nil {
		return err
	}

	insertIndex := 0
	for i, existingKey := range keys {
		if compareKeys(key, existingKey) <= 0 {
			insertIndex = i
			break
		}
		insertIndex = i + 1
	}

	err = m.keys.Clear()
	if err != nil {
		return err
	}

	for i := 0; i < insertIndex; i++ {
		existingKeyData, err := m.serializer.Serialize(keys[i])
		if err != nil {
			return err
		}
		err = m.keys.Push(existingKeyData)
		if err != nil {
			return err
		}
	}

	err = m.keys.Push(keyData)
	if err != nil {
		return err
	}

	for i := insertIndex; i < len(keys); i++ {
		existingKeyData, err := m.serializer.Serialize(keys[i])
		if err != nil {
			return err
		}
		err = m.keys.Push(existingKeyData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *TreeMap[K, V]) Get(key K) (V, error) {
	var value V
	storageKey := m.createKey(key)
	data, err := m.storage.Read(storageKey)
	if err != nil {
		return value, err
	}

	err = m.serializer.Deserialize(data, &value)
	return value, err
}

func (m *TreeMap[K, V]) Remove(key K) error {
	storageKey := m.createKey(key)
	_, err := m.storage.Delete(storageKey)
	if err != nil {
		return err
	}

	length := m.keys.Length()
	for i := uint64(0); i < length; i++ {
		keyData, err := m.keys.Get(i)
		if err != nil {
			return err
		}

		var storedKey K
		err = m.serializer.Deserialize(keyData, &storedKey)
		if err != nil {
			return err
		}

		if storedKey == key {
			for j := i; j < length-1; j++ {
				nextKey, err := m.keys.Get(j + 1)
				if err != nil {
					return err
				}
				err = m.keys.Set(j, nextKey)
				if err != nil {
					return err
				}
			}
			_, err := m.keys.Pop()
			return err
		}
	}

	return nil
}

func (m *TreeMap[K, V]) Contains(key K) (bool, error) {
	storageKey := m.createKey(key)
	return m.storage.HasKey(storageKey)
}

func (m *TreeMap[K, V]) Keys() ([]K, error) {
	length := m.keys.Length()
	keys := make([]K, length)

	for i := uint64(0); i < length; i++ {
		keyData, err := m.keys.Get(i)
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

func (m *TreeMap[K, V]) Values() ([]V, error) {
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

func (m *TreeMap[K, V]) Items() ([]struct {
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

func (m *TreeMap[K, V]) MinKey() (K, error) {
	if m.keys.Length() == 0 {
		var zero K
		return zero, errors.New(CollectionErrMapEmpty)
	}

	var keyData []byte
	keyData, err := m.keys.Get(0)
	if err != nil {
		var zero K
		return zero, err
	}

	var key K
	err = m.serializer.Deserialize(keyData, &key)
	if err != nil {
		var zero K
		return zero, err
	}

	return key, nil
}

func (m *TreeMap[K, V]) MaxKey() (K, error) {
	length := m.keys.Length()
	if length == 0 {
		var zero K
		return zero, errors.New(CollectionErrMapEmpty)
	}

	var keyData []byte
	keyData, err := m.keys.Get(length - 1)
	if err != nil {
		var zero K
		return zero, err
	}

	var key K
	err = m.serializer.Deserialize(keyData, &key)
	if err != nil {
		var zero K
		return zero, err
	}

	return key, nil
}

func (m *TreeMap[K, V]) FloorKey(key K) (K, error) {
	env.LogString("TreeMap FloorKey: Starting floor key search for: " + fmt.Sprint(key))
	keys, err := m.Keys()
	if err != nil {
		var zero K
		return zero, err
	}

	env.LogString("TreeMap FloorKey: Found " + fmt.Sprint(len(keys)) + " keys")
	for i := len(keys) - 1; i >= 0; i-- {
		env.LogString("TreeMap FloorKey: Checking key: " + fmt.Sprint(keys[i]))
		if compareKeys(keys[i], key) <= 0 {
			env.LogString("TreeMap FloorKey: Found floor key: " + fmt.Sprint(keys[i]))
			return keys[i], nil
		}
	}

	var zero K
	env.LogString("TreeMap FloorKey: No floor key found")
	return zero, errors.New(CollectionErrNoKeyLessOrEqual)
}

func (m *TreeMap[K, V]) CeilingKey(key K) (K, error) {
	env.LogString("TreeMap CeilingKey: Starting ceiling key search for: " + fmt.Sprint(key))
	keys, err := m.Keys()
	if err != nil {
		var zero K
		return zero, err
	}

	env.LogString("TreeMap CeilingKey: Found " + fmt.Sprint(len(keys)) + " keys")
	for _, k := range keys {
		env.LogString("TreeMap CeilingKey: Checking key: " + fmt.Sprint(k))
		if compareKeys(k, key) >= 0 {
			env.LogString("TreeMap CeilingKey: Found ceiling key: " + fmt.Sprint(k))
			return k, nil
		}
	}

	var zero K
	env.LogString("TreeMap CeilingKey: No ceiling key found")
	return zero, errors.New(CollectionErrNoKeyGreaterOrEqual)
}

func (m *TreeMap[K, V]) Clear() error {
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
