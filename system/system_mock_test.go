package system

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/vlmoon99/near-sdk-go/types"
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

	expected := mockSys.AccountBalanceSys.ToLE()
	if string(data[:]) != string(expected) {
		t.Errorf("expected %v, got %v", expected, data[:])
	}
}

func TestAccountLockedBalance(t *testing.T) {
	mockSys := NewMockSystem()
	var data [16]byte
	mockSys.AccountLockedBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))

	expected := mockSys.AccountLockedBalanceSys.ToLE()
	if string(data[:]) != string(expected) {
		t.Errorf("expected %v, got %v", expected, data[:])
	}
}

func TestAttachedDeposit(t *testing.T) {
	mockSys := NewMockSystem()
	var data [16]byte
	mockSys.AttachedDeposit(uint64(uintptr(unsafe.Pointer(&data[0]))))

	expected := mockSys.AttachedDepositSys.ToLE()
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

func TestKeccak512(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	data := []byte("test data")
	dataPtr := uintptr(unsafe.Pointer(&data[0]))

	mockSys.Keccak512(uint64(len(data)), uint64(dataPtr), registerId)

	expected := "hash"
	if string(mockSys.Registers[registerId]) != expected {
		t.Errorf("expected %s, got %s", expected, string(mockSys.Registers[registerId]))
	}
}

func TestRipemd160(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	data := []byte("test data")
	dataPtr := uintptr(unsafe.Pointer(&data[0]))

	mockSys.Ripemd160(uint64(len(data)), uint64(dataPtr), registerId)

	expected := "hash"
	if string(mockSys.Registers[registerId]) != expected {
		t.Errorf("expected %s, got %s", expected, string(mockSys.Registers[registerId]))
	}
}

func TestAltBn128G1Multiexp(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	data := []byte{1, 2, 3, 4, 5}
	dataPtr := uintptr(unsafe.Pointer(&data[0]))

	mockSys.AltBn128G1Multiexp(uint64(len(data)), uint64(dataPtr), registerId)

	expected := "simpleMultiexp"

	if string(mockSys.Registers[registerId]) != expected {
		t.Errorf("expected 'simpleMultiexp', got %s", string(mockSys.Registers[registerId]))
	}
}

func TestAltBn128G1SumSystem(t *testing.T) {
	mockSys := NewMockSystem()
	registerId := uint64(1)
	data := []byte{1, 2, 3, 4, 5}
	dataPtr := uintptr(unsafe.Pointer(&data[0]))

	mockSys.AltBn128G1SumSystem(uint64(len(data)), uint64(dataPtr), registerId)

	expected := "simpleSum"

	if string(mockSys.Registers[registerId]) != expected {
		t.Errorf("expected 'simpleSum', got %s", string(mockSys.Registers[registerId]))
	}
}

// Math API

// Validator API

func TestValidatorStake(t *testing.T) {
	mockSys := NewMockSystem()
	accountId := "validatorAccountId"
	accountIdLen := uint64(len(accountId))
	accountIdPtr := uintptr(unsafe.Pointer(&accountId))

	var stakeData uint64
	stakePtr := uintptr(unsafe.Pointer(&stakeData))

	mockSys.ValidatorStake(accountIdLen, uint64(accountIdPtr), uint64(stakePtr))

	expectedStake := uint64(100000)
	if stakeData != expectedStake {
		t.Errorf("expected stake %d, got %d", expectedStake, stakeData)
	}
}

func TestValidatorTotalStake(t *testing.T) {
	mockSys := NewMockSystem()
	var stakeData uint64
	stakePtr := uintptr(unsafe.Pointer(&stakeData))

	mockSys.ValidatorTotalStake(uint64(stakePtr))

	expectedStake := uint64(100000)
	if stakeData != expectedStake {
		t.Errorf("expected total stake %d, got %d", expectedStake, stakeData)
	}
}

// Validator API

// Miscellaneous API

func TestValueReturn(t *testing.T) {
	mockSys := NewMockSystem()
	data := []byte("test data")
	var buffer = make([]byte, len(data))
	copy(buffer, data)

	ptr := uintptr(unsafe.Pointer(&buffer[0]))

	mockSys.ValueReturn(uint64(len(data)), uint64(ptr))

	if string(mockSys.Registers[0]) != string(data) {
		t.Errorf("expected %s, got %s", string(data), string(mockSys.Registers[0]))
	}
}

// Miscellaneous API

// Promises API

func TestPromiseCreate(t *testing.T) {
	mockSys := NewMockSystem()
	accountId := "accountId"
	functionName := "functionName"
	arguments := []byte("arguments")
	amount := types.Uint128{Lo: 0, Hi: 0}
	gas := uint64(5000)

	accountIdPtr := uintptr(unsafe.Pointer(&accountId))
	functionNamePtr := uintptr(unsafe.Pointer(&functionName))
	argumentsPtr := uintptr(unsafe.Pointer(&arguments[0]))
	amountPtr := uintptr(unsafe.Pointer(&amount))

	promiseIndex := mockSys.PromiseCreate(uint64(len(accountId)), uint64(accountIdPtr),
		uint64(len(functionName)), uint64(functionNamePtr),
		uint64(len(arguments)), uint64(argumentsPtr),
		uint64(amountPtr), gas)

	expectedIndex := uint64(0)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}

	if len(mockSys.Promises) != 1 {
		t.Errorf("expected 1 promise, got %d", len(mockSys.Promises))
	}

	promise := mockSys.Promises[0]
	if promise.AccountId != accountId {
		t.Errorf("expected account id %s, got %s", accountId, promise.AccountId)
	}
	if promise.FunctionName != functionName {
		t.Errorf("expected function name %s, got %s", functionName, promise.FunctionName)
	}
	if string(promise.Arguments) != string(arguments) {
		t.Errorf("expected arguments %s, got %s", string(arguments), string(promise.Arguments))
	}
	if promise.Amount != amount {
		t.Errorf("expected amount %d, got %d", amount, promise.Amount)
	}
	if promise.Gas != gas {
		t.Errorf("expected gas %d, got %d", gas, promise.Gas)
	}
}

