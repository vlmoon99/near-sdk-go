package collections

import (
	"fmt"
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

// cleanupStorage clears the storage after each test
func cleanupStorage(t *testing.T) {
	t.Helper()
	mockSys := env.NearBlockchainImports.(*system.MockSystem)
	for k := range mockSys.Storage {
		delete(mockSys.Storage, k)
	}
}

// Vector Tests
func TestVector_Push_Get(t *testing.T) {
	defer cleanupStorage(t)
	v := NewVector("test_vector")

	testValue := "test_value"
	err := v.Push(testValue)
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	var retrievedValue string
	err = v.Get(0, &retrievedValue)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrievedValue != testValue {
		t.Fatalf("Expected value %v, got %v", testValue, retrievedValue)
	}
}

func TestVector_Length(t *testing.T) {
	defer cleanupStorage(t)
	v := NewVector("test_vector")

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
	v := NewVector("test_vector")

	var value string
	err := v.Pop(&value)
	if err == nil {
		t.Fatalf("Expected error on empty vector pop")
	}

	testValue := "test_value"
	err = v.Push(testValue)
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	fmt.Printf("testValue: %v\n", testValue)

	err = v.Pop(&value)
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

// UnorderedMap Tests
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

	fmt.Printf("Test data: %+v\n", testData)

	for k, v := range testData {
		fmt.Printf("Inserting key: %s, value: %s\n", k, v)
		err := m.Insert(k, v)
		if err != nil {
			t.Fatalf("Insert failed for key %s: %v", k, err)
		}
	}

	// Verify storage after insert
	mockSys := env.NearBlockchainImports.(*system.MockSystem)
	fmt.Printf("Storage after insert: %+v\n", mockSys.Storage)

	// Check keys Vector state
	fmt.Printf("Keys Vector length: %d\n", m.keys.Length())
	fmt.Printf("Keys Vector prefix: %s\n", m.keys.GetPrefix())

	// Print all storage keys for debugging
	fmt.Println("All storage keys:")
	for k := range mockSys.Storage {
		fmt.Printf("  %s\n", k)
	}

	keys, err := m.Keys()
	if err != nil {
		t.Fatalf("Keys failed: %v", err)
	}
	fmt.Printf("Retrieved keys: %+v\n", keys)
	if len(keys) != len(testData) {
		t.Fatalf("Expected %d keys, got %d", len(testData), len(keys))
	}

	values, err := m.Values()
	if err != nil {
		t.Fatalf("Values failed: %v", err)
	}
	fmt.Printf("Retrieved values: %+v\n", values)
	if len(values) != len(testData) {
		t.Fatalf("Expected %d values, got %d", len(testData), len(values))
	}
}

// // LookupSet Tests
// func TestLookupSet_Insert_Contains(t *testing.T) {
// 	s := NewLookupSet("test_set")

// 	value := "test_value"

// 	exists, err := s.Contains(value)
// 	if err != nil {
// 		t.Fatalf("Contains failed: %v", err)
// 	}
// 	if exists {
// 		t.Fatalf("Expected value to not exist")
// 	}

// 	err = s.Insert(value)
// 	if err != nil {
// 		t.Fatalf("Insert failed: %v", err)
// 	}

// 	exists, err = s.Contains(value)
// 	if err != nil {
// 		t.Fatalf("Contains failed: %v", err)
// 	}
// 	if !exists {
// 		t.Fatalf("Expected value to exist")
// 	}
// }

// func TestLookupSet_Remove(t *testing.T) {
// 	s := NewLookupSet("test_set")

// 	value := "test_value"

// 	err := s.Insert(value)
// 	if err != nil {
// 		t.Fatalf("Insert failed: %v", err)
// 	}

// 	err = s.Remove(value)
// 	if err != nil {
// 		t.Fatalf("Remove failed: %v", err)
// 	}

// 	exists, err := s.Contains(value)
// 	if err != nil {
// 		t.Fatalf("Contains failed: %v", err)
// 	}
// 	if exists {
// 		t.Fatalf("Expected value to not exist after removal")
// 	}
// }

// // UnorderedSet Tests
// func TestUnorderedSet_Insert_Contains(t *testing.T) {
// 	s := NewUnorderedSet("test_set")

// 	value := "test_value"

// 	// Test before insert
// 	exists, err := s.Contains(value)
// 	if err != nil {
// 		t.Fatalf("Contains failed: %v", err)
// 	}
// 	if exists {
// 		t.Fatalf("Expected value to not exist")
// 	}

// 	// Test after insert
// 	err = s.Insert(value)
// 	if err != nil {
// 		t.Fatalf("Insert failed: %v", err)
// 	}

// 	exists, err = s.Contains(value)
// 	if err != nil {
// 		t.Fatalf("Contains failed: %v", err)
// 	}
// 	if !exists {
// 		t.Fatalf("Expected value to exist")
// 	}
// }

// func TestUnorderedSet_Values(t *testing.T) {
// 	s := NewUnorderedSet("test_set")

// 	// Insert test data
// 	testValues := []string{"value1", "value2", "value3"}
// 	for _, v := range testValues {
// 		err := s.Insert(v)
// 		if err != nil {
// 			t.Fatalf("Insert failed for value %s: %v", v, err)
// 		}
// 	}

// 	// Test Values
// 	values, err := s.Values()
// 	if err != nil {
// 		t.Fatalf("Values failed: %v", err)
// 	}
// 	if len(values) != len(testValues) {
// 		t.Fatalf("Expected %d values, got %d", len(testValues), len(values))
// 	}
// }

// // TreeMap Tests
// func TestTreeMap_Insert_Get(t *testing.T) {
// 	m := NewTreeMap("test_map")

// 	key := "test_key"
// 	value := "test_value"

// 	err := m.Insert(key, value)
// 	if err != nil {
// 		t.Fatalf("Insert failed: %v", err)
// 	}

// 	var retrievedValue string
// 	err = m.Get(key, &retrievedValue)
// 	if err != nil {
// 		t.Fatalf("Get failed: %v", err)
// 	}

// 	if retrievedValue != value {
// 		t.Fatalf("Expected value %v, got %v", value, retrievedValue)
// 	}
// }

// func TestTreeMap_MinKey_MaxKey(t *testing.T) {
// 	m := NewTreeMap("test_map")

// 	// Test empty map
// 	_, err := m.MinKey()
// 	if err == nil {
// 		t.Fatalf("Expected error for empty map MinKey")
// 	}

// 	_, err = m.MaxKey()
// 	if err == nil {
// 		t.Fatalf("Expected error for empty map MaxKey")
// 	}

// 	// Insert test data
// 	testData := map[string]string{
// 		"key1": "value1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}

// 	for k, v := range testData {
// 		err := m.Insert(k, v)
// 		if err != nil {
// 			t.Fatalf("Insert failed for key %s: %v", k, err)
// 		}
// 	}

// 	// Test MinKey
// 	minKey, err := m.MinKey()
// 	if err != nil {
// 		t.Fatalf("MinKey failed: %v", err)
// 	}
// 	if minKey != "key1" {
// 		t.Fatalf("Expected min key 'key1', got %v", minKey)
// 	}

// 	// Test MaxKey
// 	maxKey, err := m.MaxKey()
// 	if err != nil {
// 		t.Fatalf("MaxKey failed: %v", err)
// 	}
// 	if maxKey != "key3" {
// 		t.Fatalf("Expected max key 'key3', got %v", maxKey)
// 	}
// }

// func TestTreeMap_FloorKey_CeilingKey(t *testing.T) {
// 	m := NewTreeMap("test_map")

// 	// Insert test data
// 	testData := map[string]string{
// 		"key1": "value1",
// 		"key3": "value3",
// 		"key5": "value5",
// 	}

// 	for k, v := range testData {
// 		err := m.Insert(k, v)
// 		if err != nil {
// 			t.Fatalf("Insert failed for key %s: %v", k, err)
// 		}
// 	}

// 	// Test FloorKey
// 	floorKey, err := m.FloorKey("key4")
// 	if err != nil {
// 		t.Fatalf("FloorKey failed: %v", err)
// 	}
// 	if floorKey != "key3" {
// 		t.Fatalf("Expected floor key 'key3', got %v", floorKey)
// 	}

// 	// Test CeilingKey
// 	ceilingKey, err := m.CeilingKey("key2")
// 	if err != nil {
// 		t.Fatalf("CeilingKey failed: %v", err)
// 	}
// 	if ceilingKey != "key3" {
// 		t.Fatalf("Expected ceiling key 'key3', got %v", ceilingKey)
// 	}
// }
