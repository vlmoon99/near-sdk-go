package system

import (
	"bytes"
	"testing"
	"unsafe"
)

// Registers API

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

// Registers API

// Storage API
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

// Storage API

// Context API

func TestCurrentAccountId(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	mockSys.CurrentAccountId(registerId)

	if data, exists := mockSys.Registers[registerId]; !exists || string(data) != mockSys.CurrentAccountIdSys {
		t.Errorf("expected %s, got %s", mockSys.CurrentAccountIdSys, string(data))
	}
}

func TestSignerAccountId(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	mockSys.SignerAccountId(registerId)

	if data, exists := mockSys.Registers[registerId]; !exists || string(data) != mockSys.SignerAccountIdSys {
		t.Errorf("expected %s, got %s", mockSys.SignerAccountIdSys, string(data))
	}
}

func TestSignerAccountPk(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	mockSys.SignerAccountPk(registerId)

	if data, exists := mockSys.Registers[registerId]; !exists || !bytes.Equal(data, mockSys.SignerAccountPkSys) {
		t.Errorf("expected %s, got %s", string(mockSys.SignerAccountPkSys), string(data))
	}
}

func TestPredecessorAccountId(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	mockSys.PredecessorAccountId(registerId)

	if data, exists := mockSys.Registers[registerId]; !exists || string(data) != mockSys.PredecessorAccountIdSys {
		t.Errorf("expected %s, got %s", mockSys.PredecessorAccountIdSys, string(data))
	}
}

func TestInput(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	mockSys.Input(registerId)

	if data, exists := mockSys.Registers[registerId]; !exists || !bytes.Equal(data, mockSys.ContractInput) {
		t.Errorf("expected %s, got %s", string(mockSys.ContractInput), string(data))
	}
}

func TestBlockIndex(t *testing.T) {
	mockSys := NewMockSystem()
	expected := mockSys.BlockIndexSys
	result := mockSys.BlockIndex()

	if expected != result {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestBlockTimestamp(t *testing.T) {
	mockSys := NewMockSystem()
	expected := mockSys.BlockTimestampSys
	result := mockSys.BlockTimestamp()

	if expected != result {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestEpochHeight(t *testing.T) {
	mockSys := NewMockSystem()
	expected := mockSys.EpochHeightSys
	result := mockSys.EpochHeight()

	if expected != result {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestStorageUsage(t *testing.T) {
	mockSys := NewMockSystem()
	expected := mockSys.StorageUsageSys
	result := mockSys.StorageUsage()

	if expected != result {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

// Context API

// Economics API

func TestAccountBalance(t *testing.T) {
	mockSys := NewMockSystem()
	var data [16]byte
	mockSys.AccountBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))

	expected := mockSys.AccountBalanceSys.ToBE()
	if string(data[:]) != string(expected) {
		t.Errorf("expected %v, got %v", expected, data[:])
	}
}

func TestAccountLockedBalance(t *testing.T) {
	mockSys := NewMockSystem()
	var data [16]byte
	mockSys.AccountLockedBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))

	expected := mockSys.AccountLockedBalanceSys.ToBE()
	if string(data[:]) != string(expected) {
		t.Errorf("expected %v, got %v", expected, data[:])
	}
}

func TestAttachedDeposit(t *testing.T) {
	mockSys := NewMockSystem()
	var data [16]byte
	mockSys.AttachedDeposit(uint64(uintptr(unsafe.Pointer(&data[0]))))

	expected := mockSys.AttachedDepositSys.ToBE()
	if string(data[:]) != string(expected) {
		t.Errorf("expected %v, got %v", expected, data[:])
	}
}

func TestPrepaidGas(t *testing.T) {
	mockSys := NewMockSystem()
	expected := mockSys.PrepaidGasSys
	result := mockSys.PrepaidGas()

	if expected != result {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestUsedGas(t *testing.T) {
	mockSys := NewMockSystem()
	expected := mockSys.UsedGasSys
	result := mockSys.UsedGas()

	if expected != result {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

// Economics API

// Math API

func TestRandomSeed(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	mockSys.RandomSeed(registerId)

	if _, exists := mockSys.Registers[registerId]; !exists {
		t.Errorf("expected random seed to be written to register %d", registerId)
	}
}

func TestSha256(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	data := []byte("test data")
	dataPtr := uintptr(unsafe.Pointer(&data[0]))

	mockSys.Sha256(uint64(len(data)), uint64(dataPtr), registerId)

	expected := "hash"
	if string(mockSys.Registers[registerId]) != expected {
		t.Errorf("expected %s, got %s", expected, string(mockSys.Registers[registerId]))
	}
}

func TestKeccak256(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	data := []byte("test data")
	dataPtr := uintptr(unsafe.Pointer(&data[0]))

	mockSys.Keccak256(uint64(len(data)), uint64(dataPtr), registerId)

	expected := "hash"
	if string(mockSys.Registers[registerId]) != expected {
		t.Errorf("expected %s, got %s", expected, string(mockSys.Registers[registerId]))
	}
}

// func TestKeccak512(t *testing.T) {
// 	mockSys := NewMockSystem()
// 	registerId := uint64(1)
// 	data := []byte("test data")
// 	dataPtr := uintptr(unsafe.Pointer(&data[0]))

// 	mockSys.Keccak512(uint64(len(data)), uint64(dataPtr), registerId)

// 	expected := "hash"
// 	if string(mockSys.Registers[registerId]) != expected {
// 		t.Errorf("expected %s, got %s", expected, string(mockSys.Registers[registerId]))
// 	}
// }

// func TestRipemd160(t *testing.T) {
// 	mockSys := NewMockSystem()
// 	registerId := uint64(1)
// 	data := []byte("test data")
// 	dataPtr := uintptr(unsafe.Pointer(&data[0]))

// 	mockSys.Ripemd160(uint64(len(data)), uint64(dataPtr), registerId)

// 	expected := "hash"
// 	if string(mockSys.Registers[registerId]) != expected {
// 		t.Errorf("expected %s, got %s", expected, string(mockSys.Registers[registerId]))
// 	}
// }

// func TestAltBn128G1Multiexp(t *testing.T) {
// 	mockSys := NewMockSystem()
// 	registerId := uint64(1)
// 	data := []byte{1, 2, 3, 4, 5}
// 	dataPtr := uintptr(unsafe.Pointer(&data[0]))

// 	mockSys.AltBn128G1Multiexp(uint64(len(data)), uint64(dataPtr), registerId)

// 	expected := "simpleMultiexp"

// 	if string(mockSys.Registers[registerId]) != expected {
// 		t.Errorf("expected 'simpleMultiexp', got %s", string(mockSys.Registers[registerId]))
// 	}
// }

// func TestAltBn128G1SumSystem(t *testing.T) {
// 	mockSys := NewMockSystem()
// 	registerId := uint64(1)
// 	data := []byte{1, 2, 3, 4, 5}
// 	dataPtr := uintptr(unsafe.Pointer(&data[0]))

// 	mockSys.AltBn128G1SumSystem(uint64(len(data)), uint64(dataPtr), registerId)

// 	expected := "simpleSum"

// 	if string(mockSys.Registers[registerId]) != expected {
// 		t.Errorf("expected 'simpleSum', got %s", string(mockSys.Registers[registerId]))
// 	}
// }

// Math API
