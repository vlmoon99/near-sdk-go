//go:build !staticcheck
// +build !staticcheck

package system

// This is the SystemNear class. It contains no internal structures or data.
// All storage and registers are provided by the NEAR Blockchain environment in the WASM runtime.
type SystemNear struct{}

// ReadRegister provides read register functionality.
//
// registerId is the ID of the register from which we want to read the data.
//
// ptr is a pointer to the buffer where this function will write the data from the register.
//
//go:wasmimport env read_register
func (SystemNear) ReadRegister(registerId, ptr uint64)

// RegisterLen provides register length retrieval functionality.
//
// registerId is the ID of the register whose length we want to obtain.
//
//go:wasmimport env register_len
func (SystemNear) RegisterLen(registerId uint64) uint64

// WriteRegister is a function that provides write register functionality.
//
// registerId is the ID of the register where we want to write the data.
//
// dataLen is the length of the data to be written.
//
// dataPtr is a pointer to the buffer containing the data to be written.
//
//go:wasmimport env write_register
func (SystemNear) WriteRegister(registerId, dataLen, dataPtr uint64)

// StorageWrite writes a key-value pair into storage.
// If a key-value pair with the same key already exists, it returns 1; otherwise, it returns 0.
// Storage functions are typically used to upgrade or migrate the contract state.
//
// keyLen is the length of the key.
//
// keyPtr is a pointer to the key.
//
// valueLen is the length of the value.
//
// valuePtr is a pointer to the value.
//
// registerId is the ID of the register where the operation result is stored.
//
//go:wasmimport env storage_write
func (SystemNear) StorageWrite(keyLen, keyPtr, valueLen, valuePtr, registerId uint64) uint64

// StorageRead reads the value stored under the given key.
// Storage functions are typically used to upgrade or migrate the contract state.
//
// keyLen is the length of the key.
//
// keyPtr is a pointer to the key.
//
// registerId is the ID of the register where the retrieved value is stored.
//
//go:wasmimport env storage_read
func (SystemNear) StorageRead(keyLen, keyPtr, registerId uint64) uint64

// StorageRemove removes the value stored under the given key.
// If the key-value pair existed, it returns 1; otherwise, it returns 0.
// Storage functions are typically used to upgrade or migrate the contract state.
//
// keyLen is the length of the key.
//
// keyPtr is a pointer to the key.
//
// registerId is the ID of the register where the operation result is stored.
//
//go:wasmimport env storage_remove
func (SystemNear) StorageRemove(keyLen, keyPtr, registerId uint64) uint64

// StorageHasKey checks if a key-value pair exists in the storage.
// Storage functions are typically used to upgrade or migrate the contract state.
//
// keyLen is the length of the key.
//
// keyPtr is a pointer to the key.
//
//go:wasmimport env storage_has_key
func (SystemNear) StorageHasKey(keyLen, keyPtr uint64) uint64

// CurrentAccountId retrieves the ID of the account that owns the current contract.
//
// registerId is the ID of the register where the current account ID will be written.
//
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env current_account_id
func (SystemNear) CurrentAccountId(registerId uint64)

// SignerAccountId retrieves the ID of the account that signed the original transaction
// or issued the initial cross-contract call.
//
// registerId is the ID of the register where the signer account ID will be written.
//
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env signer_account_id
func (SystemNear) SignerAccountId(registerId uint64)

// SignerAccountPk retrieves the public key of the account that performed the signing.
//
// registerId is the ID of the register where the signer account public key will be written.
//
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env signer_account_pk
func (SystemNear) SignerAccountPk(registerId uint64)

// PredecessorAccountId retrieves the ID of the account that was the previous contract
// in the chain of cross-contract calls.
// If this is the first contract, it is equal to `signer_account_id`.
//
// registerId is the ID of the register where the predecessor account ID will be written.
//
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env predecessor_account_id
func (SystemNear) PredecessorAccountId(registerId uint64)

// Input retrieves the smart contract function input for the contract call, serialized as bytes.
//
// registerId is the ID of the register where the smart contract input will be written.
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env input
func (SystemNear) Input(registerId uint64)

// BlockIndex retrieves the current block index (aka height of the current block)
//
//go:wasmimport env block_index
func (SystemNear) BlockIndex() uint64

