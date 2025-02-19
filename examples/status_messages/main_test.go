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
// 	env.SetPredecessorAccountID(accountId)
// 	env.SetContractInput([]byte(`{"message": "` + message + `"}`))

// 	// Call the SetStatus function
// 	SetStatus()

// 	// Retrieve the stored status message
// 	state := GetState()
// 	storedMessage, err := state.Data.Get([]byte(accountId))
// 	if err != nil {
// 		t.Fatalf("Failed to get stored message: %v", err)
// 	}

// 	if storedMessage != message {
// 		t.Fatalf("Expected message %v, got %v", message, storedMessage)
// 	}
// }

// func TestGetStatus(t *testing.T) {
// 	accountId := "test_account"
// 	message := "Hello, NEAR!"

// 	// Mock the environment and set the initial state
// 	env.SetPredecessorAccountID(accountId)
// 	state := GetState()
// 	state.Data.Insert([]byte(accountId), message)

// 	// Mock contract input
// 	env.SetContractInput([]byte(`{"account_id": "` + accountId + `"}`))

// 	// Call the GetStatus function
// 	GetStatus()

// 	// Retrieve the returned value
// 	returnedValue := env.GetContractValueReturn()
// 	if string(returnedValue) != message {
// 		t.Fatalf("Expected return value %v, got %v", message, string(returnedValue))
// 	}
// }
