package main

import (
	_ "embed"
	"strconv"

	contractBuilder "github.com/vlmoon99/near-sdk-go/contract"
	"github.com/vlmoon99/near-sdk-go/env"
	promiseBuilder "github.com/vlmoon99/near-sdk-go/promise"
	"github.com/vlmoon99/near-sdk-go/types"
)

// Transfers & Actions

// Example 1: Transfer NEAR Ⓝ
//
//go:export ExampleTransferToken
func ExampleTransferToken() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		to, err := input.JSON.GetString("to")
		if err != nil {
			return err
		}
		rawAmount, err := input.JSON.GetString("amount")
		if err != nil {
			return err
		}

		amount, _ := types.U128FromString(rawAmount)

		promiseBuilder.CreateBatch(to).
			Transfer(amount)
		return nil
	})
}

// Example 2: Function Call
//
//go:export ExampleFunctionCall
func ExampleFunctionCall() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		gas := uint64(types.ONE_TERA_GAS * 10)
		accountId := "hello-nearverse.testnet"

		promiseBuilder.NewCrossContract(accountId).
			Gas(gas).
			Call("set_greeting", map[string]string{
				"message": "howdy",
			}).
			Then("ExampleFunctionCallCallback", map[string]string{})
		return nil
	})
}

//go:export ExampleFunctionCallCallback
func ExampleFunctionCallCallback() {
	contractBuilder.HandlePromiseResult(func(result *promiseBuilder.PromiseResult) error {
		if result.Success {
			env.LogString("Callback success")
			if len(result.Data) > 0 {
				env.LogString("Result: " + string(result.Data))
			}
		} else {
			env.LogString("Callback failed")
			if len(result.Data) > 0 {
				env.LogString("Error: " + string(result.Data))
			}
		}
		return nil
	})
}

//go:export ExamplePromiseResultTesting
func ExamplePromiseResultTesting() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		gas := uint64(types.ONE_TERA_GAS * 10)
		accountId := "hello-nearverse.testnet"

		env.LogString("Starting promise result testing example")

		promiseBuilder.NewCrossContract(accountId).
			Gas(gas).
			Call("get_greeting", map[string]string{}).
			Then("ExamplePromiseResultTestingCallback", map[string]string{
				"test_type": "comprehensive",
			})
		return nil
	})
}

//go:export ExamplePromiseResultTestingCallback
func ExamplePromiseResultTestingCallback() {
	contractBuilder.HandlePromiseResult(func(result *promiseBuilder.PromiseResult) error {
		env.LogString("=== Promise Result Testing Callback ===")

		// Log basic callback information
		count := env.PromiseResultsCount()
		env.LogString("Promise results count: " + types.IntToString(int(count)))

		env.LogString("--- Test Result Details ---")
		env.LogString("Success: " + strconv.FormatBool(result.Success))
		env.LogString("Status Code: " + types.IntToString(int(result.StatusCode)))

		if len(result.Data) > 0 {
			env.LogString("Data: " + string(result.Data))
		} else {
			env.LogString("No data returned")
		}

		env.LogString("=== Promise Result Testing Complete ===")
		return nil
	})
}

// Example 3: Create a Sub Account
//
//go:export ExampleCreateSubaccount
func ExampleCreateSubaccount() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		prefix, err := input.JSON.GetString("prefix")
		if err != nil {
			return err
		}

		currentAccountId, _ := env.GetCurrentAccountId()
		subaccountId := prefix + "." + currentAccountId
		amount, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ

		promiseBuilder.CreateBatch(subaccountId).
			CreateAccount().
			Transfer(amount)
		return nil
	})
}

// Example 4: Creating .testnet / .near Accounts
//
//go:export ExampleCreateAccount
func ExampleCreateAccount() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		accountId, err := input.JSON.GetString("account_id")
		if err != nil {
			return err
		}

		publicKey, err := input.JSON.GetString("public_key")
		if err != nil {
			return err
		}

		amount, _ := types.U128FromString("2000000000000000000000") //0.002Ⓝ
		gas := uint64(200 * types.ONE_TERA_GAS)

		//publicKey (base58) - EG7JhmQybCXrbXiitxsCNStPoLwakvFjgHGCNf1Wwfnt (generate your own for testing)
		//accountId - nearsdkdocs.testnet (write your own for testing)

		createArgs := map[string]string{
			"new_account_id": accountId,
			"new_public_key": publicKey,
		}

		promiseBuilder.CreateBatch("testnet").
			FunctionCall("create_account", createArgs, amount, gas)
		return nil
	})
}

