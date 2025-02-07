package sdk

import (
	"errors"
	"math"
	"unsafe"

	"github.com/vlmoon99/jsonparser"
)

const RegisterExpectedErr = "Register was expected to have data because we just wrote it into it."

const AtomicOpRegister uint64 = ^uint64(2)

const EvictedRegister uint64 = math.MaxUint64 - 1

const DataIdRegister = 0

var StateKey = []byte("STATE")

const MinAccountIDLen uint64 = 2

const MaxAccountIDLen uint64 = 64

// Registers

func tryMethodIntoRegister(method func(uint64)) ([]byte, error) {
	method(AtomicOpRegister)

	return readRegisterSafe(AtomicOpRegister)
}

func methodIntoRegister(method func(uint64)) ([]byte, error) {
	data, err := tryMethodIntoRegister(method)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("expected data in register, but found none")
	}
	return data, nil
}

func readRegisterSafe(registerId uint64) ([]byte, error) {
	length := RegisterLen(registerId)
	if length == 0 {
		return []byte{}, errors.New("expected data in register, but found none")
	}

	buffer := make([]byte, length)

	ptr := uint64(uintptr(unsafe.Pointer(&buffer[0])))

	ReadRegister(registerId, ptr)

	return buffer, nil
}

func writeRegisterSafe(registerId uint64, data []byte) {
	if len(data) == 0 {
		return
	}

	ptr := uint64(uintptr(unsafe.Pointer(&data[0])))

	WriteRegister(registerId, uint64(len(data)), ptr)
}

// Registers

// Context API

func assertValidAccountId(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("invalid account ID")
	}
	return string(data), nil
}

