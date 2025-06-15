package main

// import (
// 	_ "embed"
// )

// // Transfers & Actions

// // Example 1: Transfer NEAR Ⓝ
// //
// //go:export Example1TransferToken
// func Example1TransferToken() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		accountId, err := input.JSON.GetString("account_id")
// 		if err != nil {
// 			return err
// 		}

// 		amount, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ
// 		CreateBatch(accountId).
// 			Transfer(amount)
// 		return nil
// 	})
// }

// // Example 2: Function Call
// //
// //go:export Example2FunctionCall
// func Example2FunctionCall() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		env.LogString("client input  ===  " + string(input.Data))
// 		gas := uint64(types.ONE_TERA_GAS * 10)
// 		accountId := "hello-nearverse.testnet"

// 		NewCrossContract(accountId).
// 			Gas(gas).
// 			Call("set_greeting", map[string]string{
// 				"message": "howdy",
// 			}).
// 			Then("Example2FunctionCallCallback", map[string]string{})
// 		return nil
// 	})
// }

// //go:export Example2FunctionCallCallback
// func Example2FunctionCallCallback() {
// 	if err := CallbackGuard(); err != nil {
// 		env.LogString("Callback rejected: not from self")
// 		return
// 	}

// 	result, err := GetPromiseResultSafe(0)
// 	if err != nil {
// 		env.LogString("Promise failed: " + err.Error())
// 		return
// 	}

// 	if result.Success {
// 		env.LogString("Callback success")
// 		if len(result.Data) > 0 {
// 			env.LogString("Result: " + string(result.Data))
// 		}
// 	} else {
// 		env.LogString("Callback failed")
// 		if len(result.Data) > 0 {
// 			env.LogString("Error: " + string(result.Data))
// 		}
// 	}
// }

// //go:export Example2PromiseResultTesting
// func Example2PromiseResultTesting() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		env.LogString("client input  ===  " + string(input.Data))

// 		gas := uint64(types.ONE_TERA_GAS * 10)
// 		accountId := "hello-nearverse.testnet"

// 		env.LogString("Starting promise result testing example")

// 		NewCrossContract(accountId).
// 			Gas(gas).
// 			Call("get_greeting", map[string]string{}).
// 			Then("ExamplePromiseResultTestingCallback", map[string]string{
// 				"test_type": "comprehensive",
// 			})
// 		return nil
// 	})
// }

// //go:export Example2PromiseResultTestingCallback
// func Example2PromiseResultTestingCallback() {
// 	env.LogString("=== Promise Result Testing Callback ===")

// 	if !IsCallbackFromSelf() {
// 		env.LogString("Error: Callback not from self")
// 		return
// 	}

// 	// Log basic callback information
// 	count := env.PromiseResultsCount()
// 	env.LogString("Promise results count: " + intToString(int(count)))

// 	if count == 0 {
// 		env.LogString("No promise results available")
// 		return
// 	}

// 	// Test 1: Try the SDK method (now fixed)
// 	env.LogString("--- Test 1: SDK PromiseResult (Fixed) ---")
// 	sdkData, sdkErr := env.PromiseResult(0)
// 	if sdkErr != nil {
// 		env.LogString("SDK method failed: " + sdkErr.Error())
// 	} else {
// 		env.LogString("SDK method success, data length: " + intToString(len(sdkData)))
// 		if len(sdkData) > 0 {
// 			env.LogString("SDK method data: " + string(sdkData))
// 		} else {
// 			env.LogString("SDK method: successful promise with no return data")
// 		}
// 	}

// 	// Test 2: Try the contract wrapper method
// 	env.LogString("--- Test 2: Contract GetPromiseResult ---")
// 	result1, err1 := GetPromiseResult(0)
// 	if err1 != nil {
// 		env.LogString("Contract method failed: " + err1.Error())
// 	} else {
// 		env.LogString("Contract method success: " + boolToString(result1.Success))
// 		env.LogString("Contract method status: " + intToString(result1.StatusCode))
// 		if len(result1.Data) > 0 {
// 			env.LogString("Contract method data: " + string(result1.Data))
// 		} else {
// 			env.LogString("Contract method: successful promise with no return data")
// 		}
// 	}