// BlockTimestamp retrieves the current block timestamp, i.e, number of non-leap-nanoseconds since January 1, 1970 0:00:00 UTC.
//
//go:wasmimport env block_timestamp
func (SystemNear) BlockTimestamp() uint64

// EpochHeight retrieves the current epoch height
//
//go:wasmimport env epoch_height
func (SystemNear) EpochHeight() uint64

// StorageUsage retrieves the current total storage usage of this smart contract that this account would be paying for.
//
//go:wasmimport env storage_usage
func (SystemNear) StorageUsage() uint64

// AccountBalance retrieves the balance attached to the given account. This includes the attached_deposit that was
// attached to the transaction.
//
// balancePtr is a pointer to the buffer containing the data to be written.
//
//go:wasmimport env account_balance
func (SystemNear) AccountBalance(balancePtr uint64)

// AccountLockedBalance retrieves the balance that was attached to the call that will be immediately deposited before the
// contract execution starts.
//
// balancePtr is a pointer to the buffer containing the data to be written.
//
//go:wasmimport env account_locked_balance
func (SystemNear) AccountLockedBalance(balancePtr uint64)

// AttachedDeposit retrieves the balance locked for potential validator staking.
//
// balancePtr is a pointer to the buffer containing the data to be written.
//
//go:wasmimport env attached_deposit
func (SystemNear) AttachedDeposit(balancePtr uint64)

// PrepaidGas retrieves the amount of gas attached to the call that can be used to pay for the gas fees.
//
//go:wasmimport env prepaid_gas
func (SystemNear) PrepaidGas() uint64

// PrepaidGas retrieves the gas that was already burnt during the contract execution (cannot exceed `prepaid_gas`).
//
//go:wasmimport env used_gas
func (SystemNear) UsedGas() uint64

// RandomSeed returns the random seed from the current block. This 32 byte hash is based on the VRF value from
// the block. This value is not modified in any way each time this function is called within the
// same method/block.
//
//go:wasmimport env random_seed
func (SystemNear) RandomSeed(registerId uint64)

// Sha256 computes the SHA-256 hash of a sequence of bytes.
//
// valueLen is the length of the input data to be hashed.
//
// valuePtr is a pointer to the input data.
//
// registerId is the ID of the register where the hash result will be written.
//
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env sha256
func (SystemNear) Sha256(valueLen, valuePtr, registerId uint64)

// Keccak256 computes the Keccak-256 hash of a sequence of bytes.
//
// valueLen is the length of the input data to be hashed.
//
// valuePtr is a pointer to the input data.
//
// registerId is the ID of the register where the hash result will be written.
//
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env keccak256
func (SystemNear) Keccak256(valueLen, valuePtr, registerId uint64)

// Keccak512 computes the Keccak-512 hash of a sequence of bytes.
//
// valueLen is the length of the input data to be hashed.
//
// valuePtr is a pointer to the input data.
//
// registerId is the ID of the register where the hash result will be written.
//
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env keccak512
func (SystemNear) Keccak512(valueLen, valuePtr, registerId uint64)

// Ripemd160 computes the RIPEMD-160 hash of a sequence of bytes.
//
// This returns a 20-byte hash.
//
// valueLen is the length of the input data to be hashed.
//
// valuePtr is a pointer to the input data.
//
// registerId is the ID of the register where the hash result will be written.
//
// For the standard implementation, it is set to `const AtomicOpRegister uint64 = math.MaxUint64 - 2`.
//
//go:wasmimport env ripemd160
func (SystemNear) Ripemd160(valueLen, valuePtr, registerId uint64)

// Ecrecover recovers an ECDSA signer address from a 32-byte message hash and a corresponding signature
// along with the `v` recovery byte.
// malleabilityFlag indicates whether the function should check for signature malleability,
// which is generally ideal for transactions.
//
// Returns 64 bytes representing the public key if the recovery was successful.
//
// hashLen is the length of the hash.
//
// hashPtr is a pointer to the hash data.
//
// sigLen is the length of the signature.
//
// sigPtr is a pointer to the signature data.
//
// v is the recovery byte.
//
// registerId is the ID of the register where the recovered public key will be written.
//
//go:wasmimport env ecrecover
func (SystemNear) Ecrecover(hashLen, hashPtr, sigLen, sigPtr, v, malleabilityFlag, registerId uint64) uint64

