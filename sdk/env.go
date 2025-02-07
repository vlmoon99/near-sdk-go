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
	length := registerLen(registerId)
	if length == 0 {
		return []byte{}, errors.New("expected data in register, but found none")
	}

	buffer := make([]byte, length)

	ptr := uint64(uintptr(unsafe.Pointer(&buffer[0])))

	readRegister(registerId, ptr)

	return buffer, nil
}

func writeRegisterSafe(registerId uint64, data []byte) {
	if len(data) == 0 {
		return
	}

	ptr := uint64(uintptr(unsafe.Pointer(&data[0])))

	writeRegister(registerId, uint64(len(data)), ptr)
}

// Registers

// Context API

func assertValidAccountId(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("invalid account ID")
	}
	return string(data), nil
}

func GetCurrentAccountId() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { currentAccountId(registerID) })
	if err != nil {
		LogString("Error in GetCurrentAccountId: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { signerAccountId(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountPK() ([]byte, error) {
	data, err := methodIntoRegister(func(registerID uint64) { signerAccountPk(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountPK: " + err.Error())
		return nil, err
	}

	return data, nil
}

func GetPredecessorAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { predecessorAccountId(registerID) })
	if err != nil {
		LogString("Error in GetPredecessorAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetCurrentBlockHeight() uint64 {
	return blockIndex()
}

func GetCurrentBlockTimeStamp() uint64 {
	return blockTimestamp()
}

func GetBlockTimeMs() uint64 {
	return blockTimestamp() / 1_000_000
}

func GetEpochHeight() uint64 {
	return epochHeight()
}

func GetStorageUsage() uint64 {
	return storageUsage()
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

func ContractInput(options ContractInputOptions) ([]byte, string, error) {
	data, err := methodIntoRegister(func(registerID uint64) {
		input(registerID)
	})
	if err != nil {
		LogString("Error in GetContractInput: " + err.Error())
		return nil, "", err
	}

	if options.IsRawBytes {
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
	valueReturn(uint64(len(inputBytes)), uint64(uintptr(unsafe.Pointer(&inputBytes[0]))))
}

func PanicStr(input string) {
	inputBytes := []byte(input)
	inputLength := uint64(len(inputBytes))

	if inputLength == 0 {
		return
	}

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	panicUtf8(inputLength, inputPtr)
}

func AbortExecution() {
	PanicStr("AbortExecution")
}

func LogString(input string) {
	inputBytes := []byte(input)
	inputLength := uint64(len(inputBytes))

	if inputLength == 0 {
		return
	}

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	logUtf8(inputLength, inputPtr)
}

func LogStringUtf8(inputBytes []byte) {

	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	logUtf8(inputLength, inputPtr)
}

func LogStringUtf16(inputBytes []byte) {

	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	logUtf16(inputLength, inputPtr)
}

// Miscellaneous API

// Economics API

func GetAccountBalance() Uint128 {
	var data [16]byte
	accountBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance := LoadUint128LE(data[:])
	return accountBalance
}

func GetAccountLockedBalance() Uint128 {
	var data [16]byte
	accountLockedBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance := LoadUint128LE(data[:])
	return accountBalance
}

func GetAttachedDepoist() Uint128 {
	var data [16]byte
	attachedDeposit(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance := LoadUint128LE(data[:])
	return accountBalance
}

func GetPrepaidGas() NearGas {
	return NearGas{prepaidGas()}
}

func GetUsedGas() NearGas {
	return NearGas{usedGas()}
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

	result := storageWriteSys(keyLen, keyPtr, valueLen, valuePtr, EvictedRegister)
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
	result := storageReadSys(keyLen, keyPtr, EvictedRegister)

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

	result := storageRemoveSys(keyLen, keyPtr, EvictedRegister)
	if result == 0 {
		return false, errors.New("Can't remove data by that key")
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

	result := storageHasKeySys(keyLen, keyPtr)
	return result == 1, nil
}

func StateRead() ([]byte, error) {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := storageReadSys(keyLen, keyPtr, 0)
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

	result := storageWriteSys(keyLen, keyPtr, valueLen, valuePtr, 0)
	if result == 0 {
		return errors.New("failed to write state to storage")
	}

	return nil
}

func StateExists() bool {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := storageHasKeySys(keyLen, keyPtr)
	return result == 1
}

// Storage API

// Math API

func GetRandomSeed() ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		randomSeed(registerID)
	})
}

func Sha256Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		sha256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Keccak256Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		keccak256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Keccak512Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		keccak512(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Ripemd160Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		ripemd160(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func EcrecoverPubKey(hash, signature []byte, v byte, malleabilityFlag bool) ([]byte, error) {
	if len(hash) == 0 || len(signature) == 0 {
		return nil, errors.New("invalid input: hash and signature must not be empty")
	}

	return methodIntoRegister(func(registerID uint64) {
		result := ecrecover(
			uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))),
			uint64(len(signature)), uint64(uintptr(unsafe.Pointer(&signature[0]))),
			uint64(v), BoolToUnit(malleabilityFlag), registerID,
		)

		if result == 0 {
			PanicStr("Ecrecover failed") // methodIntoRegister should catch this and return an error
		}
	})
}

func Ed25519VerifySig(signature [64]byte, message []byte, publicKey [32]byte) bool {
	result := ed25519Verify(
		uint64(len(signature)), uint64(uintptr(unsafe.Pointer(&signature[0]))),
		uint64(len(message)), uint64(uintptr(unsafe.Pointer(&message[0]))),
		uint64(len(publicKey)), uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
	return result == 1
}

func AltBn128G1MultiExp(value []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		altBn128G1Multiexp(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), registerID)
	})
}

func AltBn128G1Sum(value []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		altBn128G1SumSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), registerID)
	})
}