// 	// Test 3: Try the safe method (fallback)
// 	env.LogString("--- Test 3: Safe GetPromiseResultSafe ---")
// 	result2, err2 := GetPromiseResultSafe(0)
// 	if err2 != nil {
// 		env.LogString("Safe method failed: " + err2.Error())
// 	} else {
// 		env.LogString("Safe method success: " + boolToString(result2.Success))
// 		env.LogString("Safe method status: " + intToString(result2.StatusCode))
// 		if len(result2.Data) > 0 {
// 			env.LogString("Safe method data: " + string(result2.Data))
// 		}
// 	}

// 	// Test 4: Try the low-level method
// 	env.LogString("--- Test 4: Low-level GetPromiseResultWithStatus ---")
// 	status3, data3, err3 := GetPromiseResultWithStatus(0)
// 	if err3 != nil {
// 		env.LogString("Low-level method failed: " + err3.Error())
// 	} else {
// 		env.LogString("Low-level method status: " + intToString(int(status3)))
// 		if len(data3) > 0 {
// 			env.LogString("Low-level method data: " + string(data3))
// 		}
// 	}

// 	// Test 5: Try direct register access
// 	env.LogString("--- Test 5: Direct Register Access ---")
// 	directStatus := env.NearBlockchainImports.PromiseResult(0, env.AtomicOpRegister)
// 	env.LogString("Direct promise result status: " + intToString(int(directStatus)))

// 	registerLength := env.NearBlockchainImports.RegisterLen(env.AtomicOpRegister)
// 	env.LogString("Register length: " + intToString(int(registerLength)))

// 	if registerLength > 0 {
// 		directData, directErr := ReadRegisterSafeWithFallback(env.AtomicOpRegister)
// 		if directErr != nil {
// 			env.LogString("Direct register read failed: " + directErr.Error())
// 		} else {
// 			env.LogString("Direct register read success, data length: " + intToString(len(directData)))
// 			if len(directData) > 0 {
// 				env.LogString("Direct register data: " + string(directData))
// 			}
// 		}
// 	}

// 	env.LogString("=== Promise Result Testing Complete ===")
// }

// // Example 3: Create a Sub Account
// //
// //go:export Example3CreateSubaccount
// func Example3CreateSubaccount() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		prefix, err := input.JSON.GetString("prefix")
// 		if err != nil {
// 			return err
// 		}

// 		currentAccountId, _ := env.GetCurrentAccountId()
// 		subaccountId := prefix + "." + currentAccountId
// 		amount, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ

// 		CreateBatch(subaccountId).
// 			CreateAccount().
// 			Transfer(amount)
// 		return nil
// 	})
// }

// // Example 4: Creating .testnet / .near Accounts
// //
// //go:export Example4CreateAccount
// func Example4CreateAccount() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		accountId, err := input.JSON.GetString("account_id")
// 		if err != nil {
// 			return err
// 		}

// 		publicKey, err := input.JSON.GetString("public_key")
// 		if err != nil {
// 			return err
// 		}

// 		amount, _ := types.U128FromString("2000000000000000000000") //0.002Ⓝ
// 		gas := uint64(200 * types.ONE_TERA_GAS)

// 		//publicKey (base58) - EG7JhmQybCXrbXiitxsCNStPoLwakvFjgHGCNf1Wwfnt (generate your own for testing)
// 		//accountId - nearsdkdocs.testnet (write your own for testing)

// 		createArgs := map[string]string{
// 			"new_account_id": accountId,
// 			"new_public_key": publicKey,
// 		}

// 		CreateBatch("testnet").
// 			FunctionCall("create_account", createArgs, amount, gas)
// 		return nil
// 	})
// }

// //go:embed status_message_go.wasm
// var contractWasm []byte

// // Example 5: Deploying a Contract
// //
// //go:export Example5DeployContract
// func Example5DeployContract() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		prefix, err := input.JSON.GetString("prefix")
// 		if err != nil {
// 			return err
// 		}

// 		currentAccountId, _ := env.GetCurrentAccountId()
// 		subaccountId := prefix + "." + currentAccountId
// 		amount, _ := types.U128FromString("1100000000000000000000000") //1.1Ⓝ

// 		CreateBatch(subaccountId).
// 			CreateAccount().
// 			Transfer(amount).
// 			DeployContract(contractWasm)
// 		return nil
// 	})
// }