// Ed25519Verify verifies the signature of a message using the provided ED25519 public key.
//
// sigLen is the length of the signature.
//
// sigPtr is a pointer to the signature data.
//
// msgLen is the length of the message.
//
// msgPtr is a pointer to the message data.
//
// pubKeyLen is the length of the public key.
//
// pubKeyPtr is a pointer to the public key data.
//
// Returns 1 if the signature is valid, 0 otherwise.
//
//go:wasmimport env ed25519_verify
func (SystemNear) Ed25519Verify(sigLen, sigPtr, msgLen, msgPtr, pubKeyLen, pubKeyPtr uint64) uint64

// AltBn128G1Multiexp computes the multi-exponentiation operation on the `alt_bn128` curve.
// `alt_bn128` is a specific curve from the Barreto-Naehrig (BN) family, particularly
// well-suited for ZK proofs.
// See also: [EIP-196](https://eips.ethereum.org/EIPS/eip-196).
//
// valueLen is the length of the input data.
//
// valuePtr is a pointer to the input data.
//
// registerId is the ID of the register where the result will be written.
//
//go:wasmimport env alt_bn128_g1_multiexp
func (SystemNear) AltBn128G1Multiexp(valueLen, valuePtr, registerId uint64)

// AltBn128G1Sum computes the sum of multiple G1 points on the `alt_bn128` curve.
// `alt_bn128` is a specific curve from the Barreto-Naehrig (BN) family, particularly
// well-suited for ZK proofs.
// See also: [EIP-196](https://eips.ethereum.org/EIPS/eip-196).
//
// valueLen is the length of the input data.
//
// valuePtr is a pointer to the input data.
//
// registerId is the ID of the register where the result will be written.
//
//go:wasmimport env alt_bn128_g1_sum
func (SystemNear) AltBn128G1SumSystem(valueLen, valuePtr, registerId uint64)

// AltBn128PairingCheck performs a pairing check on the `alt_bn128` curve to validate
// cryptographic proofs.
// `alt_bn128` is a specific curve from the Barreto-Naehrig (BN) family, particularly
// well-suited for ZK proofs.
// See also: [EIP-197](https://eips.ethereum.org/EIPS/eip-197).
//
// valueLen is the length of the input data.
//
// valuePtr is a pointer to the input data.
//
// Returns 1 if the pairing check is valid, 0 otherwise.
//
//go:wasmimport env alt_bn128_pairing_check
func (SystemNear) AltBn128PairingCheckSystem(valueLen, valuePtr uint64) uint64

// ValidatorStake returns the current stake of a given account.
// If the account is not a validator, it returns 0.
//
// accountIdLen is the length of the account ID.
//
// accountIdPtr is a pointer to the account ID string.
//
// stakePtr is a pointer to the register where the stake amount will be written.
//
//go:wasmimport env validator_stake
func (SystemNear) ValidatorStake(accountIdLen, accountIdPtr, stakePtr uint64)

// ValidatorTotalStake returns the total stake of all validators in the current epoch.
//
// stakePtr is a pointer to the register where the total stake amount will be written.
//
//go:wasmimport env validator_total_stake
func (SystemNear) ValidatorTotalStake(stakePtr uint64)

// ValueReturn sets the blob of data as the return value of the contract.
//
// valueLen is the length of the value.
//
// valuePtr is a pointer to the value.
//
//go:wasmimport env value_return
func (SystemNear) ValueReturn(valueLen, valuePtr uint64)

// PanicUtf8 terminates the execution of the program with a UTF-8 encoded message.
//
// len is the length of the message.
//
// ptr is a pointer to the UTF-8 encoded message.
//
//go:wasmimport env panic_utf8
func (SystemNear) PanicUtf8(len, ptr uint64)

// LogUtf8 logs a message encoded in UTF-8.
//
// len is the length of the message.
//
// ptr is a pointer to the UTF-8 encoded message.
//
//go:wasmimport env log_utf8
func (SystemNear) LogUtf8(len, ptr uint64)