//go:embed status_message_go.wasm
var contractWasm []byte

// Example 5: Deploying a Contract
//
//go:export ExampleDeployContract
func ExampleDeployContract() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		prefix, err := input.JSON.GetString("prefix")
		if err != nil {
			return err
		}

		currentAccountId, _ := env.GetCurrentAccountId()
		subaccountId := prefix + "." + currentAccountId
		amount, _ := types.U128FromString("1100000000000000000000000") //1.1Ⓝ

		promiseBuilder.CreateBatch(subaccountId).
			CreateAccount().
			Transfer(amount).
			DeployContract(contractWasm)
		return nil
	})
}

// Example 6: Add Keys to Subaccount
//
//go:export ExampleAddKeys
func ExampleAddKeys() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		prefix, err := input.JSON.GetString("prefix")
		if err != nil {
			return err
		}
		publicKey, err := input.JSON.GetString("public_key")
		if err != nil {
			return err
		}
		currentAccountId, _ := env.GetCurrentAccountId()
		subaccountId := prefix + "." + currentAccountId
		amount, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ

		promiseBuilder.CreateBatch(subaccountId).
			CreateAccount().
			Transfer(amount).
			AddFullAccessKey([]byte(publicKey), 0)
		return nil
	})
}

// Example 7: Delete Account
//
//go:export ExampleCreateDeleteAccount
func ExampleCreateDeleteAccount() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		prefix, err := input.JSON.GetString("prefix")
		if err != nil {
			return err
		}

		beneficiary, err := input.JSON.GetString("beneficiary")
		if err != nil {
			return err
		}
		currentAccountId, _ := env.GetCurrentAccountId()
		subaccountId := prefix + "." + currentAccountId
		amount, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ

		promiseBuilder.CreateBatch(subaccountId).
			CreateAccount().
			Transfer(amount).
			DeleteAccount(beneficiary)
		return nil
	})
}

// Example 7: Delete Account
//
//go:export ExampleSelfDeleteAccount
func ExampleSelfDeleteAccount() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {

		beneficiary, err := input.JSON.GetString("beneficiary")
		if err != nil {
			return err
		}
		currentAccountId, _ := env.GetCurrentAccountId()

		promiseBuilder.CreateBatch(currentAccountId).
			DeleteAccount(beneficiary)
		return nil
	})
}

// Transfers & Actions

// Cross-Contract Calls

// Example 8: Cross-Contract Query with Callback

//go:export ExampleQueryingGreetingInfo
func ExampleQueryingGreetingInfo() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		helloAccount := "hello-nearverse.testnet"
		gas := uint64(10 * types.ONE_TERA_GAS)

		promiseBuilder.NewCrossContract(helloAccount).
			Gas(gas).
			Call("get_greeting", map[string]string{}).
			Value()
		return nil
	})
}

//go:export ExampleQueryingInformation
func ExampleQueryingInformation() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		helloAccount := "hello-nearverse.testnet"
		gas := uint64(10 * types.ONE_TERA_GAS)

		promiseBuilder.NewCrossContract(helloAccount).
			Gas(gas).
			Call("get_greeting", map[string]string{}).
			Then("ExampleQueryingInformationResponse", map[string]string{})
		return nil
	})
}

//go:export ExampleQueryingInformationResponse
func ExampleQueryingInformationResponse() {
	contractBuilder.HandlePromiseResult(func(result *promiseBuilder.PromiseResult) error {
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
		return nil
	})
}

