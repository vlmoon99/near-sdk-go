package main

import (
	"bytes"
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

type EcrecoverTest struct {
	M   [32]byte
	V   uint8
	Sig [64]byte
	Mc  bool
	Res [64]byte
}

// TODO : tests unstable feature
//
//go:export TestEcrecoverPubKey
func TestEcrecoverPubKey() {
	data := EcrecoverTest{}

	m, _ := hex.DecodeString("ce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	copy(data.M[:], m)

	data.V = 1

	sig, _ := hex.DecodeString("90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc93")
	copy(data.Sig[:], sig)

	data.Mc = true

	res, _ := hex.DecodeString("e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	copy(data.Res[:], res)

	pubKey, err := env.EcrecoverPubKey(data.M[:], data.Sig[:], data.V, data.Mc)
	if err != nil {
		env.PanicStr("Failed to recover public key: " + err.Error())
	}
	env.LogString("Recovered public key: " + fmt.Sprintf("%x", pubKey))
	if hex.EncodeToString(pubKey) == "e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652" {
		env.ContractValueReturn([]byte("1"))
	} else {
		env.ContractValueReturn([]byte("0"))
	}
}

//go:export TestEd25519VerifySig
func TestEd25519VerifySig() {
	signature := [64]byte{145, 193, 203, 18, 114, 227, 14, 117, 33, 213, 121, 66, 130, 14, 25, 4, 36, 120, 46,
		142, 226, 215, 7, 66, 122, 112, 97, 30, 249, 135, 61, 165, 221, 249, 252, 23, 105, 40,
		56, 70, 31, 152, 236, 141, 154, 122, 207, 20, 75, 118, 79, 90, 168, 6, 221, 122, 213,
		29, 126, 196, 216, 104, 191, 6}
	message := []byte{107, 97, 106, 100, 108, 102, 107, 106, 97, 108, 107, 102, 106, 97, 107, 108, 102, 106,
		100, 107, 108, 97, 100, 106, 102, 107, 108, 106, 97, 100, 115, 107,
	}
	publicKey := [32]byte{32, 122, 6, 120, 146, 130, 30, 37, 215, 112, 241, 251, 160, 196, 124, 17, 255, 75, 129,
		62, 84, 22, 46, 206, 158, 184, 57, 224, 118, 35, 26, 182}
	valid := env.Ed25519VerifySig(signature, message, publicKey)
	env.LogString("Ed25519 signature valid: " + fmt.Sprintf("%t", valid))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestAltBn128G1MultiExp
func TestAltBn128G1MultiExp() {
	signature := []byte{16, 238, 91, 161, 241, 22, 172, 158, 138, 252, 202, 212, 136, 37, 110, 231, 118, 220,
		8, 45, 14, 153, 125, 217, 227, 87, 238, 238, 31, 138, 226, 8, 238, 185, 12, 155, 93,
		126, 144, 248, 200, 177, 46, 245, 40, 162, 169, 80, 150, 211, 157, 13, 10, 36, 44, 232,
		173, 32, 32, 115, 123, 2, 9, 47, 190, 148, 181, 91, 69, 6, 83, 40, 65, 222, 251, 70,
		81, 73, 60, 142, 130, 217, 176, 20, 69, 75, 40, 167, 41, 180, 244, 5, 142, 215, 135,
		35}

	expectedResult := []byte{150, 94, 159, 52, 239, 226, 181, 150, 77, 86, 90, 186, 102, 219, 243, 204, 36, 128,
		164, 209, 106, 6, 62, 124, 235, 104, 223, 195, 30, 204, 42, 20, 13, 158, 14, 197,
		133, 73, 43, 171, 28, 68, 82, 116, 244, 164, 36, 251, 244, 8, 234, 40, 118, 55,
		216, 187, 242, 39, 213, 160, 192, 184, 28, 23}

	result, err := env.AltBn128G1MultiExp(signature)
	if err != nil {
		env.PanicStr("Failed to perform AltBn128G1MultiExp: " + err.Error())
	}

	if bytes.Equal(expectedResult, result) {
		env.LogString("The result matches the expected result.")
		env.ContractValueReturn([]byte("1"))
	} else {
		env.LogString("The result does not match the expected result.")
		env.ContractValueReturn([]byte("0"))
	}
}

//go:export TestAltBn128G1Sum
func TestAltBn128G1Sum() {
	signature := []byte{0, 11, 49, 94, 29, 152, 111, 116, 138, 248, 2, 184, 8, 159, 80, 169, 45, 149, 48, 32,
		49, 37, 6, 133, 105, 171, 194, 120, 44, 195, 17, 180, 35, 137, 154, 4, 192, 211, 244,
		93, 200, 2, 44, 0, 64, 26, 108, 139, 147, 88, 235, 242, 23, 253, 52, 110, 236, 67, 99,
		176, 2, 186, 198, 228, 25}

	expectedResult := []byte{11, 49, 94, 29, 152, 111, 116, 138, 248, 2, 184, 8, 159, 80, 169, 45, 149, 48, 32,
		49, 37, 6, 133, 105, 171, 194, 120, 44, 195, 17, 180, 35, 137, 154, 4, 192, 211,
		244, 93, 200, 2, 44, 0, 64, 26, 108, 139, 147, 88, 235, 242, 23, 253, 52, 110, 236,
		67, 99, 176, 2, 186, 198, 228, 25}

	result, err := env.AltBn128G1Sum(signature)
	if err != nil {
		env.PanicStr("Failed to perform AltBn128G1Sum: " + err.Error())
	}
	if bytes.Equal(expectedResult, result) {
		env.LogString("The result matches the expected result.")
		env.ContractValueReturn([]byte("1"))
	} else {
		env.LogString("The result does not match the expected result.")
		env.ContractValueReturn([]byte("0"))
	}

}

//go:export TestAltBn128PairingCheck
func TestAltBn128PairingCheck() {
	signature := []byte{117, 10, 217, 99, 113, 78, 234, 67, 183, 90, 26, 58, 200, 86, 195, 123, 42, 184, 213,
		88, 224, 248, 18, 200, 108, 6, 181, 6, 28, 17, 99, 7, 36, 134, 53, 115, 192, 180, 3,
		113, 76, 227, 174, 147, 50, 174, 79, 74, 151, 195, 172, 10, 211, 210, 26, 92, 117, 246,
		65, 237, 168, 104, 16, 4, 1, 26, 3, 219, 6, 13, 193, 115, 77, 230, 27, 13, 242, 214,
		195, 9, 213, 99, 135, 12, 160, 202, 114, 135, 175, 42, 116, 172, 79, 234, 26, 41, 212,
		111, 192, 129, 124, 112, 57, 107, 38, 244, 230, 222, 240, 36, 65, 238, 133, 188, 19,
		43, 148, 59, 205, 40, 161, 179, 173, 228, 88, 169, 231, 29, 17, 67, 163, 51, 165, 187,
		101, 44, 250, 24, 68, 101, 92, 128, 203, 190, 51, 85, 9, 43, 58, 136, 68, 180, 92, 110,
		185, 168, 107, 129, 45, 30, 187, 22, 100, 17, 75, 93, 216, 125, 23, 212, 11, 186, 199,
		204, 1, 140, 133, 11, 82, 44, 65, 222, 20, 26, 48, 26, 132, 220, 25, 213, 93, 25, 79,
		176, 4, 149, 151, 243, 11, 131, 253, 233, 121, 38, 222, 15, 118, 117, 200, 214, 175,
		233, 130, 181, 193, 167, 255, 153, 169, 240, 207, 235, 28, 31, 83, 74, 69, 179, 6, 150,
		72, 67, 74, 166, 130, 83, 82, 115, 123, 111, 208, 221, 64, 43, 237, 213, 186, 235, 7,
		56, 251, 179, 95, 233, 159, 23, 109, 173, 85, 103, 8, 165, 235, 226, 218, 79, 72, 120,
		172, 251, 20, 83, 121, 201, 140, 98, 170, 246, 121, 218, 19, 115, 42, 135, 60, 239, 30,
		32, 49, 170, 171, 204, 196, 197, 160, 158, 168, 47, 23, 110, 139, 123, 222, 222, 245,
		98, 125, 208, 70, 39, 110, 186, 146, 254, 66, 185, 118, 3, 78, 32, 47, 179, 197, 93,
		79, 240, 204, 78, 236, 133, 213, 173, 117, 94, 63, 154, 68, 89, 236, 138, 0, 247, 242,
		212, 245, 33, 249, 0, 35, 246, 233, 0, 124, 86, 198, 162, 201, 54, 19, 26, 196, 75,
		254, 71, 70, 238, 51, 2, 23, 185, 152, 139, 134, 65, 107, 129, 114, 244, 47, 251, 240,
		80, 193, 23,
	}
	valid := env.AltBn128PairingCheck(signature)
	env.LogString("AltBn128PairingCheck valid: " + fmt.Sprintf("%t", valid))
	env.ContractValueReturn([]byte("1"))
}

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
			env.LogString("Promise result at index: " + fmt.Sprintf("%x", result))
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

//go:export TestPromiseBatchActionAddKeyWithFullAccess
func TestPromiseBatchActionAddKeyWithFullAccess() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte(accountId)
	promiseIdx := env.PromiseBatchCreate(accountID)
	publicKey, _ := types.PublicKeyFromString("ed25519:ExeqWPvjcUjLX3NfTk3JzisaXLjsCqJNZFCj7ub92RQW")
	env.LogString("publicKey publicKey.Curve: " + fmt.Sprintf("%d", publicKey.Curve) + " " + "publicKey.String(): " + publicKey.ToBase58String())
	nonce := uint64(0)
	env.PromiseBatchActionAddKeyWithFullAccess(promiseIdx, publicKey.Bytes(), nonce)
	env.LogString("Promise batch action add key with full access with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchActionAddKeyWithFunctionCall
func TestPromiseBatchActionAddKeyWithFunctionCall() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte(accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)
	publicKey, _ := types.PublicKeyFromString("ed25519:BeWDy6pKWCiTkHewFcNunbg883abSqVCW42tpUpCrCVU")
	env.LogString("publicKey publicKey.Curve: " + fmt.Sprintf("%d", publicKey.Curve) + " " + "publicKey.String(): " + publicKey.ToBase58String())
	nonce := uint64(0)
	amount := types.Uint128{Hi: 0, Lo: 1000}
	receiverId := []byte("receiver.near")
	functionName := []byte("TestLogStringUtf8")

	env.PromiseBatchActionAddKeyWithFunctionCall(promiseIdx, publicKey.Bytes(), nonce, amount, receiverId, functionName)
	env.LogString("Promise batch action add key with function call with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

//go:export TestPromiseBatchActionDeleteKey
func TestPromiseBatchActionDeleteKey() {
	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	accountID := []byte(accountId)

	promiseIdx := env.PromiseBatchCreate(accountID)
	publicKey, _ := types.PublicKeyFromString("ed25519:BeWDy6pKWCiTkHewFcNunbg883abSqVCW42tpUpCrCVU")
	env.LogString("publicKey publicKey.Curve: " + fmt.Sprintf("%d", publicKey.Curve) + " " + "publicKey.String(): " + publicKey.ToBase58String())

	env.PromiseBatchActionDeleteKey(promiseIdx, publicKey.Bytes())
	env.LogString("Promise batch action delete key with index: " + fmt.Sprintf("%d", promiseIdx))
	env.ContractValueReturn([]byte("1"))
}

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

//go:export TestPromiseYieldCreateYieldResume
func TestPromiseYieldCreateYieldResume() {
	prepaidGas := env.GetPrepaidGas()
	env.LogString(fmt.Sprintf("Prepaid gas: %d", prepaidGas.Inner))

	functionName := []byte("TestContractValueReturn")
	arguments := []byte("{}")
	gasWeight := uint64(0)

	promiseIdx := env.PromiseYieldCreate(functionName, arguments, prepaidGas.Inner/3, gasWeight)
	env.LogString("Promise yield create with index: " + fmt.Sprintf("%d", promiseIdx))
	data, err := env.ReadRegisterSafe(env.DataIdRegister)
	if err != nil {
		env.LogString("Error : " + err.Error())
	}
	env.LogString("hex.EncodeToString(data) : " + hex.EncodeToString(data))
	env.LogString("raw(data) : " + fmt.Sprintf("%x", data))
	result := env.PromiseYieldResume(data, arguments)

	env.LogString("result : " + fmt.Sprintf("%d", result))

	env.ContractValueReturn([]byte("1"))
}

// Promises API Action

// Promise API Results

//go:export TestPromiseResultsCount
func TestPromiseResultsCount() {
	count := env.PromiseResultsCount()
	env.LogString("Promise results count: " + fmt.Sprintf("%d", count))
	env.ContractValueReturn([]byte("1"))
}

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

	env.PromiseReturn(chainedPromiseIdx)

	env.LogString("Promise returned with ID: " + fmt.Sprintf("%d", chainedPromiseIdx))
	env.ContractValueReturn([]byte("1"))
}

// Promise API Results