func TestPromiseThen(t *testing.T) {
	mockSys := NewMockSystem()

	// Create the first promise
	accountId := "accountId"
	functionName := "functionName"
	arguments := []byte("arguments")
	amount := types.Uint128{Lo: 0, Hi: 0}
	gas := uint64(5000)

	accountIdPtr := uintptr(unsafe.Pointer(&accountId))
	functionNamePtr := uintptr(unsafe.Pointer(&functionName))
	argumentsPtr := uintptr(unsafe.Pointer(&arguments[0]))
	amountPtr := uintptr(unsafe.Pointer(&amount))

	// Create the first promise
	mockSys.PromiseCreate(
		uint64(len(accountId)), uint64(accountIdPtr),
		uint64(len(functionName)), uint64(functionNamePtr),
		uint64(len(arguments)), uint64(argumentsPtr),
		uint64(amountPtr), gas,
	)

	promiseIdx := uint64(0)
	promiseIndex := mockSys.PromiseThen(
		promiseIdx,
		uint64(len(accountId)), uint64(accountIdPtr),
		uint64(len(functionName)), uint64(functionNamePtr),
		uint64(len(arguments)), uint64(argumentsPtr),
		uint64(amountPtr), gas,
	)

	expectedIndex := uint64(1)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}

	if len(mockSys.Promises) != 2 {
		t.Errorf("expected 2 promises, got %d", len(mockSys.Promises))
	}

	promise := mockSys.Promises[1]
	if promise.AccountId != accountId {
		t.Errorf("expected account id %s, got %s", accountId, promise.AccountId)
	}
	if promise.FunctionName != functionName {
		t.Errorf("expected function name %s, got %s", functionName, promise.FunctionName)
	}
	if string(promise.Arguments) != string(arguments) {
		t.Errorf("expected arguments %s, got %s", string(arguments), string(promise.Arguments))
	}
	if promise.Amount.Lo != amount.Lo {
		t.Errorf("expected amount %d, got %d", amount.Lo, promise.Amount)
	}
	if promise.Gas != gas {
		t.Errorf("expected gas %d, got %d", gas, promise.Gas)
	}
}

func TestPromiseAnd(t *testing.T) {
	mockSys := NewMockSystem()
	promiseIndices := []uint64{0, 1}
	promiseIdxPtr := uintptr(unsafe.Pointer(&promiseIndices[0]))

	promiseIndex := mockSys.PromiseAnd(uint64(promiseIdxPtr), uint64(len(promiseIndices)))

	expectedIndex := uint64(2)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}
}

