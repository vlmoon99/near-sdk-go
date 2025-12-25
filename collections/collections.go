package collections

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/vlmoon99/near-sdk-go/env"
)

var (
	ErrIndexOutOfBounds   = errors.New("collections: index out of bounds")
	ErrKeyNotFound        = errors.New("collections: key not found")
	ErrVectorEmpty        = errors.New("collections: vector is empty")
	ErrUnsupportedKeyType = errors.New("collections: unsupported key type")
	ErrInconsistentState  = errors.New("collections: inconsistent state")
	ErrMapEmpty           = errors.New("collections: map is empty")
)

func createKey(prefix string, key interface{}) string {
	var keyStr string

	switch k := key.(type) {
	case string:
		keyStr = k
	case uint64:
		keyStr = strconv.FormatUint(k, 10)
	case int:
		keyStr = strconv.Itoa(k)
	case int64:
		keyStr = strconv.FormatInt(k, 10)
	case uint:
		keyStr = strconv.FormatUint(uint64(k), 10)
	case int32:
		keyStr = strconv.FormatInt(int64(k), 10)
	case uint32:
		keyStr = strconv.FormatUint(uint64(k), 10)
	case []byte:
		keyStr = string(k)
	default:
		env.PanicStr("collections: unsupported key type")
		return ""
	}

	return prefix + ":" + keyStr
}

func compareKeys(a, b interface{}) int {
	switch va := a.(type) {
	case string:
		if vb, ok := b.(string); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case int:
		if vb, ok := b.(int); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case uint64:
		if vb, ok := b.(uint64); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	}
	env.PanicStr("collections: unsupported key type for comparison")
	return 0
}

// ==============================================================================
// Vector
// ==============================================================================

type Vector[T any] struct {
	Prefix string `json:"prefix"`
	Len    uint64 `json:"len"`
}

func NewVector[T any](prefix string) *Vector[T] {
	return &Vector[T]{
		Prefix: prefix,
		Len:    0,
	}
}

func (v *Vector[T]) Length() uint64 {
	return v.Len
}

func (v *Vector[T]) Push(value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	key := createKey(v.Prefix, v.Len)
	_, err = env.StorageWrite([]byte(key), data)
	if err != nil {
		return err
	}
	v.Len++
	return nil
}

func (v *Vector[T]) Get(index uint64) (T, error) {
	var zero T
	if index >= v.Len {
		return zero, ErrIndexOutOfBounds
	}
	key := createKey(v.Prefix, index)
	data, err := env.StorageRead([]byte(key))
	if err != nil {
		return zero, err
	}
	err = json.Unmarshal(data, &zero)
	return zero, err
}

