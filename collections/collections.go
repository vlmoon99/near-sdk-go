// This package provides collections for data manipulation on the blockchain.
package collections

import (
	"errors"

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
	CollectionErrKeyNotFound         = "key not found"
	CollectionErrValueNotFound       = "value not found in set"
	CollectionErrInconsistentState   = "inconsistent contract state: value index not found"
)

const (
	KeySeparator  = ":"
	KeysPrefix    = "keys"
	ValuesPrefix  = "values"
	IndicesPrefix = "indices"
)

const (
	LogNoFloorKey   = "TreeMap FloorKey: No floor key found"
	LogNoCeilingKey = "TreeMap CeilingKey: No ceiling key found"
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
	indices    LookupMap[K, uint64]
}

func NewUnorderedMap[K comparable, V any](prefix string) *UnorderedMap[K, V] {
	return &UnorderedMap[K, V]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		keys:           *NewVector[[]byte](prefix + KeySeparator + KeysPrefix),
		indices:        *NewLookupMap[K, uint64](prefix + KeySeparator + IndicesPrefix),
	}
}

func (m *UnorderedMap[K, V]) Insert(key K, value V) error {
	storageKey := m.createKey(key)
	exists, err := m.storage.HasKey(storageKey)
	if err != nil {
		return err
	}

	data, err := m.serializer.Serialize(value)
	if err != nil {
		return err
	}
	_, err = m.storage.Write(storageKey, data)
	if err != nil {
		return err
	}

	if !exists {
		index := m.keys.Length()
		keyData, err := m.serializer.Serialize(key)
		if err != nil {
			return err
		}
		err = m.keys.Push(keyData)
		if err != nil {
			return err
		}
		err = m.indices.Insert(key, index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *UnorderedMap[K, V]) Remove(key K) error {
	storageKey := m.createKey(key)
	exists, err := m.storage.HasKey(storageKey)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(CollectionErrKeyNotFound)
	}

	index, err := m.indices.Get(key)
	if err != nil {
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
					var lastKeyObj K
					err = m.serializer.Deserialize(lastKey, &lastKeyObj)
					if err != nil {
						return err
					}
					err = m.indices.Insert(lastKeyObj, i)
					if err != nil {
						return err
					}
				}
				_, err = m.keys.Pop()
				if err != nil {
					return err
				}
				break
			}
		}
	} else {
		length := m.keys.Length()
		if index < length {
			if index != length-1 {
				lastKey, err := m.keys.Get(length - 1)
				if err != nil {
					return err
				}
				err = m.keys.Set(index, lastKey)
				if err != nil {
					return err
				}
				var lastKeyObj K
				err = m.serializer.Deserialize(lastKey, &lastKeyObj)
				if err != nil {
					return err
				}
				err = m.indices.Insert(lastKeyObj, index)
				if err != nil {
					return err
				}
			}
			_, err = m.keys.Pop()
			if err != nil {
				return err
			}
		}
	}

	_, err = m.storage.Delete(storageKey)
	if err != nil {
		return err
	}
	err = m.indices.Remove(key)
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

type KeyValuePair[K comparable, V any] struct {
	Key   K
	Value V
}

func (m *UnorderedMap[K, V]) Items(startIndex uint64, limit *uint64) ([]KeyValuePair[K, V], error) {
	length := m.keys.Length()
	if startIndex >= length {
		return []KeyValuePair[K, V]{}, nil
	}

	endIndex := length
	if limit != nil {
		if *limit < length-startIndex {
			endIndex = startIndex + *limit
		}
	}

	items := make([]KeyValuePair[K, V], endIndex-startIndex)
	for i := startIndex; i < endIndex; i++ {
		keyData, err := m.keys.Get(i)
		if err != nil {
			return nil, err
		}

		var key K
		err = m.serializer.Deserialize(keyData, &key)
		if err != nil {
			return nil, err
		}

		value, err := m.Get(key)
		if err != nil {
			return nil, err
		}

		items[i-startIndex] = KeyValuePair[K, V]{Key: key, Value: value}
	}

	return items, nil
}

func (m *UnorderedMap[K, V]) Seek(startIndex uint64, limit *uint64) ([]KeyValuePair[K, V], error) {
	return m.Items(startIndex, limit)
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

	return nil
}

type LookupSet[T comparable] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	length     uint64
}

func NewLookupSet[T comparable](prefix string) *LookupSet[T] {
	return &LookupSet[T]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		length:         0,
	}
}

func (s *LookupSet[T]) Length() uint64 {
	return s.length
}

func (s *LookupSet[T]) Contains(value T) (bool, error) {
	storageKey := s.createKey(value)
	return s.storage.HasKey(storageKey)
}

