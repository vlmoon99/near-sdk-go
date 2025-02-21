package main

import (
	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/types"
)

// This type represents the internal state which we are using in our smart contract.
// We can have stored state on the blockchain or we can simply store it inside the wasm contract in memory.
// In this case, I store it inside memory and do not save my state inside the blockchain, because we have a very simple structure and type.
type StatusMessage struct {
	// Represents a proxy collection to the original methods inside /near-sdk-go/env/env.go file such as StorageWrite, StorageRead, StorageRemove, StorageHasKey.
	// Before using any collections or top-level abstractions, it is highly recommended to learn how env methods work.
	Data *collections.LookupMap
}

func GetState() StatusMessage {
	return StatusMessage{
		// []byte("b") - represents a prefix which will be added for each key inside this collection. So if I put a key with the name "test", in the blockchain I will have (b + test) as the key,
		// but the value remains the same.
		Data: collections.NewLookupMap([]byte("b")),
	}
}

// //go:export - This is a commentary which we need to use in order to export our functions to the smart contract clients. If we do not mark our methods as //go:export, we cannot call them after deployment.
// If we mark our methods with this commentary, it will be exported in our wasm file and will be visible to our clients.
// Exported functions cannot have any input and output parameters. For input from the user side, we need to use the env.ContractInput method to receive user input.
// For output, we need to use the env.ContractValueReturn function in order to provide the return value to the user.

//go:export SetStatus
func SetStatus() {
	accountId, _ := env.GetPredecessorAccountID()
	options := types.ContractInputOptions{IsRawBytes: false}
	contractInput, _, _ := env.ContractInput(options)
	parser := json.NewParser(contractInput)
	message, _ := parser.GetString("message")
	state := GetState()
	state.Data.Insert([]byte(accountId), string(message))
	// env.LogString("Message : " + message + " was insterted")
	env.ContractValueReturn([]byte(message))
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
	// env.LogString("Status : " + status + " on account id : " + accountId)
	env.ContractValueReturn([]byte(status))
}