// LogUtf16 logs a message encoded in UTF-16.
//
// len is the length of the message.
//
// ptr is a pointer to the UTF-16 encoded message.
//
//go:wasmimport env log_utf16
func (SystemNear) LogUtf16(len, ptr uint64)

// PromiseCreate creates a promise to execute a method on a specified account with the given arguments,
// attaching the specified amount and gas.
//
// accountIdLen is the length of the target account ID.
//
// accountIdPtr is a pointer to the target account ID.
//
// functionNameLen is the length of the function name to be called.
//
// functionNamePtr is a pointer to the function name.
//
// argumentsLen is the length of the arguments for the function call.
//
// argumentsPtr is a pointer to the arguments.
//
// amountPtr is a pointer to the attached amount (must be a Unit128 value).
//
// gas is the amount of gas allocated for the execution (must be a u64 value).
//
//go:wasmimport env promise_create
func (SystemNear) PromiseCreate(accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64

// PromiseThen attaches a callback that executes after the promise specified by `promiseIndex` completes.
//
// promiseIndex is the index of the promise to attach the callback to.
//
// accountIdLen is the length of the target account ID.
//
// accountIdPtr is a pointer to the target account ID.
//
// functionNameLen is the length of the function name to be called.
//
// functionNamePtr is a pointer to the function name.
//
// argumentsLen is the length of the arguments for the function call.
//
// argumentsPtr is a pointer to the arguments.
//
// amountPtr is a pointer to the attached amount (must be a Unit128 value).
//
// gas is the amount of gas allocated for the execution (must be a u64 value).
//
//go:wasmimport env promise_then
func (SystemNear) PromiseThen(promiseIndex, accountIdLen, accountIdPtr, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64) uint64

// PromiseAnd creates a new promise that completes when all promises passed as arguments complete.
//
// promiseIdxPtr is a pointer to an array of promise indexes.
//
// promiseIdxCount is the number of promises in the array.
//
//go:wasmimport env promise_and
func (SystemNear) PromiseAnd(promiseIdxPtr, promiseIdxCount uint64) uint64

// PromiseBatchCreate creates a new batch promise for a specified account.
//
// accountIdLen is the length of the target account ID.
//
// accountIdPtr is a pointer to the target account ID.
//
// Returns the index of the created promise batch.
//
//go:wasmimport env promise_batch_create
func (SystemNear) PromiseBatchCreate(accountIdLen, accountIdPtr uint64) uint64

// PromiseBatchThen creates a "then" callback function that executes after the specified promise completes.
//
// promiseIndex is the index of the existing promise to attach the callback to.
//
// accountIdLen is the length of the target account ID for the callback execution.
//
// accountIdPtr is a pointer to the target account ID for the callback execution.
//
// Returns the index of the created callback promise.
//
//go:wasmimport env promise_batch_then
func (SystemNear) PromiseBatchThen(promiseIndex, accountIdLen, accountIdPtr uint64) uint64

// PromiseBatchActionCreateAccount represents the action of creating a new account as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the account creation action will be attached.
//
//go:wasmimport env promise_batch_action_create_account
func (SystemNear) PromiseBatchActionCreateAccount(promiseIndex uint64)

// PromiseBatchActionDeployContract represents the action of deploying a contract as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the contract deployment action will be attached.
//
// codeLen is the length of the contract code in bytes.
//
// codePtr is a pointer to the contract code to be deployed.
//
// Returns nothing, but performs the contract deployment action as part of the promise batch.
//
//go:wasmimport env promise_batch_action_deploy_contract
func (SystemNear) PromiseBatchActionDeployContract(promiseIndex, codeLen, codePtr uint64)

// PromiseBatchActionFunctionCall represents the action of invoking a function on a contract as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the function call action will be attached.
//
// functionNameLen is the length of the function name.
//
// functionNamePtr is a pointer to the function name.
//
// argumentsLen is the length of the arguments for the function call.
//
// argumentsPtr is a pointer to the arguments to be passed.
//
// amountPtr represents the amount to attach to the call, which should be a Unit128 value.
//
// gas is the amount of gas to attach to the function call, represented as u64.
//
// Returns nothing, but performs the function call action as part of the promise batch.
//
//go:wasmimport env promise_batch_action_function_call
func (SystemNear) PromiseBatchActionFunctionCall(promiseIndex, functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, amountPtr, gas uint64)

// PromiseBatchActionFunctionCallWeight represents the action of invoking a function on a contract with a specified weight as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the function call action will be attached.
//
// functionNameLen is the length of the function name.
//
// functionNamePtr is a pointer to the function name.
//
// argumentsLen is the length of the arguments for the function call.
//
// argumentsPtr is a pointer to the arguments to be passed.
//
// amountPtr represents the amount to attach to the call, which should be a Unit128 value.
//
// gas is the amount of gas to attach to the function call, represented as u64.
//
// weight represents the weight for the function call action.
//
// Returns nothing, but performs the function call action with weight as part of the promise batch.
//
//go:wasmimport env promise_batch_action_function_call_weight
func (SystemNear) PromiseBatchActionFunctionCallWeight(promise_index, function_name_len, function_name_ptr, arguments_len, arguments_ptr, amount_ptr, gas, weight uint64)

// PromiseBatchActionTransfer represents the action of transferring tokens as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the transfer action will be attached.
//
// amountPtr represents the amount to transfer, which should be a Unit128 value.
//
// Returns nothing, but performs the transfer action as part of the promise batch.
//
//go:wasmimport env promise_batch_action_transfer
func (SystemNear) PromiseBatchActionTransfer(promiseIndex, amountPtr uint64)

// PromiseBatchActionStake represents the action of staking tokens as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the staking action will be attached.
//
// amountPtr represents the amount to stake, which should be a Unit128 value.
//
// publicKeyLen is the length of the public key to associate with the staking action.
//
// publicKeyPtr is a pointer to the public key to associate with the staking action.
//
// Returns nothing, but performs the staking action as part of the promise batch.
//
//go:wasmimport env promise_batch_action_stake
func (SystemNear) PromiseBatchActionStake(promiseIndex, amountPtr, publicKeyLen, publicKeyPtr uint64)

// PromiseBatchActionAddKeyWithFullAccess represents the action of adding a key with full access as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the key addition action will be attached.
//
// publicKeyLen is the length of the public key to associate with the full access key.
//
// publicKeyPtr is a pointer to the public key to associate with the full access key.
//
// nonce is the nonce to associate with the key addition action.
//
// Returns nothing, but performs the key addition with full access as part of the promise batch.
//
//go:wasmimport env promise_batch_action_add_key_with_full_access
func (SystemNear) PromiseBatchActionAddKeyWithFullAccess(promiseIndex, publicKeyLen, publicKeyPtr, nonce uint64)

// PromiseBatchActionAddKeyWithFunctionCall represents the action of adding a key with function call access as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the key addition action will be attached.
//
// publicKeyLen is the length of the public key to associate with the function call access key.
//
// publicKeyPtr is a pointer to the public key to associate with the function call access key.
//
// nonce is the nonce to associate with the key addition action.
//
// allowancePtr represents the allowance for the function call access, which should be a Unit128 value.
//
// receiverIdLen is the length of the receiver account ID for the function call access.
//
// receiverIdPtr is a pointer to the receiver account ID for the function call access.
//
// functionNamesLen is the length of the function names list to associate with the key.
//
// functionNamesPtr is a pointer to the list of function names to associate with the key.
//
// Returns nothing, but performs the key addition with function call access as part of the promise batch.
//
//go:wasmimport env promise_batch_action_add_key_with_function_call
func (SystemNear) PromiseBatchActionAddKeyWithFunctionCall(promiseIndex, publicKeyLen, publicKeyPtr, nonce, allowancePtr, receiverIdLen, receiverIdPtr, functionNamesLen, functionNamesPtr uint64)

// PromiseBatchActionDeleteKey represents the action of deleting a key as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the key deletion action will be attached.
//
// publicKeyLen is the length of the public key to delete.
//
// publicKeyPtr is a pointer to the public key to delete.
//
// Returns nothing, but performs the key deletion as part of the promise batch.
//
//go:wasmimport env promise_batch_action_delete_key
func (SystemNear) PromiseBatchActionDeleteKey(promiseIndex, publicKeyLen, publicKeyPtr uint64)

// PromiseBatchActionDeleteAccount represents the action of deleting an account as part of a promise batch.
// This action can be called after a promise batch has already been created with the represented accountId.
//
// promiseIndex is the index of the promise batch to which the account deletion action will be attached.
//
// beneficiaryIdLen is the length of the beneficiary account ID to receive any remaining balance after deletion.
//
// beneficiaryIdPtr is a pointer to the beneficiary account ID to receive any remaining balance after deletion.
//
// Returns nothing, but performs the account deletion as part of the promise batch.
//
//go:wasmimport env promise_batch_action_delete_account
func (SystemNear) PromiseBatchActionDeleteAccount(promiseIndex, beneficiaryIdLen, beneficiaryIdPtr uint64)

// PromiseYieldCreate creates a promise that will execute a method on the current account with the given arguments.
// It writes a resumption token (data id) to the specified register. The callback method will execute
// after `promise_yield_resume` is called with the data id, or after enough blocks have passed. The timeout
// length is specified by the protocol-level parameter `yield_timeout_length_in_blocks = 200`.
//
// The callback method will execute with a single promise input. The input will either be a payload
// provided by the user when calling `promise_yield_resume`, or a `PromiseError` in case of timeout.
// Resumption tokens are specific to the local account, and `promise_yield_resume` must be called from
// a method of the same contract.
//
// functionNameLen is the length of the function name to execute on the current account.
//
// functionNamePtr is a pointer to the function name to execute on the current account.
//
// argumentsLen is the length of the arguments to pass to the function.
//
// argumentsPtr is a pointer to the arguments to pass to the function.
//
// gas is the amount of gas to attach to the promise creation, specified as `u64`.
//
// gasWeight is the weight of the gas to apply to the promise, specified as `u64`.
//
// registerId is the ID of the register where the resumption token will be written.
//
// Returns the index of the created promise.
//
//go:wasmimport env promise_yield_create
func (SystemNear) PromiseYieldCreate(functionNameLen, functionNamePtr, argumentsLen, argumentsPtr, gas, gasWeight, registerId uint64) uint64

// PromiseYieldResume accepts a resumption token `data_id` created by `promise_yield_create` on the local account.
// The `data` is a payload to be passed to the callback method as a promise input. Returns false if
// no promise yield with the specified `data_id` is found. Returns true otherwise, guaranteeing
// that the callback method will be executed with a user-provided payload.
//
// If `promise_yield_resume` is called multiple times with the same `data_id`, it is possible to get
// multiple 'true' results. The payload from the first successful call is passed to the callback.
//
// dataIdLen is the length of the `data_id` (resumption token) to resume.
//
// dataIdPtr is a pointer to the `data_id` (resumption token) to resume.
//
// payloadLen is the length of the payload to pass to the callback method.
//
// payloadPtr is a pointer to the payload to pass to the callback method.
//
// Returns a `uint32`, where '1' indicates success, and '0' indicates failure (if no promise yield is found).
//
//go:wasmimport env promise_yield_resume
func (SystemNear) PromiseYieldResume(dataIdLen, dataIdPtr, payloadLen, payloadPtr uint64) uint32

// PromiseResultsCount returns the number of complete and incomplete callback results from the promises
// that triggered the current callback execution. This function can only be called if the current function
// was invoked by a callback. It helps in checking how many promise results are available for processing.
//
// Returns the count of complete and incomplete callback results.
//
//go:wasmimport env promise_results_count
func (SystemNear) PromiseResultsCount() uint64

// PromiseResult retrieves the execution result of a promise identified by `resultIdx` that caused the
// current callback. It allows access to the outcome of the promise. This function can only be called if the
// current function was invoked by a callback.
//
// resultIdx is the index of the result to retrieve.
//
// registerId is the register where the result will be stored or accessed.
//
// Returns the execution result of the specified promise.
//
//go:wasmimport env promise_result
func (SystemNear) PromiseResult(resultIdx uint64, registerId uint64) uint64

// PromiseReturn considers the execution result of the promise specified by `promiseId` as the execution
// result of the current function. This allows the callback to finalize or return the result of a promise as
// the outcome of the function.
//
// promiseId is the index of the promise whose result will be treated as the execution result of this function.
//
//go:wasmimport env promise_return
func (SystemNear) PromiseReturn(promiseId uint64)
