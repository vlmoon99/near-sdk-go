package collections

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func TestLookupMap_Insert_Get(t *testing.T) {
	m := NewLookupMap([]byte("prefix"))

	key := []byte("key")
	value := "value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	retrievedValue, err := m.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	t.Logf("Retrieved Value: %v", retrievedValue)

	if retrievedValue != value {
		t.Fatalf("Expected value %v, got %v", value, retrievedValue)
	}
}

func TestLookupMap_ContainsKey(t *testing.T) {
	m := NewLookupMap([]byte("prefix"))

	key := []byte("key")
	value := "value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	exists, err := m.ContainsKey(key)
	if err != nil {
		t.Fatalf("ContainsKey failed: %v", err)
	}

	if !exists {
		t.Fatalf("Expected key to exist")
	}

	nonExistentKey := []byte("nonExistentKey")
	exists, err = m.ContainsKey(nonExistentKey)
	if err != nil {
		t.Fatalf("ContainsKey failed for non-existent key: %v", err)
	}

	if exists {
		t.Fatalf("Expected key to not exist")
	}
}

func TestLookupMap_Remove(t *testing.T) {
	m := NewLookupMap([]byte("prefix"))

	key := []byte("key")
	value := "value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	exists, err := m.ContainsKey(key)
	if err != nil {
		t.Fatalf("ContainsKey failed: %v", err)
	}

	if !exists {
		t.Fatalf("Expected key to exist")
	}

	err = m.Remove(key)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	_, err = m.Get(key)
	if err == nil {
		t.Fatalf("Get not failed failed")
	}

}

func TestLookupMap_SerializeDeserialize(t *testing.T) {
	m := NewLookupMap([]byte("prefix"))

	key := []byte("key")
	value := "value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	serializedData, err := m.Serialize()
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	deserializedLookupMap, err := DeserializeLookupMap(serializedData)
	if err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	retrievedValue, err := deserializedLookupMap.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Expected value %v, got %v", value, retrievedValue)
	}
}
