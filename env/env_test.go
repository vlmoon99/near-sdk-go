package env

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/system"
)

func TestRegisterAPI(t *testing.T) {
	NearBlockchainImports = system.NewMockSystem()

	NearBlockchainImports.WriteRegisterSys(1, 5, 0)

	length := NearBlockchainImports.RegisterLenSys(1)
	if length != 5 {
		t.Errorf("Expected length 5, got %d", length)
	}

	data, err := readRegisterSafe(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(data) != 5 {
		t.Errorf("Expected 5 bytes, got %d", len(data))
	}
}
