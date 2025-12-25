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
	AccountID string `json:"account_id"`
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
// @contract:mutating
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
// @contract:mutating
func (c *Contract) ExampleFunctionCall() {
	gas := uint64(types.ONE_TERA_GAS * 10)
	accountId := "hello-nearverse.testnet"

	promise.NewCrossContract(accountId).
		Gas(gas).
		Call("set_greeting", map[string]string{
			"message": "howdy",
		}).
		Then("ExampleFunctionCallCallback", map[string]string{})
}

// @contract:mutating
func (c *Contract) ExampleFunctionCallCallback() {
	result, err := promise.GetPromiseResult(0)
	if err != nil {
		env.LogString("Callback error: " + err.Error())
		return
	}

	if result.Success {
		env.LogString("Callback success")
		if len(result.Data) > 0 {
			env.LogString("Result: " + string(result.Data))
		}
	} else {
		env.LogString("Callback failed")
	}
}

// Example: Promise Result Testing
// @contract:mutating
func (c *Contract) ExamplePromiseResultTesting() {
	gas := uint64(types.ONE_TERA_GAS * 10)
	accountId := "hello-nearverse.testnet"

	env.LogString("Starting promise result testing example")

	promise.NewCrossContract(accountId).
		Gas(gas).
		Call("get_greeting", map[string]string{}).
		Then("ExamplePromiseResultTestingCallback", map[string]string{
			"test_type": "comprehensive",
		})
}

// @contract:mutating
func (c *Contract) ExamplePromiseResultTestingCallback() {
	env.LogString("=== Promise Result Testing Callback ===")

	count := env.PromiseResultsCount()
	env.LogString("Promise results count: " + types.IntToString(int(count)))

	result, err := promise.GetPromiseResult(0)
	if err != nil {
		env.LogString("Error getting result: " + err.Error())
		return
	}

	env.LogString("--- Test Result Details ---")
	env.LogString("Success: " + strconv.FormatBool(result.Success))
	env.LogString("Status Code: " + types.IntToString(result.StatusCode))

	if len(result.Data) > 0 {
		env.LogString("Data: " + string(result.Data))
	} else {
		env.LogString("No data returned")
	}

	env.LogString("=== Promise Result Testing Complete ===")
}

// ============================================================================
// Account Management
// ============================================================================

// Example 3: Create a Sub Account
// @contract:mutating
func (c *Contract) ExampleCreateSubaccount(input SubaccountInput) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := input.Prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1000000000000000000000")

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount)
}

// Example 4: Creating .testnet / .near Accounts
// @contract:mutating
func (c *Contract) ExampleCreateAccount(input CreateAccountInput) {
	amount, _ := types.U128FromString("2000000000000000000000")
	gas := uint64(200 * types.ONE_TERA_GAS)

	createArgs := map[string]string{
		"new_account_id": input.AccountID,
		"new_public_key": input.PublicKey,
	}

	promise.CreateBatch("testnet").
		FunctionCall("create_account", createArgs, amount, gas)
}

// Example 5: Deploying a Contract
// @contract:mutating
func (c *Contract) ExampleDeployContract(input DeployContractInput) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := input.Prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1100000000000000000000000")

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount).
		DeployContract(contractWasm)
}

// Example 6: Add Keys to Subaccount
// @contract:mutating
func (c *Contract) ExampleAddKeys(input AddKeysInput) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := input.Prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1000000000000000000000")

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount).
		AddFullAccessKey([]byte(input.PublicKey), 0)
}

// Example 7: Delete Account
// @contract:mutating
func (c *Contract) ExampleCreateDeleteAccount(input DeleteAccountInput) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := input.Prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1000000000000000000000")

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
// @contract:mutating
func (c *Contract) ExampleQueryingGreetingInfo() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)

	promise.NewCrossContract(helloAccount).
		Gas(gas).
		Call("get_greeting", map[string]string{}).
		Value()
}

