package env

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	SetEnv(system.NewMockSystem())
}

func TestRegisterAPI(t *testing.T) {
	initData := []byte("Bytes")
	writeRegisterSafe(1, initData)

	data, err := readRegisterSafe(1)

	if err != nil {
		t.Errorf("Error: readRegisterSafe failed")
	}

	if string(data) != string(initData) {
		t.Errorf("Error: incorrect data, got %s, expected %s", string(data), string(initData))
	}
}
