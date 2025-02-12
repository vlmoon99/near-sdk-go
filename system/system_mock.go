package system

// For some env limitation reason we can't use crypto/* or golang.org/x/crypto/* packages
import (
	"fmt"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/vlmoon99/near-sdk-go/types"
)

type MockPromise struct {
	AccountId    string
	FunctionName string
	Arguments    []byte
	Amount       uint64
	Gas          uint64
	PromiseIndex uint64
}

// Test Mock impl of the System interface
// TODO : improve it + tests this system , make it very simmilar to the Near Blockchain behaviour
type MockSystem struct {
	Promises             []MockPromise
	Registers            map[uint64][]byte
	Storage              map[string][]byte
	currentAccountId     string
	signerAccountId      string
	signerAccountPk      []byte
	predecessorAccountId string
	input                []byte
	blockIndex           uint64
	blockTimestamp       uint64
	epochHeight          uint64
	storageUsage         uint64
	accountBalance       types.Uint128
	accountLockedBalance types.Uint128
	attachedDeposit      types.Uint128
	prepaidGas           uint64
	usedGas              uint64
}

func NewMockSystem() *MockSystem {
	return &MockSystem{
		Registers:            make(map[uint64][]byte),
		Storage:              make(map[string][]byte),
		currentAccountId:     "currentAccountId.near",
		signerAccountId:      "signerAccountId.near",
		signerAccountPk:      []byte("signerAccountPk"),
		predecessorAccountId: "predecessorAccountId.near",
		input:                []byte("Test Input"),
		blockIndex:           1,
		blockTimestamp:       uint64(time.Now().UnixNano()),
		epochHeight:          1,
		storageUsage:         0,
		accountBalance:       types.Uint128{Hi: 0, Lo: 0},
		accountLockedBalance: types.Uint128{Hi: 0, Lo: 0},
		attachedDeposit:      types.Uint128{Hi: 0, Lo: 0},
		prepaidGas:           5000,
		usedGas:              2500,
	}
}

// Work with env tests and with my impl env methods, but not work with system_mock_tests.go
// Registers API

func (m *MockSystem) WriteRegister(registerId, dataLen, dataPtr uint64) {
	dataSlice := make([]byte, dataLen)
	copy(dataSlice, unsafe.Slice((*byte)(unsafe.Pointer(uintptr(dataPtr))), dataLen)) // ✅ Safe conversion

	m.Registers[registerId] = dataSlice
}

func (m *MockSystem) ReadRegister(registerId, ptr uint64) {
	if data, exists := m.Registers[registerId]; exists {
		copy(unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), len(data)), data) // ✅ Safe
	}
}

func (m *MockSystem) RegisterLen(registerId uint64) uint64 {
	if data, exists := m.Registers[registerId]; exists {
		return uint64(len(data))
	}
	return 0
}

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
	m.WriteRegister(registerId, uint64(len(m.currentAccountId)), uint64(uintptr(unsafe.Pointer(&m.currentAccountId))))
}

func (m *MockSystem) SignerAccountId(registerId uint64) {
	m.WriteRegister(registerId, uint64(len(m.signerAccountId)), uint64(uintptr(unsafe.Pointer(&m.signerAccountId))))
}

func (m *MockSystem) SignerAccountPk(registerId uint64) {
	m.WriteRegister(registerId, uint64(len(m.signerAccountPk)), uint64(uintptr(unsafe.Pointer(&m.signerAccountPk))))
}

func (m *MockSystem) PredecessorAccountId(registerId uint64) {
	m.WriteRegister(registerId, uint64(len(m.predecessorAccountId)), uint64(uintptr(unsafe.Pointer(&m.predecessorAccountId))))
}

func (m *MockSystem) Input(registerId uint64) {
	m.WriteRegister(registerId, uint64(len(m.input)), uint64(uintptr(unsafe.Pointer(&m.input))))
}

func (m *MockSystem) BlockIndex() uint64 {
	return m.blockIndex
}

func (m *MockSystem) BlockTimestamp() uint64 {
	return m.blockTimestamp
}

func (m *MockSystem) EpochHeight() uint64 {
	return m.epochHeight
}

func (m *MockSystem) StorageUsage() uint64 {
	return m.storageUsage
}