// // Example 6: Add Keys to Subaccount
// //
// //go:export Example6AddKeys
// func Example6AddKeys() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		prefix, err := input.JSON.GetString("prefix")
// 		if err != nil {
// 			return err
// 		}

// 		publicKey, _ := env.GetSignerAccountPK()
// 		currentAccountId, _ := env.GetCurrentAccountId()
// 		subaccountId := prefix + "." + currentAccountId
// 		amount, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ

// 		CreateBatch(subaccountId).
// 			CreateAccount().
// 			Transfer(amount).
// 			AddFullAccessKey(publicKey, 0)
// 		return nil
// 	})
// }

// // Example 7: Delete Account
// //
// //go:export Example7DeleteAccount
// func Example7DeleteAccount() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		prefix, err := input.JSON.GetString("prefix")
// 		if err != nil {
// 			return err
// 		}

// 		currentAccountId, _ := env.GetCurrentAccountId()
// 		subaccountId := prefix + "." + currentAccountId
// 		amount, _ := types.U128FromString("1000000000000000000000") //0.001Ⓝ

// 		CreateBatch(subaccountId).
// 			CreateAccount().
// 			Transfer(amount).
// 			DeleteAccount(currentAccountId)
// 		return nil
// 	})
// }

// // Transfers & Actions

// // Cross-Contract Calls

// // Example 8: Cross-Contract Query with Callback
// //
// //go:export Example8QueryingInformation
// func Example8QueryingInformation() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		env.LogString("client input  ===  " + string(input.Data))

// 		helloAccount := "hello-nearverse.testnet"
// 		gas := uint64(10 * types.ONE_TERA_GAS)

// 		NewCrossContract(helloAccount).
// 			Gas(gas).
// 			Call("get_greeting", map[string]string{}).
// 			Then("Example8QueryingInformationResponse", map[string]string{})
// 		return nil
// 	})
// }

// //go:export Example8QueryingInformationResponse
// func Example8QueryingInformationResponse() {
// 	if !IsCallbackFromSelf() {
// 		env.LogString("Error: Callback not from self")
// 		return
// 	}
// 	env.LogString("Callback was executed successfully")
// }

// // Example 9: Improved Cross-Contract Call with State Change
// //
// //go:export Example9SendingInformation
// func Example9SendingInformation() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		env.LogString("Client input: " + string(input.Data))

// 		helloAccount := "hello-nearverse.testnet"
// 		gas := uint64(30 * types.ONE_TERA_GAS)

// 		args := map[string]string{
// 			"message": "New Greeting",
// 		}

// 		NewCrossContract(helloAccount).
// 			Gas(gas).
// 			Call("set_greeting", args).
// 			Then("Example9ChangeGreetingCallback", map[string]string{})
// 		return nil
// 	})
// }

// // Example 9 Response: Improved State Change Handling
// //
// //go:export Example9ChangeGreetingCallback
// func Example9ChangeGreetingCallback() {
// 	if !IsCallbackFromSelf() {
// 		env.LogString("Error: Callback not from self")
// 		return
// 	}

// 	result, err := GetPromiseResultSafe(0)
// 	if err != nil {
// 		env.LogString("Error getting promise result: " + err.Error())
// 		return
// 	}

// 	if result.Success {
// 		env.LogString("State change completed successfully")
// 	} else {
// 		env.LogString("State change failed")
// 	}

// 	// Additional logging for debugging purposes
// 	env.LogString("Promise result status: " + intToString(result.StatusCode))
// 	if len(result.Data) > 0 {
// 		env.LogString("Returned data: " + string(result.Data))
// 	} else {
// 		env.LogString("No return data from state change")
// 	}
// }

// // Example 10: High-Level Cross-Contract API Equivalent in Go
// //
// //go:export Example10CrossContractCall
// func Example10CrossContractCall() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		env.LogString("Client input: " + string(input.Data))

// 		externalAccount := "hello-nearverse.testnet"
// 		gas := uint64(5 * types.ONE_TERA_GAS)

// 		args := map[string]string{
// 			"message": "New Greeting",
// 		}

// 		NewCrossContract(externalAccount).
// 			Gas(gas).
// 			Call("set_greeting", args).
// 			Then("Example10CrossContractCallback", map[string]string{
// 				"context_data": "saved_for_callback",
// 			})
// 		return nil
// 	})
// }

