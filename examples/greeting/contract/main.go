package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/types"
)

const GreetingKey = "greeting_message"

const DefaultGreeting = "greeting_message"

func getStoredGreeting() string {
	greeting, err := env.StorageRead([]byte(GreetingKey))
	if err == nil {
		return string(greeting)
	}
	return DefaultGreeting
}

//go:export GetGreeting
func GetGreeting() {
	greeting := getStoredGreeting()
	env.LogString(greeting)
	env.ContractValueReturn([]byte(greeting))
}

//go:export SetGreeting
func SetGreeting() {
	options := types.ContractInputOptions{IsRawBytes: false}
	contractInput, _, _ := env.ContractInput(options)
	parser := json.NewParser(contractInput)
	message, _ := parser.GetString(GreetingKey)
	env.LogString("Updating greeting to: " + message)
	env.StorageWrite([]byte(GreetingKey), []byte(message))
	env.ContractValueReturn([]byte("Greeting updated successfully"))
}
