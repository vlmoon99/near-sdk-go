package env

import (
	"errors"
	"math"
	"unsafe"

	"github.com/vlmoon99/jsonparser"
	"github.com/vlmoon99/near-sdk-go/system"
	"github.com/vlmoon99/near-sdk-go/types"
)

const RegisterExpectedErr = "Register was expected to have data because we just wrote it into it."

const AtomicOpRegister uint64 = math.MaxUint64 - 2

const EvictedRegister uint64 = math.MaxUint64 - 1

const DataIdRegister = 0

var StateKey = []byte("STATE")

const MinAccountIDLen uint64 = 2

const MaxAccountIDLen uint64 = 64

var nearBlockchainImports system.System = system.SystemNear{}

const (
	ErrExpectedDataInRegister            = "(REGISTER_ERROR): expected data in register, but found none"
	ErrInvalidAccountID                  = "(ACCOUNT_ERROR): invalid account ID"
	ErrKeyNotFound                       = "(STORAGE_ERROR): key not found"
	ErrValueNotFound                     = "(STORAGE_ERROR): value not found"
	ErrFailedToParseInput                = "(INPUT_ERROR): failed to parse input"
	ErrUnsupportedDataFormat             = "(FORMAT_ERROR): unsupported data format"
	ErrGettingAccountBalance             = "(BALANCE_ERROR): error while getting account balance"
	ErrGettingLockedAccountBalance       = "(BALANCE_ERROR): error while getting locked account balance"
	ErrGettingAttachedDeposit            = "(DEPOSIT_ERROR): error while getting attached deposit"
	ErrFailedToWriteValueInStorage       = "(STORAGE_ERROR): failed to write value in the storage by provided key, result of operation is 0"
	ErrKeyIsEmpty                        = "(STORAGE_ERROR): key is empty"
	ErrFailedToReadKey                   = "(STORAGE_ERROR): failed to read the key"
	ErrFailedToReadRegister              = "(REGISTER_ERROR): failed to read register"
	ErrCantRemoveDataByKey               = "(STORAGE_ERROR): can't remove data by that key"
	ErrFailedToReadEvictedRegister       = "(REGISTER_ERROR): failed to read evicted register"
	ErrStateNotFound                     = "(STATE_ERROR): state not found"
	ErrFailedToWriteStateToStorage       = "(STORAGE_ERROR): failed to write state to storage"
	ErrInvalidInputHashAndSignatureEmpty = "(INPUT_ERROR): invalid input: hash and signature must not be empty"
	PanicStrEcrecoverFailed              = "(PANIC): Ecrecover failed"
	ErrAccountIDMustNotBeEmpty           = "(ACCOUNT_ERROR): account ID must not be empty"
	ErrGettingValidatorStakeAmount       = "(STAKE_ERROR): error while getting validator stake amount"
	ErrGettingValidatorTotalStakeAmount  = "(STAKE_ERROR): error while getting validator total stake amount"
)

func SetEnv(system system.System) {
	nearBlockchainImports = system
}

// Registers API

func tryMethodIntoRegister(method func(uint64)) ([]byte, error) {
	method(AtomicOpRegister)

	return ReadRegisterSafe(AtomicOpRegister)
}

func methodIntoRegister(method func(uint64)) ([]byte, error) {
	data, err := tryMethodIntoRegister(method)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New(ErrExpectedDataInRegister)
	}
	return data, nil
}

func ReadRegisterSafe(registerId uint64) ([]byte, error) {
	length := nearBlockchainImports.RegisterLen(registerId)
	//TODO : If len == 0 - ExecutionError("WebAssembly trap: An `unreachable` opcode was executed.") for some reason, if we convert value into string erroe gone
	assertValidAccountId([]byte(string(length)))
	if length == 0 {
		return []byte{}, errors.New(ErrExpectedDataInRegister)
	}
	buffer := make([]byte, length)
	ptr := uint64(uintptr(unsafe.Pointer(&buffer[0])))
	nearBlockchainImports.ReadRegister(registerId, ptr)
	return buffer, nil
}

func WriteRegisterSafe(registerId uint64, data []byte) {
	if len(data) == 0 {
		return
	}

	ptr := uint64(uintptr(unsafe.Pointer(&data[0])))

	nearBlockchainImports.WriteRegister(registerId, uint64(len(data)), ptr)
}

// Registers API

// Storage API

func StorageWrite(key, value []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New(ErrKeyNotFound)
	}

	if len(value) == 0 {
		return false, errors.New(ErrValueNotFound)
	}

	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	valueLen := uint64(len(value))
	valuePtr := uint64(uintptr(unsafe.Pointer(&value[0])))

	return storageWriteRecursive(keyLen, keyPtr, valueLen, valuePtr, 0)
}

