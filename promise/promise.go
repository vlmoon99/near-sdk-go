// Package promise provides functions for creating and managing promises.
package promise

import (
	"encoding/json"
	"errors"
	"unsafe"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/types"
)

const (
	DefaultGas = 5 * types.ONE_TERA_GAS
	MaxGas     = 300 * types.ONE_TERA_GAS
)

const (
	ErrPromiseFailedPrefix         = "promise failed with status: "
	ErrMarshalingArgsInThen        = "Error marshaling args in Then: "
	ErrGettingCurrentAccountInThen = "Error getting current account in Then: "
	ErrMarshalingArgsInThenCall    = "Error marshaling args in ThenCall: "
	ErrMarshalingArgsInJoin        = "Error marshaling args in Join: "
	ErrGettingCurrentAccountInJoin = "Error getting current account in Join: "
	ErrMarshalingArgsInCall        = "Error marshaling args in FunctionCall: "
	ErrMarshalingArgsInCrossCall   = "Error marshaling args in Call: "
	ErrFailedToGetPromiseResult    = "failed to get promise result at index "
	ErrNoPromiseResults            = "no promise results available"
	ErrPromiseResultNotReady       = "promise result not ready"
	ErrFailedToReadPromiseResult   = "failed to read successful promise result: "
	ErrUnknownPromiseResultStatus  = "unknown promise result status: "
	ErrCallbackOnly                = "this function can only be called as a callback"
	ErrCallbackFromSelfOnly        = "callbacks can only be called by the contract itself"
	ErrRegisterEmpty               = "Register is empty, returning empty byte array"
	ErrStandardRegisterReadFailed  = "Standard register read failed: "
	ErrSuccessfullyReadBytes       = "Successfully read "
	ErrBytesFromRegister           = " bytes from register"
	ErrDirectRegisterReadCompleted = "Direct register read completed for "
	ErrFailedToCreateBuffer        = "failed to create buffer for register read"
	ErrRegisterLength              = "Register "
	ErrLength                      = " length: "
)

type BatchAction string

const (
	CreateAccountAction      BatchAction = "create_account"
	DeployContractAction     BatchAction = "deploy_contract"
	FunctionCallAction       BatchAction = "function_call"
	TransferAction           BatchAction = "transfer"
	StakeAction              BatchAction = "stake"
	AddKeyFullAccessAction   BatchAction = "add_key_full_access"
	AddKeyFunctionCallAction BatchAction = "add_key_function_call"
	DeleteKeyAction          BatchAction = "delete_key"
	DeleteAccountAction      BatchAction = "delete_account"
)

type PromiseResult struct {
	StatusCode int
	Data       []byte
	Success    bool
}

func NewPromiseResult(statusCode int, data []byte) PromiseResult {
	return PromiseResult{
		StatusCode: statusCode,
		Data:       data,
		Success:    statusCode == 1,
	}
}

func (pr PromiseResult) Unwrap() ([]byte, error) {
	if !pr.Success {
		statusName := getStatusName(pr.StatusCode)
		return nil, errors.New(ErrPromiseFailedPrefix + statusName)
	}
	return pr.Data, nil
}

func (pr PromiseResult) UnwrapOr(defaultValue []byte) []byte {
	if pr.Success {
		return pr.Data
	}
	return defaultValue
}

func getStatusName(statusCode int) string {
	switch statusCode {
	case 0:
		return "NotReady"
	case 1:
		return "Successful"
	case 2:
		return "Failed"
	default:
		return "Unknown"
	}
}

type Promise struct {
	promiseID uint64
	gas       uint64
	deposit   types.Uint128
}

func NewPromise(promiseID uint64) *Promise {
	return &Promise{
		promiseID: promiseID,
		gas:       DefaultGas,
		deposit:   types.Uint128{Hi: 0, Lo: 0},
	}
}

func CreateBatch(accountID string) *PromiseBatch {
	promiseID := env.PromiseBatchCreate([]byte(accountID))
	return NewPromiseBatch(promiseID)
}

func (p *Promise) Gas(amount uint64) *Promise {
	return &Promise{
		promiseID: p.promiseID,
		gas:       amount,
		deposit:   p.deposit,
	}
}

