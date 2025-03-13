# **Near SDK GO**  
[![Telegram](https://img.shields.io/badge/Telegram-join%20chat-blue.svg)](https://t.me/go_near_sdk)  
[![Discord](https://img.shields.io/badge/Discord-join%20chat-blue.svg)](https://discord.gg/UBUPuBm2)  
[![Pkg Go Dev](https://img.shields.io/badge/Pkg%20Go%20Dev-view%20docs-blue.svg)](https://pkg.go.dev/github.com/vlmoon99/near-sdk-go)  

**Near SDK GO** is a library designed for smart contract development on the Near Blockchain. This SDK provides all the necessary tools for creating, deploying, and testing smart contracts written in Go for the Near Blockchain.


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

If you're a beginner, it's highly recommended to start with the tutorials in `doc/beginners/`, which take you through three step-by-step lessons. These lessons cover:  
[1. Core concepts of the Blockchain technology and Near Blockchain.](doc/beginners/1.Blockchain.md)  
[2. Smart contract basics.](doc/beginners/2.SmartContract.md)  
[3. Client(Frontend) WEB 3 basics.](doc/beginners/3.Client(Frontend).md)  

Once you've mastered these foundational concepts, proceed to [Full-Stack App Tutorial](doc/Full-Stack_App_Tutorial.md). This tutorial guides you through developing and deploying a full-stack application, demonstrating end-to-end development on the Near Blockchain.

---
