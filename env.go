package main

import (
	"errors"
	"math"
	"unsafe"
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
		return nil, errors.New("expected data in register, but found none")
	}
	return data, nil
}

func ReadRegisterSafe(registerId uint64) ([]byte, error) {
	length := RegisterLen(registerId)
	if length == 0 {
		return []byte{}, errors.New("expected data in register, but found none")
	}

	buffer := make([]byte, length)

	ptr := uint64(uintptr(unsafe.Pointer(&buffer[0])))

	ReadRegister(registerId, ptr)

	return buffer, nil
}

func WriteRegisterSafe(registerId uint64, data []byte) {
	if len(data) == 0 {
		return
	}

	ptr := uint64(uintptr(unsafe.Pointer(&data[0])))

	WriteRegister(registerId, uint64(len(data)), ptr)
}

func SmartContractLog(input string) {
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
		SmartContractLog("Error in GetCurrentAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { SignerAccountId(registerID) })
	if err != nil {
		SmartContractLog("Error in GetSignerAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSignerAccountPK() ([]byte, error) {
	data, err := methodIntoRegister(func(registerID uint64) { SignerAccountPk(registerID) })
	if err != nil {
		SmartContractLog("Error in GetSignerAccountPK: " + err.Error())
		return nil, err
	}

	return data, nil
}

func GetPredecessorAccountID() (string, error) {
	data, err := methodIntoRegister(func(registerID uint64) { PredecessorAccountId(registerID) })
	if err != nil {
		SmartContractLog("Error in GetPredecessorAccountID: " + err.Error())
		return "", err
	}

	return assertValidAccountId(data)
}

func GetSmartContractInput() ([]byte, error) {
	Input(AtomicOpRegister)
	return ReadRegisterSafe(AtomicOpRegister)
}
