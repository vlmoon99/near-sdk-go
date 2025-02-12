package system

import (
	"testing"
	"unsafe"
)

func TestReadRegister(t *testing.T) {
	mockSys := NewMockSystem()
	mockSys.Registers[1] = []byte("test data")

	var buffer = make([]byte, len(mockSys.Registers[1]))
	ptr := uintptr(unsafe.Pointer(&buffer[0]))
	mockSys.ReadRegister(1, uint64(ptr))

	if string(buffer) != "test data" {
		t.Errorf("expected 'test data', got %s", string(buffer))
	}
}

func TestRegisterLen(t *testing.T) {
	mockSys := NewMockSystem()
	mockSys.Registers[1] = []byte("test data")

	length := mockSys.RegisterLen(1)
	expectedLength := uint64(len(mockSys.Registers[1]))

	if length != expectedLength {
		t.Errorf("expected length %d, got %d", expectedLength, length)
	}
}

func TestWriteRegister(t *testing.T) {
	mockSys := NewMockSystem()
	data := []byte("test data")
	var buffer = make([]byte, len(data))
	copy(buffer, data)

	ptr := uintptr(unsafe.Pointer(&buffer[0]))
	mockSys.WriteRegister(1, uint64(len(data)), uint64(ptr))

	if string(mockSys.Registers[1]) != "test data" {
		t.Errorf("expected 'test data', got %s", string(mockSys.Registers[1]))
	}
}

func TestWriteAndReadRegister(t *testing.T) {
	mockSys := NewMockSystem()
	data := []byte("test data")
	var buffer = make([]byte, len(data))
	copy(buffer, data)

	ptr := uintptr(unsafe.Pointer(&buffer[0]))
	mockSys.WriteRegister(1, uint64(len(data)), uint64(ptr))

	var readBuffer = make([]byte, len(data))
	readPtr := uintptr(unsafe.Pointer(&readBuffer[0]))
	mockSys.ReadRegister(1, uint64(readPtr))

	if string(readBuffer) != "test data" {
		t.Errorf("expected 'test data', got %s", string(readBuffer))
	}
}

func TestStorageWrite(t *testing.T) {
	mockSys := NewMockSystem()
	key := "testKey"
	value := "testValue"
	var keyBuffer = []byte(key)
	var valueBuffer = []byte(value)
	keyPtr := uintptr(unsafe.Pointer(&keyBuffer[0]))
	valuePtr := uintptr(unsafe.Pointer(&valueBuffer[0]))

	result := mockSys.StorageWrite(uint64(len(keyBuffer)), uint64(keyPtr), uint64(len(valueBuffer)), uint64(valuePtr), 0)
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}

	if string(mockSys.Storage[key]) != value {
		t.Errorf("expected '%s', got '%s'", value, mockSys.Storage[key])
	}

	// Update existing key
	newValue := "newValue"
	var newValueBuffer = []byte(newValue)
	valuePtr = uintptr(unsafe.Pointer(&newValueBuffer[0]))
	result = mockSys.StorageWrite(uint64(len(keyBuffer)), uint64(keyPtr), uint64(len(newValueBuffer)), uint64(valuePtr), 0)
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}

	if string(mockSys.Storage[key]) != newValue {
		t.Errorf("expected '%s', got '%s'", newValue, mockSys.Storage[key])
	}
}

func TestStorageRead(t *testing.T) {
	mockSys := NewMockSystem()
	key := "testKey"
	value := "testValue"
	mockSys.Storage[key] = []byte(value)
	var keyBuffer = []byte(key)
	keyPtr := uintptr(unsafe.Pointer(&keyBuffer[0]))

	registerId := uint64(1)
	result := mockSys.StorageRead(uint64(len(keyBuffer)), uint64(keyPtr), registerId)
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}

	readValue := mockSys.Registers[registerId]
	if string(readValue) != value {
		t.Errorf("expected '%s', got '%s'", value, readValue)
	}

	// Read non-existing key
	nonExistingKey := "nonExistingKey"
	var nonExistingKeyBuffer = []byte(nonExistingKey)
	keyPtr = uintptr(unsafe.Pointer(&nonExistingKeyBuffer[0]))
	result = mockSys.StorageRead(uint64(len(nonExistingKeyBuffer)), uint64(keyPtr), registerId)
	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func TestStorageRemove(t *testing.T) {
	mockSys := NewMockSystem()
	key := "testKey"
	value := "testValue"
	mockSys.Storage[key] = []byte(value)
	var keyBuffer = []byte(key)
	keyPtr := uintptr(unsafe.Pointer(&keyBuffer[0]))

	registerId := uint64(1)
	result := mockSys.StorageRemove(uint64(len(keyBuffer)), uint64(keyPtr), registerId)
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}

	if _, exists := mockSys.Storage[key]; exists {
		t.Errorf("expected key to be removed")
	}

	// Remove non-existing key
	nonExistingKey := "nonExistingKey"
	var nonExistingKeyBuffer = []byte(nonExistingKey)
	keyPtr = uintptr(unsafe.Pointer(&nonExistingKeyBuffer[0]))
	result = mockSys.StorageRemove(uint64(len(nonExistingKeyBuffer)), uint64(keyPtr), registerId)
	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func TestStorageHasKey(t *testing.T) {
	mockSys := NewMockSystem()
	key := "testKey"
	value := "testValue"
	mockSys.Storage[key] = []byte(value)
	var keyBuffer = []byte(key)
	keyPtr := uintptr(unsafe.Pointer(&keyBuffer[0]))

	result := mockSys.StorageHasKey(uint64(len(keyBuffer)), uint64(keyPtr))
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}

	// Check non-existing key
	nonExistingKey := "nonExistingKey"
	var nonExistingKeyBuffer = []byte(nonExistingKey)
	keyPtr = uintptr(unsafe.Pointer(&nonExistingKeyBuffer[0]))
	result = mockSys.StorageHasKey(uint64(len(nonExistingKeyBuffer)), uint64(keyPtr))
	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}
