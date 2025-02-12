package env

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/system"
	"github.com/vlmoon99/near-sdk-go/types"
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

// Registers

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

// Registers

// Storage API
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

// Storage API

func TestGetCurrentAccountId(t *testing.T) {
	accountId, err := GetCurrentAccountId()
	if err != nil {
		t.Fatalf("GetCurrentAccountId failed: %v", err)
	}

	expected := "currentAccountId.near"
	if accountId != expected {
		t.Fatalf("Expected account ID %s, got %s", expected, accountId)
	}
}

func TestGetSignerAccountID(t *testing.T) {
	accountId, err := GetSignerAccountID()
	if err != nil {
		t.Fatalf("GetSignerAccountID failed: %v", err)
	}

	expected := "signerAccountId.near"
	if accountId != expected {
		t.Fatalf("Expected account ID %s, got %s", expected, accountId)
	}
}

func TestGetSignerAccountPK(t *testing.T) {
	accountPk, err := GetSignerAccountPK()
	if err != nil {
		t.Fatalf("GetSignerAccountPK failed: %v", err)
	}

	expected := "signerAccountPk"
	if string(accountPk) != expected {
		t.Fatalf("Expected account PK %s, got %s", expected, string(accountPk))
	}
}

func TestGetPredecessorAccountID(t *testing.T) {
	accountId, err := GetPredecessorAccountID()
	if err != nil {
		t.Fatalf("GetPredecessorAccountID failed: %v", err)
	}

	expected := "predecessorAccountId.near"
	if accountId != expected {
		t.Fatalf("Expected account ID %s, got %s", expected, accountId)
	}
}

func TestGetCurrentBlockHeight(t *testing.T) {
	blockHeight := GetCurrentBlockHeight()
	expected := system.NewMockSystem().BlockTimestamp()

	if blockHeight != expected {
		t.Fatalf("Expected block height %d, got %d", expected, blockHeight)
	}
}

func TestGetBlockTimeMs(t *testing.T) {
	blockTimeMs := GetBlockTimeMs()
	expected := uint64(system.NewMockSystem().BlockTimestamp() / 1_000_000)

	if blockTimeMs != expected {
		t.Fatalf("Expected block time in ms %d, got %d", expected, blockTimeMs)
	}
}

func TestGetEpochHeight(t *testing.T) {
	epochHeight := GetEpochHeight()
	expected := uint64(system.NewMockSystem().EpochHeight())

	if epochHeight != expected {
		t.Fatalf("Expected epoch height %d, got %d", expected, epochHeight)
	}
}

func TestGetStorageUsage(t *testing.T) {
	storageUsage := GetStorageUsage()
	expected := uint64(system.NewMockSystem().StorageUsage())

	if storageUsage != expected {
		t.Fatalf("Expected storage usage %d, got %d", expected, storageUsage)
	}
}

func TestContractInputRawBytes(t *testing.T) {
	options := types.ContractInputOptions{IsRawBytes: true}
	data, dataType, err := ContractInput(options)
	if err != nil {
		t.Fatalf("ContractInput failed: %v", err)
	}

	expectedData := []byte("Test Input")
	expectedType := "rawBytes"

	if string(data) != string(expectedData) {
		t.Fatalf("Expected data %s, got %s", string(expectedData), string(data))
	}

	if dataType != expectedType {
		t.Fatalf("Expected data type %s, got %s", expectedType, dataType)
	}
}

func TestContractInputJSON(t *testing.T) {
	// Build JSON data
	builder := json.NewBuilder()
	jsonData := builder.AddString("key1", "value1").
		AddInt("key2", 42).
		AddBool("key3", true).
		Build()

	mockSys, _ := nearBlockchainImports.(*system.MockSystem)
	mockSys.ContractInput = jsonData
	mockSys.Input(1)

	// Call ContractInput with IsRawBytes set to false
	options := types.ContractInputOptions{IsRawBytes: false}
	data, dataType, err := ContractInput(options)
	if err != nil {
		t.Fatalf("ContractInput failed: %v", err)
	}

	// Parse the JSON input and verify the values
	parser := json.NewParser(data)

	// Verify "key1"
	value1, err := parser.GetString("key1")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	expectedValue1 := "value1"
	if value1 != expectedValue1 {
		t.Fatalf("Expected value %s, got %s", expectedValue1, value1)
	}

	// Verify "key2"
	value2, err := parser.GetInt("key2")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	expectedValue2 := int64(42)
	if value2 != expectedValue2 {
		t.Fatalf("Expected value %d, got %d", expectedValue2, value2)
	}

	// Verify "key3"
	value3, err := parser.GetBoolean("key3")
	if err != nil {
		t.Fatalf("GetBoolean failed: %v", err)
	}
	expectedValue3 := true
	if value3 != expectedValue3 {
		t.Fatalf("Expected value %v, got %v", expectedValue3, value3)
	}

	// Verify detected type
	expectedType := "object"
	if dataType != expectedType {
		t.Fatalf("Expected data type %s, got %s", expectedType, dataType)
	}
}

func TestGetAccountBalance(t *testing.T) {
	expected := types.Uint128{Hi: 0, Lo: 0}
	balance, err := GetAccountBalance()
	if err != nil {
		t.Fatalf("GetAccountBalance failed: %v", err)
	}

	if balance != expected {
		t.Fatalf("Expected balance %v, got %v", expected, balance)
	}
}

func TestGetAccountLockedBalance(t *testing.T) {
	expected := types.Uint128{Hi: 0, Lo: 0}
	balance, err := GetAccountLockedBalance()
	if err != nil {
		t.Fatalf("GetAccountLockedBalance failed: %v", err)
	}

	if balance != expected {
		t.Fatalf("Expected balance %v, got %v", expected, balance)
	}
}

func TestGetAttachedDepoist(t *testing.T) {
	expected := types.Uint128{Hi: 0, Lo: 0}
	deposit, err := GetAttachedDepoist()
	if err != nil {
		t.Fatalf("GetAttachedDepoist failed: %v", err)
	}

	if deposit != expected {
		t.Fatalf("Expected deposit %v, got %v", expected, deposit)
	}
}

func TestGetPrepaidGas(t *testing.T) {
	expected := types.NearGas{Inner: 5000}
	prepaidGas := GetPrepaidGas()

	if prepaidGas != expected {
		t.Fatalf("Expected prepaid gas %v, got %v", expected, prepaidGas)
	}
}

func TestGetUsedGas(t *testing.T) {
	expected := types.NearGas{Inner: 2500}
	usedGas := GetUsedGas()

	if usedGas != expected {
		t.Fatalf("Expected used gas %v, got %v", expected, usedGas)
	}
}
