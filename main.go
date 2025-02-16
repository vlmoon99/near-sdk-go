package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
)

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
	env.ContractValueReturn([]byte("1"))
}
