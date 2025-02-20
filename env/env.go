// This package provides implementations of low-level blockchain functions.
// Under the hood, it uses the system package, which provides raw environment imports from the Near Blockchain environment.
// Here, we wrap these functions and provide simpler methods for using low-level environment functions for smart contract development.
// For examples, please go to examples/integration_tests/main.go to see how these functions are called.
// All these functions are tested in integration_tests/src/main.rs. It's the best way to understand how they work in real examples.
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
	ErrPromiseResult                     = "(PROMISE_ERROR): no promise results available"
)

// SetEnv sets the environment to be used for Near Blockchain imports.
// It can be a mocked environment for unit tests or the default Near Blockchain imports for production.
//
// Parameters:
// - system: The system environment to be set.
func SetEnv(system system.System) {
	NearBlockchainImports = system
}

// Registers API

// tryMethodIntoRegister tries to execute the given method and reads the data from the register.
//
// Parameters:
// - method: The method to be executed.
//
// Returns:
// - []byte: The data read from the register.
// - error: An error if the method execution or data reading fails.
func tryMethodIntoRegister(method func(uint64)) ([]byte, error) {
	method(AtomicOpRegister)

	return ReadRegisterSafe(AtomicOpRegister)
}

// methodIntoRegister executes the given method and ensures the data is read from the register.
//
// Parameters:
// - method: The method to be executed.
//
// Returns:
// - []byte: The data read from the register.
// - error: An error if the method execution or data reading fails.
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

// ReadRegisterSafe reads the data from the specified register safely.
//
// Parameters:
// - registerId: The ID of the register to read from.
//
// Returns:
// - []byte: The data read from the register, or an error if the register reading fails.
func ReadRegisterSafe(registerId uint64) ([]byte, error) {
	length := NearBlockchainImports.RegisterLen(registerId)
	//TODO: If len == 0 - ExecutionError("WebAssembly trap: An `unreachable` opcode was executed.") for some reason, if we convert value into string, error gone
	assertValidAccountId([]byte(string(length)))
	if length == 0 {
		return []byte{}, errors.New(ErrExpectedDataInRegister)
	}
	buffer := make([]byte, length)
	ptr := uint64(uintptr(unsafe.Pointer(&buffer[0])))
	NearBlockchainImports.ReadRegister(registerId, ptr)
	return buffer, nil
}

// WriteRegisterSafe writes the given data to the specified register safely.
//
// Parameters:
// - registerId: The ID of the register to write to.
// - data: The data to be written to the register.
func WriteRegisterSafe(registerId uint64, data []byte) {
	if len(data) == 0 {
		return
	}

	ptr := uint64(uintptr(unsafe.Pointer(&data[0])))

	NearBlockchainImports.WriteRegister(registerId, uint64(len(data)), ptr)
}

// Registers API

// Storage API

// StorageWrite writes the given value to the specified key in storage.
//
// Parameters:
// - key: The key to write the value to.
// - value: The value to write.
//
// Returns:
// - bool: True if the value was successfully written, false otherwise.
// - error: An error if the key or value is empty or if the write operation fails.
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

// storageWriteRecursive attempts to write the value to the specified key in storage recursively.
//
// Parameters:
// - keyLen: The length of the key.
// - keyPtr: The pointer to the key.
// - valueLen: The length of the value.
// - valuePtr: The pointer to the value.
// - attempt: The current attempt number.
//
// Returns:
// - bool: True if the value was successfully written, false otherwise.
// - error: An error if the write operation fails after the allowed attempts.
func storageWriteRecursive(keyLen uint64, keyPtr uint64, valueLen uint64, valuePtr uint64, attempt int) (bool, error) {
	result := NearBlockchainImports.StorageWrite(keyLen, keyPtr, valueLen, valuePtr, EvictedRegister)

	if result == 1 {
		return true, nil
	}

	if result == 0 && attempt < 1 {
		return storageWriteRecursive(keyLen, keyPtr, valueLen, valuePtr, attempt+1)
	}

	return false, errors.New(ErrFailedToWriteValueInStorage)
}

