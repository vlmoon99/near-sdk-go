## **Full-Stack App Tutorial**

### 🚨 **IMPORTANT PREREQUISITES** 🚨

**Before creating a smart contract in Go, make sure you have the following tools installed on your PC:**

1. **Go:**
   - Before choosing a Go version, see the [Go Compatibility Matrix](https://tinygo.org/docs/reference/go-compat-matrix/).
   - **Highly recommended:** Use the latest supported version, Go 1.23.6.

2. **Rust:**
   - Install Rust from [Rust Lang](https://www.rust-lang.org/tools/install).
   - **_Required for running integration tests._**

3. **[TinyGo](https://tinygo.org/getting-started/install/):**
   - **_Required for building smart contracts._**

⚠️ **Ensure these tools are installed to avoid errors!** ⚠️

*In this tutorial, I will explain how to create a smart contract called "Status Message" in Go from scratch using raw TinyGo and NEAR CLI RS. Additionally, to simplify the development process, you can use the [NEAR Go CLI](https://github.com/vlmoon99/near-cli-go), which provides simplified versions of all commands for creating a project, building, creating a developer account, deploying on the testnet, importing a production account, running tests, and deploying it to production.*

---

## **Project Creation**

🚨 **Highly recommended to see how it works without CLI the first time, and after that, you can use CLI to simplify development** 🚨

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
    go get github.com/vlmoon99/near-sdk-go@v0.0.8
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
    near-go create -p "status_messages" -m "github.com/{your-github-account-id}/status_messages" -t "smart-contract-empty"
    ```

---

## **Project Building Process**

### **1. Without CLI**

1. Build the project:
    ```bash
    tinygo build -size short -no-debug -panic=trap -scheduler=none -gc=leaking -o main.wasm -target wasm-unknown ./ && ls -lh main.wasm
    ```

2. If the build is successful, it will complete 100% correctly. However, if you encounter errors like this one:
    ```bash
    ../../../../go/pkg/mod/github.com/vlmoon99/near-sdk-go@v0.0.8/system/system_near.go:17:7:
    ```
    This error is not a real problem but a temporary bug that will be fixed in the future. In such cases, you need to run the build command again.

### **2. With CLI**

1. Simply call the following command:
    ```bash
    near-go build
    ```
   The CLI will handle all other logic under the hood.

---

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

For now, we have some unit tests and mocks for system and env methods in SDK, which you can find in:

- `near-sdk-go/system/system_mock.go`
- `near-sdk-go/system/system_mock_test.go`
- `near-sdk-go/env/env_test.go`

And sometimes you can use them, but for some functionality, it will be better to create your own mocks of system methods as in `near-sdk-go/system/system_mock.go` and initialize this system into the environment like here:

```go
func init() {
	SetEnv(system.NewMockSystem())
}
```

For this Smart Contract, we will be using standard mocks for unit tests and [NEAR Workspaces](https://github.com/near/near-workspaces-rs) for integration tests in the emulated Blockchain environment.

---

### **Unit Testing Process**

#### **1. Run Unit Tests Without CLI**

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
    ../../../../go/pkg/mod/github.com/vlmoon99/near-sdk-go@v0.0.8/system/system_near.go:17:7:
    ```
    This error is not a real problem but a temporary bug that will be fixed in the future. In such cases, you need to run the test command again.

### **2. Run Unit Tests With CLI**

You can easily run unit tests with the **NEAR Go CLI** using the following commands:

#### **For Unit Tests Inside Your Package:**

```bash
near-go test package
```

#### **For Unit Tests Inside Your Project:**

```bash
near-go test project
```

The **NEAR Go CLI** will automatically handle all necessary logic under the hood.

---

### **Integration Testing**

After testing functionality in a mocked environment, we can move on to writing integration tests. But first, we need to build the WASM file. For this, we will use **Near Workspaces RS**. Start by creating a new empty Rust project, adding necessary dependencies, and writing our tests.

#### **Initialize the Rust Project:**

```bash
mkdir integration_tests && cd integration_tests && cargo init --bin
```

#### **Add Dependencies to `Cargo.toml`:**

```bash
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

#### **Helper Functions for Integration Testing:**

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
            println!("result.is_success: {:#?}", result.clone().is_success());
            println!("Functions Logs: {:#?}", result.logs());
            Ok(())
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

#### **Write Integration Tests for Smart Contract:**

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

#### **Run the Integration Tests:**

Before running integration tests, build your **main.go** file with the **TinyGo** or **NEAR Go build** command:

```bash
cd integration_tests && cargo run
```

#### **Sample Logs:**

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

You can see that the tests are passing, and everything is working correctly. After testing in the development environment, you can deploy your contract to production.

---

### **Deployment Process**

#### **1. Create a Development Account on Testnet**

To create an account on the **Testnet**, you can use either of the following methods:

#### **Without CLI:**

```bash
near account create-account sponsor-by-faucet-service your-smart-contract-account-id.testnet autogenerate-new-keypair save-to-legacy-keychain network-config testnet create
```

#### **With NEAR-GO CLI:**

```bash
near-go account create -n "testnet" -a "your-smart-contract-account-id.near"
```

---

#### **2. Push Smart Contract to Testnet**

#### **Without CLI:**

```bash
near contract deploy your-smart-contract-account-id.testnet use-file ./main.wasm without-init-call network-config testnet sign-with-legacy-keychain send
```

#### **With CLI:**

```bash
near-go deploy -id "your-smart-contract-account-id.near" -n "testnet"
```


### **3. Test Smart Contract on Testnet**

In this step, we cannot use the **NEAR Go CLI** because there are no smart contract calls available. You will need to execute the calls manually for now. In our "Status Message" example, we have two functions:

#### **Without CLI:**

##### **SetStatus:**

```bash
near contract call-function as-transaction your-smart-contract-account-id.testnet SetStatus json-args '{"message" : "tutorial"}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as your-smart-contract-account-id.testnet network-config testnet sign-with-legacy-keychain send
```

##### **GetStatus:**

```bash
near contract call-function as-read-only your-smart-contract-account-id.testnet GetStatus json-args '{"account_id":"your-smart-contract-account-id.testnet"}' network-config testnet now
```

---

#### **With CLI:**

##### **SetStatus:**

```bash
near-go call \
  --from your-account-id.testnet \
  --to your-smart-contract-account-id.testnet \
  --function SetStatus \
  --args '{"message": "tutorial"}' \
  --gas '100 Tgas' \
  --deposit '0 NEAR' \
  --network testnet
```

##### **GetStatus:**

```bash
near-go call \
  --from your-account-id.testnet \
  --to your-smart-contract-account-id.testnet \
  --function GetStatus \
  --args '{"account_id": "your-account-id.testnet"}' \
  --network testnet
```

---

### **4. Create Mainnet Account**

To create a **Mainnet** account, you can use various options. For example, you can use **near-cli-rs**, generate a mnemonic using your own cryptography, import it to the CLI, and fund it with NEAR. However, we advise you to try **web wallets** of NEAR to experience the client-side process. For example, use the [Meteor Wallet](https://wallet.meteorwallet.app/). After that, you can import this account.

#### **Without CLI:**

```bash
near account import-account
```

#### **With CLI:**

```bash
near-go account import
```

---

### **5. Deploy Smart Contract to Mainnet**

#### **Without CLI:**

```bash
near contract deploy your-smart-contract-account-id.near use-file ./main.wasm without-init-call network-config mainnet sign-with-legacy-keychain send
```

#### **With CLI:**

```bash
near-go deploy -id "your-smart-contract-account-id.near" -n "mainnet"
```

---

### **6. Test Smart Contract on Mainnet**

#### **Without CLI:**

##### **SetStatus:**

```bash
near contract call-function as-transaction your-smart-contract-account-id.near SetStatus json-args '{"message" : "tutorial"}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as your-smart-contract-account-id.near network-config mainnet sign-with-legacy-keychain send
```

##### **GetStatus:**

```bash
near contract call-function as-read-only your-smart-contract-account-id.near GetStatus json-args '{"account_id":"your-smart-contract-account-id.near"}' network-config mainnet now
```

---

#### **With CLI:**

##### **SetStatus:**

```bash
go run main.go call \
  --from your-account-id.near \
  --to your-smart-contract-account-id.near \
  --function SetStatus \
  --args '{"message": "tutorial"}' \
  --gas '100 Tgas' \
  --deposit '0 NEAR' \
  --network mainnet
```

##### **GetStatus:**

```bash
go run main.go call \
  --from your-account-id.near \
  --to your-smart-contract-account-id.near \
  --function GetStatus \
  --args '{"account_id": "your-account-id.near"}' \
  --network mainnet
```