func (p *Promise) Deposit(amount types.Uint128) *Promise {
	return &Promise{
		promiseID: p.promiseID,
		gas:       p.gas,
		deposit:   amount,
	}
}

func (p *Promise) DepositYocto(amount uint64) *Promise {
	return p.Deposit(types.U64ToUint128(amount))
}

func (p *Promise) Then(method string, args interface{}) *Promise {
	argsBytes, err := json.Marshal(args)
	if err != nil {
		env.LogString(ErrMarshalingArgsInThen + err.Error())
		return p
	}

	currentAccount, err := env.GetCurrentAccountId()
	if err != nil {
		env.LogString(ErrGettingCurrentAccountInThen + err.Error())
		return p
	}

	promiseID := env.PromiseThen(
		p.promiseID,
		[]byte(currentAccount),
		[]byte(method),
		argsBytes,
		types.Uint128{Hi: 0, Lo: 0},
		p.gas,
	)

	return NewPromise(promiseID)
}

func (p *Promise) ThenCall(contractID, method string, args interface{}) *Promise {
	argsBytes, err := json.Marshal(args)
	if err != nil {
		env.LogString(ErrMarshalingArgsInThenCall + err.Error())
		return p
	}

	promiseID := env.PromiseThen(
		p.promiseID,
		[]byte(contractID),
		[]byte(method),
		argsBytes,
		p.deposit,
		p.gas,
	)

	return NewPromise(promiseID)
}

func (p *Promise) ThenBatch(accountID string) *PromiseBatch {
	promiseID := env.PromiseBatchThen(p.promiseID, []byte(accountID))
	return NewPromiseBatch(promiseID).Gas(p.gas)
}

func (p *Promise) Join(otherPromises []*Promise, callback string, args interface{}) *Promise {
	promiseIDs := make([]uint64, len(otherPromises)+1)
	promiseIDs[0] = p.promiseID
	for i, promise := range otherPromises {
		promiseIDs[i+1] = promise.promiseID
	}

	combinedPromise := env.PromiseAnd(promiseIDs)

	argsBytes, err := json.Marshal(args)
	if err != nil {
		env.LogString(ErrMarshalingArgsInJoin + err.Error())
		return p
	}

	currentAccount, err := env.GetCurrentAccountId()
	if err != nil {
		env.LogString(ErrGettingCurrentAccountInJoin + err.Error())
		return p
	}

	promiseID := env.PromiseThen(
		combinedPromise,
		[]byte(currentAccount),
		[]byte(callback),
		argsBytes,
		types.Uint128{Hi: 0, Lo: 0},
		p.gas,
	)

	return NewPromise(promiseID)
}

func All(promises []*Promise) uint64 {
	promiseIDs := make([]uint64, len(promises))
	for i, promise := range promises {
		promiseIDs[i] = promise.promiseID
	}
	return env.PromiseAnd(promiseIDs)
}

func (p *Promise) Value() {
	env.PromiseReturn(p.promiseID)
}

type PromiseBatch struct {
	promiseID uint64
	gas       uint64
}

func NewPromiseBatch(promiseID uint64) *PromiseBatch {
	return &PromiseBatch{
		promiseID: promiseID,
		gas:       DefaultGas,
	}
}

func (pb *PromiseBatch) Gas(amount uint64) *PromiseBatch {
	return &PromiseBatch{
		promiseID: pb.promiseID,
		gas:       amount,
	}
}

func (pb *PromiseBatch) CreateAccount() *PromiseBatch {
	env.PromiseBatchActionCreateAccount(pb.promiseID)
	return pb
}

func (pb *PromiseBatch) DeployContract(code []byte) *PromiseBatch {
	env.PromiseBatchActionDeployContract(pb.promiseID, code)
	return pb
}

func (pb *PromiseBatch) FunctionCall(method string, args interface{}, amount types.Uint128, gas uint64) *PromiseBatch {
	if gas == 0 {
		gas = pb.gas
	}

	argsBytes, err := json.Marshal(args)
	if err != nil {
		env.LogString(ErrMarshalingArgsInCall + err.Error())
		return pb
	}

	env.PromiseBatchActionFunctionCall(
		pb.promiseID,
		[]byte(method),
		argsBytes,
		amount,
		gas,
	)
	return pb
}

