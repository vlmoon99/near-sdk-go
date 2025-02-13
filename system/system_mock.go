package system

// For some env limitation reason we can't use crypto/* or golang.org/x/crypto/* packages
import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
	"unsafe"

	"github.com/vlmoon99/near-sdk-go/types"
)

type MockPromise struct {
	AccountId    string
	FunctionName string
	Arguments    []byte
	Amount       types.Uint128
	Gas          uint64
	PromiseIndex uint64
}

// Test Mock impl of the System interface
type MockSystem struct {
	Promises                []MockPromise
	Registers               map[uint64][]byte
	Storage                 map[string][]byte
	CurrentAccountIdSys     string
	SignerAccountIdSys      string
	SignerAccountPkSys      []byte
	PredecessorAccountIdSys string
	ContractInput           []byte
	BlockIndexSys           uint64
	BlockTimestampSys       uint64
	EpochHeightSys          uint64
	StorageUsageSys         uint64
	AccountBalanceSys       types.Uint128
	AccountLockedBalanceSys types.Uint128
	AttachedDepositSys      types.Uint128
	PrepaidGasSys           uint64
	UsedGasSys              uint64
}

func NewMockSystem() *MockSystem {
	return &MockSystem{
		Registers:               make(map[uint64][]byte),
		Storage:                 make(map[string][]byte),
		CurrentAccountIdSys:     "currentAccountId.near",
		SignerAccountIdSys:      "signerAccountId.near",
		SignerAccountPkSys:      []byte("signerAccountPk"),
		PredecessorAccountIdSys: "predecessorAccountId.near",
		ContractInput:           []byte("Test Input"),
		BlockIndexSys:           1,
		BlockTimestampSys:       uint64(1739394085901002712),
		EpochHeightSys:          1,
		StorageUsageSys:         0,
		AccountBalanceSys:       types.Uint128{Hi: 0, Lo: 0},
		AccountLockedBalanceSys: types.Uint128{Hi: 0, Lo: 0},
		AttachedDepositSys:      types.Uint128{Hi: 0, Lo: 0},
		PrepaidGasSys:           5000,
		UsedGasSys:              2500,
	}
}

// Registers API

func (m *MockSystem) WriteRegister(registerId, dataLen, dataPtr uint64) {
	dataSlice := make([]byte, dataLen)
	copy(dataSlice, unsafe.Slice((*byte)(unsafe.Pointer(uintptr(dataPtr))), dataLen))

	m.Registers[registerId] = dataSlice
}

func (m *MockSystem) ReadRegister(registerId, ptr uint64) {
	if data, exists := m.Registers[registerId]; exists {
		copy(unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), len(data)), data)
	}
}

func (m *MockSystem) RegisterLen(registerId uint64) uint64 {
	if data, exists := m.Registers[registerId]; exists {
		return uint64(len(data))
	}
	return 0
}

// Registers API

// Storage API
func (m *MockSystem) StorageWrite(keyLen, keyPtr, valueLen, valuePtr, registerId uint64) uint64 {
	key := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(keyPtr))), keyLen)
	value := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(valuePtr))), valueLen)
	keyStr := string(key)

	m.Storage[keyStr] = value
	return 1
}

func (m *MockSystem) StorageRead(keyLen, keyPtr, registerId uint64) uint64 {
	key := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(keyPtr))), keyLen)
	keyStr := string(key)

	if value, exists := m.Storage[keyStr]; exists {
		m.WriteRegister(registerId, uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))))
		return 1
	}
	return 0
}

func (m *MockSystem) StorageRemove(keyLen, keyPtr, registerId uint64) uint64 {
	key := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(keyPtr))), keyLen)
	keyStr := string(key)

	if value, exists := m.Storage[keyStr]; exists {
		delete(m.Storage, keyStr)
		m.WriteRegister(registerId, uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))))
		return 1
	}
	return 0
}

func (m *MockSystem) StorageHasKey(keyLen, keyPtr uint64) uint64 {
	key := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(keyPtr))), keyLen)
	keyStr := string(key)

	if _, exists := m.Storage[keyStr]; exists {
		return 1
	}
	return 0
}

// Storage API

// Context API
func (m *MockSystem) CurrentAccountId(registerId uint64) {
	data := []byte(m.CurrentAccountIdSys)
	m.WriteRegister(registerId, uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))))
}

func (m *MockSystem) SignerAccountId(registerId uint64) {
	data := []byte(m.SignerAccountIdSys)
	m.WriteRegister(registerId, uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))))
}

func (m *MockSystem) SignerAccountPk(registerId uint64) {
	m.WriteRegister(registerId, uint64(len(m.SignerAccountPkSys)), uint64(uintptr(unsafe.Pointer(&m.SignerAccountPkSys[0]))))
}

