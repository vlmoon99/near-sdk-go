package main

import (
	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/types"
)

type StatusMessage struct {
	Data *collections.LookupMap
}

func GetState() StatusMessage {
	return StatusMessage{
		Data: collections.NewLookupMap([]byte("b")),
	}
}

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
