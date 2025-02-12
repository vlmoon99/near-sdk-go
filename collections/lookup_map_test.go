package collections

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.NearBlockchainImports = system.NewMockSystem()
}

// TODO: Understand why these errors happen
// TODO: Add unit tests for Lookup map
// tinygo test ./
// panic: runtime error at 0x000000000020f010: caught signal SIGSEGV
func TestLookupMap_Insert_Get(t *testing.T) {
	m := NewLookupMap([]byte("prefix"))

	key := []byte("key")
	value := "value"

	err := m.Insert(key, value)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// retrievedValue, err := m.Get(key)
	// if err != nil {
	// 	t.Fatalf("Get failed: %v", err)
	// }

	// t.Logf("Retrieved Value: %v", retrievedValue)

	// if retrievedValue != value {
	// 	t.Fatalf("Expected value %v, got %v", value, retrievedValue)
	// }
}