// StorageRead reads the value associated with the given key from storage.
//
// Parameters:
// - key: The key to read the value for.
//
// Returns:
// - []byte: The value associated with the key, or an error if the read operation fails.
func StorageRead(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New(ErrKeyIsEmpty)
	}
	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))
	result := NearBlockchainImports.StorageRead(keyLen, keyPtr, AtomicOpRegister)

	if result == 0 {
		return nil, errors.New(ErrFailedToReadKey)
	}

	value, err := ReadRegisterSafe(AtomicOpRegister)
	if err != nil {
		return nil, errors.New(ErrFailedToReadRegister)
	}

	return value, nil
}

// StorageRemove removes the value associated with the given key from storage.
//
// Parameters:
// - key: The key to remove.
//
// Returns:
// - bool: True if the value was successfully removed, false otherwise.
// - error: An error if the key is empty or if the remove operation fails.
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

// StorageGetEvicted reads the value from the evicted register.
//
// Returns:
// - []byte: The value read from the evicted register, or an error if the read operation fails.
func StorageGetEvicted() ([]byte, error) {
	value, err := ReadRegisterSafe(EvictedRegister)
	if err != nil {
		return nil, errors.New(ErrFailedToReadEvictedRegister + " " + err.Error())
	}

	return value, nil
}

// StorageHasKey checks if the given key exists in storage.
//
// Parameters:
// - key: The key to check for existence.
//
// Returns:
// - bool: True if the key exists, false otherwise.
// - error: An error if the key is empty.
func StorageHasKey(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errors.New(ErrKeyIsEmpty)
	}

	keyLen := uint64(len(key))
	keyPtr := uint64(uintptr(unsafe.Pointer(&key[0])))

	result := NearBlockchainImports.StorageHasKey(keyLen, keyPtr)

	return result == 1, nil
}

// StateWrite writes the given data to the state.
//
// Parameters:
// - data: The data to write.
//
// Returns:
// - error: An error if the write operation fails.
func StateWrite(data []byte) error {

	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	valueLen := uint64(len(data))
	valuePtr := uint64(uintptr(unsafe.Pointer(&data[0])))

	_, err := storageWriteRecursive(keyLen, keyPtr, valueLen, valuePtr, 0)

	return err
}

// StateRead reads the data from the state.
//
// Returns:
// - []byte: The data read from the state, or an error if the read operation fails.
func StateRead() ([]byte, error) {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := NearBlockchainImports.StorageRead(keyLen, keyPtr, EvictedRegister)
	if result == 0 {
		return nil, errors.New(ErrStateNotFound)
	}

	data, err := ReadRegisterSafe(EvictedRegister)
	if err != nil {
		return nil, errors.New(ErrFailedToReadRegister)
	}

	return data, nil
}

// StateExists checks if the state exists.
//
// Returns:
// - bool: True if the state exists, false otherwise.
func StateExists() bool {
	keyLen := uint64(len(StateKey))
	keyPtr := uint64(uintptr(unsafe.Pointer(&StateKey[0])))

	result := NearBlockchainImports.StorageHasKey(keyLen, keyPtr)
	return result == 1
}

// Storage API

// Context API

// assertValidAccountId checks if the provided account ID is valid.
//
// Parameters:
// - data: The account ID data to validate.
//
// Returns:
// - string: The valid account ID as a string.
// - error: An error if the account ID is invalid.
func assertValidAccountId(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New(ErrInvalidAccountID)
	}
	return string(data), nil
}