func GetCurrentAccountID() (string, error) {
	CurrentAccountId(AtomicOpRegister)
	data, err := methodIntoRegister(func(registerID uint64) { CurrentAccountId(registerID) })
	if err != nil {
		LogString("Error in GetCurrentAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { SignerAccountId(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountPK() ([]byte, error) {
	data, err := methodIntoRegister(func(registerID uint64) { SignerAccountPk(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountPK: " + err.Error())
		return nil, err
	}

	return data, nil
}

func GetPredecessorAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { PredecessorAccountId(registerID) })
	if err != nil {
		LogString("Error in GetPredecessorAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetCurrentBlockHeight() uint64 {
	return BlockIndex()
}

func GetCurrentBlockTimeStamp() uint64 {
	return BlockTimestamp()
}

func GetBlockTimeMs() uint64 {
	return BlockTimestamp() / 1_000_000
}

func GetEpochHeight() uint64 {
	return EpochHeight()
}

func GetStorageUsage() uint64 {
	return StorageUsage()
}

func detectInputType(decodedData []byte, keyPath ...string) ([]byte, string, error) {
	value, dataType, _, err := jsonparser.Get(decodedData, keyPath...)

	if err != nil {
		if dataType == jsonparser.NotExist {
			return nil, "not_exist", errors.New("key not found")
		}
		return nil, "unknown", errors.New("failed to parse input")
	}

	switch dataType {
	case jsonparser.String:
		return value, "string", nil
	case jsonparser.Number:
		return value, "number", nil
	case jsonparser.Boolean:
		return value, "boolean", nil
	case jsonparser.Array:
		return value, "array", nil
	case jsonparser.Object:
		return value, "object", nil
	case jsonparser.Null:
		return nil, "null", nil
	default:
		return nil, "unknown", errors.New("unsupported data format")
	}
}

func ContractInput(isRaw bool) ([]byte, string, error) {

	data, err := methodIntoRegister(func(registerID uint64) {
		Input(registerID)
	})
	if err != nil {
		LogString("Error in GetContractInput: " + err.Error())
		return nil, "", err
	}
	if isRaw {
		return data, "rawBytes", nil
	}

	parsedData, detectedType, err := detectInputType(data)
	if err != nil {
		LogString("Failed to detect input type: " + err.Error())
		return nil, "", err
	}

	return parsedData, detectedType, nil
}

// Context API

// Miscellaneous API

func ContractValueReturn(inputBytes []byte) {
	ValueReturn(uint64(len(inputBytes)), uint64(uintptr(unsafe.Pointer(&inputBytes[0]))))
}

func PanicStr(input string) {
	inputBytes := []byte(input)
	inputLength := uint64(len(inputBytes))

	if inputLength == 0 {
		return
	}

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	PanicUtf8(inputLength, inputPtr)
}

func AbortExecution() {
	Panic()
}

func LogString(input string) {
	inputBytes := []byte(input)
	inputLength := uint64(len(inputBytes))

	if inputLength == 0 {
		return
	}

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	LogUtf8(inputLength, inputPtr)
}

func LogStringUtf8(inputBytes []byte) {

	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	LogUtf8(inputLength, inputPtr)
}

func LogStringUtf16(inputBytes []byte) {

	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	LogUtf16(inputLength, inputPtr)
}

// Miscellaneous API

// Economics API

func GetAccountBalance() Uint128 {
	var data [16]byte
	AccountBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance := LoadUint128LE(data[:])
	return accountBalance
}

func GetAccountLockedBalance() Uint128 {
	var data [16]byte
	AccountLockedBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance := LoadUint128LE(data[:])
	return accountBalance
}

func GetAttachedDepoist() Uint128 {
	var data [16]byte
	AttachedDeposit(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance := LoadUint128LE(data[:])
	return accountBalance
}

func GetPrepaidGas() NearGas {
	return NearGas{PrepaidGas()}
}

func GetUsedGas() NearGas {
	return NearGas{UsedGas()}
}

// Economics API

// Storage API

func StorageWrite(key, value []byte) (bool, error) {
	if len(key) == 0 || len(value) == 0 {
		return false, errors.New("key not found")
	}

	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	valueLen := uint64(len(value))
	valuePtr := uint64(uintptr(unsafe.Pointer(&value[0])))

	result := StorageWriteSys(keyLen, keyPtr, valueLen, valuePtr, EvictedRegister)
	if result == 0 {
		return false, errors.New("Failed to Write value in the storage by provided key")
	}

	return true, nil
}

func StorageRead(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("key is empty")
	}
	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))
	result := StorageReadSys(keyLen, keyPtr, EvictedRegister)

	if result == 0 {
		return nil, errors.New("Failed to Read the key")
	}

	value, err := readRegisterSafe(EvictedRegister)
	if err != nil {
		return nil, errors.New("Failed to Read Register")
	}

	return value, nil
}

func StorageRemove(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New("key is empty")
	}

	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	result := StorageRemoveSys(keyLen, keyPtr, EvictedRegister)
	if result == 0 {
		return false, nil
	}

	return true, nil
}

func StorageGetEvicted() ([]byte, error) {
	value, err := readRegisterSafe(EvictedRegister)
	if err != nil {
		return nil, errors.New("failed to read evicted register")
	}

	return value, nil
}

func StorageHasKey(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New("key is empty")
	}

	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	result := StorageHasKeySys(keyLen, keyPtr)
	return result == 1, nil
}

func StateRead() ([]byte, error) {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := StorageReadSys(keyLen, keyPtr, 0)
	if result == 0 {
		return nil, errors.New("state not found")
	}

	data, err := readRegisterSafe(0)
	if err != nil {
		return nil, errors.New("failed to read register")
	}

	return data, nil
}

func StateWrite(data []byte) error {

	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	valueLen := uint64(len(data))
	valuePtr := uint64(uintptr(unsafe.Pointer(&data[0])))

	result := StorageWriteSys(keyLen, keyPtr, valueLen, valuePtr, 0)
	if result == 0 {
		return errors.New("failed to write state to storage")
	}

	return nil
}