// TODO : try to undersatnd why on the first SotrageWrite we have 0 and on the second one we have 1, each other writing in this key - works well
func storageWriteRecursive(keyLen uint64, keyPtr uint64, valueLen uint64, valuePtr uint64, attempt int) (bool, error) {
	result := nearBlockchainImports.StorageWrite(keyLen, keyPtr, valueLen, valuePtr, EvictedRegister)

	if result == 1 {
		return true, nil
	}

	if result == 0 && attempt < 1 {
		return storageWriteRecursive(keyLen, keyPtr, valueLen, valuePtr, attempt+1)
	}

	return false, errors.New(ErrFailedToWriteValueInStorage)
}

func StorageRead(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New(ErrKeyIsEmpty)
	}
	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))
	result := nearBlockchainImports.StorageRead(keyLen, keyPtr, AtomicOpRegister)

	if result == 0 {
		return nil, errors.New(ErrFailedToReadKey)
	}

	value, err := ReadRegisterSafe(AtomicOpRegister)
	if err != nil {
		return nil, errors.New(ErrFailedToReadRegister)
	}

	return value, nil
}

func StorageRemove(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New(ErrKeyIsEmpty)
	}

	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	result := nearBlockchainImports.StorageRemove(keyLen, keyPtr, EvictedRegister)
	if result == 0 {
		return false, errors.New(ErrCantRemoveDataByKey)
	}

	return true, nil
}

func StorageGetEvicted() ([]byte, error) {
	value, err := ReadRegisterSafe(EvictedRegister)
	if err != nil {
		return nil, errors.New(ErrFailedToReadEvictedRegister + " " + err.Error())
	}

	return value, nil
}

func StorageHasKey(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New(ErrKeyIsEmpty)
	}

	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	result := nearBlockchainImports.StorageHasKey(keyLen, keyPtr)

	return result == 1, nil
}

func StateWrite(data []byte) error {

	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	valueLen := uint64(len(data))
	valuePtr := uint64(uintptr(unsafe.Pointer(&data[0])))

	_, err := storageWriteRecursive(keyLen, keyPtr, valueLen, valuePtr, 0)

	return err
}

func StateRead() ([]byte, error) {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := nearBlockchainImports.StorageRead(keyLen, keyPtr, EvictedRegister)
	if result == 0 {
		return nil, errors.New(ErrStateNotFound)
	}

	data, err := ReadRegisterSafe(EvictedRegister)
	if err != nil {
		return nil, errors.New(ErrFailedToReadRegister)
	}

	return data, nil
}

func StateExists() bool {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := nearBlockchainImports.StorageHasKey(keyLen, keyPtr)
	return result == 1
}

// Storage API

// Context API

func assertValidAccountId(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New(ErrInvalidAccountID)
	}
	return string(data), nil
}

func GetCurrentAccountId() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { nearBlockchainImports.CurrentAccountId(registerID) })
	if err != nil {
		LogString("Error in GetCurrentAccountId: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { nearBlockchainImports.SignerAccountId(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountPK() ([]byte, error) {
	data, err := methodIntoRegister(func(registerID uint64) { nearBlockchainImports.SignerAccountPk(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountPK: " + err.Error())
		return nil, err
	}

	return data, nil
}

func GetPredecessorAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { nearBlockchainImports.PredecessorAccountId(registerID) })
	if err != nil {
		LogString("Error in GetPredecessorAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetCurrentBlockHeight() uint64 {
	return nearBlockchainImports.BlockTimestamp()
}

func GetBlockTimeMs() uint64 {
	return nearBlockchainImports.BlockTimestamp() / 1_000_000
}

func GetEpochHeight() uint64 {
	return nearBlockchainImports.EpochHeight()
}

func GetStorageUsage() uint64 {
	return nearBlockchainImports.StorageUsage()
}

func detectInputType(decodedData []byte, keyPath ...string) ([]byte, string, error) {
	value, dataType, _, err := jsonparser.Get(decodedData, keyPath...)

	if err != nil {
		if dataType == jsonparser.NotExist {
			return nil, "not_exist", errors.New(ErrKeyNotFound)
		}
		return nil, "unknown", errors.New(ErrFailedToParseInput)
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
		return nil, "unknown", errors.New(ErrUnsupportedDataFormat)
	}
}

func ContractInput(options types.ContractInputOptions) ([]byte, string, error) {
	data, err := methodIntoRegister(func(registerID uint64) {
		nearBlockchainImports.Input(registerID)
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

// Economics API

func GetAccountBalance() (types.Uint128, error) {
	var data [16]byte
	nearBlockchainImports.AccountBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingAccountBalance)
	}
	return accountBalance, nil
}

func GetAccountLockedBalance() (types.Uint128, error) {
	var data [16]byte
	nearBlockchainImports.AccountLockedBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingLockedAccountBalance)
	}
	return accountBalance, nil
}

func GetAttachedDepoist() (types.Uint128, error) {
	var data [16]byte
	nearBlockchainImports.AttachedDeposit(uint64(uintptr(unsafe.Pointer(&data[0]))))
	attachedDeposit, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingAttachedDeposit)
	}
	return attachedDeposit, nil
}

