//WORKED GET STATE , SET STATE
package main

import (
	contractBuilder "github.com/vlmoon99/near-sdk-go/contract"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/borsh"
	"strconv"
)

// @contract:state
type ContractState struct {
	Count int
	Owner string
}

// @contract:public
// @contract:view
func (c *ContractState) GetCount() int {
	return c.Count
}

// @contract:public
// @contract:mutating
func (c *ContractState) Increment(amount int) int {
	c.Count += amount
	return c.Count
}

// @contract:public
// @contract:mutating
func (c *ContractState) Decrement(amount int) int {
	c.Count -= amount
	return c.Count
}

// @contract:public
// @contract:mutating
// @contract:payable min_deposit=0.001
func (c *ContractState) Reset() string {
	c.Count = 0
	return "Counter reset"
}


//go:export init
func init() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		
		newState := ContractState{1,"vlmoon.near"}

		val, err := borsh.Serialize(newState)
		if err != nil {
			env.LogString("Serialization failed")
		}
		size := len(val)

		env.LogString("Serialized len --> " + strconv.Itoa(size))

		env.StateWrite(val);

		val,err = env.StateRead();
		
		if err != nil {
			env.LogString("StateRead failed")
		}

		var deserialized ContractState
		err = borsh.Deserialize(val, &deserialized)
		if err != nil {
			env.LogString("Deserialization failed")
		}

		val, err = borsh.Serialize(deserialized)
		if err != nil {
			env.LogString("Serialization failed")
		}

		size = len(val)

		env.LogString("Deserialized len --> " + strconv.Itoa(size))

		contractBuilder.ReturnValue("OK");
		
		return nil
	})
}

//go:export getState
func getState() {
		contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		
		val,err := env.StateRead();

		var deserialized ContractState
		err = borsh.Deserialize(val, &deserialized)
		if err != nil {
			env.LogString("Deserialization failed")
		}

		val, err = borsh.Serialize(deserialized)
		if err != nil {
			env.LogString("Serialization failed")
		}

		size := len(val)

		env.LogString("Deserialized len --> " + strconv.Itoa(size))

		contractBuilder.ReturnValue("Count --> " + strconv.Itoa(deserialized.GetCount()));
		
		return nil
	})
}

//go:export setState
func setState() {
		contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		
		val,err := env.StateRead();

		var deserialized ContractState
		err = borsh.Deserialize(val, &deserialized)
		if err != nil {
			env.LogString("Deserialization failed")
		}

		deserialized.Increment(1);

		val, err = borsh.Serialize(deserialized)
		if err != nil {
			env.LogString("Serialization failed")
		}

		size := len(val)

		env.LogString("Deserialized len --> " + strconv.Itoa(size))

		err = env.StateWrite(val);
		if err != nil {
			env.LogString("StateWrite failed")
		}

		contractBuilder.ReturnValue("Count --> " + strconv.Itoa(deserialized.GetCount()));
		
		return nil
	})
}