func (m *MockSystem) PredecessorAccountId(registerId uint64) {
	data := []byte(m.PredecessorAccountIdSys)
	m.WriteRegister(registerId, uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))))
}

func (m *MockSystem) Input(registerId uint64) {
	m.WriteRegister(registerId, uint64(len(m.ContractInput)), uint64(uintptr(unsafe.Pointer(&m.ContractInput[0]))))
}

func (m *MockSystem) BlockIndex() uint64 {
	return m.BlockIndexSys
}

func (m *MockSystem) BlockTimestamp() uint64 {
	return m.BlockTimestampSys
}

func (m *MockSystem) EpochHeight() uint64 {
	return m.EpochHeightSys
}

func (m *MockSystem) StorageUsage() uint64 {
	return m.StorageUsageSys
}

// Context API

// Economics API
func (m *MockSystem) AccountBalance(balancePtr uint64) {
	balance := *(*[]byte)(unsafe.Pointer(&balancePtr))
	balance = m.AccountBalanceSys.ToBE()
	fmt.Printf("balance: %v\n", balance)
}

func (m *MockSystem) AccountLockedBalance(balancePtr uint64) {
	balance := *(*[]byte)(unsafe.Pointer(&balancePtr))
	balance = m.AccountLockedBalanceSys.ToBE()
	fmt.Printf("balance: %v\n", balance)
}

func (m *MockSystem) AttachedDeposit(balancePtr uint64) {
	balance := *(*[]byte)(unsafe.Pointer(&balancePtr))
	balance = m.AttachedDepositSys.ToBE()
	fmt.Printf("balance: %v\n", balance)
}

func (m *MockSystem) PrepaidGas() uint64 {
	return m.PrepaidGasSys
}

func (m *MockSystem) UsedGas() uint64 {
	return m.UsedGasSys
}

// Math API

func (m *MockSystem) RandomSeed(registerId uint64) {
	seed := []byte("randomSeed")
	m.WriteRegister(registerId, uint64(len(seed)), uint64(uintptr(unsafe.Pointer(&seed[0]))))
}

func (m *MockSystem) Sha256(valueLen, valuePtr, registerId uint64) {
	hash := []byte("hash")
	m.WriteRegister(registerId, uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))))
}

func (m *MockSystem) Keccak256(valueLen, valuePtr, registerId uint64) {
	hash := []byte("hash")
	m.WriteRegister(registerId, uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))))
}

func (m *MockSystem) Keccak512(valueLen, valuePtr, registerId uint64) {
	hash := []byte("hash")
	m.WriteRegister(registerId, uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))))
}

func (m *MockSystem) Ripemd160(valueLen, valuePtr, registerId uint64) {
	hash := []byte("hash")
	m.WriteRegister(registerId, uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))))
}

func (m *MockSystem) Ecrecover(hashLen, hashPtr, sigLen, sigPtr, v, malleabilityFlag, registerId uint64) uint64 {
	return 1
}

func (m *MockSystem) Ed25519Verify(sigLen, sigPtr, msgLen, msgPtr, pubKeyLen, pubKeyPtr uint64) uint64 {
	return 1
}

func (m *MockSystem) AltBn128G1Multiexp(valueLen, valuePtr, registerId uint64) {
	simpleMultiexp := []byte("simpleMultiexp")
	m.WriteRegister(registerId, uint64(len(simpleMultiexp)), uint64(uintptr(unsafe.Pointer(&simpleMultiexp[0]))))
}

func (m *MockSystem) AltBn128G1SumSystem(valueLen, valuePtr, registerId uint64) {
	simpleSum := []byte("simpleSum")
	m.WriteRegister(registerId, uint64(len(simpleSum)), uint64(uintptr(unsafe.Pointer(&simpleSum[0]))))
}

func (m *MockSystem) AltBn128PairingCheckSystem(valueLen, valuePtr uint64) uint64 {
	return 1
}

// Math API

// Validator API

func (m *MockSystem) ValidatorStake(accountIdLen, accountIdPtr, stakePtr uint64) {
	accountId := make([]byte, accountIdLen)
	copy(accountId, *(*[]byte)(unsafe.Pointer(uintptr(accountIdPtr))))

	stakeAmount := uint64(0)
	if string(accountId) == "validatorAccountId" {
		stakeAmount = 100000
	}

	binary.LittleEndian.PutUint64((*(*[8]byte)(unsafe.Pointer(uintptr(stakePtr))))[:], stakeAmount)
	fmt.Printf("stake: %v\n", stakeAmount)
}

func (m *MockSystem) ValidatorTotalStake(stakePtr uint64) {
	totalStake := uint64(1000000)

	binary.LittleEndian.PutUint64((*(*[8]byte)(unsafe.Pointer(uintptr(stakePtr))))[:], totalStake)
	fmt.Printf("stake: %v\n", totalStake)
}