// GetCurrentAccountId retrieves the current account ID.
//
// Returns:
// - string: The current account ID.
// - error: An error if the retrieval fails.
func GetCurrentAccountId() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { NearBlockchainImports.CurrentAccountId(registerID) })
	if err != nil {
		LogString("Error in GetCurrentAccountId: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

// GetSignerAccountID retrieves the signer account ID.
//
// Returns:
// - string: The signer account ID.
// - error: An error if the retrieval fails.
func GetSignerAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { NearBlockchainImports.SignerAccountId(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

// GetSignerAccountPK retrieves the public key of the signer account.
//
// Returns:
// - []byte: The public key of the signer account.
// - error: An error if the retrieval fails.
func GetSignerAccountPK() ([]byte, error) {
	data, err := methodIntoRegister(func(registerID uint64) { NearBlockchainImports.SignerAccountPk(registerID) })
	if err != nil {
		LogString("Error in GetSignerAccountPK: " + err.Error())
		return nil, err
	}

	return data, nil
}

// GetPredecessorAccountID retrieves the predecessor account ID.
//
// Returns:
// - string: The predecessor account ID.
// - error: An error if the retrieval fails.
func GetPredecessorAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { NearBlockchainImports.PredecessorAccountId(registerID) })
	if err != nil {
		LogString("Error in GetPredecessorAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

// GetCurrentBlockHeight retrieves the current block height.
//
// Returns:
// - uint64: The current block height.
func GetCurrentBlockHeight() uint64 {
	return NearBlockchainImports.BlockTimestamp()
}

// GetBlockTimeMs retrieves the block time in milliseconds.
//
// Returns:
// - uint64: The block time in milliseconds.
func GetBlockTimeMs() uint64 {
	return NearBlockchainImports.BlockTimestamp() / 1_000_000
}

// GetEpochHeight retrieves the current epoch height.
//
// Returns:
// - uint64: The current epoch height.
func GetEpochHeight() uint64 {
	return NearBlockchainImports.EpochHeight()
}

// GetStorageUsage retrieves the storage usage.
//
// Returns:
// - uint64: The storage usage.
func GetStorageUsage() uint64 {
	return NearBlockchainImports.StorageUsage()
}

// detectInputType detects the type of input data based on the provided key path.
//
// Parameters:
// - decodedData: The decoded data to analyze.
// - keyPath: The key path to locate the specific data element.
//
// Returns:
// - []byte: The detected value.
// - string: The type of the detected value.
// - error: An error if the detection fails.
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

// ContractInput retrieves the input data for the contract.
//
// Parameters:
// - options: Options specifying how to handle the input data.
//
// Returns:
// - []byte: The input data.
// - string: The type of the input data.
// - error: An error if the retrieval or processing fails.
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

// Economics API

// GetAccountBalance retrieves the current account balance.
//
// Returns:
// - types.Uint128: The current account balance.
// - error: An error if the retrieval fails.
func GetAccountBalance() (types.Uint128, error) {
	var data [16]byte
	NearBlockchainImports.AccountBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingAccountBalance)
	}
	return accountBalance, nil
}

// GetAccountLockedBalance retrieves the locked balance of the account.
//
// Returns:
// - types.Uint128: The locked balance of the account.
// - error: An error if the retrieval fails.
func GetAccountLockedBalance() (types.Uint128, error) {
	var data [16]byte
	NearBlockchainImports.AccountLockedBalance(uint64(uintptr(unsafe.Pointer(&data[0]))))
	accountBalance, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingLockedAccountBalance)
	}
	return accountBalance, nil
}

// GetAttachedDepoist retrieves the attached deposit.
//
// Returns:
// - types.Uint128: The attached deposit.
// - error: An error if the retrieval fails.
func GetAttachedDepoist() (types.Uint128, error) {
	var data [16]byte
	NearBlockchainImports.AttachedDeposit(uint64(uintptr(unsafe.Pointer(&data[0]))))
	attachedDeposit, err := types.LoadUint128LE(data[:])
	if err != nil {
		return types.Uint128{Hi: 0, Lo: 0}, errors.New(ErrGettingAttachedDeposit)
	}
	return attachedDeposit, nil
}

// GetPrepaidGas retrieves the prepaid gas.
//
// Returns:
// - types.NearGas: The prepaid gas.
func GetPrepaidGas() types.NearGas {
	return types.NearGas{Inner: NearBlockchainImports.PrepaidGas()}
}

// GetUsedGas retrieves the used gas.
//
// Returns:
// - types.NearGas: The used gas.
func GetUsedGas() types.NearGas {
	return types.NearGas{Inner: NearBlockchainImports.UsedGas()}
}

// Economics API

// Math API

// GetRandomSeed retrieves a random seed from the blockchain.
//
// Returns:
// - []byte: The random seed.
// - error: An error if the retrieval fails.
func GetRandomSeed() ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.RandomSeed(registerID)
	})
}

