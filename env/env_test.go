package env

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	SetEnv(system.NewMockSystem())
}

func TestSetEnv(t *testing.T) {
	mockSys := system.NewMockSystem()
	SetEnv(mockSys)

	if nearBlockchainImports != mockSys {
		t.Errorf("expected nearBlockchainImports to be set to mockSys, got %v", nearBlockchainImports)
	}
}

func TestTryMethodIntoRegister(t *testing.T) {
	mockSys := system.NewMockSystem()
	SetEnv(mockSys)

	data := []byte("test data")
	mockSys.Registers[AtomicOpRegister] = data

	method := func(registerId uint64) {
		writeRegisterSafe(registerId, data)
	}

	result, err := tryMethodIntoRegister(method)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if string(result) != string(data) {
		t.Errorf("expected '%s', got '%s'", data, result)
	}
}

func TestMethodIntoRegister(t *testing.T) {
	mockSys := system.NewMockSystem()
	SetEnv(mockSys)

	data := []byte("test data")
	mockSys.Registers[AtomicOpRegister] = data

	method := func(registerId uint64) {
		writeRegisterSafe(registerId, data)
	}

	result, err := methodIntoRegister(method)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if string(result) != string(data) {
		t.Errorf("expected '%s', got '%s'", data, result)
	}
}

func TestReadRegisterSafe(t *testing.T) {
	mockSys := system.NewMockSystem()
	SetEnv(mockSys)

	data := []byte("test data")
	mockSys.Registers[AtomicOpRegister] = data

	result, err := readRegisterSafe(AtomicOpRegister)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if string(result) != string(data) {
		t.Errorf("expected '%s', got '%s'", data, result)
	}

	// Test with an empty register
	result, err = readRegisterSafe(1)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got '%s'", result)
	}
}

func TestWriteRegisterSafe(t *testing.T) {
	mockSys := system.NewMockSystem()
	SetEnv(mockSys)

	data := []byte("test data")
	writeRegisterSafe(1, data)

	if string(mockSys.Registers[1]) != string(data) {
		t.Errorf("expected '%s', got '%s'", data, mockSys.Registers[1])
	}

	// Test with empty data
	writeRegisterSafe(2, []byte{})
	if _, exists := mockSys.Registers[2]; exists {
		t.Errorf("expected register 2 to be empty")
	}
}

func TestStorageWrite(t *testing.T) {
	initSystem := system.NewMockSystem()
	SetEnv(initSystem)

	key := []byte("testKey")
	value := []byte("testValue")

	success, err := StorageWrite(key, value)
	if !success || err != nil {
		t.Errorf("expected successful write, got error: %v", err)
	}

	if string(initSystem.Storage["testKey"]) != "testValue" {
		t.Errorf("expected value 'testValue', got '%s'", string(initSystem.Storage["testKey"]))
	}
}

func TestStorageRead(t *testing.T) {
	initSystem := system.NewMockSystem()
	SetEnv(initSystem)

	key := []byte("testKey")
	expectedValue := []byte("testValue")
	initSystem.Storage[string(key)] = expectedValue

	value, err := StorageRead(key)
	if err != nil {
		t.Errorf("expected successful read, got error: %v", err)
	}

	if string(value) != string(expectedValue) {
		t.Errorf("expected value '%s', got '%s'", string(expectedValue), string(value))
	}
}

func TestStorageRemove(t *testing.T) {
	initSystem := system.NewMockSystem()
	SetEnv(initSystem)

	key := []byte("testKey")
	value := []byte("testValue")
	initSystem.Storage[string(key)] = value

	success, err := StorageRemove(key)
	if !success || err != nil {
		t.Errorf("expected successful remove, got error: %v", err)
	}

	if _, exists := initSystem.Storage[string(key)]; exists {
		t.Errorf("expected key to be removed, but it still exists")
	}
}

func TestStorageHasKey(t *testing.T) {
	initSystem := system.NewMockSystem()
	SetEnv(initSystem)

	key := []byte("testKey")
	value := []byte("testValue")
	initSystem.Storage[string(key)] = value

	hasKey, err := StorageHasKey(key)
	if err != nil {
		t.Errorf("expected successful has key check, got error: %v", err)
	}

	if !hasKey {
		t.Errorf("expected key to exist, but it does not")
	}

	// Check non-existing key
	nonExistingKey := []byte("nonExistingKey")
	hasKey, err = StorageHasKey(nonExistingKey)
	if err != nil {
		t.Errorf("expected successful has key check, got error: %v", err)
	}

	if hasKey {
		t.Errorf("expected key not to exist, but it does")
	}
}
