// Package main is used as an internal smart contract for testing features during development.
// All completed features must be transferred to the examples/* directory and be structured as separate smart contracts.
package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
)

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
	env.ContractValueReturn([]byte("1"))
}
