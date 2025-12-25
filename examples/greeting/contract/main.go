package main

import (
	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
)

// @contract:state
type GreetingContract struct {
	Greetings *collections.UnorderedMap[string, string]
}

// @contract:init
func (c *GreetingContract) Init() {
	c.Greetings = collections.NewUnorderedMap[string, string]("g")
	env.LogString("Hello from Init Method")
	c.Greetings.Insert("default", "Hello from NEAR!")
}

// @contract:view
func (c *GreetingContract) GetGreeting() string {

	val, err := c.Greetings.Get("default")
	if err != nil {
		return "Default greeting not found"
	}
	return val
}

// @contract:mutating
func (c *GreetingContract) SetGreeting(greeting string) {
	c.Greetings.Insert("default", greeting)
}
