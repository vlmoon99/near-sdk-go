package main

import (
	"errors"
	"fmt"
	"math"
	"unsafe"

	"github.com/buger/jsonparser"
)

// Error message when a register is expected to have data but does not.
const RegisterExpectedErr = "Register was expected to have data because we just wrote it into it."

// Register used internally for atomic operations. This register is safe to use by the user,
// since it only needs to be untouched while methods of `Environment` execute, which is guaranteed
// as guest code is not parallel.
const AtomicOpRegister uint64 = ^uint64(2)

// Register used to record evicted values from the storage.
const EvictedRegister uint64 = math.MaxUint64 - 1

// Key used to store the state of the contract.
var StateKey = []byte("STATE")

// The minimum length of a valid account ID.
const MinAccountIDLen uint64 = 2

// The maximum length of a valid account ID.
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
		return nil, "unknown", fmt.Errorf("failed to parse input: %v", err)
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

func GetSmartContractInput() ([]byte, string, error) {

	data, err := methodIntoRegister(func(registerID uint64) {
		Input(registerID)
	})
	if err != nil {
		LogString("Error in GetSmartContractInput: " + err.Error())
		return nil, "", err
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

func SmartContractValueReturn(inputBytes []byte) {
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

// ###############
// # Math API #
// ###############
// pub fn random_seed() -> Vec<u8> {
//     random_seed_array().to_vec()
// }
// pub fn random_seed_array() -> [u8; 32] {
//     //* SAFETY: random_seed syscall will always generate 32 bytes inside of the atomic op register
//     //*         so the read will have a sufficient buffer of 32, and can transmute from uninit
//     //*         because all bytes are filled. This assumes a valid random_seed implementation.
//     unsafe {
//         sys::random_seed(ATOMIC_OP_REGISTER);
//         read_register_fixed_32(ATOMIC_OP_REGISTER)
//     }
// }
// pub fn sha256(value: &[u8]) -> Vec<u8> {
//     sha256_array(value).to_vec()
// }
// pub fn keccak256(value: &[u8]) -> Vec<u8> {
//     keccak256_array(value).to_vec()
// }
// pub fn keccak512(value: &[u8]) -> Vec<u8> {
//     keccak512_array(value).to_vec()
// }
// pub fn sha256_array(value: &[u8]) -> [u8; 32] {
//     //* SAFETY: sha256 syscall will always generate 32 bytes inside of the atomic op register
//     //*         so the read will have a sufficient buffer of 32, and can transmute from uninit
//     //*         because all bytes are filled. This assumes a valid sha256 implementation.
//     unsafe {
//         sys::sha256(value.len() as _, value.as_ptr() as _, ATOMIC_OP_REGISTER);
//         read_register_fixed_32(ATOMIC_OP_REGISTER)
//     }
// }
// pub fn keccak256_array(value: &[u8]) -> [u8; 32] {
//     //* SAFETY: keccak256 syscall will always generate 32 bytes inside of the atomic op register
//     //*         so the read will have a sufficient buffer of 32, and can transmute from uninit
//     //*         because all bytes are filled. This assumes a valid keccak256 implementation.
//     unsafe {
//         sys::keccak256(value.len() as _, value.as_ptr() as _, ATOMIC_OP_REGISTER);
//         read_register_fixed_32(ATOMIC_OP_REGISTER)
//     }
// }
// pub fn keccak512_array(value: &[u8]) -> [u8; 64] {
//     //* SAFETY: keccak512 syscall will always generate 64 bytes inside of the atomic op register
//     //*         so the read will have a sufficient buffer of 64, and can transmute from uninit
//     //*         because all bytes are filled. This assumes a valid keccak512 implementation.
//     unsafe {
//         sys::keccak512(value.len() as _, value.as_ptr() as _, ATOMIC_OP_REGISTER);
//         read_register_fixed_64(ATOMIC_OP_REGISTER)
//     }
// }
// pub fn ripemd160_array(value: &[u8]) -> [u8; 20] {
//     //* SAFETY: ripemd160 syscall will always generate 20 bytes inside of the atomic op register
//     //*         so the read will have a sufficient buffer of 20, and can transmute from uninit
//     //*         because all bytes are filled. This assumes a valid ripemd160 implementation.
//     unsafe {
//         sys::ripemd160(value.len() as _, value.as_ptr() as _, ATOMIC_OP_REGISTER);
//         read_register_fixed_20(ATOMIC_OP_REGISTER)
//     }
// }
// #[cfg(feature = "unstable")]
// pub fn ecrecover(
//     hash: &[u8],
//     signature: &[u8],
//     v: u8,
//     malleability_flag: bool,
// ) -> Option<[u8; 64]> {
//     unsafe {
//         let return_code = sys::ecrecover(
//             hash.len() as _,
//             hash.as_ptr() as _,
//             signature.len() as _,
//             signature.as_ptr() as _,
//             v as u64,
//             malleability_flag as u64,
//             ATOMIC_OP_REGISTER,
//         );
//         if return_code == 0 {
//             None
//         } else {
//             Some(read_register_fixed_64(ATOMIC_OP_REGISTER))
//         }
//     }
// }
// pub fn ed25519_verify(signature: &[u8; 64], message: &[u8], public_key: &[u8; 32]) -> bool {
//     unsafe {
//         sys::ed25519_verify(
//             signature.len() as _,
//             signature.as_ptr() as _,
//             message.len() as _,
//             message.as_ptr() as _,
//             public_key.len() as _,
//             public_key.as_ptr() as _,
//         ) == 1
//     }
// }
// pub fn alt_bn128_g1_multiexp(value: &[u8]) -> Vec<u8> {
//     unsafe {
//         sys::alt_bn128_g1_multiexp(value.len() as _, value.as_ptr() as _, ATOMIC_OP_REGISTER);
//     };
//     match read_register(ATOMIC_OP_REGISTER) {
//         Some(result) => result,
//         None => panic_str(REGISTER_EXPECTED_ERR),
//     }
// }
// pub fn alt_bn128_g1_sum(value: &[u8]) -> Vec<u8> {
//     unsafe {
//         sys::alt_bn128_g1_sum(value.len() as _, value.as_ptr() as _, ATOMIC_OP_REGISTER);
//     };
//     match read_register(ATOMIC_OP_REGISTER) {
//         Some(result) => result,
//         None => panic_str(REGISTER_EXPECTED_ERR),
//     }
// }
// pub fn alt_bn128_pairing_check(value: &[u8]) -> bool {
//     unsafe { sys::alt_bn128_pairing_check(value.len() as _, value.as_ptr() as _) == 1 }
// }
// ###############
// # Storage API #
// ###############

// pub fn storage_write(key: &[u8], value: &[u8]) -> bool {
//     match unsafe {
//         sys::storage_write(
//             key.len() as _,
//             key.as_ptr() as _,
//             value.len() as _,
//             value.as_ptr() as _,
//             EVICTED_REGISTER,
//         )
//     } {
//         0 => false,
//         1 => true,
//         _ => abort(),
//     }
// }
// pub fn storage_read(key: &[u8]) -> Option<Vec<u8>> {
//     match unsafe { sys::storage_read(key.len() as _, key.as_ptr() as _, ATOMIC_OP_REGISTER) } {
//         0 => None,
//         1 => Some(expect_register(read_register(ATOMIC_OP_REGISTER))),
//         _ => abort(),
//     }
// }
// pub fn storage_remove(key: &[u8]) -> bool {
//     match unsafe { sys::storage_remove(key.len() as _, key.as_ptr() as _, EVICTED_REGISTER) } {
//         0 => false,
//         1 => true,
//         _ => abort(),
//     }
// }
// pub fn storage_get_evicted() -> Option<Vec<u8>> {
//     read_register(EVICTED_REGISTER)
// }
// pub fn storage_has_key(key: &[u8]) -> bool {
//     match unsafe { sys::storage_has_key(key.len() as _, key.as_ptr() as _) } {
//         0 => false,
//         1 => true,
//         _ => abort(),
//     }
// }

// ############################################
// # Saving and loading of the contract state #
// ############################################
/// Load the state of the given object.
// pub fn state_read<T: borsh::BorshDeserialize>() -> Option<T> {
//     storage_read(STATE_KEY).map(|data| {
//         T::try_from_slice(&data)
//             .unwrap_or_else(|_| panic_str("Cannot deserialize the contract state."))
//     })
// }
// pub fn state_write<T: borsh::BorshSerialize>(state: &T) {
//     let data = match borsh::to_vec(state) {
//         Ok(serialized) => serialized,
//         Err(_) => panic_str("Cannot serialize the contract state."),
//     };
//     storage_write(STATE_KEY, &data);
// }
// pub fn state_exists() -> bool {
//     storage_has_key(STATE_KEY)
// }
// #####################################
// # Parameters exposed by the runtime #
// #####################################
// pub fn storage_byte_cost() -> NearToken {
//     NearToken::from_yoctonear(10_000_000_000_000_000_000u128)
// }