// Economics API
func (m *MockSystem) AccountBalance(balancePtr uint64) {
	balance := *(*[]byte)(unsafe.Pointer(&balancePtr))
	balance = m.accountBalance.ToBE()
	fmt.Printf("balance: %v\n", balance)
}

func (m *MockSystem) AccountLockedBalance(balancePtr uint64) {
	balance := *(*[]byte)(unsafe.Pointer(&balancePtr))
	balance = m.accountBalance.ToBE()
	fmt.Printf("balance: %v\n", balance)
}

func (m *MockSystem) AttachedDeposit(balancePtr uint64) {
	balance := *(*[]byte)(unsafe.Pointer(&balancePtr))
	balance = m.attachedDeposit.ToBE()
	fmt.Printf("balance: %v\n", balance)
}

func (m *MockSystem) PrepaidGas() uint64 {
	return m.prepaidGas
}

func (m *MockSystem) UsedGas() uint64 {
	return m.usedGas
}

// Math API
func (m *MockSystem) RandomSeed(registerId uint64) {
	// seed := make([]byte, 32)
	// rand.Read(seed)
	// m.WriteRegister(registerId, uint64(len(seed)), uint64(uintptr(unsafe.Pointer(&seed[0]))))
}

func (m *MockSystem) Sha256(valueLen, valuePtr, registerId uint64) {
	// data := *(*[]byte)(unsafe.Pointer(&valuePtr))
	// hash := sha256.Sum256(data[:valueLen])
	// m.WriteRegister(registerId, uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))))
}

func (m *MockSystem) Keccak256(valueLen, valuePtr, registerId uint64) {
	// data := *(*[]byte)(unsafe.Pointer(&valuePtr))
	// hash := sha3.Sum256(data[:valueLen])
	// m.WriteRegister(registerId, uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))))
}

func (m *MockSystem) Keccak512(valueLen, valuePtr, registerId uint64) {
	// data := *(*[]byte)(unsafe.Pointer(&valuePtr))
	// hash := sha3.Sum512(data[:valueLen])
	// m.WriteRegister(registerId, uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))))
}

func (m *MockSystem) Ripemd160(valueLen, valuePtr, registerId uint64) {
	// data := *(*[]byte)(unsafe.Pointer(&valuePtr))
	// hasher := ripemd160.New()
	// hasher.Write(data[:valueLen])
	// hash := hasher.Sum(nil)
	// m.WriteRegister(registerId, uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))))
}

func (m *MockSystem) Ecrecover(hashLen, hashPtr, sigLen, sigPtr, v, malleabilityFlag, registerId uint64) uint64 {
	// Implement ECDSA recover functionality if necessary.
	return 0 // Placeholder return value
}

func (m *MockSystem) Ed25519Verify(sigLen, sigPtr, msgLen, msgPtr, pubKeyLen, pubKeyPtr uint64) uint64 {
	// Implement ED25519 verify functionality if necessary.
	return 0 // Placeholder return value
}

func (m *MockSystem) AltBn128G1Multiexp(valueLen, valuePtr, registerId uint64) {
	// Implement AltBn128G1Multiexp functionality if necessary.
}

func (m *MockSystem) AltBn128G1SumSystem(valueLen, valuePtr, registerId uint64) {
	// Implement AltBn128G1Sum functionality if necessary.
}

func (m *MockSystem) AltBn128PairingCheckSystem(valueLen, valuePtr uint64) uint64 {
	// Implement AltBn128PairingCheck functionality if necessary.
	return 0 // Placeholder return value
}

// Math API

// Miscellaneous API
func (m *MockSystem) ValueReturn(valueLen, valuePtr uint64) {
	value := *(*[]byte)(unsafe.Pointer(&valuePtr))
	// Normally this would return the value from the smart contract
	// Since this is a mock implementation, we can store it or just simulate the return
	m.Registers[0] = value[:valueLen]
}

func (m *MockSystem) PanicUtf8(len, ptr uint64) {
	message := *(*string)(unsafe.Pointer(&ptr))
	fmt.Printf("Panic: %s", message[:len])
}

func (m *MockSystem) LogUtf8(len, ptr uint64) {
	message := *(*string)(unsafe.Pointer(&ptr))
	fmt.Printf("Log: %s", message[:len])
}