func AltBn128PairingCheck(value []byte) bool {
	return altBn128PairingCheckSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0])))) == 1
}

// Math API

// Validator API

func ValidatorStakeAmount(accountID []byte) (Uint128, error) {
	if len(accountID) == 0 {
		return Uint128{0, 0}, errors.New("account ID must not be empty")
	}

	var stakeData [16]byte
	validatorStake(uint64(len(accountID)), uint64(uintptr(unsafe.Pointer(&accountID[0]))), uint64(uintptr(unsafe.Pointer(&stakeData[0]))))

	return LoadUint128LE(stakeData[:]), nil
}

func ValidatorTotalStakeAmount() Uint128 {
	var stakeData [16]byte
	validatorTotalStake(uint64(uintptr(unsafe.Pointer(&stakeData[0]))))

	return LoadUint128LE(stakeData[:])
}

// Validator API

// Promises API

func PromiseCreate(accountId []byte, functionName []byte, arguments []byte, amount Uint128, gas uint64) uint64 {
	return promiseCreateSys(
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
	return promiseThenSys(
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
	return promiseAndSys(uint64(uintptr(unsafe.Pointer(&promiseIndices[0]))), uint64(len(promiseIndices)))
}

func PromiseBatchCreate(accountId []byte) uint64 {
	return promiseBatchCreateSys(uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

func PromiseBatchThen(promiseIdx uint64, accountId []byte) uint64 {
	return promiseBatchThenSys(promiseIdx, uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

// Promises API

// Promises API Action

func PromiseBatchActionCreateAccount(promiseIdx uint64) {
	promiseBatchActionCreateAccountSys(promiseIdx)
}

func PromiseBatchActionDeployContract(promiseIdx uint64, bytes []byte) {
	promiseBatchActionDeployContractSys(promiseIdx, uint64(len(bytes)), uint64(uintptr(unsafe.Pointer(&bytes[0]))))
}

func PromiseBatchActionFunctionCall(promiseIdx uint64, functionName []byte, arguments []byte, amount Uint128, gas uint64) {
	promiseBatchActionFunctionCallSys(promiseIdx,
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),
		gas,
	)
}

func PromiseBatchActionFunctionCallWeight(promiseIdx uint64, functionName []byte, arguments []byte, amount Uint128, gas uint64, weight uint64) {
	promiseBatchActionFunctionCallWeightSys(promiseIdx,
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
	promiseBatchActionTransferSys(promiseIdx, uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))))
}

func PromiseBatchActionStake(promiseIdx uint64, amount Uint128, publicKey []byte) {
	promiseBatchActionStakeSys(
		promiseIdx,
		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

func PromiseBatchActionAddKeyWithFullAccess(promiseIdx uint64, publicKey []byte, nonce uint64) {
	promiseBatchActionAddKeyWithFullAccessSys(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),

		nonce,
	)
}

func PromiseBatchActionAddKeyWithFunctionCall(promiseIdx uint64, publicKey []byte, nonce uint64, amount Uint128, receiverId []byte, functionName []byte) {
	promiseBatchActionAddKeyWithFunctionCallSys(
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
	promiseBatchActionDeleteKeySys(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

func PromiseBatchActionDeleteAccount(promiseIdx uint64, beneficiaryId []byte) {
	promiseBatchActionDeleteAccountSys(
		promiseIdx,

		uint64(len(beneficiaryId)),
		uint64(uintptr(unsafe.Pointer(&beneficiaryId[0]))),
	)
}

func PromiseYieldCreate(functionName []byte, arguments []byte, gas uint64, gasWeight uint64) uint64 {
	return promiseYieldCreateSys(
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
	return promiseYieldResumeSys(
		uint64(len(data)),
		uint64(uintptr(unsafe.Pointer(&data[0]))),

		uint64(len(data)),
		uint64(uintptr(unsafe.Pointer(&data[0]))),
	)
}

// Promises API Action

// Promise API Results
func PromiseResultsCount(data []byte, payload []byte) uint64 {
	return promiseResultsCountSys()
}

func PromiseResult(resultIdx uint64) uint64 {
	return promiseResultSys(resultIdx, AtomicOpRegister)
}

func PromiseReturn(promiseId uint64) {
	promiseReturnSys(promiseId)
}

// Promise API Results
