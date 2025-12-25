package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func setupTest(t *testing.T) *GreetingContract {
	mockSys, ok := env.NearBlockchainImports.(*system.MockSystem)
	if !ok {
		t.Fatal("Environment is not set to MockSystem")
	}
	mockSys.Storage = make(map[string][]byte)

	contract := &GreetingContract{}

	contract.Init()

	return contract
}

func TestGreeting_Init_Default(t *testing.T) {
	contract := setupTest(t)

	expected := "Hello from NEAR!"
	result := contract.GetGreeting()

	if result != expected {
		t.Errorf("Expected default greeting '%s', got '%s'", expected, result)
	}
}

func TestGreeting_SetGreeting(t *testing.T) {
	contract := setupTest(t)

	if contract.GetGreeting() != "Hello from NEAR!" {
		t.Fatal("Initial state incorrect")
	}

	newGreeting := "Welcome to Web3 Go!"
	contract.SetGreeting(newGreeting)

	result := contract.GetGreeting()
	if result != newGreeting {
		t.Errorf("Expected new greeting '%s', got '%s'", newGreeting, result)
	}
}

func TestGreeting_GetGreeting_NotFound(t *testing.T) {
	mockSys := env.NearBlockchainImports.(*system.MockSystem)
	mockSys.Storage = make(map[string][]byte)

	contract := &GreetingContract{}
	contract.Greetings = collections.NewUnorderedMap[string, string]("g")

	result := contract.GetGreeting()
	expected := "Default greeting not found"

	if result != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, result)
	}
}

func TestGreeting_Persistence_Simulation(t *testing.T) {
	contract := setupTest(t)

	savedMessage := "I persist across execution"
	contract.SetGreeting(savedMessage)

	newContractInstance := &GreetingContract{}
	newContractInstance.Greetings = collections.NewUnorderedMap[string, string]("g")

	result := newContractInstance.GetGreeting()
	if result != savedMessage {
		t.Errorf("Persistence check failed. Expected '%s', got '%s'", savedMessage, result)
	}
}
