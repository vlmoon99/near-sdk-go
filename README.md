# **Near SDK GO**

[![Telegram](https://img.shields.io/badge/Telegram-join%20chat-blue.svg)](https://t.me/go_near_sdk)  
[![Discord](https://img.shields.io/badge/Discord-join%20chat-blue.svg)](https://discord.gg/UBUPuBm2)  
[![Pkg Go Dev](https://img.shields.io/badge/Pkg%20Go%20Dev-view%20docs-blue.svg)](https://pkg.go.dev/github.com/vlmoon99/near-sdk-go)  

**Near SDK GO** is a library designed for smart contract development on the Near Blockchain. This SDK provides the necessary tools for creating, deploying, and testing smart contracts written in Go for the Near Blockchain.

To simplify development on NEAR using Go, you can use the [NEAR Go CLI](https://github.com/vlmoon99/near-cli-go), which streamlines common tasks such as building, deploying, and interacting with smart contracts.

---

## 🚨 **Important Prerequisites** 🚨  

For an optimal development experience, use **Go** version **1.23.7** and **TinyGo** version **0.36.0**.

For more information about the compatibility between Go and TinyGo, please refer to the following:  
- [Go Compatibility Matrix](https://tinygo.org/docs/reference/go-compat-matrix/)  
- [Go Language Features Support](https://tinygo.org/docs/reference/lang-support/)  

Please note that smart contracts run in an isolated environment without access to sockets, files, or other APIs. Therefore, you must test them separately in the Near Blockchain environment before deploying them to production.

---

## **Quick Start**  

**TODO:**
1. Add steps for CLI and manual methods
2. Add instructions for running tests, deploying, and calling functions
3. Add links on managing keys, account system, etc.

---

## **1. Create Project**

```bash
near-go create -p "test1" -m "test1" -t "smart-contract-empty"
```

---

## **2. Build Code**

```bash
near-go build
```

---

## **3. Test Code**

```bash
# Run tests inside the current Go package
near-go test package

# Run tests for the entire project
near-go test project
```

---

## **4. Create/Import Account**

```bash
# Create a testnet account with 10N for future test transactions
# Change accountid.testnet to your preferred account ID
near-go account create -n "testnet" -a "accountid.testnet"

# To deploy your smart contract on the mainnet, import your mainnet account
# Create a mainnet wallet using [Near Wallets](https://wallet.near.org/) providers (e.g., Meteor Wallet)
# Then import it using your seed phrase or private key
near-go account import
```

---

## **5. Deploy Contract**

```bash
# Deploy the smart contract on the testnet or mainnet

# Deploy to testnet
near-go deploy -id "accountid.testnet" -n "testnet"

# Deploy to mainnet
near-go deploy -id "accountid.near" -n "mainnet"
```

---

## **6. Call Contract**

After successful deployment, you can call your smart contract functions and test them:

```bash
# Call a method with arguments:
near-go call \
  --from fromaccountid.testnet \
  --to toaccountid.testnet \
  --function WriteData \
  --args '{"key": "testKey", "data": "test1"}' \
  --gas '100 Tgas' \
  --deposit '0 NEAR' \
  --network testnet

# Or call a method without arguments:
near-go call \
  --from fromaccountid.testnet \
  --to toaccountid.testnet \
  --function ReadIncommingTxData \
  --network testnet
```

---

## **7. Managing in Production**

Before deploying to production, review the following resources to ensure your understanding of the NEAR ecosystem:

1. [NEAR Accounts](https://docs.near.org/protocol/account-model)
2. [NEAR Access Keys](https://docs.near.org/protocol/access-keys)
3. [NEAR Storage](https://docs.near.org/protocol/storage/storage-staking)
4. [Integration Tests on Rust](https://github.com/near/near-workspaces-rs)
5. [Integration Tests on JS](https://github.com/near/near-workspaces-js)
6. [Build Your Own Indexer on Rust](https://github.com/near/near-lake-framework-rs)
7. [Build Your Own Indexer on JS](https://github.com/near/near-lake-framework-js)
8. [Build Your Web3 Client (Frontend)](https://github.com/near/wallet-selector)
9. [Near API JS](https://github.com/near/near-api-js)

---

## **Documentation**

For detailed documentation, refer to the `doc` folder. It contains all the essential information required for building smart contracts using Go.

If you're a beginner, it's highly recommended to start with the tutorials on the [Official NEAR Docs](https://docs.near.org), which guide you through step-by-step lessons. These lessons cover:

1. [What is NEAR?](https://docs.near.org/protocol/basics)  
2. [NEAR Accounts](https://docs.near.org/protocol/account-model)  
3. [Transactions](https://docs.near.org/protocol/transactions)  
4. [What is a Smart Contract?](https://docs.near.org/smart-contracts/what-is)  
5. [What are Web3 Apps?](https://docs.near.org/web3-apps/what-is)

---

## **Status Message Tutorial**

For a hands-on experience, go through the [Status_Message_Tutorial](doc/Status_Message_Tutorial.md). This tutorial will teach you how the NEAR-Go CLI works under the hood. You will learn how to:

- Build Go files into WASM smart contracts  
- Sign and send transactions  
- Deploy code to the blockchain

This will give you a deeper understanding of the full development and deployment process on NEAR using Go.
