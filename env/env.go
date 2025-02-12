package env

import (
	"errors"
	"fmt"
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

var NearBlockchainImports system.System = system.SystemNear{}

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
	ErrFailedToWriteValueInStorage       = "(STORAGE_ERROR): failed to write value in the storage by provided key"
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
		return nil, errors.New(ErrExpectedDataInRegister)
	}
	return data, nil
}

func readRegisterSafe(registerId uint64) ([]byte, error) {

	length := NearBlockchainImports.RegisterLen(registerId)

	//TODO : If len == 0 - ExecutionError("WebAssembly trap: An `unreachable` opcode was executed.") for some reason, if we convert value into string erroe gone
	assertValidAccountId([]byte(string(length)))

	if length == 0 {
		return []byte{}, errors.New(ErrExpectedDataInRegister)
	}

	buffer := make([]byte, length)

	ptr := uint64(uintptr(unsafe.Pointer(&buffer[0])))

	NearBlockchainImports.ReadRegister(registerId, ptr)

	return buffer, nil
}

func writeRegisterSafe(registerId uint64, data []byte) {
	if len(data) == 0 {
		return
	}

	ptr := uint64(uintptr(unsafe.Pointer(&data[0])))

	NearBlockchainImports.WriteRegister(registerId, uint64(len(data)), ptr)
}

// Registers

// Context API

func assertValidAccountId(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New(ErrInvalidAccountID)
	}
	return string(data), nil
}

func GetCurrentAccountId() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { NearBlockchainImports.CurrentAccountId(registerID) })
	if err != nil {
		LogString("Error in GetCurrentAccountId: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { NearBlockchainImports.SignerAccountId(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountPK() ([]byte, error) {
	data, err := methodIntoRegister(func(registerID uint64) { NearBlockchainImports.SignerAccountPk(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountPK: " + err.Error())
		return nil, err
	}

	return data, nil
}

func GetPredecessorAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { NearBlockchainImports.PredecessorAccountId(registerID) })
	if err != nil {
		LogString("Error in GetPredecessorAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetCurrentBlockHeight() uint64 {
	return NearBlockchainImports.BlockTimestamp()
}

func GetBlockTimeMs() uint64 {
	return NearBlockchainImports.BlockTimestamp() / 1_000_000
}

func GetEpochHeight() uint64 {
	return NearBlockchainImports.EpochHeight()
}

func GetStorageUsage() uint64 {
	return NearBlockchainImports.StorageUsage()
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
		NearBlockchainImports.Input(registerID)
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
	NearBlockchainImports.ValueReturn(uint64(len(inputBytes)), uint64(uintptr(unsafe.Pointer(&inputBytes[0]))))
}

func PanicStr(input string) {
	inputBytes := []byte(input)
	inputLength := uint64(len(inputBytes))

	if inputLength == 0 {
		return
	}

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	NearBlockchainImports.PanicUtf8(inputLength, inputPtr)
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

	NearBlockchainImports.LogUtf8(inputLength, inputPtr)
}

func LogStringUtf8(inputBytes []byte) {

	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	NearBlockchainImports.LogUtf8(inputLength, inputPtr)
}

func LogStringUtf16(inputBytes []byte) {

	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	NearBlockchainImports.LogUtf16(inputLength, inputPtr)
}

// Miscellaneous API

// Economics API

func GetAccountBalance() (types.Uint128, error) {
	var data [16]byte
	NearBlockchainImports.AccountBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingAccountBalance)
	}
	return accountBalance, nil
}

func GetAccountLockedBalance() (types.Uint128, error) {
	var data [16]byte
	NearBlockchainImports.AccountLockedBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingLockedAccountBalance)
	}
	return accountBalance, nil
}

func GetAttachedDepoist() (types.Uint128, error) {
	var data [16]byte
	NearBlockchainImports.AttachedDeposit(uint64(uintptr(unsafe.Pointer(&data[0]))))
	attachedDeposit, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingAttachedDeposit)
	}
	return attachedDeposit, nil
}

func GetPrepaidGas() types.NearGas {
	return types.NearGas{Inner: NearBlockchainImports.PrepaidGas()}
}

func GetUsedGas() types.NearGas {
	return types.NearGas{Inner: NearBlockchainImports.UsedGas()}
}

// Economics API

// Storage API

