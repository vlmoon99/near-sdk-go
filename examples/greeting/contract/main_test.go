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
	systemMock.SetPredecessorAccountID(accountId)
	systemMock.SetContractInput([]byte(`{"greeting_message": "` + message + `"}`))
	SetGreeting()
	storedGreeting := getStoredGreeting()
	if storedGreeting != message {
		t.Fatalf("Expected message %v, got %v", message, storedGreeting)
	}

}

func TestGetStatus(t *testing.T) {
	message := "Hello, NEAR!"

	GetGreeting()

	storedGreeting := getStoredGreeting()
	if storedGreeting != message {
		t.Fatalf("Expected message %v, got %v", message, storedGreeting)
	}
}
