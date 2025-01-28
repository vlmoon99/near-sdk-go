package main

import (
	"math"
	"unsafe"
)

// Error message when a register is expected to have data but does not.
const RegisterExpectedErr = "Register was expected to have data because we just wrote it into it."

// Register used internally for atomic operations. This register is safe to use by the user,
// since it only needs to be untouched while methods of `Environment` execute, which is guaranteed
// as guest code is not parallel.
const AtomicOpRegister uint64 = math.MaxUint64 - 2

// Register used to record evicted values from the storage.
const EvictedRegister uint64 = math.MaxUint64 - 1

// Key used to store the state of the contract.
var StateKey = []byte("STATE")

// The minimum length of a valid account ID.
const MinAccountIDLen uint64 = 2

// The maximum length of a valid account ID.
const MaxAccountIDLen uint64 = 64

func EnvReadRegister(registerID uint64) ([]byte, error) {
	len := RegisterLen(registerID)
	if len == 0 {
		return nil, nil
	}

	buffer := make([]byte, len)

	ReadRegister(registerID, *(*uint64)(unsafe.Pointer(&buffer[0])))

	return buffer, nil
}

func SmartContractInput() ([]byte, error) {
	Input(AtomicOpRegister)
	return EnvReadRegister(AtomicOpRegister)
}
