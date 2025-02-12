// Package "system" provides system bindings to interact with the NEAR Blockchain environment.
// These bindings facilitate the creation of smart contracts and provide all necessary functions for working with them.
// The implementation of these system methods can be found in the `sdk` package, specifically in the `env.go` file.
// It is not recommended to use this package directly unless you fully understand what you are doing.
//
// system_mock.go - represents mocks for unit tests
//
// system_near.go - represents NEAR System methods which the blockchain will provide inside the WASM environment.
package system

type System interface {
	//Registers API
	ReadRegister(registerId, ptr uint64)
	RegisterLen(registerId uint64) uint64
	WriteRegister(registerId, dataLen, dataPtr uint64)
	//Registers API

	// Storage API
	StorageWrite(keyLen, keyPtr, valueLen, valuePtr, registerId uint64) uint64
	StorageRead(keyLen uint64, keyPtr uint64, registerId uint64) uint64
	StorageRemove(keyLen, keyPtr, registerId uint64) uint64
	StorageHasKey(keyLen, keyPtr uint64) uint64
	// Storage API

	//Context API
	CurrentAccountId(registerId uint64)
	SignerAccountId(registerId uint64)
	SignerAccountPk(registerId uint64)
	PredecessorAccountId(registerId uint64)
	Input(registerId uint64)
	BlockIndex() uint64
	BlockTimestamp() uint64
	EpochHeight() uint64
	StorageUsage() uint64
	//Context API

	// Economics API
	AccountBalance(balancePtr uint64)
	AccountLockedBalance(balancePtr uint64)
	AttachedDeposit(balancePtr uint64)
	PrepaidGas() uint64
	UsedGas() uint64
	// Economics API

	// Math API
	RandomSeed(registerId uint64)
	Sha256(valueLen, valuePtr, registerId uint64)
	Keccak256(valueLen, valuePtr, registerId uint64)
	Keccak512(valueLen, valuePtr, registerId uint64)
	Ripemd160(valueLen, valuePtr, registerId uint64)
	Ecrecover(hashLen, hashPtr, sigLen, sigPtr, v, malleabilityFlag, registerId uint64) uint64
	Ed25519Verify(sigLen, sigPtr, msgLen, msgPtr, pubKeyLen, pubKeyPtr uint64) uint64
	AltBn128G1Multiexp(valueLen, valuePtr, registerId uint64)
	AltBn128G1SumSystem(valueLen, valuePtr, registerId uint64)
	AltBn128PairingCheckSystem(valueLen, valuePtr uint64) uint64
	// Math API

	// Validator API
	ValidatorStake(accountIdLen, accountIdPtr, stakePtr uint64)
	ValidatorTotalStake(stakePtr uint64)
	// Validator API

	// Miscellaneous API
	ValueReturn(valueLen, valuePtr uint64)
	PanicUtf8(len, ptr uint64)
	LogUtf8(len, ptr uint64)
	LogUtf16(len, ptr uint64)
	// Abort(msgPtr, filenamePtr, line, col uint32)
	// Panic()
	// Miscellaneous API

	// Promises API
	PromiseCreate(accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64
	PromiseThen(promiseIndex, accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64
	PromiseAnd(promiseIdxPtr, promiseIdxCount uint64) uint64
	PromiseBatchCreate(accountIdLen, accountIdPtr uint64) uint64
	PromiseBatchThen(promiseIndex, accountIdLen, accountIdPtr uint64) uint64
	// Promises API

	// Promise API Actions
	PromiseBatchActionCreateAccount(promiseIndex uint64)
	PromiseBatchActionDeployContract(promiseIndex, codeLen, codePtr uint64)
	PromiseBatchActionFunctionCall(promiseIndex, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64)
	PromiseBatchActionFunctionCallWeight(promise_index, function_name_len, function_name_ptr, arguments_len, arguments_ptr, amount_ptr, gas, weight uint64)
	PromiseBatchActionTransfer(promiseIndex, amountPtr uint64)
	PromiseBatchActionStake(promiseIndex, amountPtr, publicKeyLen, publicKeyPtr uint64)
	PromiseBatchActionAddKeyWithFullAccess(promiseIndex, publicKeyLen, publicKeyPtr, nonce uint64)
	PromiseBatchActionAddKeyWithFunctionCall(promiseIndex, publicKeyLen, publicKeyPtr, nonce, allowancePtr, receiverIdLen, receiverIdPtr, functionNamesLen, functionNamesPtr uint64)
	PromiseBatchActionDeleteKey(promiseIndex, publicKeyLen, publicKeyPtr uint64)
	PromiseBatchActionDeleteAccount(promiseIndex, beneficiaryIdLen, beneficiaryIdPtr uint64)
	PromiseYieldCreate(functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, gas, gasWeight, registerId uint64) uint64
	PromiseYieldResume(dataIdLen, dataIdPtr, payloadLen, payloadPtr uint64) uint32
	// Promise API Actions

	// Promise API Results
	PromiseResultsCount() uint64
	PromiseResult(resultIdx uint64, registerId uint64) uint64
	PromiseReturn(promiseId uint64)
	// Promise API Results

}
