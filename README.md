<h1><code>Near-GO-SDK</code></h1>

### **Description**
- This is an alpha version of the Near SDK for Go, built on TinyGo (https://tinygo.org/).
- It is designed to work in a blockchain environment where traditional I/O, networking, and other system-level operations are unavailable.
- The SDK provides system methods in `system.go` and a more user-friendly wrapper in `env.go`.

---

Sure! Here's the structured and formatted text for your README:

### **Benchmarks**

Compare Status Message Smart Contract execution.

**Near Go SDK**
https://github.com/vlmoon99/near-sdk-go/blob/main/examples/status_messages/main.go

**Near Rust SDK**
https://github.com/near/near-sdk-rs/blob/master/examples/status-message/src/lib.rs

**Parameters for Tests:**
1. Size (Gas to deploy, cost to store smart contract)
2. Set Status (Gas, Speed of execution)
3. Get Status (Gas, Speed of execution)
4. Code simplicity (Size, Basic Abstractions)

| Library | Storage Used | Transaction Fee | Gas Limit & Usage | Burnt Gas & Tokens | 
| ------- | -------------| ----------------| ------------------| -------------------|
| **Near Go SDK** | | | | |
| Deploy | 40.35 KB | 0.000344 Ⓝ | 0.00 gas | 3.44 Tgas (0%) | 🔥565 Ggas | 0.000057 Ⓝ |
| Set Status | - | 0.00017 Ⓝ | 100 Tgas | 1.93 Tgas (1.93%) | 🔥308 Ggas | 0.000031 Ⓝ |
| Get Status | - | 0.000133 Ⓝ | 100 Tgas | 1.55 Tgas (1.55%) | 🔥308 Ggas | 0.000031 Ⓝ |
| **Near Rust SDK** | | | | |
| Deploy | 118.85 KB | 0.000904 Ⓝ | 0.00 gas | 9.04 Tgas (0%) | 🔥1.10 Tgas | 0.00011 Ⓝ |
| Set Status | - | 0.000196 Ⓝ | 100 Tgas | 2.19 Tgas (2.19%) | 🔥308 Ggas | 0.000031 Ⓝ |
| Get Status | - | 0.000152 Ⓝ | 100 Tgas | 1.74 Tgas (1.74%) | 🔥308 Ggas | 0.000031 Ⓝ |

### Commands to Reproduce

**Near Go SDK**

1. Deploy
```sh
near contract deploy statusmsggosdk.testnet use-file ./status_message_go.wasm without-init-call network-config testnet sign-with-legacy-keychain send
```

2. Set Status
```sh
near contract call-function as-transaction statusmsggosdk.testnet SetStatus json-args '{"message" : "Qwerty1234"}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as statusmsggosdk.testnet network-config testnet sign-with-keychain send
```

3. Get Status
```sh
near contract call-function as-transaction statusmsggosdk.testnet GetStatus json-args '{"account_id" : "statusmsggosdk.testnet"}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as statusmsggosdk.testnet network-config testnet sign-with-keychain send
```

**Near Rust SDK**

1. Deploy
```sh
near contract deploy statusmsgrssdk.testnet use-file ./status_message_rust.wasm without-init-call network-config testnet sign-with-legacy-keychain send
```

2. Set Status
```sh
near contract call-function as-transaction statusmsgrssdk.testnet set_status json-args '{"message" : "Qwerty1234"}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as statusmsgrssdk.testnet network-config testnet sign-with-keychain send
```

3. Get Status
```sh
near contract call-function as-transaction statusmsgrssdk.testnet get_status json-args '{"account_id" : "statusmsgrssdk.testnet"}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as statusmsgrssdk.testnet network-config testnet sign-with-keychain send
```

### **Benchmarks Summary**
We can see that the Near Go SDK is more lightweight and faster than the Rust SDK. However, the Rust SDK can be more concise in code stroke size and simplified development with rich developer tools. 



### **Limitations of the Environment**
1. **No I/O or Networking**: Blockchain environments do not support traditional I/O or networking operations.
2. **Supported Packages**: Only a subset of Go's standard library is supported. Refer to TinyGo's [stdlib support reference](https://tinygo.org/docs/reference/lang-support/stdlib/).
3. **Avoid Unsupported Packages**:
   - Do not use `fmt` in production builds.
   - Avoid `math/big` and any packages that depend on I/O, reflection, or encoding/json.
   - If you need a package, inspect its source code to ensure it doesn't rely on unsupported dependencies.

---

### **Tips for Building**
1. **Avoid Unsupported Packages**: Stick to packages that pass TinyGo's compatibility tests.
2. **Optimization**: Don't use fmt in production build's , use it only for debug. Create integration-tests , examples u can find in integration_tests/src/main.rs and examples/integration_tests/main.go.

---

### **Instructions for Building and Deploying**

#### **Build By**
1. **Full Optimization**:
   ```bash
   tinygo build -size short -no-debug -panic=trap -scheduler=none -gc=leaking -o main.wasm -target wasm-unknown ./ && ls -lh main.wasm
   ```
   - This minimizes the binary size and removes debugging information.

2. **Minimum Optimization**:
   ```bash
   tinygo build -size short -no-debug -o main.wasm -target wasm-unknown ./ && ls -lh main.wasm
   ```

#### **Deploy**
```bash
near contract deploy {your smart-contract account id} use-file ./main.wasm with-init-call InitContract json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' network-config testnet sign-with-legacy-keychain send
```
- Deploys the contract to the testnet with an initialization call to `InitContract`.

#### **Test Function Calls**
1. **SetStatus**:
   ```bash
   near contract call-function as-transaction {your smart-contract account id} SetStatus text-args TestikTinyGo19999 prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as {your signer account id} network-config testnet sign-with-keychain send
   ```

2. **GetStatus**:
   ```bash
   near contract call-function as-transaction {your smart-contract account id} GetStatus json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as {your signer account id} network-config testnet sign-with-keychain send
   ```

#### **Build and Deploy in One Command**
```bash
tinygo build -size short -no-debug -panic=trap -scheduler=none -gc=leaking -o main.wasm -target wasm-unknown ./ && ls -lh main.wasm && near contract deploy {your smart-contract account id} use-file ./main.wasm with-init-call InitContract json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' network-config testnet sign-with-legacy-keychain send
```

---

### **TODO List**
1. **Mock Env Package and Test IT**: Improves mocks in system/system_mock_test.go and unit test's in system/system_mock_test.go, env/env_test.go
2. **Examples**: Add more real-world examples
3. **Smart Contract Standarts**: Implement basic Smart contract standarts as in  https://github.com/near/near-sdk-rs/tree/master/near-contract-standards

---

### **Additional Notes**
- **Testing and Debugging**: Since the environment is restrictive, thorough testing is crucial. Use the provided test function calls to validate your contract.
- **Optimization**: Always aim for the smallest possible WASM binary to reduce gas costs and improve performance.
- **Community Contributions**: If you fork or modify unsupported packages, consider contributing back to the community.