func (pb *PromiseBatch) FunctionCallSimple(method string, args interface{}) *PromiseBatch {
	return pb.FunctionCall(method, args, types.Uint128{Hi: 0, Lo: 0}, pb.gas)
}

func (pb *PromiseBatch) Transfer(amount types.Uint128) *PromiseBatch {
	env.PromiseBatchActionTransfer(pb.promiseID, amount)
	return pb
}

func (pb *PromiseBatch) TransferYocto(amount uint64) *PromiseBatch {
	return pb.Transfer(types.U64ToUint128(amount))
}

func (pb *PromiseBatch) Stake(amount types.Uint128, publicKey []byte) *PromiseBatch {
	env.PromiseBatchActionStake(pb.promiseID, amount, publicKey)
	return pb
}

func (pb *PromiseBatch) AddFullAccessKey(publicKey []byte, nonce uint64) *PromiseBatch {
	env.PromiseBatchActionAddKeyWithFullAccess(pb.promiseID, publicKey, nonce)
	return pb
}

func (pb *PromiseBatch) AddAccessKey(publicKey []byte, allowance types.Uint128, receiverID string, methodNames []string, nonce uint64) *PromiseBatch {
	methodsStr := ""
	for i, method := range methodNames {
		if i > 0 {
			methodsStr += ","
		}
		methodsStr += method
	}

	env.PromiseBatchActionAddKeyWithFunctionCall(
		pb.promiseID,
		publicKey,
		nonce,
		allowance,
		[]byte(receiverID),
		[]byte(methodsStr),
	)
	return pb
}

func (pb *PromiseBatch) DeleteKey(publicKey []byte) *PromiseBatch {
	env.PromiseBatchActionDeleteKey(pb.promiseID, publicKey)
	return pb
}

func (pb *PromiseBatch) DeleteAccount(beneficiaryID string) *PromiseBatch {
	env.PromiseBatchActionDeleteAccount(pb.promiseID, []byte(beneficiaryID))
	return pb
}

func (pb *PromiseBatch) Then(accountID string) *PromiseBatch {
	promiseID := env.PromiseBatchThen(pb.promiseID, []byte(accountID))
	return NewPromiseBatch(promiseID).Gas(pb.gas)
}

func (pb *PromiseBatch) Value() {
	env.PromiseReturn(pb.promiseID)
}

type CrossContract struct {
	accountID string
	gas       uint64
	deposit   types.Uint128
}

func NewCrossContract(accountID string) *CrossContract {
	return &CrossContract{
		accountID: accountID,
		gas:       DefaultGas,
		deposit:   types.Uint128{Hi: 0, Lo: 0},
	}
}

func (cc *CrossContract) Gas(amount uint64) *CrossContract {
	return &CrossContract{
		accountID: cc.accountID,
		gas:       amount,
		deposit:   cc.deposit,
	}
}

func (cc *CrossContract) Deposit(amount types.Uint128) *CrossContract {
	return &CrossContract{
		accountID: cc.accountID,
		gas:       cc.gas,
		deposit:   amount,
	}
}

func (cc *CrossContract) DepositYocto(amount uint64) *CrossContract {
	return cc.Deposit(types.U64ToUint128(amount))
}

func (cc *CrossContract) Call(method string, args interface{}) *Promise {

	argsBytes, err := json.Marshal(args)

	if err != nil {
		env.LogString(ErrMarshalingArgsInCall + err.Error())
		return NewPromise(0)
	}

	promiseID := env.PromiseCreate(
		[]byte(cc.accountID),
		[]byte(method),
		argsBytes,
		cc.deposit,
		cc.gas,
	)

	return NewPromise(promiseID).Gas(cc.gas)
}

func (cc *CrossContract) Batch() *PromiseBatch {
	promiseID := env.PromiseBatchCreate([]byte(cc.accountID))
	return NewPromiseBatch(promiseID).Gas(cc.gas)
}