func (m *MockSystem) LogUtf16(len, ptr uint64) {
	utf16Bytes := *(*[]uint16)(unsafe.Pointer(&ptr))
	message := string(utf16.Decode(utf16Bytes[:len]))
	fmt.Printf("Log: %s", message)
}

// Miscellaneous API

// Validator API
func (m *MockSystem) ValidatorStake(accountIdLen, accountIdPtr, stakePtr uint64) {
	accountId := *(*string)(unsafe.Pointer(&accountIdPtr))
	stakeAmount := uint64(0)

	// Simulate checking validator stake
	if accountId[:accountIdLen] == "validatorAccountId" {
		stakeAmount = 100000 // Example stake amount
	}

	stake := *(*uint64)(unsafe.Pointer(&stakePtr))
	stake = stakeAmount

	fmt.Printf("stake: %v\n", stake)
}

func (m *MockSystem) ValidatorTotalStake(stakePtr uint64) {
	totalStake := uint64(1000000) // Example total stake amount
	stake := *(*uint64)(unsafe.Pointer(&stakePtr))
	stake = totalStake

	fmt.Printf("stake: %v\n", stake)
}

// Promise API Actions
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
	amount := *(*uint64)(unsafe.Pointer(&amountPtr))

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
	amount := *(*uint64)(unsafe.Pointer(&amountPtr))
	promise := &m.Promises[promiseIndex]
	promise.Amount = amount
}

func (m *MockSystem) PromiseBatchActionStake(promiseIndex, amountPtr, publicKeyLen, publicKeyPtr uint64) {
	amount := *(*uint64)(unsafe.Pointer(&amountPtr))
	publicKey := *(*string)(unsafe.Pointer(&publicKeyPtr))

	promise := &m.Promises[promiseIndex]
	promise.Amount = amount
	promise.Arguments = append(promise.Arguments, publicKey[:publicKeyLen]...)
}

// Promises API
func (m *MockSystem) PromiseCreate(accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64 {
	accountId := *(*string)(unsafe.Pointer(&accountIdPtr))
	functionName := *(*string)(unsafe.Pointer(&functionNamePtr))
	arguments := *(*[]byte)(unsafe.Pointer(&argumentsPtr))
	amount := *(*uint64)(unsafe.Pointer(&amountPtr))

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
	accountId := *(*string)(unsafe.Pointer(&accountIdPtr))
	functionName := *(*string)(unsafe.Pointer(&functionNamePtr))
	arguments := *(*[]byte)(unsafe.Pointer(&argumentsPtr))
	amount := *(*uint64)(unsafe.Pointer(&amountPtr))

	promise := MockPromise{
		AccountId:    accountId[:accountIdLen],
		FunctionName: functionName[:functionNameLen],
		Arguments:    arguments[:argumentsLen],
		Amount:       amount,
		Gas:          gas,
		PromiseIndex: promiseIndex,
	}

	m.Promises = append(m.Promises, promise)

	return uint64(len(m.Promises)) - 1
}

func (m *MockSystem) PromiseAnd(promiseIdxPtr, promiseIdxCount uint64) uint64 {
	promiseIndexes := *(*[]uint64)(unsafe.Pointer(&promiseIdxPtr))

	for i := uint64(0); i < promiseIdxCount; i++ {
		// Here we can handle the promises that need to be combined
		fmt.Printf("promiseIndexes[i]: %v\n", promiseIndexes[i])
	}

	return uint64(len(m.Promises))
}

func (m *MockSystem) PromiseBatchCreate(accountIdLen, accountIdPtr uint64) uint64 {
	accountId := *(*string)(unsafe.Pointer(&accountIdPtr))

	promise := MockPromise{
		AccountId:    accountId[:accountIdLen],
		PromiseIndex: uint64(len(m.Promises)),
	}

	m.Promises = append(m.Promises, promise)

	return promise.PromiseIndex
}

func (m *MockSystem) PromiseBatchThen(promiseIndex, accountIdLen, accountIdPtr uint64) uint64 {
	accountId := *(*string)(unsafe.Pointer(&accountIdPtr))

	promise := MockPromise{
		AccountId:    accountId[:accountIdLen],
		PromiseIndex: promiseIndex,
	}

	m.Promises = append(m.Promises, promise)

	return uint64(len(m.Promises)) - 1
}
