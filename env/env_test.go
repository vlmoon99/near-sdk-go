package env

import (
	"encoding/json"
	"testing"

	"github.com/vlmoon99/near-sdk-go/system"
	"github.com/vlmoon99/near-sdk-go/types"
)

func init() {
	SetEnv(system.NewMockSystem())
}

func TestSetEnv(t *testing.T) {
	mockSys := system.NewMockSystem()
	SetEnv(mockSys)

	if NearBlockchainImports != mockSys {
		t.Errorf("expected NearBlockchainImports to be set to mockSys, got %v", NearBlockchainImports)
	}
}

// Registers API

func TestTryMethodIntoRegister(t *testing.T) {
	mockSys := system.NewMockSystem()
	SetEnv(mockSys)

	data := []byte("test data")
	mockSys.Registers[AtomicOpRegister] = data

	method := func(registerId uint64) {
		WriteRegisterSafe(registerId, data)
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
		WriteRegisterSafe(registerId, data)
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

	result, err := ReadRegisterSafe(AtomicOpRegister)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if string(result) != string(data) {
		t.Errorf("expected '%s', got '%s'", data, result)
	}

	result, err = ReadRegisterSafe(1)
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
	WriteRegisterSafe(1, data)

	if string(mockSys.Registers[1]) != string(data) {
		t.Errorf("expected '%s', got '%s'", data, mockSys.Registers[1])
	}

	WriteRegisterSafe(2, []byte{})
	if _, exists := mockSys.Registers[2]; exists {
		t.Errorf("expected register 2 to be empty")
	}
}

// Registers API

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

// Context API
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
	type TestPayload struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
		Key3 bool   `json:"key3"`
	}

	inputData := TestPayload{
		Key1: "value1",
		Key2: 42,
		Key3: true,
	}

	jsonData, err := json.Marshal(inputData)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	mockSys, _ := NearBlockchainImports.(*system.MockSystem)
	mockSys.ContractInput = jsonData
	mockSys.Input(1)

	options := types.ContractInputOptions{IsRawBytes: false}
	data, dataType, err := ContractInput(options)
	if err != nil {
		t.Fatalf("ContractInput failed: %v", err)
	}

	var result TestPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if result.Key1 != "value1" {
		t.Fatalf("Expected value %s, got %s", "value1", result.Key1)
	}

	if result.Key2 != 42 {
		t.Fatalf("Expected value %d, got %d", 42, result.Key2)
	}

	if result.Key3 != true {
		t.Fatalf("Expected value %v, got %v", true, result.Key3)
	}

	expectedType := "object"
	if dataType != expectedType {
		t.Fatalf("Expected data type %s, got %s", expectedType, dataType)
	}
}

// Context API

// Economics API

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