func TestPromiseBatchCreate(t *testing.T) {
	mockSys := NewMockSystem()
	accountId := "accountId"
	accountIdPtr := uintptr(unsafe.Pointer(&accountId))

	promiseIndex := mockSys.PromiseBatchCreate(uint64(len(accountId)), uint64(accountIdPtr))

	expectedIndex := uint64(0)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}
}

func TestPromiseBatchThen(t *testing.T) {
	mockSys := NewMockSystem()
	accountId := "accountId"
	accountIdPtr := uintptr(unsafe.Pointer(&accountId))
	promiseIdx := uint64(0)

	promiseIndex := mockSys.PromiseBatchThen(promiseIdx, uint64(len(accountId)), uint64(accountIdPtr))

	expectedIndex := uint64(1)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}
}

// Promises API

// Promises API Action

func TestPromiseBatchActionCreateAccount(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)

	mockSys.PromiseBatchActionCreateAccount(promiseIdx)

}

func TestPromiseBatchActionDeployContract(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	contractBytes := []byte("sample contract bytes")
	codeLen := uint64(len(contractBytes))
	codePtr := uintptr(unsafe.Pointer(&contractBytes[0]))

	mockSys.PromiseBatchActionDeployContract(promiseIdx, codeLen, uint64(codePtr))

}

func TestPromiseBatchActionFunctionCall(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	functionName := []byte("TestLogStringUtf8")
	functionNameLen := uint64(len(functionName))
	functionNamePtr := uintptr(unsafe.Pointer(&functionName[0]))
	arguments := []byte("{}")
	argumentsLen := uint64(len(arguments))
	argumentsPtr := uintptr(unsafe.Pointer(&arguments[0]))
	amount := types.Uint128{Hi: 0, Lo: 0}
	amountPtr := uintptr(unsafe.Pointer(&amount))
	gas := uint64(3000000000)

	mockSys.PromiseBatchActionFunctionCall(promiseIdx, functionNameLen, uint64(functionNamePtr), argumentsLen, uint64(argumentsPtr), uint64(amountPtr), gas)

}

func TestPromiseBatchActionFunctionCallWeight(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	functionName := []byte("TestLogStringUtf8")
	functionNameLen := uint64(len(functionName))
	functionNamePtr := uintptr(unsafe.Pointer(&functionName[0]))
	arguments := []byte("{}")
	argumentsLen := uint64(len(arguments))
	argumentsPtr := uintptr(unsafe.Pointer(&arguments[0]))
	amount := types.Uint128{Hi: 0, Lo: 0}
	amountPtr := uintptr(unsafe.Pointer(&amount))
	gas := uint64(3000000000)
	weight := uint64(1)

	mockSys.PromiseBatchActionFunctionCallWeight(promiseIdx, functionNameLen, uint64(functionNamePtr), argumentsLen, uint64(argumentsPtr), uint64(amountPtr), gas, weight)
}

func TestPromiseBatchActionTransfer(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	amount, _ := types.U128FromString("10000000000000000000000") // 0.01 Near
	amountPtr := uintptr(unsafe.Pointer(&amount))
	mockSys.PromiseBatchActionTransfer(promiseIdx, uint64(amountPtr))
}

func TestPromiseBatchActionStake(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	amount, _ := types.U128FromString("10000000000000000000000") // 0.01 Near
	amountPtr := uintptr(unsafe.Pointer(&amount))
	signerPK := []byte("sample_signer_public_key")
	publicKeyLen := uint64(len(signerPK))
	publicKeyPtr := uintptr(unsafe.Pointer(&signerPK[0]))

	mockSys.PromiseBatchActionStake(promiseIdx, uint64(amountPtr), publicKeyLen, uint64(publicKeyPtr))
}

func TestPromiseBatchActionAddKeyWithFullAccess(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	publicKey, _ := types.PublicKeyFromString("ed25519:ExeqWPvjcUjLX3NfTk3JzisaXLjsCqJNZFCj7ub92RQW")
	publicKeyLen := uint64(len(publicKey.Bytes()))
	publicKeyPtr := uintptr(unsafe.Pointer(&publicKey.Bytes()[0]))
	nonce := uint64(0)

	mockSys.PromiseBatchActionAddKeyWithFullAccess(promiseIdx, publicKeyLen, uint64(publicKeyPtr), nonce)
}

