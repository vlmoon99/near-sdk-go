package sdk

// Registers
//
//go:wasmimport env read_register
func ReadRegister(registerId, ptr uint64)

//go:wasmimport env register_len
func RegisterLen(registerId uint64) uint64

//go:wasmimport env write_register
func WriteRegister(registerId, dataLen, dataPtr uint64)

// Context API
//
//go:wasmimport env current_account_id
func CurrentAccountId(registerId uint64)

//go:wasmimport env signer_account_id
func SignerAccountId(registerId uint64)

//go:wasmimport env signer_account_pk
func SignerAccountPk(registerId uint64)

//go:wasmimport env predecessor_account_id
func PredecessorAccountId(registerId uint64)

//go:wasmimport env input
func Input(registerId uint64)

//go:wasmimport env block_index
func BlockIndex() uint64

//go:wasmimport env block_timestamp
func BlockTimestamp() uint64

//go:wasmimport env epoch_height
func EpochHeight() uint64

//go:wasmimport env storage_usage
func StorageUsage() uint64

// Economics API
//
//go:wasmimport env account_balance
func AccountBalance(balancePtr uint64)

//go:wasmimport env account_locked_balance
func AccountLockedBalance(balancePtr uint64)

//go:wasmimport env attached_deposit
func AttachedDeposit(balancePtr uint64)

//go:wasmimport env prepaid_gas
func PrepaidGas() uint64

//go:wasmimport env used_gas
func UsedGas() uint64

// Math API
//
//go:wasmimport env random_seed
func RandomSeed(registerId uint64)

//go:wasmimport env sha256
func Sha256(valueLen, valuePtr, registerId uint64)

//go:wasmimport env keccak256
func Keccak256(valueLen, valuePtr, registerId uint64)

//go:wasmimport env keccak512
func Keccak512(valueLen, valuePtr, registerId uint64)

//go:wasmimport env ripemd160
func Ripemd160(valueLen, valuePtr, registerId uint64)

//go:wasmimport env ecrecover
func Ecrecover(hashLen, hashPtr, sigLen, sigPtr, v, malleabilityFlag, registerId uint64) uint64

//go:wasmimport env ed25519_verify
func Ed25519Verify(sigLen, sigPtr, msgLen, msgPtr, pubKeyLen, pubKeyPtr uint64) uint64

//go:wasmimport env alt_bn128_g1_multiexp
func AltBn128G1Multiexp(valueLen, valuePtr, registerId uint64)

//go:wasmimport env alt_bn128_g1_sum
func AltBn128G1SumSystem(valueLen, valuePtr, registerId uint64)

//go:wasmimport env alt_bn128_pairing_check
func AltBn128PairingCheckSystem(valueLen, valuePtr uint64) uint64

// Validator API
//
//go:wasmimport env validator_stake
func ValidatorStake(accountIdLen, accountIdPtr, stakePtr uint64)

//go:wasmimport env validator_total_stake
func ValidatorTotalStake(stakePtr uint64)

// Miscellaneous API
//
//go:wasmimport env value_return
func ValueReturn(valueLen, valuePtr uint64)

//go:wasmimport env panic
func Panic()

//go:wasmimport env panic_utf8
func PanicUtf8(len, ptr uint64)

//go:wasmimport env log_utf8
func LogUtf8(len, ptr uint64)

//go:wasmimport env log_utf16
func LogUtf16(len, ptr uint64)

//go:wasmimport env abort
func Abort(msgPtr, filenamePtr, line, col uint32)

// Storage API
//
//go:wasmimport env storage_write
func StorageWriteSys(keyLen, keyPtr, valueLen, valuePtr, registerId uint64) uint64

//go:wasmimport env storage_read
func StorageReadSys(keyLen uint64, keyPtr uint64, registerId uint64) uint64

//go:wasmimport env storage_remove
func StorageRemoveSys(keyLen, keyPtr, registerId uint64) uint64

//go:wasmimport env storage_has_key
func StorageHasKeySys(keyLen, keyPtr uint64) uint64

//go:wasmimport env storage_iter_prefix
func StorageIterPrefix(prefixLen, prefixPtr uint64) uint64

//go:wasmimport env storage_iter_range
func StorageIterRange(startLen, startPtr, endLen, endPtr uint64) uint64

//go:wasmimport env storage_iter_next
func StorageIterNext(iteratorId, keyRegisterId, valueRegisterId uint64) uint64

// Promises API
//
//go:wasmimport env promise_create
func PromiseCreateSys(accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64

//go:wasmimport env promise_then
func PromiseThenSys(promiseIndex, accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64

//go:wasmimport env promise_and
func PromiseAndSys(promiseIdxPtr, promiseIdxCount uint64) uint64

//go:wasmimport env promise_batch_create
func PromiseBatchCreateSys(accountIdLen, accountIdPtr uint64) uint64

//go:wasmimport env promise_batch_then
func PromiseBatchThenSys(promiseIndex, accountIdLen, accountIdPtr uint64) uint64

// Promise API Actions
//
//go:wasmimport env promise_batch_action_create_account
func PromiseBatchActionCreateAccountSys(promiseIndex uint64)

//go:wasmimport env promise_batch_action_deploy_contract
func PromiseBatchActionDeployContractSys(promiseIndex, codeLen, codePtr uint64)

//go:wasmimport env promise_batch_action_function_call
func PromiseBatchActionFunctionCallSys(promiseIndex, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64)

//go:wasmimport env promise_batch_action_function_call_weight
func PromiseBatchActionFunctionCallWeightSys(promise_index, function_name_len, function_name_ptr, arguments_len, arguments_ptr, amount_ptr, gas, weight uint64)

//go:wasmimport env promise_batch_action_transfer
func PromiseBatchActionTransferSys(promiseIndex, amountPtr uint64)

//go:wasmimport env promise_batch_action_stake
func PromiseBatchActionStakeSys(promiseIndex, amountPtr, publicKeyLen, publicKeyPtr uint64)

//go:wasmimport env promise_batch_action_add_key_with_full_access
func PromiseBatchActionAddKeyWithFullAccessSys(promiseIndex, publicKeyLen, publicKeyPtr, nonce uint64)

//go:wasmimport env promise_batch_action_add_key_with_function_call
func PromiseBatchActionAddKeyWithFunctionCallSys(promiseIndex, publicKeyLen, publicKeyPtr, nonce, allowancePtr, receiverIdLen, receiverIdPtr, functionNamesLen, functionNamesPtr uint64)

//go:wasmimport env promise_batch_action_delete_key
func PromiseBatchActionDeleteKeySys(promiseIndex, publicKeyLen, publicKeyPtr uint64)

//go:wasmimport env promise_batch_action_delete_account
func PromiseBatchActionDeleteAccountSys(promiseIndex, beneficiaryIdLen, beneficiaryIdPtr uint64)

//go:wasmimport env promise_yield_create
func PromiseYieldCreateSys(functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, gas, gasWeight, registerId uint64) uint64

//go:wasmimport env promise_yield_resume
func PromiseYieldResumeSys(dataIdLen, dataIdPtr, payloadLen, payloadPtr uint64) uint32

// Promise API Results
//
//go:wasmimport env promise_results_count
func PromiseResultsCountSys() uint64

//go:wasmimport env promise_result
func PromiseResultSys(resultIdx uint64, registerId uint64) uint64

//go:wasmimport env promise_return
func PromiseReturnSys(promiseId uint64)