func TestGetAttachedDeposit(t *testing.T) {
	expected := types.Uint128{Hi: 0, Lo: 0}
	deposit, err := GetAttachedDeposit()
	if err != nil {
		t.Fatalf("GetAttachedDeposit failed: %v", err)
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

// Economics API

// Math API

func TestGetRandomSeed(t *testing.T) {
	expected := []byte("randomSeed")
	seed, err := GetRandomSeed()
	if err != nil {
		t.Fatalf("GetRandomSeed failed: %v", err)
	}

	if string(seed) != string(expected) {
		t.Fatalf("Expected %s, got %s", string(expected), string(seed))
	}
}

func TestSha256Hash(t *testing.T) {
	data := []byte("test data")
	expected := []byte("hash")
	hash, err := Sha256Hash(data)
	if err != nil {
		t.Fatalf("Sha256Hash failed: %v", err)
	}

	if string(hash) != string(expected) {
		t.Fatalf("Expected %s, got %s", string(expected), string(hash))
	}
}

func TestKeccak256Hash(t *testing.T) {
	data := []byte("test data")
	expected := []byte("hash")
	hash, err := Keccak256Hash(data)
	if err != nil {
		t.Fatalf("Keccak256Hash failed: %v", err)
	}

	if string(hash) != string(expected) {
		t.Fatalf("Expected %s, got %s", string(expected), string(hash))
	}
}

func TestKeccak512Hash(t *testing.T) {
	data := []byte("test data")
	expected := []byte("hash")
	hash, err := Keccak512Hash(data)
	if err != nil {
		t.Fatalf("Keccak512Hash failed: %v", err)
	}

	if string(hash) != string(expected) {
		t.Fatalf("Expected %s, got %s", string(expected), string(hash))
	}
}

func TestRipemd160Hash(t *testing.T) {
	data := []byte("test data")
	expected := []byte("hash")
	hash, err := Ripemd160Hash(data)
	if err != nil {
		t.Fatalf("Ripemd160Hash failed: %v", err)
	}

	if string(hash) != string(expected) {
		t.Fatalf("Expected %s, got %s", string(expected), string(hash))
	}
}

func TestEcrecoverPubKey(t *testing.T) {
	hash := []byte("test hash")
	signature := []byte("test signature")
	v := byte(1)
	malleabilityFlag := true
	expected := []byte{1, 2, 3, 4}

	mockSys, _ := NearBlockchainImports.(*system.MockSystem)
	mockSys.Registers[AtomicOpRegister] = expected

	pubKey, err := EcrecoverPubKey(hash, signature, v, malleabilityFlag)
	if err != nil {
		t.Fatalf("EcrecoverPubKey failed: %v", err)
	}

	if string(pubKey) != string(expected) {
		t.Fatalf("Expected %s, got %s", string(expected), string(pubKey))
	}
}

func TestEd25519VerifySig(t *testing.T) {
	signature := [64]byte{}
	message := []byte("test message")
	publicKey := [32]byte{}

	result := Ed25519VerifySig(signature, message, publicKey)
	if !result {
		t.Fatalf("Expected true , got false")
	}
}

func TestAltBn128G1MultiExp(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	expected := []byte("simpleMultiexp")
	result, err := AltBn128G1MultiExp(data)
	if err != nil {
		t.Fatalf("AltBn128G1MultiExp failed: %v", err)
	}

	if string(result) != string(expected) {
		t.Fatalf("Expected %s, got %s", string(expected), string(result))
	}
}

func TestAltBn128G1Sum(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	expected := []byte("simpleSum")
	result, err := AltBn128G1Sum(data)
	if err != nil {
		t.Fatalf("AltBn128G1Sum failed: %v", err)
	}

	if string(result) != string(expected) {
		t.Fatalf("Expected %s, got %s", string(expected), string(result))
	}
}

func TestAltBn128PairingCheck(t *testing.T) {
	data := []byte("test data")
	result := AltBn128PairingCheck(data)
	if !result {
		t.Fatalf("Expected true , got false")
	}
}

// Math API

// Validator API

func TestValidatorStakeAmount(t *testing.T) {
	accountID := []byte("validatorAccountId")
	expectedStake := types.Uint128{Hi: 0, Lo: 100000}

	stakeAmount, err := ValidatorStakeAmount(accountID)
	if err != nil {
		t.Fatalf("ValidatorStakeAmount failed: %v", err)
	}

	if stakeAmount != expectedStake {
		t.Fatalf("expected stake %v, got %v", expectedStake, stakeAmount)
	}
}

func TestValidatorStakeAmount_EmptyAccountID(t *testing.T) {
	accountID := []byte("")

	_, err := ValidatorStakeAmount(accountID)
	if err == nil || err.Error() != ErrAccountIDMustNotBeEmpty {
		t.Fatalf("expected error %v, got %v", ErrAccountIDMustNotBeEmpty, err)
	}
}

func TestValidatorTotalStakeAmount(t *testing.T) {
	expectedTotalStake := types.Uint128{Hi: 0, Lo: 100000}

	totalStakeAmount, err := ValidatorTotalStakeAmount()
	if err != nil {
		t.Fatalf("ValidatorTotalStakeAmount failed: %v", err)
	}

	if totalStakeAmount != expectedTotalStake {
		t.Fatalf("expected total stake %v, got %v", expectedTotalStake, totalStakeAmount)
	}
}

// Validator API

// Miscellaneous API

func TestContractValueReturn(t *testing.T) {
	input := []byte("test value")
	mockSys, _ := NearBlockchainImports.(*system.MockSystem)
	ContractValueReturn(input)
	if string(mockSys.Registers[0]) != string(input) {
		t.Fatalf("expected %s, got %s", string(input), string(mockSys.Registers[0]))
	}
}

// Miscellaneous API

// Promises API
func TestPromiseCreate(t *testing.T) {
	accountId := []byte("accountId")
	functionName := []byte("functionName")
	arguments := []byte("arguments")
	amount := types.Uint128{Lo: 0, Hi: 0}
	gas := uint64(5000)

	promiseIndex := PromiseCreate(accountId, functionName, arguments, amount, gas)

	expectedIndex := uint64(0)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}

	mockSys, ok := NearBlockchainImports.(*system.MockSystem)
	if !ok {
		t.Fatalf("Failed to cast NearBlockchainImports to *system.MockSystem")
	}

	if len(mockSys.Promises) != 1 {
		t.Errorf("expected 1 promise, got %d", len(mockSys.Promises))
	}

	promise := mockSys.Promises[0]
	if string(promise.AccountId) != string(accountId) {
		t.Errorf("expected account id %s, got %s", string(accountId), promise.AccountId)
	}
	if string(promise.FunctionName) != string(functionName) {
		t.Errorf("expected function name %s, got %s", string(functionName), promise.FunctionName)
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

func TestPromiseThen(t *testing.T) {
	accountId := []byte("accountId")
	functionName := []byte("functionName")
	arguments := []byte("arguments")
	amount := types.Uint128{Lo: 0, Hi: 0}
	gas := uint64(5000)

	// Create the first promise
	PromiseCreate(accountId, functionName, arguments, amount, gas)

	promiseIdx := uint64(0)
	promiseIndex := PromiseThen(promiseIdx, accountId, functionName, arguments, amount, gas)

	expectedIndex := uint64(2)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}

	mockSys, ok := NearBlockchainImports.(*system.MockSystem)
	if !ok {
		t.Fatalf("Failed to cast NearBlockchainImports to *system.MockSystem")
	}

	if len(mockSys.Promises) != 3 {
		t.Errorf("expected 3 promises, got %d", len(mockSys.Promises))
	}

	promise := mockSys.Promises[1]
	if string(promise.AccountId) != string(accountId) {
		t.Errorf("expected account id %s, got %s", string(accountId), promise.AccountId)
	}
	if string(promise.FunctionName) != string(functionName) {
		t.Errorf("expected function name %s, got %s", string(functionName), promise.FunctionName)
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
	promiseIndices := []uint64{0, 1}
	promiseIndex := PromiseAnd(promiseIndices)

	expectedIndex := uint64(2)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}
}

func TestPromiseBatchCreate(t *testing.T) {
	accountId := []byte("accountId")
	promiseIndex := PromiseBatchCreate(accountId)

	expectedIndex := uint64(0)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}
}

func TestPromiseBatchThen(t *testing.T) {
	accountId := []byte("accountId")
	promiseIdx := uint64(0)
	promiseIndex := PromiseBatchThen(promiseIdx, accountId)

	expectedIndex := uint64(1)
	if promiseIndex != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIndex)
	}
}