func (v *Vector[T]) Set(index uint64, value T) error {
	if index >= v.Len {
		return ErrIndexOutOfBounds
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	key := createKey(v.Prefix, index)
	_, err = env.StorageWrite([]byte(key), data)
	return err
}

func (v *Vector[T]) Pop() (T, error) {
	var zero T
	if v.Len == 0 {
		return zero, ErrVectorEmpty
	}
	lastIndex := v.Len - 1
	item, err := v.Get(lastIndex)
	if err != nil {
		return zero, err
	}
	key := createKey(v.Prefix, lastIndex)
	env.StorageRemove([]byte(key))
	v.Len--
	return item, nil
}

func (v *Vector[T]) Clear() error {
	for i := uint64(0); i < v.Len; i++ {
		key := createKey(v.Prefix, i)
		env.StorageRemove([]byte(key))
	}
	v.Len = 0
	return nil
}

func (v *Vector[T]) ToSlice() ([]T, error) {
	result := make([]T, v.Len)
	for i := uint64(0); i < v.Len; i++ {
		val, err := v.Get(i)
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}

// ==============================================================================
// LookupMap
// ==============================================================================

type LookupMap[K comparable, V any] struct {
	Prefix string `json:"prefix"`
}

func NewLookupMap[K comparable, V any](prefix string) *LookupMap[K, V] {
	return &LookupMap[K, V]{Prefix: prefix}
}

func (m *LookupMap[K, V]) Insert(key K, value V) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	storageKey := createKey(m.Prefix, key)
	_, err = env.StorageWrite([]byte(storageKey), data)
	return err
}

func (m *LookupMap[K, V]) Get(key K) (V, error) {
	var val V
	storageKey := createKey(m.Prefix, key)
	data, err := env.StorageRead([]byte(storageKey))
	if err != nil {
		return val, ErrKeyNotFound
	}
	err = json.Unmarshal(data, &val)
	return val, err
}

func (m *LookupMap[K, V]) Contains(key K) (bool, error) {
	storageKey := createKey(m.Prefix, key)
	return env.StorageHasKey([]byte(storageKey))
}

func (m *LookupMap[K, V]) Remove(key K) error {
	storageKey := createKey(m.Prefix, key)
	_, err := env.StorageRemove([]byte(storageKey))
	return err
}

// ==============================================================================
// LookupSet
// ==============================================================================

type LookupSet[T comparable] struct {
	Prefix string `json:"prefix"`
}

func NewLookupSet[T comparable](prefix string) *LookupSet[T] {
	return &LookupSet[T]{Prefix: prefix}
}

func (s *LookupSet[T]) Insert(value T) error {
	data, _ := json.Marshal(true)
	key := createKey(s.Prefix, value)
	_, err := env.StorageWrite([]byte(key), data)
	return err
}

func (s *LookupSet[T]) Contains(value T) (bool, error) {
	key := createKey(s.Prefix, value)
	return env.StorageHasKey([]byte(key))
}

func (s *LookupSet[T]) Remove(value T) error {
	key := createKey(s.Prefix, value)
	_, err := env.StorageRemove([]byte(key))
	return err
}

// ==============================================================================
// UnorderedMap
// ==============================================================================

type UnorderedMap[K comparable, V any] struct {
	Prefix string `json:"prefix"`
	Len    uint64 `json:"len"`
}

func NewUnorderedMap[K comparable, V any](prefix string) *UnorderedMap[K, V] {
	return &UnorderedMap[K, V]{
		Prefix: prefix,
		Len:    0,
	}
}

func (m *UnorderedMap[K, V]) keyPrefix() string { return m.Prefix + ":k" }
func (m *UnorderedMap[K, V]) valPrefix() string { return m.Prefix + ":v" }
func (m *UnorderedMap[K, V]) idxPrefix() string { return m.Prefix + ":i" }

func (m *UnorderedMap[K, V]) Length() uint64 {
	return m.Len
}

func (m *UnorderedMap[K, V]) Insert(key K, value V) error {
	valData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	valKey := createKey(m.valPrefix(), key)
	idxKey := createKey(m.idxPrefix(), key)

	exists, _ := env.StorageHasKey([]byte(valKey))

	_, err = env.StorageWrite([]byte(valKey), valData)
	if err != nil {
		return err
	}

	if !exists {
		currentIdx := m.Len
		keyVectorKey := createKey(m.keyPrefix(), currentIdx)
		keyData, err := json.Marshal(key)
		if err != nil {
			return err
		}
		env.StorageWrite([]byte(keyVectorKey), keyData)

		idxData, _ := json.Marshal(currentIdx)
		env.StorageWrite([]byte(idxKey), idxData)

		m.Len++
	}
	return nil
}

func (m *UnorderedMap[K, V]) Get(key K) (V, error) {
	var val V
	valKey := createKey(m.valPrefix(), key)
	data, err := env.StorageRead([]byte(valKey))
	if err != nil {
		return val, ErrKeyNotFound
	}
	err = json.Unmarshal(data, &val)
	return val, err
}

func (m *UnorderedMap[K, V]) Remove(key K) error {
	idxKey := createKey(m.idxPrefix(), key)
	valKey := createKey(m.valPrefix(), key)

	idxData, err := env.StorageRead([]byte(idxKey))
	if err != nil {
		return ErrKeyNotFound
	}
	var indexToRemove uint64
	json.Unmarshal(idxData, &indexToRemove)

	lastIndex := m.Len - 1
	if indexToRemove != lastIndex {
		lastKeyVectorKey := createKey(m.keyPrefix(), lastIndex)
		lastKeyData, _ := env.StorageRead([]byte(lastKeyVectorKey))

		var lastKey K
		json.Unmarshal(lastKeyData, &lastKey)

		keyVectorKeyToRemove := createKey(m.keyPrefix(), indexToRemove)
		env.StorageWrite([]byte(keyVectorKeyToRemove), lastKeyData)

		idxKeyForLast := createKey(m.idxPrefix(), lastKey)
		newIdxData, _ := json.Marshal(indexToRemove)
		env.StorageWrite([]byte(idxKeyForLast), newIdxData)
	}

	env.StorageRemove([]byte(createKey(m.keyPrefix(), lastIndex)))
	env.StorageRemove([]byte(idxKey))
	env.StorageRemove([]byte(valKey))

	m.Len--
	return nil
}

func (m *UnorderedMap[K, V]) Keys() ([]K, error) {
	result := make([]K, m.Len)
	for i := uint64(0); i < m.Len; i++ {
		keyVectorKey := createKey(m.keyPrefix(), i)
		data, err := env.StorageRead([]byte(keyVectorKey))
		if err != nil {
			return nil, err
		}
		var k K
		json.Unmarshal(data, &k)
		result[i] = k
	}
	return result, nil
}

func (m *UnorderedMap[K, V]) Values() ([]V, error) {
	keys, err := m.Keys()
	if err != nil {
		return nil, err
	}
	values := make([]V, len(keys))
	for i, k := range keys {
		val, err := m.Get(k)
		if err != nil {
			return nil, err
		}
		values[i] = val
	}
	return values, nil
}

func (m *UnorderedMap[K, V]) Clear() error {
	keys, err := m.Keys()
	if err != nil {
		return err
	}
	for _, k := range keys {
		m.Remove(k)
	}
	m.Len = 0
	return nil
}

// ==============================================================================
// UnorderedSet
// ==============================================================================

type UnorderedSet[T comparable] struct {
	Prefix string `json:"prefix"`
	Len    uint64 `json:"len"`
}

func NewUnorderedSet[T comparable](prefix string) *UnorderedSet[T] {
	return &UnorderedSet[T]{
		Prefix: prefix,
		Len:    0,
	}
}

func (s *UnorderedSet[T]) elemPrefix() string { return s.Prefix + ":e" }
func (s *UnorderedSet[T]) idxPrefix() string  { return s.Prefix + ":i" }

func (s *UnorderedSet[T]) Length() uint64 {
	return s.Len
}

func (s *UnorderedSet[T]) Insert(value T) error {
	idxKey := createKey(s.idxPrefix(), value)
	exists, _ := env.StorageHasKey([]byte(idxKey))
	if exists {
		return nil
	}

	currentIdx := s.Len
	elemKey := createKey(s.elemPrefix(), currentIdx)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	env.StorageWrite([]byte(elemKey), data)

	idxData, _ := json.Marshal(currentIdx)
	env.StorageWrite([]byte(idxKey), idxData)

	s.Len++
	return nil
}

func (s *UnorderedSet[T]) Contains(value T) (bool, error) {
	idxKey := createKey(s.idxPrefix(), value)
	return env.StorageHasKey([]byte(idxKey))
}

func (s *UnorderedSet[T]) Remove(value T) error {
	idxKey := createKey(s.idxPrefix(), value)

	idxData, err := env.StorageRead([]byte(idxKey))
	if err != nil {
		return ErrKeyNotFound
	}
	var indexToRemove uint64
	json.Unmarshal(idxData, &indexToRemove)

	lastIndex := s.Len - 1

	if indexToRemove != lastIndex {
		lastElemKey := createKey(s.elemPrefix(), lastIndex)
		lastElemData, _ := env.StorageRead([]byte(lastElemKey))
		var lastElem T
		json.Unmarshal(lastElemData, &lastElem)

		elemKeyToRemove := createKey(s.elemPrefix(), indexToRemove)
		env.StorageWrite([]byte(elemKeyToRemove), lastElemData)

		idxKeyForLast := createKey(s.idxPrefix(), lastElem)
		newIdxData, _ := json.Marshal(indexToRemove)
		env.StorageWrite([]byte(idxKeyForLast), newIdxData)
	}

	env.StorageRemove([]byte(createKey(s.elemPrefix(), lastIndex)))
	env.StorageRemove([]byte(idxKey))

	s.Len--
	return nil
}

func (s *UnorderedSet[T]) All() ([]T, error) {
	result := make([]T, s.Len)
	for i := uint64(0); i < s.Len; i++ {
		elemKey := createKey(s.elemPrefix(), i)
		data, err := env.StorageRead([]byte(elemKey))
		if err != nil {
			return nil, err
		}
		var val T
		json.Unmarshal(data, &val)
		result[i] = val
	}
	return result, nil
}

func (s *UnorderedSet[T]) Clear() error {
	items, err := s.All()
	if err != nil {
		return err
	}
	for _, item := range items {
		s.Remove(item)
	}
	return nil
}

// ==============================================================================
// TreeMap
// ==============================================================================

type TreeMap[K comparable, V any] struct {
	Prefix string `json:"prefix"`
	Len    uint64 `json:"len"`
}

func NewTreeMap[K comparable, V any](prefix string) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		Prefix: prefix,
		Len:    0,
	}
}