func (s *LookupSet[T]) Insert(value T) error {
	storageKey := s.createKey(value)
	exists, err := s.storage.HasKey(storageKey)
	if err != nil {
		return err
	}

	if !exists {
		data, err := s.serializer.Serialize(true)
		if err != nil {
			return err
		}
		_, err = s.storage.Write(storageKey, data)
		if err != nil {
			return err
		}
		s.length++
	}

	return nil
}

func (s *LookupSet[T]) Remove(value T) error {
	storageKey := s.createKey(value)
	exists, err := s.storage.HasKey(storageKey)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(CollectionErrValueNotFound)
	}

	_, err = s.storage.Delete(storageKey)
	if err != nil {
		return err
	}

	s.length--
	return nil
}

func (s *LookupSet[T]) Clear() error {
	s.length = 0
	return nil
}

type UnorderedSet[T comparable] struct {
	BaseCollection
	storage    Storage
	serializer DefaultSerializer
	values     Vector[[]byte]
	indices    LookupMap[T, uint64]
	length     uint64
}

func NewUnorderedSet[T comparable](prefix string) *UnorderedSet[T] {
	return &UnorderedSet[T]{
		BaseCollection: NewBaseCollection(prefix),
		storage:        DefaultStorage{},
		serializer:     DefaultSerializer{},
		values:         *NewVector[[]byte](prefix + KeySeparator + ValuesPrefix),
		indices:        *NewLookupMap[T, uint64](prefix + KeySeparator + IndicesPrefix),
		length:         0,
	}
}

func (s *UnorderedSet[T]) Length() uint64 {
	return s.length
}

func (s *UnorderedSet[T]) Insert(value T) error {
	storageKey := s.createKey(value)
	exists, err := s.storage.HasKey(storageKey)
	if err != nil {
		return err
	}

	if !exists {
		data, err := s.serializer.Serialize(true)
		if err != nil {
			return err
		}
		_, err = s.storage.Write(storageKey, data)
		if err != nil {
			return err
		}

		index := s.values.Length()
		valueData, err := s.serializer.Serialize(value)
		if err != nil {
			return err
		}
		err = s.values.Push(valueData)
		if err != nil {
			return err
		}

		err = s.indices.Insert(value, index)
		if err != nil {
			return err
		}

		s.length++
	}

	return nil
}

func (s *UnorderedSet[T]) Remove(value T) error {
	storageKey := s.createKey(value)
	exists, err := s.storage.HasKey(storageKey)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(CollectionErrValueNotFound)
	}

	index, err := s.indices.Get(value)
	if err != nil {
		return errors.New(CollectionErrInconsistentState)
	}

	length := s.values.Length()
	if index < length {
		if index != length-1 {
			lastValue, err := s.values.Get(length - 1)
			if err != nil {
				return err
			}
			err = s.values.Set(index, lastValue)
			if err != nil {
				return err
			}
			var lastValueObj T
			err = s.serializer.Deserialize(lastValue, &lastValueObj)
			if err != nil {
				return err
			}
			err = s.indices.Insert(lastValueObj, index)
			if err != nil {
				return err
			}
		}
		_, err = s.values.Pop()
		if err != nil {
			return err
		}
	}

	err = s.indices.Remove(value)
	if err != nil {
		return err
	}
	_, err = s.storage.Delete(storageKey)
	if err != nil {
		return err
	}

	s.length--
	return nil
}

func (s *UnorderedSet[T]) Contains(value T) (bool, error) {
	storageKey := s.createKey(value)
	return s.storage.HasKey(storageKey)
}