// Promises API

// Promises API Action

func TestPromiseBatchActionCreateAccount(t *testing.T) {
	promiseIdx := uint64(0)
	PromiseBatchActionCreateAccount(promiseIdx)
	// Verify actions
}

func TestPromiseBatchActionDeployContract(t *testing.T) {
	promiseIdx := uint64(0)
	contractBytes := []byte("sample contract bytes")
	PromiseBatchActionDeployContract(promiseIdx, contractBytes)
	// Verify actions
}

func TestPromiseBatchActionFunctionCall(t *testing.T) {
	promiseIdx := uint64(0)
	functionName := []byte("TestLogStringUtf8")
	arguments := []byte("{}")
	amount := types.Uint128{Hi: 0, Lo: 0}
	gas := uint64(3000000000)
	PromiseBatchActionFunctionCall(promiseIdx, functionName, arguments, amount, gas)
	// Verify actions
}

func TestPromiseBatchActionFunctionCallWeight(t *testing.T) {
	promiseIdx := uint64(0)
	functionName := []byte("TestLogStringUtf8")
	arguments := []byte("{}")
	amount := types.Uint128{Hi: 0, Lo: 0}
	gas := uint64(3000000000)
	weight := uint64(1)
	PromiseBatchActionFunctionCallWeight(promiseIdx, functionName, arguments, amount, gas, weight)
	// Verify actions
}