func StateExists() bool {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := StorageHasKeySys(keyLen, keyPtr)
	return result == 1
}

// Storage API

// Math API

func GetRandomSeed() ([]byte, error) {
	RandomSeed(AtomicOpRegister)
	return readRegisterSafe(AtomicOpRegister)
}

func Sha256Hash(data []byte) ([]byte, error) {
	Sha256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), AtomicOpRegister)
	return readRegisterSafe(AtomicOpRegister)
}

func Keccak256Hash(data []byte) ([]byte, error) {
	Keccak256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), AtomicOpRegister)
	return readRegisterSafe(AtomicOpRegister)
}

func Keccak512Hash(data []byte) ([]byte, error) {
	Keccak512(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), AtomicOpRegister)
	return readRegisterSafe(AtomicOpRegister)
}

func Ripemd160Hash(data []byte) ([]byte, error) {
	Ripemd160(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), AtomicOpRegister)
	return readRegisterSafe(AtomicOpRegister)
}

func EcrecoverPubKey(hash, signature []byte, v byte, malleabilityFlag bool) ([]byte, error) {
	if len(hash) == 0 || len(signature) == 0 {
		return nil, errors.New("invalid input: hash and signature must not be empty")
	}

	result := Ecrecover(
		uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))),
		uint64(len(signature)), uint64(uintptr(unsafe.Pointer(&signature[0]))),
		uint64(v), BoolToUnit(malleabilityFlag), AtomicOpRegister,
	)

	if result == 0 {
		return nil, errors.New("Ecrecover failed")
	}

	return readRegisterSafe(AtomicOpRegister)
}

func Ed25519VerifySig(signature [64]byte, message []byte, publicKey [32]byte) bool {
	result := Ed25519Verify(
		uint64(len(signature)), uint64(uintptr(unsafe.Pointer(&signature[0]))),
		uint64(len(message)), uint64(uintptr(unsafe.Pointer(&message[0]))),
		uint64(len(publicKey)), uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)

	return result == 1
}

func AltBn128G1MultiExp(value []byte) ([]byte, error) {
	AltBn128G1Multiexp(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), AtomicOpRegister)
	return readRegisterSafe(AtomicOpRegister)
}

func AltBn128G1Sum(value []byte) ([]byte, error) {
	AltBn128G1SumSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), AtomicOpRegister)
	return readRegisterSafe(AtomicOpRegister)
}

func AltBn128PairingCheck(value []byte) bool {
	return AltBn128PairingCheckSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0])))) == 1
}

// Math API

// Validator API

func ValidatorStakeAmount(accountID []byte) (Uint128, error) {
	if len(accountID) == 0 {
		return Uint128{0, 0}, errors.New("account ID must not be empty")
	}

	var stakeData [16]byte
	ValidatorStake(uint64(len(accountID)), uint64(uintptr(unsafe.Pointer(&accountID[0]))), uint64(uintptr(unsafe.Pointer(&stakeData[0]))))

	return LoadUint128LE(stakeData[:]), nil
}

func ValidatorTotalStakeAmount() Uint128 {
	var stakeData [16]byte
	ValidatorTotalStake(uint64(uintptr(unsafe.Pointer(&stakeData[0]))))

	return LoadUint128LE(stakeData[:])
}

// Validator API

// Promises API

func PromiseCreate(accountId []byte, functionName []byte, arguments []byte, amount Uint128, gas uint64) uint64 {
	return PromiseCreateSys(
		uint64(len(accountId)),
		uint64(uintptr(unsafe.Pointer(&accountId[0]))),

		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),
		gas,
	)
}

func PromiseThen(promiseIdx uint64, accountId []byte, functionName []byte, arguments []byte, amount Uint128, gas uint64) uint64 {
	return PromiseThenSys(
		promiseIdx,
		uint64(len(accountId)),
		uint64(uintptr(unsafe.Pointer(&accountId[0]))),

		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),
		gas,
	)
}

