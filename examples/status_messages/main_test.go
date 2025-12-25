package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

// Initialize the MockSystem globally for the tests
func init() {
	env.SetEnv(system.NewMockSystem())
}

// Helper to setup the test environment and initialize the contract
func setupTest(t *testing.T) *StatusMessage {
	mockSys, ok := env.NearBlockchainImports.(*system.MockSystem)
	if !ok {
		t.Fatal("Environment is not set to MockSystem")
	}
	// Reset storage between tests
	mockSys.Storage = make(map[string][]byte)

	contract := &StatusMessage{}
	contract.Init()

	return contract
}

func TestStatusMessage_Init(t *testing.T) {
	contract := setupTest(t)

	// Since the contract starts empty, checking a random account should return empty string
	// This confirms Init() didn't crash and the map is ready.
	result := contract.GetStatus("any.user")
	expected := ""

	if result != expected {
		t.Errorf("Expected empty string for initialized contract, got '%s'", result)
	}
}

func TestStatusMessage_SetStatus(t *testing.T) {
	contract := setupTest(t)
	mockSys := env.NearBlockchainImports.(*system.MockSystem)

	// Simulate "bob.near" calling the contract
	caller := "bob.near"
	mockSys.PredecessorAccountIdSys = caller

	message := "Hello form Bob!"
	contract.SetStatus(message)

	// Verify the status was saved for Bob
	result := contract.GetStatus(caller)
	if result != message {
		t.Errorf("Expected status '%s', got '%s'", message, result)
	}
}

func TestStatusMessage_GetStatus_NotFound(t *testing.T) {
	// We don't use setupTest here to simulate a manual setup or partial state
	mockSys := env.NearBlockchainImports.(*system.MockSystem)
	mockSys.Storage = make(map[string][]byte)

	contract := &StatusMessage{}
	// Manually initialize the map with prefix "r"
	contract.Records = collections.NewLookupMap[string, string]("r")

	// Check for an account that hasn't set a status
	result := contract.GetStatus("nonexistent.user")
	expected := "" // Your contract logic returns "" on error/not found

	if result != expected {
		t.Errorf("Expected default empty response '%s', got '%s'", expected, result)
	}
}

func TestStatusMessage_MultipleUsers(t *testing.T) {
	contract := setupTest(t)
	mockSys := env.NearBlockchainImports.(*system.MockSystem)

	// 1. Alice sets status
	mockSys.PredecessorAccountIdSys = "alice.near"
	contract.SetStatus("I am Alice")

	// 2. Bob sets status
	mockSys.PredecessorAccountIdSys = "bob.near"
	contract.SetStatus("I am Bob")

	// 3. Verify Alice's data is distinct
	aliceVal := contract.GetStatus("alice.near")
	if aliceVal != "I am Alice" {
		t.Errorf("Alice's data corrupted. Got: %s", aliceVal)
	}

	// 4. Verify Bob's data is distinct
	bobVal := contract.GetStatus("bob.near")
	if bobVal != "I am Bob" {
		t.Errorf("Bob's data corrupted. Got: %s", bobVal)
	}
}

func TestStatusMessage_Persistence_Simulation(t *testing.T) {
	contract := setupTest(t)
	mockSys := env.NearBlockchainImports.(*system.MockSystem)

	// 1. Set state with the original contract instance
	user := "persistent.user"
	message := "Data that should survive"

	mockSys.PredecessorAccountIdSys = user
	contract.SetStatus(message)

	// 2. Simulate a new WASM execution (new struct instance)
	// We do NOT call Init() here, we manually hydrate the map like a loaded contract would.
	newContractInstance := &StatusMessage{}
	newContractInstance.Records = collections.NewLookupMap[string, string]("r")

	// 3. Verify data exists in the new instance (reading from the same MockSystem storage)
	result := newContractInstance.GetStatus(user)
	if result != message {
		t.Errorf("Persistence check failed. Expected '%s', got '%s'", message, result)
	}
}