func StorageWrite(key, value []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New(ErrKeyNotFound)
	}

	if len(value) == 0 {
		return false, errors.New(ErrValueNotFound + " " + string(value) + " " + fmt.Sprintf("%d", len(value)))
	}
	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	valueLen := uint64(len(value))
	valuePtr := uint64(uintptr(unsafe.Pointer(&value[0])))

	result := NearBlockchainImports.StorageWrite(keyLen, keyPtr, valueLen, valuePtr, EvictedRegister)
	if result == 0 {
		return false, errors.New(ErrFailedToWriteValueInStorage)
	}

	return true, nil
}

func StorageRead(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New(ErrKeyIsEmpty)
	}
	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))
	result := NearBlockchainImports.StorageRead(keyLen, keyPtr, EvictedRegister)

	if result == 0 {
		return nil, errors.New(ErrFailedToReadKey)
	}

	value, err := readRegisterSafe(EvictedRegister)
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

	result := NearBlockchainImports.StorageRemove(keyLen, keyPtr, EvictedRegister)
	if result == 0 {
		return false, errors.New(ErrCantRemoveDataByKey)
	}

	return true, nil
}

func StorageGetEvicted() ([]byte, error) {
	value, err := readRegisterSafe(EvictedRegister)
	if err != nil {
		return nil, errors.New(ErrFailedToReadEvictedRegister)
	}

	return value, nil
}

func StorageHasKey(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New(ErrKeyIsEmpty)
	}

	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	result := NearBlockchainImports.StorageHasKey(keyLen, keyPtr)
	return result == 1, nil
}

func StateRead() ([]byte, error) {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := NearBlockchainImports.StorageRead(keyLen, keyPtr, 0)
	if result == 0 {
		return nil, errors.New(ErrStateNotFound)
	}

	data, err := readRegisterSafe(0)
	if err != nil {
		return nil, errors.New(ErrFailedToReadRegister)
	}

	return data, nil
}

func StateWrite(data []byte) error {

	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	valueLen := uint64(len(data))
	valuePtr := uint64(uintptr(unsafe.Pointer(&data[0])))

	result := NearBlockchainImports.StorageWrite(keyLen, keyPtr, valueLen, valuePtr, 0)
	if result == 0 {
		return errors.New(ErrFailedToWriteStateToStorage)
	}

	return nil
}

func StateExists() bool {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := NearBlockchainImports.StorageHasKey(keyLen, keyPtr)
	return result == 1
}

// Storage API

// Math API

func GetRandomSeed() ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.RandomSeed(registerID)
	})
}

func Sha256Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.Sha256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Keccak256Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.Keccak256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Keccak512Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.Keccak512(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func Ripemd160Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.Ripemd160(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

func EcrecoverPubKey(hash, signature []byte, v byte, malleabilityFlag bool) ([]byte, error) {
	if len(hash) == 0 || len(signature) == 0 {
		return nil, errors.New(ErrInvalidInputHashAndSignatureEmpty)
	}

	return methodIntoRegister(func(registerID uint64) {
		result := NearBlockchainImports.Ecrecover(
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
	result := NearBlockchainImports.Ed25519Verify(
		uint64(len(signature)), uint64(uintptr(unsafe.Pointer(&signature[0]))),
		uint64(len(message)), uint64(uintptr(unsafe.Pointer(&message[0]))),
		uint64(len(publicKey)), uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
	return result == 1
}

func AltBn128G1MultiExp(value []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.AltBn128G1Multiexp(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), registerID)
	})
}

func AltBn128G1Sum(value []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.AltBn128G1SumSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), registerID)
	})
}

func AltBn128PairingCheck(value []byte) bool {
	return NearBlockchainImports.AltBn128PairingCheckSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0])))) == 1
}

// Math API

// Validator API

func ValidatorStakeAmount(accountID []byte) (types.Uint128, error) {
	if len(accountID) == 0 {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrAccountIDMustNotBeEmpty)
	}

	var stakeData [16]byte
	NearBlockchainImports.ValidatorStake(uint64(len(accountID)), uint64(uintptr(unsafe.Pointer(&accountID[0]))), uint64(uintptr(unsafe.Pointer(&stakeData[0]))))

	validatorStakeAmount, err := types.LoadUint128LE(stakeData[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingValidatorStakeAmount)
	}

	return validatorStakeAmount, nil
}