// @contract:mutating
func (c *Contract) ExampleQueryingInformation() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)

	promise.NewCrossContract(helloAccount).
		Gas(gas).
		Call("get_greeting", map[string]string{}).
		Then("ExampleQueryingInformationResponse", map[string]string{})
}

// @contract:mutating
func (c *Contract) ExampleQueryingInformationResponse() {
	result, err := promise.GetPromiseResult(0)
	if err != nil {
		env.LogString("Error retrieving result: " + err.Error())
		return
	}

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
// @contract:mutating
func (c *Contract) ExampleSendingInformation() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(30 * types.ONE_TERA_GAS)

	args := map[string]string{
		"message": "New Greeting",
	}

	promise.NewCrossContract(helloAccount).
		Gas(gas).
		Call("set_greeting", args).
		Then("ExampleChangeGreetingCallback", map[string]string{})
}

// @contract:mutating
func (c *Contract) ExampleChangeGreetingCallback() {
	c.ExampleQueryingInformationResponse()
}

// Example 10: High-Level Cross-Contract API
// @contract:mutating
func (c *Contract) ExampleCrossContractCall() {
	externalAccount := "hello-nearverse.testnet"
	gas := uint64(5 * types.ONE_TERA_GAS)

	args := map[string]string{
		"message": "New Greeting",
	}

	promise.NewCrossContract(externalAccount).
		Gas(gas).
		Call("set_greeting", args).
		Then("ExampleCrossContractCallback", map[string]string{
			"context_data": "saved_for_callback",
		}).
		Value()
}

// @contract:mutating
func (c *Contract) ExampleCrossContractCallback() {
	env.LogString("Executing callback")

	result, err := promise.GetPromiseResult(0)
	if err == nil && result.Success {
		env.LogString("Cross-contract call executed successfully")
	} else {
		env.LogString("Cross-contract call failed")
	}
}

// ============================================================================
// Batch & Parallel Calls
// ============================================================================

// Example 11: Batch Calls (Multiple Actions on One Contract)
// @contract:mutating
func (c *Contract) ExampleBatchCallsSameContract() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)
	amount, _ := types.U128FromString("0")

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
		FunctionCall("ExampleBatchCallsCallback", map[string]string{
			"original_data": "[Greeting One, Greeting Two]",
		}, amount, gas)

	env.LogString("Batch call created successfully")
}

// @contract:mutating
func (c *Contract) ExampleBatchCallsCallback() {
	env.LogString("Processing batch call results")

	result, err := promise.GetPromiseResult(0)
	if err != nil {
		env.LogString("Error: " + err.Error())
		return
	}

	env.LogString("Batch call success: " + strconv.FormatBool(result.Success))
	if len(result.Data) > 0 {
		env.LogString("Batch call data: " + string(result.Data))
	}
}

// Example 12: Parallel Calls (Different Contracts)
// @contract:mutating
func (c *Contract) ExampleParallelCallsDifferentContracts() {
	contractA := "hello-nearverse.testnet"
	contractB := "statusmessage.neargocli.testnet"

	promiseA := promise.NewCrossContract(contractA).
		Call("get_greeting", map[string]string{})

	promiseB := promise.NewCrossContract(contractB).
		Call("SetStatus", map[string]string{"message": "Hello, World!"})

	promiseA.Join([]*promise.Promise{promiseB}, "ExampleParallelContractsCallback", map[string]string{
		"contract_ids": contractA + "," + contractB,
	}).Value()

	env.LogString("Parallel contract calls initialized")
}

// @contract:mutating
func (c *Contract) ExampleParallelContractsCallback() {
	env.LogString("Processing results from multiple contracts")

	results, err := promise.GetAllPromiseResults()
	if err != nil {
		env.LogString("Error fetching results: " + err.Error())
		return
	}

	for i, result := range results {
		env.LogString("Processing result " + types.IntToString(i))
		env.LogString("Success: " + strconv.FormatBool(result.Success))
		if len(result.Data) > 0 {
			env.LogString("Data: " + string(result.Data))
		}
	}

	env.LogString("Processed " + types.IntToString(len(results)) + " contract responses")
}
