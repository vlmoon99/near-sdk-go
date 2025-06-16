package main

import (
	"sync"

	"github.com/vlmoon99/near-sdk-go/collections"
	contractBuilder "github.com/vlmoon99/near-sdk-go/contract"
)

var (
	contractInstance interface{}
	contractOnce     sync.Once
)

type GreetingState struct {
	greetings *collections.UnorderedMap[string, string]
}

func NewGreetingState() *GreetingState {
	return &GreetingState{
		greetings: collections.NewUnorderedMap[string, string]("greetings"),
	}
}

type GreetingContract struct {
	state *GreetingState
}

func NewGreetingContract() *GreetingContract {
	return &GreetingContract{
		state: NewGreetingState(),
	}
}

func GetContract() interface{} {
	contractOnce.Do(func() {
		if contractInstance == nil {
			contractInstance = NewGreetingContract()
		}
	})
	return contractInstance
}

//go:export get_greeting
func GetGreeting() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {

		contract := GetContract().(*GreetingContract)

		greeting, err := contract.state.greetings.Get("default")
		if err != nil {
			greeting = "Error getting greeting"
		}

		contractBuilder.ReturnValue(greeting)
		return nil
	})
}

//go:export set_greeting
func SetGreeting() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		greeting, err := input.JSON.GetString("greeting")
		if err != nil {
			return err
		}

		contract := GetContract().(*GreetingContract)

		if err := contract.state.greetings.Insert("default", greeting); err != nil {
			return err
		}

		contractBuilder.ReturnValue("Greeting updated successfully")
		return nil
	})
}
