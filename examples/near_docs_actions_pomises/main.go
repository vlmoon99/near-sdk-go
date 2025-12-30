package main

import (
	_ "embed"
	"strconv"

	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/promise"
	"github.com/vlmoon99/near-sdk-go/types"
)

//go:embed status_message_go.wasm
var contractWasm []byte

// ============================================================================
// Input Models (DTOs)
// ============================================================================

type TransferTokenInput struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type SubaccountInput struct {
	Prefix string `json:"prefix"`
}

type CreateAccountInput struct {
	AccountId string `json:"account_id"`
	PublicKey string `json:"public_key"`
}

type DeployContractInput struct {
	Prefix string `json:"prefix"`
}

type AddKeysInput struct {
	Prefix    string `json:"prefix"`
	PublicKey string `json:"public_key"`
}

type DeleteAccountInput struct {
	Prefix      string `json:"prefix"`
	Beneficiary string `json:"beneficiary"`
}

type SelfDeleteInput struct {
	Beneficiary string `json:"beneficiary"`
}

type MessageInput struct {
	Message string `json:"message"`
}

type PromiseCallbackInputData struct {
	Data string `json:"data"`
}

// ============================================================================
// Contract State
// ============================================================================

// @contract:state
type Contract struct {
	Registry *collections.UnorderedMap[string, string]
}

// ============================================================================
// Initialization
// ============================================================================

// @contract:init
func (c *Contract) Init() string {
	c.Registry = collections.NewUnorderedMap[string, string]("registry")
	env.LogString("Examples Contract Initialized")
	return "Inited"
}

// ============================================================================
// Transfers & Actions
// ============================================================================

// Example 1: Transfer NEAR Ⓝ
// @contract:payable min_deposit=1NEAR
func (c *Contract) ExampleTransferToken(input TransferTokenInput) error {
	amount, err := types.U128FromString(input.Amount)
	if err != nil {
		return err
	}

	promise.CreateBatch(input.To).
		Transfer(amount)

	return nil
}

// Example 2: Function Call
// @contract:payable min_deposit=0.00001NEAR
func (c *Contract) ExampleFunctionCall() {
	gas := uint64(types.ONE_TERA_GAS * 10)
	accountId := "hello-nearverse.testnet"
	args := map[string]string{
		"message": "howdy",
	}
	promise.NewCrossContract(accountId).
		Gas(gas).
		Call("set_greeting", args).
		Then("example_function_call_callback", args)
}

// @contract:view
// @contract:promise_callback
func (c *Contract) ExampleFunctionCallCallback(input MessageInput, result promise.PromiseResult) MessageInput {
	env.LogString("Executing callback")
	env.LogString("Input Message : " + input.Message)

	if result.Success {
		env.LogString("Cross-contract call executed successfully")
		env.LogString("Promise Result Status --> " + strconv.FormatInt(int64(result.StatusCode), 10))
		if len(result.Data) > 0 {
			env.LogString("Batch call data: " + string(result.Data))
		}
	} else {
		env.LogString("Cross-contract call failed")
	}
	return input
}

// ============================================================================
// Account Management
// ============================================================================

// Example 3: Create a Sub Account
// @contract:payable min_deposit=0.001NEAR
func (c *Contract) ExampleCreateSubaccount(prefix string) {
	currentAccountId, err := env.GetCurrentAccountId()
	if err != nil {
		env.PanicStr("Failed to get current account")
	}

	subaccountId := prefix + "." + currentAccountId

	amount, err := types.U128FromString("1000000000000000000000") //0.001Ⓝ
	if err != nil {
		env.PanicStr("Bad amount format")
	}

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount)
}

// Example 4: Creating .testnet / .near Accounts
// @contract:payable min_deposit=0.002NEAR
func (c *Contract) ExampleCreateAccount(args CreateAccountInput) {
	amount, _ := types.U128FromString("2000000000000000000000") // 0.002 NEAR
	gas := uint64(200 * types.ONE_TERA_GAS)

	//publicKey (base58) - 4omJwNS1WbniWtbPkLYBrFwN3YLeffXCkpvriYgeLhst (generate your own for testing)
	//accountId - nearsdkdocs1.testnet (write your own for testing)
	createArgs := map[string]string{
		"new_account_id": args.AccountId,
		"new_public_key": args.PublicKey,
	}

	promise.CreateBatch("testnet").
		FunctionCall("create_account", createArgs, amount, gas)
}

// Example 5: Deploying a Contract
// @contract:payable min_deposit=1.1NEAR
func (c *Contract) ExampleDeployContract(prefix string) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1100000000000000000000000") // 1.1Ⓝ

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount).
		DeployContract(contractWasm)
}

// Example 6: Add Keys to Subaccount
// @contract:payable min_deposit=0.001NEAR
func (c *Contract) ExampleAddKeys(input AddKeysInput) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := input.Prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1000000000000000000000") // 0.001Ⓝ

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount).
		AddFullAccessKey([]byte(input.PublicKey), 0)
}

// Example 7: Delete Account
// @contract:payable min_deposit=0.001NEAR
func (c *Contract) ExampleCreateDeleteAccount(input DeleteAccountInput) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := input.Prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1000000000000000000000") // 0.001Ⓝ

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount).
		DeleteAccount(input.Beneficiary)
}

// @contract:mutating
func (c *Contract) ExampleSelfDeleteAccount(input SelfDeleteInput) {
	currentAccountId, _ := env.GetCurrentAccountId()

	promise.CreateBatch(currentAccountId).
		DeleteAccount(input.Beneficiary)
}

// ============================================================================
// Cross-Contract Calls & Queries
// ============================================================================