func (s *UnorderedSet[T]) Values(startIndex uint64, limit *uint64) ([]T, error) {
	length := s.values.Length()
	if startIndex >= length {
		return []T{}, nil
	}

	endIndex := length
	if limit != nil {
		if *limit < length-startIndex {
			endIndex = startIndex + *limit
		}
	}

	values := make([]T, endIndex-startIndex)
	for i := startIndex; i < endIndex; i++ {
		valueData, err := s.values.Get(i)
		if err != nil {
			return nil, err
		}

		err = s.serializer.Deserialize(valueData, &values[i-startIndex])
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func (s *UnorderedSet[T]) Seek(startIndex uint64, limit *uint64) ([]T, error) {
	return s.Values(startIndex, limit)
}

func (s *UnorderedSet[T]) Clear() error {
	values, err := s.Values(0, nil)
	if err != nil {
		return err
	}

	for _, value := range values {
		storageKey := s.createKey(value)
		_, err = s.storage.Delete(storageKey)
		if err != nil {
			return err
		}

		err = s.indices.Remove(value)
		if err != nil {
			return err
		}
	}

	err = s.values.Clear()
	if err != nil {
		return err
	}

	s.length = 0
	return nil
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
		keys:           *NewVector[[]byte](prefix + KeySeparator + KeysPrefix),
	}
}

func (m *TreeMap[K, V]) findKeyIndex(key K) (int, bool, error) {
	length := m.keys.Length()
	left, right := 0, int(length)-1
	for left <= right {
		mid := (left + right) / 2
		keyData, err := m.keys.Get(uint64(mid))
		if err != nil {
			return 0, false, err
		}
		var midKey K
		err = m.serializer.Deserialize(keyData, &midKey)
		if err != nil {
			return 0, false, err
		}
		cmp := compareKeys(midKey, key)
		if cmp == 0 {
			return mid, true, nil
		} else if cmp < 0 {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left, false, nil
}

func (m *TreeMap[K, V]) insertAtIndex(index int, key K) error {
	length := m.keys.Length()
	keyData, err := m.serializer.Serialize(key)
	if err != nil {
		return err
	}
	err = m.keys.Push(keyData)
	if err != nil {
		return err
	}
	for i := length; i > uint64(index); i-- {
		prev, err := m.keys.Get(i - 1)
		if err != nil {
			return err
		}
		err = m.keys.Set(i, prev)
		if err != nil {
			return err
		}
	}
	return m.keys.Set(uint64(index), keyData)
}

func (m *TreeMap[K, V]) removeAtIndex(index int) error {
	length := m.keys.Length()
	for i := uint64(index); i < length-1; i++ {
		next, err := m.keys.Get(i + 1)
		if err != nil {
			return err
		}
		err = m.keys.Set(i, next)
		if err != nil {
			return err
		}
	}
	_, err := m.keys.Pop()
	return err
}

func (m *TreeMap[K, V]) Set(key K, value V) error {
	return m.Insert(key, value)
}

func (m *TreeMap[K, V]) Insert(key K, value V) error {
	data, err := m.serializer.Serialize(value)
	if err != nil {
		return err
	}
	storageKey := m.createKey(key)
	_, err = m.storage.Write(storageKey, data)
	if err != nil {
		return err
	}
	index, exists, err := m.findKeyIndex(key)
	if err != nil {
		return err
	}
	if !exists {
		return m.insertAtIndex(index, key)
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
	index, exists, err := m.findKeyIndex(key)
	if err != nil {
		return err
	}
	if exists {
		return m.removeAtIndex(index)
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

func (m *TreeMap[K, V]) MinKey() (K, error) {
	if m.keys.Length() == 0 {
		var zero K
		return zero, errors.New(CollectionErrMapEmpty)
	}
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
	length := m.keys.Length()
	if length == 0 {
		var zero K
		return zero, nil
	}
	left, right := 0, int(length)-1
	var result K
	found := false
	for left <= right {
		mid := (left + right) / 2
		keyData, err := m.keys.Get(uint64(mid))
		if err != nil {
			var zero K
			return zero, err
		}
		var midKey K
		err = m.serializer.Deserialize(keyData, &midKey)
		if err != nil {
			var zero K
			return zero, err
		}
		cmp := compareKeys(midKey, key)
		if cmp == 0 {
			return midKey, nil
		} else if cmp < 0 {
			result = midKey
			found = true
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	if found {
		return result, nil
	}
	var zero K
	return zero, nil
}

func (m *TreeMap[K, V]) CeilingKey(key K) (K, error) {
	length := m.keys.Length()
	if length == 0 {
		var zero K
		return zero, nil
	}
	left, right := 0, int(length)-1
	var result K
	found := false
	for left <= right {
		mid := (left + right) / 2
		keyData, err := m.keys.Get(uint64(mid))
		if err != nil {
			var zero K
			return zero, err
		}
		var midKey K
		err = m.serializer.Deserialize(keyData, &midKey)
		if err != nil {
			var zero K
			return zero, err
		}
		cmp := compareKeys(midKey, key)
		if cmp == 0 {
			return midKey, nil
		} else if cmp > 0 {
			result = midKey
			found = true
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	if found {
		return result, nil
	}
	var zero K
	return zero, nil
}

func (m *TreeMap[K, V]) Range(fromKey *K, toKey *K) ([]K, error) {
	allKeys, err := m.Keys()
	if err != nil {
		return nil, err
	}
	startIdx := 0
	endIdx := len(allKeys)
	if fromKey != nil {
		for i, key := range allKeys {
			if compareKeys(key, *fromKey) >= 0 {
				startIdx = i
				break
			}
		}
	}
	if toKey != nil {
		for i, key := range allKeys[startIdx:] {
			if compareKeys(key, *toKey) >= 0 {
				endIdx = startIdx + i
				break
			}
		}
	}
	return allKeys[startIdx:endIdx], nil
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
