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
	promiseCount := env.PromiseResultsCount()
	if promiseCount != 0 {
		env.LogString("Promise Count is : " + fmt.Sprintf("%d", promiseCount))
		result, err := env.PromiseResult(0)
		if err != nil {
			env.LogString("Promise result err : " + err.Error())
		} else {
			env.LogString("Promise result at index: " + string(result))
		}
	}

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

// Promises API

//go:export TestPromiseCreate
func TestPromiseCreate() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}

	arguments := []byte("")
	accountID := []byte(accountId)
	functionName := []byte("TestLogStringUtf16")
	amount := types.Uint128{Hi: 0, Lo: 0}
	gas := uint64(3000000000)

	promiseIdx := env.PromiseCreate(accountID, functionName, arguments, amount, gas)

	env.LogString("Promise created with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseThen
func TestPromiseThen() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	arguments1 := []byte("")
	accountID1 := []byte(accountId)
	functionName1 := []byte("TestLogStringUtf16")
	amount1 := types.Uint128{Hi: 0, Lo: 0}
	gas1 := uint64(3000000000)

	promiseIdx := env.PromiseCreate(accountID1, functionName1, arguments1, amount1, gas1)

	arguments2 := []byte("")
	accountID2 := []byte(accountId)
	functionName2 := []byte("TestLogStringUtf8")
	amount2 := types.Uint128{Hi: 0, Lo: 0}
	gas2 := uint64(3000000000)

	chainedPromiseIdx := env.PromiseThen(promiseIdx, accountID2, functionName2, arguments2, amount2, gas2)
	env.LogString("Chained promise created with index: " + fmt.Sprintf("%d", chainedPromiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseAnd
func TestPromiseAnd() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	arguments1 := []byte("")
	accountID1 := []byte(accountId)
	functionName1 := []byte("TestLogString")
	amount1 := types.Uint128{Hi: 0, Lo: 0}
	gas1 := uint64(3000000000)

	promiseIdx1 := env.PromiseCreate(accountID1, functionName1, arguments1, amount1, gas1)

	arguments2 := []byte("")
	accountID2 := []byte(accountId)
	functionName2 := []byte("TestLogStringUtf8")
	amount2 := types.Uint128{Hi: 0, Lo: 0}
	gas2 := uint64(3000000000)

	promiseIdx2 := env.PromiseCreate(accountID2, functionName2, arguments2, amount2, gas2)

	promiseIndices := []uint64{promiseIdx1, promiseIdx2}
	chainedPromiseIdx := env.PromiseAnd(promiseIndices)
	env.LogString("Chained promise created with index: " + fmt.Sprintf("%d", chainedPromiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchCreate
func TestPromiseBatchCreate() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte(accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)
	env.LogString("Promise batch created with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchThen
func TestPromiseBatchThen() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	arguments := []byte("")
	accountID := []byte(accountId)
	functionName := []byte("TestLogString")
	amount := types.Uint128{Hi: 0, Lo: 0}
	gas := uint64(3000000000)

	promiseIdx := env.PromiseCreate(accountID, functionName, arguments, amount, gas)

	accountID2 := []byte(accountId)
	chainedPromiseIdx := env.PromiseBatchThen(promiseIdx, accountID2)
	env.LogString("Chained promise batch created with index: " + fmt.Sprintf("%d", chainedPromiseIdx))
	env.ContractValueReturn([]byte("1"))
}

// Promises API

// Promises API Action

//go:export TestPromiseBatchActionCreateAccount
func TestPromiseBatchActionCreateAccount() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	env.LogString("accountId : " + accountId)
	accountID := []byte("app." + accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)
	env.PromiseBatchActionCreateAccount(promiseIdx)
	env.LogString("Promise batch action create account with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchActionDeployContract
func TestPromiseBatchActionDeployContract() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte("app." + accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)
	contractBytes := []byte("sample contract bytes")
	env.PromiseBatchActionDeployContract(promiseIdx, contractBytes)
	env.LogString("Promise batch action deploy contract with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchActionFunctionCall
func TestPromiseBatchActionFunctionCall() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte(accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)
	functionName := []byte("TestLogStringUtf8")
	arguments := []byte("{}")
	amount := types.Uint128{Hi: 0, Lo: 0}
	gas := uint64(3000000000)

	env.PromiseBatchActionFunctionCall(promiseIdx, functionName, arguments, amount, gas)
	env.LogString("Promise batch action function call with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchActionFunctionCallWeight
func TestPromiseBatchActionFunctionCallWeight() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte(accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)
	functionName := []byte("TestLogStringUtf8")
	arguments := []byte("{}")
	amount := types.Uint128{Hi: 0, Lo: 0}
	gas := uint64(3000000000)
	weight := uint64(1)

	env.PromiseBatchActionFunctionCallWeight(promiseIdx, functionName, arguments, amount, gas, weight)
	env.LogString("Promise batch action function call with weight with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchActionTransfer
func TestPromiseBatchActionTransfer() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}

	accountID := []byte("app." + accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)

	// amount, _ := types.U128FromString("1000000000000000000000000") // 1 Near
	amount, _ := types.U128FromString("10000000000000000000000") // 0.01 Near

	env.PromiseBatchActionTransfer(promiseIdx, amount)
	env.LogString("Promise batch action transfer with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchActionStake
func TestPromiseBatchActionStake() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte(accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)
	// amount, _ := types.U128FromString("1000000000000000000000000") // 1 Near
	amount, _ := types.U128FromString("10000000000000000000000") // 0.01 Near
	signerPK, err := env.GetSignerAccountPK()
	if err != nil || signerPK == nil {
		env.PanicStr("Failed to get signer account PK: " + err.Error())
	}
	env.LogString("signerPK: " + hex.EncodeToString(signerPK))

	env.PromiseBatchActionStake(promiseIdx, amount, signerPK)
	env.LogString("Promise batch action stake with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//TODO : add crypto for createing ed25519 keys
// //go:export TestPromiseBatchActionAddKeyWithFullAccess
// func TestPromiseBatchActionAddKeyWithFullAccess() {
// 	accountId, err := env.GetCurrentAccountId()
// 	if err != nil || accountId == "" {
// 		env.PanicStr("Failed to get current account ID: " + err.Error())
// 	}
// 	accountID := []byte(accountId)

// 	promiseIdx := env.PromiseBatchCreate(accountID)
// 	publicKey := []byte("sample_public_key")
// 	nonce := uint64(0)

// 	env.PromiseBatchActionAddKeyWithFullAccess(promiseIdx, publicKey, nonce)
// 	env.LogString("Promise batch action add key with full access with index: " + fmt.Sprintf("%d", promiseIdx))
// 	env.ContractValueReturn([]byte("1"))
// }

//TODO : add crypto for createing ed25519 keys
// //go:export TestPromiseBatchActionAddKeyWithFunctionCall
// func TestPromiseBatchActionAddKeyWithFunctionCall() {
// 	accountId, err := env.GetCurrentAccountId()
// 	if err != nil || accountId == "" {
// 		env.PanicStr("Failed to get current account ID: " + err.Error())
// 	}
// 	accountID := []byte(accountId)

// 	promiseIdx := env.PromiseBatchCreate(accountID)
// 	publicKey := []byte("sample_public_key")
// 	nonce := uint64(0)
// 	amount := types.Uint128{Hi: 0, Lo: 1000}
// 	receiverId := []byte("receiver.near")
// 	functionName := []byte("TestLogStringUtf8")

// 	env.PromiseBatchActionAddKeyWithFunctionCall(promiseIdx, publicKey, nonce, amount, receiverId, functionName)
// 	env.LogString("Promise batch action add key with function call with index: " + fmt.Sprintf("%d", promiseIdx))
// 	env.ContractValueReturn([]byte("1"))
// }

//TODO : add crypto for createing ed25519 keys
// //go:export TestPromiseBatchActionDeleteKey
// func TestPromiseBatchActionDeleteKey() {
// 	accountId, err := env.GetCurrentAccountId()
// 	if err != nil || accountId == "" {
// 		env.PanicStr("Failed to get current account ID: " + err.Error())
// 	}
// 	accountID := []byte(accountId)

// 	promiseIdx := env.PromiseBatchCreate(accountID)
// 	publicKey := []byte("sample_public_key")

// 	env.PromiseBatchActionDeleteKey(promiseIdx, publicKey)
// 	env.LogString("Promise batch action delete key with index: " + fmt.Sprintf("%d", promiseIdx))
// 	env.ContractValueReturn([]byte("1"))
// }

//go:export TestPromiseBatchActionDeleteAccount
func TestPromiseBatchActionDeleteAccount() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte("app." + accountId)
	beneficiaryId := []byte("beneficiary.near")

	promiseIdx := env.PromiseBatchCreate(accountID)
	env.PromiseBatchActionDeleteAccount(promiseIdx, beneficiaryId)
	env.LogString("Promise batch action delete account with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//TODO : read near sdk rs about it
// //go:export TestPromiseYieldCreate
// func TestPromiseYieldCreate() {
// 	// accountId, err := env.GetCurrentAccountId()
// 	// if err != nil || accountId == "" {
// 	// 	env.PanicStr("Failed to get current account ID: " + err.Error())
// 	// }
// 	// arguments1 := []byte("")
// 	// accountID1 := []byte(accountId)
// 	// functionName1 := []byte("TestLogString")
// 	// amount1 := types.Uint128{Hi: 0, Lo: 0}
// 	// gas1 := uint64(3000000000)

// 	// env.PromiseCreate(accountID1, functionName1, arguments1, amount1, gas1)

// 	functionName := []byte("TestLogStringUtf8")
// 	arguments := []byte("{}")
// 	gas := uint64(3000000000)
// 	gasWeight := uint64(1)

// 	promiseIdx := env.PromiseYieldCreate(functionName, arguments, gas, gasWeight)
// 	env.LogString("Promise yield create with index: " + fmt.Sprintf("%d", promiseIdx))
// 	env.ContractValueReturn([]byte("1"))
// }

//TODO : read near sdk rs about it
// //go:export TestPromiseYieldResume
// func TestPromiseYieldResume() {
// 	data := []byte("sample data")
// 	payload := []byte("sample payload")

// 	result := env.PromiseYieldResume(data, payload)
// 	env.LogString("Promise yield resume result: " + fmt.Sprintf("%d", result))
// 	env.ContractValueReturn([]byte("1"))
// }

// Promises API Action

// Promise API Results

//go:export TestPromiseResultsCount
func TestPromiseResultsCount() {
	count := env.PromiseResultsCount()
	env.LogString("Promise results count: " + fmt.Sprintf("%d", count))
	env.ContractValueReturn([]byte("1"))
}

// TODO : Learn about this flow execution
//
//go:export TestPromiseResult
func TestPromiseResult() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}

	prepaidGas := env.GetPrepaidGas()
	env.LogString(fmt.Sprintf("Prepaid gas: %d", prepaidGas.Inner))

	arguments1 := []byte("")
	accountID1 := []byte(accountId)
	functionName1 := []byte("TestContractValueReturn")
	amount1 := types.Uint128{Hi: 0, Lo: 0}
	gas1 := prepaidGas.Inner / 3

	promiseIdx := env.PromiseCreate(accountID1, functionName1, arguments1, amount1, gas1)

	arguments2 := []byte("")
	accountID2 := []byte(accountId)
	functionName2 := []byte("TestPromiseResultsCount")
	amount2 := types.Uint128{Hi: 0, Lo: 0}
	gas2 := prepaidGas.Inner / 3

	chainedPromiseIdx := env.PromiseThen(promiseIdx, accountID2, functionName2, arguments2, amount2, gas2)
	env.LogString("Chained promise created with index: " + fmt.Sprintf("%d", chainedPromiseIdx))

	env.ContractValueReturn([]byte("1"))
}

// TODO : Learn about this flow execution
//
//go:export TestPromiseReturn
func TestPromiseReturn() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	prepaidGas := env.GetPrepaidGas()
	env.LogString(fmt.Sprintf("Prepaid gas: %d", prepaidGas.Inner))

	arguments1 := []byte("")
	accountID1 := []byte(accountId)
	functionName1 := []byte("TestContractValueReturn")
	amount1 := types.Uint128{Hi: 0, Lo: 0}
	gas1 := prepaidGas.Inner / 3

	promiseIdx := env.PromiseCreate(accountID1, functionName1, arguments1, amount1, gas1)

	arguments2 := []byte("")
	accountID2 := []byte(accountId)
	functionName2 := []byte("TestPromiseResultsCount")
	amount2 := types.Uint128{Hi: 0, Lo: 0}
	gas2 := prepaidGas.Inner / 3

	chainedPromiseIdx := env.PromiseThen(promiseIdx, accountID2, functionName2, arguments2, amount2, gas2)
	env.LogString("Chained promise created with index: " + fmt.Sprintf("%d", chainedPromiseIdx))

	// Return the promise result
	env.PromiseReturn(chainedPromiseIdx)

	env.LogString("Promise returned with ID: " + fmt.Sprintf("%d", chainedPromiseIdx))
	env.ContractValueReturn([]byte("1"))
}

// Promise API Results
