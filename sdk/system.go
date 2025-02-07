package sdk

// Registers
//
//go:wasmimport env read_register
func readRegister(registerId, ptr uint64)

//go:wasmimport env register_len
func registerLen(registerId uint64) uint64

//go:wasmimport env write_register
func writeRegister(registerId, dataLen, dataPtr uint64)

// Context API
//
//go:wasmimport env current_account_id
func currentAccountId(registerId uint64)

//go:wasmimport env signer_account_id
func signerAccountId(registerId uint64)

//go:wasmimport env signer_account_pk
func signerAccountPk(registerId uint64)

//go:wasmimport env predecessor_account_id
func predecessorAccountId(registerId uint64)

//go:wasmimport env input
func input(registerId uint64)

//go:wasmimport env block_index
func blockIndex() uint64

//go:wasmimport env block_timestamp
func blockTimestamp() uint64

//go:wasmimport env epoch_height
func epochHeight() uint64

//go:wasmimport env storage_usage
func storageUsage() uint64

// Economics API
//
//go:wasmimport env account_balance
func accountBalance(balancePtr uint64)

//go:wasmimport env account_locked_balance
func accountLockedBalance(balancePtr uint64)

//go:wasmimport env attached_deposit
func attachedDeposit(balancePtr uint64)

//go:wasmimport env prepaid_gas
func prepaidGas() uint64

//go:wasmimport env used_gas
func usedGas() uint64

// Math API
//
//go:wasmimport env random_seed
func randomSeed(registerId uint64)

//go:wasmimport env sha256
func sha256(valueLen, valuePtr, registerId uint64)

//go:wasmimport env keccak256
func keccak256(valueLen, valuePtr, registerId uint64)

//go:wasmimport env keccak512
func keccak512(valueLen, valuePtr, registerId uint64)

//go:wasmimport env ripemd160
func ripemd160(valueLen, valuePtr, registerId uint64)

//go:wasmimport env ecrecover
func ecrecover(hashLen, hashPtr, sigLen, sigPtr, v, malleabilityFlag, registerId uint64) uint64

//go:wasmimport env ed25519_verify
func ed25519Verify(sigLen, sigPtr, msgLen, msgPtr, pubKeyLen, pubKeyPtr uint64) uint64

//go:wasmimport env alt_bn128_g1_multiexp
func altBn128G1Multiexp(valueLen, valuePtr, registerId uint64)

//go:wasmimport env alt_bn128_g1_sum
func altBn128G1SumSystem(valueLen, valuePtr, registerId uint64)

//go:wasmimport env alt_bn128_pairing_check
func altBn128PairingCheckSystem(valueLen, valuePtr uint64) uint64

// Validator API
//
//go:wasmimport env validator_stake
func validatorStake(accountIdLen, accountIdPtr, stakePtr uint64)

//go:wasmimport env validator_total_stake
func validatorTotalStake(stakePtr uint64)

// Miscellaneous API
//
//go:wasmimport env value_return
func valueReturn(valueLen, valuePtr uint64)

//go:wasmimport env panic
func panic()

//go:wasmimport env panic_utf8
func panicUtf8(len, ptr uint64)

//go:wasmimport env log_utf8
func logUtf8(len, ptr uint64)

//go:wasmimport env log_utf16
func logUtf16(len, ptr uint64)

//go:wasmimport env abort
func abort(msgPtr, filenamePtr, line, col uint32)

// Storage API
//
//go:wasmimport env storage_write
func storageWriteSys(keyLen, keyPtr, valueLen, valuePtr, registerId uint64) uint64

//go:wasmimport env storage_read
func storageReadSys(keyLen uint64, keyPtr uint64, registerId uint64) uint64

//go:wasmimport env storage_remove
func storageRemoveSys(keyLen, keyPtr, registerId uint64) uint64

//go:wasmimport env storage_has_key
func storageHasKeySys(keyLen, keyPtr uint64) uint64

//go:wasmimport env storage_iter_prefix
func storageIterPrefix(prefixLen, prefixPtr uint64) uint64

//go:wasmimport env storage_iter_range
func storageIterRange(startLen, startPtr, endLen, endPtr uint64) uint64

//go:wasmimport env storage_iter_next
func storageIterNext(iteratorId, keyRegisterId, valueRegisterId uint64) uint64

// Promises API
//
//go:wasmimport env promise_create
func promiseCreateSys(accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64

//go:wasmimport env promise_then
func promiseThenSys(promiseIndex, accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64

//go:wasmimport env promise_and
func promiseAndSys(promiseIdxPtr, promiseIdxCount uint64) uint64

//go:wasmimport env promise_batch_create
func promiseBatchCreateSys(accountIdLen, accountIdPtr uint64) uint64

//go:wasmimport env promise_batch_then
func promiseBatchThenSys(promiseIndex, accountIdLen, accountIdPtr uint64) uint64

// Promise API Actions
//
//go:wasmimport env promise_batch_action_create_account
func promiseBatchActionCreateAccountSys(promiseIndex uint64)

//go:wasmimport env promise_batch_action_deploy_contract
func promiseBatchActionDeployContractSys(promiseIndex, codeLen, codePtr uint64)

//go:wasmimport env promise_batch_action_function_call
func promiseBatchActionFunctionCallSys(promiseIndex, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64)

//go:wasmimport env promise_batch_action_function_call_weight
func promiseBatchActionFunctionCallWeightSys(promise_index, function_name_len, function_name_ptr, arguments_len, arguments_ptr, amount_ptr, gas, weight uint64)

//go:wasmimport env promise_batch_action_transfer
func promiseBatchActionTransferSys(promiseIndex, amountPtr uint64)

//go:wasmimport env promise_batch_action_stake
func promiseBatchActionStakeSys(promiseIndex, amountPtr, publicKeyLen, publicKeyPtr uint64)

//go:wasmimport env promise_batch_action_add_key_with_full_access
func promiseBatchActionAddKeyWithFullAccessSys(promiseIndex, publicKeyLen, publicKeyPtr, nonce uint64)

//go:wasmimport env promise_batch_action_add_key_with_function_call
func promiseBatchActionAddKeyWithFunctionCallSys(promiseIndex, publicKeyLen, publicKeyPtr, nonce, allowancePtr, receiverIdLen, receiverIdPtr, functionNamesLen, functionNamesPtr uint64)

//go:wasmimport env promise_batch_action_delete_key
func promiseBatchActionDeleteKeySys(promiseIndex, publicKeyLen, publicKeyPtr uint64)

//go:wasmimport env promise_batch_action_delete_account
func promiseBatchActionDeleteAccountSys(promiseIndex, beneficiaryIdLen, beneficiaryIdPtr uint64)

//go:wasmimport env promise_yield_create
func promiseYieldCreateSys(functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, gas, gasWeight, registerId uint64) uint64

//go:wasmimport env promise_yield_resume
func promiseYieldResumeSys(dataIdLen, dataIdPtr, payloadLen, payloadPtr uint64) uint32

// Promise API Results
//
//go:wasmimport env promise_results_count
func promiseResultsCountSys() uint64

//go:wasmimport env promise_result
func promiseResultSys(resultIdx uint64, registerId uint64) uint64

//go:wasmimport env promise_return
func promiseReturnSys(promiseId uint64)