func ValidatorTotalStakeAmount() (types.Uint128, error) {
	var stakeData [16]byte
	NearBlockchainImports.ValidatorTotalStake(uint64(uintptr(unsafe.Pointer(&stakeData[0]))))

	validatorTotalStakeAmount, err := types.LoadUint128LE(stakeData[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingValidatorTotalStakeAmount)
	}

	return validatorTotalStakeAmount, nil
}

// Validator API

// Promises API

func PromiseCreate(accountId []byte, functionName []byte, arguments []byte, amount types.Uint128, gas uint64) uint64 {
	return NearBlockchainImports.PromiseCreate(
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
	return NearBlockchainImports.PromiseThen(
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
	return NearBlockchainImports.PromiseAnd(uint64(uintptr(unsafe.Pointer(&promiseIndices[0]))), uint64(len(promiseIndices)))
}

func PromiseBatchCreate(accountId []byte) uint64 {
	return NearBlockchainImports.PromiseBatchCreate(uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

func PromiseBatchThen(promiseIdx uint64, accountId []byte) uint64 {
	return NearBlockchainImports.PromiseBatchThen(promiseIdx, uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

// Promises API

// Promises API Action

func PromiseBatchActionCreateAccount(promiseIdx uint64) {
	NearBlockchainImports.PromiseBatchActionCreateAccount(promiseIdx)
}

func PromiseBatchActionDeployContract(promiseIdx uint64, bytes []byte) {
	NearBlockchainImports.PromiseBatchActionDeployContract(promiseIdx, uint64(len(bytes)), uint64(uintptr(unsafe.Pointer(&bytes[0]))))
}

func PromiseBatchActionFunctionCall(promiseIdx uint64, functionName []byte, arguments []byte, amount types.Uint128, gas uint64) {
	NearBlockchainImports.PromiseBatchActionFunctionCall(promiseIdx,
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),
		gas,
	)
}

func PromiseBatchActionFunctionCallWeight(promiseIdx uint64, functionName []byte, arguments []byte, amount types.Uint128, gas uint64, weight uint64) {
	NearBlockchainImports.PromiseBatchActionFunctionCallWeight(promiseIdx,
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
	NearBlockchainImports.PromiseBatchActionTransfer(promiseIdx, uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))))
}

func PromiseBatchActionStake(promiseIdx uint64, amount types.Uint128, publicKey []byte) {
	NearBlockchainImports.PromiseBatchActionStake(
		promiseIdx,
		uint64(uintptr(unsafe.Pointer(&amount.ToBE()[0]))),

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

func PromiseBatchActionAddKeyWithFullAccess(promiseIdx uint64, publicKey []byte, nonce uint64) {
	NearBlockchainImports.PromiseBatchActionAddKeyWithFullAccess(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),

		nonce,
	)
}

func PromiseBatchActionAddKeyWithFunctionCall(promiseIdx uint64, publicKey []byte, nonce uint64, amount types.Uint128, receiverId []byte, functionName []byte) {
	NearBlockchainImports.PromiseBatchActionAddKeyWithFunctionCall(
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
	NearBlockchainImports.PromiseBatchActionDeleteKey(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

func PromiseBatchActionDeleteAccount(promiseIdx uint64, beneficiaryId []byte) {
	NearBlockchainImports.PromiseBatchActionDeleteAccount(
		promiseIdx,

		uint64(len(beneficiaryId)),
		uint64(uintptr(unsafe.Pointer(&beneficiaryId[0]))),
	)
}

func PromiseYieldCreate(functionName []byte, arguments []byte, gas uint64, gasWeight uint64) uint64 {
	return NearBlockchainImports.PromiseYieldCreate(
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
	return NearBlockchainImports.PromiseYieldResume(
		uint64(len(data)),
		uint64(uintptr(unsafe.Pointer(&data[0]))),

		uint64(len(data)),
		uint64(uintptr(unsafe.Pointer(&data[0]))),
	)
}

// Promises API Action

// Promise API Results
func PromiseResultsCount(data []byte, payload []byte) uint64 {
	return NearBlockchainImports.PromiseResultsCount()
}

func PromiseResult(resultIdx uint64) uint64 {
	return NearBlockchainImports.PromiseResult(resultIdx, AtomicOpRegister)
}

func PromiseReturn(promiseId uint64) {
	NearBlockchainImports.PromiseReturn(promiseId)
}

// Promise API Results
