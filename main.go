package main

import (
	"encoding/hex"
	"fmt"

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
	blockHeight := env.GetCurrentBlockHeight()
	env.LogString("Current block height: " + fmt.Sprintf("%d", blockHeight))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetBlockTimeMs
func TestGetBlockTimeMs() {
	blockTimeMs := env.GetBlockTimeMs()
	env.LogString("Block time in ms: " + fmt.Sprintf("%d", blockTimeMs))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetEpochHeight
func TestGetEpochHeight() {
	epochHeight := env.GetEpochHeight()
	env.LogString("Epoch height: " + fmt.Sprintf("%d", epochHeight))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetStorageUsage
func TestGetStorageUsage() {
	storageUsage := env.GetStorageUsage()
	env.LogString("Storage usage: " + fmt.Sprintf("%d", storageUsage))
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

// Economics API

//go:export TestGetAccountBalance
func TestGetAccountBalance() {
	accountBalance, err := env.GetAccountBalance()
	if err != nil {
		env.PanicStr("Failed to get account balance: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Account balance: Hi: %d, Lo: %d", accountBalance.Hi, accountBalance.Lo))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetAccountLockedBalance
func TestGetAccountLockedBalance() {
	lockedBalance, err := env.GetAccountLockedBalance()
	if err != nil {
		env.PanicStr("Failed to get account locked balance: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Account locked balance: Hi: %d, Lo: %d", lockedBalance.Hi, lockedBalance.Lo))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetAttachedDeposit
func TestGetAttachedDeposit() {
	attachedDeposit, err := env.GetAttachedDepoist()
	if err != nil {
		env.PanicStr("Failed to get attached deposit: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Attached deposit: Hi: %d, Lo: %d", attachedDeposit.Hi, attachedDeposit.Lo))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetPrepaidGas
func TestGetPrepaidGas() {
	prepaidGas := env.GetPrepaidGas()
	env.LogString(fmt.Sprintf("Prepaid gas: %d", prepaidGas.Inner))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestGetUsedGas
func TestGetUsedGas() {
	usedGas := env.GetUsedGas()
	env.LogString(fmt.Sprintf("Used gas: %d", usedGas.Inner))
	env.ContractValueReturn([]byte("1"))
}

// Economics API

// Math API

//go:export TestGetRandomSeed
func TestGetRandomSeed() {
	seed, err := env.GetRandomSeed()
	if err != nil {
		env.PanicStr("Failed to get random seed: " + err.Error())
	}
	env.LogString("Random seed: " + fmt.Sprintf("%x", seed))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestSha256Hash
func TestSha256Hash() {
	data := []byte("test data")
	hash, err := env.Sha256Hash(data)
	if err != nil {
		env.PanicStr("Failed to get SHA-256 hash: " + err.Error())
	}
	env.LogString("SHA-256 hash: " + fmt.Sprintf("%x", hash))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestKeccak256Hash
func TestKeccak256Hash() {
	data := []byte("test data")
	hash, err := env.Keccak256Hash(data)
	if err != nil {
		env.PanicStr("Failed to get Keccak-256 hash: " + err.Error())
	}
	env.LogString("Keccak-256 hash: " + fmt.Sprintf("%x", hash))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestKeccak512Hash
func TestKeccak512Hash() {
	data := []byte("test data")
	hash, err := env.Keccak512Hash(data)
	if err != nil {
		env.PanicStr("Failed to get Keccak-512 hash: " + err.Error())
	}
	env.LogString("Keccak-512 hash: " + fmt.Sprintf("%x", hash))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestRipemd160Hash
func TestRipemd160Hash() {
	data := []byte("test data")
	hash, err := env.Ripemd160Hash(data)
	if err != nil {
		env.PanicStr("Failed to get RIPEMD-160 hash: " + err.Error())
	}
	env.LogString("RIPEMD-160 hash: " + fmt.Sprintf("%x", hash))
	env.ContractValueReturn([]byte("1"))
}

//TODO : add test cases
// //go:export TestEcrecoverPubKey
// func TestEcrecoverPubKey() {
// 	hash := []byte{0x1c, 0x29, 0x37, 0xbf, 0x5e, 0x3b, 0x5b, 0xd7, 0x0a, 0x47, 0x29, 0x28, 0x62, 0x10, 0x33, 0x6d, 0x7d, 0x6b, 0x5b, 0x6c, 0x61, 0x17, 0x2a, 0xe2, 0x69, 0x1a, 0x51, 0x37, 0x3e, 0x7d, 0x1c, 0xf9}
// 	signature := []byte{0x1f, 0x6e, 0x3a, 0x58, 0x91, 0x29, 0x55, 0x8a, 0xf6, 0x88, 0x5f, 0x51, 0xe8, 0x5b, 0x5d, 0x59, 0x42, 0x9b, 0x51, 0x2f, 0x19, 0x71, 0x5a, 0xa1, 0xd9, 0x6e, 0x5a, 0xd2, 0xbe, 0x59, 0xad, 0x56, 0x1f, 0x6e, 0x3a, 0x58, 0x91, 0x29, 0x55, 0x8a, 0xf6, 0x88, 0x5f, 0x51, 0xe8, 0x5b, 0x5d, 0x59, 0x42, 0x9b, 0x51, 0x2f, 0x19, 0x71, 0x5a, 0xa1, 0xd9, 0x6e, 0x5a, 0xd2, 0xbe, 0x59, 0xad, 0x56}
// 	v := byte(0x1b)
// 	malleabilityFlag := false
// 	pubKey, err := env.EcrecoverPubKey(hash, signature, v, malleabilityFlag)
// 	if err != nil {
// 		env.PanicStr("Failed to recover public key: " + err.Error())
// 	}
// 	env.LogString("Recovered public key: " + fmt.Sprintf("%x", pubKey))
// 	env.ContractValueReturn([]byte("1"))
// }

// //go:export TestEd25519VerifySig
// func TestEd25519VerifySig() {
// 	signature := [64]byte{0x1f, 0x6e, 0x3a, 0x58, 0x91, 0x29, 0x55, 0x8a, 0xf6, 0x88, 0x5f, 0x51, 0xe8, 0x5b, 0x5d, 0x59, 0x42, 0x9b, 0x51, 0x2f, 0x19, 0x71, 0x5a, 0xa1, 0xd9, 0x6e, 0x5a, 0xd2, 0xbe, 0x59, 0xad, 0x56, 0x1f, 0x6e, 0x3a, 0x58, 0x91, 0x29, 0x55, 0x8a, 0xf6, 0x88, 0x5f, 0x51, 0xe8, 0x5b, 0x5d, 0x59, 0x42, 0x9b, 0x51, 0x2f, 0x19, 0x71, 0x5a, 0xa1, 0xd9, 0x6e, 0x5a, 0xd2, 0xbe, 0x59, 0xad, 0x56}
// 	message := []byte("test message")
// 	publicKey := [32]byte{0x9b, 0x51, 0x2f, 0x19, 0x71, 0x5a, 0xa1, 0xd9, 0x6e, 0x5a, 0xd2, 0xbe, 0x59, 0xad, 0x56, 0x1f, 0x6e, 0x3a, 0x58, 0x91, 0x29, 0x55, 0x8a, 0xf6, 0x88, 0x5f, 0x51, 0xe8, 0x5b, 0x5d, 0x59, 0x42}
// 	valid := env.Ed25519VerifySig(signature, message, publicKey)
// 	env.LogString("Ed25519 signature valid: " + fmt.Sprintf("%t", valid))
// 	env.ContractValueReturn([]byte("1"))
// }

// //go:export TestAltBn128G1MultiExp
// func TestAltBn128G1MultiExp() {
// 	data := types.Uint128{Hi: 2, Lo: 1}
// 	result, err := env.AltBn128G1MultiExp(data.ToLE())
// 	if err != nil {
// 		env.PanicStr("Failed to perform AltBn128G1MultiExp: " + err.Error())
// 	}
// 	env.LogString("AltBn128G1MultiExp result: " + fmt.Sprintf("%x", result))
// 	env.ContractValueReturn([]byte("1"))
// }

// //go:export TestAltBn128G1Sum
// func TestAltBn128G1Sum() {
// 	data := types.Uint128{Hi: 2, Lo: 1}
// 	result, err := env.AltBn128G1Sum(data.ToLE())
// 	if err != nil {
// 		env.PanicStr("Failed to perform AltBn128G1Sum: " + err.Error())
// 	}
// 	env.LogString("AltBn128G1Sum result: " + fmt.Sprintf("%x", result))
// 	env.ContractValueReturn([]byte("1"))
// }

// //go:export TestAltBn128PairingCheck
// func TestAltBn128PairingCheck() {
// 	data := types.Uint128{Hi: 2, Lo: 1}
// 	valid := env.AltBn128PairingCheck(data.ToLE())
// 	env.LogString("AltBn128PairingCheck valid: " + fmt.Sprintf("%t", valid))
// 	env.ContractValueReturn([]byte("1"))
// }

// Math API

// Validator API

//go:export TestValidatorStakeAmount
func TestValidatorStakeAmount() {
	accountID := []byte("testaccount")
	stakeAmount, err := env.ValidatorStakeAmount(accountID)
	if err != nil {
		env.PanicStr("Failed to get validator stake amount: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Validator stake amount: Hi: %d, Lo: %d", stakeAmount.Hi, stakeAmount.Lo))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestValidatorTotalStakeAmount
func TestValidatorTotalStakeAmount() {
	totalStakeAmount, err := env.ValidatorTotalStakeAmount()
	if err != nil {
		env.PanicStr("Failed to get validator total stake amount: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Validator total stake amount: Hi: %d, Lo: %d", totalStakeAmount.Hi, totalStakeAmount.Lo))
	env.ContractValueReturn([]byte("1"))
}

// Validator API

// Miscellaneous API

//go:export TestContractValueReturn
func TestContractValueReturn() {
	inputBytes := []byte("test value")
	env.ContractValueReturn(inputBytes)
	env.LogString("Contract value returned: " + string(inputBytes))
	env.ContractValueReturn([]byte("1"))
}

// It's unnesessary to test, but if u want unncoment and tests it
//
// //go:export TestPanicStr
// func TestPanicStr() {
// 	env.PanicStr("Test panic")
// }

// It's unnesessary to test, but if u want unncoment and tests it
//
// //go:export TestAbortExecution
// func TestAbortExecution() {
// 	env.AbortExecution()
// }

//go:export TestLogString
func TestLogString() {
	input := "test log"
	env.LogString(input)
	env.LogString("Logged string: " + input)
	env.ContractValueReturn([]byte("1"))
}

//go:export TestLogStringUtf8
func TestLogStringUtf8() {
	input := []byte("test log utf8")
	env.LogStringUtf8(input)
	env.LogString("Logged string UTF-8: " + string(input))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestLogStringUtf16
func TestLogStringUtf16() {
	input := []byte("test log utf16")
	env.LogStringUtf16(input)
	env.LogString("Logged string UTF-16: " + string(input))
	env.ContractValueReturn([]byte("1"))
}

// Miscellaneous API