func GetPromiseResult(index uint64) (PromiseResult, error) {
	if env.PromiseResultsCount() == 0 {
		return PromiseResult{}, errors.New(ErrNoPromiseResults)
	}

	data, err := env.PromiseResult(index)
	if err != nil {
		if err.Error() == "(PROMISE_ERROR): promise execution failed with data" {
			return NewPromiseResult(2, data), nil
		} else if err.Error() == "(PROMISE_ERROR): promise execution failed" {
			return NewPromiseResult(2, []byte{}), nil
		}
		return PromiseResult{}, err
	}

	return NewPromiseResult(1, data), nil
}

func GetPromiseResultSafe(index uint64) (PromiseResult, error) {
	count := env.PromiseResultsCount()
	if count == 0 {
		return PromiseResult{}, errors.New(ErrNoPromiseResults)
	}

	if index >= count {
		return PromiseResult{}, errors.New(ErrFailedToGetPromiseResult + types.IntToString(int(index)))
	}

	status, data, err := GetPromiseResultWithStatus(index)
	if err != nil {
		return PromiseResult{}, err
	}

	return NewPromiseResult(int(status), data), nil
}

func GetPromiseResultWithStatus(index uint64) (uint64, []byte, error) {
	status := env.NearBlockchainImports.PromiseResult(index, env.AtomicOpRegister)

	switch status {
	case 0:
		return status, nil, errors.New(ErrPromiseResultNotReady)
	case 1:
		data, err := ReadRegisterSafeWithFallback(env.AtomicOpRegister)
		if err != nil {
			return status, nil, errors.New(ErrFailedToReadPromiseResult + err.Error())
		}
		return status, data, nil
	case 2:
		data, err := ReadRegisterSafeWithFallback(env.AtomicOpRegister)
		if err != nil {
			return status, []byte{}, nil
		}
		return status, data, nil
	default:
		return status, nil, errors.New(ErrUnknownPromiseResultStatus + types.IntToString(int(status)))
	}
}

func GetAllPromiseResults() ([]PromiseResult, error) {
	count := env.PromiseResultsCount()
	if count == 0 {
		return nil, errors.New(ErrNoPromiseResults)
	}

	results := make([]PromiseResult, count)
	for i := uint64(0); i < count; i++ {
		result, err := GetPromiseResult(i)
		if err != nil {
			return nil, errors.New(ErrFailedToGetPromiseResult + types.IntToString(int(i)) + ": " + err.Error())
		}
		results[i] = result
	}

	return results, nil
}

func IsCallback() bool {
	return env.PromiseResultsCount() > 0
}

func IsCallbackFromSelf() bool {
	if !IsCallback() {
		return false
	}

	predecessor, err := env.GetPredecessorAccountID()
	if err != nil {
		return false
	}

	current, err := env.GetCurrentAccountId()
	if err != nil {
		return false
	}

	return predecessor == current
}

func CallbackGuard() error {
	if !IsCallback() {
		return errors.New(ErrCallbackOnly)
	}

	if !IsCallbackFromSelf() {
		return errors.New(ErrCallbackFromSelfOnly)
	}

	return nil
}

func ReadRegisterSafeWithFallback(registerId uint64) ([]byte, error) {
	length := env.NearBlockchainImports.RegisterLen(registerId)
	env.LogString(ErrRegisterLength + types.IntToString(int(registerId)) + ErrLength + types.IntToString(int(length)))

	if length == 0 {
		env.LogString(ErrRegisterEmpty)
		return []byte{}, nil
	}

	data, err := env.ReadRegisterSafe(registerId)
	if err != nil {
		env.LogString(ErrStandardRegisterReadFailed + err.Error())

		return ReadRegisterDirect(registerId, length)
	}

	env.LogString(ErrSuccessfullyReadBytes + types.IntToString(len(data)) + ErrBytesFromRegister)
	return data, nil
}

func ReadRegisterDirect(registerId uint64, length uint64) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	buffer := make([]byte, length)
	if len(buffer) > 0 {
		ptr := uint64(uintptr(unsafe.Pointer(&buffer[0])))
		env.NearBlockchainImports.ReadRegister(registerId, ptr)
		env.LogString(ErrDirectRegisterReadCompleted + types.IntToString(int(length)) + ErrBytesFromRegister)
		return buffer, nil
	}

	return []byte{}, errors.New(ErrFailedToCreateBuffer)
}