func (m *TreeMap[K, V]) keyPrefix() string { return m.Prefix + ":k" }
func (m *TreeMap[K, V]) valPrefix() string { return m.Prefix + ":v" }

func (m *TreeMap[K, V]) Length() uint64 { return m.Len }

func (m *TreeMap[K, V]) getKeyAt(index uint64) (K, error) {
	var zero K
	keyVecKey := createKey(m.keyPrefix(), index)
	data, err := env.StorageRead([]byte(keyVecKey))
	if err != nil {
		return zero, err
	}
	err = json.Unmarshal(data, &zero)
	return zero, err
}

func (m *TreeMap[K, V]) setKeyAt(index uint64, key K) error {
	keyVecKey := createKey(m.keyPrefix(), index)
	data, err := json.Marshal(key)
	if err != nil {
		return err
	}
	_, err = env.StorageWrite([]byte(keyVecKey), data)
	return err
}

func (m *TreeMap[K, V]) findKeyIndex(key K) (uint64, bool, error) {
	if m.Len == 0 {
		return 0, false, nil
	}
	low, high := uint64(0), m.Len-1

	for low <= high {
		mid := low + (high-low)/2
		midKey, err := m.getKeyAt(mid)
		if err != nil {
			return 0, false, err
		}

		cmp := compareKeys(midKey, key)
		if cmp == 0 {
			return mid, true, nil
		} else if cmp < 0 {
			low = mid + 1
		} else {
			if mid == 0 {
				break
			}
			high = mid - 1
		}
	}
	return low, false, nil
}

