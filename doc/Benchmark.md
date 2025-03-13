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
