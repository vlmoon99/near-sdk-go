package main

import (
	"encoding/hex"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/types"
)

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
	env.ContractValueReturn([]byte("1"))
}

//This is integration tests which will be executred on the testnet , u need to call it step by step in order to reproduce

// Registers API

//go:export TestWriteReadRegisterSafe
func TestWriteReadRegisterSafe() {
	data := []byte("test data")
	registerId := uint64(1)

	env.WriteRegisterSafe(registerId, data)

	readData, err := env.ReadRegisterSafe(registerId)
	if err != nil {
		env.PanicStr("Failed to read from register:  " + err.Error())
	}

	if string(readData) != "test data" {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

// Registers API

// Storage API

//go:export TestStorageWrite
func TestStorageWrite() {
	key := []byte("testKey")
	value := []byte("testValue")

	success, err := env.StorageWrite(key, value)
	if err != nil {
		env.PanicStr("Failed to write to storage: " + err.Error())
	}

	if !success {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

//go:export TestStorageRead
func TestStorageRead() {
	key := []byte("testKey")
	value := []byte("testValue")

	// Read from storage
	readValue, err := env.StorageRead(key)
	if err != nil {
		env.PanicStr("Failed to read from storage: " + err.Error())
	}

	if string(readValue) != string(value) {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

//go:export TestStorageHasKey
func TestStorageHasKey() {
	key := []byte("testKey")

	exists, err := env.StorageHasKey(key)
	if err != nil {
		env.PanicStr("Failed to check key existence: " + err.Error())
	}

	if !exists {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

//go:export TestStorageRemove
func TestStorageRemove() {
	key := []byte("testKey")

	success, err := env.StorageRemove(key)
	if err != nil {
		env.PanicStr("Failed to remove from storage: " + err.Error())
	}

	if !success {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

//go:export TestStateWrite
func TestStateWrite() {
	data := []byte("stateData")

	err := env.StateWrite(data)
	if err != nil {
		env.PanicStr("Failed to write state: " + err.Error())
	}

	// Verify the state was written correctly
	readData, err := env.StateRead()
	if err != nil {
		env.PanicStr("Failed to read state: " + err.Error())
	}

	if string(readData) != string(data) {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

//go:export TestStateRead
func TestStateRead() {
	data, err := env.StateRead()
	if err != nil {
		env.PanicStr("Failed to read state: " + err.Error())
	}

	if data == nil {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

//go:export TestStateExists
func TestStateExists() {
	exists := env.StateExists()

	if !exists {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

// NOT Working for some Reason

// //go:export TestStorageGetEvicted
// func TestStorageGetEvicted() {
// 	value, err := env.StorageGetEvicted()
// 	if err != nil {
// 		env.PanicStr("Failed to get evicted value: " + err.Error())
// 	}

// 	if value == nil {
// 		env.ContractValueReturn([]byte("0"))
// 	} else {
// 		env.ContractValueReturn([]byte("1"))
// 	}
// }

// Storage API

// Context API

//go:export TestGetCurrentAccountId
func TestGetCurrentAccountId() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	env.LogString("Current account ID: " + accountId)
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetSignerAccountID
func TestGetSignerAccountID() {
	signerId, err := env.GetSignerAccountID()
	if err != nil || signerId == "" {
		env.PanicStr("Failed to get signer account ID: " + err.Error())
	}
	env.LogString("Signer account ID: " + signerId)
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetSignerAccountPK
func TestGetSignerAccountPK() {
	signerPK, err := env.GetSignerAccountPK()
	if err != nil || signerPK == nil {
		env.PanicStr("Failed to get signer account PK: " + err.Error())
	}
	env.LogString("Signer account PK: " + hex.EncodeToString(signerPK))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetPredecessorAccountID
func TestGetPredecessorAccountID() {
	predecessorId, err := env.GetPredecessorAccountID()
	if err != nil || predecessorId == "" {
		env.PanicStr("Failed to get predecessor account ID: " + err.Error())
	}
	env.LogString("Predecessor account ID: " + predecessorId)
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetCurrentBlockHeight
func TestGetCurrentBlockHeight() {
	env.GetCurrentBlockHeight()
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetBlockTimeMs
func TestGetBlockTimeMs() {
	env.GetBlockTimeMs()
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetEpochHeight
func TestGetEpochHeight() {
	env.GetEpochHeight()
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetStorageUsage
func TestGetStorageUsage() {
	env.GetStorageUsage()
	env.ContractValueReturn([]byte("1"))
}

//go:export TestContractInputRawBytes
func TestContractInputRawBytes() {
	options := types.ContractInputOptions{IsRawBytes: true}
	input, detectedType, err := env.ContractInput(options)
	if err != nil {
		env.PanicStr("Failed to get contract input: " + err.Error())
	}
	env.LogString("Contract input (raw bytes): " + string(input))
	env.LogString("Detected input type: " + detectedType)
	env.ContractValueReturn([]byte("1"))
}

//go:export TestContractInputJSON
func TestContractInputJSON() {
	options := types.ContractInputOptions{IsRawBytes: false}
	input, detectedType, err := env.ContractInput(options)
	if err != nil {
		env.PanicStr("Failed to get contract input: " + err.Error())
	}
	env.LogString("Contract input (JSON): " + string(input))
	env.LogString("Detected input type: " + detectedType)
	env.ContractValueReturn([]byte("1"))
}

// Context API