func (m *TreeMap[K, V]) Insert(key K, value V) error {
	valKey := createKey(m.valPrefix(), key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = env.StorageWrite([]byte(valKey), data)
	if err != nil {
		return err
	}

	idx, exists, err := m.findKeyIndex(key)
	if err != nil {
		return err
	}

	if !exists {
		for i := m.Len; i > idx; i-- {
			prevKey, err := m.getKeyAt(i - 1)
			if err != nil {
				return err
			}
			m.setKeyAt(i, prevKey)
		}
		m.setKeyAt(idx, key)
		m.Len++
	}

	return nil
}

func (m *TreeMap[K, V]) Get(key K) (V, error) {
	var val V
	valKey := createKey(m.valPrefix(), key)
	data, err := env.StorageRead([]byte(valKey))
	if err != nil {
		return val, ErrKeyNotFound
	}
	err = json.Unmarshal(data, &val)
	return val, err
}

func (m *TreeMap[K, V]) Remove(key K) error {
	idx, exists, err := m.findKeyIndex(key)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	valKey := createKey(m.valPrefix(), key)
	env.StorageRemove([]byte(valKey))

	for i := idx; i < m.Len-1; i++ {
		nextKey, err := m.getKeyAt(i + 1)
		if err != nil {
			return err
		}
		m.setKeyAt(i, nextKey)
	}

	lastVecKey := createKey(m.keyPrefix(), m.Len-1)
	env.StorageRemove([]byte(lastVecKey))
	m.Len--

	return nil
}

func (m *TreeMap[K, V]) MinKey() (K, error) {
	if m.Len == 0 {
		var zero K
		return zero, ErrMapEmpty
	}
	return m.getKeyAt(0)
}

func (m *TreeMap[K, V]) MaxKey() (K, error) {
	if m.Len == 0 {
		var zero K
		return zero, ErrMapEmpty
	}
	return m.getKeyAt(m.Len - 1)
}

func (m *TreeMap[K, V]) Keys() ([]K, error) {
	result := make([]K, m.Len)
	for i := uint64(0); i < m.Len; i++ {
		k, err := m.getKeyAt(i)
		if err != nil {
			return nil, err
		}
		result[i] = k
	}
	return result, nil
}

func (m *TreeMap[K, V]) Clear() error {
	keys, err := m.Keys()
	if err != nil {
		return err
	}
	for _, k := range keys {
		env.StorageRemove([]byte(createKey(m.valPrefix(), k)))
	}
	for i := uint64(0); i < m.Len; i++ {
		env.StorageRemove([]byte(createKey(m.keyPrefix(), i)))
	}
	m.Len = 0
	return nil
}
