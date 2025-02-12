package env

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/system"
)

func TestRegisterAPI(t *testing.T) {
	NearBlockchainImports = system.NewMockSystem()

	NearBlockchainImports.WriteRegister(1, 5, 0)

	length := NearBlockchainImports.RegisterLen(1)
	if length != 5 {
		t.Errorf("Expected length 5, got %d", length)
	}

}