// Example 9: Improved Cross-Contract Call with State Change
//
//go:export ExampleSendingInformation
func ExampleSendingInformation() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		helloAccount := "hello-nearverse.testnet"
		gas := uint64(30 * types.ONE_TERA_GAS)

		args := map[string]string{
			"message": "New Greeting",
		}

		promiseBuilder.NewCrossContract(helloAccount).
			Gas(gas).
			Call("set_greeting", args).
			Then("ExampleChangeGreetingCallback", map[string]string{})
		return nil
	})
}

// Example 9 Response: Improved State Change Handling
//
//go:export ExampleChangeGreetingCallback
func ExampleChangeGreetingCallback() {
	contractBuilder.HandlePromiseResult(func(result *promiseBuilder.PromiseResult) error {
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
		return nil
	})
}

// Example 10: High-Level Cross-Contract API Equivalent in Go
//
//go:export ExampleCrossContractCall
func ExampleCrossContractCall() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		externalAccount := "hello-nearverse.testnet"
		gas := uint64(5 * types.ONE_TERA_GAS)

		args := map[string]string{
			"message": "New Greeting",
		}

		promiseBuilder.NewCrossContract(externalAccount).
			Gas(gas).
			Call("set_greeting", args).
			Then("ExampleCrossContractCallback", map[string]string{
				"context_data": "saved_for_callback",
			}).
			Value()
		return nil
	})
}

// Example 10: Callback Handling
//
//go:export ExampleCrossContractCallback
func ExampleCrossContractCallback() {
	contractBuilder.HandlePromiseResult(func(result *promiseBuilder.PromiseResult) error {
		env.LogString("Executing callback")

		if result.Success {
			env.LogString("Cross-contract call executed successfully")
		} else {
			env.LogString("Cross-contract call failed")
		}
		return nil
	})
}

// Example 11: Batch Calls - Multiple Functions, Same Contract
//
//go:export ExampleBatchCallsSameContract
func ExampleBatchCallsSameContract() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		helloAccount := "hello-nearverse.testnet"
		gas := uint64(10 * types.ONE_TERA_GAS)
		amount, _ := types.U128FromString("0")
		promiseBuilder.NewCrossContract(helloAccount).
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
		return nil
	})
}

// Example 11: Batch Callback Handling
//
//go:export ExampleBatchCallsCallback
func ExampleBatchCallsCallback() {
	contractBuilder.HandlePromiseResult(func(result *promiseBuilder.PromiseResult) error {
		env.LogString("Processing batch call results")
		env.LogString("Batch call success: " + strconv.FormatBool(result.Success))
		if len(result.Data) > 0 {
			env.LogString("Batch call data: " + string(result.Data))
		}
		return nil
	})
}

// Example 12: Parallel Calls - Multiple Functions, Different Contracts
//
//go:export ExampleParallelCallsDifferentContracts
func ExampleParallelCallsDifferentContracts() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		contractA := "hello-nearverse.testnet"
		contractB := "statusmessage.neargocli.testnet"

		promiseA := promiseBuilder.NewCrossContract(contractA).
			Call("get_greeting", map[string]string{})

		promiseB := promiseBuilder.NewCrossContract(contractB).
			Call("SetStatus", map[string]string{"message": "Hello, World!"})

		// Join the promises and assign a callback
		promiseA.Join([]*promiseBuilder.Promise{promiseB}, "ExampleParallelContractsCallback", map[string]string{
			"contract_ids": contractA + "," + contractB,
		}).Value()

		env.LogString("Parallel contract calls initialized")
		return nil
	})
}

// Example 12: Handling Results from Parallel Calls
//
//go:export ExampleParallelContractsCallback
func ExampleParallelContractsCallback() {
	env.LogString("Processing results from multiple contracts")

	count := env.PromiseResultsCount()
	for i := uint64(0); i < count; i++ {
		contractBuilder.HandlePromiseResult(func(result *promiseBuilder.PromiseResult) error {
			env.LogString("Processing result " + types.IntToString(int(i)))
			env.LogString("Success: " + strconv.FormatBool(result.Success))
			if len(result.Data) > 0 {
				env.LogString("Data: " + string(result.Data))
			}
			return nil
		})
	}

	env.LogString("Processed " + types.IntToString(int(count)) + " contract responses")
}

// Cross-Contract Calls