// Validator API

// Miscellaneous API

func (m *MockSystem) ValueReturn(valueLen, valuePtr uint64) {
	m.WriteRegister(0, valueLen, valuePtr)
}

func (m *MockSystem) PanicUtf8(len, ptr uint64) {
	value := make([]byte, len)
	copy(value, *(*[]byte)(unsafe.Pointer(uintptr(ptr))))

	fmt.Printf("Panic: %s", value[:len])
}

func (m *MockSystem) LogUtf8(len, ptr uint64) {
	value := make([]byte, len)
	copy(value, *(*[]byte)(unsafe.Pointer(uintptr(ptr))))
	fmt.Printf("Log: %s", value[:len])
}

func (m *MockSystem) LogUtf16(len, ptr uint64) {
	utf16Bytes := make([]uint16, len/2)
	for i := 0; i < int(len)/2; i++ {
		utf16Bytes[i] = *(*uint16)(unsafe.Pointer(uintptr(ptr) + uintptr(i*2)))
	}
	message := string(utf16.Decode(utf16Bytes))
	fmt.Printf("Log: %s", message)
}

// Miscellaneous API

// Promises API
func (m *MockSystem) PromiseCreate(accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64 {
	accountId := "accountId"
	functionName := "functionName"
	arguments := []byte("arguments")
	amount := types.Uint128{Lo: 0, Hi: 0}

	promise := MockPromise{
		AccountId:    accountId[:accountIdLen],
		FunctionName: functionName[:functionNameLen],
		Arguments:    arguments[:argumentsLen],
		Amount:       amount,
		Gas:          gas,
		PromiseIndex: uint64(len(m.Promises)),
	}

	m.Promises = append(m.Promises, promise)

	return promise.PromiseIndex
}

func (m *MockSystem) PromiseThen(promiseIndex, accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64 {
	accountId := "accountId"
	functionName := "functionName"
	arguments := []byte("arguments")
	amount := types.Uint128{Lo: 0, Hi: 0}

	promise := MockPromise{
		AccountId:    accountId,
		FunctionName: functionName,
		Arguments:    arguments,
		Amount:       amount,
		Gas:          gas,
		PromiseIndex: uint64(len(m.Promises)),
	}

	m.Promises = append(m.Promises, promise)

	return promise.PromiseIndex
}

func (m *MockSystem) PromiseAnd(promiseIdxPtr, promiseIdxCount uint64) uint64 {
	return uint64(2)
}

func (m *MockSystem) PromiseBatchCreate(accountIdLen, accountIdPtr uint64) uint64 {
	return 0
}

func (m *MockSystem) PromiseBatchThen(promiseIndex, accountIdLen, accountIdPtr uint64) uint64 {
	return 1
}

// Promises API

// Promise API Actions
func (m *MockSystem) PromiseBatchActionAddKeyWithFullAccess(promiseIndex, publicKeyLen, publicKeyPtr, nonce uint64) {
	publicKey := *(*string)(unsafe.Pointer(&publicKeyPtr))
	promise := &m.Promises[promiseIndex]
	promise.Arguments = append(promise.Arguments, publicKey[:publicKeyLen]...)
	promise.Arguments = append(promise.Arguments, byte(nonce))
}

func (m *MockSystem) PromiseBatchActionAddKeyWithFunctionCall(promiseIndex, publicKeyLen, publicKeyPtr, nonce, allowancePtr, receiverIdLen, receiverIdPtr, functionNamesLen, functionNamesPtr uint64) {
	publicKey := *(*string)(unsafe.Pointer(&publicKeyPtr))
	allowance := *(*uint64)(unsafe.Pointer(&allowancePtr))
	receiverId := *(*string)(unsafe.Pointer(&receiverIdPtr))
	functionNames := *(*string)(unsafe.Pointer(&functionNamesPtr))

	promise := &m.Promises[promiseIndex]
	promise.Arguments = append(promise.Arguments, publicKey[:publicKeyLen]...)
	promise.Arguments = append(promise.Arguments, byte(nonce))
	promise.Arguments = append(promise.Arguments, byte(allowance))
	promise.Arguments = append(promise.Arguments, receiverId[:receiverIdLen]...)
	promise.Arguments = append(promise.Arguments, functionNames[:functionNamesLen]...)
}

func (m *MockSystem) PromiseBatchActionDeleteKey(promiseIndex, publicKeyLen, publicKeyPtr uint64) {
	publicKey := *(*string)(unsafe.Pointer(&publicKeyPtr))
	promise := &m.Promises[promiseIndex]
	promise.Arguments = append(promise.Arguments, publicKey[:publicKeyLen]...)
}

func (m *MockSystem) PromiseBatchActionDeleteAccount(promiseIndex, beneficiaryIdLen, beneficiaryIdPtr uint64) {
	beneficiaryId := *(*string)(unsafe.Pointer(&beneficiaryIdPtr))
	promise := &m.Promises[promiseIndex]
	promise.Arguments = append(promise.Arguments, beneficiaryId[:beneficiaryIdLen]...)
}

func (m *MockSystem) PromiseYieldCreate(functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, gas, gasWeight, registerId uint64) uint64 {
	functionName := *(*string)(unsafe.Pointer(&functionNamePtr))
	arguments := *(*[]byte)(unsafe.Pointer(&argumentsPtr))

	promise := MockPromise{
		FunctionName: functionName[:functionNameLen],
		Arguments:    arguments[:argumentsLen],
		Gas:          gas,
	}

	m.Promises = append(m.Promises, promise)
	m.WriteRegister(registerId, uint64(len(m.Promises)), uint64(uintptr(unsafe.Pointer(&promise.PromiseIndex))))

	return promise.PromiseIndex
}

func (m *MockSystem) PromiseYieldResume(dataIdLen, dataIdPtr, payloadLen, payloadPtr uint64) uint32 {
	dataId := *(*string)(unsafe.Pointer(&dataIdPtr))
	payload := *(*[]byte)(unsafe.Pointer(&payloadPtr))

	for _, promise := range m.Promises {
		if string(promise.Arguments) == dataId[:dataIdLen] {
			promise.Arguments = append(promise.Arguments, payload[:payloadLen]...)
			return 1
		}
	}
	return 0
}

func (m *MockSystem) PromiseBatchActionCreateAccount(promiseIndex uint64) {
	// Simulate creating an account in a promise batch
}

func (m *MockSystem) PromiseBatchActionDeployContract(promiseIndex, codeLen, codePtr uint64) {
	code := *(*[]byte)(unsafe.Pointer(&codePtr))
	promise := &m.Promises[promiseIndex]
	promise.Arguments = append(promise.Arguments, code[:codeLen]...)
}

func (m *MockSystem) PromiseBatchActionFunctionCall(promiseIndex, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) {
	functionName := *(*string)(unsafe.Pointer(&functionNamePtr))
	arguments := *(*[]byte)(unsafe.Pointer(&argumentsPtr))
	amountBytes := *(*[16]byte)(unsafe.Pointer(&amountPtr))
	amount, _ := types.LoadUint128BE(amountBytes[:])

	promise := &m.Promises[promiseIndex]
	promise.FunctionName = functionName[:functionNameLen]
	promise.Arguments = append(promise.Arguments, arguments[:argumentsLen]...)
	promise.Amount = amount
	promise.Gas = gas
}

func (m *MockSystem) PromiseBatchActionFunctionCallWeight(promiseIndex, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas, weight uint64) {
	// Assuming weight is used somehow in the call, though not specified how exactly
	m.PromiseBatchActionFunctionCall(promiseIndex, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas)
}

func (m *MockSystem) PromiseBatchActionTransfer(promiseIndex, amountPtr uint64) {
	amountBytes := *(*[16]byte)(unsafe.Pointer(&amountPtr))
	amount, _ := types.LoadUint128BE(amountBytes[:])

	promise := &m.Promises[promiseIndex]
	promise.Amount = amount
}

func (m *MockSystem) PromiseBatchActionStake(promiseIndex, amountPtr, publicKeyLen, publicKeyPtr uint64) {
	amountBytes := *(*[16]byte)(unsafe.Pointer(&amountPtr))
	amount, _ := types.LoadUint128BE(amountBytes[:])
	publicKey := *(*string)(unsafe.Pointer(&publicKeyPtr))

	promise := &m.Promises[promiseIndex]
	promise.Amount = amount
	promise.Arguments = append(promise.Arguments, publicKey[:publicKeyLen]...)
}

// Promise API Actions

// Promise API Results
func (m *MockSystem) PromiseResultsCount() uint64 {
	return uint64(len(m.Promises))
}

func (m *MockSystem) PromiseResult(resultIdx uint64, registerId uint64) uint64 {
	if resultIdx < uint64(len(m.Promises)) {
		result := m.Promises[resultIdx].Arguments
		m.WriteRegister(registerId, uint64(len(result)), uint64(uintptr(unsafe.Pointer(&result[0]))))
		return 1 // Return a success indicator
	}
	return 0 // Return a failure indicator
}

func (m *MockSystem) PromiseReturn(promiseId uint64) {
	if promiseId < uint64(len(m.Promises)) {
		m.Registers[0] = m.Promises[promiseId].Arguments
	}
}

// Promise API Results