// Sha256Hash computes the SHA-256 hash of the given data.
//
// Parameters:
// - data: The data to hash.
//
// Returns:
// - []byte: The SHA-256 hash of the data.
// - error: An error if the hashing fails.
func Sha256Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.Sha256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

// Keccak256Hash computes the Keccak-256 hash of the given data.
//
// Parameters:
// - data: The data to hash.
//
// Returns:
// - []byte: The Keccak-256 hash of the data.
// - error: An error if the hashing fails.
func Keccak256Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.Keccak256(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

// Keccak512Hash computes the Keccak-512 hash of the given data.
//
// Parameters:
// - data: The data to hash.
//
// Returns:
// - []byte: The Keccak-512 hash of the data.
// - error: An error if the hashing fails.
func Keccak512Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.Keccak512(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

// Ripemd160Hash computes the RIPEMD-160 hash of the given data.
//
// Parameters:
// - data: The data to hash.
//
// Returns:
// - []byte: The RIPEMD-160 hash of the data.
// - error: An error if the hashing fails.
func Ripemd160Hash(data []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.Ripemd160(uint64(len(data)), uint64(uintptr(unsafe.Pointer(&data[0]))), registerID)
	})
}

// EcrecoverPubKey recovers the public key from the given hash and signature using the ECDSA algorithm.
//
// Parameters:
// - hash: The hash of the data.
// - signature: The signature of the data.
// - v: The recovery id (v).
// - malleabilityFlag: Indicates if malleable.
//
// Returns:
// - []byte: The recovered public key.
// - error: An error if the input hash or signature is empty, or if the recovery fails.
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

// Ed25519VerifySig verifies the Ed25519 signature of the given message with the public key.
//
// Parameters:
// - signature: The Ed25519 signature.
// - message: The message to verify.
// - publicKey: The Ed25519 public key.
//
// Returns:
// - bool: True if the signature is valid, false otherwise.
func Ed25519VerifySig(signature [64]byte, message []byte, publicKey [32]byte) bool {
	result := NearBlockchainImports.Ed25519Verify(
		uint64(len(signature)), uint64(uintptr(unsafe.Pointer(&signature[0]))),
		uint64(len(message)), uint64(uintptr(unsafe.Pointer(&message[0]))),
		uint64(len(publicKey)), uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
	return result == 1
}

// AltBn128G1MultiExp performs a multi-exponentiation on the given value using the alt_bn128 curve.
//
// Parameters:
// - value: The value to perform multi-exponentiation on.
//
// Returns:
// - []byte: The result of the multi-exponentiation.
// - error: An error if the operation fails.
func AltBn128G1MultiExp(value []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.AltBn128G1Multiexp(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), registerID)
	})
}

// AltBn128G1Sum performs a summation on the given value using the alt_bn128 curve.
//
// Parameters:
// - value: The value to perform summation on.
//
// Returns:
// - []byte: The result of the summation.
// - error: An error if the operation fails.
func AltBn128G1Sum(value []byte) ([]byte, error) {
	return methodIntoRegister(func(registerID uint64) {
		NearBlockchainImports.AltBn128G1SumSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0]))), registerID)
	})
}

// AltBn128PairingCheck performs a pairing check on the given value using the alt_bn128 curve.
//
// Parameters:
// - value: The value to perform pairing check on.
//
// Returns:
// - bool: True if the pairing check is successful, false otherwise.
func AltBn128PairingCheck(value []byte) bool {
	return NearBlockchainImports.AltBn128PairingCheckSystem(uint64(len(value)), uint64(uintptr(unsafe.Pointer(&value[0])))) == 1
}

// Math API

// Validator API

// ValidatorStakeAmount retrieves the stake amount for a given validator account ID.
//
// Parameters:
// - accountID: The account ID of the validator.
//
// Returns:
// - types.Uint128: The stake amount of the validator.
// - error: An error if the account ID is empty or if the retrieval fails.
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

