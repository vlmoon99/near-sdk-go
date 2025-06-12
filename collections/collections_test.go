package collections

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func cleanupStorage(t *testing.T) {
	t.Helper()
	mockSys := env.NearBlockchainImports.(*system.MockSystem)
	for k := range mockSys.Storage {
		delete(mockSys.Storage, k)
	}
}

func TestVector_Push_Get(t *testing.T) {
	defer cleanupStorage(t)
	v := NewVector[string]("test_vector")

	testValue := "test_value"
	err := v.Push(testValue)
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	retrievedValue, err := v.Get(0)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrievedValue != testValue {
		t.Fatalf("Expected value %v, got %v", testValue, retrievedValue)
	}
}

func TestVector_Length(t *testing.T) {
	defer cleanupStorage(t)
	v := NewVector[string]("test_vector")

	if v.Length() != 0 {
		t.Fatalf("Expected initial length 0, got %d", v.Length())
	}

	err := v.Push("value1")
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	if v.Length() != 1 {
		t.Fatalf("Expected length 1, got %d", v.Length())
	}

	err = v.Push("value2")
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	if v.Length() != 2 {
		t.Fatalf("Expected length 2, got %d", v.Length())
	}
}

func TestVector_Pop(t *testing.T) {
	defer cleanupStorage(t)
	v := NewVector[string]("test_vector")

	_, err := v.Pop()
	if err == nil {
		t.Fatalf("Expected error on empty vector pop")
	}

	testValue := "test_value"
	err = v.Push(testValue)
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	value, err := v.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}

	if value != testValue {
		t.Fatalf("Expected value %v, got %v", testValue, value)
	}

	if v.Length() != 0 {
		t.Fatalf("Expected length 0 after pop, got %d", v.Length())
	}
}

func TestLookupMap_Insert_Get(t *testing.T) {
	defer cleanupStorage(t)
	m := NewLookupMap[string, string]("test_map")

	key := "test_key"
	value := "test_value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	retrievedValue, err := m.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Expected value %v, got %v", value, retrievedValue)
	}
}

func TestLookupMap_Contains(t *testing.T) {
	defer cleanupStorage(t)
	m := NewLookupMap[string, string]("test_map")

	key := "test_key"
	value := "test_value"

	exists, err := m.Contains(key)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected key to not exist")
	}

	err = m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	exists, err = m.Contains(key)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if !exists {
		t.Fatalf("Expected key to exist")
	}
}

func TestLookupMap_Remove(t *testing.T) {
	defer cleanupStorage(t)
	m := NewLookupMap[string, string]("test_map")

	key := "test_key"
	value := "test_value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	err = m.Remove(key)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	exists, err := m.Contains(key)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected key to not exist after removal")
	}
}

func TestLookupMap_MultipleValues(t *testing.T) {
	defer cleanupStorage(t)
	m := NewLookupMap[string, string]("test_map")

	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range testData {
		err := m.Insert(k, v)
		if err != nil {
			t.Fatalf("Insert failed for key %s: %v", k, err)
		}
	}

	for k, expected := range testData {
		value, err := m.Get(k)
		if err != nil {
			t.Fatalf("Get failed for key %s: %v", k, err)
		}
		if value != expected {
			t.Fatalf("Expected value %s for key %s, got %s", expected, k, value)
		}
	}
}

func TestLookupMap_StructValues(t *testing.T) {
	defer cleanupStorage(t)
	type TestValue struct {
		Data string
		Num  int
	}

	m := NewLookupMap[string, TestValue]("test_map")

	value := TestValue{Data: "test", Num: 42}
	err := m.Insert("key1", value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	retrieved, err := m.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved != value {
		t.Fatalf("Expected value %+v, got %+v", value, retrieved)
	}
}

func TestLookupMap_ErrorCases(t *testing.T) {
	defer cleanupStorage(t)
	m := NewLookupMap[string, int]("test_map")

	_, err := m.Get("non_existent")
	if err == nil {
		t.Fatalf("Expected error when getting non-existent key")
	}

	err = m.Remove("non_existent")
	if err == nil {
		t.Fatalf("Expected error when removing non-existent key: %v", err)
	}

	exists, err := m.Contains("non_existent")
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected non-existent key to not exist")
	}
}

func TestUnorderedMap_Insert_Get(t *testing.T) {
	defer cleanupStorage(t)
	m := NewUnorderedMap[string, string]("test_map")

	key := "test_key"
	value := "test_value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	retrievedValue, err := m.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Expected value %v, got %v", value, retrievedValue)
	}
}

func TestUnorderedMap_Contains(t *testing.T) {
	defer cleanupStorage(t)
	m := NewUnorderedMap[string, string]("test_map")

	key := "test_key"
	value := "test_value"

	exists, err := m.Contains(key)

	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected key to not exist")
	}

	err = m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	exists, err = m.Contains(key)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if !exists {
		t.Fatalf("Expected key to exist")
	}
}

func TestUnorderedMap_Remove(t *testing.T) {
	defer cleanupStorage(t)

	m := NewUnorderedMap[string, string]("test_map")

	key := "test_key"
	value := "test_value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	err = m.Remove(key)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	exists, err := m.Contains(key)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected key to not exist after removal")
	}
}

