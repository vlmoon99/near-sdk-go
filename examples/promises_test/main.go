package main

import (
	"github.com/vlmoon99/near-sdk-go/sdk"
)

//go:export InitContract
func InitContract() {
	sdk.LogString("Init Smart Contract")
}