// ValidatorTotalStakeAmount retrieves the total stake amount of all validators.
//
// Returns:
// - types.Uint128: The total stake amount of all validators.
// - error: An error if the retrieval fails.
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

// Miscellaneous API

// ContractValueReturn returns the specified value to the contract caller.
//
// Parameters:
// - inputBytes: The value to return.
func ContractValueReturn(inputBytes []byte) {
	NearBlockchainImports.ValueReturn(uint64(len(inputBytes)), uint64(uintptr(unsafe.Pointer(&inputBytes[0]))))
}

// PanicStr triggers a panic with the specified message.
//
// Parameters:
// - input: The panic message.
func PanicStr(input string) {
	inputBytes := []byte(input)
	inputLength := uint64(len(inputBytes))

	if inputLength == 0 {
		return
	}

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	NearBlockchainImports.PanicUtf8(inputLength, inputPtr)
}

// AbortExecution aborts the execution of the contract.
func AbortExecution() {
	PanicStr("AbortExecution")
}

// LogString logs the specified string message.
//
// Parameters:
// - input: The string message to log.
func LogString(input string) {
	inputBytes := []byte(input)
	inputLength := uint64(len(inputBytes))

	if inputLength == 0 {
		return
	}

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	NearBlockchainImports.LogUtf8(inputLength, inputPtr)
}

// LogStringUtf8 logs the specified UTF-8 encoded string message.
//
// Parameters:
// - inputBytes: The UTF-8 encoded string message to log.
func LogStringUtf8(inputBytes []byte) {
	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	NearBlockchainImports.LogUtf8(inputLength, inputPtr)
}

// LogStringUtf16 logs the specified UTF-16 encoded string message.
//
// Parameters:
// - inputBytes: The UTF-16 encoded string message to log.
func LogStringUtf16(inputBytes []byte) {
	inputLength := uint64(len(inputBytes))

	inputPtr := uint64(uintptr(unsafe.Pointer(&inputBytes[0])))

	NearBlockchainImports.LogUtf16(inputLength, inputPtr)
}

// Miscellaneous API

// Promises API

// PromiseCreate creates a promise to call a specified function on a specified account.
//
// Parameters:
// - accountId: The ID of the account to call the function on.
// - functionName: The name of the function to call.
// - arguments: The arguments to pass to the function.
// - amount: The amount to attach to the call.
// - gas: The amount of gas to attach to the call.
//
// Returns:
// - uint64: The ID of the created promise.
func PromiseCreate(accountId []byte, functionName []byte, arguments []byte, amount types.Uint128, gas uint64) uint64 {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}
	return NearBlockchainImports.PromiseCreate(
		uint64(len(accountId)),
		uint64(uintptr(unsafe.Pointer(&accountId[0]))),

		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToLE()[0]))),
		gas,
	)
}

// PromiseThen creates a dependent promise that will be executed after the initial promise resolves.
//
// Parameters:
// - promiseIdx: The ID of the initial promise.
// - accountId: The ID of the account to call the function on.
// - functionName: The name of the function to call.
// - arguments: The arguments to pass to the function.
// - amount: The amount to attach to the call.
// - gas: The amount of gas to attach to the call.
//
// Returns:
// - uint64: The ID of the created dependent promise.
func PromiseThen(promiseIdx uint64, accountId []byte, functionName []byte, arguments []byte, amount types.Uint128, gas uint64) uint64 {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}

	return NearBlockchainImports.PromiseThen(
		promiseIdx,
		uint64(len(accountId)),
		uint64(uintptr(unsafe.Pointer(&accountId[0]))),

		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToLE()[0]))),
		gas,
	)
}

// PromiseAnd combines multiple promises into a single promise that will resolve when all the combined promises resolve.
//
// Parameters:
// - promiseIndices: The IDs of the promises to combine.
//
// Returns:
// - uint64: The ID of the combined promise.
func PromiseAnd(promiseIndices []uint64) uint64 {
	return NearBlockchainImports.PromiseAnd(uint64(uintptr(unsafe.Pointer(&promiseIndices[0]))), uint64(len(promiseIndices)))
}

