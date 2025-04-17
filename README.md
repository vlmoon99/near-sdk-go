# **Near SDK GO**  
[![Telegram](https://img.shields.io/badge/Telegram-join%20chat-blue.svg)](https://t.me/go_near_sdk)  
[![Discord](https://img.shields.io/badge/Discord-join%20chat-blue.svg)](https://discord.gg/UBUPuBm2)  
[![Pkg Go Dev](https://img.shields.io/badge/Pkg%20Go%20Dev-view%20docs-blue.svg)](https://pkg.go.dev/github.com/vlmoon99/near-sdk-go)  

**Near SDK GO** is a library designed for smart contract development on the Near Blockchain. This SDK provides all the necessary tools for creating, deploying, and testing smart contracts written in Go for the Near Blockchain.

To simplify development on NEAR using Go, you can use the [NEAR Go CLI](https://github.com/vlmoon99/near-cli-go). This tool streamlines common tasks such as building, deploying, and interacting with smart contracts written in Go.

### 🚨 **IMPORTANT PREREQUISITES** 🚨  

For an optimal development experience, use **Go** version **1.23.7** and **Tinygo** version **0.36.0**.  

For more information about the compatibility between Go and Tinygo, please refer to the following:  
- [Go Compatibility Matrix](https://tinygo.org/docs/reference/go-compat-matrix/)  
- [Go language features Support](https://tinygo.org/docs/reference/lang-support/)  

Please note that smart contracts run in an isolated environment without access to sockets, files, or other APIs. Therefore, you must test them separately in the Near Blockchain environment before deploying them to production.  

---

## **Quick Start**  
**1.Create project**
```bash
go mod init github.com/vlmoon99/near-sdk-go/examples/quick_intro
```
**2.Add dependencies**
```bash
go get github.com/vlmoon99/near-sdk-go@v0.0.8
```
**3.Add code**

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

**4.Build code**
```bash
tinygo build -size short -no-debug -panic=trap -scheduler=none -gc=leaking -o main.wasm -target wasm-unknown ./ && ls -lh main.wasm
```
**IMPORTANT** - If your build will produce some errors like this one  -
```bash
unsupported parameter type "xyz" 
```
You will need to build it again, this problem solved in CLI which u can use instead of raw tinygo builds

---

## **Documentation**

For detailed documentation, refer to the `doc` folder. It contains all the essential information required for building smart contracts using Go.

If you're a beginner, it's highly recommended to start with the tutorials on the [Official NEAR Docs](https://docs.near.org), which guide you through step-by-step lessons. These lessons cover:

1. [What is NEAR?](https://docs.near.org/protocol/basics)  
2. [NEAR Accounts](https://docs.near.org/protocol/account-model)  
3. [Transactions](https://docs.near.org/protocol/transactions)  
4. [What is a Smart Contract?](https://docs.near.org/smart-contracts/what-is)  
5. [What are Web3 Apps?](https://docs.near.org/web3-apps/what-is)


I highly recommend going through the [Status_Message_Tutorial](doc/Status_Message_Tutorial.md), where you'll explore how the NEAR-Go CLI works under the hood. You'll learn how to:

- Build Go files into WASM smart contracts  
- Sign and send transactions  
- Deploy code to the blockchain

This hands-on experience is key to understanding the full development and deployment process on NEAR using Go.
