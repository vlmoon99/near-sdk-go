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
    go get github.com/vlmoon99/near-sdk-go@v0.0.6
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
    ../../../../go/pkg/mod/github.com/vlmoon99/near-sdk-go@v0.0.6/system/system_near.go:17:7:
    ```
    This error is not a real problem but a temporary bug that will be fixed in the future. In such cases, you need to run the build command again.

### **2. With CLI**

1. Simply call the following command:
    ```bash
    near-go build
    ```
   I will handle all other logic under the hood.


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
    env.ContractValueReturn([]byte(contractInput))
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

Let's start from the unit tests :



## **Deploy**

---
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
