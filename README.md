This is a guide for working with the alpha version of the Near SDK on Go, specifically using TinyGo for WebAssembly (WASM) compilation. Below is a breakdown of the key points, limitations, and instructions provided:

---

### **Description**
- This is an alpha version of the Near SDK for Go, built on TinyGo (https://tinygo.org/).
- It is designed to work in a blockchain environment where traditional I/O, networking, and other system-level operations are unavailable.
- The SDK provides system methods in `system.go` and a more user-friendly wrapper in `env.go`.

---

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
2. **Custom Serialization**: Consider writing a custom Borsh serializer without reflection (e.g., using code generation for serializable/deserializable structs).
3. **Optimization**: Use TinyGo's optimization flags to reduce the size of the WASM binary.

---

### **Instructions for Building and Deploying**

#### **Build**
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
near contract deploy testiktinygo.testnet use-file ./main.wasm with-init-call InitContract json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' network-config testnet sign-with-legacy-keychain send
```
- Deploys the contract to the testnet with an initialization call to `InitContract`.

#### **Test Function Calls**
1. **SetStatus**:
   ```bash
   near contract call-function as-transaction testiktinygo.testnet SetStatus text-args TestikTinyGo123 prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as testiktinygo.testnet network-config testnet sign-with-keychain send
   ```

2. **GetStatus**:
   ```bash
   near contract call-function as-transaction testiktinygo.testnet GetStatus json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' sign-as testiktinygo.testnet network-config testnet sign-with-keychain send
   ```

#### **Build and Deploy in One Command**
```bash
tinygo build -size short -no-debug -panic=trap -scheduler=none -gc=leaking -o main.wasm -target wasm-unknown ./ && ls -lh main.wasm && near contract deploy testiktinygo.testnet use-file ./main.wasm with-init-call InitContract json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' network-config testnet sign-with-legacy-keychain send
```

---

### **TODO List**
1. **Mock Env Package and Test IT**: Mock env pacakge and provide isoladet VM-like env for testing env methods
2. **Add More Types**: Introduce additional types to simplify development.
3. **Write More Tests**: Expand unit and integration tests.

---

### **Additional Notes**
- **Testing and Debugging**: Since the environment is restrictive, thorough testing is crucial. Use the provided test function calls to validate your contract.
- **Optimization**: Always aim for the smallest possible WASM binary to reduce gas costs and improve performance.
- **Community Contributions**: If you fork or modify unsupported packages, consider contributing back to the community.