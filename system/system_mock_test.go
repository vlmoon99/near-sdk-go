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
