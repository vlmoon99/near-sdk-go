package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/types"
)

const greetingKey = "greeting_message"

func getStoredGreeting() string {
	greeting, err := env.StorageRead([]byte(greetingKey))
	if err == nil {
		return string(greeting)
	}
	return "Hello"
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
	message, _ := parser.GetString(greetingKey)
	env.LogString("Updating greeting to: " + message)
	env.StorageWrite([]byte(greetingKey), []byte(message))
	env.ContractValueReturn([]byte("Greeting updated successfully"))
}
