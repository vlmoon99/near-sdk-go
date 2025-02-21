# **Near-GO-SDK**

## **Community**

[![Telegram](https://img.shields.io/badge/Telegram-join%20chat-blue.svg)](https://t.me/go_near_sdk)
[![Discord](https://img.shields.io/badge/Discord-join%20chat-blue.svg)](https://discord.gg/UBUPuBm2)
[![Pkg Go Dev](https://img.shields.io/badge/Pkg%20Go%20Dev-view%20docs-blue.svg)](https://pkg.go.dev/github.com/vlmoon99/near-sdk-go)

## **Description**
- **This is an alpha version of the Near SDK for Go, built on [TinyGo](https://tinygo.org/).**
- **It is designed to work in a blockchain environment where traditional I/O, networking, and other system-level operations are unavailable.**
- **The SDK provides system methods in `system/system_near.go` and a more user-friendly wrapper in `env/env.go`.**

## **Tutorial**

### 🚨 **IMPORTANT PREREQUISITES** 🚨

**Before creating a smart contract in Go, make sure you have the following tools installed on your PC:**

1. **[TinyGo](https://tinygo.org/getting-started/install/)** - **_Required for building smart contracts._**
2. **[near-cli-rs](https://github.com/near/near-cli-rs)** - **_Required for interacting with the NEAR network._**

⚠️ **Ensure these tools are installed to avoid errors!** ⚠️

*In this tutorial, I will explain how to create a smart contract called "Status Message" in Go from scratch using raw TinyGo and NEAR CLI RS. Additionally, to simplify the development process, you can use the [NEAR Go CLI](https://github.com/vlmoon99/near-cli-go), which provides simplified versions of all commands for creating a project, building, creating a developer account, deploying on the testnet, importing a production account, and deploying it to production.*


## **Project Creation**

### **1. Without CLI**

1. Create a directory and navigate into it:
    ```bash
    mkdir status_messages
    cd status_messages
    ```
2. Initialize the project module:
    ```bash
    go mod init github.com/{your-github-account-id}/status_messages
    ```
3. Get the required SDK:
    ```bash
    go get github.com/vlmoon99/near-sdk-go@v0.0.7
    ```
4. Create the `main.go` file and add the following code:
    ```bash
    touch main.go && echo 'package main

    import (
        "github.com/vlmoon99/near-sdk-go/env"
    )

    //go:export InitContract
    func InitContract() {
        env.LogString("Init Smart Contract")
    }' > main.go
    ```

### **2. With CLI**

1. Create the project using the NEAR Go CLI:
    ```bash
    near-go create -p "status_messages" -m "github.com/{your-github-account-id}/status_messages"
    ```

## **Project Building Process**

### **1. Without CLI**

1. Build the project:
    ```bash
    tinygo build -size short -no-debug -panic=trap -scheduler=none -gc=leaking -o main.wasm -target wasm-unknown ./ && ls -lh main.wasm
    ```
2. If the build is successful, it will complete 100% correctly. However, if you encounter errors like this one:
    ```bash
    ../../../../go/pkg/mod/github.com/vlmoon99/near-sdk-go@v0.0.7/system/system_near.go:17:7:
    ```
    This error is not a real problem but a temporary bug that will be fixed in the future. In such cases, you need to run the build command again.

### **2. With CLI**

1. Simply call the following command:
    ```bash
    near-go build
    ```
   Cli will handle all other logic under the hood.


## **Tests**

Before running tests, let's write a simple logic inside this smart contract:

```go
package main

import (
	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/types"
)

// This type represents the internal state which we are using in our smart contract.
// We can have stored state on the blockchain or we can simply store it inside the wasm contract in memory.
// In this case, I store it inside memory and do not save my state inside the blockchain, because we have a very simple structure and type.
type StatusMessage struct {
	// Represents a proxy collection to the original methods inside /near-sdk-go/env/env.go file such as StorageWrite, StorageRead, StorageRemove, StorageHasKey.
	// Before using any collections or top-level abstractions, it is highly recommended to learn how env methods work.
	Data *collections.LookupMap
}

func GetState() StatusMessage {
	return StatusMessage{
		// []byte("b") - represents a prefix which will be added for each key inside this collection. So if I put a key with the name "test", in the blockchain I will have (b + test) as the key,
		// but the value remains the same.
		Data: collections.NewLookupMap([]byte("b")),
	}
}

// //go:export - This is a commentary which we need to use in order to export our functions to the smart contract clients. If we do not mark our methods as //go:export, we cannot call them after deployment.
// If we mark our methods with this commentary, it will be exported in our wasm file and will be visible to our clients.
// Exported functions cannot have any input and output parameters. For input from the user side, we need to use the env.ContractInput method to receive user input.
// For output, we need to use the env.ContractValueReturn function in order to provide the return value to the user.

//go:export SetStatus
func SetStatus() {
	accountId, _ := env.GetPredecessorAccountID()
	options := types.ContractInputOptions{IsRawBytes: false}
	contractInput, _, _ := env.ContractInput(options)
	parser := json.NewParser(contractInput)
	message, _ := parser.GetString("message")
	state := GetState()
	state.Data.Insert([]byte(accountId), string(message))
	env.LogString("Message : " + message + " was insterted")
	env.ContractValueReturn([]byte(message))
}

//go:export GetStatus
func GetStatus() {
	options := types.ContractInputOptions{IsRawBytes: false}
	contractInput, _, _ := env.ContractInput(options)
	parser := json.NewParser(contractInput)
	accountId, _ := parser.GetString("account_id")
	state := GetState()
	val, _ := state.Data.Get([]byte(accountId))
	status, _ := val.(string)
	env.LogString("Status : " + status + " on account id : " + accountId)
	env.ContractValueReturn([]byte(status))
}

```

For now we have some unit test and mocks for system and env methods in SDK which u can find in 
near-sdk-go/system/system_mock.go, 
near-sdk-go/system/system_mock_test.go,
near-sdk-go/env/env_test.go.

And sometimes you also can use them , but for some functionality it will be better to create your own mocks of system methods as in near-sdk-go/system/system_mock.go
and init this system into enviroment like here

```go
func init() {
	SetEnv(system.NewMockSystem())
}
```

For this Smart Contract we will be using standart mocks for unit tests and  [NEAR Workspaces](https://github.com/near/near-workspaces-rs) for integration tests in emulated Blockchain enviroment.

Let's start from the unit tests which we will write in main_test.go :

```go
package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	systemMock := system.NewMockSystem()
	env.SetEnv(systemMock)
}

func TestSetStatus(t *testing.T) {
	accountId := "test_account"
	message := "Hello, NEAR!"

	systemMock := env.NearBlockchainImports.(*system.MockSystem)
	systemMock.SetPredecessorAccountID(accountId)
	systemMock.SetContractInput([]byte(`{"message": "` + message + `"}`))

	SetStatus()

	state := GetState()
	storedMessageInterface, err := state.Data.Get([]byte(accountId))
	if err != nil {
		t.Fatalf("Failed to get stored message: %v", err)
	}

	storedMessage, ok := storedMessageInterface.(string)
	if !ok {
		t.Fatalf("Stored message is not a string")
	}

	if string(storedMessage) != message {
		t.Fatalf("Expected message %v, got %v", message, string(storedMessage))
	}
}

func TestGetStatus(t *testing.T) {
	accountId := "test_account"
	message := "Hello, NEAR!"

	systemMock := env.NearBlockchainImports.(*system.MockSystem)
	systemMock.SetPredecessorAccountID(accountId)
	state := GetState()
	state.Data.Insert([]byte(accountId), []byte(message))

	systemMock.SetContractInput([]byte(`{"account_id": "` + accountId + `"}`))

	GetStatus()

	storedMessageInterface, err := state.Data.Get([]byte(accountId))
	if err != nil {
		t.Fatalf("Failed to get stored message: %v", err)
	}

	storedMessage, ok := storedMessageInterface.(string)
	if !ok {
		t.Fatalf("Stored message is not a string")
	}

	if string(storedMessage) != message {
		t.Fatalf("Expected message %v, got %v", message, string(storedMessage))
	}
}

```


Here's the standardized and prettified version of your documentation, following the example format:

---

## **Unit Testing Process**

### **1. Run Unit Tests Without CLI**

1. To run unit tests, use the following commands:
    ```bash
    tinygo test ./
    ```
    - To run unit tests inside your package.

    ```bash
    tinygo test ./...
    ```
    - To run unit tests in all your packages.

2. If the tests are successful, they will complete 100% correctly. However, if you encounter errors like this one:
    ```bash
    ../../../../go/pkg/mod/github.com/vlmoon99/near-sdk-go@v0.0.7/system/system_near.go:17:7:
    ```
    This error is not a real problem but a temporary bug that will be fixed in the future. In such cases, you need to run the test command again.

### **2. Run Unit Tests With CLI**

1. Simply call the following commands:
    ```bash
    near-go test package
    ```
    - To run unit tests inside your package.

    ```bash
    near-go test project
    ```
    - To run unit tests inside your project.

   The CLI will handle all other logic under the hood.

---

This should make your documentation look more polished and consistent. Let me know if you need any further adjustments!



## **Integration Testing**

After we test our functionality in a mocked environment, we can start to write an integration test. For this task, we will use [Near Workspaces RS](https://github.com/near/near-workspaces-rs). First, we need to create an empty Rust project, add the necessary dependencies, and write our tests:

```bash
mkdir integration_tests
cd integration_tests
cargo init --bin

echo '[package]
name = "integration_tests"
version = "0.1.0"
edition = "2021"

[dependencies]
anyhow = "1.0.93"
json-patch = "3.0.1"
near-workspaces = "0.15.0"
serde = "1.0.215"
serde_json = "1.0.133"
tokio = "1.41.1"
near-gas = "0.3.0"' > Cargo.toml
```

After we initialize the project, we need to add the code for our integration tests for the "Status Message" Smart Contract. Let's start with the helper functions:

```rust
use near_gas::NearGas;
use near_workspaces::types::NearToken;
use serde_json::json;

async fn deploy_contract(
    worker: &near_workspaces::Worker<near_workspaces::network::Sandbox>,
) -> anyhow::Result<near_workspaces::Contract> {
    const WASM_FILEPATH: &str = "../main.wasm";
    let wasm = std::fs::read(WASM_FILEPATH)?;
    let contract = worker.dev_deploy(&wasm).await?;
    Ok(contract)
}

async fn call_integration_test_function(
    contract: &near_workspaces::Contract,
    function_name: &str,
    args: serde_json::Value,
    deposit: NearToken,
    gas: NearGas,
) -> anyhow::Result<()> {
    let outcome = contract
        .call(function_name)
        .args_json(args)
        .deposit(deposit)
        .gas(gas)
        .transact()
        .await;

    match outcome {
        Ok(result) => {
            let result_value: i8 = result.clone().json::<i8>()?;
            if result_value == 1 {
                println!("{} result: Test succeeded", function_name);
                println!("Result: Functions Logs: {:#?}", result.logs());
                Ok(())
            } else {
                println!("{} result: Test failed", function_name);
                Ok(())
            }
        }
        Err(err) => {
            println!(
                "{} result: Test failed with error: {:#?}",
                function_name, err
            );
            Err(err.into())
        }
    }
}
```

After we have this in place, we can start to add integration test calls to our smart contract:

```rust
#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let worker = near_workspaces::sandbox().await?;
    let contract = deploy_contract(&worker).await?;
    let standard_deposit = NearToken::from_near(3);
    let standard_gas = NearGas::from_tgas(300);
    println!("Dev Account ID: {}", contract.id());

    let smart_contract_calls = vec![
        call_integration_test_function(
            &contract,
            "SetStatus",
            json!({ "message": "testInputValue" }),
            standard_deposit,
            standard_gas,
        )
        .await,
        call_integration_test_function(
            &contract,
            "GetStatus",
            json!({ "account_id": contract.id() }),
            standard_deposit,
            standard_gas,
        )
        .await,
    ];

    for result in smart_contract_calls {
        if let Err(e) = result {
            eprintln!("Error: {:?}", e);
        }
    }

    Ok(())
}
```

Let's run and see how our function will execute and check their logs:

```bash
Dev Account ID: dev-20250220113022-97536040932569
result.is_success: true
Functions Logs: [
    "Message: testInputValue was inserted",
]
result.is_success: true
Functions Logs: [
    "Status: testInputValue on account id: dev-20250220113022-97536040932569",
]
```

We can see that everything is fine, and we can easily deploy our smart contract on the development environment first. After testing on the development side, we can move to production.

## **Deployment Process**

### **1. Create a Development Account on Testnet**

In order to create an account on the testnet, you can use either of the following options:

**Without CLI:**

```bash
near account create-account sponsor-by-faucet-service your-smart-contract-account-id.testnet autogenerate-new-keypair save-to-legacy-keychain network-config testnet create
```

**With NEAR-GO CLI:**

```bash
near-go create-dev-account
```

### **2. Push Smart Contract to Testnet**

**Without CLI:**

```bash
near contract deploy your-smart-contract-account-id.testnet use-file ./main.wasm without-init-call network-config testnet sign-with-legacy-keychain send
```

**With CLI:**

```bash
near-go deploy
```

### **3. Test Smart Contract on Testnet**

On this step, we can't use NEAR-GO CLI because there are no smart contract calls. You need to do it manually for now. In our "Status Message" example, we have two functions:

**SetStatus:**

```bash
near contract call-function as-transaction your-smart-contract-account-id.testnet SetStatus json-args '{"message" : "tutorial"}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as your-smart-contract-account-id.testnet network-config testnet sign-with-legacy-keychain send
```

**GetStatus:**

```bash
near contract call-function as-read-only your-smart-contract-account-id.testnet GetStatus json-args '{"account_id":"your-smart-contract-account-id.testnet"}' network-config testnet now
```

### **4. Create Mainnet Account**

To create a mainnet account, you can use various options. For example, you can use near-cli-rs, generate a mnemonic using your own cryptography, import it to the CLI, and fund it with NEAR. However, we advise you to try web wallets of NEAR to see how it works on the client side. For example, [Meteor Wallet](https://wallet.meteorwallet.app/). After that, you can import this account.

**Without CLI:**

```bash
near account import-account using-seed-phrase "your mnemonic" --seed-phrase-hd-path m/44'/397'/0' network-config mainnet
```

**With NEAR-GO CLI:**

```bash
near-go import-mainnet-account
```

### **5. Deploy Smart Contract to Mainnet**

**Without CLI:**

```bash
near contract deploy your-smart-contract-account-id.near use-file ./main.wasm without-init-call network-config mainnet sign-with-legacy-keychain send
```

**With CLI:**

```bash
near-go deploy --prod
```

### **6. Test Smart Contract on Mainnet**

**SetStatus:**

```bash
near contract call-function as-transaction your-smart-contract-account-id.near SetStatus json-args '{"message" : "tutorial"}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as your-smart-contract-account-id.near network-config mainnet sign-with-legacy-keychain send
```

**GetStatus:**

```bash
near contract call-function as-read-only your-smart-contract-account-id.near GetStatus json-args '{"account_id":"your-smart-contract-account-id.near"}' network-config mainnet now
```


## **Benchmarks**

**Compare "Status Message" Smart Contract execution.**

- [**Near Go SDK "Status Message" Contract Example**](https://github.com/vlmoon99/near-sdk-go/blob/main/examples/status_messages/main.go)
- [**Near Rust SDK "Status Message" Contract Example**](https://github.com/near/near-sdk-rs/blob/master/examples/status-message/src/lib.rs)

**Parameters for Tests:**
1. **Size (Gas to deploy, cost to store smart contract)**
2. **Set Status (Gas, Speed of execution)**
3. **Get Status (Gas, Speed of execution)**
4. **Code simplicity (Size, Basic Abstractions)**

| Library          | Storage Used | Transaction Fee | Gas Limit & Usage  | Burnt Gas & Tokens |
| ---------------- | ------------ | --------------- | ------------------ | ------------------ |
| **Near Go SDK**  |              |                 |                    |                    |
| Deploy           | 40.35 KB     | 0.000344 Ⓝ      | 0.00 gas           | 3.44 Tgas (0%)     | 🔥565 Ggas | 0.000057 Ⓝ |
| Set Status       | -            | 0.00017 Ⓝ       | 100 Tgas           | 1.93 Tgas (1.93%)  | 🔥308 Ggas | 0.000031 Ⓝ |
| Get Status       | -            | 0.000133 Ⓝ      | 100 Tgas           | 1.55 Tgas (1.55%)  | 🔥308 Ggas | 0.000031 Ⓝ |
| **Near Rust SDK**|              |                 |                    |                    |
| Deploy           | 118.85 KB    | 0.000904 Ⓝ      | 0.00 gas           | 9.04 Tgas (0%)     | 🔥1.10 Tgas | 0.00011 Ⓝ  |
| Set Status       | -            | 0.000196 Ⓝ      | 100 Tgas           | 2.19 Tgas (2.19%)  | 🔥308 Ggas | 0.000031 Ⓝ |
| Get Status       | -            | 0.000152 Ⓝ      | 100 Tgas           | 1.74 Tgas (1.74%)  | 🔥308 Ggas | 0.000031 Ⓝ |

## **Benchmarks Summary**
**We can see that the Near Go SDK is more lightweight and faster than the Rust SDK. However, the Rust SDK can be more concise in code stroke size and simplified development with rich developer tools.**


## **TODO List**
1. **Mock Env Package and Test IT:** Improves mocks in `system/system_mock_test.go` and unit test's in `system/system_mock_test.go`, `env/env_test.go`
2. **Examples:** Add more real-world examples
3. **Smart Contract Standards:** Implement basic Smart contract standards as in [**Near RS SDK**](https://github.com/near/near-sdk-rs/tree/master/near-contract-standards)

---