// Example 8: Cross-Contract Query
// @contract:payable min_deposit=0.001NEAR
func (c *Contract) ExampleQueryingGreetingInfo() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)

	promise.NewCrossContract(helloAccount).
		Gas(gas).
		Call("get_greeting", map[string]string{}).
		Value()
}

// @contract:payable min_deposit=0.001NEAR
func (c *Contract) ExampleQueryingInformation() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)

	promise.
		NewCrossContract(helloAccount).
		Gas(gas).
		Call("get_greeting", map[string]string{}).
		Then("example_querying_information_response", map[string]string{})
}

// @contract:view
// @contract:promise_callback
func (c *Contract) ExampleQueryingInformationResponse(result promise.PromiseResult) {

	if result.Success {
		env.LogString("State change/Query completed successfully")
	} else {
		env.LogString("State change/Query failed")
	}

	env.LogString("Promise result status: " + types.IntToString(result.StatusCode))
	if len(result.Data) > 0 {
		env.LogString("Returned data: " + string(result.Data))
	} else {
		env.LogString("No return data")
	}
}

// Example 9: Improved Cross-Contract Call
// @contract:payable min_deposit=0.00001NEAR
func (c *Contract) ExampleSendingInformation() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(30 * types.ONE_TERA_GAS)

	args := map[string]string{
		"message": "New Greeting",
	}

	promise.NewCrossContract(helloAccount).
		Gas(gas).
		Call("set_greeting", args).
		Then("example_change_greeting_callback", map[string]string{})
}

// @contract:view
// @contract:promise_callback
func (c *Contract) ExampleChangeGreetingCallback(result promise.PromiseResult) {
	if result.Success {
		env.LogString("State change completed successfully")
	} else {
		env.LogString("State change failed")
	}

	env.LogString("Promise result status: " + types.IntToString(int(result.StatusCode)))
	if len(result.Data) > 0 {
		env.LogString("Returned data: " + string(result.Data))
	} else {
		env.LogString("No return data from state change")
	}
}

// Example 10: High-Level Cross-Contract API
// @contract:payable min_deposit=0.00001NEAR
func (c *Contract) ExampleCrossContractCall() {
	externalAccount := "hello-nearverse.testnet"
	gas := uint64(5 * types.ONE_TERA_GAS)

	args := map[string]string{
		"message": "New Greeting",
	}
	callback_args := map[string]string{
		"data": "saved_for_callback",
	}
	promise.NewCrossContract(externalAccount).
		Gas(gas).
		Call("set_greeting", args).
		Then("example_cross_contract_callback", callback_args).
		Value()
}

// @contract:view
// @contract:promise_callback
func (c *Contract) ExampleCrossContractCallback(input PromiseCallbackInputData, result promise.PromiseResult) {
	env.LogString("Executing callback")

	env.LogString("Input CrossContractCallback : " + input.Data)

	if result.Success {
		env.LogString("Cross-contract call executed successfully")
	} else {
		env.LogString("Cross-contract call failed")
	}
}

// ============================================================================
// Batch & Parallel Calls
// ============================================================================

// Example 11: Batch Calls (Multiple Actions on One Contract)
// @contract:payable min_deposit=0.00001NEAR
func (c *Contract) ExampleBatchCallsSameContract() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)
	amount, _ := types.U128FromString("0")
	callback_args := map[string]string{
		"data": "[Greeting One, Greeting Two]",
	}

	promise.NewCrossContract(helloAccount).
		Batch().
		Gas(gas).
		FunctionCall("set_greeting", map[string]string{
			"message": "Greeting One",
		}, amount, gas).
		FunctionCall("another_method", map[string]string{
			"arg1": "val1",
		}, amount, gas).
		Then(helloAccount).
		FunctionCall("example_batch_calls_callback", callback_args, amount, gas)

	env.LogString("Batch call created successfully")
}

// @contract:view
// @contract:promise_callback
func (c *Contract) ExampleBatchCallsCallback(input PromiseCallbackInputData, result promise.PromiseResult) {
	env.LogString("Processing batch call results")
	env.LogString("Input CrossContractCallback : " + input.Data)

	env.LogString("Batch call success: " + strconv.FormatBool(result.Success))
	if len(result.Data) > 0 {
		env.LogString("Batch call data: " + string(result.Data))
	}
}

// Example 12: Parallel Calls (Different Contracts)
// @contract:payable min_deposit=0.00001NEAR
func (c *Contract) ExampleParallelCallsDifferentContracts() {
	contractA := "hello-nearverse.testnet"
	contractB := "child.neargopromises1.testnet"

	promiseA := promise.NewCrossContract(contractA).
		Call("get_greeting", map[string]string{})

	promiseB := promise.NewCrossContract(contractB).
		Call("SetStatus", map[string]string{"message": "Hello, World!"})

	promiseA.Join([]*promise.Promise{promiseB}, "example_parallel_contracts_callback", map[string]string{
		"data": contractA + "," + contractB,
	}).Value()

	env.LogString("Parallel contract calls initialized")
}

// @contract:view
// @contract:promise_callback
func (c *Contract) ExampleParallelContractsCallback(input PromiseCallbackInputData, results []promise.PromiseResult) {
	env.LogString("Processing results from multiple contracts")
	env.LogString("Input CrossContractCallback : " + input.Data)

	for i, result := range results {
		env.LogString("Processing result " + types.IntToString(i))
		env.LogString("Success: " + strconv.FormatBool(result.Success))
		if len(result.Data) > 0 {
			env.LogString("Data: " + string(result.Data))
		}
	}

	env.LogString("Processed " + types.IntToString(len(results)) + " contract responses")
}