func TestPromiseBatchActionTransfer(t *testing.T) {
	promiseIdx := uint64(0)
	amount := types.Uint128{Hi: 0, Lo: 1000}
	PromiseBatchActionTransfer(promiseIdx, amount)
	// Verify actions
}

func TestPromiseBatchActionStake(t *testing.T) {
	promiseIdx := uint64(0)
	amount := types.Uint128{Hi: 0, Lo: 1000}
	publicKey := []byte("sample_public_key")
	PromiseBatchActionStake(promiseIdx, amount, publicKey)
	// Verify actions
}

func TestPromiseBatchActionAddKeyWithFullAccess(t *testing.T) {
	promiseIdx := uint64(0)
	publicKey := []byte("sample_public_key")
	nonce := uint64(0)
	PromiseBatchActionAddKeyWithFullAccess(promiseIdx, publicKey, nonce)
	// Verify actions
}

func TestPromiseBatchActionAddKeyWithFunctionCall(t *testing.T) {
	promiseIdx := uint64(0)
	publicKey := []byte("sample_public_key")
	nonce := uint64(0)
	amount := types.Uint128{Hi: 0, Lo: 1000}
	receiverId := []byte("receiver.near")
	functionName := []byte("TestLogStringUtf8")
	PromiseBatchActionAddKeyWithFunctionCall(promiseIdx, publicKey, nonce, amount, receiverId, functionName)
	// Verify actions
}

func TestPromiseBatchActionDeleteKey(t *testing.T) {
	promiseIdx := uint64(0)
	publicKey := []byte("sample_public_key")
	PromiseBatchActionDeleteKey(promiseIdx, publicKey)
	// Verify actions
}

func TestPromiseBatchActionDeleteAccount(t *testing.T) {
	promiseIdx := uint64(0)
	beneficiaryId := []byte("beneficiary.near")
	PromiseBatchActionDeleteAccount(promiseIdx, beneficiaryId)
	// Verify actions
}

func TestPromiseYieldCreate(t *testing.T) {
	functionName := []byte("TestContractValueReturn")
	arguments := []byte("{}")
	gas := uint64(3000000000)
	gasWeight := uint64(0)
	promiseIdx := PromiseYieldCreate(functionName, arguments, gas, gasWeight)
	expectedIndex := uint64(1)
	if promiseIdx != expectedIndex {
		t.Errorf("expected promise index %d, got %d", expectedIndex, promiseIdx)
	}
}

func TestPromiseYieldResume(t *testing.T) {
	data := []byte("sample data")
	payload := []byte("sample payload")
	result := PromiseYieldResume(data, payload)
	expectedResult := uint32(1)
	if result != expectedResult {
		t.Errorf("expected result %d, got %d", expectedResult, result)
	}
}

// Promises API Action

// Promise API Results
func TestPromiseResultsCount(t *testing.T) {
	count := PromiseResultsCount()
	expectedCount := uint64(3)
	if count != expectedCount {
		t.Errorf("expected promise count %d, got %d", expectedCount, count)
	}
}

func TestPromiseResult(t *testing.T) {
	resultIdx := uint64(0)
	_, err := PromiseResult(resultIdx)
	if err != nil {
		t.Fatalf("PromiseResult failed: %v", err)
	}

}

func TestPromiseReturn(t *testing.T) {
	promiseId := uint64(0)
	PromiseReturn(promiseId)
}

// Promise API Results