// // Example 10: Callback Handling
// //
// //go:export Example10CrossContractCallback
// func Example10CrossContractCallback() {
// 	env.LogString("Executing callback")

// 	if !IsCallbackFromSelf() {
// 		env.LogString("Error: Callback not from self")
// 		return
// 	}

// 	count := env.PromiseResultsCount()
// 	env.LogString("Promise results count: " + intToString(int(count)))

// 	if count == 0 {
// 		env.LogString("No promise results available")
// 		return
// 	}

// 	result, err := GetPromiseResultSafe(0)
// 	if err != nil {
// 		env.LogString("Error retrieving promise result: " + err.Error())
// 		return
// 	}

// 	if result.Success {
// 		env.LogString("Cross-contract call executed successfully")
// 	} else {
// 		env.LogString("Cross-contract call failed")
// 	}
// }

// // Example 11: Batch Calls - Multiple Functions, Same Contract
// //
// //go:export Example11BatchCallsSameContract
// func Example11BatchCallsSameContract() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		env.LogString("Client input: " + string(input.Data))

// 		helloAccount := "hello-nearverse.testnet"
// 		gas := uint64(10 * types.ONE_TERA_GAS)

// 		NewCrossContract(helloAccount).
// 			Gas(gas).
// 			Call("set_greeting", map[string]string{
// 				"message": "Greeting One",
// 			}).
// 			ThenCall(helloAccount, "get_greeting", map[string]string{}).
// 			Then("Example11BatchCallsCallback", map[string]string{
// 				"original_data": "[Greeting One, Greeting Two]",
// 			})

// 		env.LogString("Batch call created successfully")
// 		return nil
// 	})
// }

// // Example 11: Batch Callback Handling
// //
// //go:export Example11BatchCallsCallback
// func Example11BatchCallsCallback() {
// 	env.LogString("Processing batch call results")

// 	if !IsCallbackFromSelf() {
// 		env.LogString("Error: Callback not from self")
// 		return
// 	}

// 	count := env.PromiseResultsCount()
// 	env.LogString("Promise results count: " + intToString(int(count)))

// 	if count == 0 {
// 		env.LogString("No promise results available")
// 		return
// 	}

// 	result, err := GetPromiseResultSafe(0)
// 	if err != nil {
// 		env.LogString("Error retrieving batch result: " + err.Error())
// 		return
// 	}

// 	env.LogString("Batch call success: " + boolToString(result.Success))
// 	env.LogString("Batch call data: " + string(result.Data))
// }

// // Example 12: Parallel Calls - Multiple Functions, Different Contracts
// //
// //go:export Example12ParallelCallsDifferentContracts
// func Example12ParallelCallsDifferentContracts() {
// 	HandleClientJSONInput(func(input *ContractInput) error {
// 		env.LogString("Client input: " + string(input.Data))

// 		contractA := "hello-nearverse.testnet"
// 		contractB := "statusmessage.neargocli.testnet"

// 		promiseA := NewCrossContract(contractA).
// 			Call("get_greeting", map[string]string{})

// 		promiseB := NewCrossContract(contractB).
// 			Call("SetStatus", map[string]string{"message": "Hello, World!"})

// 		// Join the promises and assign a callback
// 		promiseA.Join([]*Promise{promiseB}, "Example12ParallelContractsCallback", map[string]string{
// 			"contract_ids": contractA + "," + contractB,
// 		})

// 		env.LogString("Parallel contract calls initialized")
// 		return nil
// 	})
// }

// // Example 12: Handling Results from Parallel Calls
// //
// //go:export Example12ParallelContractsCallback
// func Example12ParallelContractsCallback() {
// 	env.LogString("Processing results from multiple contracts")

// 	if !IsCallbackFromSelf() {
// 		env.LogString("Error: Callback not from self")
// 		return
// 	}

// 	results := []PromiseResult{}
// 	count := env.PromiseResultsCount()
// 	for i := 0; i < int(count); i++ {
// 		result, err := GetPromiseResultSafe(uint64(i))
// 		if err != nil {
// 			env.LogString("Error processing result " + intToString(i) + ": " + err.Error())
// 			continue
// 		}
// 		results = append(results, result)
// 	}

// 	env.LogString("Processed " + intToString(len(results)) + " contract responses")
// }

// // Cross-Contract Calls