func TestUnorderedMap_Keys_Values(t *testing.T) {
	defer cleanupStorage(t)

	m := NewUnorderedMap[string, string]("test_map")

	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range testData {
		err := m.Insert(k, v)
		if err != nil {
			t.Fatalf("Insert failed for key %s: %v", k, err)
		}
	}

	keys, err := m.Keys()
	if err != nil {
		t.Fatalf("Keys failed: %v", err)
	}
	if len(keys) != len(testData) {
		t.Fatalf("Expected %d keys, got %d", len(testData), len(keys))
	}

	values, err := m.Values()
	if err != nil {
		t.Fatalf("Values failed: %v", err)
	}
	if len(values) != len(testData) {
		t.Fatalf("Expected %d values, got %d", len(testData), len(values))
	}
}

func TestLookupSet_Insert_Contains(t *testing.T) {
	defer cleanupStorage(t)

	s := NewLookupSet[string]("test_set")

	value := "test_value"

	exists, err := s.Contains(value)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected value to not exist")
	}

	err = s.Insert(value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	exists, err = s.Contains(value)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if !exists {
		t.Fatalf("Expected value to exist")
	}
}

func TestLookupSet_Remove(t *testing.T) {
	defer cleanupStorage(t)

	s := NewLookupSet[string]("test_set")

	value := "test_value"

	err := s.Insert(value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	err = s.Remove(value)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	exists, err := s.Contains(value)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected value to not exist after removal")
	}
}

func TestUnorderedSet_Insert_Contains(t *testing.T) {
	defer cleanupStorage(t)

	s := NewUnorderedSet[string]("test_set")

	value := "test_value"

	exists, err := s.Contains(value)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if exists {
		t.Fatalf("Expected value to not exist")
	}

	err = s.Insert(value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	exists, err = s.Contains(value)
	if err != nil {
		t.Fatalf("Contains failed: %v", err)
	}
	if !exists {
		t.Fatalf("Expected value to exist")
	}
}

func TestUnorderedSet_Values(t *testing.T) {
	defer cleanupStorage(t)

	s := NewUnorderedSet[string]("test_set")

	testValues := []string{"value1", "value2", "value3"}

	for _, v := range testValues {
		err := s.Insert(v)
		if err != nil {
			t.Fatalf("Insert failed for value %s: %v", v, err)
		}
	}

	values, err := s.Values()
	if err != nil {
		t.Fatalf("Values failed: %v", err)
	}
	if len(values) != len(testValues) {
		t.Fatalf("Expected %d values, got %d", len(testValues), len(values))
	}
}

func TestTreeMap_Insert_Get(t *testing.T) {
	defer cleanupStorage(t)

	m := NewTreeMap[string, string]("test_map")

	key := "test_key"
	value := "test_value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	retrievedValue, err := m.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Expected value %v, got %v", value, retrievedValue)
	}
}

func TestTreeMap_MinKey_MaxKey(t *testing.T) {
	defer cleanupStorage(t)

	m := NewTreeMap[string, string]("test_map")

	env.LogString("TestTreeMap_MinKey_MaxKey: Testing empty map MinKey")
	_, err := m.MinKey()
	if err == nil {
		t.Fatalf("Expected error for empty map MinKey")
	}
	if err.Error() != CollectionErrMapEmpty {
		t.Fatalf("Expected error %s, got %s", CollectionErrMapEmpty, err.Error())
	}

	env.LogString("TestTreeMap_MinKey_MaxKey: Testing empty map MaxKey")
	_, err = m.MaxKey()
	if err == nil {
		t.Fatalf("Expected error for empty map MaxKey")
	}
	if err.Error() != CollectionErrMapEmpty {
		t.Fatalf("Expected error %s, got %s", CollectionErrMapEmpty, err.Error())
	}

	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	env.LogString("TestTreeMap_MinKey_MaxKey: Inserting test data")
	for k, v := range testData {
		err := m.Insert(k, v)
		if err != nil {
			t.Fatalf("Insert failed for key %s: %v", k, err)
		}
	}

	env.LogString("TestTreeMap_MinKey_MaxKey: Testing MinKey")
	minKey, err := m.MinKey()
	if err != nil {
		t.Fatalf("MinKey failed: %v", err)
	}

	if minKey != "key1" {
		t.Fatalf("Expected min key 'key1', got %v", minKey)
	}

	env.LogString("TestTreeMap_MinKey_MaxKey: Testing MaxKey")
	maxKey, err := m.MaxKey()
	if err != nil {
		t.Fatalf("MaxKey failed: %v", err)
	}
	if maxKey != "key3" {
		t.Fatalf("Expected max key 'key3', got %v", maxKey)
	}
}

func TestTreeMap_FloorKey_CeilingKey(t *testing.T) {
	defer cleanupStorage(t)

	m := NewTreeMap[string, string]("test_map")

	testData := map[string]string{
		"key1": "value1",
		"key3": "value3",
		"key5": "value5",
	}

	for k, v := range testData {
		err := m.Insert(k, v)
		if err != nil {
			t.Fatalf("Insert failed for key %s: %v", k, err)
		}
	}

	floorKey, err := m.FloorKey("key4")
	if err != nil {
		t.Fatalf("FloorKey failed: %v", err)
	}
	if floorKey != "key3" {
		t.Fatalf("Expected floor key 'key3', got %v", floorKey)
	}

	ceilingKey, err := m.CeilingKey("key2")
	if err != nil {
		t.Fatalf("CeilingKey failed: %v", err)
	}
	if ceilingKey != "key3" {
		t.Fatalf("Expected ceiling key 'key3', got %v", ceilingKey)
	}
}