// PromiseBatchCreate creates a batch promise for a specified account.
//
// Parameters:
// - accountId: The ID of the account to create the batch promise for.
//
// Returns:
// - uint64: The ID of the created batch promise.
func PromiseBatchCreate(accountId []byte) uint64 {
	return NearBlockchainImports.PromiseBatchCreate(uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

// PromiseBatchThen creates a dependent batch promise that will be executed after the initial promise resolves.
//
// Parameters:
// - promiseIdx: The ID of the initial promise.
// - accountId: The ID of the account to create the dependent batch promise for.
//
// Returns:
// - uint64: The ID of the created dependent batch promise.
func PromiseBatchThen(promiseIdx uint64, accountId []byte) uint64 {
	return NearBlockchainImports.PromiseBatchThen(promiseIdx, uint64(len(accountId)), uint64(uintptr(unsafe.Pointer(&accountId[0]))))
}

// Promises API

// Promises API Action

// PromiseBatchActionCreateAccount creates a promise batch action to create a new account.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
func PromiseBatchActionCreateAccount(promiseIdx uint64) {
	NearBlockchainImports.PromiseBatchActionCreateAccount(promiseIdx)
}

// PromiseBatchActionDeployContract creates a promise batch action to deploy a contract.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - bytes: The contract code to deploy.
func PromiseBatchActionDeployContract(promiseIdx uint64, bytes []byte) {
	NearBlockchainImports.PromiseBatchActionDeployContract(promiseIdx, uint64(len(bytes)), uint64(uintptr(unsafe.Pointer(&bytes[0]))))
}

// PromiseBatchActionFunctionCall creates a promise batch action to call a function.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - functionName: The name of the function to call.
// - arguments: The arguments to pass to the function.
// - amount: The amount to attach to the call.
// - gas: The amount of gas to attach to the call.
func PromiseBatchActionFunctionCall(promiseIdx uint64, functionName []byte, arguments []byte, amount types.Uint128, gas uint64) {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}
	NearBlockchainImports.PromiseBatchActionFunctionCall(promiseIdx,
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToLE()[0]))),
		gas,
	)
}

// PromiseBatchActionFunctionCallWeight creates a promise batch action to call a function with a specified weight.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - functionName: The name of the function to call.
// - arguments: The arguments to pass to the function.
// - amount: The amount to attach to the call.
// - gas: The amount of gas to attach to the call.
// - weight: The weight of the call.
func PromiseBatchActionFunctionCallWeight(promiseIdx uint64, functionName []byte, arguments []byte, amount types.Uint128, gas uint64, weight uint64) {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}
	NearBlockchainImports.PromiseBatchActionFunctionCallWeight(promiseIdx,
		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),

		uint64(len(arguments)),
		uint64(uintptr(unsafe.Pointer(&arguments[0]))),

		uint64(uintptr(unsafe.Pointer(&amount.ToLE()[0]))),
		gas,
		weight,
	)
}

// PromiseBatchActionTransfer creates a promise batch action to transfer an amount.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - amount: The amount to transfer.
func PromiseBatchActionTransfer(promiseIdx uint64, amount types.Uint128) {
	NearBlockchainImports.PromiseBatchActionTransfer(promiseIdx, uint64(uintptr(unsafe.Pointer(&amount.ToLE()[0]))))
}

