package main

import (
	_ "embed"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/types"
)

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
}

//Transfers & Actions

//go:export TransferToken
func TransferToken() {
	contractInput, _, _ := env.ContractInput(types.ContractInputOptions{IsRawBytes: false})
	accountId, _ := json.NewParser(contractInput).GetString("account_id")
	promiseId := env.PromiseBatchCreate([]byte(accountId))
	oneNear, _ := types.U128FromString("1000000000000000000000000")
	env.PromiseBatchActionTransfer(promiseId, oneNear)
}

//go:export FunctionCall
func FunctionCall() {
	contractInput, _, _ := env.ContractInput(types.ContractInputOptions{IsRawBytes: false})
	parsedAccountId, _ := json.NewParser(contractInput).GetString("account_id")
	accountId := []byte(parsedAccountId)
	gas := uint64(1_000_000_000_000)
	amount, _ := types.U128FromString("0")
	functionName := []byte("CreateSubaccount")
	functionName1 := []byte("CreateAccount")
	arguments := []byte("{}")
	promiseCreateIdx := env.PromiseCreate(accountId, functionName, arguments, amount, gas)
	env.PromiseThen(promiseCreateIdx, accountId, functionName1, arguments, amount, gas)
}

//go:export CreateSubaccount
func CreateSubaccount() {
	minStorage, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ
	contractInput, _, _ := env.ContractInput(types.ContractInputOptions{IsRawBytes: false})
	prefix, _ := json.NewParser(contractInput).GetString("prefix")
	currentAccountId, _ := env.GetCurrentAccountId()
	promiseId := env.PromiseBatchCreate([]byte(prefix + "." + currentAccountId))
	env.PromiseBatchActionCreateAccount(promiseId)
	env.PromiseBatchActionTransfer(promiseId, minStorage)
}

//go:export CreateAccount
func CreateAccount() {
	minStorage, _ := types.U128FromString("2000000000000000000000") //0.002Ⓝ
	contractInput, _, _ := env.ContractInput(types.ContractInputOptions{IsRawBytes: false})
	accountId, _ := json.NewParser(contractInput).GetString("account_id")
	publicKey, _ := json.NewParser(contractInput).GetString("public_key")
	builder := json.NewBuilder().
		AddString("new_account_id", accountId).
		AddString("new_public_key", publicKey)
	args := builder.Build()
	promiseId := env.PromiseBatchCreate([]byte("testnet"))
	functionName := []byte("create_account")
	gas := uint64(types.ONE_TERA_GAS) * 200
	env.PromiseBatchActionFunctionCall(promiseId, functionName, args, minStorage, gas)
}

//go:embed greeting_contract.wasm
var wasmBytes []byte

//go:export DeployContract
func DeployContract() {
	contractInput, _, _ := env.ContractInput(types.ContractInputOptions{IsRawBytes: false})
	minStorage, _ := types.U128FromString("1100000000000000000000000") // 1.1 Near
	prefix, _ := json.NewParser(contractInput).GetString("prefix")
	currentAccountId, _ := env.GetCurrentAccountId()
	promiseId := env.PromiseBatchCreate([]byte(prefix + "." + currentAccountId))
	env.PromiseBatchActionCreateAccount(promiseId)
	env.PromiseBatchActionTransfer(promiseId, minStorage)
	env.PromiseBatchActionDeployContract(promiseId, wasmBytes)
}

//go:export AddKeys
func AddKeys() {
	minStorage, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ
	contractInput, _, _ := env.ContractInput(types.ContractInputOptions{IsRawBytes: false})
	prefix, _ := json.NewParser(contractInput).GetString("prefix")
	publicKey, _ := env.GetSignerAccountPK()
	currentAccountId, _ := env.GetCurrentAccountId()
	promiseId := env.PromiseBatchCreate([]byte(prefix + "." + currentAccountId))
	env.PromiseBatchActionCreateAccount(promiseId)
	env.PromiseBatchActionTransfer(promiseId, minStorage)
	env.PromiseBatchActionAddKeyWithFullAccess(promiseId, publicKey, 0)
}

//go:export DeleteAccount
func DeleteAccount() {
	minStorage, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ
	contractInput, _, _ := env.ContractInput(types.ContractInputOptions{IsRawBytes: false})
	prefix, _ := json.NewParser(contractInput).GetString("prefix")
	currentAccountId, _ := env.GetCurrentAccountId()
	promiseId := env.PromiseBatchCreate([]byte(prefix + "." + currentAccountId))
	env.PromiseBatchActionCreateAccount(promiseId)
	env.PromiseBatchActionTransfer(promiseId, minStorage)
	env.PromiseBatchActionDeleteAccount(promiseId, []byte(currentAccountId))
}

//Cross-Contract Calls

//go:export QueryingInformation
func QueryingInformation() {
	env.LogString("QueryingInformation")
}

//go:export SendingInformation
func SendingInformation() {
	env.LogString("SendingInformation")
}

//go:export Promises
func Promises() {
	env.LogString("Promises")
}

//go:export CreatingCrossContractCall
func CreatingCrossContractCall() {
	env.LogString("CreatingCrossContractCall")
}

//go:export CallbackFunction
func CallbackFunction() {
	env.LogString("CallbackFunction")
}

//go:export MultipleFunctionsSameContract
func MultipleFunctionsSameContract() {
	env.LogString("MultipleFunctionsSameContract")
}

//go:export MultipleFunctionsDifferentContracts
func MultipleFunctionsDifferentContracts() {
	env.LogString("MultipleFunctionsDifferentContracts")
}

//go:export YieldingPromise
func YieldingPromise() {
	env.LogString("YieldingPromise")
}

//go:export SignalingResume
func SignalingResume() {
	env.LogString("SignalingResume")
}

//go:export FunctionResumes
func FunctionResumes() {
	env.LogString("FunctionResumes")
}
