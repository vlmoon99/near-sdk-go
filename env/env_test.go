package env

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/system"
)

//TODO : add env tests, but before create system mocks them in according to the Near Runtime functionality

func init() {
	NearBlockchainImports = system.NewMockSystem()
}

func TestRegisterAPI(t *testing.T) {

	NearBlockchainImports.WriteRegister(1, 5, 0)

	length := NearBlockchainImports.RegisterLen(1)
	if length != 5 {
		t.Errorf("Expected length 5, got %d", length)
	}

}