// PromiseBatchActionStake creates a promise batch action to stake an amount.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - amount: The amount to stake.
// - publicKey: The public key to stake with.
func PromiseBatchActionStake(promiseIdx uint64, amount types.Uint128, publicKey []byte) {
	NearBlockchainImports.PromiseBatchActionStake(
		promiseIdx,
		uint64(uintptr(unsafe.Pointer(&amount.ToLE()[0]))),

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

// PromiseBatchActionAddKeyWithFullAccess creates a promise batch action to add a key with full access.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - publicKey: The public key to add.
// - nonce: The nonce for the key.
func PromiseBatchActionAddKeyWithFullAccess(promiseIdx uint64, publicKey []byte, nonce uint64) {
	NearBlockchainImports.PromiseBatchActionAddKeyWithFullAccess(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),

		nonce,
	)
}

// PromiseBatchActionAddKeyWithFunctionCall creates a promise batch action to add a key with function call access.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - publicKey: The public key to add.
// - nonce: The nonce for the key.
// - amount: The amount to attach to the call.
// - receiverId: The ID of the receiver account.
// - functionName: The name of the function to call.
func PromiseBatchActionAddKeyWithFunctionCall(promiseIdx uint64, publicKey []byte, nonce uint64, amount types.Uint128, receiverId []byte, functionName []byte) {
	NearBlockchainImports.PromiseBatchActionAddKeyWithFunctionCall(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),

		nonce,
		uint64(uintptr(unsafe.Pointer(&amount.ToLE()[0]))),

		uint64(len(receiverId)),
		uint64(uintptr(unsafe.Pointer(&receiverId[0]))),

		uint64(len(functionName)),
		uint64(uintptr(unsafe.Pointer(&functionName[0]))),
	)
}

// PromiseBatchActionDeleteKey creates a promise batch action to delete a key.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - publicKey: The public key to delete.
func PromiseBatchActionDeleteKey(promiseIdx uint64, publicKey []byte) {
	NearBlockchainImports.PromiseBatchActionDeleteKey(
		promiseIdx,

		uint64(len(publicKey)),
		uint64(uintptr(unsafe.Pointer(&publicKey[0]))),
	)
}

// PromiseBatchActionDeleteAccount creates a promise batch action to delete an account.
//
// Parameters:
// - promiseIdx: The ID of the promise batch.
// - beneficiaryId: The ID of the beneficiary account.
func PromiseBatchActionDeleteAccount(promiseIdx uint64, beneficiaryId []byte) {
	NearBlockchainImports.PromiseBatchActionDeleteAccount(
		promiseIdx,

		uint64(len(beneficiaryId)),
		uint64(uintptr(unsafe.Pointer(&beneficiaryId[0]))),
	)
}

// PromiseYieldCreate creates a yield promise to call a specified function.
//
// Parameters:
// - functionName: The name of the function to call.
// - arguments: The arguments to pass to the function.
// - gas: The amount of gas to attach to the call.
// - gasWeight: The weight of the gas.
//
// Returns:
// - uint64: The ID of the created yield promise.
func PromiseYieldCreate(functionName []byte, arguments []byte, gas uint64, gasWeight uint64) uint64 {
	if len(arguments) == 0 {
		arguments = []byte("{}")
	}
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

// PromiseYieldResume resumes a yield promise with the specified data and payload.
//
// Parameters:
// - data: The data to resume the promise with.
// - payload: The payload to resume the promise with.
//
// Returns:
// - uint32: The status of the resumed promise.
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

// PromiseResultsCount retrieves the count of promise results.
//
// Returns:
// - uint64: The count of promise results.
func PromiseResultsCount() uint64 {
	return NearBlockchainImports.PromiseResultsCount()
}

// PromiseResult retrieves the result of a specified promise.
//
// Parameters:
// - resultIdx: The index of the promise result to retrieve.
//
// Returns:
// - []byte: The result of the specified promise.
// - error: An error if there are no promise results or if the retrieval fails.
func PromiseResult(resultIdx uint64) ([]byte, error) {
	if PromiseResultsCount() == 0 {
		return nil, errors.New(ErrPromiseResult)
	}

	NearBlockchainImports.PromiseResult(resultIdx, AtomicOpRegister)

	value, err := ReadRegisterSafe(AtomicOpRegister)
	if err != nil {
		return nil, errors.New(ErrFailedToReadRegister)
	}

	return value, nil
}

// PromiseReturn returns the result of a specified promise.
//
// Parameters:
// - promiseId: The ID of the promise to return the result for.
func PromiseReturn(promiseId uint64) {
	NearBlockchainImports.PromiseReturn(promiseId)
}

// Promise API Results