func GetPrepaidGas() types.NearGas {
	return types.NearGas{Inner: nearBlockchainImports.PrepaidGas()}
}

func GetUsedGas() types.NearGas {
	return types.NearGas{Inner: nearBlockchainImports.UsedGas()}
}

// Economics API

// Math API

func GetRandomSeed() ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		nearBlockchainImports.RandomSeed(registerID)
	})
}

func Sha256Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		nearBlockchainImports.Sha256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Keccak256Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		nearBlockchainImports.Keccak256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Keccak512Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		nearBlockchainImports.Keccak512(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Ripemd160Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		nearBlockchainImports.Ripemd160(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func EcrecoverPubKey(hash, signature []byte, v byte, malleabilityFlag bool) ([]byte, error) {
	if len(hash) == 0 || len(signature) == 0 {
		return nil, errors.New(ErrInvalidInputHashAndSignatureEmpty)
	}

	return methodIntoRegister(func(registerID uint64) {
		result := nearBlockchainImports.Ecrecover(
			uint64(len(hash)), uint64(uintptr(unsafe.Pointer(&hash[0]))),
			uint64(len(signature)), uint64(uintptr(unsafe.Pointer(&signature[0]))),
			uint64(v), types.BoolToUnit(malleabilityFlag), registerID,
		)

		if result == 0 {
			PanicStr(PanicStrEcrecoverFailed)
		}
	})
}

func Ed25519VerifySig(signature [64]byte, message []byte, publicKey [32]byte) bool {
	result := nearBlockchainImports.Ed25519Verify(
		uint64(len(signature)), uint64(uintptr(unsafe.Pointer(&signature[0]))),
		uint64(len(message)), uint64(uintptr(unsafe.Pointer(&message[0]))),
		uint64(len(publicKey)), uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
	return result == 1
}

func AltBn128G1MultiExp(value []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		nearBlockchainImports.AltBn128G1Multiexp(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), registerID)
	})
}

func AltBn128G1Sum(value []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		nearBlockchainImports.AltBn128G1SumSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), registerID)
	})
}

func AltBn128PairingCheck(value []byte) bool {
	return nearBlockchainImports.AltBn128PairingCheckSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0])))) == 1
}

// Math API

// Validator API

func ValidatorStakeAmount(accountID []byte) (types.Uint128, error) {
	if len(accountID) == 0 {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrAccountIDMustNotBeEmpty)
	}

	var stakeData [16]byte
	nearBlockchainImports.ValidatorStake(uint64(len(accountID)), uint64(uintptr(unsafe.Pointer(&accountID[0]))), uint64(uintptr(unsafe.Pointer(&stakeData[0]))))

	validatorStakeAmount, err := types.LoadUint128LE(stakeData[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingValidatorStakeAmount)
	}

	return validatorStakeAmount, nil
}

func ValidatorTotalStakeAmount() (types.Uint128, error) {
	var stakeData [16]byte
	nearBlockchainImports.ValidatorTotalStake(uint64(uintptr(unsafe.Pointer(&stakeData[0]))))

	validatorTotalStakeAmount, err := types.LoadUint128LE(stakeData[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingValidatorTotalStakeAmount)
	}

	return validatorTotalStakeAmount, nil
}

// Validator API

// Miscellaneous API

func ContractValueReturn(inputBytes []byte) {
	nearBlockchainImports.ValueReturn(uint64(len(inputBytes)), uint64(uintptr(unsafe.Pointer(&inputBytes[0]))))
}

func PanicStr(input string) {
	inputBytes := []byte(input)
	inputLength := uint64(len(inputBytes))

	if inputLength == 0 {
		return
	}

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	nearBlockchainImports.PanicUtf8(inputLength, inputPtr)
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

	nearBlockchainImports.LogUtf8(inputLength, inputPtr)
}