func PromiseAnd(promiseIndices []uint64) uint64 {
	return PromiseAndSys(uint64(uintptr(unsafe.Pointer(&promiseIndices[0]))), uint64(len(promiseIndices)))
}

func PromiseBatchCreate(accountId []byte) uint64 {
	return PromiseBatchCreateSys(uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

func PromiseBatchThen(promiseIdx uint64, accountId []byte) uint64 {
	return PromiseBatchThenSys(promiseIdx, uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

// Promises API

// Promises API Action

func PromiseBatchActionCreateAccount(promiseIdx uint64) {
	PromiseBatchActionCreateAccountSys(promiseIdx)
}

func PromiseBatchActionDeployContract(promiseIdx uint64, bytes []byte) {
	PromiseBatchActionDeployContractSys(promiseIdx, uint64(len(bytes)), uint64(uintptr(unsafe.Pointer(&bytes[0]))))
}

func PromiseBatchActionFunctionCall(promiseIdx uint64, functionName []byte, arguments []byte, amount Uint128, gas uint64) {
	PromiseBatchActionFunctionCallSys(promiseIdx,
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),
		gas,
	)
}

func PromiseBatchActionFunctionCallWeight(promiseIdx uint64, functionName []byte, arguments []byte, amount Uint128, gas uint64, weight uint64) {
	PromiseBatchActionFunctionCallWeightSys(promiseIdx,
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),
		gas,
		weight,
	)
}

func PromiseBatchActionTransfer(promiseIdx uint64, amount Uint128) {
	PromiseBatchActionTransferSys(promiseIdx, uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))))
}

func PromiseBatchActionStake(promiseIdx uint64, amount Uint128, publicKey []byte) {
	PromiseBatchActionStakeSys(
		promiseIdx,
		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

func PromiseBatchActionAddKeyWithFullAccess(promiseIdx uint64, publicKey []byte, nonce uint64) {
	PromiseBatchActionAddKeyWithFullAccessSys(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),

		nonce,
	)
}

func PromiseBatchActionAddKeyWithFunctionCall(promiseIdx uint64, publicKey []byte, nonce uint64, amount Uint128, receiverId []byte, functionName []byte) {
	PromiseBatchActionAddKeyWithFunctionCallSys(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),

		nonce,
		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),

		uint64(len(receiverId)),
		uint64(uintptr(unsafe.Pointer(&receiverId[0]))),

		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),
	)
}

func PromiseBatchActionDeleteKey(promiseIdx uint64, publicKey []byte) {
	PromiseBatchActionDeleteKeySys(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

func PromiseBatchActionDeleteAccount(promiseIdx uint64, beneficiaryId []byte) {
	PromiseBatchActionDeleteAccountSys(
		promiseIdx,

		uint64(len(beneficiaryId)),
		uint64(uintptr(unsafe.Pointer(&beneficiaryId[0]))),
	)
}

func PromiseYieldCreate(functionName []byte, arguments []byte, gas uint64, gasWeight uint64) uint64 {
	return PromiseYieldCreateSys(
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),
		gas,
		gasWeight,
		DataIdRegister,
	)
}

func PromiseYieldResume(data []byte, payload []byte) uint32 {
	return PromiseYieldResumeSys(
		uint64(len(data)),
		uint64(uintptr(unsafe.Pointer(&data[0]))),

		uint64(len(data)),
		uint64(uintptr(unsafe.Pointer(&data[0]))),
	)
}

// Promises API Action

// Promise API Results
func PromiseResultsCount(data []byte, payload []byte) uint64 {
	return PromiseResultsCountSys()
}

func PromiseResult(resultIdx uint64) uint64 {
	return PromiseResultSys(resultIdx, AtomicOpRegister)
}

func PromiseReturn(promiseId uint64) {
	PromiseReturnSys(promiseId)
}

// Promise API Results
