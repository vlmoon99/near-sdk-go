package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	systemMock := system.NewMockSystem()
	env.SetEnv(systemMock)
}

func TestSetGreeting(t *testing.T) {
	accountId := "test_account"
	message := "Hello, NEAR!"

	systemMock := env.NearBlockchainImports.(*system.MockSystem)
	systemMock.PredecessorAccountIdSys = accountId
	systemMock.ContractInput = []byte(`{"greeting": "` + message + `"}`)
	SetGreeting()
	contract := GetContract().(*GreetingContract)

	greeting, err := contract.state.greetings.Get("default")
	if err != nil {
		greeting = "Error getting greeting"
	}

	if greeting != message {
		t.Fatalf("Expected message %v, got %v", message, greeting)
	}

}

func TestGetStatus(t *testing.T) {
	message := "Hello, NEAR!"

	GetGreeting()

	contract := GetContract().(*GreetingContract)

	greeting, err := contract.state.greetings.Get("default")
	if err != nil {
		greeting = "Error getting greeting"
	}

	if greeting != message {
		t.Fatalf("Expected message %v, got %v", message, greeting)
	}
}