func LogStringUtf8(inputBytes []byte) {

	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	nearBlockchainImports.LogUtf8(inputLength, inputPtr)
}

func LogStringUtf16(inputBytes []byte) {

	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	nearBlockchainImports.LogUtf16(inputLength, inputPtr)
}

// Miscellaneous API

// Promises API

func PromiseCreate(accountId []byte, functionName []byte, arguments []byte, amount types.Uint128, gas uint64) uint64 {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}
	return nearBlockchainImports.PromiseCreate(
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

func PromiseThen(promiseIdx uint64, accountId []byte, functionName []byte, arguments []byte, amount types.Uint128, gas uint64) uint64 {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}

	return nearBlockchainImports.PromiseThen(
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
	return nearBlockchainImports.PromiseAnd(uint64(uintptr(unsafe.Pointer(&promiseIndices[0]))), uint64(len(promiseIndices)))
}

func PromiseBatchCreate(accountId []byte) uint64 {
	return nearBlockchainImports.PromiseBatchCreate(uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

func PromiseBatchThen(promiseIdx uint64, accountId []byte) uint64 {
	return nearBlockchainImports.PromiseBatchThen(promiseIdx, uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

// Promises API

// Promises API Action

func PromiseBatchActionCreateAccount(promiseIdx uint64) {
	nearBlockchainImports.PromiseBatchActionCreateAccount(promiseIdx)
}

func PromiseBatchActionDeployContract(promiseIdx uint64, bytes []byte) {
	nearBlockchainImports.PromiseBatchActionDeployContract(promiseIdx, uint64(len(bytes)), uint64(uintptr(unsafe.Pointer(&bytes[0]))))
}

func PromiseBatchActionFunctionCall(promiseIdx uint64, functionName []byte, arguments []byte, amount types.Uint128, gas uint64) {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}
	nearBlockchainImports.PromiseBatchActionFunctionCall(promiseIdx,
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),
		gas,
	)
}

func PromiseBatchActionFunctionCallWeight(promiseIdx uint64, functionName []byte, arguments []byte, amount types.Uint128, gas uint64, weight uint64) {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}
	nearBlockchainImports.PromiseBatchActionFunctionCallWeight(promiseIdx,
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),
		gas,
		weight,
	)
}

func PromiseBatchActionTransfer(promiseIdx uint64, amount types.Uint128) {
	nearBlockchainImports.PromiseBatchActionTransfer(promiseIdx, uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))))
}

func PromiseBatchActionStake(promiseIdx uint64, amount types.Uint128, publicKey []byte) {
	nearBlockchainImports.PromiseBatchActionStake(
		promiseIdx,
		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

func PromiseBatchActionAddKeyWithFullAccess(promiseIdx uint64, publicKey []byte, nonce uint64) {
	nearBlockchainImports.PromiseBatchActionAddKeyWithFullAccess(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),

		nonce,
	)
}

func PromiseBatchActionAddKeyWithFunctionCall(promiseIdx uint64, publicKey []byte, nonce uint64, amount types.Uint128, receiverId []byte, functionName []byte) {
	nearBlockchainImports.PromiseBatchActionAddKeyWithFunctionCall(
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
	nearBlockchainImports.PromiseBatchActionDeleteKey(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

func PromiseBatchActionDeleteAccount(promiseIdx uint64, beneficiaryId []byte) {
	nearBlockchainImports.PromiseBatchActionDeleteAccount(
		promiseIdx,

		uint64(len(beneficiaryId)),
		uint64(uintptr(unsafe.Pointer(&beneficiaryId[0]))),
	)
}

func PromiseYieldCreate(functionName []byte, arguments []byte, gas uint64, gasWeight uint64) uint64 {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}
	return nearBlockchainImports.PromiseYieldCreate(
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
	return nearBlockchainImports.PromiseYieldResume(
		uint64(len(data)),
		uint64(uintptr(unsafe.Pointer(&data[0]))),

		uint64(len(data)),
		uint64(uintptr(unsafe.Pointer(&data[0]))),
	)
}

// Promises API Action

// Promise API Results

func PromiseResultsCount(data []byte, payload []byte) uint64 {
	return nearBlockchainImports.PromiseResultsCount()
}

func PromiseResult(resultIdx uint64) uint64 {
	return nearBlockchainImports.PromiseResult(resultIdx, AtomicOpRegister)
}

func PromiseReturn(promiseId uint64) {
	nearBlockchainImports.PromiseReturn(promiseId)
}

// Promise API Results