func TestPromiseBatchActionAddKeyWithFunctionCall(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	publicKey, _ := types.PublicKeyFromString("ed25519:ExeqWPvjcUjLX3NfTk3JzisaXLjsCqJNZFCj7ub92RQW")
	publicKeyLen := uint64(len(publicKey.Bytes()))
	publicKeyPtr := uintptr(unsafe.Pointer(&publicKey.Bytes()[0]))
	nonce := uint64(0)
	amount := types.Uint128{Hi: 0, Lo: 1000}
	amountPtr := uintptr(unsafe.Pointer(&amount))
	receiverId := []byte("receiver.near")
	receiverIdLen := uint64(len(receiverId))
	receiverIdPtr := uintptr(unsafe.Pointer(&receiverId[0]))
	functionName := []byte("TestLogStringUtf8")
	functionNameLen := uint64(len(functionName))
	functionNamePtr := uintptr(unsafe.Pointer(&functionName[0]))

	mockSys.PromiseBatchActionAddKeyWithFunctionCall(promiseIdx, publicKeyLen, uint64(publicKeyPtr), nonce, uint64(amountPtr), receiverIdLen, uint64(receiverIdPtr), functionNameLen, uint64(functionNamePtr))

}

func TestPromiseBatchActionDeleteKey(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	publicKey, _ := types.PublicKeyFromString("ed25519:ExeqWPvjcUjLX3NfTk3JzisaXLjsCqJNZFCj7ub92RQW")
	publicKeyLen := uint64(len(publicKey.Bytes()))
	publicKeyPtr := uintptr(unsafe.Pointer(&publicKey.Bytes()[0]))
	mockSys.PromiseBatchActionDeleteKey(promiseIdx, publicKeyLen, uint64(publicKeyPtr))
}

func TestPromiseBatchActionDeleteAccount(t *testing.T) {
	mockSys := MockSystem{}
	promiseIdx := uint64(0)
	beneficiaryId := []byte("beneficiary.near")
	beneficiaryIdLen := uint64(len(beneficiaryId))
	beneficiaryIdPtr := uintptr(unsafe.Pointer(&beneficiaryId[0]))
	mockSys.PromiseBatchActionDeleteAccount(promiseIdx, beneficiaryIdLen, uint64(beneficiaryIdPtr))
}

func TestPromiseYieldCreate(t *testing.T) {
	mockSys := MockSystem{}
	functionName := []byte("TestContractValueReturn")
	functionNameLen := uint64(len(functionName))
	functionNamePtr := uintptr(unsafe.Pointer(&functionName[0]))
	arguments := []byte("{}")
	argumentsLen := uint64(len(arguments))
	argumentsPtr := uintptr(unsafe.Pointer(&arguments[0]))
	gas := uint64(3000000000)
	gasWeight := uint64(0)
	registerId := uint64(0)

	promiseIdx := mockSys.PromiseYieldCreate(functionNameLen, uint64(functionNamePtr), argumentsLen, uint64(argumentsPtr), gas, gasWeight, registerId)

	expectedIndex := uint64(1)
	if promiseIdx != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIdx)
	}
}

func TestPromiseYieldResume(t *testing.T) {
	mockSys := MockSystem{}
	data := []byte("sample data")
	dataIdLen := uint64(len(data))
	dataIdPtr := uintptr(unsafe.Pointer(&data[0]))
	payload := []byte("sample payload")
	payloadLen := uint64(len(payload))
	payloadPtr := uintptr(unsafe.Pointer(&payload[0]))

	result := mockSys.PromiseYieldResume(dataIdLen, uint64(dataIdPtr), payloadLen, uint64(payloadPtr))

	expectedResult := uint32(1)
	if result != expectedResult {
		t.Errorf("expected result %d, got %d", expectedResult, result)
	}
}

// Promises API Action

// Promise API Results

func TestPromiseResultsCount(t *testing.T) {
	mockSys := MockSystem{}

	count := mockSys.PromiseResultsCount()
	expectedCount := uint64(0)
	if count != expectedCount {
		t.Errorf("expected promise count %d, got %d", expectedCount, count)
	}
}

func TestPromiseResult(t *testing.T) {
	mockSys := MockSystem{}

	resultIdx := uint64(0)
	registerId := uint64(0)

	status := mockSys.PromiseResult(resultIdx, registerId)
	expectedStatus := uint64(0)
	if status != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, status)
	}
}

func TestPromiseReturn(t *testing.T) {
	mockSys := MockSystem{}
	promiseId := uint64(0)

	mockSys.PromiseReturn(promiseId)
}

func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Promise API Results
