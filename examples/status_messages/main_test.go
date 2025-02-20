package main

// import (
// 	"testing"

// 	"github.com/vlmoon99/near-sdk-go/env"
// 	"github.com/vlmoon99/near-sdk-go/system"
// )

// func init() {
// 	systemMock := system.NewMockSystem()
// 	env.SetEnv(systemMock)
// }

// func TestSetStatus(t *testing.T) {
// 	accountId := "test_account"
// 	message := "Hello, NEAR!"

// 	// Mock the environment
// 	systemMock := env.NearBlockchainImports.(*system.MockSystem)
// 	systemMock.SetPredecessorAccountID(accountId)
// 	systemMock.SetContractInput([]byte(`{"message": "` + message + `"}`))

// 	// Call the SetStatus function
// 	SetStatus()

// 	// Retrieve the stored status message
// 	state := GetState()
// 	storedMessageInterface, err := state.Data.Get([]byte(accountId))
// 	if err != nil {
// 		t.Fatalf("Failed to get stored message: %v", err)
// 	}

// 	storedMessage, ok := storedMessageInterface.(string)
// 	if !ok {
// 		t.Fatalf("Stored message is not a string")
// 	}

// 	if string(storedMessage) != message {
// 		t.Fatalf("Expected message %v, got %v", message, string(storedMessage))
// 	}
// }

// func TestGetStatus(t *testing.T) {
// 	accountId := "test_account"
// 	message := "Hello, NEAR!"

// 	// Mock the environment and set the initial state
// 	systemMock := env.NearBlockchainImports.(*system.MockSystem)
// 	systemMock.SetPredecessorAccountID(accountId)
// 	state := GetState()
// 	state.Data.Insert([]byte(accountId), []byte(message))

// 	// Mock contract input
// 	systemMock.SetContractInput([]byte(`{"account_id": "` + accountId + `"}`))

// 	// Call the GetStatus function
// 	GetStatus()

// 	// Retrieve the stored status message
// 	storedMessageInterface, err := state.Data.Get([]byte(accountId))
// 	if err != nil {
// 		t.Fatalf("Failed to get stored message: %v", err)
// 	}

// 	storedMessage, ok := storedMessageInterface.(string)
// 	if !ok {
// 		t.Fatalf("Stored message is not a string")
// 	}

// 	if string(storedMessage) != message {
// 		t.Fatalf("Expected message %v, got %v", message, string(storedMessage))
// 	}
// }
